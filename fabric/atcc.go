package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ===================== 数据结构定义 =====================
// 注意：所有结构体添加 docType 字段用于CouchDB查询分类

// Evaluation 测评记录结构
type Evaluation struct {
	DocType      string `json:"docType"`      // 文档类型标识
	EvaluationID string `json:"Evaluation_ID"` // 测评唯一ID
	UserID       string `json:"User_ID"`      // 关联用户ID
	PointsDegree string `json:"Points_Degree"`// 评分等级
	Feedback     string `json:"Feedback"`     // 详细反馈
}

// TestResult 测试结果结构
type TestResult struct {
	DocType     string `json:"docType"`      // 文档类型标识
	TestID      string `json:"Test_ID"`      // 测试唯一ID
	UserID      string `json:"User_ID"`      // 关联用户ID
	ScoreSum    string `json:"Score_Sum"`    // 总分
	PaperNumber string `json:"Paper_Number"` // 试卷编号
	Answer      string `json:"Answer"`       // 答案内容
}

// Judgement 评价记录结构
type Judgement struct {
	DocType            string `json:"docType"`              // 文档类型标识
	JudgementID        string `json:"Judgement_ID"`         // 评价唯一ID
	UserID             string `json:"User_ID"`              // 关联用户ID
	JudgementObjection string `json:"Judgement_Objection"`  // 异议内容
	JudgementObjectID  string `json:"Judgement_ObjectID"`   // 关联对象ID
	JudgementRating    string `json:"Judgement_Rating"`     // 评分等级
	JudgementContent   string `json:"Judgement_Content"`    // 评价内容
	JudgementTime      string `json:"Judgement_Time"`       // 评价时间
}

// ===================== 智能合约结构 =====================
type SmartContract struct {
	contractapi.Contract
}

// ===================== 测评记录管理 =====================

// UploadEvaluation 上传测评记录
// 参数：测评记录JSON字符串
// 返回值：错误信息
func (s *SmartContract) UploadEvaluation(ctx contractapi.TransactionContextInterface, evaluationJSON string) error {
	var evaluation Evaluation
	// 解析输入数据
	if err := json.Unmarshal([]byte(evaluationJSON), &evaluation); err != nil {
		return fmt.Errorf("解析测评记录失败: %v", err)
	}
	
	// 数据校验
	if evaluation.EvaluationID == "" || evaluation.UserID == "" {
		return fmt.Errorf("缺少必要字段（EvaluationID/UserID）")
	}
	
	// 设置文档类型
	evaluation.DocType = "Evaluation"
	compositeKey := fmt.Sprintf("Evaluation-%s", evaluation.EvaluationID)
	
	// 检查重复记录
	existing, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return fmt.Errorf("状态数据库查询失败: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("测评记录 %s 已存在", evaluation.EvaluationID)
	}
	
	// 存储数据
	data, err := json.Marshal(evaluation)
	if err != nil {
		return fmt.Errorf("数据序列化失败: %v", err)
	}
	return ctx.GetStub().PutState(compositeKey, data)
}

// ModifyEvaluation 修改测评记录
// 参数：测评ID，新测评记录JSON
// 返回值：错误信息
func (s *SmartContract) ModifyEvaluation(ctx contractapi.TransactionContextInterface, evaluationID string, newEvaluationJSON string) error {
	// 获取原记录
	compositeKey := fmt.Sprintf("Evaluation-%s", evaluationID)
	existingData, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return fmt.Errorf("状态查询失败: %v", err)
	}
	if existingData == nil {
		return fmt.Errorf("找不到指定测评记录")
	}
	
	// 解析新数据
	var newEval Evaluation
	if err := json.Unmarshal([]byte(newEvaluationJSON), &newEval); err != nil {
		return fmt.Errorf("解析新记录失败: %v", err)
	}
	
	// ID一致性检查
	if newEval.EvaluationID != evaluationID {
		return fmt.Errorf("禁止修改测评ID")
	}
	
	// 保留原始文档类型
	newEval.DocType = "Evaluation"
	
	// 存储更新
	data, err := json.Marshal(newEval)
	if err != nil {
		return fmt.Errorf("数据序列化失败: %v", err)
	}
	return ctx.GetStub().PutState(compositeKey, data)
}

// GetEvaluationByID 根据ID获取测评记录
// 参数：测评ID，用户ID（用于权限验证）
// 返回值：测评记录指针，错误信息
func (s *SmartContract) GetEvaluationByID(ctx contractapi.TransactionContextInterface, evaluationID string, userID string) (*Evaluation, error) {
	if evaluationID == "" || userID == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	
	compositeKey := fmt.Sprintf("Evaluation-%s", evaluationID)
	data, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("状态数据库查询失败: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("找不到指定测评记录")
	}
	
	var evaluation Evaluation
	if err := json.Unmarshal(data, &evaluation); err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}
	
	// 权限验证
	if evaluation.UserID != userID {
		return nil, fmt.Errorf("无权访问该记录")
	}
	return &evaluation, nil
}

// GetEvaluationByUser 根据用户ID获取所有测评记录
// 参数：用户ID
// 返回值：测评记录切片，错误信息
func (s *SmartContract) GetEvaluationByUser(ctx contractapi.TransactionContextInterface, userID string) ([]*Evaluation, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 构建安全查询（只返回属于该用户的记录）
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"docType": "Evaluation",
			"User_ID": userID,
		},
	}
	queryBytes, _ := json.Marshal(query)

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("查询执行失败: %v", err)
	}
	defer resultsIterator.Close()

	var evaluations []*Evaluation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("结果迭代失败: %v", err)
		}

		var eval Evaluation
		if err := json.Unmarshal(queryResponse.Value, &eval); err != nil {
			return nil, fmt.Errorf("数据解析失败: %v", err)
		}
		evaluations = append(evaluations, &eval)
	}
	return evaluations, nil
}

// GetAllEvaluations 获取所有测评记录（谨慎使用，大数据量时需要分页）
// 参数：无
// 返回值：全部测评记录切片，错误信息
func (s *SmartContract) GetAllEvaluations(ctx contractapi.TransactionContextInterface) ([]*Evaluation, error) {
	// 构建类型查询
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"docType": "Evaluation",
		},
	}
	queryBytes, _ := json.Marshal(query)

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("查询执行失败: %v", err)
	}
	defer resultsIterator.Close()

	var evaluations []*Evaluation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("结果迭代失败: %v", err)
		}

		var eval Evaluation
		if err := json.Unmarshal(queryResponse.Value, &eval); err != nil {
			return nil, fmt.Errorf("数据解析失败: %v", err)
		}
		evaluations = append(evaluations, &eval)
	}
	return evaluations, nil
}

// ===================== 测试结果管理 =====================

// UploadTestResult 上传测试结果
// 参数：测试结果JSON字符串
// 返回值：错误信息
func (s *SmartContract) UploadTestResult(ctx contractapi.TransactionContextInterface, testJSON string) error {
	var testResult TestResult
	if err := json.Unmarshal([]byte(testJSON), &testResult); err != nil {
		return fmt.Errorf("解析测试结果失败: %v", err)
	}
	
	// 数据校验
	if testResult.TestID == "" || testResult.UserID == "" {
		return fmt.Errorf("缺少必要字段（TestID/UserID）")
	}
	
	// 设置文档类型
	testResult.DocType = "TestResult"
	compositeKey := fmt.Sprintf("TestResult-%s", testResult.TestID)
	
	// 检查重复记录
	existing, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return fmt.Errorf("状态数据库查询失败: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("测试结果 %s 已存在", testResult.TestID)
	}
	
	// 存储数据
	data, err := json.Marshal(testResult)
	if err != nil {
		return fmt.Errorf("数据序列化失败: %v", err)
	}
	return ctx.GetStub().PutState(compositeKey, data)
}

// GetTestResultsByUser 获取用户所有测试结果
// 参数：用户ID
// 返回值：测试结果切片，错误信息
func (s *SmartContract) GetTestResultsByUser(ctx contractapi.TransactionContextInterface, userID string) ([]*TestResult, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	
	// 构建安全查询
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"docType": "TestResult",
			"User_ID": userID,
		},
	}
	queryBytes, _ := json.Marshal(query)
	
	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("查询执行失败: %v", err)
	}
	defer resultsIterator.Close()
	
	var results []*TestResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("结果迭代失败: %v", err)
		}
		
		var test TestResult
		if err := json.Unmarshal(queryResponse.Value, &test); err != nil {
			return nil, fmt.Errorf("数据解析失败: %v", err)
		}
		results = append(results, &test)
	}
	return results, nil
}


// GetTestResultsByTestID 根据测试ID获取测试结果
// 参数：测试ID
// 返回值：测试结果指针，错误信息
// 注意：此方法不验证用户权限，需根据业务需求决定是否公开访问
func (s *SmartContract) GetTestResultsByTestID(ctx contractapi.TransactionContextInterface, testID string) (*TestResult, error) {
	if testID == "" {
		return nil, fmt.Errorf("测试ID不能为空")
	}

	compositeKey := fmt.Sprintf("TestResult-%s", testID)
	data, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("状态数据库查询失败: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("找不到指定测试结果")
	}

	var testResult TestResult
	if err := json.Unmarshal(data, &testResult); err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}
	return &testResult, nil
}

// GetTestResultsByID 根据用户ID和测试ID联合查询
// 参数：用户ID（用于权限验证），测试ID
// 返回值：测试结果指针，错误信息
func (s *SmartContract) GetTestResultsByID(ctx contractapi.TransactionContextInterface, userID string, testID string) (*TestResult, error) {
	if userID == "" || testID == "" {
		return nil, fmt.Errorf("参数不能为空")
	}

	// 先通过测试ID获取记录
	testResult, err := s.GetTestResultsByTestID(ctx, testID)
	if err != nil {
		return nil, err
	}

	// 权限验证
	if testResult.UserID != userID {
		return nil, fmt.Errorf("无权访问该测试记录")
	}
	return testResult, nil
}

// ===================== 评价记录管理 =====================

// UploadJudgement 上传评价记录
// 参数：评价记录JSON字符串
// 返回值：错误信息
func (s *SmartContract) UploadJudgement(ctx contractapi.TransactionContextInterface, judgementJSON string) error {
	var judgement Judgement
	if err := json.Unmarshal([]byte(judgementJSON), &judgement); err != nil {
		return fmt.Errorf("解析评价记录失败: %v", err)
	}
	
	// 数据校验
	if judgement.JudgementID == "" || judgement.UserID == "" {
		return fmt.Errorf("缺少必要字段（JudgementID/UserID）")
	}
	
	// 设置文档类型
	judgement.DocType = "Judgement"
	compositeKey := fmt.Sprintf("Judgement-%s", judgement.JudgementID)
	
	// 存储数据
	data, err := json.Marshal(judgement)
	if err != nil {
		return fmt.Errorf("数据序列化失败: %v", err)
	}
	return ctx.GetStub().PutState(compositeKey, data)
}

// GetJudgementByUser 获取用户所有评价记录
// 参数：用户ID
// 返回值：评价记录切片，错误信息
func (s *SmartContract) GetJudgementByUser(ctx contractapi.TransactionContextInterface, userID string) ([]*Judgement, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 构建安全查询
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"docType": "Judgement",
			"User_ID": userID,
		},
	}
	queryBytes, _ := json.Marshal(query)

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("查询执行失败: %v", err)
	}
	defer resultsIterator.Close()

	var judgements []*Judgement
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("结果迭代失败: %v", err)
		}

		var judgement Judgement
		if err := json.Unmarshal(queryResponse.Value, &judgement); err != nil {
			return nil, fmt.Errorf("数据解析失败: %v", err)
		}
		judgements = append(judgements, &judgement)
	}
	return judgements, nil
}

// GetJudgementByID 根据用户ID和评价ID联合查询
// 参数：用户ID（用于权限验证），评价ID
// 返回值：评价记录指针，错误信息
func (s *SmartContract) GetJudgementByID(ctx contractapi.TransactionContextInterface, userID string, judgementID string) (*Judgement, error) {
	if userID == "" || judgementID == "" {
		return nil, fmt.Errorf("参数不能为空")
	}

	// 先获取评价记录
	compositeKey := fmt.Sprintf("Judgement-%s", judgementID)
	data, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("状态数据库查询失败: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("找不到指定评价记录")
	}

	var judgement Judgement
	if err := json.Unmarshal(data, &judgement); err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}

	// 权限验证
	if judgement.UserID != userID {
		return nil, fmt.Errorf("无权访问该评价记录")
	}
	return &judgement, nil
}


// GetJudgementByJudgementID 根据评价ID查询（无权限验证）
// 参数：评价ID
// 返回值：评价记录指针，错误信息
// 注意：此方法不验证用户权限，需根据业务需求决定是否公开访问
func (s *SmartContract) GetJudgementByJudgementID(ctx contractapi.TransactionContextInterface, judgementID string) (*Judgement, error) {
	if judgementID == "" {
		return nil, fmt.Errorf("评价ID不能为空")
	}

	compositeKey := fmt.Sprintf("Judgement-%s", judgementID)
	data, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("状态数据库查询失败: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("找不到指定评价记录")
	}

	var judgement Judgement
	if err := json.Unmarshal(data, &judgement); err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}
	return &judgement, nil
}


// ===================== 通用功能 =====================

// DeleteRecord 通用删除方法
// 参数：记录类型（Evaluation/TestResult/Judgement），记录ID
// 返回值：错误信息
func (s *SmartContract) DeleteRecord(ctx contractapi.TransactionContextInterface, recordType string, recordID string) error {
	if recordID == "" {
		return fmt.Errorf("记录ID不能为空")
	}
	
	var compositeKey string
	switch recordType {
	case "Evaluation":
		compositeKey = fmt.Sprintf("Evaluation-%s", recordID)
	case "TestResult":
		compositeKey = fmt.Sprintf("TestResult-%s", recordID)
	case "Judgement":
		compositeKey = fmt.Sprintf("Judgement-%s", recordID)
	default:
		return fmt.Errorf("不支持的记录类型")
	}
	
	return ctx.GetStub().DelState(compositeKey)
}

// ===================== 初始化方法 =====================

// InitLedger 初始化示例数据（仅限开发环境使用）
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 示例测评记录
	evaluation := Evaluation{
		DocType:      "Evaluation",
		EvaluationID: "eval_001",
		UserID:       "user_001",
		PointsDegree: "A",
		Feedback:     "Excellent performance in all aspects",
	}
	evalKey := fmt.Sprintf("Evaluation-%s", evaluation.EvaluationID)
	if data, err := json.Marshal(evaluation); err == nil {
		ctx.GetStub().PutState(evalKey, data)
	}
	
	// 示例测试结果
	testResult := TestResult{
		DocType:     "TestResult",
		TestID:      "test_001",
		UserID:      "user_001",
		ScoreSum:    "98",
		PaperNumber: "2023-FINAL-01",
		Answer:      "Correct answers for all questions",
	}
	testKey := fmt.Sprintf("TestResult-%s", testResult.TestID)
	if data, err := json.Marshal(testResult); err == nil {
		ctx.GetStub().PutState(testKey, data)
	}
	
	// 示例评价记录
	judgement := Judgement{
		DocType:            "Judgement",
		JudgementID:        "judge_001",
		UserID:             "user_002",
		JudgementObjection: "None",
		JudgementObjectID:  "eval_001",
		JudgementRating:    "5",
		JudgementContent:   "Very fair evaluation",
		JudgementTime:      time.Now().Format(time.RFC3339),
	}
	judgeKey := fmt.Sprintf("Judgement-%s", judgement.JudgementID)
	if data, err := json.Marshal(judgement); err == nil {
		ctx.GetStub().PutState(judgeKey, data)
	}
	
	return nil
}

// ===================== 主函数 =====================
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("链码初始化失败: %v", err)
		return
	}
	
	if err := chaincode.Start(); err != nil {
		fmt.Printf("链码服务启动失败: %v", err)
	}
}
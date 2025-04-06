package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"
	"strings"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/Users/zhangyibin/Fabric/hyperledger-fabric/fabric/scripts/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem" // 完整证书路径
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	channelName  = "mychannel"
	chaincodeID  = "atcc"
)

// 数据结构定义（必须与链码中的结构匹配）
type Evaluation struct {
	DocType      string `json:"docType"`
	EvaluationID string `json:"Evaluation_ID"`
	UserID       string `json:"User_ID"`
	PointsDegree string `json:"Points_Degree"`
	Feedback     string `json:"Feedback"`
}

type TestResult struct {
	DocType     string `json:"docType"`
	TestID      string `json:"Test_ID"`
	UserID      string `json:"User_ID"`
	ScoreSum    string `json:"Score_Sum"`
	PaperNumber string `json:"Paper_Number"`
	Answer      string `json:"Answer"`
}

type Judgement struct {
	DocType            string `json:"docType"`
	JudgementID        string `json:"Judgement_ID"`
	UserID             string `json:"User_ID"`
	JudgementObjection string `json:"Judgement_Objection"`
	JudgementObjectID  string `json:"Judgement_ObjectID"`
	JudgementRating    string `json:"Judgement_Rating"`
	JudgementContent   string `json:"Judgement_Content"`
	JudgementTime      string `json:"Judgement_Time"`
}

type Client struct {
	contract *client.Contract
}

func NewClient() (*Client, error) {
	// 创建gRPC客户端连接
	connection, err := newGrpcConnection()
	if err != nil {
		return nil, fmt.Errorf("创建gRPC连接失败: %v", err)
	}

	// 创建Gateway客户端
	id, err := newIdentity()
	if err != nil {
		return nil, err
	}

	sign, err := newSign()
	if err != nil {
		return nil, err
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(connection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("连接网关失败: %v", err)
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeID)

	return &Client{contract: contract}, nil
}

// ===================== 测评记录操作 =====================
func (c *Client) UploadEvaluation(evaluation Evaluation) (string, error) {
	evalJSON, err := json.Marshal(evaluation)
	if err != nil {
		return "", fmt.Errorf("序列化测评记录失败: %v", err)
	}

	result, err := c.contract.SubmitTransaction("UploadEvaluation", string(evalJSON))
	if err != nil {
		return "", fmt.Errorf("提交交易失败: %v", err)
	}
	return string(result), nil
}

func (c *Client) ModifyEvaluation(evaluationID string, newEvaluation Evaluation) error {
	newEvalJSON, err := json.Marshal(newEvaluation)
	if err != nil {
		return fmt.Errorf("序列化新测评记录失败: %v", err)
	}

	_, err = c.contract.SubmitTransaction("ModifyEvaluation", evaluationID, string(newEvalJSON))
	return err
}

func (c *Client) GetEvaluationByID(evaluationID, userID string) (*Evaluation, error) {
	result, err := c.contract.EvaluateTransaction("GetEvaluationByID", evaluationID, userID)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var evaluation Evaluation
	if err := json.Unmarshal(result, &evaluation); err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}
	return &evaluation, nil
}

func (c *Client) GetEvaluationByUser(userID string) ([]Evaluation, error) {
	result, err := c.contract.EvaluateTransaction("GetEvaluationByUser", userID)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var evaluations []Evaluation
	if err := json.Unmarshal(result, &evaluations); err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}
	return evaluations, nil
}

// ===================== 测试结果操作 =====================
func (c *Client) UploadTestResult(test TestResult) (string, error) {
	testJSON, err := json.Marshal(test)
	if err != nil {
		return "", fmt.Errorf("序列化测试结果失败: %v", err)
	}

	result, err := c.contract.SubmitTransaction("UploadTestResult", string(testJSON))
	if err != nil {
		return "", fmt.Errorf("提交交易失败: %v", err)
	}
	return string(result), nil
}

func (c *Client) GetTestResultsByUser(userID string) ([]TestResult, error) {
	result, err := c.contract.EvaluateTransaction("GetTestResultsByUser", userID)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var tests []TestResult
	if err := json.Unmarshal(result, &tests); err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}
	return tests, nil
}

func (c *Client) GetTestResultsByID(userID, testID string) (*TestResult, error) {
	result, err := c.contract.EvaluateTransaction("GetTestResultsByID", userID, testID)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var test TestResult
	if err := json.Unmarshal(result, &test); err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}
	return &test, nil
}

// ===================== 评价记录操作 =====================
func (c *Client) UploadJudgement(judgement Judgement) (string, error) {
	judgeJSON, err := json.Marshal(judgement)
	if err != nil {
		return "", fmt.Errorf("序列化评价记录失败: %v", err)
	}

	result, err := c.contract.SubmitTransaction("UploadJudgement", string(judgeJSON))
	if err != nil {
		return "", fmt.Errorf("提交交易失败: %v", err)
	}
	return string(result), nil
}

func (c *Client) GetJudgementByUser(userID string) ([]Judgement, error) {
	result, err := c.contract.EvaluateTransaction("GetJudgementByUser", userID)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var judgements []Judgement
	if err := json.Unmarshal(result, &judgements); err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}
	return judgements, nil
}

func (c *Client) GetJudgementByID(userID, judgementID string) (*Judgement, error) {
	result, err := c.contract.EvaluateTransaction("GetJudgementByID", userID, judgementID)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	var judgement Judgement
	if err := json.Unmarshal(result, &judgement); err != nil {
		return nil, fmt.Errorf("解析结果失败: %v", err)
	}
	return &judgement, nil
}

// ===================== 通用操作 =====================
func (c *Client) DeleteRecord(recordType, recordID string) error {
	_, err := c.contract.SubmitTransaction("DeleteRecord", recordType, recordID)
	return err
}

// ===================== 连接工具函数 =====================
func newGrpcConnection() (*grpc.ClientConn, error) {
	certBytes, err := ioutil.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("读取TLS证书失败: %v", err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certBytes); !ok {
		return nil, fmt.Errorf("解析TLS证书失败")
	}

	transportCredentials := credentials.NewClientTLSFromCert(certPool, "peer0.org1.example.com")
	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("创建gRPC连接失败: %v", err)
	}

	return connection, nil
}

func newIdentity() (*identity.X509Identity, error) {
	// 直接使用配置的证书路径
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书文件失败: %v", err)
	}

	cert, err := identity.CertificateFromPEM(certBytes)
	if err != nil {
		return nil, err
	}

	return identity.NewX509Identity(mspID, cert)
}

func newSign() (identity.Sign, error) {
	// 自动查找私钥文件（通常以_sk结尾）
	files, err := ioutil.ReadDir(keyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥目录失败: %v", err)
	}

	// 查找第一个以_sk结尾的文件
	var pkFile string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "_sk") {
			pkFile = f.Name()
			break
		}
	}
	if pkFile == "" {
		return nil, fmt.Errorf("未找到私钥文件")
	}

	pkPath := path.Join(keyPath, pkFile)
	pkBytes, err := ioutil.ReadFile(pkPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %v", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(pkBytes)
	if err != nil {
		return nil, err
	}

	return identity.NewPrivateKeySign(privateKey)
}

// ===================== 示例使用 =====================
func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 示例：上传测评记录
	eval := Evaluation{
		EvaluationID: "eval_002",
		UserID:       "user_002",
		PointsDegree: "B+",
		Feedback:     "Good performance with room for improvement",
	}
	if _, err := client.UploadEvaluation(eval); err != nil {
		log.Printf("上传测评记录失败: %v", err)
	}

	// 示例：查询测评记录
	result, err := client.GetEvaluationByID("eval_002", "user_002")
	if err != nil {
		log.Printf("查询测评记录失败: %v", err)
	} else {
		fmt.Printf("查询结果: %+v\n", result)
	}

	// 示例：删除记录
	if err := client.DeleteRecord("Evaluation", "eval_002"); err != nil {
		log.Printf("删除记录失败: %v", err)
	}
}
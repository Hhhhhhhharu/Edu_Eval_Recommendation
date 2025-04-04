<template>
  <div class="class-info-container">
    <el-card class="class-card">
      <!-- 标题 & 退出按钮 -->
      <el-row justify="space-between" align="middle" class="header-row">
        <h2 class="title">班级信息</h2>
        <el-button type="danger" class="logout-btn" @click="logout">退出登录</el-button>
      </el-row>

      <el-divider></el-divider>

      <!-- 选择班级 -->
      <el-form-item label="选择班级">
        <el-select v-model="selectedClass" placeholder="请选择班级" @change="fetchClassData">
          <el-option v-for="classItem in classList" :key="classItem.id" :label="classItem.name" :value="classItem.id" />
        </el-select>
      </el-form-item>

      <el-divider></el-divider>

      <!-- 班级总体评估 -->
      <h3>班级总体评估</h3>
      <el-descriptions title="评估结果" :column="1" border>
        <el-descriptions-item label="平均分">{{ classEvaluation.averageScore }} 分</el-descriptions-item>
        <el-descriptions-item label="总体评估结果">{{ classEvaluation.evaluation_result }}</el-descriptions-item>
      </el-descriptions>

      <el-divider></el-divider>

      <!-- 学生列表 -->
      <h3>学生列表</h3>
      <el-empty v-if="students.length === 0" description="暂无学生信息" />
      <el-table :data="students" border stripe style="width: 100%">
        <el-table-column prop="studentId" label="学号" width="120" />
        <el-table-column prop="name" label="姓名" width="150" />
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button type="primary" @click="viewStudent(row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 学生详情对话框 -->
    <el-dialog v-model="studentDialogVisible" title="学生详情" width="500px">
      <el-descriptions title="基本信息" :column="1" border>
        <el-descriptions-item label="学号">{{ selectedStudent.studentId }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ selectedStudent.name }}</el-descriptions-item>
      </el-descriptions>
      <el-divider></el-divider>
      <h3>已完成科目</h3>
      <el-empty v-if="selectedStudent.completedSubjects.length === 0" description="暂无评测记录" />
      <el-table :data="selectedStudent.completedSubjects" border stripe>
        <el-table-column prop="name" label="科目" width="200" />
        <el-table-column prop="score" label="得分" width="100" />
        <el-table-column prop="evaluation_result" label="评估结果" />
        <el-table-column label="反馈">
          <template #default="{ row }">
            <el-button type="success" @click="openFeedbackDialog(row)">反馈</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="studentDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 反馈对话框 -->
    <el-dialog v-model="feedbackDialogVisible" title="提交反馈" width="400px">
      <el-input v-model="feedbackText" type="textarea" placeholder="请输入反馈内容" />
      <template #footer>
        <el-button @click="feedbackDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitFeedback">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const classList = ref([{ id: 1, name: "区块链211班" }, { id: 2, name: "区块链212班" }, { id: 3, name: "区块链213班" }, { id: 4, name: "区块链214班" }]);
const selectedClass = ref<number | null>(null);
const students = ref<any[]>([]);
const classEvaluation = ref({ averageScore: "0", evaluation_result: "暂无评估结果" });
const studentDialogVisible = ref(false);
const selectedStudent = ref({ studentId: "", name: "", completedSubjects: [] });
const feedbackDialogVisible = ref(false);
const feedbackText = ref("");
const feedbackTarget = ref<any>(null);

const fetchClassData = () => {
  if (!selectedClass.value) return;
  students.value = [
    { studentId: "2021131001", name: "张三", completedSubjects: [{ name: "数据结构", score: 85, evaluation_result: "{ 排序：90%, 树：80%}" }] },
    { studentId: "2021131002", name: "李四", completedSubjects: [{ name: "Python", score: 90, evaluation_result: "{ 排序：95%, 树：85%}" }] },
    { studentId: "2021131003", name: "王五", completedSubjects: [{ name: "操作系统", score: 75, evaluation_result: "{ 排序：70%, 树：80%}" }] }
  ];

  const totalScore = students.value.reduce((sum, student) => sum + (student.completedSubjects[0]?.score || 0), 0);
  classEvaluation.value.averageScore = (totalScore / students.value.length).toFixed(1);

  let totalSort = 0, totalTree = 0, count = 0;
  students.value.forEach(student => {
    const evalResult = student.completedSubjects[0]?.evaluation_result;
    if (evalResult) {
      const sortMatch = evalResult.match(/排序：(\d+)%/);
      const treeMatch = evalResult.match(/树：(\d+)%/);
      if (sortMatch && treeMatch) {
        totalSort += parseInt(sortMatch[1]);
        totalTree += parseInt(treeMatch[1]);
        count++;
      }
    }
  });
  classEvaluation.value.evaluation_result = `{ 排序：${count > 0 ? (totalSort / count).toFixed(1) : "0"}%, 树：${count > 0 ? (totalTree / count).toFixed(1) : "0"}% }`;
};

const viewStudent = (student: any) => {
  selectedStudent.value = student;
  studentDialogVisible.value = true;
};

const openFeedbackDialog = (row: any) => {
  feedbackTarget.value = row;
  feedbackText.value = "";
  feedbackDialogVisible.value = true;
};

const submitFeedback = () => {
  console.log(`反馈：${feedbackText.value}，针对 ${feedbackTarget.value.name}`);
  feedbackDialogVisible.value = false;
};

const logout = () => {
  localStorage.removeItem("userToken");
  router.push("/login");
};

onMounted(() => {
  if (classList.value.length > 0) {
    selectedClass.value = classList.value[0].id;
    fetchClassData();
  }
});
</script>

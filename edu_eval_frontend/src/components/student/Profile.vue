<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <h2>个人信息</h2>
      <el-divider></el-divider>
      <el-descriptions title="用户信息" :column="1" border>
        <el-descriptions-item label="学号">{{ user.studentId }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ user.name }}</el-descriptions-item>
      </el-descriptions>

      <el-divider></el-divider>

      <h3>已完成的科目</h3>
      <el-empty v-if="completedSubjects.length === 0" description="暂无已完成的科目" />

      <el-row gutter="20" v-else>
        <el-col v-for="subject in completedSubjects" :key="subject.id" :span="8">
          <el-button class="subject-button" type="primary" @click="viewEvaluation(subject)">
            {{ subject.name }}
          </el-button>
        </el-col>
      </el-row>

      <el-divider></el-divider>

      <!-- 退出登录按钮 -->
      <el-button type="danger" class="logout-button" @click="logout">退出登录</el-button>
    </el-card>

    <!-- 评测详情对话框 -->
    <el-dialog v-model="evaluationDialogVisible" :title="selectedSubject.name + ' 评估结果'" width="500px">
      <el-descriptions title="评估详情" :column="1" border>
        <el-descriptions-item label="总得分">{{ evaluationResult.score }} 分</el-descriptions-item>
        <el-descriptions-item label="评估结果">{{ evaluationResult.evaluation_result }}</el-descriptions-item>
      </el-descriptions>
      
      <!-- 反馈输入框 -->
      <el-input v-model="evaluationFeedback" type="textarea" placeholder="输入您的反馈" />
      <el-button type="primary" @click="submitFeedback">提交反馈</el-button>
      
      <template #footer>
        <el-button @click="evaluationDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";

const router = useRouter();

// 用户信息
const user = ref({
  studentId: "",
  name: ""
});

// 已完成的科目
const completedSubjects = ref([
  { id: 1, name: "数据结构", score: 85, evaluation_result: "{ 排序：90%, 树：80%}" },
  { id: 2, name: "Python", score: 90, evaluation_result: "{ 函数：95%, 类：85%}" },
]);

// 评测详情
const evaluationDialogVisible = ref(false);
const selectedSubject = ref({ id: 0, name: "" });
const evaluationResult = ref({ score: 0, evaluation_result: "" });
const evaluationFeedback = ref("");

// 读取用户信息
onMounted(() => {
  const storedUser = localStorage.getItem("userInfo");
  if (storedUser) {
    user.value = JSON.parse(storedUser);
  } else {
    ElMessage.error("未找到用户信息，请重新登录");
    router.push("/login");
  }
});

// 选择科目查看评测结果
const viewEvaluation = (subject: any) => {
  selectedSubject.value = subject;
  evaluationResult.value = {
    score: subject.score,
    evaluation_result: subject.evaluation_result
  };
  evaluationDialogVisible.value = true;
};

// 提交反馈
const submitFeedback = () => {
  ElMessage.success("反馈已提交: " + evaluationFeedback.value);
  evaluationFeedback.value = ""; // 清空输入框
};

// 退出登录
const logout = () => {
  localStorage.removeItem("userInfo"); // 清除用户信息
  ElMessage.success("退出成功");
  router.push("/login"); // 跳转到登录页
};
</script>

<style scoped>
.profile-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(to right, #f8f9fa, #e9ecef);
}

.profile-card {
  width: 600px;
  padding: 30px;
  text-align: center;
  border-radius: 12px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

.subject-button {
  width: 100%;
  font-size: 16px;
  padding: 12px 0;
  margin-top: 10px;
  border-radius: 8px;
}

.logout-button {
  margin-top: 20px;
  width: 100%;
  font-size: 16px;
  padding: 10px;
}
</style>

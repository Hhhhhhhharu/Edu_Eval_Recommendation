<template>
  <div class="evaluation-container">
    <el-card class="evaluation-card">
      <h2>试题列表</h2>
      <el-divider></el-divider>
      <el-row justify="space-around" gutter="30">
        <el-col v-for="course in courses" :key="course.id" :span="7">
          <el-button type="primary" class="course-button" @click="selectCourse(course)">
            {{ course.name }}
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 试题界面 -->
    <el-dialog v-model="examDialogVisible" :title="selectedCourse.name" width="600px">
      <el-scrollbar height="400px">
        <el-card v-for="exam in selectedExams" :key="exam.id" class="exam-card">
          <h3>{{ exam.title }}</h3>
          <p>{{ exam.question }}</p>
          <el-radio-group v-model="exam.selectedAnswer">
            <el-radio v-for="(option, index) in exam.options" :key="index" :label="option">
              {{ option }}
            </el-radio>
          </el-radio-group>
        </el-card>
      </el-scrollbar>
      <template #footer>
        <el-button @click="examDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAnswers">提交答案</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";

// 课程列表
const courses = ref([
  { id: 1, name: "数据结构" },
  { id: 2, name: "Python" },
  { id: 3, name: "操作系统" }
]);

// 试题数据
const examsData = {
  数据结构: [
    { id: 1, title: "树结构", question: "树的遍历方式有几种？", options: ["2种", "3种", "4种", "5种"], selectedAnswer: "" },
    { id: 2, title: "堆", question: "堆的特点是什么？", options: ["先进先出", "完全二叉树", "有序数组", "无序集合"], selectedAnswer: "" }
  ],
  Python: [
    { id: 3, title: "Python 基础", question: "Python 的数据类型有哪些？", options: ["int", "str", "list", "以上都是"], selectedAnswer: "" },
    { id: 4, title: "Python 进阶", question: "Python 中的生成器是什么？", options: ["一种函数", "一种类", "一种装饰器", "一种数据结构"], selectedAnswer: "" }
  ],
  操作系统: [
    { id: 5, title: "进程管理", question: "进程与线程的区别是？", options: ["无区别", "线程是轻量级的进程", "进程是轻量级的线程", "以上都不对"], selectedAnswer: "" }
  ]
};

const examDialogVisible = ref(false);
const selectedCourse = ref({ id: 0, name: "" });
const selectedExams = ref([] as any[]);

// 选择课程，显示试题
const selectCourse = (course: any) => {
  selectedCourse.value = course;
  selectedExams.value = examsData[course.name] || [];
  examDialogVisible.value = true;
};

// 提交答案
const submitAnswers = () => {
  if (selectedExams.value.some(exam => !exam.selectedAnswer)) {
    ElMessage.warning("请完成所有试题");
    return;
  }

  console.log("提交的答案：", selectedExams.value.map(exam => ({
    examId: exam.id,
    selectedAnswer: exam.selectedAnswer
  })));

  ElMessage.success("作答完成");
  examDialogVisible.value = false;
};
</script>

<style scoped>
/* 让容器占满整个页面 */
.evaluation-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(to right, #eef2f3, #dfe9f3);
}

/* 加大卡片大小，使其更显眼 */
.evaluation-card {
  width: 600px;
  padding: 30px;
  text-align: center;
  border-radius: 12px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

/* 课程按钮样式 */
.course-button {
  width: 100%;
  font-size: 18px;
  padding: 15px 0;
  margin-top: 15px;
  border-radius: 8px;
}

/* 试题卡片样式 */
.exam-card {
  margin-bottom: 15px;
  padding: 15px;
  border-radius: 8px;
}
</style>

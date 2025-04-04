<template> 
  <div class="login-container">
    <el-card class="login-card">
      <h2>用户登录</h2>
      <el-form :model="form" ref="loginForm" label-width="80px">
        
        <!-- 开发模式下的角色切换 -->
        <el-form-item v-if="isDev" label="选择角色">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>

        <el-form-item label="学号">
          <el-input v-model="form.studentId" placeholder="请输入学号" />
        </el-form-item>
        
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleLogin">登录</el-button>
          <el-button type="text" @click="goToRegister">没有账号？去注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";

const router = useRouter();

// 是否为开发模式
const isDev = computed(() => import.meta.env.MODE === "development");

// 登录表单
const form = reactive({
  studentId: "",
  password: "",
  role: "student", // 默认角色
});

// 模拟数据库用户数据
const users = [
  { studentId: "2021131090", name: "学生", password: "123456", role: "student" },
  { studentId: "2021132090", name: "老师", password: "123456", role: "teacher" },
  { studentId: "2021133090", name: "教学单位", password: "123456", role: "admin" },
];

// 处理登录
const handleLogin = (): void => {
  if (!form.studentId || !form.password) {
    ElMessage.error("请输入学号和密码");
    return;
  }

  // 查找匹配的用户
  const user = users.find(
    (u) => u.studentId === form.studentId && u.password === form.password
  );

  if (user) {
    ElMessage.success("登录成功");

    // 存储用户信息
    localStorage.setItem("userInfo", JSON.stringify(user));


    const finalRole = user.role;


    // 跳转到对应的主页
    switch (finalRole) {
      case "student":
        router.push("/student");
        break;
      case "teacher":
        router.push("/teacher");
        break;
      case "admin":
        router.push("/admin");
        break;
      default:
        router.push("/student");
    }
  } else {
    ElMessage.error("学号或密码错误");
  }
};

const goToRegister = (): void => {
  router.push("/register");
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}

.login-card {
  width: 400px;
  padding: 25px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

h2 {
  margin-bottom: 20px;
}
</style>

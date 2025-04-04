<template>
    <div class="register-container">
      <el-card class="register-card">
        <h2>用户注册</h2>
        <el-form :model="form" ref="registerForm" label-width="80px">
            <!-- 角色选择 -->
          <el-form-item label="选择角色">
            <el-select v-model="form.role" placeholder="请选择角色">
              <el-option label="学生" value="student" />
              <el-option label="教师" value="teacher" />
              <el-option label="管理员" value="admin" />
            </el-select>
          </el-form-item>
          <el-form-item label="学号">
            <el-input v-model="form.studentId" placeholder="请输入学号" />
          </el-form-item>
          <el-form-item label="姓名">
            <el-input v-model="form.name" placeholder="请输入姓名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
          </el-form-item>
          <el-form-item label="确认密码">
            <el-input v-model="form.confirmPassword" type="password" show-password placeholder="请确认密码" />
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="handleRegister">注册</el-button>
            <el-button type="text" @click="goToLogin">已有账号？去登录</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </template>
  
  <script setup lang="ts">
  import { reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { ElMessage } from 'element-plus';
  
  const router = useRouter();
  
  interface RegisterForm {
    studentId: string;
    name: string;
    password: string;
    confirmPassword: string;
    role: string; // 新增角色字段
  }
  
  const form = reactive<RegisterForm>({
    studentId: '',
    name: '',
    password: '',
    confirmPassword: '',
    role: '' // 默认空
  });
  
  const handleRegister = (): void => {
    if (!form.studentId || !form.name || !form.password || !form.confirmPassword || !form.role) {
      ElMessage.error('请填写完整信息');
      return;
    }
    if (form.password !== form.confirmPassword) {
      ElMessage.error('两次输入的密码不一致');
      return;
    }
  
    // 存储角色到 localStorage
    localStorage.setItem('userRole', form.role);
  
    ElMessage.success('注册成功');
  
    // 跳转到对应的角色页面
    switch (form.role) {
      case 'student':
        router.push('/student');
        break;
      case 'teacher':
        router.push('/teacher');
        break;
      case 'admin':
        router.push('/admin');
        break;
      default:
        router.push('/student');
    }
  };
  
  const goToLogin = (): void => {
    router.push('/login');
  };
  </script>
  
  <style scoped>
  .register-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f5f5f5;
  }
  
  .register-card {
    width: 350px;
    padding: 20px;
    text-align: center;
  }
  </style>
  
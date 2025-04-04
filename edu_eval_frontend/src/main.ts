import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { createApp } from 'vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import App from './App.vue';
import StudentHome from './views/StudentHome.vue';
import TeacherHome from './views/TeacherHome.vue';
import AdminHome from './views/AdminHome.vue';
import Login from './components/Login.vue';
import Register from './components/Register.vue';
import { ElMessage } from 'element-plus';


// 模拟登录状态（手动设置角色）
const simulateLogin = (role: 'student' | 'teacher' | 'admin') => {
  localStorage.setItem('userRole', role);
  return true;
};

// 检查是否已登录（开发模式下直接返回 true，生产环境需取消注释真实逻辑）
const isAuthenticated = () => {
  // 开发模式：绕过验证
  return true;

  // 生产环境：真实验证（注释掉上面的 return true）
  // return !!localStorage.getItem('userRole');
};

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/login' },
  { 
    path: '/login', 
    component: Login,
    // 在登录页面手动模拟角色（开发模式）
    beforeEnter: (to, from, next) => {
        const isDevMode = import.meta.env.MODE === 'development'; //判断是否是开发模式
        if (isDevMode) {
            // 模拟登录为 student（可改为 teacher 或 admin）
            // simulateLogin('student');
            // simulateLogin('teacher');
            simulateLogin('admin');
        }
        next();
    }
  },
  { path: '/register', component: Register },
  { 
    path: '/student', 
    component: StudentHome, 
    beforeEnter: (to, from, next) => {
      if (isAuthenticated() && localStorage.getItem('userRole') === 'student') {
        next();
      } else {
        ElMessage.error('仅限学生访问');
        next('/login');
      }
    }
  },
  { 
    path: '/teacher', 
    component: TeacherHome, 
    beforeEnter: (to, from, next) => {
      if (isAuthenticated() && localStorage.getItem('userRole') === 'teacher') {
        next();
      } else {
        ElMessage.error('仅限教师访问');
        next('/login');
      }
    }
  },
  { 
    path: '/admin', 
    component: AdminHome, 
    beforeEnter: (to, from, next) => {
      if (isAuthenticated() && localStorage.getItem('userRole') === 'admin') {
        next();
      } else {
        ElMessage.error('仅限管理员访问');
        next('/login');
      }
    }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

const app = createApp(App);
app.use(ElementPlus);
app.use(router);
app.mount('#app');
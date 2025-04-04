<template>
  <div class="student-home">
    <!-- 顶部标题 -->
    <el-header class="header">
      <h2 class="header-title">学生首页</h2>
    </el-header>

    <!-- 内容区 -->
    <el-main class="main-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="评测" name="evaluation">
          <Evaluation />
        </el-tab-pane>
        <el-tab-pane label="学习资料" name="resources">
          <Resources />
        </el-tab-pane>
        <el-tab-pane label="个人信息" name="profile">
          <Profile />
        </el-tab-pane>
      </el-tabs>
    </el-main>

    <!-- 底部导航栏 -->
    <el-footer class="footer">
      <el-row justify="space-around" class="nav-container">
        <el-col :span="8" class="nav-item" @click="activeTab = 'evaluation'">
          <el-icon><Document /></el-icon>
          <span>评测</span>
        </el-col>
        <el-col :span="8" class="nav-item" @click="activeTab = 'resources'">
          <el-icon><Folder /></el-icon>
          <span>学习资料</span>
        </el-col>
        <el-col :span="8" class="nav-item" @click="activeTab = 'profile'">
          <el-icon><User /></el-icon>
          <span>个人信息</span>
        </el-col>
      </el-row>
    </el-footer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Evaluation from '@/components/student/Evaluation.vue';
import Resources from '@/components/student/Resources.vue';
import Profile from '@/components/student/Profile.vue';
import { Document, Folder, User } from '@element-plus/icons-vue';

// 默认选中的标签页
const activeTab = ref('evaluation');
</script>

<style scoped>
.student-home {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

/* 顶部标题 */
.header {
  background-color: #409eff;
  color: white;
  text-align: center;
  padding: 15px 0;
  font-size: 18px;
}

.header-title {
  margin: 0;
  font-weight: bold;
}

/* 主体部分 */
.main-content {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
  margin-bottom: 70px; /* 避免内容被底部导航遮挡 */
}

/* 底部导航栏 */
.footer {
  background-color: white;
  border-top: 2px solid #ddd;
  position: fixed;
  bottom: 0;
  width: 100%;
  height: 60px;
  display: flex;
  align-items: center;
  z-index: 1000;
}

/* 导航按钮 */
.nav-container {
  width: 100%;
  text-align: center;
}

.nav-item {
  text-align: center;
  cursor: pointer;
  padding: 5px;
  font-size: 14px;
  color: #666;
  transition: color 0.3s ease-in-out;
}

.nav-item:hover {
  color: #409eff;
}

.nav-item span {
  display: block;
  margin-top: 5px;
  font-size: 13px;
}
</style>

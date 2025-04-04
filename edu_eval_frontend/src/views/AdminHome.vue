<template>
  <div class="admin-home">
    <!-- 顶部标题 -->
    <el-header class="header">
      <h2 class="header-title">教学单位首页</h2>
    </el-header>

    <!-- 内容区 -->
    <el-main class="main-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="教学资料" name="teaching">
          <TeachingMaterials />
        </el-tab-pane>
        <el-tab-pane label="班级信息" name="grade-info">
          <GradeInfo />
        </el-tab-pane>
      </el-tabs>
    </el-main>

    <!-- 底部导航栏 -->
    <el-footer class="footer">
      <el-row justify="space-around" class="nav-container">
        <el-col :span="12" class="nav-item" @click="activeTab = 'teaching'">
          <el-icon><Folder /></el-icon>
          <span>教学资料</span>
        </el-col>
        <el-col :span="12" class="nav-item" @click="activeTab = 'grade-info'">
          <el-icon><User /></el-icon>
          <span>班级信息</span>
        </el-col>
      </el-row>
    </el-footer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import TeachingMaterials from '@/components/admin/TeachingMaterials.vue';
import GradeInfo from '@/components/admin/GradeInfo.vue';
import { Folder, User } from '@element-plus/icons-vue';

// 默认选中的标签页
const activeTab = ref('teaching');
</script>

<style scoped>
.admin-home {
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

<template>
  <div class="resources-container">
    <el-card class="resources-card">
      <h2>学习资料</h2>
      
      <!-- 搜索框 -->
      <el-input v-model="searchQuery" placeholder="搜索学习资料" clearable @input="filterResources" />
      
      <el-divider></el-divider>
      
      <el-table :data="filteredResources" style="width: 100%">
        <el-table-column prop="title" label="资料名称" />
        <el-table-column prop="subject" label="科目" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" @click="viewResource(scope.row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";

const searchQuery = ref("");
const resources = ref([
  { id: 1, title: "数据结构学习指南", subject: "数据结构" },
  { id: 2, title: "Python 编程基础", subject: "Python" },
  { id: 3, title: "操作系统原理", subject: "操作系统" }
]);
const filteredResources = ref([...resources.value]);

const filterResources = () => {
  if (!searchQuery.value) {
    filteredResources.value = resources.value;
  } else {
    filteredResources.value = resources.value.filter(resource =>
      resource.title.includes(searchQuery.value) ||
      resource.subject.includes(searchQuery.value)
    );
  }
};

const viewResource = (resource) => {
  ElMessage.info(`打开 ${resource.title}`);
};

onMounted(() => {
  filterResources();
});
</script>

<style scoped>
.resources-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(to right, #f8f9fa, #e9ecef);
}

.resources-card {
  width: 600px;
  padding: 30px;
  text-align: center;
  border-radius: 12px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}
</style>

<template>
  <div class="materials-container">
    <el-card class="materials-card">
      <h2>教学资料推荐</h2>
      
      <!-- 搜索框 -->
      <el-input v-model="searchQuery" placeholder="搜索教学资料" clearable />
      
      <el-divider></el-divider>
      
      <el-empty v-if="filteredMaterials.length === 0" description="暂无相关资料" />
      
      <el-row gutter="20" v-else>
        <el-col v-for="material in filteredMaterials" :key="material.id" :span="12">
          <el-card shadow="hover" class="material-item">
            <h3>{{ material.title }}</h3>
            <p>{{ material.description }}</p>
            <el-button type="primary" @click="viewMaterial(material)">查看详情</el-button>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  
    <!-- 教学资料详情对话框 -->
    <el-dialog v-model="dialogVisible" :title="selectedMaterial.title" width="500px">
      <p>{{ selectedMaterial.description }}</p>
      <el-button type="primary" @click="downloadMaterial">下载</el-button>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";

const materials = ref([
  { id: 1, title: "数据结构教材", description: "全面解析数据结构与算法。", link: "#" },
  { id: 2, title: "Python 编程基础", description: "入门 Python，掌握编程技能。", link: "#" },
  { id: 3, title: "操作系统原理", description: "学习操作系统的核心概念。", link: "#" }
]);

const searchQuery = ref("");
const dialogVisible = ref(false);
const selectedMaterial = ref({ title: "", description: "", link: "" });

// 过滤教学资料
const filteredMaterials = computed(() => {
  return materials.value.filter(material =>
    material.title.includes(searchQuery.value) ||
    material.description.includes(searchQuery.value)
  );
});

// 选择资料查看详情
const viewMaterial = (material: any) => {
  selectedMaterial.value = material;
  dialogVisible.value = true;
};

// 下载资料（模拟）
const downloadMaterial = () => {
  window.open(selectedMaterial.value.link, "_blank");
};
</script>

<style scoped>
.materials-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
  background: #f9f9f9;
}

.materials-card {
  width: 80%;
  padding: 20px;
  text-align: center;
}

.material-item {
  margin-bottom: 20px;
  padding: 15px;
}
</style>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import Login from './components/Login.vue';
import Dashboard from './components/Dashboard.vue';

const isAuthenticated = ref(false);

const checkAuth = () => {
  isAuthenticated.value = !!localStorage.getItem('token');
};

onMounted(() => {
  checkAuth();
});

watch(() => localStorage.getItem('token'), () => {
  checkAuth();
});
</script>

<template>
  <Login v-if="!isAuthenticated" />
  <Dashboard v-else />
</template>

<style scoped>
</style>

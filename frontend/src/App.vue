<script setup>
  import AppFooter from './components/AppFooter.vue';
import AppHeader from './components/AppHeader.vue';
</script>

<template>
  <AppHeader />
  <main class="main-content">
    <RouterView />
  </main>
  <AppFooter />
</template>

<style scoped>
.main-content {
  flex: 1;
  padding-bottom: 80px;
}
</style>

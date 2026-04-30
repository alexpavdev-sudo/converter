<template>
  <div class="app-wrapper">
    <header class="app-header">
      <div class="container header-container">
        <router-link to="/" class="logo">Converter</router-link>
        <nav class="nav-menu">
          <router-link to="/converter-images" class="nav-link">Конвертер файлов</router-link>
          <router-link to="/download" class="nav-link">Мои файлы</router-link>
          <router-link to="/contact" class="nav-link">Контакты</router-link>
        </nav>
      </div>
    </header>

    <main class="main-content">
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component"/>
        </keep-alive>
      </router-view>
    </main>

    <footer class="app-footer">
      <div class="container footer-container">
        <div>©2026 Converter</div>
      </div>
    </footer>
  </div>
</template>

<style lang="scss">
</style>
<script setup>
</script>
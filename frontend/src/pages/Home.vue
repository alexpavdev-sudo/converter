<!-- src/views/Home.vue -->
<template>
  <div class="container">
    <div class="section">
      <h1 class="section-title">Welcome to Our App</h1>

      <div class="card-grid">
        <div class="card" v-for="item in items" :key="item.id">
          <div class="card-header">
            <h3 class="card-title">{{ item.title }}</h3>
          </div>
          <div class="card-content">
            <p>{{ item.description }}</p>
          </div>
          <div class="card-footer">
            <button class="btn btn-primary btn-sm">Learn More</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const items = ref([
  { id: 1, title: 'Feature 1', description: 'Description for feature 1' },
  { id: 2, title: 'Feature 2', description: 'Description for feature 2' },
  { id: 3, title: 'Feature 3', description: 'Description for feature 3' }
])
</script>
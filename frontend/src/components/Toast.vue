<template>
  <div class="toast-container">
    <transition-group name="toast" tag="div">
      <div
          v-for="toast in toasts"
          :key="toast.id"
          class="toast-item"
          :class="{
          'toast-success': toast.type === 'success',
          'toast-error': toast.type === 'error',
        }"
          @click="removeToast(toast.id)"
      >
        <strong class="toast-title">{{ toast.title }}</strong>
        <p class="toast-detail">{{ toast.detail }}</p>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { toast } from '@/services/toast';

const { toasts, removeToast } = toast();
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 20px;
  left: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 400px;
  width: 100%;
}

.toast-item {
  padding: 12px 16px;
  margin: 10px 0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  color: #fff;
  cursor: pointer;
  backdrop-filter: blur(8px);
  transition: all 0.3s ease;
}

.toast-success {
  background-color: rgba(56, 142, 60, 0.9);
}

.toast-error {
  background-color: rgba(211, 47, 47, 0.9);
}

.toast-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
  display: block;
}

.toast-detail {
  font-size: 14px;
  margin: 0;
  line-height: 1.4;
}

/* Анимации transition-group */
.toast-enter-active {
  transition: all 0.4s ease-out;
}
.toast-leave-active {
  transition: all 0.3s ease-in;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
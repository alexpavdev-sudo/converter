<template>
  <div class="image-list">
    <div class="list-header">
      <h3>📋 Сконвертированные файлы</h3>
      <p>Файлы будут храниться в течение 24 часов</p>
    </div>

    <div class="list-items">
      <div v-for="image in images" :key="image.id" class="image-item">
        <!-- Превью -->
        <div class="preview">
          <img src="" :alt="image.original_name"/>
        </div>
        <div class="info">
          <div class="name">{{ image.original_name }}</div>
        </div>
        <div class="info">
          <div>
            <div class="size">Исходный</div>
            <div class="original-format">
              {{ image.extension.toUpperCase() }} / {{ func.formatFileSize(image.size) }}
            </div>
          </div>
        </div>
        <div class="info">
          <div>
            <div
                :class="getStatusType(image.status)"
                :title="image.status == 3 ? 'Нажмите чтобы увидеть ошибку' : ''"
                @click="image.status == 3 ? showError(image.id) : null">{{ image.status_label }}
            </div>
            <div class="original-format">
              {{ image.format.toUpperCase() }} / {{ func.formatFileSize(image.size_processed) }}
            </div>
          </div>
        </div>
        <div class="info">
          <button @click="downloadImage(image.id)" v-if="image.status == 2" class="btn-sm btn-primary"
                  title="Скачать">Скачать
          </button>
        </div>
        <button @click="removeImage(image.id)" class="remove-btn" title="Удалить">✕</button>
      </div>
    </div>

    <div v-if="images.length === 0" class="empty-state">
      <p>Нет сконвертированных изображений</p>
      <p class="hint">
        <router-link to="/converter-images">➡️ Загрузите изображения для конвертации</router-link>
      </p>
    </div>

    <!-- Модальное окно ошибки -->
    <div v-if="errorModalVisible" class="modal-overlay" @click="closeErrorModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h4>Ошибка конвертации</h4>
          <button @click="closeErrorModal" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <p>{{ currentError }}</p>
        </div>
        <div class="modal-footer">
          <button @click="closeErrorModal" class="btn-sm btn-primary">Закрыть</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import func from "@/services/functionHelper.js";
import type {File} from '@/types/file'
import api from '@/services/api';

const props = defineProps<{
  images: File[]
}>()

const emit = defineEmits(['remove', 'download', 'update-format'])

const errorModalVisible = ref(false)
const currentError = ref('')

const getStatusType = (status: number) => {
  switch (status) {
    case 0:
      return 'status-queued'
    case 1:
      return 'status-processing'
    case 2:
      return 'status-processed'
    case 3:
      return 'status-error'
  }
}

const showError = async (id: number) => {
  const {data} = await api.get(`/api/files/error/${id}`);
  currentError.value = data.data
  errorModalVisible.value = true
}

const closeErrorModal = () => {
  errorModalVisible.value = false
}

const removeImage = (id) => {
  emit('remove', id)
}
const downloadImage = (id) => {
  emit('download', id)
}

const updateFormat = (id, format) => {
  emit('update-format', id, format)
}
</script>

<style lang="scss" scoped>
.image-list {
  margin: 20px 0;
}

.list-header {
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid $gray-200;
}

.list-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.image-item {
  display: grid;
  grid-template-columns: 80px 2fr 2fr 2fr 1fr 40px;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: $white;
  border-radius: $border-radius-md;
  box-shadow: $box-shadow-sm;
  transition: all $transition-base;

  &:hover {
    box-shadow: $box-shadow-md;
  }
}

.preview {
  width: 40px;
  height: 40px;
  overflow: hidden;
  border-radius: $border-radius-sm;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.info {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 5px;
  align-items: center;

  .name {
    font-weight: 500;
    margin-bottom: 5px;
    @include text-ellipsis(1);
  }

  .size {
    font-size: $font-size-sm;
    color: $gray-600;
    margin-bottom: 3px;
  }

  .original-format {
    font-size: $font-size-sm;
    color: $gray-500;
  }
}

.format-select {
  select {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid $gray-300;
    border-radius: $border-radius-sm;
    background: $white;
    cursor: pointer;

    &:focus {
      outline: none;
      border-color: $primary-color;
    }
  }
}

.btn-error {
  border: 1px solid $danger-color;
  color: $danger-color;
}

.remove-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: $danger-color;
  color: white;
  border: none;
  cursor: pointer;
  font-size: 18px;
  transition: all $transition-base;

  &:hover {
    background: darken($danger-color, 10%);
    transform: scale(1.1);
  }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: $gray-100;
  border-radius: $border-radius-lg;

  .hint {
    color: $gray-600;
    font-size: $font-size-sm;
    margin-top: 10px;
  }
}

@include respond-to(sm) {
  .image-item {
    grid-template-columns: 60px 1fr;
    gap: 10px;

    .format-select {
      grid-column: 1 / -1;
    }

    .remove-btn {
      position: absolute;
      top: 10px;
      right: 10px;
    }
  }
}

//стили модального окна
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: $white;
  border-radius: $border-radius-md;
  box-shadow: $box-shadow-lg;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow: hidden;
  animation: modalFadeIn 0.3s ease-out;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid $gray-200;
  display: flex;
  justify-content: space-between;
  align-items: center;

  h4 {
    margin: 0;
    color: $gray-800;
  }
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: $gray-500;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    color: $gray-800;
  }
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  max-height: 400px;

  p {
    margin: 0;
    color: $danger-color;
    white-space: pre-wrap;
  }
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid $gray-200;
  text-align: right;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
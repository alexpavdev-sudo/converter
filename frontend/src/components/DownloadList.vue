<!-- src/components/ImageList.vue -->
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
            <div class="size">{{ image.status_label }}</div>
            <div class="original-format">
              {{ image.format.toUpperCase() }} / {{ func.formatFileSize(image.size_processed) }}
            </div>
          </div>
        </div>
        <div class="info">
          <button @click="downloadImage(image.id)" v-if="image.status == 2" class="btn-sm btn-primary" title="Скачать">Скачать</button>
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
  </div>
</template>

<script setup lang="ts">
import func from "@/services/functionHelper.js";
import type {File} from '@/types/file'

const props = defineProps<{
  images: File[]
}>()

const emit = defineEmits(['remove', 'download', 'update-format'])

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
</style>
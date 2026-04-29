<template>
  <div class="image-converter">
    <!-- Форма загрузки -->
    <div class="upload-section">
      <div class="upload-area"
           @dragover.prevent
           @drop.prevent="handleDrop"
           :class="{ 'drag-over': isDragging }">

        <input type="file"
               ref="fileInput"
               multiple
               accept="image/*, video/*"
               @change="handleFileSelect"
               style="display: none"/>

        <button @click="triggerFileSelect" class="btn btn-primary">
          📁 Выбрать файлы
        </button>

        <button @click="handleAddMore" class="btn btn-secondary"
                v-if="uploadFiles.length > 0">
          ➕ Добавить ещё
        </button>

        <p class="upload-hint" @click="fetchPing">или перетащите файлы сюда</p>
      </div>
    </div>

    <!-- Список изображений -->
    <ImageList
        :files="uploadFiles"
        :formats="formats"
        @remove="removeImage"
        @update-format="updateFormat"
    />

    <!-- Кнопка конвертации -->
    <div class="actions" v-if="uploadFiles.length > 0">
      <button @click="convertAll" class="btn btn-success btn-large">
        <span v-if="!isLoading">🔄 Конвертировать</span>
        <span v-else>⏳ Загрузка...</span>
        ({{ uploadFiles.length }} файлов)
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onActivated, ref} from 'vue'
import {useRouter} from 'vue-router'
import ImageList from './ImageList.vue'
import api from '@/services/api';
import {Format} from '@/types/format';

const router = useRouter()
const fileInput = ref(null)
const isDragging = ref(false)
const isLoading = ref(false)
const uploadFiles = ref([])
const formats = ref<Format[]>([]);

onActivated(async () => {
  await fetchFormats()
})

const fetchFormats = async () => {
  try {
    const response = await api.get('/api/formats/');
    const data = response.data.data;
    formats.value = Object.values(data);
  } catch (error) {
    console.error(error);
    formats.value = [];
  }
}

// Обработка выбора файлов
const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)
  addFiles(files)
}

// Обработка drag & drop
const handleDrop = (event) => {
  isDragging.value = false
  const files = Array.from(event.dataTransfer.files)
  addFiles(files)
}

// Добавление файлов в список
const addFiles = (newFiles) => {
  const files = newFiles.filter(
      file => (file.type.startsWith('image/') || file.type.startsWith('video/'))
  )

  files.forEach(file => {
    // Проверка на дубликаты
    if (!uploadFiles.value.some(img => img.name === file.name && img.size === file.size)) {
      uploadFiles.value.push({
        id: uploadFiles.value.length + 1,
        file: file,
        original_name: file.name,
        size: file.size,
        extension: file.type.split('/')[1],
        format: 'webp', // формат по умолчанию
      })
    }
  })

  // Очищаем input
  if (fileInput.value) fileInput.value.value = ''
}

// Триггер выбора файлов
const triggerFileSelect = () => {
  fileInput.value.click()
}

// Добавление новых файлов к существующим
const handleAddMore = () => {
  fileInput.value.click()
}

// Удаление изображения
const removeImage = (id) => {
  const index = uploadFiles.value.findIndex(img => img.id === id)
  if (index !== -1) {
    uploadFiles.value.splice(index, 1)
  }
}

// Обновление формата конвертации
const updateFormat = (id, format) => {
  const image = uploadFiles.value.find(img => img.id === id)
  if (image) {
    image.format = format
  }
}

// Конвертация всех изображений
const convertAll = async () => {
  try {
    isLoading.value = true

    const formData = new FormData()
    uploadFiles.value.forEach(img => {
      formData.append('formats', img.format)
      formData.append('images', img.file)
    })

    const response = await api.post('/api/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      }
    });
    uploadFiles.value = []
    isLoading.value = false
  } catch (error) {
    console.error('Произошла ошибка при отправке файлов:', error)
    alert('Произошла ошибка при отправке файлов')
  }
}
</script>

<style lang="scss" scoped>

.image-converter {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.upload-section {
  margin-bottom: 30px;
}

.upload-area {
  border: 2px dashed $gray-400;
  border-radius: $border-radius-lg;
  padding: 40px;
  text-align: center;
  background: $gray-100;
  transition: all $transition-base;

  &.drag-over {
    border-color: $primary-color;
    background: rgba($primary-color, 0.05);
  }
}

.upload-hint {
  margin-top: 15px;
  color: $gray-600;
  font-size: $font-size-sm;
}

.actions {
  text-align: center;
  margin-top: 30px;
}

.btn-large {
  padding: $spacing-md $spacing-xl;
  font-size: $font-size-lg;
}
</style>
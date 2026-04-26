<template>
  <div class="image-download">
    <DownloadList
        :images="images"
        @remove="removeImage"
    />
  </div>
</template>

<script setup lang="ts">
import {onActivated, ref} from 'vue'
import DownloadList from "./DownloadList.vue";
import api from '@/services/api';
import type {File} from '@/types/file';

const images = ref<File[]>([]);

onActivated(async () => {
  await fetchImages()
})

const fetchImages = async () => {
  try {
    const response = await api.get('/api/files/');
    const data = response.data.data;
    if (Array.isArray(data)) {
      images.value = data as FileResponse[];
    } else {
      images.value = [];
    }
  } catch (error) {
    console.error(error);
    images.value = [];
  }
}

const removeImage = async (id) => {
  try {
    const { data } = await api.delete(`/api/files/${id}`);

    if (data.success) {
      const index = images.value.findIndex(img => img.id === id);
      if (index !== -1) {
        images.value.splice(index, 1);
      }
    }
  } catch (error) {
    console.error('Failed to remove image:', error);
  }
}

</script>

<style lang="scss" scoped>
.image-download {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
</style>
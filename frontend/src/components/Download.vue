<template>
  <div class="image-download">
    <DownloadList
        :files="files"
        @remove="removeImage"
        @download="downloadFile"
    />
  </div>
</template>

<script setup lang="ts">
import {onActivated, onMounted, onUnmounted, ref} from 'vue'
import DownloadList from "./DownloadList.vue";
import api from '@/services/api';
import type {File} from '@/types/file';
import {NotificationDto, WSMessage} from '@/types/ws';
import {toast} from '@/services/toast';
import {useWebSocket} from "@/services/useWebSocket";

const files = ref<File[]>([]);

const {showToast} = toast();

onActivated(async () => {
  await fetchImages()
})

const fetchImages = async () => {
  try {
    const response = await api.get('/api/files/');
    const data = response.data.data;
    if (Array.isArray(data)) {
      files.value = data as FileResponse[];
    } else {
      files.value = [];
    }
  } catch (error) {
    console.error(error);
    files.value = [];
  }
}

const removeImage = async (id) => {
  try {
    const {data} = await api.delete(`/api/files/${id}`);

    if (data.success) {
      const index = files.value.findIndex(img => img.id === id);
      if (index !== -1) {
        files.value.splice(index, 1);
      }
    }
  } catch (error) {
    console.error('Failed to remove image:', error);
  }
}

const downloadFile = async (id) => {
  try {
    const url = `/api/files/download/${id}`;
    window.open(url, '_blank');
  } catch (error) {
  }
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null;

const handleMessage = (msg: WSMessage) => {
  if (msg.type === 1) {
    const payload = msg.payload as NotificationDto;
    let detail = JSON.parse(payload.detail)
    showToast('', detail.data, detail.success ? 'success' : 'error');

    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      fetchImages();
      debounceTimer = null;
    }, 1000);
  }
};

const {onMessage} = useWebSocket();
onMounted(() => {
  const unsubscribe = onMessage(handleMessage);
  onUnmounted(unsubscribe); // отписываемся при размонтировании
});
</script>

<style lang="scss" scoped>
.image-download {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
</style>
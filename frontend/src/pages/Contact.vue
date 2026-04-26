<!-- src/views/Contact.vue -->
<template>
  <div class="container">
    <div class="section">
      <h1 class="section-title">Contact Us</h1>

      <form @submit.prevent="handleSubmit" class="form">
        <div class="form-group">
          <label class="form-label">Name</label>
          <input
              type="text"
              v-model="form.name"
              class="form-control"
              :class="{ 'is-invalid': errors.name }"
          />
          <span v-if="errors.name" class="form-text error">{{ errors.name }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">Email</label>
          <input
              type="email"
              v-model="form.email"
              class="form-control"
              :class="{ 'is-invalid': errors.email }"
          />
          <span v-if="errors.email" class="form-text error">{{ errors.email }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">Message</label>
          <textarea
              v-model="form.message"
              class="form-control"
              rows="5"
          ></textarea>
        </div>

        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Send Message</button>
          <button type="reset" class="btn btn-outline-primary">Clear</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'

defineOptions({
  name: 'Contact'
})

const form = reactive({
  name: '',
  email: '',
  message: ''
})

const errors = reactive({
  name: '',
  email: ''
})

const handleSubmit = () => {
  // Логика отправки
  console.log('Form submitted:', form)
}
</script>
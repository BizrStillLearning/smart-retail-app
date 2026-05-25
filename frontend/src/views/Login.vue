<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref(null)

const prosesLogin = async () => {
  if (!username.value || !password.value) {
    error.value = "Username dan password wajib diisi"
    return
  }

  loading.value = true
  error.value = null

  try {
    const response = await axios.post('http://localhost:8082/api/login', {
      username: username.value,
      password: password.value
    })

    authStore.setAuth(response.data.token, response.data.user)

    alert(`Selamat datang, ${response.data.user.nama_lengkap}!`)

    router.push('/admin')
  } catch (err) {
    error.value = err.response?.data?.error || 'Gagal terhubung ke server'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-2xl shadow-lg border border-gray-100 p-8">

      <div class="text-center mb-8">
        <h1 class="text-3xl font-black text-gray-900">Smart Retail</h1>
        <p class="text-gray-500 mt-2">Login untuk masuk ke Dasbor Admin</p>
      </div>

      <div v-if="error" class="mb-4 p-3 bg-red-50 text-red-600 rounded-lg text-sm border border-red-100">
        {{ error }}
      </div>

      <form @submit.prevent="prosesLogin" class="space-y-5">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input
              v-model="username"
              type="text"
              placeholder="Masukkan username"
              class="w-full px-4 py-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
          >
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
              v-model="password"
              type="password"
              placeholder="••••••••"
              class="w-full px-4 py-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
          >
        </div>

        <button
            type="submit"
            :disabled="loading"
            class="w-full bg-blue-600 text-white font-bold py-3 rounded-lg hover:bg-blue-700 transition disabled:bg-blue-300"
        >
          {{ loading ? 'Memproses...' : 'Masuk ke Dasbor' }}
        </button>
      </form>

      <div class="mt-6 text-center">
        <router-link to="/" class="text-sm text-gray-500 hover:text-blue-600 font-medium">
          ← Kembali ke Katalog Toko
        </router-link>
      </div>

    </div>
  </div>
</template>
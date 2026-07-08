<script setup>
import { ref } from 'vue'
import { useRouter } from 'nuxt/app'
import { useCustomerAuthStore } from '~/stores/customerAuth'
import { useSwal } from '~/composables/useSwal'

const authStore = useCustomerAuthStore()
const router = useRouter()
const swal = useSwal()

const form = ref({ email: '', password: '' })
const isSubmitting = ref(false)

const prosesLogin = async () => {
  isSubmitting.value = true
  try {
    const response = await $fetch('http://localhost:8082/api/pelanggan/login', {
      method: 'POST',
      body: form.value
    })

    if (response.status === 'sukses') {
      authStore.setAuth(response.data.token, {
        id_pelanggan: response.data.id_pelanggan,
        nama_pelanggan: response.data.nama_pelanggan
      })
      swal.sukses('Berhasil masuk ke sistem')
      router.push('/')
    }
  } catch (error) {
    swal.error(error.response?._data?.error || 'Kredensial tidak valid.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full bg-white p-8 rounded-3xl border border-slate-100 shadow-sm">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-black text-slate-900 tracking-tight">Selamat Datang</h2>
        <p class="text-sm font-medium text-slate-500 mt-2">Akses akun untuk melanjutkan transaksi.</p>
      </div>

      <form @submit.prevent="prosesLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1">Alamat Email</label>
          <input v-model="form.email" type="email" required class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-2 outline-none transition" placeholder="email@anda.com">
        </div>
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1">Kata Sandi</label>
          <input v-model="form.password" type="password" required class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-2 outline-none transition" placeholder="••••••••">
        </div>

        <button type="submit" :disabled="isSubmitting" class="w-full bg-slate-900 text-white font-black py-3.5 rounded-xl hover:bg-slate-800 transition disabled:opacity-50 mt-4 shadow-md">
          {{ isSubmitting ? 'Memvalidasi...' : 'Masuk Sekarang' }}
        </button>
      </form>

      <div class="text-center mt-6 text-sm font-medium text-slate-500">
        Belum mendaftar?
        <NuxtLink to="/register" class="text-indigo-600 hover:text-indigo-800 font-bold transition">Buat akun baru</NuxtLink>
      </div>
    </div>
  </div>
</template>


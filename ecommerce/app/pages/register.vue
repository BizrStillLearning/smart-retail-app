<script setup>
import { ref } from 'vue'
import { useRouter } from 'nuxt/app'
import { useSwal } from '~/composables/useSwal'

const router = useRouter()
const swal = useSwal()

const form = ref({ nama_pelanggan: '', email: '', password: '' })
const isSubmitting = ref(false)

const prosesRegister = async () => {
  isSubmitting.value = true
  try {
    const response = await $fetch('http://localhost:8082/api/pelanggan/register', {
      method: 'POST',
      body: form.value
    })

    if (response.status === 'sukses') {
      swal.sukses('Akun berhasil dibuat! Silakan masuk.')
      router.push('/login')
    }
  } catch (error) {
    swal.error(error.response?._data?.error || 'Gagal mendaftarkan akun. Silakan coba lagi.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full bg-white p-8 rounded-3xl border border-slate-100 shadow-sm">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-black text-slate-900 tracking-tight">Buat Akun Baru</h2>
        <p class="text-sm font-medium text-slate-500 mt-2">Daftar sekarang untuk mulai berbelanja.</p>
      </div>

      <form @submit.prevent="prosesRegister" class="space-y-4">
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1">Nama Lengkap</label>
          <input v-model="form.nama_pelanggan" type="text" required class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-2 outline-none transition" placeholder="Contoh: John Doe">
        </div>
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1">Alamat Email</label>
          <input v-model="form.email" type="email" required class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-2 outline-none transition" placeholder="email@anda.com">
        </div>
        <div>
          <label class="block text-sm font-bold text-slate-700 mb-1">Kata Sandi</label>
          <input v-model="form.password" type="password" required class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-2 outline-none transition" placeholder="••••••••">
        </div>

        <button type="submit" :disabled="isSubmitting" class="w-full bg-indigo-600 text-white font-black py-3.5 rounded-xl hover:bg-indigo-700 transition disabled:opacity-50 mt-4 shadow-md shadow-indigo-100">
          {{ isSubmitting ? 'Memproses Data...' : 'Daftar Akun' }}
        </button>
      </form>

      <div class="text-center mt-6 text-sm font-medium text-slate-500">
        Sudah memiliki akun?
        <NuxtLink to="/login" class="text-indigo-600 hover:text-indigo-800 font-bold transition">Masuk di sini</NuxtLink>
      </div>
    </div>
  </div>
</template>
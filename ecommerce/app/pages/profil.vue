<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'nuxt/app'
import { useCustomerAuthStore } from '~/stores/customerAuth'
import { useSwal } from '~/composables/useSwal'

const authStore = useCustomerAuthStore()
const router = useRouter()
const swal = useSwal()

const isLoading = ref(true)
const isSubmitting = ref(false)

const form = ref({
  nama_pelanggan: '',
  email: '',
  alamat: ''
})

onMounted(async () => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }

  form.value.nama_pelanggan = authStore.customerName
  isLoading.value = false
})

const simpanProfil = async () => {
  isSubmitting.value = true
  try {
    const response = await $fetch('http://localhost:8082/api/shop/profil', {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: {
        nama_pelanggan: form.value.nama_pelanggan,
        alamat: form.value.alamat
      }
    })

    if (response.status === 'sukses') {
      swal.sukses('Pembaruan data logistik berhasil disimpan.')
      authStore.setAuth(authStore.token, { ...authStore.customer, nama_pelanggan: form.value.nama_pelanggan })
    }
  } catch (error) {
    swal.error('Gagal memperbarui profil.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-8">Pengaturan Akun</h1>

    <div v-if="isLoading" class="animate-pulse flex gap-4">
      <div class="h-64 bg-slate-200 w-full rounded-3xl"></div>
    </div>

    <div v-else class="bg-white p-8 sm:p-10 rounded-3xl border border-slate-200 shadow-sm">
      <form @submit.prevent="simpanProfil" class="space-y-6">

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-bold text-slate-700 mb-2">Nama Lengkap</label>
            <input v-model="form.nama_pelanggan" type="text" class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none transition">
          </div>
          <div>
            <label class="block text-sm font-bold text-slate-700 mb-2">Email Akun (Terkunci)</label>
            <input v-model="form.email" type="email" disabled class="w-full px-4 py-3 bg-slate-100 border border-slate-200 rounded-xl text-slate-500 cursor-not-allowed placeholder-slate-400" placeholder="Diatur saat pendaftaran">
          </div>
        </div>

        <div>
          <label class="block text-sm font-bold text-slate-700 mb-2">Alamat Pengiriman Lengkap</label>
          <textarea v-model="form.alamat" required rows="4" class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 outline-none transition" placeholder="Masukkan nama jalan, nomor rumah, RT/RW, kelurahan, kecamatan..."></textarea>
          <p class="text-xs font-medium text-slate-500 mt-2">Alamat ini akan digunakan sebagai destinasi pengiriman pesanan Anda secara default.</p>
        </div>

        <div class="flex justify-end pt-4 border-t border-slate-100">
          <button type="submit" :disabled="isSubmitting" class="bg-indigo-600 text-white font-bold px-8 py-3.5 rounded-xl hover:bg-indigo-700 transition disabled:opacity-50 shadow-md shadow-indigo-100">
            {{ isSubmitting ? 'Menyinkronkan...' : 'Simpan Perubahan' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>



<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useToast } from 'vue-toastification'
import { useAuthStore } from '../stores/auth'
import AdminNavbar from '../components/AdminNavbar.vue'

const toast = useToast()
const authStore = useAuthStore()
const pelanggans = ref([])
const loading = ref(true)

const form = ref({
  nama_pelanggan: '',
  alamat: '',
  tipe_member: 'Silver'
})

const api = axios.create({
  baseURL: 'http://localhost:8082/api/admin'
})

api.interceptors.request.use((config) => {
  if (authStore.token) config.headers.Authorization = `Bearer ${authStore.token}`
  return config
})

const fetchPelanggan = async () => {
  loading.value = true
  try {
    const response = await api.get('/pelanggan')
    pelanggans.value = response.data.data
  } catch (err) {
    toast.error('Gagal mengambil data pelanggan.')
  } finally {
    loading.value = false
  }
}

const tambahPelanggan = async () => {
  try {
    await api.post('/pelanggan', form.value)
    toast.success('Member baru berhasil diregistrasi!')
    form.value = { nama_pelanggan: '', alamat: '', tipe_member: 'Silver' }
    fetchPelanggan()
  } catch (err) {
    toast.error('Gagal menyimpan profil pelanggan.')
  }
}

onMounted(() => {
  fetchPelanggan()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <AdminNavbar />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="mb-6">
        <h2 class="text-2xl font-black text-gray-800 tracking-tight">Manajemen Pelanggan</h2>
        <p class="text-sm text-gray-500 mt-1 font-medium">Registrasi dan kelola status keanggotaan ritel.</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">

        <div class="lg:col-span-1 bg-white p-6 rounded-2xl shadow-sm border border-gray-100 h-fit">
          <h3 class="text-md font-bold text-gray-800 mb-4 border-b pb-2">Registrasi Member Baru</h3>
          <form @submit.prevent="tambahPelanggan" class="space-y-4">
            <div>
              <label class="block text-xs font-bold text-gray-600 mb-1">Nama Lengkap</label>
              <input v-model="form.nama_pelanggan" type="text" required class="w-full p-2.5 text-sm border border-gray-200 rounded-lg focus:border-indigo-500 outline-none transition">
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-600 mb-1">Domisili / Kota</label>
              <input v-model="form.alamat" type="text" required class="w-full p-2.5 text-sm border border-gray-200 rounded-lg focus:border-indigo-500 outline-none transition">
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-600 mb-1">Tipe Awal Keanggotaan</label>
              <select v-model="form.tipe_member" class="w-full p-2.5 text-sm border border-gray-200 rounded-lg focus:border-indigo-500 outline-none transition">
                <option value="Silver">Member Silver</option>
                <option value="Gold">Member Gold</option>
              </select>
            </div>
            <button type="submit" class="w-full bg-indigo-600 text-white font-bold py-3 rounded-lg hover:bg-indigo-700 transition mt-2">Daftarkan Pelanggan</button>
          </form>
        </div>

        <div class="lg:col-span-2 bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          <div v-if="loading" class="text-center py-12 text-gray-500 font-medium animate-pulse">Sinkronisasi data profil...</div>
          <table v-else class="w-full text-left">
            <thead class="bg-gray-50 border-b border-gray-100 text-gray-500 text-xs uppercase tracking-wider">
            <tr>
              <th class="p-4 font-bold">ID</th>
              <th class="p-4 font-bold">Profil Pelanggan</th>
              <th class="p-4 font-bold">Tipe Keanggotaan</th>
              <th class="p-4 font-bold text-right">Poin Aktif</th>
            </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
            <tr v-for="p in pelanggans" :key="p.id_pelanggan" class="hover:bg-gray-50 transition">
              <td class="p-4 text-sm font-black text-gray-400">#{{ p.id_pelanggan }}</td>
              <td class="p-4">
                <div class="font-bold text-gray-800">{{ p.nama_pelanggan }}</div>
                <div class="text-xs text-gray-500 mt-0.5">{{ p.alamat }}</div>
              </td>
              <td class="p-4">
                  <span :class="p.tipe_member === 'Gold' ? 'bg-yellow-100 text-yellow-800' : (p.tipe_member === 'Silver' ? 'bg-slate-100 text-slate-700' : 'bg-gray-100 text-gray-500')" class="px-2.5 py-1 rounded-md text-xs font-black uppercase tracking-wide">
                    {{ p.tipe_member }}
                  </span>
              </td>
              <td class="p-4 text-right">
                <div class="text-lg font-black text-indigo-600">{{ p.poin_loyalitas }}</div>
              </td>
            </tr>
            </tbody>
          </table>
        </div>

      </div>
    </div>
  </div>
</template>
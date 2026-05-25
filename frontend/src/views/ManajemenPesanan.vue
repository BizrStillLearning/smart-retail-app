<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const api = axios.create({
  baseURL: 'http://localhost:8082/api',
  headers: { Authorization: `Bearer ${authStore.token}` }
})

const pesanans = ref([])
const loading = ref(true)

const fetchPesanan = async () => {
  try {
    const response = await api.get('/admin/pesanan')
    pesanans.value = response.data.data
  } catch (err) {
    alert('Gagal memuat data pesanan')
  } finally {
    loading.value = false
  }
}

const ubahStatus = async (id, statusBaru) => {
  if (confirm(`Ubah status pesanan #${id} menjadi ${statusBaru}?`)) {
    try {
      await api.put(`/admin/pesanan/${id}/status`, { status: statusBaru })
      fetchPesanan()
    } catch (err) {
      alert('Gagal merubah status')
    }
  }
}

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)
const formatTanggal = (tanggalStr) => new Date(tanggalStr).toLocaleString('id-ID')
const hitungTotal = (detailArray) => detailArray.reduce((total, item) => total + (item.kuantitas * item.produk.harga), 0)

const prosesLogout = () => {
  if (confirm('Apakah Anda yakin ingin keluar?')) {
    authStore.logout()
    router.push('/login')
  }
}

onMounted(() => fetchPesanan())
</script>

<template>
  <div class="bg-gray-50 min-h-screen font-sans">

    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center mr-8">
              <span class="text-2xl font-black text-indigo-600 tracking-tight">Smart<span class="text-gray-800">Admin</span></span>
            </div>
            <div class="hidden sm:flex sm:space-x-8">
              <router-link to="/admin" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition">Dasbor</router-link>
              <router-link to="/admin/produk" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition">Produk</router-link>
              <router-link to="/admin/pesanan" class="border-indigo-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-bold">Pesanan</router-link>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <router-link to="/" class="text-sm font-medium text-gray-500 hover:text-indigo-600 transition">Lihat Toko ↗</router-link>
            <div class="h-6 w-px bg-gray-200"></div>
            <button @click="prosesLogout" class="bg-red-50 text-red-600 px-4 py-2 rounded-lg text-sm font-bold hover:bg-red-100 transition border border-red-100">Logout</button>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Manajemen Pesanan</h1>
        <p class="text-gray-500 mt-1">Pantau dan proses pesanan pelanggan toko Anda.</p>
      </div>

      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <div v-if="loading" class="p-12 text-center text-gray-500 animate-pulse">Memuat data pesanan...</div>

        <table v-else class="w-full text-left border-collapse">
          <thead>
          <tr class="bg-gray-50 text-gray-600 border-b border-gray-200 text-sm">
            <th class="p-4 font-bold">ID / Tanggal</th>
            <th class="p-4 font-bold">Item Belanja</th>
            <th class="p-4 font-bold">Total Tagihan</th>
            <th class="p-4 font-bold">Status</th>
            <th class="p-4 font-bold text-right">Aksi</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="pesanan in pesanans" :key="pesanan.id_pesanan" class="border-b border-gray-50 hover:bg-gray-50 transition">
            <td class="p-4">
              <div class="font-black text-gray-800 text-lg">#{{ pesanan.id_pesanan }}</div>
              <div class="text-xs text-gray-500 font-medium">{{ formatTanggal(pesanan.tanggal_pesan) }}</div>
            </td>
            <td class="p-4 text-sm">
              <ul class="list-disc pl-4 text-gray-600 font-medium">
                <li v-for="item in pesanan.detail" :key="item.id_produk">{{ item.produk.nama_produk }} <span class="text-gray-400">(x{{ item.kuantitas }})</span></li>
              </ul>
            </td>
            <td class="p-4 font-black text-gray-900 text-lg">Rp {{ formatRupiah(hitungTotal(pesanan.detail)) }}</td>
            <td class="p-4">
                <span :class="{
                  'bg-yellow-100 text-yellow-800 border-yellow-200': pesanan.status === 'Menunggu',
                  'bg-blue-100 text-blue-800 border-blue-200': pesanan.status === 'Diproses',
                  'bg-green-100 text-green-800 border-green-200': pesanan.status === 'Selesai',
                  'bg-red-100 text-red-800 border-red-200': pesanan.status === 'Dibatalkan'
                }" class="px-3 py-1.5 rounded-md text-xs font-bold border">
                  {{ pesanan.status }}
                </span>
            </td>
            <td class="p-4 text-right space-x-2">
              <button v-if="pesanan.status === 'Selesai'" @click="ubahStatus(pesanan.id_pesanan, 'Diproses')" class="text-blue-600 bg-blue-50 hover:bg-blue-100 font-bold px-3 py-1.5 rounded-md text-xs transition">Set ke Proses</button>
              <button v-if="pesanan.status === 'Diproses'" @click="ubahStatus(pesanan.id_pesanan, 'Selesai')" class="text-green-600 bg-green-50 hover:bg-green-100 font-bold px-3 py-1.5 rounded-md text-xs transition">Set ke Selesai</button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>
    </main>
  </div>
</template>
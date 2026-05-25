<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { Bar } from 'vue-chartjs'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const authStore = useAuthStore()
const router = useRouter()
const laporan = ref([])
const statistik = ref({
  pesanan_selesai: 0,
  pesanan_dibatalkan: 0
})
const loading = ref(true)
const error = ref(null)

const chartData = computed(() => {
  return {
    labels: laporan.value.map(item => item.kategori),
    datasets: [{
      label: 'Pendapatan (Rp)',
      backgroundColor: '#4f46e5',
      borderRadius: 6,
      data: laporan.value.map(item => item.total_pendapatan)
    }]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: { y: { beginAtZero: true } }
}

const statCards = computed(() => {
  const totalPendapatan = laporan.value.reduce((sum, item) => sum + item.total_pendapatan, 0)

  return [
    { title: 'Total Pendapatan', value: `Rp ${formatRupiah(totalPendapatan)}`, icon: '💰', color: 'text-green-600', bg: 'bg-green-100' },
    { title: 'Total Pesanan Selesai', value: statistik.value.pesanan_selesai.toString(), icon: '📦', color: 'text-blue-600', bg: 'bg-blue-100' },
    { title: 'Pesanan Dibatalkan', value: statistik.value.pesanan_dibatalkan.toString(), icon: '❌', color: 'text-red-600', bg: 'bg-red-100' },
    { title: 'Kategori Aktif', value: laporan.value.length.toString(), icon: '🏷️', color: 'text-purple-600', bg: 'bg-purple-100' }
  ]
})

const fetchLaporan = async () => {
  try {
    const response = await axios.get('http://localhost:8082/api/laporan/top-kategori')
    laporan.value = response.data.data
    if (response.data.stats) {
      statistik.value = response.data.stats
    }
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Gagal memuat laporan'
  } finally {
    loading.value = false
  }
}

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)

const prosesLogout = () => {
  if (confirm('Apakah Anda yakin ingin keluar?')) {
    authStore.logout()
    router.push('/login')
  }
}

const unduhPDF = () => {
  window.print()
}

onMounted(() => {
  fetchLaporan()
})
</script>

<template>
  <div class="bg-gray-50 min-h-screen font-sans">

    <nav class="bg-white shadow-sm border-b border-gray-200 print:hidden">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">

          <div class="flex">
            <div class="flex-shrink-0 flex items-center mr-8">
              <span class="text-2xl font-black text-indigo-600 tracking-tight">Smart<span class="text-gray-800">Admin</span></span>
            </div>
            <div class="hidden sm:flex sm:space-x-8">
              <router-link to="/admin" class="border-indigo-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-bold">
                Dasbor
              </router-link>
              <router-link to="/admin/produk" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition">
                Produk
              </router-link>
              <router-link to="/admin/pesanan" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition">
                Pesanan
              </router-link>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <router-link to="/" class="text-sm font-medium text-gray-500 hover:text-indigo-600 transition">
              Lihat Toko ↗
            </router-link>
            <div class="h-6 w-px bg-gray-200"></div> <button @click="prosesLogout" class="bg-red-50 text-red-600 px-4 py-2 rounded-lg text-sm font-bold hover:bg-red-100 transition border border-red-100">
            Logout
          </button>
          </div>

        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">

      <div class="flex justify-between items-end mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Dasbor Analitik</h1>
          <p class="text-gray-500 mt-1">Ringkasan performa penjualan harian toko Anda.</p>
        </div>

        <button @click="unduhPDF" class="print:hidden bg-indigo-600 text-white px-5 py-2.5 rounded-lg hover:bg-indigo-700 transition font-semibold shadow-sm flex items-center gap-2 border border-indigo-700">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"></path></svg>
          Cetak PDF
        </button>
      </div>

      <div v-if="loading" class="text-gray-500 animate-pulse bg-white p-8 rounded-xl border border-gray-100 text-center">Menarik data dari database...</div>
      <div v-else-if="error" class="text-red-500 bg-red-50 p-4 rounded-xl border border-red-200">{{ error }}</div>

      <div v-else id="area-laporan" class="space-y-8">

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="(stat, index) in statCards" :key="index" class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex items-center gap-4 hover:shadow-md transition duration-200">
            <div :class="[stat.bg, stat.color, 'w-14 h-14 rounded-2xl flex items-center justify-center text-2xl shadow-inner']">
              {{ stat.icon }}
            </div>
            <div>
              <p class="text-sm text-gray-500 font-semibold mb-1">{{ stat.title }}</p>
              <h3 class="text-2xl font-black text-gray-800 tracking-tight">{{ stat.value }}</h3>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

          <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex flex-col">
            <h2 class="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
              <span class="w-1.5 h-6 bg-indigo-500 rounded-full"></span>
              Grafik Pendapatan
            </h2>
            <div class="flex-grow h-72">
              <Bar v-if="laporan.length > 0" :data="chartData" :options="chartOptions" />
              <div v-else class="h-full flex items-center justify-center text-gray-400 bg-gray-50 rounded-xl border border-dashed border-gray-200">
                Data belum tersedia.
              </div>
            </div>
          </div>

          <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
            <h2 class="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
              <span class="w-1.5 h-6 bg-green-500 rounded-full"></span>
              Rincian Top Kategori
            </h2>
            <div v-if="laporan.length === 0" class="text-gray-500 text-center py-12 bg-gray-50 rounded-xl border border-dashed border-gray-200">
              Belum ada transaksi selesai.
            </div>
            <div v-else class="overflow-x-auto rounded-xl border border-gray-100">
              <table class="w-full text-left border-collapse">
                <thead>
                <tr class="bg-gray-50 text-gray-600 text-sm border-b border-gray-200">
                  <th class="p-4 font-bold">Kategori Produk</th>
                  <th class="p-4 font-bold text-right">Total Pendapatan</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item, index) in laporan" :key="index" class="border-b border-gray-50 hover:bg-gray-50 transition">
                  <td class="p-4 font-medium text-gray-700">
                      <span class="bg-gray-100 text-gray-800 px-3 py-1 rounded-md text-sm border border-gray-200 font-semibold shadow-sm">
                        {{ item.kategori }}
                      </span>
                  </td>
                  <td class="p-4 font-black text-gray-900 text-right text-lg">
                    Rp {{ formatRupiah(item.total_pendapatan) }}
                  </td>
                </tr>
                </tbody>
              </table>
            </div>
          </div>

        </div>
      </div>
    </main>
  </div>
</template>
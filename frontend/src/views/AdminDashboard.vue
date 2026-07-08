<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { Bar } from 'vue-chartjs'
import { useAuthStore } from '../stores/auth'
import Layout from '../components/Layout.vue'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const authStore = useAuthStore()
const laporan = ref([])
const statistik = ref({
  pendapatan_kotor: 0,
  pendapatan_bersih: 0,
  total_diskon_diberikan: 0,
  pesanan_selesai: 0
})
const loading = ref(true)
const error = ref(null)

const chartData = computed(() => {
  return {
    labels: laporan.value.map(item => item.kategori),
    datasets: [{
      label: 'Pendapatan Bersih (Rp)',
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

const statCards = computed(() => [
  { title: 'Pendapatan Kotor', value: `Rp ${formatRupiah(statistik.value.pendapatan_kotor)}`, icon: '📈', color: 'text-slate-600', bg: 'bg-slate-100' },
  { title: 'Pendapatan Bersih (Real)', value: `Rp ${formatRupiah(statistik.value.pendapatan_bersih)}`, icon: '💰', color: 'text-green-600', bg: 'bg-green-100' },
  { title: 'Diskon Terpakai', value: `Rp ${formatRupiah(statistik.value.total_diskon_diberikan)}`, icon: '🎟️', color: 'text-orange-600', bg: 'bg-orange-100' },
  { title: 'Transaksi Sukses', value: statistik.value.pesanan_selesai.toString(), icon: '📦', color: 'text-blue-600', bg: 'bg-blue-100' }
])

const fetchLaporan = async () => {
  try {
    const response = await axios.get('http://localhost:8082/api/admin/dashboard', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    laporan.value = response.data.data.top_kategori
    statistik.value = response.data.data.statistik
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Gagal memuat analitik operasional'
  } finally {
    loading.value = false
  }
}

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka || 0)
const unduhPDF = () => window.print()

const loadingAI = ref(false)
const latihUlangAI = async () => {
  loadingAI.value = true
  try {
    const response = await axios.post('http://localhost:8082/api/admin/ai/sync', {}, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    alert(response.data.message)
  } catch (err) {
    alert(`Gagal menyinkronkan data ke Python Service: ${err.response?.data?.error || err.message}`)
  } finally {
    loadingAI.value = false
  }
}

onMounted(() => {
  fetchLaporan()
})
</script>

<template>
  <Layout>
    <div class="flex justify-between items-end mb-8 print:mb-4">
      <div>
        <h1 class="text-3xl font-black text-gray-900 tracking-tight">Dasbor Analitik</h1>
        <p class="text-sm font-medium text-gray-500 mt-1">Evaluasi performa penjualan dan utilisasi diskon.</p>
      </div>

      <div class="flex gap-3 print:hidden">
        <button @click="latihUlangAI" :disabled="loadingAI" class="bg-purple-50 text-purple-700 px-4 py-2 rounded-lg hover:bg-purple-100 transition font-bold text-sm shadow-sm flex items-center gap-2 border border-purple-200 disabled:opacity-50">
          <span v-if="loadingAI">⚙️ Memproses...</span>
          <span v-else>🧠 Latih AI</span>
        </button>
        <button @click="unduhPDF" class="bg-slate-900 text-white px-4 py-2 rounded-lg hover:bg-slate-800 transition font-bold text-sm shadow-sm flex items-center gap-2">
          🖨️ Cetak Laporan
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-gray-500 animate-pulse bg-white p-8 rounded-xl border border-gray-100 text-center font-bold">Mengagregasi data finansial...</div>
    <div v-else-if="error" class="text-red-500 bg-red-50 p-4 rounded-xl border border-red-200 font-bold">{{ error }}</div>

    <div v-else id="area-laporan" class="space-y-8">

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div v-for="(stat, index) in statCards" :key="index" class="bg-white p-5 rounded-2xl shadow-sm border border-gray-100 flex items-center gap-4">
          <div :class="[stat.bg, stat.color, 'w-12 h-12 rounded-xl flex items-center justify-center text-xl shadow-inner']">
            {{ stat.icon }}
          </div>
          <div>
            <p class="text-xs text-gray-500 font-bold uppercase tracking-wider mb-0.5">{{ stat.title }}</p>
            <h3 class="text-xl font-black text-gray-800">{{ stat.value }}</h3>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex flex-col">
          <h2 class="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
            <span class="w-2 h-2 bg-indigo-500 rounded-full"></span> Tren Kategori (Net)
          </h2>
          <div class="flex-grow h-72">
            <Bar v-if="laporan.length > 0" :data="chartData" :options="chartOptions" />
            <div v-else class="h-full flex items-center justify-center text-gray-400 bg-gray-50 rounded-xl border border-dashed border-gray-200 text-sm font-bold">Tidak ada data.</div>
          </div>
        </div>

        <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
          <h2 class="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
            <span class="w-2 h-2 bg-green-500 rounded-full"></span> Distribusi Pendapatan
          </h2>
          <div v-if="laporan.length === 0" class="text-gray-500 text-center py-12 bg-gray-50 rounded-xl border border-dashed border-gray-200 text-sm font-bold">Belum ada transaksi.</div>
          <div v-else class="overflow-x-auto rounded-xl border border-gray-100">
            <table class="w-full text-left">
              <thead class="bg-slate-50 border-b border-gray-200 text-slate-500 text-xs uppercase tracking-wider">
              <tr>
                <th class="p-4 font-black">Kategori</th>
                <th class="p-4 font-black text-right">Pendapatan Bersih</th>
              </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
              <tr v-for="(item, index) in laporan" :key="index" class="hover:bg-slate-50 transition">
                <td class="p-4 font-bold text-slate-700">{{ item.kategori }}</td>
                <td class="p-4 font-black text-indigo-600 text-right">Rp {{ formatRupiah(item.total_pendapatan) }}</td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

    </div>
  </Layout>
</template>
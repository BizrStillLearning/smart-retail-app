<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useToast } from 'vue-toastification'
import { useAuthStore} from "../stores/auth.js";
import Navbar from "../components/Navbar.vue";

const toast = useToast()
const authStore = useAuthStore()
const promos = ref([])
const loading = ref(true)

const form = ref({
  kode_promo: '',
  nama_promo: '',
  tipe_diskon: 'nominal',
  nilai_diskon: 0,
  minimal_belanja: 0,
  khusus_member: false
})

const api = axios.create({
  baseURL: 'http://localhost:8082/api/admin'
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

const fetchPromos = async () => {
  loading.value = true
  try {
    const response = await api.get('/promo')
    promos.value = response.data.data
  } catch (err) {
    toast.error('Gagal mengambil data master promo.')
  } finally {
    loading.value = false
  }
}

const tambahPromo = async () => {
  try {
    const payload = {
      ...form.value,
      kode_promo: form.value.kode_promo.toUpperCase().replace(/\s+/g, '')
    }

    await api.post('/promo', payload)
    toast.success('Kupon promo berhasil diterbitkan!')

    form.value = { kode_promo: '', nama_promo: '', tipe_diskon: 'nominal', nilai_diskon: 0, minimal_belanja: 0, khusus_member: false }
    fetchPromos()
  } catch (err) {
    const pesanError = err.response?.data?.error || 'Terjadi kesalahan sistem'
    toast.error(`Gagal menyimpan: ${pesanError}`)
  }
}

const toggleStatus = async (id, statusSaatIni) => {
  const statusBaru = statusSaatIni === 'aktif' ? 'nonaktif' : 'aktif'
  try {
    await api.put(`/promo/${id}/status`, { status: statusBaru })
    toast.info(`Status promo diubah menjadi ${statusBaru.toUpperCase()}`)
    fetchPromos()
  } catch (err) {
    toast.error('Gagal mengubah status operasional promo.')
  }
}

const hapusPromo = async (id, kode) => {
  if (!confirm(`Tindakan ini permanen! Yakin ingin menghapus kupon ${kode}?`)) return

  try {
    await api.delete(`/promo/${id}`)
    toast.success('Data promo berhasil dihapus dari sistem.')
    fetchPromos()
  } catch (err) {
    toast.error('Penghapusan ditolak oleh server.')
  }
}

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)

onMounted(() => {
  fetchPromos()
})
</script>

<template>
  <Navbar />

  <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
    <div class="mb-6 border-b border-gray-100 pb-4">
      <h2 class="text-2xl font-black text-gray-800 tracking-tight flex items-center gap-2">
        <span class="text-indigo-600">🎟️</span> Manajemen Promo & Diskon
      </h2>
      <p class="text-sm text-gray-500 mt-1 font-medium">Buat dan kendalikan aturan potongan harga untuk transaksi ritel.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">

      <div class="lg:col-span-1 bg-gray-50 p-5 rounded-xl border border-gray-200 h-fit">
        <h3 class="text-md font-bold text-gray-700 mb-4 border-b pb-2">Buat Kupon Baru</h3>

        <form @submit.prevent="tambahPromo" class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-gray-600 mb-1">Kode Kupon (Unik)</label>
            <input v-model="form.kode_promo" type="text" required placeholder="Cth: SEMBAKOMURAH" class="w-full p-2.5 text-sm border border-gray-300 rounded-lg uppercase focus:ring-2 focus:ring-indigo-500 outline-none">
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-600 mb-1">Nama / Deskripsi Promo</label>
            <input v-model="form.nama_promo" type="text" required placeholder="Cth: Diskon Akhir Bulan" class="w-full p-2.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 outline-none">
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-gray-600 mb-1">Tipe Diskon</label>
              <select v-model="form.tipe_diskon" class="w-full p-2.5 text-sm border border-gray-300 rounded-lg outline-none">
                <option value="nominal">Nominal (Rp)</option>
                <option value="persentase">Persentase (%)</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-600 mb-1">Nilai Potongan</label>
              <input v-model="form.nilai_diskon" type="number" min="1" required class="w-full p-2.5 text-sm border border-gray-300 rounded-lg outline-none">
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-600 mb-1">Minimal Belanja (Rp)</label>
            <input v-model="form.minimal_belanja" type="number" min="0" required class="w-full p-2.5 text-sm border border-gray-300 rounded-lg outline-none">
          </div>

          <div class="flex items-center gap-2 bg-white p-3 border border-gray-200 rounded-lg">
            <input v-model="form.khusus_member" type="checkbox" id="memberOnly" class="w-4 h-4 text-indigo-600 rounded">
            <label for="memberOnly" class="text-sm font-bold text-gray-700 cursor-pointer">Eksklusif Member (Non-Guest)</label>
          </div>

          <button type="submit" class="w-full bg-indigo-600 text-white font-bold py-3 rounded-lg hover:bg-indigo-700 transition shadow-sm mt-2">
            Terbitkan Promo
          </button>
        </form>
      </div>

      <!-- TABEL DAFTAR PROMO -->
      <div class="lg:col-span-2">
        <div v-if="loading" class="text-center py-10 text-gray-500 font-medium animate-pulse">Memuat data operasional...</div>

        <div v-else class="overflow-x-auto border border-gray-200 rounded-xl">
          <table class="w-full text-left border-collapse">
            <thead>
            <tr class="bg-gray-100 text-gray-600 text-xs uppercase tracking-wider">
              <th class="p-4 font-bold border-b">Kode & Detail</th>
              <th class="p-4 font-bold border-b">Aturan Diskon</th>
              <th class="p-4 font-bold border-b text-center">Status</th>
              <th class="p-4 font-bold border-b text-center">Aksi</th>
            </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
            <tr v-if="promos.length === 0">
              <td colspan="4" class="p-8 text-center text-gray-400 font-medium">Belum ada promo yang terdaftar.</td>
            </tr>
            <tr v-for="p in promos" :key="p.id_promo" class="hover:bg-gray-50 transition">
              <td class="p-4">
                <div class="font-black text-indigo-700 text-sm tracking-wide">{{ p.kode_promo }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ p.nama_promo }}</div>
                <div v-if="p.khusus_member" class="inline-block mt-1 bg-yellow-100 text-yellow-800 text-[10px] px-2 py-0.5 rounded font-bold">Khusus Member</div>
              </td>
              <td class="p-4">
                <div class="text-sm font-bold text-gray-800">
                  {{ p.tipe_diskon === 'persentase' ? `${p.nilai_diskon}%` : `Rp ${formatRupiah(p.nilai_diskon)}` }}
                </div>
                <div class="text-xs text-gray-500 mt-0.5">Min: Rp {{ formatRupiah(p.minimal_belanja) }}</div>
              </td>
              <td class="p-4 text-center">
                <button
                    @click="toggleStatus(p.id_promo, p.status)"
                    :class="p.status === 'aktif' ? 'bg-green-100 text-green-700 border-green-200' : 'bg-red-100 text-red-700 border-red-200'"
                    class="px-3 py-1 rounded-full text-xs font-bold border transition hover:opacity-80 uppercase tracking-wider"
                    title="Klik untuk mengubah status"
                >
                  {{ p.status }}
                </button>
              </td>
              <td class="p-4 text-center">
                <button @click="hapusPromo(p.id_promo, p.kode_promo)" class="text-red-500 hover:text-red-700 hover:bg-red-50 p-2 rounded-lg transition" title="Hapus Permanen">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                </button>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  </div>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'nuxt/app'
import { useCustomerAuthStore } from '~/stores/customerAuth'
import { useSwal } from '~/composables/useSwal'

const authStore = useCustomerAuthStore()
const router = useRouter()
const swal = useSwal()

const daftarRiwayat = ref([])
const isLoading = ref(true)
const fileInput = ref(null)
const idPesananAktif = ref(null)

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)
const formatTanggal = (tanggal) => new Date(tanggal).toLocaleDateString('id-ID', {
  year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit'
})

const muatRiwayat = async () => {
  try {
    const response = await $fetch('http://localhost:8082/api/shop/riwayat', {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (response.status === 'sukses') daftarRiwayat.value = response.data
  } catch (error) {
    console.error("Gagal memuat riwayat:", error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  if (!authStore.isLoggedIn) {
    router.push('/login')
    return
  }
  muatRiwayat()
})

const triggerUpload = (id) => {
  idPesananAktif.value = id
  fileInput.value.click()
}

const prosesUnggahBerkas = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('bukti_transfer', file)
  formData.append('id_pesanan', idPesananAktif.value)

  try {
    swal.info('Mengunggah berkas...')
    const response = await $fetch('http://localhost:8082/api/shop/upload-bukti', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: formData
    })

    if (response.status === 'sukses') {
      swal.sukses('Bukti transfer berhasil dikirim ke Admin.')
      muatRiwayat()
    }
  } catch (error) {
    swal.error('Gagal mengunggah gambar. Pastikan format valid (JPG/PNG).')
  } finally {
    event.target.value = null
  }
}

const getStatusWarna = (status) => {
  switch (status.toLowerCase()) {
    case 'menunggu pembayaran': return 'bg-amber-100 text-amber-700 border-amber-200'
    case 'verifikasi pembayaran': return 'bg-orange-100 text-orange-700 border-orange-200'
    case 'diproses': return 'bg-blue-100 text-blue-700 border-blue-200'
    case 'dikirim': return 'bg-indigo-100 text-indigo-700 border-indigo-200'
    case 'selesai': return 'bg-emerald-100 text-emerald-700 border-emerald-200'
    default: return 'bg-slate-100 text-slate-700 border-slate-200'
  }
}
</script>

<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-3xl font-black text-slate-900 tracking-tight">Riwayat Pesanan</h1>
      <NuxtLink to="/" class="text-sm font-bold text-indigo-600 hover:text-indigo-800 transition">← Belanja Lagi</NuxtLink>
    </div>

    <input type="file" ref="fileInput" @change="prosesUnggahBerkas" accept="image/*" class="hidden">

    <div v-if="isLoading" class="text-center py-12 text-slate-500 font-bold animate-pulse">Memuat data transaksi...</div>
    <div v-else-if="daftarRiwayat.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 text-center shadow-sm">
      <div class="text-5xl mb-4">🧾</div>
      <h2 class="text-xl font-bold text-slate-700 mb-2">Belum Ada Transaksi</h2>
    </div>

    <div v-else class="space-y-6">
      <div v-for="pesanan in daftarRiwayat" :key="pesanan.id_pesanan" class="bg-white rounded-3xl border border-slate-200 overflow-hidden shadow-sm hover:shadow-md transition">

        <div class="bg-slate-50 px-6 py-4 border-b border-slate-100 flex flex-wrap items-center justify-between gap-4">
          <div>
            <div class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-1">ID Pesanan: #{{ pesanan.id_pesanan }}</div>
            <div class="text-sm font-bold text-slate-800">{{ formatTanggal(pesanan.tanggal_pesan) }}</div>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-xs font-bold px-3 py-1 rounded-full border" :class="getStatusWarna(pesanan.status)">
              {{ pesanan.status.toUpperCase() }}
            </span>
          </div>
        </div>

        <div class="px-6 py-6 flex flex-col md:flex-row justify-between items-center gap-6">
          <div class="w-full md:w-auto">
            <div class="text-sm text-slate-500 font-medium mb-1">Metode Pembayaran</div>
            <div class="text-base font-bold text-slate-800">{{ pesanan.metode_pembayaran }}</div>
            <div class="text-sm text-slate-500 font-medium mt-3 mb-1">Total Transaksi</div>
            <div class="text-2xl font-black text-indigo-600">Rp {{ formatRupiah(pesanan.total_akhir) }}</div>
          </div>

          <div class="w-full md:w-auto flex flex-col gap-3">
            <button v-if="pesanan.status === 'menunggu pembayaran' && pesanan.metode_pembayaran === 'Transfer Bank'"
                    @click="triggerUpload(pesanan.id_pesanan)"
                    class="bg-amber-500 text-white font-bold px-6 py-3 rounded-xl hover:bg-amber-600 transition shadow-md shadow-amber-100">
              Unggah Bukti Transfer
            </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>


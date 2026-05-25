<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useKeranjangStore } from '../stores/keranjang'

const keranjangStore = useKeranjangStore()
const produks = ref([])
const loading = ref(true)
const error = ref(null)

const api = axios.create({
  baseURL: 'http://localhost:8082/api'
})

const pesananBerhasil = ref(false)
const dataNota = ref({
  idPesanan: null,
  total: 0,
  items: []
})

const fetchProduk = async () => {
  try {
    const response = await api.get('/produk')
    produks.value = response.data.data
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Gagal mengambil data produk'
  } finally {
    loading.value = false
  }
}

const prosesCheckout = async () => {
  if (keranjangStore.items.length === 0) return

  try {
    const snapshotKeranjang = [...keranjangStore.items]
    const snapshotTotal = keranjangStore.totalHarga

    const payload = {
      items: keranjangStore.items.map(item => ({
        id_produk: item.id_produk,
        kuantitas: item.kuantitas
      }))
    }

    const response = await api.post('/pesanan', payload)

    dataNota.value = {
      idPesanan: response.data.id_pesanan,
      total: snapshotTotal,
      items: snapshotKeranjang
    }
    pesananBerhasil.value = true

    keranjangStore.kosongkanKeranjang()
  } catch (err) {
    const pesanError = err.response?.data?.error || err.message
    alert(`Terjadi kesalahan: ${pesanError}`)
  }
}

const batalkanPesanan = async () => {
  try {
    await api.put(`/pesanan/${dataNota.value.idPesanan}/batal`)
    alert('Pesanan berhasil dibatalkan!')
    pesananBerhasil.value = false
  } catch (err) {
    const pesanError = err.response?.data?.error || err.message
    alert(`Gagal membatalkan pesanan: ${pesanError}`)
  }
}

const formatRupiah = (angka) => {
  return new Intl.NumberFormat('id-ID').format(angka)
}

const cetakStruk = () => {
  window.print()
}

onMounted(() => {
  fetchProduk()
})
</script>

<template>
  <div class="p-8 bg-gray-50 min-h-screen">
    <header class="flex justify-between items-center mb-8 pb-4 border-b border-gray-200 print:hidden">
      <div>
        <h1 class="text-3xl font-black text-indigo-600 tracking-tight">Smart<span class="text-gray-800">Retail</span></h1>
        <p class="text-sm text-gray-500 font-medium">Sistem Pemesanan Cepat & Terintegrasi</p>
      </div>
      <router-link to="/login" class="text-sm font-bold text-gray-500 hover:text-indigo-600 flex items-center gap-1 transition">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path></svg>
        Portal Admin
      </router-link>
    </header>
    <div v-if="!pesananBerhasil">
      <h1 class="text-3xl font-bold text-gray-800 mb-8">Sistem Smart Retail</h1>

      <div class="flex flex-col lg:flex-row gap-8">

        <div class="lg:w-2/3">
          <h2 class="text-xl font-semibold mb-4 text-gray-700">Katalog Produk</h2>

          <div v-if="loading" class="text-gray-500 animate-pulse">Memuat data...</div>
          <div v-else-if="error" class="text-red-500">{{ error }}</div>

          <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
            <div v-for="produk in produks" :key="produk.id_produk" class="bg-white p-5 rounded-xl shadow-sm border border-gray-100 flex flex-col justify-between">
              <div>
                <img v-if="produk.gambar" :src="`http://localhost:8082${produk.gambar}`" class="w-full h-40 object-cover rounded-lg mb-4 border border-gray-100" alt="Produk" />
                <div v-else class="w-full h-40 bg-gray-100 rounded-lg mb-4 flex items-center justify-center text-gray-400 text-sm">No Image</div>

                <span class="text-xs font-semibold text-blue-600 bg-blue-50 px-2 py-1 rounded-full">{{ produk.kategori }}</span>
                <h3 class="text-lg font-bold text-gray-800 mt-3">{{ produk.nama_produk }}</h3>
                <div class="text-xl font-black text-gray-900 mt-1">Rp {{ formatRupiah(produk.harga) }}</div>
              </div>
              <button @click="keranjangStore.tambahItem(produk)" class="mt-4 w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition">
                Tambah
              </button>
            </div>
          </div>
        </div>

        <div class="lg:w-1/3">
          <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 sticky top-8">
            <h2 class="text-xl font-semibold mb-4 text-gray-700">Keranjang Belanja</h2>

            <div v-if="keranjangStore.items.length === 0" class="text-center text-gray-400 py-8">Belum ada produk di keranjang.</div>

            <div v-else class="space-y-4 max-h-96 overflow-y-auto pr-2">
              <div v-for="item in keranjangStore.items" :key="item.id_produk" class="flex justify-between items-center border-b pb-3">
                <div>
                  <div class="font-semibold text-gray-800">{{ item.nama_produk }}</div>
                  <div class="text-sm text-gray-500">{{ item.kuantitas }} x Rp {{ formatRupiah(item.harga) }}</div>
                </div>
                <button @click="keranjangStore.hapusItem(item.id_produk)" class="text-red-500 hover:text-red-700 text-sm font-bold p-2">✕</button>
              </div>
            </div>

            <div class="mt-6 pt-4 border-t-2 border-gray-100">
              <div class="flex justify-between items-center mb-6">
                <span class="text-gray-600 font-medium">Total Harga:</span>
                <span class="text-2xl font-black text-gray-900">Rp {{ formatRupiah(keranjangStore.totalHarga) }}</span>
              </div>
              <button @click="prosesCheckout" :disabled="keranjangStore.items.length === 0" class="w-full bg-green-600 text-white py-3 rounded-lg font-bold hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition">
                Proses Checkout
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="max-w-2xl mx-auto bg-white p-8 rounded-xl shadow-lg border border-gray-200 mt-10 print:shadow-none print:border-none print:mt-0 print:p-0 print:bg-transparent">

      <div id="area-struk">
        <div class="text-center mb-8">
          <div class="text-green-500 text-5xl mb-4 print:hidden">✓</div>
          <h2 class="text-3xl font-bold text-gray-800">Smart Retail</h2>
          <p class="text-gray-500 mt-2">ID Pesanan: <span class="font-bold text-gray-800">#{{ dataNota.idPesanan }}</span></p>
        </div>

        <div class="border-t border-b py-4 mb-6 space-y-3 border-gray-300">
          <div v-for="item in dataNota.items" :key="item.id_produk" class="flex justify-between text-gray-800">
            <span>{{ item.nama_produk }} <span class="text-gray-500">(x{{ item.kuantitas }})</span></span>
            <span class="font-medium">Rp {{ formatRupiah(item.harga * item.kuantitas) }}</span>
          </div>
        </div>

        <div class="flex justify-between items-center text-xl font-bold mb-8 text-gray-900">
          <span>Total Dibayar</span>
          <span>Rp {{ formatRupiah(dataNota.total) }}</span>
        </div>

        <div class="text-center text-sm text-gray-500 hidden print:block mt-8">
          Terima kasih telah berbelanja di Smart Retail!<br>
          Barang yang sudah dibeli tidak dapat ditukar/dikembalikan.
        </div>
      </div>

      <div class="flex gap-4 print:hidden">
        <button @click="pesananBerhasil = false" class="w-1/3 bg-gray-200 text-gray-800 py-3 rounded-lg font-semibold hover:bg-gray-300 transition">Belanja Lagi</button>
        <button @click="batalkanPesanan" class="w-1/3 bg-red-50 text-red-600 py-3 rounded-lg font-semibold border border-red-200 hover:bg-red-100 transition">Batalkan Pesanan</button>
        <button @click="cetakStruk" class="w-1/3 bg-indigo-600 text-white py-3 rounded-lg font-semibold shadow-md hover:bg-indigo-700 transition flex items-center justify-center gap-2">
          🖨️ Cetak Struk
        </button>
      </div>

    </div>

  </div>
</template>
<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import axios from 'axios'
import { useKeranjangStore } from '../stores/keranjang.js'
import { useToast } from 'vue-toastification'

import PromoCarousel from '../components/PromoCarousel.vue'
import KategoriTabs from '../components/KategoriTabs.vue'
import ScarcitySection from '../components/ScarcitySection.vue'
import TrendingProducts from '../components/TrendingProducts.vue'

const keranjangStore = useKeranjangStore()
const toast = useToast()
const produks = ref([])
const loading = ref(true)
const error = ref(null)

const dataMemberAktif = ref(null)
const loadingMember = ref(false)
const errorMember = ref('')
let timeoutCariMember = null

const api = axios.create({
  baseURL: 'http://localhost:8082/api'
})

const kategoriAktif = ref('Semua')

const idPelangganInput = ref('')
const kodePromoInput = ref('')

const pesananBerhasil = ref(false)
const dataNota = ref({
  idPesanan: null,
  subtotal: 0,
  totalDiskon: 0,
  totalAkhir: 0,
  items: []
})

const daftarKategori = computed(() => {
  const list = new Set(produks.value.map(p => p.kategori))
  return ['Semua', ...list]
})

const produkTerfilter = computed(() => {
  if (kategoriAktif.value === 'Semua') return produks.value
  return produks.value.filter(p => p.kategori === kategoriAktif.value)
})

const fetchProduk = async () => {
  try {
    const response = await api.get('/produk')
    produks.value = response.data.data
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Gagal mengambil data produk'
    toast.error('Gagal memuat katalog produk.')
  } finally {
    loading.value = false
  }
}

const cekKuantitasDiKeranjang = (id_produk) => {
  const item = keranjangStore.items.find(i => i.id_produk === id_produk)
  return item ? item.kuantitas : 0
}

const tambahKeKeranjang = (produk) => {
  const qtySaatIni = cekKuantitasDiKeranjang(produk.id_produk)

  if (qtySaatIni >= produk.stok) {
    toast.warning(`Maksimal! Sisa stok ${produk.nama_produk} hanya ${produk.stok}.`)
    return
  }

  keranjangStore.tambahItem(produk)
  toast.success(`${produk.nama_produk} masuk ke keranjang!`, { timeout: 1500 })
}

const prosesCheckout = async () => {
  if (keranjangStore.items.length === 0) return

  try {
    const snapshotKeranjang = [...keranjangStore.items]

    const payload = {
      id_pelanggan: parseInt(idPelangganInput.value) || 1,
      kode_promo: kodePromoInput.value.toUpperCase().trim(),
      items: keranjangStore.items.map(item => ({
        id_produk: item.id_produk,
        kuantitas: item.kuantitas
      }))
    }

    const response = await api.post('/pesanan', payload)

    dataNota.value = {
      idPesanan: response.data.id_pesanan,
      subtotal: response.data.subtotal,
      totalDiskon: response.data.total_diskon,
      totalAkhir: response.data.total_akhir,
      items: snapshotKeranjang
    }

    pesananBerhasil.value = true
    keranjangStore.kosongkanKeranjang()
    idPelangganInput.value = ''
    kodePromoInput.value = ''

    toast.success('Transaksi berhasil diproses!')
    fetchProduk()
  } catch (err) {
    const pesanError = err.response?.data?.error || err.message
    toast.error(`Gagal Checkout: ${pesanError}`)
  }
}

const batalkanPesanan = async () => {
  try {
    await api.put(`/pesanan/${dataNota.value.idPesanan}/batal`)
    toast.success('Pesanan berhasil dibatalkan dan stok dikembalikan.')
    pesananBerhasil.value = false
    fetchProduk()
  } catch (err) {
    const pesanError = err.response?.data?.error || err.message
    toast.error(`Gagal membatalkan: ${pesanError}`)
  }
}

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)

const rekomendasiAI = ref([])
const loadingRekomendasi = ref(false)

const fetchRekomendasi = async (id_produk) => {
  loadingRekomendasi.value = true
  try {
    const response = await api.get(`/produk/${id_produk}/rekomendasi`)
    rekomendasiAI.value = response.data.data
  } catch (err) {
    toast.warning('Layanan rekomendasi AI sedang tidak tersedia.')
  } finally {
    loadingRekomendasi.value = false
  }
}

const dataStruk = ref(null)

const prosesCheckoutSukses = (responsDariBackend) => {
  dataStruk.value = {
    no_pesanan: responsDariBackend.id_pesanan || `TRX-${Math.floor(Math.random() * 100000)}`,
    waktu: new Date().toLocaleString('id-ID'),
    kasir: authStore.namaLengkap || 'Admin Toko',
    member: dataMemberAktif.value ? dataMemberAktif.value.nama_pelanggan : 'Guest (Umum)',
    items: [...keranjang.value],
    subtotal: totalHarga.value,
    diskon: nilaiDiskon.value,
    totalAkhir: totalSetelahDiskon.value,
    metode_pembayaran: 'Tunai'
  }

  keranjang.value = []
  kodePromo.value = ''
  idPelangganInput.value = ''
  dataMemberAktif.value = null
  diskonAktif.value = null
}

const cetakStruk = () => {
  window.print()
}

const tutupStruk = () => {
  dataStruk.value = null
}

watch(() => keranjangStore.items, (newItems) => {
  if (newItems.length > 0) {
    const itemTerakhir = newItems[newItems.length - 1]
    fetchRekomendasi(itemTerakhir.id_produk)
  } else {
    rekomendasiAI.value = []
  }
}, { deep: true })

watch(idPelangganInput, (newVal) => {
  clearTimeout(timeoutCariMember)
  dataMemberAktif.value = null
  errorMember.value = ''

  if (!newVal || newVal == 1) return

  loadingMember.value = true

  timeoutCariMember = setTimeout(async () => {
    try {
      const response = await axios.get(`http://localhost:8082/api/pelanggan/${newVal}`)
      dataMemberAktif.value = response.data.data
    } catch (err) {
      errorMember.value = 'Data Member tidak ditemukan.'
    } finally {
      loadingMember.value = false
    }
  }, 500)
})

onMounted(() => {
  fetchProduk()
})
</script>

<template>
  <div class="p-8 bg-gray-50 min-h-screen font-sans">

    <header class="flex justify-between items-center mb-8 pb-4 border-b border-gray-200 print:hidden">
      <div>
        <h1 class="text-3xl font-black text-indigo-600 tracking-tight">Smart<span class="text-gray-800">Retail</span></h1>
        <p class="text-sm text-gray-500 font-medium">Sistem Pemesanan Cepat & Terintegrasi</p>
      </div>
      <router-link to="/login" class="text-sm font-bold text-gray-500 hover:text-indigo-600 flex items-center gap-1 transition bg-white px-4 py-2 rounded-lg shadow-sm border border-gray-100">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path></svg>
        Portal Admin
      </router-link>
    </header>

    <div v-if="!pesananBerhasil">
      <PromoCarousel class="print:hidden" />

      <ScarcitySection :produks="produks" @tambah-ke-keranjang="tambahKeKeranjang" class="print:hidden" />

      <div class="flex flex-col lg:flex-row gap-8">

        <div class="lg:w-2/3">
          <div class="flex flex-col sm:flex-row justify-between sm:items-center gap-4 mb-4">
            <h2 class="text-xl font-bold text-gray-800 flex items-center gap-2">
              <span class="w-1.5 h-5 bg-indigo-500 rounded-full"></span>
              Etalase Produk
            </h2>

            <KategoriTabs :categories="daftarKategori" v-model="kategoriAktif" />
          </div>

          <div v-if="loading" class="text-gray-500 animate-pulse bg-white p-8 rounded-xl border border-gray-100 text-center font-medium">Memuat etalase...</div>
          <div v-else-if="error" class="text-red-500 bg-red-50 p-4 rounded-xl border border-red-200">{{ error }}</div>

          <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
            <div v-for="produk in produkTerfilter" :key="produk.id_produk" class="bg-white p-5 rounded-2xl shadow-sm border border-gray-100 flex flex-col justify-between hover:shadow-md transition">
              <div>
                <img v-if="produk.gambar" :src="`http://localhost:8082${produk.gambar}`" class="w-full h-40 object-cover rounded-xl mb-4 border border-gray-100" alt="Produk" />
                <div v-else class="w-full h-40 bg-gray-50 rounded-xl mb-4 flex items-center justify-center text-gray-400 text-sm border border-dashed border-gray-200">No Image</div>

                <div class="flex justify-between items-center mb-2">
                  <span class="text-xs font-bold text-blue-700 bg-blue-50 px-2.5 py-1 rounded-md border border-blue-100">{{ produk.kategori }}</span>
                  <span :class="{
                    'text-red-700 bg-red-50 border-red-100': produk.stok === 0,
                    'text-orange-700 bg-orange-50 border-orange-100': produk.stok > 0 && produk.stok <= 5,
                    'text-green-700 bg-green-50 border-green-100': produk.stok > 5
                  }" class="text-[10px] font-black px-2 py-1 rounded-md border tracking-wide uppercase">
                    {{ produk.stok === 0 ? 'Habis' : `Stok: ${produk.stok}` }}
                  </span>
                </div>

                <h3 class="text-lg font-bold text-gray-800 mt-1 leading-tight">{{ produk.nama_produk }}</h3>
                <div class="text-xl font-black text-gray-900 mt-2">Rp {{ formatRupiah(produk.harga) }}</div>
              </div>

              <button
                  @click="tambahKeKeranjang(produk)"
                  :disabled="produk.stok === 0 || cekKuantitasDiKeranjang(produk.id_produk) >= produk.stok"
                  :class="{
                  'bg-gray-100 text-gray-400 cursor-not-allowed border border-gray-200': produk.stok === 0 || cekKuantitasDiKeranjang(produk.id_produk) >= produk.stok,
                  'bg-indigo-600 text-white hover:bg-indigo-700 shadow-sm': produk.stok > 0 && cekKuantitasDiKeranjang(produk.id_produk) < produk.stok
                }"
                  class="mt-5 w-full py-2.5 rounded-xl font-bold transition flex justify-center items-center gap-2"
              >
                <span v-if="produk.stok === 0">Stok Habis</span>
                <span v-else-if="cekKuantitasDiKeranjang(produk.id_produk) >= produk.stok">Maksimal Stok</span>
                <span v-else>+ Keranjang</span>
              </button>
            </div>
          </div>
        </div>

        <div class="lg:w-1/3">
          <div class="sticky top-8">

            <TrendingProducts :produks="produks" class="print:hidden" />

            <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
              <h2 class="text-xl font-bold mb-6 text-gray-800 flex items-center gap-2">
                <span class="w-1.5 h-5 bg-green-500 rounded-full"></span>
                Keranjang Belanja
              </h2>

              <div v-if="keranjangStore.items.length === 0" class="text-center text-gray-400 py-12 bg-gray-50 rounded-xl border border-dashed border-gray-200 font-medium">Belum ada pesanan.</div>

              <div v-else class="space-y-4 max-h-64 overflow-y-auto pr-2 custom-scrollbar">
                <div v-for="item in keranjangStore.items" :key="item.id_produk" class="flex justify-between items-center border-b border-gray-100 pb-4">
                  <div>
                    <div class="font-bold text-gray-800">{{ item.nama_produk }}</div>
                    <div class="text-sm font-medium text-gray-500 mt-0.5">{{ item.kuantitas }} x Rp {{ formatRupiah(item.harga) }}</div>
                  </div>
                  <button @click="keranjangStore.hapusItem(item.id_produk)" class="text-red-500 hover:text-red-700 hover:bg-red-50 p-2 rounded-lg font-bold transition">✕</button>
                </div>
              </div>

              <div class="mt-4 space-y-3 bg-gray-50 p-4 rounded-xl border border-gray-100">
                <div>
                  <label class="block text-xs font-bold text-gray-600 mb-1">ID Member (Opsional)</label>
                  <input
                      v-model="idPelangganInput"
                      type="number"
                      min="1"
                      placeholder="Contoh: 2 (Biarkan kosong untuk Guest)"
                      class="w-full py-2 px-3 text-sm text-gray-800 outline-none bg-white border border-gray-200 rounded-lg focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
                  >

                  <div v-if="loadingMember" class="text-xs text-indigo-500 font-bold mt-2 animate-pulse flex items-center gap-1">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                    Mencari data ke sistem...
                  </div>
                  <div v-else-if="errorMember" class="text-xs text-red-500 font-bold mt-2 flex items-center gap-1">
                    ⚠️ {{ errorMember }}
                  </div>

                  <div v-else-if="dataMemberAktif && dataMemberAktif.id_pelanggan !== 1" class="mt-3 bg-gradient-to-br from-slate-800 to-slate-900 p-3.5 rounded-xl shadow-md border border-slate-700 flex justify-between items-center text-white transform transition-all">
                    <div class="flex items-center gap-3">
                      <div class="w-8 h-8 rounded-full bg-slate-700 flex items-center justify-center text-slate-300 font-black">
                        {{ dataMemberAktif.nama_pelanggan.charAt(0) }}
                      </div>
                      <div>
                        <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider mb-0.5">Pelanggan Valid</div>
                        <div class="text-sm font-black leading-none">{{ dataMemberAktif.nama_pelanggan }}</div>
                      </div>
                    </div>
                    <div class="text-right">
                      <div :class="dataMemberAktif.tipe_member === 'Gold' ? 'text-yellow-400' : 'text-slate-300'" class="text-xs font-black uppercase tracking-widest">{{ dataMemberAktif.tipe_member }}</div>
                      <div class="text-xs font-bold text-slate-400 mt-0.5">{{ dataMemberAktif.poin_loyalitas }} Poin</div>
                    </div>
                  </div>
                </div>


                <div>
                  <label class="block text-xs font-bold text-gray-600 mb-1">Kode Kupon / Promo</label>
                  <input
                      v-model="kodePromoInput"
                      type="text"
                      placeholder="Masukkan kode promo..."
                      class="w-full py-2 px-3 text-sm text-gray-800 outline-none bg-white border border-gray-200 rounded-lg uppercase placeholder:normal-case focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
                  >
                </div>
              </div>

              <div class="mt-4 pt-4 border-t-2 border-dashed border-gray-200">
                <div class="flex justify-between items-center mb-6">
                  <span class="text-gray-500 font-bold">Estimasi Subtotal:</span>
                  <span class="text-2xl font-black text-gray-900">Rp {{ formatRupiah(keranjangStore.totalHarga) }}</span>
                </div>
                <button @click="prosesCheckout" :disabled="keranjangStore.items.length === 0" class="w-full bg-green-600 text-white py-3.5 rounded-xl font-black hover:bg-green-700 disabled:bg-gray-200 disabled:text-gray-400 disabled:cursor-not-allowed transition shadow-sm">
                  Proses Pembayaran
                </button>
              </div>

              <div v-if="rekomendasiAI.length > 0" class="bg-gradient-to-br from-indigo-50 to-blue-50 p-6 rounded-2xl shadow-sm border border-indigo-100 mt-6 transition-all">
                <h2 class="text-sm font-black mb-4 text-indigo-800 flex items-center gap-2 uppercase tracking-wide">
                  <span class="text-xl">✨</span> Sering Dibeli Bersamaan
                </h2>

                <div class="space-y-3">
                  <div v-for="rek in rekomendasiAI" :key="rek.id_produk" class="bg-white p-3 rounded-xl border border-indigo-50 flex justify-between items-center shadow-sm">
                    <div class="flex items-center gap-3">
                      <img v-if="rek.gambar" :src="`http://localhost:8082${rek.gambar}`" class="w-10 h-10 object-cover rounded-lg border border-gray-100" />
                      <div v-else class="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center text-[10px] text-gray-400">No Img</div>

                      <div>
                        <div class="font-bold text-gray-800 text-sm">{{ rek.nama_produk }}</div>
                        <div class="text-xs font-medium text-indigo-600">Rp {{ formatRupiah(rek.harga) }}</div>
                      </div>
                    </div>

                    <button
                        @click="tambahKeKeranjang(rek)"
                        :disabled="rek.stok === 0 || cekKuantitasDiKeranjang(rek.id_produk) >= rek.stok"
                        class="bg-indigo-100 text-indigo-700 hover:bg-indigo-600 hover:text-white p-2 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed"
                        title="Tambah ke Keranjang"
                    >
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4"></path></svg>
                    </button>
                  </div>
                </div>
              </div>

            </div>
          </div>
        </div>

      </div>
    </div>

    <div v-else class="max-w-2xl mx-auto bg-white p-8 rounded-2xl shadow-xl border border-gray-100 mt-10 print:shadow-none print:border-none print:mt-0 print:p-0 print:bg-transparent">
      <div id="area-struk">
        <div class="text-center mb-8">
          <div class="w-20 h-20 bg-green-100 text-green-600 rounded-full flex items-center justify-center text-4xl mx-auto mb-4 print:hidden">✓</div>
          <h2 class="text-3xl font-black text-gray-800 tracking-tight">SmartRetail</h2>
          <p class="text-gray-500 mt-2 font-medium">ID Pesanan: <span class="font-bold text-gray-900">#{{ dataNota.idPesanan }}</span></p>
        </div>

        <div class="border-t-2 border-dashed border-gray-200 pt-6 pb-2 space-y-4">
          <div v-for="item in dataNota.items" :key="item.id_produk" class="flex justify-between text-gray-800 font-medium">
            <span>{{ item.nama_produk }} <span class="text-gray-400 ml-1">x{{ item.kuantitas }}</span></span>
            <span class="font-bold">Rp {{ formatRupiah(item.harga * item.kuantitas) }}</span>
          </div>
        </div>

        <div class="border-t-2 border-b-2 border-gray-200 py-4 mb-6 space-y-2">
          <div class="flex justify-between text-gray-600 font-bold">
            <span>Subtotal</span>
            <span>Rp {{ formatRupiah(dataNota.subtotal) }}</span>
          </div>
          <div v-if="dataNota.totalDiskon > 0" class="flex justify-between text-red-500 font-bold">
            <span>Diskon Promo</span>
            <span>- Rp {{ formatRupiah(dataNota.totalDiskon) }}</span>
          </div>
        </div>

        <div class="flex justify-between items-center text-2xl font-black mb-10 text-gray-900">
          <span>TOTAL DIBAYAR</span>
          <span>Rp {{ formatRupiah(dataNota.totalAkhir) }}</span>
        </div>

        <div class="text-center text-sm text-gray-500 hidden print:block mt-8 font-medium">
          Terima kasih telah berbelanja di SmartRetail!<br>
          Dicetak pada: {{ new Date().toLocaleString('id-ID') }}
        </div>
      </div>

      <div class="flex gap-4 print:hidden">
        <button @click="pesananBerhasil = false" class="w-1/3 bg-gray-100 text-gray-700 py-3 rounded-xl font-bold hover:bg-gray-200 transition border border-gray-200">Belanja Lagi</button>
        <button @click="batalkanPesanan" class="w-1/3 bg-white text-red-600 py-3 rounded-xl font-bold border border-red-200 hover:bg-red-50 transition">Batalkan Pesanan</button>
        <button @click="cetakStruk" class="w-1/3 bg-indigo-600 text-white py-3 rounded-xl font-bold shadow-md hover:bg-indigo-700 transition flex items-center justify-center gap-2">
          🖨️ Cetak Struk
        </button>
      </div>
    </div>

  </div>
  <!-- MODAL STRUK (Muncul setelah checkout berhasil) -->
  <div v-if="dataStruk" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900 bg-opacity-75 print:bg-white print:static print:inset-auto print:block">

    <!-- Kertas Struk -->
    <div id="area-struk" class="bg-white w-96 p-6 rounded-2xl shadow-2xl print:w-[80mm] print:p-0 print:m-0 print:shadow-none print:rounded-none text-slate-800">

      <!-- Header Toko -->
      <div class="text-center mb-4 border-b-2 border-dashed border-slate-300 pb-4">
        <h2 class="text-xl font-black uppercase tracking-widest">SmartRetail</h2>
        <p class="text-xs font-medium text-slate-500">Jl. Teknologi No. 404, Kota Inovasi</p>
        <p class="text-xs font-medium text-slate-500">Telp: 0812-3456-7890</p>
      </div>

      <!-- Meta Transaksi -->
      <div class="text-xs font-bold font-mono space-y-1 mb-4 border-b-2 border-dashed border-slate-300 pb-4">
        <div class="flex justify-between"><span>No:</span> <span>{{ dataStruk.no_pesanan }}</span></div>
        <div class="flex justify-between"><span>Tgl:</span> <span>{{ dataStruk.waktu }}</span></div>
        <div class="flex justify-between"><span>Kasir:</span> <span>{{ dataStruk.kasir }}</span></div>
        <div class="flex justify-between"><span>Plgn:</span> <span>{{ dataStruk.member }}</span></div>
      </div>

      <div class="mb-4 border-b-2 border-dashed border-slate-300 pb-4">
        <table class="w-full text-xs font-mono">
          <tbody>
          <tr v-for="(item, idx) in dataStruk.items" :key="idx">
            <td class="py-1">
              <div class="font-bold">{{ item.nama_produk }}</div>
              <div class="text-slate-500">{{ item.kuantitas }} x {{ formatRupiah(item.harga) }}</div>
            </td>
            <td class="py-1 text-right align-bottom font-bold">
              {{ formatRupiah(item.harga * item.kuantitas) }}
            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <div class="text-xs font-bold font-mono space-y-1 mb-6 border-b-2 border-dashed border-slate-300 pb-4">
        <div class="flex justify-between"><span>Subtotal:</span> <span>Rp {{ formatRupiah(dataStruk.subtotal) }}</span></div>
        <div v-if="dataStruk.diskon > 0" class="flex justify-between text-red-600 print:text-black">
          <span>Diskon Promo:</span> <span>- Rp {{ formatRupiah(dataStruk.diskon) }}</span>
        </div>
        <div class="flex justify-between text-base font-black mt-2 pt-2 border-t border-slate-200">
          <span>TOTAL:</span> <span>Rp {{ formatRupiah(dataStruk.totalAkhir) }}</span>
        </div>
      </div>

      <div class="text-center">
        <p class="text-xs font-bold text-slate-500 mb-6 uppercase tracking-wider">Terima Kasih Atas Kunjungan Anda</p>
        <div class="flex gap-2 print:hidden">
          <button @click="tutupStruk" class="flex-1 bg-slate-100 text-slate-700 py-2 rounded-lg font-bold hover:bg-slate-200 transition">Tutup</button>
          <button @click="cetakStruk" class="flex-1 bg-indigo-600 text-white py-2 rounded-lg font-bold hover:bg-indigo-700 transition">🖨️ Cetak Struk</button>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
.scrollbar-none {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

@media print {
  body * {
    visibility: hidden;
  }

  #area-struk, #area-struk * {
    visibility: visible;
  }

  #area-struk {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    font-family: "Courier New", Courier, monospace !important;
  }
}
</style>




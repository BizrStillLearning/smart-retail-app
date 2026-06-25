<script setup>
import { ref, onMounted, watch } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import Navbar from "../components/Navbar.vue";

const authStore = useAuthStore()

const api = axios.create({
  baseURL: 'http://localhost:8082/api',
  headers: {
    Authorization: `Bearer ${authStore.token}`
  }
})

const produks = ref([])
const loading = ref(true)

const searchQuery = ref('')
const selectedKategori = ref('')
const currentPage = ref(1)
const limit = ref(5)
const totalPages = ref(1)
const totalData = ref(0)

const showModal = ref(false)
const isEditing = ref(false)
const form = ref({ id_produk: null, nama_produk: '', kategori: '', harga: 0 })

const fetchProduk = async () => {
  loading.value = true
  try {
    const response = await api.get('/produk', {
      params: {
        search: searchQuery.value,
        kategori: selectedKategori.value,
        page: currentPage.value,
        limit: limit.value
      }
    })
    produks.value = response.data.data
    totalPages.value = response.data.last_page
    totalData.value = response.data.total
  } catch (err) {
    alert('Gagal mengambil data produk')
  } finally {
    loading.value = false
  }
}

watch([searchQuery, selectedKategori, limit], () => {
  currentPage.value = 1
  fetchProduk()
})

const nextPage = () => { if (currentPage.value < totalPages.value) { currentPage.value++; fetchProduk() } }
const prevPage = () => { if (currentPage.value > 1) { currentPage.value--; fetchProduk() } }

const fileGambar = ref(null)

const handleFileChange = (event) => {
  fileGambar.value = event.target.files[0]
}

const bukaModalTambah = () => {
  isEditing.value = false;
  form.value = { id_produk: null, nama_produk: '', kategori: '', harga: 0 };
  fileGambar.value = null;
  showModal.value = true
}

const bukaModalEdit = (produk) => {
  isEditing.value = true;
  form.value = { ...produk };
  fileGambar.value = null;
  showModal.value = true
}

const simpanProduk = async () => {
  try {
    const formData = new FormData()
    formData.append('nama_produk', form.value.nama_produk)
    formData.append('kategori', form.value.kategori)
    formData.append('harga', form.value.harga)

    if (fileGambar.value) {
      formData.append('gambar', fileGambar.value)
    }

    if (isEditing.value) {
      await api.put(`/admin/produk/${form.value.id_produk}`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
    } else {
      await api.post('/admin/produk', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
    }

    showModal.value = false
    fetchProduk()
  } catch (err) {
    alert(`Error: ${err.response?.data?.error || err.message}`)
  }
}

const hapusProduk = async (id) => {
  if (confirm('Yakin ingin menghapus produk ini?')) {
    try {
      await api.delete(`/admin/produk/${id}`)
      fetchProduk()
    } catch (err) {
      alert(`Error: ${err.response?.data?.error || err.message}`)
    }
  }
}

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)

onMounted(() => fetchProduk())
</script>

<template>
  <Navbar />

  <div class="p-8 bg-gray-50 min-h-screen">
    <div class="max-w-6xl mx-auto">

      <div class="flex justify-between items-center mb-6 bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <div>
          <h1 class="text-3xl font-bold text-gray-800">Manajemen Produk</h1>
          <p class="text-gray-500 mt-1">Kelola data katalog, filter, dan cari produk</p>
        </div>
        <div class="flex gap-3">
          <router-link to="/admin" class="bg-gray-100 text-gray-700 px-5 py-2.5 rounded-lg hover:bg-gray-200 transition font-semibold border border-gray-300">
            ← Dasbor
          </router-link>
          <button @click="bukaModalTambah" class="bg-blue-600 text-white px-5 py-2.5 rounded-lg hover:bg-blue-700 transition font-semibold shadow-md">
            + Tambah Produk
          </button>
          <router-link to="/admin/pesanan" class="bg-indigo-600 text-white px-5 py-2.5 rounded-lg hover:bg-indigo-700 transition font-semibold shadow-md">
            Daftar Pesanan
          </router-link>
        </div>
      </div>

      <div class="flex flex-col md:flex-row gap-4 mb-6 bg-white p-4 rounded-xl shadow-sm border border-gray-100">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama produk..."
            class="flex-1 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
        >
        <select v-model="selectedKategori" class="px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none bg-white">
          <option value="">Semua Kategori</option>
          <option value="Elektronik">Elektronik</option>
          <option value="Pakaian">Pakaian</option>
          <option value="Makanan">Makanan</option>
        </select>
        <select v-model="limit" class="px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none bg-white">
          <option :value="5">5 Per Halaman</option>
          <option :value="10">10 Per Halaman</option>
          <option :value="50">50 Per Halaman</option>
        </select>
      </div>

      <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <div v-if="loading" class="p-8 text-center text-gray-500 animate-pulse">Menarik data dari database...</div>

        <div v-else-if="produks.length === 0" class="p-12 text-center text-gray-500">
          Tidak ada produk yang sesuai dengan pencarian Anda.
        </div>

        <table v-else class="w-full text-left border-collapse">
          <thead>
          <tr class="bg-gray-50 text-gray-600 border-b border-gray-200 text-sm">
            <th class="p-4 font-semibold">ID</th>
            <th class="p-4 font-semibold">Nama Produk</th>
            <th class="p-4 font-semibold">Kategori</th>
            <th class="p-4 font-semibold">Harga</th>
            <th class="p-4 font-semibold text-right">Aksi</th>
            <th class="p-4 font-semibold">Gambar</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="produk in produks" :key="produk.id_produk" class="border-b border-gray-50 hover:bg-gray-50 transition">
            <td class="p-4 text-gray-400 text-sm">#{{ produk.id_produk }}</td>
            <td class="p-4 font-bold text-gray-800">{{ produk.nama_produk }}</td>
            <td class="p-4"><span class="bg-blue-50 text-blue-700 border border-blue-100 px-3 py-1 rounded-md text-xs font-semibold">{{ produk.kategori }}</span></td>
            <td class="p-4 font-medium text-gray-900">Rp {{ formatRupiah(produk.harga) }}</td>
            <td class="p-4 text-right space-x-2">
              <button @click="bukaModalEdit(produk)" class="text-blue-600 hover:text-blue-800 font-bold px-3 py-1.5 bg-blue-50 rounded border border-blue-100 text-sm">Edit</button>
              <button @click="hapusProduk(produk.id_produk)" class="text-red-600 hover:text-red-800 font-bold px-3 py-1.5 bg-red-50 rounded border border-red-100 text-sm">Hapus</button>
            </td>
            <td class="p-4">
              <img v-if="produk.gambar" :src="`http://localhost:8082${produk.gambar}`" class="w-12 h-12 object-cover rounded-lg border border-gray-200" alt="Produk" />
              <div v-else class="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center text-xs text-gray-400">No Img</div>
            </td>

            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-1">Upload Gambar</label>
              <input type="file" @change="handleFileChange" accept="image/*" class="w-full px-4 py-2 border rounded-lg file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 transition">
              <p v-if="isEditing" class="text-xs text-gray-500 mt-1">*Kosongkan jika tidak ingin mengubah gambar</p>
            </div>
          </tr>
          </tbody>
        </table>

        <div class="bg-gray-50 p-4 border-t border-gray-100 flex justify-between items-center">
          <span class="text-sm text-gray-500">
            Menampilkan halaman {{ currentPage }} dari {{ totalPages }} (Total: {{ totalData }} produk)
          </span>
          <div class="flex gap-2">
            <button @click="prevPage" :disabled="currentPage === 1" class="px-4 py-2 border rounded bg-white text-gray-700 disabled:bg-gray-100 disabled:text-gray-400 font-medium hover:bg-gray-50 transition">Sebelumnya</button>
            <button @click="nextPage" :disabled="currentPage === totalPages" class="px-4 py-2 border rounded bg-white text-gray-700 disabled:bg-gray-100 disabled:text-gray-400 font-medium hover:bg-gray-50 transition">Selanjutnya</button>
          </div>
        </div>
      </div>

    </div>

    <div v-if="showModal" class="fixed inset-0 bg-gray-900/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-2xl p-8 max-w-md w-full shadow-2xl">
        <h2 class="text-2xl font-bold mb-6 text-gray-800">{{ isEditing ? 'Edit Produk' : 'Tambah Produk Baru' }}</h2>

        <form @submit.prevent="simpanProduk" class="space-y-5">
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">Nama Produk</label>
            <input v-model="form.nama_produk" required type="text" placeholder="Misal: Kemeja Polos" class="w-full px-4 py-3 rounded-lg border focus:ring-2 focus:ring-blue-500 outline-none transition">
          </div>
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">Kategori</label>
            <select v-model="form.kategori" required class="w-full px-4 py-3 rounded-lg border focus:ring-2 focus:ring-blue-500 outline-none bg-white transition">
              <option value="Elektronik">Elektronik</option>
              <option value="Pakaian">Pakaian</option>
              <option value="Makanan">Makanan</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-1">Harga (Rp)</label>
            <input v-model="form.harga" required type="number" min="0" class="w-full px-4 py-3 rounded-lg border focus:ring-2 focus:ring-blue-500 outline-none transition">
          </div>

          <div class="flex gap-3 mt-8">
            <button type="button" @click="showModal = false" class="w-1/2 bg-gray-100 text-gray-700 py-3 rounded-lg font-bold hover:bg-gray-200 transition">Batal</button>
            <button type="submit" class="w-1/2 bg-blue-600 text-white py-3 rounded-lg font-bold hover:bg-blue-700 transition">Simpan Produk</button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>
<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'nuxt/app'
import { useCartStore } from '~/stores/cart'
import { useSwal } from '~/composables/useSwal'

const route = useRoute()
const cartStore = useCartStore()
const swal = useSwal()

const produkId = ref(route.params.id)
const produk = ref(null)
const produkRekomendasi = ref([])
const isLoading = ref(true)
const isLoadingRekomendasi = ref(true)
const kuantitas = ref(1)

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)

const muatDataProdukDanRekomendasi = async (id) => {
  isLoading.value = true
  isLoadingRekomendasi.value = true
  try {
    const resProduk = await $fetch(`http://localhost:8082/api/pelanggan/produk/${id}`)
    if (resProduk.status === 'sukses') produk.value = resProduk.data

    const resRekomendasi = await $fetch(`http://localhost:8082/api/pelanggan/produk/${id}/rekomendasi`)
    if (resRekomendasi.status === 'sukses') produkRekomendasi.value = resRekomendasi.data
  } catch (error) {
    console.error("Gagal memuat komponen halaman:", error)
  } finally {
    isLoading.value = false
    isLoadingRekomendasi.value = false
  }
}

onMounted(() => {
  muatDataProdukDanRekomendasi(produkId.value)
})

watch(() => route.params.id, (newId) => {
  if (newId) {
    produkId.value = newId
    kuantitas.value = 1
    muatDataProdukDanRekomendasi(newId)
  }
})

const tambahKuantitas = () => { if (kuantitas.value < produk.value.stok) kuantitas.value++ }
const kurangKuantitas = () => { if (kuantitas.value > 1) kuantitas.value-- }

const prosesTambahKeranjang = () => {
  cartStore.tambahKeKeranjang({
    id: produk.value.id_produk,
    nama: produk.value.nama_produk,
    harga: produk.value.harga,
    ikon: produk.value.ikon || '📦'
  }, kuantitas.value)
  swal.sukses(`${kuantitas.value} item dimasukkan ke keranjang.`)
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <NuxtLink to="/" class="inline-flex items-center gap-2 text-sm font-bold text-slate-500 hover:text-indigo-600 transition mb-8">
      ← Kembali ke Katalog
    </NuxtLink>

    <div v-if="isLoading" class="text-center py-12 text-slate-500 font-bold animate-pulse">Memuat data produk...</div>
    <div v-else-if="!produk" class="bg-white p-12 text-center rounded-3xl border border-slate-100 shadow-sm">
      <h2 class="text-xl font-bold text-slate-700">Produk Tidak Tersedia</h2>
    </div>

    <div v-else class="space-y-16">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-12 bg-white p-8 sm:p-12 rounded-3xl border border-slate-200 shadow-sm">
        <div class="bg-slate-50 rounded-2xl h-96 flex items-center justify-center text-9xl border border-slate-100">
          {{ produk.ikon || '📦' }}
        </div>

        <div class="flex flex-col justify-between">
          <div>
            <span class="text-xs font-black text-indigo-600 uppercase tracking-widest bg-indigo-50 px-3 py-1 rounded-full border border-indigo-100">
              {{ produk.kategori || 'Katalog' }}
            </span>
            <h1 class="text-3xl font-black text-slate-900 tracking-tight mt-4 mb-2">{{ produk.nama_produk }}</h1>
            <div class="text-2xl font-black text-slate-900 tracking-tight mb-6">Rp {{ formatRupiah(produk.harga) }}</div>

            <div class="border-t border-slate-100 pt-6">
              <h3 class="text-sm font-bold text-slate-800 uppercase tracking-wider mb-2">Deskripsi</h3>
              <p class="text-slate-500 text-sm leading-relaxed font-medium">{{ produk.deskripsi || 'Sembako pilihan berkualitas tinggi.' }}</p>
            </div>
          </div>

          <div class="border-t border-slate-100 pt-6 space-y-4">
            <div class="flex items-center justify-between text-sm font-bold text-slate-600">
              <span>Jumlah</span>
              <span class="text-xs text-slate-400">Stok: {{ produk.stok }} unit</span>
            </div>
            <div class="flex flex-col sm:flex-row gap-4">
              <div class="flex items-center justify-between bg-slate-100 p-1.5 rounded-xl border border-slate-200 sm:w-32">
                <button @click="kurangKuantitas" class="w-8 h-8 bg-white rounded-lg flex items-center justify-center font-bold text-slate-700 shadow-sm shadow-slate-100">-</button>
                <span class="font-black text-slate-800 text-sm">{{ kuantitas }}</span>
                <button @click="tambahKuantitas" class="w-8 h-8 bg-white rounded-lg flex items-center justify-center font-bold text-slate-700 shadow-sm shadow-slate-100">+</button>
              </div>
              <button @click="prosesTambahKeranjang" class="flex-1 bg-indigo-600 text-white font-bold px-6 py-3.5 rounded-xl hover:bg-indigo-700 transition shadow-md">
                Tambah ke Keranjang
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="border-t border-slate-200 pt-12">
        <h2 class="text-2xl font-black text-slate-900 tracking-tight mb-6 flex items-center gap-2">
          <span>🧠</span> Sering Dibeli Bersamaan
        </h2>

        <div v-if="isLoadingRekomendasi" class="grid grid-cols-2 md:grid-cols-4 gap-6 animate-pulse">
          <div v-for="n in 4" :key="n" class="h-48 bg-slate-100 rounded-2xl"></div>
        </div>

        <div v-else-if="produkRekomendasi.length === 0" class="text-slate-400 text-sm font-medium bg-slate-50 p-8 rounded-2xl border border-dashed border-slate-200 text-center">
          Belum ada pola transaksi yang cukup untuk merender rekomendasi pada produk ini.
        </div>

        <div v-else class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <NuxtLink v-for="item in produkRekomendasi" :key="item.id_produk" :to="`/produk/${item.id_produk}`" class="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm hover:border-indigo-300 hover:shadow-md transition group text-left block">
            <div class="aspect-square bg-slate-50 rounded-xl flex items-center justify-center text-5xl mb-4 group-hover:scale-105 transition-transform">
              {{ item.ikon || '📦' }}
            </div>
            <h3 class="font-bold text-slate-800 text-sm line-clamp-1 mb-1">{{ item.nama_produk }}</h3>
            <div class="text-indigo-600 font-black text-sm">Rp {{ formatRupiah(item.harga) }}</div>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>


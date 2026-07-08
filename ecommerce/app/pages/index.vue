<script setup>
import { ref, onMounted } from 'vue'

const produks = ref([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    const response = await $fetch('http://localhost:8082/api/pelanggan/produk')
    if (response.status === 'sukses') {
      produks.value = response.data
    }
  } catch (error) {
    console.error("Gagal memuat katalog:", error)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">

    <div class="w-full bg-gradient-to-br from-indigo-900 to-indigo-700 rounded-3xl p-10 text-white mb-12 shadow-xl flex justify-between items-center overflow-hidden relative">
      <div class="absolute -right-10 -top-10 w-64 h-64 bg-white opacity-5 rounded-full blur-3xl"></div>
      <div class="relative z-10">
        <span class="bg-indigo-500 text-white text-xs font-black px-3 py-1 rounded-full uppercase tracking-widest mb-4 inline-block">Promo Spesial</span>
        <h2 class="text-4xl font-black mb-3 tracking-tight">Belanja Hemat Akhir Bulan!</h2>
        <p class="text-indigo-200 font-medium text-lg">Gunakan kode <span class="bg-white text-indigo-900 px-2 py-0.5 rounded font-black mx-1 shadow-sm">HEMAT50</span> untuk potongan Rp 50.000</p>
      </div>
      <div class="hidden md:block text-8xl relative z-10 drop-shadow-2xl hover:scale-110 transition-transform cursor-pointer">🛍️</div>
    </div>

    <div class="flex justify-between items-end mb-8">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Eksplorasi Katalog</h2>
        <p class="text-sm font-medium text-slate-500 mt-1">Sembako segar langsung diantar ke rumah Anda.</p>
      </div>
      <div class="hidden sm:flex gap-2">
        <button class="px-4 py-2 bg-slate-900 text-white text-sm font-bold rounded-full">Semua</button>
        <button class="px-4 py-2 bg-white text-slate-600 text-sm font-bold rounded-full border border-slate-200 hover:border-indigo-500 transition">Kebutuhan Pokok</button>
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <ProductCard v-for="produk in produks" :key="produk.id" :produk="produk" />
    </div>

  </div>
</template>


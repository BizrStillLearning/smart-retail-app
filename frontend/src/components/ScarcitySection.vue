<script setup>
import { computed } from 'vue'

const props = defineProps({
  produks: { type: Array, required: true }
})

defineEmits(['tambah-ke-keranjang'])

const produkKritis = computed(() => {
  return props.produks.filter(p => p.stok > 0 && p.stok <= 5)
})

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)
</script>

<template>
  <div v-if="produkKritis.length > 0" class="mb-8 bg-red-50/70 border border-red-100 rounded-2xl p-5">
    <h2 class="text-md font-black text-red-700 flex items-center gap-2 uppercase tracking-wide mb-4">
      <span class="animate-pulse flex h-2 w-2 rounded-full bg-red-600"></span>
      Jangan Sampai Kehabisan! (Stok Menipis)
    </h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="produk in produkKritis" :key="produk.id_produk" class="bg-white p-4 rounded-xl border border-red-100 flex gap-4 items-center shadow-sm relative">
        <div class="w-16 h-16 flex-shrink-0 bg-gray-50 rounded-lg overflow-hidden border border-gray-100">
          <img v-if="produk.gambar" :src="`http://localhost:8082${produk.gambar}`" class="w-full h-full object-cover" />
          <div v-else class="w-full h-full flex items-center justify-center text-[10px] text-gray-400">No Img</div>
        </div>

        <div class="flex-1 min-w-0">
          <h4 class="text-sm font-bold text-gray-800 truncate">{{ produk.nama_produk }}</h4>
          <p class="text-xs text-gray-900 font-black mt-0.5">Rp {{ formatRupiah(produk.harga) }}</p>

          <div class="mt-2">
            <div class="flex justify-between text-[10px] text-red-600 font-bold mb-0.5">
              <span>Sisa {{ produk.stok }} item!</span>
            </div>
            <div class="w-full bg-gray-100 h-1.5 rounded-full overflow-hidden">
              <div class="bg-red-500 h-full rounded-full transition-all duration-500" :style="{ width: `${(produk.stok / 20) * 100}%` }"></div>
            </div>
          </div>
        </div>

        <button
            @click="$emit('tambah-ke-keranjang', produk)"
            class="bg-red-600 hover:bg-red-700 text-white p-2 rounded-lg transition shadow-sm"
            title="Serbu Produk"
        >
          🛒
        </button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed } from 'vue'

const props = defineProps({
  produks: { type: Array, required: true }
})

const topTerlaris = computed(() => {
  return [...props.produks].slice(0, 3)
})

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)
</script>

<template>
  <div v-if="topTerlaris.length > 0" class="bg-white p-5 rounded-2xl border border-gray-100 shadow-sm mb-6">
    <h3 class="text-sm font-black text-gray-800 flex items-center gap-2 uppercase tracking-wider mb-4">
      👑 Top 3 Terlaris Minggu Ini
    </h3>
    <div class="space-y-3">
      <div v-for="(produk, idx) in topTerlaris" :key="produk.id_produk" class="flex items-center justify-between p-2 rounded-xl bg-gray-50 border border-gray-100">
        <div class="flex items-center gap-3 min-w-0">
          <div :class="{
            'bg-amber-400 text-white': idx === 0,
            'bg-slate-400 text-white': idx === 1,
            'bg-amber-600 text-white': idx === 2
          }" class="w-6 h-6 rounded-full flex items-center justify-center font-black text-xs shadow-sm">
            {{ idx + 1 }}
          </div>
          <div class="truncate">
            <div class="font-bold text-gray-800 text-xs truncate">{{ produk.nama_produk }}</div>
            <div class="text-[11px] font-medium text-gray-500">Rp {{ formatRupiah(produk.harga) }}</div>
          </div>
        </div>
        <span class="text-[10px] font-bold text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-md border border-indigo-100 whitespace-nowrap">
          🔥 Kebutuhan Utama
        </span>
      </div>
    </div>
  </div>
</template>
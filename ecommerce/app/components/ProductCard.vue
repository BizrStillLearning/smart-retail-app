<script setup>
defineProps({
  produk: {
    type: Object,
    required: true
  }
})

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)
</script>

<template>
  <div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden hover:shadow-md transition-all group flex flex-col cursor-pointer">

    <NuxtLink :to="`/produk/${produk.id}`" class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden hover:shadow-md transition-all group flex flex-col cursor-pointer">

      <div class="h-48 bg-slate-50 relative overflow-hidden flex items-center justify-center border-b border-slate-100">
        <span class="text-6xl group-hover:scale-110 transition-transform duration-300">{{ produk.ikon || '📦' }}</span>
        <div v-if="produk.stok < 10" class="absolute top-3 right-3 bg-red-500 text-white text-[10px] font-black px-2.5 py-1 rounded-full uppercase tracking-wider">
          Sisa {{ produk.stok }}
        </div>
      </div>

    </NuxtLink>

    <div class="p-5 flex flex-col flex-grow">
      <div class="text-xs font-black text-indigo-500 mb-1.5 uppercase tracking-wider">{{ produk.kategori }}</div>
      <h3 class="text-slate-800 font-bold leading-snug mb-3 flex-grow">{{ produk.nama_produk }}</h3>

      <div class="flex items-end justify-between mt-auto">
        <div>
          <div v-if="produk.diskon" class="text-xs font-bold text-slate-400 line-through mb-0.5">Rp {{ formatRupiah(produk.harga_asli) }}</div>
          <div class="text-lg font-black text-slate-900 tracking-tight">Rp {{ formatRupiah(produk.harga) }}</div>
        </div>
        <button class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center hover:bg-indigo-600 hover:text-white transition-colors shadow-sm">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"></path></svg>
        </button>
      </div>
    </div>
  </div>
</template>
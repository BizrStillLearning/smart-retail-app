<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useToast } from 'vue-toastification'

const toast = useToast()
const activeSlide = ref(0)
let timer = null

const listPromo = ref([
  { id: 1, kode: 'SEMBAKOMURAH', deskripsi: 'Potongan Rp 5.000 untuk beras dan minyak!', bg: 'from-orange-500 to-red-600' },
  { id: 2, kode: 'JUMATBERKAH', deskripsi: 'Diskon kilat 10% khusus hari Jumat minumal 50rb.', bg: 'from-indigo-600 to-purple-600' },
  { id: 3, kode: 'MEMBERBARU', deskripsi: 'Gratis ongkir tanpa minimal belanja bagi member baru.', bg: 'from-emerald-500 to-teal-600' }
])

const nextSlide = () => {
  activeSlide.value = (activeSlide.value + 1) % listPromo.value.length
}

const prevSlide = () => {
  activeSlide.value = (activeSlide.value - 1 + listPromo.value.length) % listPromo.value.length
}

const salinKode = (kode) => {
  navigator.clipboard.writeText(kode)
  toast.success(`Kode promo "${kode}" berhasil disalin!`, { timeout: 2000 })
}

onMounted(() => {
  timer = setInterval(nextSlide, 5000) // Geser otomatis tiap 5 detik
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <div class="relative overflow-hidden rounded-2xl bg-gradient-to-r shadow-lg mb-8 h-44 flex items-center transition-all duration-500" :class="listPromo[activeSlide].bg">
    <div class="p-8 text-white w-full flex justify-between items-center z-10">
      <div class="max-w-md">
        <span class="bg-white/20 text-xs font-black uppercase tracking-widest px-2.5 py-1 rounded-md border border-white/20">Promo Aktif 🔥</span>
        <h2 class="text-2xl font-black mt-2 leading-tight tracking-tight">{{ listPromo[activeSlide].deskripsi }}</h2>
      </div>

      <div class="flex flex-col items-end gap-2">
        <div class="bg-white text-gray-900 font-mono font-black text-lg px-4 py-2 rounded-xl shadow-sm border border-white/50 select-all">
          {{ listPromo[activeSlide].kode }}
        </div>
        <button @click="salinKode(listPromo[activeSlide].kode)" class="bg-white/20 hover:bg-white text-white hover:text-gray-900 font-bold text-xs px-4 py-2 rounded-lg transition border border-white/30 backdrop-blur-sm">
          📋 Salin Kode
        </button>
      </div>
    </div>

    <button @click="prevSlide" class="absolute left-3 p-2 rounded-full bg-black/20 hover:bg-black/40 text-white transition z-20">❮</button>
    <button @click="nextSlide" class="absolute right-3 p-2 rounded-full bg-black/20 hover:bg-black/40 text-white transition z-20">❯</button>

    <div class="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-1.5 z-20">
      <span
          v-for="(promo, index) in listPromo"
          :key="promo.id"
          @click="activeSlide = index"
          :class="index === activeSlide ? 'w-6 bg-white' : 'w-2 bg-white/50'"
          class="h-2 rounded-full cursor-pointer transition-all duration-300"
      ></span>
    </div>
  </div>
</template>
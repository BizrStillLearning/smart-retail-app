<script setup>
import { ref } from 'vue'
import { useRouter } from 'nuxt/app'
import { useCartStore } from '~/stores/cart'
import { useCustomerAuthStore } from '~/stores/customerAuth'
import { useSwal } from '~/composables/useSwal'

const cartStore = useCartStore()
const authStore = useCustomerAuthStore()
const router = useRouter()
const swal = useSwal()

const isMemproses = ref(false)
// State reaktif untuk metode pembayaran
const metodePembayaran = ref('Transfer Bank')

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)

const prosesCheckout = async () => {
  if (cartStore.items.length === 0) {
    swal.error('Keranjang belanja Anda masih kosong.')
    return
  }

  if (!authStore.isLoggedIn) {
    swal.info('Akses ditolak: Harap masuk ke akun Anda untuk melanjutkan pembayaran.')
    router.push('/login')
    return
  }

  isMemproses.value = true

  try {
    // Menyertakan metode pembayaran ke dalam payload JSON
    const payload = {
      items: cartStore.items.map(item => ({
        id: item.id,
        kuantitas: item.kuantitas,
        harga: item.harga
      })),
      metode_pembayaran: metodePembayaran.value
    }

    const response = await $fetch('http://localhost:8082/api/shop/checkout', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: payload
    })

    if (response.status === 'sukses') {
      swal.sukses(`Pesanan #${response.id_pesanan} berhasil dicatat!`)
      cartStore.items = []
      router.push('/riwayat')
    }
  } catch (error) {
    const pesanError = error.response?._data?.error || 'Gagal memproses transaksi. Server tidak merespons.'
    swal.error(pesanError)
  } finally {
    isMemproses.value = false
  }
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <h1 class="text-3xl font-black text-slate-900 tracking-tight mb-8">Keranjang Belanja</h1>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">

      <div class="lg:col-span-2 space-y-4">
        <div v-if="cartStore.items.length === 0" class="bg-white p-12 rounded-3xl border border-slate-200 text-center shadow-sm">
          <div class="text-6xl mb-4">🛒</div>
          <h2 class="text-xl font-bold text-slate-700 mb-2">Keranjang masih kosong</h2>
          <p class="text-slate-500 mb-6 text-sm">Silakan kembali ke katalog untuk memilih produk sembako Anda.</p>
          <NuxtLink to="/" class="bg-indigo-600 text-white font-bold px-6 py-2.5 rounded-full hover:bg-indigo-700 transition">
            Lihat Katalog
          </NuxtLink>
        </div>

        <div v-else v-for="item in cartStore.items" :key="item.id" class="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-4 hover:border-indigo-200 transition">
          <div class="w-20 h-20 bg-slate-50 rounded-xl flex items-center justify-center text-4xl border border-slate-100">
            {{ item.ikon }}
          </div>
          <div class="flex-grow">
            <h3 class="font-bold text-slate-800 text-lg">{{ item.nama }}</h3>
            <div class="text-indigo-600 font-black text-sm">Rp {{ formatRupiah(item.harga) }}</div>
          </div>

          <div class="flex items-center gap-4">
            <div class="text-sm font-bold text-slate-500 bg-slate-100 px-3 py-1 rounded-lg">
              Qty: {{ item.kuantitas }}
            </div>
            <div class="font-black text-slate-900 text-lg w-28 text-right">
              Rp {{ formatRupiah(item.harga * item.kuantitas) }}
            </div>
            <button @click="cartStore.hapusDariKeranjang(item.id)" class="w-10 h-10 bg-red-50 text-red-500 rounded-xl flex items-center justify-center hover:bg-red-500 hover:text-white transition">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
        </div>
      </div>

      <div class="lg:col-span-1">
        <div class="bg-white p-6 rounded-3xl border border-slate-200 shadow-sm sticky top-24">
          <h2 class="text-lg font-black text-slate-800 mb-6 border-b border-slate-100 pb-4">Ringkasan Belanja</h2>

          <div class="space-y-3 text-sm font-bold text-slate-600 mb-6">
            <div class="flex justify-between">
              <span>Total Item</span>
              <span class="text-slate-900">{{ cartStore.totalItem }} Produk</span>
            </div>
            <div class="flex justify-between">
              <span>Subtotal Harga</span>
              <span class="text-slate-900">Rp {{ formatRupiah(cartStore.totalHarga) }}</span>
            </div>
          </div>

          <div class="border-t border-slate-100 pt-6 mb-6 space-y-3">
            <h3 class="text-sm font-bold text-slate-800">Metode Pembayaran</h3>

            <label class="flex items-center gap-3 p-3 border border-slate-200 rounded-xl cursor-pointer hover:bg-slate-50 transition" :class="{'border-indigo-500 bg-indigo-50': metodePembayaran === 'Transfer Bank'}">
              <input type="radio" v-model="metodePembayaran" value="Transfer Bank" class="w-4 h-4 text-indigo-600 focus:ring-indigo-500">
              <span class="text-sm font-bold text-slate-700">Transfer Bank Manual</span>
            </label>

            <label class="flex items-center gap-3 p-3 border border-slate-200 rounded-xl cursor-pointer hover:bg-slate-50 transition" :class="{'border-indigo-500 bg-indigo-50': metodePembayaran === 'COD'}">
              <input type="radio" v-model="metodePembayaran" value="COD" class="w-4 h-4 text-indigo-600 focus:ring-indigo-500">
              <span class="text-sm font-bold text-slate-700">Bayar di Tempat (COD)</span>
            </label>
          </div>

          <div class="border-t border-slate-100 pt-6 mb-6">
            <div class="flex justify-between items-end">
              <span class="text-sm font-bold text-slate-500">Total Pembayaran</span>
              <span class="text-2xl font-black text-indigo-600 tracking-tight">Rp {{ formatRupiah(cartStore.totalHarga) }}</span>
            </div>
          </div>

          <button @click="prosesCheckout" :disabled="cartStore.items.length === 0 || isMemproses" class="w-full bg-indigo-600 text-white font-black py-4 rounded-xl hover:bg-indigo-700 transition disabled:opacity-50 disabled:cursor-not-allowed shadow-md shadow-indigo-100">
            {{ isMemproses ? 'Menghubungi Server...' : 'Lanjutkan Pembayaran' }}
          </button>
        </div>
      </div>

    </div>
  </div>
</template>
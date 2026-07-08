<script setup>
import { useAuthStore } from '../stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'vue-toastification'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const toast = useToast()

const prosesLogout = () => {
  if (confirm('Sesi operasional akan diakhiri. Lanjutkan?')) {
    authStore.logout()
    toast.info('Anda telah keluar dari sistem.')
    router.push('/login')
  }
}

const menuItems = [
  { name: 'Dasbor Analitik', path: '/admin', icon: '📊' },
  { name: 'Katalog Produk', path: '/admin/produk', icon: '📦' },
  { name: 'Riwayat Pesanan', path: '/admin/pesanan', icon: '🛒' },
  { name: 'Manajemen Promo', path: '/admin/promo', icon: '🎟️' },
  { name: 'Data Pelanggan', path: '/admin/pelanggan', icon: '👥' }
]

const isActive = (path) => route.path === path
</script>

<template>
  <div class="flex h-screen bg-slate-50 font-sans overflow-hidden">

    <aside class="w-64 bg-slate-900 text-slate-300 flex flex-col transition-all shadow-xl z-20 print:hidden">
      <div class="h-16 flex items-center px-6 border-b border-slate-800 bg-slate-950">
        <span class="text-2xl font-black text-indigo-500 tracking-tight">Smart<span class="text-white">Admin</span></span>
      </div>

      <nav class="flex-1 overflow-y-auto py-6 px-3 space-y-1">
        <p class="px-3 text-xs font-bold text-slate-500 uppercase tracking-widest mb-3">Menu Operasional</p>

        <router-link
            v-for="item in menuItems"
            :key="item.path"
            :to="item.path"
            :class="[isActive(item.path) ? 'bg-indigo-600 text-white shadow-md' : 'hover:bg-slate-800 hover:text-white', 'flex items-center gap-3 px-3 py-2.5 rounded-lg font-medium transition-all']"
        >
          <span class="text-lg">{{ item.icon }}</span>
          {{ item.name }}
        </router-link>
      </nav>

      <div class="p-4 border-t border-slate-800">
        <div class="flex items-center gap-3 mb-4 px-2">
          <div class="w-8 h-8 rounded-full bg-indigo-500 flex items-center justify-center text-white font-bold">
            {{ authStore.namaLengkap.charAt(0) || 'A' }}
          </div>
          <div class="text-sm">
            <p class="text-white font-bold truncate">{{ authStore.namaLengkap || 'Administrator' }}</p>
            <p class="text-xs text-slate-400 capitalize">{{ authStore.role || 'Admin' }}</p>
          </div>
        </div>
      </div>
    </aside>

    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">

      <header class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-8 shadow-sm z-10 print:hidden">
        <div class="text-sm font-medium text-gray-500">
          Sistem Manajemen Ritel Terpadu
        </div>
        <div class="flex items-center gap-4">
          <router-link to="/" class="text-sm font-bold text-indigo-600 hover:text-indigo-800 transition">Lihat Mode Kasir ↗</router-link>
          <div class="h-5 w-px bg-gray-300"></div>
          <button @click="prosesLogout" class="text-sm font-bold text-red-600 hover:text-red-800 transition">Logout</button>
        </div>
      </header>

      <main class="flex-1 overflow-y-auto p-8 relative custom-scrollbar">
        <slot />
      </main>

      <footer class="bg-white border-t border-gray-200 py-4 px-8 text-xs text-gray-500 font-medium flex justify-between print:hidden">
        <span>&copy; {{ new Date().getFullYear() }} SmartRetail System.</span>
        <span>Versi 2.0.1 (Enterprise)</span>
      </footer>

    </div>
  </div>
</template>

<style>
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: #f8fafc; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>
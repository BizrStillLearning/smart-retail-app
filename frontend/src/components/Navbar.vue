<script setup>
import { useAuthStore } from '../stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'vue-toastification'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const toast = useToast()

const prosesLogout = () => {
  if (confirm('Apakah Anda yakin ingin keluar?')) {
    authStore.logout()
    toast.info('Sesi diakhiri. Anda telah logout.')
    router.push('/login')
  }
}

const isActive = (path) => route.path === path
</script>

<template>
  <nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">

        <div class="flex">
          <div class="flex-shrink-0 flex items-center mr-8">
            <span class="text-2xl font-black text-indigo-600 tracking-tight">Smart<span class="text-gray-800">Admin</span></span>
          </div>
          <div class="hidden sm:flex sm:space-x-8">
            <router-link to="/admin" :class="[isActive('/admin') ? 'border-indigo-500 text-gray-900 font-bold' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700', 'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition']">Dasbor</router-link>
            <router-link to="/admin/produk" :class="[isActive('/admin/produk') ? 'border-indigo-500 text-gray-900 font-bold' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700', 'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition']">Produk</router-link>
            <router-link to="/admin/pesanan" :class="[isActive('/admin/pesanan') ? 'border-indigo-500 text-gray-900 font-bold' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700', 'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition']">Pesanan</router-link>
            <router-link to="/admin/promo" :class="[isActive('/admin/promo') ? 'border-indigo-500 text-gray-900 font-bold' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700', 'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition']">Promo</router-link>
          </div>
        </div>

        <div class="flex items-center gap-4">
          <router-link to="/" class="text-sm font-medium text-gray-500 hover:text-indigo-600 transition">Lihat Toko ↗</router-link>
          <div class="h-6 w-px bg-gray-200"></div>
          <button @click="prosesLogout" class="bg-red-50 text-red-600 px-4 py-2 rounded-lg text-sm font-bold hover:bg-red-100 transition border border-red-100">Logout</button>
        </div>

      </div>
    </div>
  </nav>
</template>
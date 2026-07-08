<script setup>
import { useCartStore } from '~/stores/cart'
import { useCustomerAuthStore } from '~/stores/customerAuth'

const cartStore = useCartStore()
const authStore = useCustomerAuthStore()

</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-50 font-sans">

    <header class="bg-white border-b border-slate-200 sticky top-0 z-50 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">

          <div class="flex items-center gap-8">
            <NuxtLink to="/public" class="text-2xl font-black text-indigo-600 tracking-tight">
              Smart<span class="text-slate-900">Retail</span>
            </NuxtLink>
            <nav class="hidden md:flex gap-6 text-sm font-bold text-slate-600">
              <NuxtLink to="/public" class="hover:text-indigo-600 transition">Beranda</NuxtLink>
              <NuxtLink to="/kategori" class="hover:text-indigo-600 transition">Kategori</NuxtLink>
              <NuxtLink to="/promo" class="hover:text-indigo-600 transition">Promo Khusus</NuxtLink>
            </nav>
          </div>

          <div class="flex items-center gap-6">
            <div class="hidden lg:block relative">
              <input type="text" placeholder="Cari sembako..." class="w-64 pl-4 pr-10 py-2 bg-slate-100 border-transparent rounded-full text-sm focus:bg-white focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 outline-none transition">
              <span class="absolute right-3 top-2.5 text-slate-400">🔍</span>
            </div>

            <div class="flex items-center gap-4 border-l border-slate-200 pl-6">
              <NuxtLink to="/keranjang" class="relative text-slate-600 hover:text-indigo-600 transition">
                <span class="text-xl">🛒</span>
                <span class="absolute -top-1.5 -right-2 bg-red-500 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full">
                  {{ cartStore.totalItem }}
                </span>
              </NuxtLink>

              <NuxtLink v-if="!authStore.isLoggedIn" to="/login" class="text-sm font-bold bg-slate-900 text-white px-4 py-2 rounded-full hover:bg-slate-800 transition shadow-sm">
                Masuk
              </NuxtLink>

              <div v-else class="flex items-center gap-3">
                <span class="text-sm font-bold text-slate-700 bg-slate-100 px-3 py-1.5 rounded-full">Hai, {{ authStore.customerName }}</span>
                <button @click="authStore.logout" class="text-xs font-bold text-red-500 hover:text-red-700 transition">Logout</button>
              </div>
            </div>

          </div>

        </div>
      </div>
    </header>

    <main class="flex-grow">
      <slot />
    </main>

    <footer class="bg-white border-t border-slate-200 mt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 text-center text-sm font-medium text-slate-500">
        &copy; {{ new Date().getFullYear() }} SmartRetail E-Commerce. Hak Cipta Dilindungi.
      </div>
    </footer>

  </div>
</template>
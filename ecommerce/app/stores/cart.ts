import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCartStore = defineStore('cart', () => {
    const items = ref([])

    const totalItem = computed(() => items.value.reduce((total, item) => total + item.kuantitas, 0))
    const totalHarga = computed(() => items.value.reduce((total, item) => total + (item.harga * item.kuantitas), 0))

    const tambahKeKeranjang = (produk, kuantitas) => {
        const itemAda = items.value.find(item => item.id === produk.id)

        if (itemAda) {
            itemAda.kuantitas += kuantitas
        } else {
            items.value.push({ ...produk, kuantitas })
        }
    }

    const hapusDariKeranjang = (id) => {
        items.value = items.value.filter(item => item.id !== id)
    }

    return { items, totalItem, totalHarga, tambahKeKeranjang, hapusDariKeranjang }
})
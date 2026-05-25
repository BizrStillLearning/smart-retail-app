import { defineStore } from 'pinia'
import { ref, watch, computed } from 'vue'

export const useKeranjangStore = defineStore('keranjang', () => {
    const dataTersimpan = localStorage.getItem('data_keranjang')
    const items = ref(dataTersimpan ? JSON.parse(dataTersimpan) : [])

    watch(items, (nilaiBaru) => {
        localStorage.setItem('data_keranjang', JSON.stringify(nilaiBaru))
    }, { deep: true })

    const totalHarga = computed(() => {
        return items.value.reduce((total, item) => total + (item.harga * item.kuantitas), 0)
    })

    const tambahItem = (produk) => {
        const produkAda = items.value.find(item => item.id_produk === produk.id_produk)
        if (produkAda) {
            produkAda.kuantitas++
        } else {
            items.value.push({ ...produk, kuantitas: 1 })
        }
        alert(`${produk.nama_produk} berhasil ditambahkan!`)
    }

    const hapusItem = (id_produk) => {
        items.value = items.value.filter(item => item.id_produk !== id_produk)
    }

    const kosongkanKeranjang = () => {
        items.value = []
    }

    return { items, totalHarga, tambahItem, hapusItem, kosongkanKeranjang }
})
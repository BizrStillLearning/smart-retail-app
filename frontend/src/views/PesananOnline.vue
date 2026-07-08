<script setup>
import { ref, onMounted, computed } from 'vue'
import Swal from 'sweetalert2'

const daftarPesanan = ref([])
const isLoading = ref(true)

const tokenAdmin = localStorage.getItem('token') || ''

const muatPesananMasuk = async () => {
  isLoading.value = true
  try {
    const response = await fetch('http://localhost:8082/api/admin/pesanan-online', {
      headers: { 'Authorization': `Bearer ${tokenAdmin}` }
    })
    const result = await response.json()
    if (result.status === 'sukses') {
      daftarPesanan.value = result.data
    }
  } catch (error) {
    console.error("Gagal menarik data pesanan:", error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  muatPesananMasuk()
})

const formatRupiah = (angka) => new Intl.NumberFormat('id-ID').format(angka)
const formatTanggal = (tgl) => new Date(tgl).toLocaleString('id-ID')

const pesananPerluTindakan = computed(() => {
  return daftarPesanan.value.filter(p => p.status === 'verifikasi pembayaran' || p.status === 'diproses')
})

const ubahStatusPesanan = async (idPesanan, statusBaru) => {
  try {
    const response = await fetch(`http://localhost:8082/api/admin/pesanan/${idPesanan}/status`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${tokenAdmin}`
      },
      body: JSON.stringify({ status: statusBaru })
    })

    const result = await response.json()

    if (result.status === 'sukses') {
      Swal.fire({
        icon: 'success',
        title: 'Status Diperbarui',
        toast: true,
        position: 'top-end',
        timer: 3000,
        showConfirmButton: false
      })
      muatPesananMasuk()
    } else {
      Swal.fire({ icon: 'error', title: 'Gagal', text: result.error })
    }
  } catch (error) {
    Swal.fire({ icon: 'error', title: 'Gagal memproses', text: 'Koneksi ke server terputus.' })
  }
}
</script>

<template>
  <div class="p-8">
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-black text-slate-800">Manajemen Pesanan Online</h1>
        <p class="text-slate-500 font-medium mt-1">Kelola transaksi dari platform E-Commerce.</p>
      </div>
      <div class="bg-indigo-50 text-indigo-700 px-4 py-2 rounded-xl font-bold text-sm border border-indigo-100">
        {{ pesananPerluTindakan.length }} Pesanan Aktif
      </div>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
      <table class="w-full text-left border-collapse">
        <thead>
        <tr class="bg-slate-50 border-b border-slate-200 text-sm font-bold text-slate-600 uppercase tracking-wider">
          <th class="px-6 py-4">ID Pesanan</th>
          <th class="px-6 py-4">Tanggal & Waktu</th>
          <th class="px-6 py-4">Metode / Status</th>
          <th class="px-6 py-4">Total Belanja</th>
          <th class="px-6 py-4 text-center">Aksi Kasir</th>
        </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
        <tr v-for="pesanan in pesananPerluTindakan" :key="pesanan.id_pesanan" class="hover:bg-slate-50 transition">
          <td class="px-6 py-4">
            <div class="font-black text-slate-800">#{{ pesanan.id_pesanan }}</div>
            <div class="text-sm font-medium text-slate-500">ID Pelanggan: {{ pesanan.id_pelanggan }}</div>
          </td>
          <td class="px-6 py-4 text-sm font-medium text-slate-600">{{ formatTanggal(pesanan.tanggal_pesan) }}</td>
          <td class="px-6 py-4">
            <div class="text-sm font-bold text-slate-800 mb-1">{{ pesanan.metode_pembayaran }}</div>
            <span v-if="pesanan.status === 'verifikasi pembayaran'" class="bg-amber-100 text-amber-700 text-xs font-bold px-2.5 py-1 rounded-md">Perlu Cek Bukti</span>
            <span v-else-if="pesanan.status === 'diproses'" class="bg-blue-100 text-blue-700 text-xs font-bold px-2.5 py-1 rounded-md">Siap Dikemas</span>
          </td>
          <td class="px-6 py-4 font-black text-slate-800">Rp {{ formatRupiah(pesanan.total_akhir) }}</td>
          <td class="px-6 py-4 text-center">

            <div v-if="pesanan.status === 'verifikasi pembayaran'" class="flex flex-col gap-2">
              <a :href="`http://localhost:8082${pesanan.bukti_transfer}`" target="_blank" class="text-xs font-bold text-slate-600 bg-slate-100 px-3 py-2 rounded-lg hover:bg-slate-200 transition">Lihat Bukti TF</a>
              <button @click="ubahStatusPesanan(pesanan.id_pesanan, 'diproses')" class="text-xs font-bold text-white bg-emerald-500 px-3 py-2 rounded-lg hover:bg-emerald-600 transition shadow-sm">Terima & Proses</button>
            </div>

            <div v-else-if="pesanan.status === 'diproses'">
              <button @click="ubahStatusPesanan(pesanan.id_pesanan, 'dikirim')" class="text-xs font-bold text-white bg-indigo-600 px-4 py-2 rounded-lg hover:bg-indigo-700 transition shadow-sm">Tandai Dikirim</button>
            </div>

          </td>
        </tr>
        <tr v-if="pesananPerluTindakan.length === 0">
          <td colspan="5" class="px-6 py-12 text-center text-slate-500 font-medium">Tidak ada pesanan baru yang memerlukan tindakan.</td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
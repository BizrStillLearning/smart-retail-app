import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('token') || '')

    const namaTersimpan = localStorage.getItem('nama_lengkap')
    const namaLengkap = ref(namaTersimpan !== 'null' && namaTersimpan ? namaTersimpan : '')

    const role = ref(localStorage.getItem('role') || '')

    const setAuth = (dataToken, dataUser) => {
        token.value = dataToken
        namaLengkap.value = dataUser.nama_lengkap || 'Pengguna'
        role.value = dataUser.role

        localStorage.setItem('token', dataToken)
        localStorage.setItem('nama_lengkap', namaLengkap.value)
        localStorage.setItem('role', dataUser.role)
    }

    const logout = () => {
        token.value = ''
        namaLengkap.value = ''
        role.value = ''
        localStorage.clear()
    }

    const isAdmin = () => role.value === 'admin'
    const isAuthenticated = () => !!token.value

    return { token, namaLengkap, role, setAuth, logout, isAdmin, isAuthenticated }
})
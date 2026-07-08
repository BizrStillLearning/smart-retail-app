import { defineStore } from 'pinia'
import { computed } from 'vue'
import { useCookie } from 'nuxt/app'

interface CustomerData {
    id_pelanggan: number;
    nama_pelanggan: string;
}

export const useCustomerAuthStore = defineStore('customerAuth', () => {
    const token = useCookie<string | null>('customer_token', { maxAge: 60 * 60 * 24 * 3 })
    const customer = useCookie<CustomerData | null>('customer_data', { default: () => null })

    const isLoggedIn = computed(() => !!token.value)
    const customerName = computed(() => customer.value ? customer.value.nama_pelanggan : 'Tamu')

    const setAuth = (newToken: string, customerData: CustomerData) => {
        token.value = newToken
        customer.value = customerData
    }

    const logout = () => {
        token.value = null
        customer.value = null
    }

    return { token, customer, isLoggedIn, customerName, setAuth, logout }
})
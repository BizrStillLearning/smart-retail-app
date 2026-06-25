import { createRouter, createWebHistory } from 'vue-router'
import KatalogProduk from '../components/KatalogProduk.vue'
import AdminDashboard from '../views/AdminDashboard.vue'
import ManajemenProduk from '../views/ManajemenProduk.vue'
import ManajemenPesanan from '../views/ManajemenPesanan.vue'
import ManajemenPromo from '../views/ManajemenPromo.vue'
import Login from '../views/Login.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
    {
        path: '/',
        name: 'Toko',
        component: KatalogProduk
    },
    {
        path: '/login',
        name: 'Login',
        component: Login
    },
    {
        path: '/admin',
        name: 'Admin',
        component: AdminDashboard,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/produk',
        name: 'ManajemenProduk',
        component: ManajemenProduk,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/pesanan',
        name: 'ManajemenPesanan',
        component: ManajemenPesanan,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/promo',
        name: 'ManajemenPromo',
        component: ManajemenPromo,
        meta: { requiresAdmin: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    if (to.meta.requiresAdmin && !authStore.token) {
        next('/login')
    } else {
        next()
    }
})

export default router
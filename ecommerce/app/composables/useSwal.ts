import Swal from 'sweetalert2'

export const useSwal = () => {
    const Toast = Swal.mixin({
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        customClass: {
            popup: 'rounded-xl shadow-lg border border-slate-100 font-sans'
        },
        didOpen: (toast) => {
            toast.onmouseenter = Swal.stopTimer
            toast.onmouseleave = Swal.resumeTimer
        }
    })

    const sukses = (pesan) => Toast.fire({ icon: 'success', title: pesan })
    const error = (pesan) => Toast.fire({ icon: 'error', title: pesan })
    const info = (pesan) => Toast.fire({ icon: 'info', title: pesan })

    return { sukses, error, info }
}
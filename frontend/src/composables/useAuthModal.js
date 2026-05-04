import { ref } from 'vue'

const authModal = ref(null)

export function useAuthModal() {
  function openAuthModal(mode = 'login') {
    authModal.value = mode
  }

  function closeAuthModal() {
    authModal.value = null
  }

  return {
    authModal,
    openAuthModal,
    closeAuthModal,
  }
}

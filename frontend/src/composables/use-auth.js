import { ref } from 'vue';

/**
 * Composable для управления аутентификацией
 */
export function useAuth() {
  const isLoggedIn = ref(false);
  const isLoading = ref(false);

  function login() {
    isLoggedIn.value = true;
  }

  function logout() {
    isLoggedIn.value = false;
  }

  return {
    isLoggedIn,
    isLoading,
    login,
    logout
  };
}

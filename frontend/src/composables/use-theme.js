import { ref, onMounted } from 'vue';

/**
 * Composable для управления темой
 */
export function useTheme() {
  const isDarkTheme = ref(true);

  function updateBodyTheme() {
    if (isDarkTheme.value) {
      document.body.classList.add('th-dark');
      document.body.classList.remove('th-light');
    } else {
      document.body.classList.add('th-light');
      document.body.classList.remove('th-dark');
    }
  }

  function toggleTheme() {
    isDarkTheme.value = !isDarkTheme.value;
    updateBodyTheme();
  }

  function setTheme(dark) {
    isDarkTheme.value = dark;
    updateBodyTheme();
  }

  onMounted(() => {
    updateBodyTheme();
  });

  return {
    isDarkTheme,
    toggleTheme,
    setTheme
  };
}

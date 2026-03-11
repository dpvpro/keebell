import { ref } from 'vue';
import { OpenDatabase, GetEntries, GetGroups, CloseDatabase, IsDatabaseOpen, GetDatabaseName } from '../../wailsjs/go/main/App';

/**
 * Composable для управления базой данных KeePass
 */
export function useDatabase() {
  const entries = ref([]);
  const groups = ref([]);
  const databaseName = ref('');
  const isDatabaseOpen = ref(false);
  const isLoading = ref(false);
  const error = ref(null);

  async function openDatabase(path, password) {
    isLoading.value = true;
    error.value = null;

    try {
      const result = await OpenDatabase(path, password);

      if (result.success) {
        entries.value = result.database.allEntries || [];
        groups.value = result.database.groups || [];
        databaseName.value = result.database.name;
        isDatabaseOpen.value = true;
        return { success: true, data: result.database };
      } else {
        error.value = result.error;
        return { success: false, error: result.error };
      }
    } catch (e) {
      error.value = e.message;
      return { success: false, error: e.message };
    } finally {
      isLoading.value = false;
    }
  }

  function closeDB() {
    CloseDatabase();
    entries.value = [];
    groups.value = [];
    databaseName.value = '';
    isDatabaseOpen.value = false;
    error.value = null;
  }

  function getEntries() {
    return entries.value;
  }

  function getGroups() {
    return groups.value;
  }

  return {
    entries,
    groups,
    databaseName,
    isDatabaseOpen,
    isLoading,
    error,
    openDatabase,
    closeDB,
    getEntries,
    getGroups
  };
}

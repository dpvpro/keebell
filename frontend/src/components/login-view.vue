<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue';

const emit = defineEmits(['login', 'open-file']);

const password = ref('');
const showPassword = ref(false);
const isLoading = ref(false);
const errorMessage = ref('');
const showOptionsModal = ref(false);
const recentDatabases = ref([]);

const selectedDatabase = ref(null);
const passwordInputRef = ref(null);

// Load recent databases from localStorage
onMounted(() => {
  try {
    const stored = localStorage.getItem('keebell-recent-dbs');
    if (stored) {
      recentDatabases.value = JSON.parse(stored);
    }
  } catch (e) {
    console.error('Failed to load recent databases:', e);
  }
  
  // Listen for login errors from App.vue
  window.addEventListener('login-error', handleLoginError);
  // Listen for global escape
  window.addEventListener('global-escape', handleGlobalEscape);
});

function handleLoginError(event) {
  const message = event.detail?.message || 'Failed to open database';
  errorMessage.value = message;
  isLoading.value = false;
  // Focus password field for retry
  focusPassword();
}

function handleGlobalEscape() {
  // Close options modal
  if (showOptionsModal.value) {
    closeOptions();
  }
}

async function focusPassword() {
  await nextTick();
  if (passwordInputRef.value) {
    passwordInputRef.value.focus();
  }
}

function handleLogin() {
  if (!password.value.trim()) {
    errorMessage.value = 'Please enter your password';
    return;
  }
  
  if (!selectedDatabase.value) {
    errorMessage.value = 'Please select a database';
    return;
  }
  
  isLoading.value = true;
  errorMessage.value = '';
  
  // Emit login event with database and password
  emit('login', { database: selectedDatabase.value, password: password.value });
}

function handleKeyDown(event, action) {
  if (event.key === 'Enter') {
    event.preventDefault();
    action();
  }
}

function selectDatabase(db) {
  selectedDatabase.value = db;
  errorMessage.value = '';
  // Focus password field after selecting database
  focusPassword();
}

async function openDatabase() {
  try {
    const result = await window.Launcher?.getOpenFileName?.('', 'Open KeePass Database', 'KeePass Database', '*.kdbx');
    if (result && result.path) {
      const db = {
        id: Date.now(),
        name: result.path.split('/').pop() || result.path.split('\\').pop(),
        path: result.path,
        modified: 'Just now'
      };
      
      // Add to recent databases
      recentDatabases.value = [db, ...recentDatabases.value.filter(d => d.path !== db.path)].slice(0, 10);
      
      // Save to localStorage
      try {
        localStorage.setItem('keebell-recent-dbs', JSON.stringify(recentDatabases.value));
      } catch (e) {
        console.error('Failed to save recent databases:', e);
      }
      
      selectDatabase(db);
      emit('open-file', db);
    }
  } catch (e) {
    console.error('Failed to open database:', e);
    errorMessage.value = 'Failed to open database file';
  }
}

function createDatabase() {
  console.log('Create new database');
}

function openOptions() {
  showOptionsModal.value = true;
}

function closeOptions() {
  showOptionsModal.value = false;
}

function handleEscapeKey(event) {
  if (event.key === 'Escape') {
    if (showOptionsModal.value) {
      closeOptions();
    }
  }
}

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscapeKey);
  window.removeEventListener('login-error', handleLoginError);
  window.removeEventListener('global-escape', handleGlobalEscape);
});
</script>

<template>
  <div class="login-view">
    <div class="login-container">
      <!-- Logo and Title -->
      <div class="login-header">
        <div class="login-logo">
          <i class="fa fa-key"></i>
        </div>
        <h1 class="login-title">Keebell</h1>
        <p class="login-subtitle">Password Manager</p>
      </div>

      <!-- Top Action Buttons -->
      <div class="top-actions">
        <button
          class="top-action-btn"
          @click="openDatabase"
          @keydown="handleKeyDown($event, openDatabase)"
          title="Open Database"
        >
          <i class="fa fa-folder-open"></i>
          <span>Open</span>
        </button>
        <button
          class="top-action-btn"
          @click="createDatabase"
          @keydown="handleKeyDown($event, createDatabase)"
          title="Create New Database"
        >
          <i class="fa fa-plus"></i>
          <span>New</span>
        </button>
        <button
          class="top-action-btn"
          @click="openOptions"
          @keydown="handleKeyDown($event, openOptions)"
          title="More Options"
        >
          <i class="fa fa-ellipsis-h"></i>
          <span>More</span>
        </button>
      </div>

      <!-- Database Selection -->
      <div class="database-section">
        <h2 class="database-section-title">Recent Databases</h2>
        <div v-if="recentDatabases.length === 0" class="database-empty">
          <i class="fa fa-folder-open"></i>
          <p>No recent databases. Click "Open" to select a database file.</p>
        </div>
        <div v-else class="database-list">
          <div
            v-for="db in recentDatabases"
            :key="db.id"
            :class="['database-item', { 'database-item--selected': selectedDatabase?.id === db.id }]"
            @click="selectDatabase(db)"
            @keydown="handleKeyDown($event, () => selectDatabase(db))"
            role="button"
            tabindex="0"
          >
            <div class="database-item__icon">
              <i class="fa fa-file"></i>
            </div>
            <div class="database-item__content">
              <div class="database-item__name">{{ db.name }}</div>
              <div class="database-item__path">{{ db.path }}</div>
            </div>
            <div class="database-item__modified">{{ db.modified }}</div>
          </div>
        </div>
      </div>

      <!-- Password Input -->
      <div class="password-section">
        <div class="password-input-wrapper">
          <label for="password" class="password-label">Password</label>
          <div class="password-input-container">
            <input
              ref="passwordInputRef"
              id="password"
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              class="password-input"
              placeholder="Enter master password"
              @keydown="handleKeyDown($event, handleLogin)"
            />
            <button
              class="password-toggle"
              @click="showPassword = !showPassword"
              @keydown="handleKeyDown($event, () => showPassword = !showPassword)"
              :title="showPassword ? 'Hide' : 'Show'"
            >
              <i :class="showPassword ? 'fa fa-eye-slash' : 'fa fa-eye'"></i>
            </button>
          </div>
          <div v-if="errorMessage" class="password-error">
            <i class="fa fa-exclamation-circle"></i>
            <span>{{ errorMessage }}</span>
          </div>
        </div>

        <button
          class="login-btn"
          :class="{ 'login-btn--loading': isLoading }"
          @click="handleLogin"
          @keydown="handleKeyDown($event, handleLogin)"
          :disabled="isLoading"
        >
          <i v-if="isLoading" class="fa fa-spinner fa-spin"></i>
          <i v-else class="fa fa-unlock-alt"></i>
          <span>{{ isLoading ? 'Unlocking...' : 'Unlock Database' }}</span>
        </button>
      </div>

    </div>

    <!-- Options Modal -->
    <div v-if="showOptionsModal" class="modal-overlay" @click="closeOptions">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2 class="modal__title">Options</h2>
          <button class="modal__close" @click="closeOptions" title="Close">
            <i class="fa fa-times"></i>
          </button>
        </div>
        <div class="modal__body">
          <div class="options-list">
            <button class="options-list__item" @click="console.log('Settings')">
              <i class="fa fa-cog"></i>
              <span>Settings</span>
            </button>
            <button class="options-list__item" @click="console.log('Change theme')">
              <i class="fa fa-palette"></i>
              <span>Change Theme</span>
            </button>
            <button class="options-list__item" @click="console.log('About')">
              <i class="fa fa-info-circle"></i>
              <span>About Keebell</span>
            </button>
            <button class="options-list__item" @click="console.log('Lock all')">
              <i class="fa fa-lock"></i>
              <span>Lock All Databases</span>
            </button>
          </div>
        </div>
        <div class="modal__footer">
          <button class="modal__btn modal__btn--primary" @click="closeOptions">
            Close
          </button>
        </div>
      </div>
    </div>

    <!-- Background Decoration -->
    <div class="login-background">
      <div class="login-bg-circle login-bg-circle--1"></div>
      <div class="login-bg-circle login-bg-circle--2"></div>
      <div class="login-bg-circle login-bg-circle--3"></div>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  position: relative;
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--background-color);
  overflow: hidden;
}

.login-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 480px;
  padding: 40px;
}

/* Header */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-logo {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, var(--action-color), var(--selected-item-color));
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(49, 126, 246, 0.3);
}

.login-logo i {
  font-size: 40px;
  color: white;
}

.login-title {
  font-size: 32px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--text-color);
}

.login-subtitle {
  font-size: 14px;
  color: var(--muted-color);
  margin: 0;
}

/* Top Actions */
.top-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 24px;
}

.top-action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-size: 13px;
  font-weight: 500;
  background-color: var(--secondary-background-color);
  border: 1px solid var(--base-border-color);
  border-radius: var(--input-border-radius);
  color: var(--text-color);
  cursor: pointer;
  transition: all var(--fast-duration) var(--base-timing);
}

.top-action-btn:hover {
  background-color: var(--hover-background-color);
  border-color: var(--action-color);
}

.top-action-btn i {
  font-size: 14px;
}

/* Database Section */
.database-section {
  margin-bottom: 30px;
  text-align: center;
}

.database-section-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--muted-color);
  margin: 0 0 12px 0;
}

.database-list {
  background-color: var(--secondary-background-color);
  border: 1px solid var(--base-border-color);
  border-radius: var(--base-border-radius);
  overflow: hidden;
  max-width: 400px;
  margin: 0 auto;
}

.database-empty {
  max-width: 400px;
  margin: 0 auto;
  padding: 40px 20px;
  text-align: center;
  color: var(--muted-color);
}

.database-empty i {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.database-empty p {
  margin: 0;
  font-size: 13px;
}

.database-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color var(--fast-duration) var(--base-timing);
  border-bottom: 1px solid var(--light-border-color);
}

.database-item:last-child {
  border-bottom: none;
}

.database-item:hover {
  background-color: var(--hover-background-color);
}

.database-item--selected {
  background-color: var(--selected-item-color) !important;
}

.database-item__icon {
  width: 40px;
  height: 40px;
  background-color: var(--intermediate-background-color);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.database-item__icon i {
  font-size: 18px;
  color: var(--action-color);
}

.database-item__content {
  flex: 1;
  min-width: 0;
}

.database-item__name {
  font-weight: 500;
  font-size: 13px;
  color: var(--text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.database-item__path {
  font-size: 11px;
  color: var(--muted-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}

.database-item__modified {
  font-size: 11px;
  color: var(--muted-color);
  white-space: nowrap;
}

/* Password Section */
.password-section {
  margin-bottom: 24px;
}

.password-input-wrapper {
  margin-bottom: 16px;
}

.password-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--muted-color);
  margin-bottom: 8px;
}

.password-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input {
  flex: 1;
  width: 100%;
  padding: 12px 44px 12px 14px;
  font-size: 14px;
  background-color: var(--secondary-background-color);
  border: 1px solid var(--base-border-color);
  border-radius: var(--input-border-radius);
  color: var(--text-color);
  transition: border-color var(--fast-duration) var(--base-timing),
              box-shadow var(--fast-duration) var(--base-timing);
}

.password-input:focus {
  outline: none;
  border-color: var(--form-box-border-color-focus);
  box-shadow: 0 0 0 var(--focus-shadow-spread) var(--form-box-shadow-color-focus);
}

.password-input::placeholder {
  color: var(--muted-color);
}

.password-toggle {
  position: absolute;
  right: 4px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--muted-color);
  cursor: pointer;
  padding: 6px 8px;
  border-radius: var(--button-border-radius);
  transition: all var(--fast-duration) var(--base-timing);
}

.password-toggle:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.password-error {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--error-color);
}

.password-error i {
  font-size: 14px;
}

.login-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 14px 20px;
  font-size: 14px;
  font-weight: 500;
  background-color: var(--action-color);
  border: 1px solid var(--action-color);
  border-radius: var(--input-border-radius);
  color: white;
  cursor: pointer;
  transition: all var(--fast-duration) var(--base-timing);
}

.login-btn:hover:not(:disabled) {
  background-color: var(--action-background-color-focus);
}

.login-btn:active:not(:disabled) {
  background-color: var(--action-background-color-active);
}

.login-btn--loading {
  opacity: 0.8;
  cursor: wait;
}

.login-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.login-btn i {
  font-size: 16px;
}

/* Footer */
.login-footer {
  display: flex;
  justify-content: center;
  gap: 24px;
}

.login-link {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--muted-color);
  transition: color var(--fast-duration) var(--base-timing);
}

.login-link:hover {
  color: var(--action-color);
  text-decoration: none;
}

.login-link i {
  font-size: 14px;
}

/* Background Decoration */
.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
  pointer-events: none;
}

.login-bg-circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--action-color), transparent);
  opacity: 0.1;
}

.login-bg-circle--1 {
  width: 400px;
  height: 400px;
  top: -100px;
  right: -100px;
}

.login-bg-circle--2 {
  width: 300px;
  height: 300px;
  bottom: -50px;
  left: -50px;
}

.login-bg-circle--3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

/* Responsive */
@media (max-width: 600px) {
  .login-container {
    padding: 24px;
  }

  .login-logo {
    width: 60px;
    height: 60px;
    border-radius: 15px;
  }

  .login-logo i {
    font-size: 30px;
  }

  .login-title {
    font-size: 24px;
  }

  .database-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .database-item__icon {
    width: 36px;
    height: 36px;
  }
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-index-modal);
}

.modal {
  background-color: var(--background-color);
  border: 1px solid var(--base-border-color);
  border-radius: var(--base-border-radius);
  width: 100%;
  max-width: 400px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
}

.modal__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--base-border-color);
}

.modal__title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: var(--text-color);
}

.modal__close {
  background: none;
  border: none;
  color: var(--muted-color);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--button-border-radius);
  font-size: 16px;
  transition: all var(--fast-duration) var(--base-timing);
}

.modal__close:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.modal__body {
  padding: 20px;
}

.modal__footer {
  padding: 16px 20px;
  border-top: 1px solid var(--base-border-color);
  display: flex;
  justify-content: flex-end;
}

.modal__btn {
  padding: 8px 20px;
  font-size: 13px;
  font-weight: 500;
  border-radius: var(--button-border-radius);
  cursor: pointer;
  transition: all var(--fast-duration) var(--base-timing);
}

.modal__btn--primary {
  background-color: var(--action-color);
  border: 1px solid var(--action-color);
  color: white;
}

.modal__btn--primary:hover {
  background-color: var(--action-background-color-focus);
}

/* Options List */
.options-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.options-list__item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 16px;
  background-color: var(--secondary-background-color);
  border: 1px solid var(--base-border-color);
  border-radius: var(--base-border-radius);
  color: var(--text-color);
  font-size: 13px;
  cursor: pointer;
  transition: all var(--fast-duration) var(--base-timing);
  text-align: left;
}

.options-list__item:hover {
  background-color: var(--hover-background-color);
  border-color: var(--action-color);
}

.options-list__item i {
  font-size: 16px;
  width: 20px;
  text-align: center;
  color: var(--action-color);
}
</style>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { Version } from '../wailsjs/go/main/App';
import LoginView from './components/LoginView.vue';
import EntryTable from './components/EntryTable.vue';
import { useDatabase, useAuth, useTheme } from './composables';
import { copyToClipboard, openUrl } from './utils/clipboard';

// Composables
const {
  entries,
  groups,
  databaseName,
  isDatabaseOpen,
  isLoading,
  error: dbError,
  openDatabase,
  closeDB
} = useDatabase();

const {
  isLoggedIn,
  login,
  logout
} = useAuth();

const {
  isDarkTheme,
  toggleTheme
} = useTheme();

// Local state
const selectedEntry = ref(null);
const showPassword = ref(false);
const showVersionModal = ref(false);
const appVersion = ref('');

function handleKeyDown(event, action) {
  // Handle Enter or Space key for accessibility
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    action();
  }
}

function getEntryIconClass(entry) {
  if (entry.color) {
    return `fa fa-${entry.icon} ${entry.color}-color`;
  }
  return `fa fa-${entry.icon}`;
}

async function showVersionInfo() {
  try {
    const version = await Version();
    appVersion.value = version;
    showVersionModal.value = true;
  } catch (err) {
    console.error('Failed to get version:', err);
    appVersion.value = 'Unknown';
    showVersionModal.value = true;
  }
}

function closeVersionModal() {
  showVersionModal.value = false;
}

async function handleLogin(data) {
  console.log('Login with database:', data.database, 'Password:', data.password ? '***' : '');

  if (!data.database || !data.database.path || !data.password) {
    emitLoginError('Please select a database and enter password');
    return;
  }

  const result = await openDatabase(data.database.path, data.password);

  if (result.success) {
    login();
    console.log('Database opened successfully:', entries.value.length, 'entries');
  } else {
    emitLoginError(result.error || 'Failed to open database');
  }
}

function emitLoginError(message) {
  window.dispatchEvent(new CustomEvent('login-error', { detail: { message } }));
}

function handleOpenFile(db) {
  console.log('Database selected:', db);
}

function handleLogout() {
  closeDB();
  logout();
  selectedEntry.value = null;
}

function handleEscapeKey(event) {
  if (event.key === 'Escape') {
    if (showVersionModal.value) {
      closeVersionModal();
    }
    window.dispatchEvent(new CustomEvent('global-escape'));
    if (document.activeElement && document.activeElement.tagName !== 'BODY') {
      document.activeElement.blur();
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscapeKey);
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscapeKey);
});
</script>

<template>
  <div class="app">
    <!-- Titlebar drag area (for frameless windows) -->
    <div class="app__titlebar-drag"></div>

    <!-- Login View -->
    <LoginView
      v-if="!isLoggedIn"
      @login="handleLogin"
      @open-file="handleOpenFile"
    />

    <!-- Main App View -->
    <div v-else class="app__body-wrapper">
      <!-- Entry Table -->
      <EntryTable
        v-if="isDatabaseOpen"
        :entries="entries"
        @select-entry="selectEntry"
        @copy-username="copyToClipboard"
        @copy-password="copyToClipboard"
        @open-url="openUrl"
        @show-version="showVersionInfo"
        @show-help="console.log('Help clicked')"
        @show-settings="console.log('Settings clicked')"
        @generate-password="console.log('Generate password clicked')"
        @logout="handleLogout"
      />

      <!-- Empty state when no database is open -->
      <div v-else class="empty-state">
        <i class="fa fa-key"></i>
        <h2>No Database Open</h2>
        <p>Open a KeePass database file to view entries</p>
      </div>
    </div>

    <!-- Version Modal -->
    <div v-if="showVersionModal" class="modal-overlay" @click="closeVersionModal">
      <div class="modal" @click.stop>
        <div class="modal__header">
          <h2 class="modal__title">KeeWeb</h2>
          <button class="modal__close" @click="closeVersionModal" title="Close">
            <i class="fa fa-times"></i>
          </button>
        </div>
        <div class="modal__body">
          <div class="version-info">
            <div class="version-info__row">
              <span class="version-info__label">Version:</span>
              <span class="version-info__value">{{ appVersion }}</span>
            </div>
          </div>
        </div>
        <div class="modal__footer">
          <button class="modal__btn modal__btn--primary" @click="closeVersionModal">
            OK
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
/* Theme variables are defined in style.css */

/* Base styles */
:global(body) {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica Neue, Helvetica, Roboto, Arial, sans-serif;
  font-size: 12px;
  line-height: 1.5;
  overflow: hidden;
  height: 100vh;
  background-color: var(--background-color);
  color: var(--text-color);
}

/* App container */
.app {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.app__titlebar-drag {
  position: absolute;
  width: 100%;
  height: 30px;
  top: 0;
  right: 0;
  pointer-events: none;
}

.app__body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* Empty State */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--muted-color);
  text-align: center;
}

.empty-state i {
  font-size: 64px;
  margin-bottom: 24px;
  opacity: 0.5;
}

.empty-state h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: var(--text-color);
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

/* Menu sidebar */
.app__menu {
  flex: 0 0 auto;
  width: 15em;
  background-color: var(--secondary-background-color);
  display: flex;
  flex-direction: column;
}

.app__menu-drag {
  flex: 0 0 auto;
  width: 1px;
  background-color: var(--base-border-color);
  cursor: col-resize;
}

.menu {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.menu__header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--base-border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.menu__header-title {
  font-weight: 600;
  font-size: 14px;
}

.menu__header-button {
  background: none;
  border: none;
  color: var(--muted-color);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 3px;
}

.menu__header-button:hover {
  background-color: var(--hover-background-color);
}

.menu__sections {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
}

.menu__section {
  padding: 0 8px;
}

.menu__section-header {
  padding: 8px 8px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border-radius: 3px;
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  font: inherit;
  color: inherit;
}

.menu__section-header:hover {
  background-color: var(--hover-background-color);
}

.menu__section-icon {
  margin-right: 8px;
  width: 16px;
  text-align: center;
  color: var(--muted-color);
}

.menu__section-name {
  flex: 1;
}

.menu__section-count {
  background-color: var(--action-color);
  color: white;
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.menu__footer {
  border-top: 1px solid var(--base-border-color);
  padding: 8px 0;
}

.menu__footer-item {
  padding: 8px 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  color: var(--muted-color);
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  font: inherit;
}

.menu__footer-item:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.menu__footer-item i {
  margin-right: 8px;
  width: 16px;
  text-align: center;
}

/* List area */
.app__list-wrap {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.app__list {
  flex: 0 0 25em;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--background-color);
}

.app__list-drag {
  flex: 0 0 auto;
  width: 1px;
  background-color: var(--base-border-color);
  cursor: col-resize;
}

.list {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.list__header {
  padding: 12px;
  border-bottom: 1px solid var(--base-border-color);
}

.list__search {
  position: relative;
  margin-bottom: 8px;
}

.list__search-icon {
  position: absolute;
  left: 8px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--muted-color);
}

.list__search-input {
  width: 100%;
  padding: 6px 8px 6px 28px;
  border: 1px solid var(--base-border-color);
  border-radius: 4px;
  background-color: var(--secondary-background-color);
  color: var(--text-color);
  font-size: 12px;
}

.list__search-input:focus {
  outline: none;
  border-color: var(--action-color);
}

.list__controls {
  display: flex;
  gap: 4px;
}

.list__control-button {
  background: none;
  border: 1px solid var(--base-border-color);
  color: var(--muted-color);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 11px;
}

.list__control-button:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.list__items {
  flex: 1;
  overflow-y: auto;
}

.list__item {
  padding: 12px;
  border-bottom: 1px solid var(--light-border-color);
  display: flex;
  align-items: center;
  cursor: pointer;
}

.list__item:hover {
  background-color: var(--hover-background-color);
}

.list__item--active {
  background-color: var(--selected-item-color) !important;
}

.list__item-icon {
  margin-right: 12px;
  font-size: 16px;
  width: 24px;
  text-align: center;
}

.list__item-content {
  flex: 1;
  min-width: 0;
}

.list__item-title {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.list__item-subtitle {
  font-size: 11px;
  color: var(--muted-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.list__item-action {
  color: var(--muted-color);
  padding: 4px;
  border-radius: 3px;
}

.list__item-action:hover {
  background-color: var(--intermediate-background-color);
  color: var(--text-color);
}

/* Details area */
.app__details {
  flex: 1;
  display: flex;
  overflow: hidden;
  background-color: var(--background-color);
}

.details {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  overflow: hidden;
}

.details--empty {
  align-items: center;
  justify-content: center;
}

.details__empty-message {
  text-align: center;
  color: var(--muted-color);
}

.details__empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.details__empty-message h2 {
  margin: 8px 0;
  font-weight: 500;
}

.details__empty-message p {
  margin: 0;
}

.details__header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--base-border-color);
}

.details__header-color {
  font-size: 24px;
  margin-right: 12px;
}

.details__header-title {
  margin: 0;
  font-weight: 500;
  font-size: 20px;
  flex: 1;
}

.details__body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.scroller {
  flex: 1;
  overflow-y: auto;
  padding-right: 8px;
}

.details__body-fields {
  margin-bottom: 24px;
}

.details__field {
  margin-bottom: 16px;
}

.details__field-label {
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  color: var(--muted-color);
  margin-bottom: 4px;
}

.details__field-value {
  display: flex;
  gap: 4px;
}

.details__field-value input,
.details__field-value textarea {
  flex: 1;
  padding: 8px;
  border: 1px solid var(--base-border-color);
  border-radius: 4px;
  background-color: var(--secondary-background-color);
  color: var(--text-color);
  font-family: inherit;
  font-size: 12px;
}

.details__field-value textarea {
  resize: vertical;
  min-height: 60px;
}

.details__field-value input:focus,
.details__field-value textarea:focus {
  outline: none;
  border-color: var(--action-color);
}

.details__field-copy,
.details__field-show,
.details__field-open {
  background: none;
  border: 1px solid var(--base-border-color);
  color: var(--muted-color);
  cursor: pointer;
  padding: 6px 10px;
  border-radius: 3px;
  font-size: 12px;
}

.details__field-copy:hover,
.details__field-show:hover,
.details__field-open:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.details__body-aside {
  flex: 0 0 200px;
  margin-left: 24px;
  border-left: 1px solid var(--base-border-color);
  padding-left: 24px;
}

.details__aside-section {
  margin-bottom: 24px;
}

.details__aside-title {
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
  color: var(--muted-color);
  margin: 0 0 8px 0;
}

.details__aside-button {
  width: 100%;
  padding: 8px;
  border: 1px solid var(--base-border-color);
  background-color: var(--secondary-background-color);
  color: var(--text-color);
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.details__aside-button:hover {
  background-color: var(--hover-background-color);
}

.details__aside-text {
  font-size: 11px;
  color: var(--muted-color);
  margin-bottom: 4px;
}

.details__buttons {
  display: flex;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid var(--base-border-color);
}

.details__button {
  padding: 8px 16px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.details__button--edit {
  background-color: var(--action-color);
  color: white;
}

.details__button--edit:hover {
  background-color: #1a64d6;
}

.details__button--delete {
  background-color: var(--error-color);
  color: white;
}

.details__button--delete:hover {
  background-color: #d84c43;
}

/* Footer */
.app__footer {
  flex: 0 0 auto;
  border-top: 1px solid var(--base-border-color);
  background-color: var(--secondary-background-color);
}

.footer {
  display: flex;
  align-items: center;
  padding: 0 8px;
  height: 40px;
}

.footer__db {
  padding: 6px 12px;
  margin: 0 4px;
  border-radius: 3px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  max-width: 150px;
  overflow: hidden;
}

.footer__db:hover {
  background-color: var(--hover-background-color);
}

.footer__db--dimmed {
  opacity: 0.6;
}

.footer__db-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.footer__db-sign {
  color: var(--action-color);
  font-size: 8px;
  margin-left: 4px;
}

.footer__db-open {
  color: var(--action-color);
}

.footer__db-text {
  margin-left: 4px;
}

.footer__btn {
  margin-left: auto;
  padding: 8px 12px;
  cursor: pointer;
  color: var(--muted-color);
  border-radius: 3px;
  background: none;
  border: none;
  font: inherit;
}

.footer__btn:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

/* Scrollbar styling */
.scroller::-webkit-scrollbar {
  width: 8px;
}

.scroller::-webkit-scrollbar-track {
  background: var(--secondary-background-color);
}

.scroller::-webkit-scrollbar-thumb {
  background: var(--muted-color);
  border-radius: 4px;
}

.scroller::-webkit-scrollbar-thumb:hover {
  background: #888;
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .app__menu {
    display: none;
  }

  .app__menu-drag {
    display: none;
  }

  .app__list {
    flex: 1;
  }

  .app__list-drag {
    display: none;
  }

  .app__details {
    display: none;
  }

  .details__body-aside {
    display: none;
  }
}

/* Modal styles */
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
  background-color: var(--secondary-background-color);
  border-radius: var(--block-border-radius);
  min-width: 300px;
  max-width: 400px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.modal__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--base-border-color);
}

.modal__title {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.modal__close {
  background: none;
  border: none;
  color: var(--muted-color);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 16px;
}

.modal__close:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.modal__body {
  padding: 20px;
}

.version-info__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.version-info__row:last-child {
  margin-bottom: 0;
}

.version-info__label {
  font-weight: 500;
  color: var(--muted-color);
}

.version-info__value {
  font-family: var(--monospace-font-family);
  color: var(--text-color);
}

.modal__footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid var(--base-border-color);
}

.modal__btn {
  padding: 8px 20px;
  border-radius: var(--button-border-radius);
  cursor: pointer;
  font-size: 12px;
  border: 1px solid var(--base-border-color);
  background-color: var(--secondary-background-color);
  color: var(--text-color);
}

.modal__btn:hover {
  background-color: var(--hover-background-color);
}

.modal__btn--primary {
  background-color: var(--action-color);
  border-color: var(--action-color);
  color: white;
}

.modal__btn--primary:hover {
  background-color: var(--action-background-color-focus);
  border-color: var(--action-background-color-focus);
}
</style>

<script setup>
import { ref, onMounted } from 'vue';
import { Version } from '../wailsjs/go/main/App';

// State
const selectedEntry = ref(null);
const isDarkTheme = ref(true);
const showPassword = ref(false);
const showVersionModal = ref(false);
const appVersion = ref('');

const entries = ref([
  { id: 1, title: 'Google Account', username: 'user@gmail.com', icon: 'globe', color: 'blue', url: 'https://google.com' },
  { id: 2, title: 'GitHub', username: 'developer', icon: 'github', color: 'violet', url: 'https://github.com' },
  { id: 3, title: 'AWS Console', username: 'admin', icon: 'cloud', color: 'orange', url: 'https://aws.amazon.com' },
  { id: 4, title: 'Dropbox', username: 'user@example.com', icon: 'dropbox', color: 'blue', url: 'https://dropbox.com' },
  { id: 5, title: 'Local Database', username: 'admin', icon: 'database', color: 'green', url: '' },
  { id: 6, title: 'SSH Server', username: 'root', icon: 'server', color: 'red', url: '' },
  { id: 7, title: 'WiFi Password', username: '', icon: 'wifi', color: 'green', url: '' },
  { id: 8, title: 'Bank Account', username: '123456789', icon: 'university', color: 'yellow', url: 'https://bank.com' },
]);

const groups = ref([
  { id: 1, name: 'Internet', icon: 'globe', count: 4 },
  { id: 2, name: 'Work', icon: 'briefcase', count: 2 },
  { id: 3, name: 'Personal', icon: 'user', count: 2 },
  { id: 4, name: 'Archived', icon: 'archive', count: 0 },
]);

const openFiles = ref([
  { id: 1, name: 'Personal.kdbx', active: true, modified: false },
  { id: 2, name: 'Work.kdbx', active: false, modified: true },
]);

// Update body theme class
function updateBodyTheme() {
  if (isDarkTheme.value) {
    document.body.classList.add('th-dark');
    document.body.classList.remove('th-light');
  } else {
    document.body.classList.add('th-light');
    document.body.classList.remove('th-dark');
  }
}

// Initialize with first entry selected and set theme
onMounted(() => {
  if (entries.value.length > 0 && !selectedEntry.value) {
    selectedEntry.value = entries.value[0];
  }
  updateBodyTheme();
});

function selectEntry(entry) {
  selectedEntry.value = entry;
}

function toggleTheme() {
  isDarkTheme.value = !isDarkTheme.value;
  updateBodyTheme();
}

function copyToClipboard(text) {
  if (!text) return;
  // Use Launcher.setClipboardText if available
  if (window.Launcher) {
    window.Launcher.setClipboardText(text);
  } else {
    // Fallback to navigator clipboard API
    navigator.clipboard.writeText(text).catch(err => {
      console.error('Failed to copy:', err);
    });
  }
  // Show some feedback (could be enhanced with a toast notification)
  console.log('Copied to clipboard:', text.substring(0, 20) + '...');
}

function togglePasswordVisibility() {
  showPassword.value = !showPassword.value;
}

function openUrl(url) {
  if (!url) return;
  // Use Launcher.openLink if available
  if (window.Launcher) {
    window.Launcher.openLink(url);
  } else {
    // Fallback to window.open
    window.open(url, '_blank');
  }
}

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
</script>

<template>
  <div class="app">
    <!-- Titlebar drag area (for frameless windows) -->
    <div class="app__titlebar-drag"></div>

    <div class="app__body">
      <!-- Left sidebar menu -->
      <div class="app__menu">
        <div class="menu">
          <div class="menu__header">
            <div class="menu__header-title">KeeWeb</div>
            <button class="menu__header-button" @click="toggleTheme" title="Toggle theme">
              <i :class="`fa fa-${isDarkTheme ? 'sun' : 'moon'}`"></i>
            </button>
          </div>

          <div class="menu__sections">
            <div v-for="group in groups" :key="group.id" class="menu__section">
              <button
                class="menu__section-header"
                @click="console.log('Group clicked:', group.name)"
                @keydown="handleKeyDown($event, () => console.log('Group clicked:', group.name))"
              >
                <i :class="`fa fa-${group.icon} menu__section-icon`"></i>
                <span class="menu__section-name">{{ group.name }}</span>
                <span v-if="group.count > 0" class="menu__section-count">{{ group.count }}</span>
              </button>
            </div>
          </div>

          <div class="menu__footer">
            <button
              class="menu__footer-item"
              @click="console.log('New Group clicked')"
              @keydown="handleKeyDown($event, () => console.log('New Group clicked'))"
            >
              <i class="fa fa-plus"></i>
              <span>New Group</span>
            </button>
            <button
              class="menu__footer-item"
              @click="console.log('Search clicked')"
              @keydown="handleKeyDown($event, () => console.log('Search clicked'))"
            >
              <i class="fa fa-search"></i>
              <span>Search</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Resize handle -->
      <div class="app__menu-drag"></div>

      <!-- Main content area -->
      <div class="app__list-wrap">
        <!-- Entries list -->
        <div class="app__list">
          <div class="list">
            <div class="list__header">
              <div class="list__search">
                <i class="fa fa-search list__search-icon"></i>
                <input type="text" class="list__search-input" placeholder="Search entries..." />
              </div>
              <div class="list__controls">
                <button
                  class="list__control-button"
                  title="New Entry"
                  @click="console.log('New Entry clicked')"
                  @keydown="handleKeyDown($event, () => console.log('New Entry clicked'))"
                >
                  <i class="fa fa-plus"></i>
                </button>
                <button
                  class="list__control-button"
                  title="Sort"
                  @click="console.log('Sort clicked')"
                  @keydown="handleKeyDown($event, () => console.log('Sort clicked'))"
                >
                  <i class="fa fa-sort"></i>
                </button>
              </div>
            </div>

            <div class="list__items">
              <div
                v-for="entry in entries"
                :key="entry.id"
                :class="['list__item', { 'list__item--active': selectedEntry?.id === entry.id }]"
                @click="selectEntry(entry)"
                @keydown="handleKeyDown($event, () => selectEntry(entry))"
                role="button"
                tabindex="0"
              >
                <div class="list__item-icon">
                  <i :class="getEntryIconClass(entry)"></i>
                </div>
                <div class="list__item-content">
                  <div class="list__item-title">{{ entry.title }}</div>
                  <div v-if="entry.username" class="list__item-subtitle">{{ entry.username }}</div>
                </div>
                <div
                  v-if="entry.url"
                  class="list__item-action"
                  @click.stop="openUrl(selectedEntry.url)"
                  @keydown="handleKeyDown($event, () => openUrl(selectedEntry.url))"
                  title="Open URL"
                >
                  <i class="fa fa-external-link"></i>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Resize handle -->
        <div class="app__list-drag"></div>

        <!-- Entry details -->
        <div class="app__details">
          <div v-if="selectedEntry" class="details">
            <div class="details__header">
              <i :class="`details__header-color ${getEntryIconClass(selectedEntry)}`"></i>
              <h1 class="details__header-title">{{ selectedEntry.title }}</h1>
            </div>

            <div class="details__body">
              <div class="scroller">
                <div class="details__body-fields">
                  <div class="details__field">
                    <div class="details__field-label">Username</div>
                    <div class="details__field-value">
                      <input type="text" :value="selectedEntry.username || ''" readonly />
                      <button
                        class="details__field-copy"
                        title="Copy"
                        @click="copyToClipboard(selectedEntry.username)"
                        @keydown="handleKeyDown($event, () => copyToClipboard(selectedEntry.username))"
                      >
                        <i class="fa fa-copy"></i>
                      </button>
                    </div>
                  </div>

                  <div class="details__field">
                    <div class="details__field-label">Password</div>
                    <div class="details__field-value">
                      <input
                        :type="showPassword ? 'text' : 'password'"
                        value="password123"
                        readonly
                      />
                      <button
                        class="details__field-copy"
                        title="Copy"
                        @click="copyToClipboard('password123')"
                        @keydown="handleKeyDown($event, () => copyToClipboard('password123'))"
                      >
                        <i class="fa fa-copy"></i>
                      </button>
                      <button
                        class="details__field-show"
                        :title="showPassword ? 'Hide' : 'Show'"
                        @click="togglePasswordVisibility"
                        @keydown="handleKeyDown($event, togglePasswordVisibility)"
                      >
                        <i :class="`fa fa-eye${showPassword ? '-slash' : ''}`"></i>
                      </button>
                    </div>
                  </div>

                  <div v-if="selectedEntry.url" class="details__field">
                    <div class="details__field-label">URL</div>
                    <div class="details__field-value">
                      <input type="text" :value="selectedEntry.url" readonly />
                      <button
                        class="details__field-copy"
                        title="Copy"
                        @click="copyToClipboard(selectedEntry.url)"
                        @keydown="handleKeyDown($event, () => copyToClipboard(selectedEntry.url))"
                      >
                        <i class="fa fa-copy"></i>
                      </button>
                      <button
                        class="details__field-open"
                        title="Open"
                        @click="openUrl(selectedEntry.url)"
                        @keydown="handleKeyDown($event, () => openUrl(selectedEntry.url))"
                      >
                        <i class="fa fa-external-link"></i>
                      </button>
                    </div>
                  </div>

                  <div class="details__field">
                    <div class="details__field-label">Notes</div>
                    <div class="details__field-value">
                      <textarea readonly>This is a sample entry for demonstration purposes. The actual KeeWeb application would show real password entries here.</textarea>
                    </div>
                  </div>
                </div>

                <div class="details__body-aside">
                  <div class="details__aside-section">
                    <h3 class="details__aside-title">Auto-Type</h3>
                    <button
                      class="details__aside-button"
                      @click="console.log('Auto-type would be performed here')"
                      @keydown="handleKeyDown($event, () => console.log('Auto-type would be performed here'))"
                    >
                      <i class="fa fa-keyboard"></i>
                      <span>Perform Auto-Type</span>
                    </button>
                  </div>

                  <div class="details__aside-section">
                    <h3 class="details__aside-title">History</h3>
                    <div class="details__aside-text">Modified 2 days ago</div>
                    <div class="details__aside-text">Created 1 week ago</div>
                  </div>
                </div>
              </div>
            </div>

            <div class="details__buttons">
              <button
                class="details__button details__button--edit"
                title="Edit"
                @click="console.log('Edit entry')"
                @keydown="handleKeyDown($event, () => console.log('Edit entry'))"
              >
                <i class="fa fa-edit"></i>
              </button>
              <button
                class="details__button details__button--delete"
                title="Delete"
                @click="console.log('Delete entry')"
                @keydown="handleKeyDown($event, () => console.log('Delete entry'))"
              >
                <i class="fa fa-trash-alt"></i>
              </button>
            </div>
          </div>
          <div v-else class="details details--empty">
            <div class="details__empty-message">
              <i class="fa fa-key details__empty-icon"></i>
              <h2>No Entry Selected</h2>
              <p>Select an entry from the list to view its details</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="app__footer">
      <div class="footer">
        <div
          v-for="file in openFiles"
          :key="file.id"
          :class="['footer__db', { 'footer__db--dimmed': !file.active }]"
          :title="file.name"
        >
          <i :class="`fa fa-${file.active ? 'unlock' : 'lock'}`"></i>
          <span class="footer__db-name">{{ file.name }}</span>
          <i v-if="file.modified" class="fa fa-circle footer__db-sign"></i>
        </div>

        <div class="footer__db footer__db--dimmed footer__db-open">
          <i class="fa fa-plus"></i>
          <span class="footer__db-text">Open Database</span>
        </div>

        <button
          class="footer__btn"
          title="Version"
          @click="showVersionInfo"
          @keydown="handleKeyDown($event, showVersionInfo)"
        >
          <i class="fa fa-info-circle"></i>
        </button>

        <button
          class="footer__btn"
          title="Help"
          @click="console.log('Help clicked')"
          @keydown="handleKeyDown($event, () => console.log('Help clicked'))"
        >
          <i class="fa fa-question"></i>
        </button>

        <button
          class="footer__btn"
          title="Settings"
          @click="console.log('Settings clicked')"
          @keydown="handleKeyDown($event, () => console.log('Settings clicked'))"
        >
          <i class="fa fa-cog"></i>
        </button>

        <button
          class="footer__btn"
          title="Generate Password"
          @click="console.log('Generate password clicked')"
          @keydown="handleKeyDown($event, () => console.log('Generate password clicked'))"
        >
          <i class="fa fa-bolt"></i>
        </button>

        <button
          class="footer__btn"
          title="Lock Database"
          @click="console.log('Lock database clicked')"
          @keydown="handleKeyDown($event, () => console.log('Lock database clicked'))"
        >
          <i class="fa fa-sign-out-alt"></i>
        </button>
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

<script setup>
import { ref, computed } from 'vue';

const props = defineProps({
  entries: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['select-entry', 'copy-username', 'copy-password', 'open-url', 'show-version', 'show-help', 'show-settings', 'generate-password', 'logout']);

const selectedEntry = ref(null);
const searchTerm = ref('');
const showPassword = ref(false);

const filteredEntries = computed(() => {
  if (!searchTerm.value.trim()) {
    return props.entries;
  }
  const term = searchTerm.value.toLowerCase();
  return props.entries.filter(entry => {
    return (
      entry.title?.toLowerCase().includes(term) ||
      entry.userName?.toLowerCase().includes(term) ||
      entry.url?.toLowerCase().includes(term) ||
      entry.notes?.toLowerCase().includes(term)
    );
  });
});

function selectEntry(entry) {
  selectedEntry.value = entry;
  emit('select-entry', entry);
}

function copyToClipboard(text, type) {
  if (!text) return;
  if (window.Launcher) {
    window.Launcher.setClipboardText(text);
  } else {
    navigator.clipboard.writeText(text);
  }
  emit('copy-' + type, text);
}

function openUrl(url) {
  if (!url) return;
  if (window.Launcher) {
    window.Launcher.openLink(url);
  } else {
    window.open(url, '_blank');
  }
}

function handleKeyDown(event, action) {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    action();
  }
}

function getEntryIcon(entry) {
  if (entry.icon && entry.icon !== '00000000000000000000000000000000') {
    return 'fa-key';
  }
  return 'fa-file';
}

function formatDate(timestamp) {
  if (!timestamp) return '';
  return new Date(timestamp * 1000).toLocaleDateString();
}
</script>

<template>
  <div class="entry-table">
    <!-- Search Bar -->
    <div class="entry-table__search">
      <i class="fa fa-search entry-table__search-icon"></i>
      <input
        v-model="searchTerm"
        type="text"
        class="entry-table__search-input"
        placeholder="Search entries..."
      />
      <span v-if="filteredEntries.length" class="entry-table__count">
        {{ filteredEntries.length }} entries
      </span>
      <div class="entry-table__actions-bar">
        <button
          v-if="selectedEntry"
          class="entry-table__search-action"
          title="Copy Username"
          @click="copyToClipboard(selectedEntry.userName, 'username')"
        >
          <i class="fa fa-user"></i>
          <span>Username</span>
        </button>
        <button
          v-if="selectedEntry"
          class="entry-table__search-action"
          title="Copy Password"
          @click="copyToClipboard(selectedEntry.password, 'password')"
        >
          <i class="fa fa-key"></i>
          <span>Password</span>
        </button>
        <button
          v-if="selectedEntry?.url"
          class="entry-table__search-action"
          title="Open URL"
          @click="openUrl(selectedEntry.url)"
        >
          <i class="fa fa-external-link"></i>
          <span>Open</span>
        </button>
        <div class="entry-table__divider"></div>
        <button
          class="entry-table__search-action"
          title="Version"
          @click="$emit('show-version')"
        >
          <i class="fa fa-info-circle"></i>
        </button>
        <button
          class="entry-table__search-action"
          title="Help"
          @click="$emit('show-help')"
        >
          <i class="fa fa-question"></i>
        </button>
        <button
          class="entry-table__search-action"
          title="Settings"
          @click="$emit('show-settings')"
        >
          <i class="fa fa-cog"></i>
        </button>
        <button
          class="entry-table__search-action"
          title="Generate Password"
          @click="$emit('generate-password')"
        >
          <i class="fa fa-bolt"></i>
        </button>
        <button
          class="entry-table__search-action entry-table__search-action--danger"
          title="Lock Database"
          @click="$emit('logout')"
        >
          <i class="fa fa-sign-out-alt"></i>
        </button>
      </div>
    </div>

    <!-- Table Header -->
    <div class="entry-table__header">
      <div class="entry-table__col entry-table__col--icon"></div>
      <div class="entry-table__col entry-table__col--title">Title</div>
      <div class="entry-table__col entry-table__col--username">Username</div>
      <div class="entry-table__col entry-table__col--url">URL</div>
      <div class="entry-table__col entry-table__col--modified">Modified</div>
      <div class="entry-table__col entry-table__col--actions">Actions</div>
    </div>

    <!-- Table Body -->
    <div class="entry-table__body">
      <div
        v-for="entry in filteredEntries"
        :key="entry.uuid"
        :class="['entry-table__row', { 'entry-table__row--selected': selectedEntry?.uuid === entry.uuid }]"
        @click="selectEntry(entry)"
        @keydown="handleKeyDown($event, () => selectEntry(entry))"
        role="button"
        tabindex="0"
      >
        <div class="entry-table__col entry-table__col--icon">
          <i :class="`fa ${getEntryIcon(entry)}`"></i>
        </div>
        <div class="entry-table__col entry-table__col--title">
          <span class="entry-table__cell-text">{{ entry.title || '(no title)' }}</span>
        </div>
        <div class="entry-table__col entry-table__col--username">
          <span class="entry-table__cell-text">{{ entry.userName || '' }}</span>
        </div>
        <div class="entry-table__col entry-table__col--url">
          <span v-if="entry.url" class="entry-table__url" @click.stop="openUrl(entry.url)">
            <i class="fa fa-external-link"></i>
            {{ entry.url }}
          </span>
        </div>
        <div class="entry-table__col entry-table__col--modified">
          <span class="entry-table__cell-text">{{ formatDate(entry.modifiedTime) }}</span>
        </div>
        <div class="entry-table__col entry-table__col--actions">
          <button
            class="entry-table__action-btn"
            title="Copy Username"
            @click.stop="copyToClipboard(entry.userName, 'username')"
            @keydown="handleKeyDown($event, () => copyToClipboard(entry.userName, 'username'))"
          >
            <i class="fa fa-user"></i>
          </button>
          <button
            class="entry-table__action-btn"
            title="Copy Password"
            @click.stop="copyToClipboard(entry.password, 'password')"
            @keydown="handleKeyDown($event, () => copyToClipboard(entry.password, 'password'))"
          >
            <i class="fa fa-key"></i>
          </button>
        </div>
      </div>

      <div v-if="filteredEntries.length === 0" class="entry-table__empty">
        <i class="fa fa-search"></i>
        <p>No entries found</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.entry-table {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--background-color);
}

/* Search Bar */
.entry-table__search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--base-border-color);
  background-color: var(--secondary-background-color);
}

.entry-table__search-icon {
  color: var(--muted-color);
  font-size: 14px;
}

.entry-table__search-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--base-border-color);
  border-radius: var(--input-border-radius);
  background-color: var(--background-color);
  color: var(--text-color);
  font-size: 13px;
}

.entry-table__search-input:focus {
  outline: none;
  border-color: var(--action-color);
}

.entry-table__count {
  font-size: 12px;
  color: var(--muted-color);
  white-space: nowrap;
}

/* Actions Bar */
.entry-table__actions-bar {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
}

.entry-table__divider {
  width: 1px;
  height: 24px;
  background-color: var(--base-border-color);
  margin: 0 8px;
}

.entry-table__search-action {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: none;
  border: 1px solid transparent;
  color: var(--muted-color);
  cursor: pointer;
  border-radius: var(--button-border-radius);
  font-size: 12px;
  transition: all var(--fast-duration) var(--base-timing);
}

.entry-table__search-action:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
  border-color: var(--base-border-color);
}

.entry-table__search-action i {
  font-size: 14px;
}

.entry-table__search-action--danger {
  color: var(--error-color);
}

.entry-table__search-action--danger:hover {
  background-color: var(--error-background-color-focus-tr);
  border-color: var(--error-color);
}

/* Table Header */
.entry-table__header {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background-color: var(--intermediate-background-color);
  border-bottom: 1px solid var(--base-border-color);
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--muted-color);
  letter-spacing: 0.5px;
}

/* Table Body */
.entry-table__body {
  flex: 1;
  overflow-y: auto;
}

/* Table Row */
.entry-table__row {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid var(--light-border-color);
  cursor: pointer;
  transition: background-color var(--fast-duration) var(--base-timing);
}

.entry-table__row:hover {
  background-color: var(--hover-background-color);
}

.entry-table__row--selected {
  background-color: var(--selected-item-color) !important;
}

/* Table Columns */
.entry-table__col {
  padding: 4px 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.entry-table__col--icon {
  width: 40px;
  text-align: center;
  color: var(--action-color);
}

.entry-table__col--title {
  flex: 2;
  min-width: 150px;
}

.entry-table__col--username {
  flex: 1;
  min-width: 120px;
  color: var(--muted-color);
}

.entry-table__col--url {
  flex: 1;
  min-width: 150px;
}

.entry-table__url {
  color: var(--action-color);
  cursor: pointer;
  font-size: 12px;
}

.entry-table__url:hover {
  text-decoration: underline;
}

.entry-table__col--modified {
  width: 100px;
  font-size: 11px;
  color: var(--muted-color);
}

.entry-table__col--actions {
  width: 80px;
  display: flex;
  gap: 4px;
  justify-content: flex-end;
}

.entry-table__cell-text {
  font-size: 13px;
}

/* Action Buttons */
.entry-table__action-btn {
  background: none;
  border: 1px solid transparent;
  color: var(--muted-color);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--button-border-radius);
  font-size: 12px;
  transition: all var(--fast-duration) var(--base-timing);
}

.entry-table__action-btn:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
  border-color: var(--base-border-color);
}

/* Empty State */
.entry-table__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--muted-color);
  text-align: center;
}

.entry-table__empty i {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.entry-table__empty p {
  margin: 0;
  font-size: 13px;
}

/* Scrollbar */
.entry-table__body::-webkit-scrollbar {
  width: 8px;
}

.entry-table__body::-webkit-scrollbar-track {
  background: transparent;
}

.entry-table__body::-webkit-scrollbar-thumb {
  background-color: var(--muted-color);
  border-radius: 4px;
}

.entry-table__body::-webkit-scrollbar-thumb:hover {
  background-color: var(--text-semi-muted-color);
}
</style>

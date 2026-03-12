<script setup>
import { ref } from 'vue';

const props = defineProps({
  entry: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['close', 'copy-username', 'copy-password', 'copy-url', 'copy-notes']);

const showPassword = ref(false);

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
  let fullUrl = url;
  if (!/^(https?:\/\/|ftp:\/\/|mailto:|sftp:\/\/)/i.test(url)) {
    fullUrl = 'https://' + url;
  }
  if (window.Launcher) {
    window.Launcher.openLink(fullUrl);
  } else {
    window.open(fullUrl, '_blank');
  }
}

function handleKeyDown(event, action) {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    action();
  }
}

function formatDate(timestamp) {
  if (!timestamp) return '';
  return new Date(timestamp * 1000).toLocaleDateString();
}

function getEntryIcon(entry) {
  if (entry.icon && entry.icon !== '00000000000000000000000000000000') {
    return 'fa-key bg-action-color';
  }
  return 'fa-file';
}
</script>

<template>
  <div class="entry-detail">
    <!-- Header -->
    <div class="entry-detail__header">
      <button
        class="entry-detail__back"
        @click="$emit('close')"
        @keydown="handleKeyDown($event, () => $emit('close'))"
        title="Back"
      >
        <i class="fa fa-arrow-left"></i>
      </button>
      <h1 class="entry-detail__title">{{ entry.title || '(no title)' }}</h1>
    </div>

    <!-- Body -->
    <div class="entry-detail__body">
      <div class="entry-detail__fields">
        <!-- Username -->
        <div class="entry-detail__field">
          <div class="entry-detail__field-label">
            <i class="fa fa-user"></i>
            <span>Username</span>
          </div>
          <div class="entry-detail__field-value">
            <input
              type="text"
              :value="entry.userName || ''"
              readonly
              class="entry-detail__input"
            />
            <button
              class="entry-detail__copy-btn"
              title="Copy"
              @click="copyToClipboard(entry.userName, 'username')"
              @keydown="handleKeyDown($event, () => copyToClipboard(entry.userName, 'username'))"
            >
              <i class="fa fa-copy"></i>
            </button>
          </div>
        </div>

        <!-- Password -->
        <div class="entry-detail__field">
          <div class="entry-detail__field-label">
            <i class="fa fa-key"></i>
            <span>Password</span>
          </div>
          <div class="entry-detail__field-value">
            <input
              :type="showPassword ? 'text' : 'password'"
              :value="entry.password || ''"
              readonly
              class="entry-detail__input"
            />
            <button
              class="entry-detail__copy-btn"
              title="Copy"
              @click="copyToClipboard(entry.password, 'password')"
              @keydown="handleKeyDown($event, () => copyToClipboard(entry.password, 'password'))"
            >
              <i class="fa fa-copy"></i>
            </button>
            <button
              class="entry-detail__copy-btn"
              :title="showPassword ? 'Hide' : 'Show'"
              @click="showPassword = !showPassword"
              @keydown="handleKeyDown($event, () => showPassword = !showPassword)"
            >
              <i :class="showPassword ? 'fa fa-eye-slash' : 'fa fa-eye'"></i>
            </button>
          </div>
        </div>

        <!-- URL -->
        <div v-if="entry.url" class="entry-detail__field">
          <div class="entry-detail__field-label">
            <i class="fa fa-link"></i>
            <span>URL</span>
          </div>
          <div class="entry-detail__field-value">
            <input
              type="text"
              :value="entry.url"
              readonly
              class="entry-detail__input"
            />
            <button
              class="entry-detail__copy-btn"
              title="Copy"
              @click="copyToClipboard(entry.url, 'url')"
              @keydown="handleKeyDown($event, () => copyToClipboard(entry.url, 'url'))"
            >
              <i class="fa fa-copy"></i>
            </button>
            <button
              class="entry-detail__copy-btn"
              title="Open"
              @click="openUrl(entry.url)"
              @keydown="handleKeyDown($event, () => openUrl(entry.url))"
            >
              <i class="fa fa-globe"></i>
            </button>
          </div>
        </div>

        <!-- Notes -->
        <div v-if="entry.notes" class="entry-detail__field">
          <div class="entry-detail__field-label">
            <i class="fa fa-sticky-note"></i>
            <span>Notes</span>
          </div>
          <div class="entry-detail__field-value">
            <textarea
              :value="entry.notes"
              readonly
              class="entry-detail__textarea"
            ></textarea>
          </div>
        </div>

        <!-- Custom Fields -->
        <div
          v-for="(value, key) in entry.customFields"
          :key="key"
          class="entry-detail__field"
        >
          <div class="entry-detail__field-label">
            <i class="fa fa-tag"></i>
            <span>{{ key }}</span>
          </div>
          <div class="entry-detail__field-value">
            <input
              type="text"
              :value="value"
              readonly
              class="entry-detail__input"
            />
            <button
              class="entry-detail__copy-btn"
              title="Copy"
              @click="copyToClipboard(value, 'custom')"
              @keydown="handleKeyDown($event, () => copyToClipboard(value, 'custom'))"
            >
              <i class="fa fa-copy"></i>
            </button>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="entry-detail__sidebar">
        <div class="entry-detail__section">
          <h3 class="entry-detail__section-title">
            <i class="fa fa-clock"></i>
            <span>History</span>
          </h3>
          <div class="entry-detail__info">
            <div class="entry-detail__info-row">
              <span class="entry-detail__info-label">Created:</span>
              <span class="entry-detail__info-value">{{ formatDate(entry.createdTime) }}</span>
            </div>
            <div class="entry-detail__info-row">
              <span class="entry-detail__info-label">Modified:</span>
              <span class="entry-detail__info-value">{{ formatDate(entry.modifiedTime) }}</span>
            </div>
            <div v-if="entry.expiryTime" class="entry-detail__info-row">
              <span class="entry-detail__info-label">Expires:</span>
              <span class="entry-detail__info-value">{{ formatDate(entry.expiryTime) }}</span>
            </div>
          </div>
        </div>

        <div v-if="entry.tags?.length" class="entry-detail__section">
          <h3 class="entry-detail__section-title">
            <i class="fa fa-tags"></i>
            <span>Tags</span>
          </h3>
          <div class="entry-detail__tags">
            <span
              v-for="tag in entry.tags"
              :key="tag"
              class="entry-detail__tag"
            >
              {{ tag }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.entry-detail {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--background-color);
}

/* Header */
.entry-detail__header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  border-bottom: 1px solid var(--base-border-color);
  background-color: var(--secondary-background-color);
}

.entry-detail__back {
  background: none;
  border: 1px solid var(--base-border-color);
  color: var(--muted-color);
  cursor: pointer;
  padding: 8px 12px;
  border-radius: var(--button-border-radius);
  font-size: 14px;
  transition: all var(--fast-duration) var(--base-timing);
}

.entry-detail__back:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
}

.entry-detail__icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, var(--action-color), var(--selected-item-color));
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(49, 126, 246, 0.3);
}

.entry-detail__icon i {
  font-size: 24px;
  color: white;
}

.entry-detail__title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-color);
  flex: 1;
}

/* Body */
.entry-detail__body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.entry-detail__fields {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.entry-detail__field {
  margin-bottom: 24px;
}

.entry-detail__field-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--muted-color);
  margin-bottom: 8px;
  letter-spacing: 0.5px;
}

.entry-detail__field-label i {
  font-size: 12px;
}

.entry-detail__field-value {
  display: flex;
  gap: 8px;
  align-items: center;
}

.entry-detail__input,
.entry-detail__textarea {
  flex: 1;
  padding: 12px;
  border: 1px solid var(--base-border-color);
  border-radius: var(--input-border-radius);
  background-color: var(--secondary-background-color);
  color: var(--text-color);
  font-family: inherit;
  font-size: 13px;
}

.entry-detail__input:focus,
.entry-detail__textarea:focus {
  outline: none;
  border-color: var(--action-color);
}

.entry-detail__textarea {
  resize: vertical;
  min-height: 250px;
}

.entry-detail__copy-btn {
  background: none;
  border: 1px solid var(--base-border-color);
  color: var(--muted-color);
  cursor: pointer;
  padding: 10px 12px;
  border-radius: var(--button-border-radius);
  font-size: 14px;
  transition: all var(--fast-duration) var(--base-timing);
}

.entry-detail__copy-btn:hover {
  background-color: var(--hover-background-color);
  color: var(--text-color);
  border-color: var(--action-color);
}

/* Sidebar */
.entry-detail__sidebar {
  width: 280px;
  border-left: 1px solid var(--base-border-color);
  padding: 24px;
  background-color: var(--intermediate-background-color);
  overflow-y: auto;
}

.entry-detail__section {
  margin-bottom: 32px;
}

.entry-detail__section:last-child {
  margin-bottom: 0;
}

.entry-detail__section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--muted-color);
  margin: 0 0 16px 0;
  letter-spacing: 0.5px;
}

.entry-detail__section-title i {
  font-size: 14px;
}

.entry-detail__info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.entry-detail__info-row {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.entry-detail__info-label {
  color: var(--muted-color);
}

.entry-detail__info-value {
  color: var(--text-color);
  font-weight: 500;
}

.entry-detail__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.entry-detail__tag {
  padding: 4px 10px;
  background-color: var(--secondary-background-color);
  border: 1px solid var(--base-border-color);
  border-radius: 12px;
  font-size: 11px;
  color: var(--text-color);
}

/* Scrollbar */
.entry-detail__fields::-webkit-scrollbar,
.entry-detail__sidebar::-webkit-scrollbar {
  width: 8px;
}

.entry-detail__fields::-webkit-scrollbar-track,
.entry-detail__sidebar::-webkit-scrollbar-track {
  background: transparent;
}

.entry-detail__fields::-webkit-scrollbar-thumb,
.entry-detail__sidebar::-webkit-scrollbar-thumb {
  background-color: var(--muted-color);
  border-radius: 4px;
}
</style>

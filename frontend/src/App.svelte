<script>
  import { onMount } from 'svelte';
  
  // State
  let selectedEntry = null;
  let isDarkTheme = true;
  let showPassword = false;
  let entries = [
    { id: 1, title: 'Google Account', username: 'user@gmail.com', icon: 'globe', color: 'blue', url: 'https://google.com' },
    { id: 2, title: 'GitHub', username: 'developer', icon: 'github', color: 'violet', url: 'https://github.com' },
    { id: 3, title: 'AWS Console', username: 'admin', icon: 'cloud', color: 'orange', url: 'https://aws.amazon.com' },
    { id: 4, title: 'Dropbox', username: 'user@example.com', icon: 'dropbox', color: 'blue', url: 'https://dropbox.com' },
    { id: 5, title: 'Local Database', username: 'admin', icon: 'database', color: 'green', url: '' },
    { id: 6, title: 'SSH Server', username: 'root', icon: 'server', color: 'red', url: '' },
    { id: 7, title: 'WiFi Password', username: '', icon: 'wifi', color: 'green', url: '' },
    { id: 8, title: 'Bank Account', username: '123456789', icon: 'university', color: 'yellow', url: 'https://bank.com' },
  ];
  
  let groups = [
    { id: 1, name: 'Internet', icon: 'globe', count: 4 },
    { id: 2, name: 'Work', icon: 'briefcase', count: 2 },
    { id: 3, name: 'Personal', icon: 'user', count: 2 },
    { id: 4, name: 'Archived', icon: 'archive', count: 0 },
  ];
  
  let openFiles = [
    { id: 1, name: 'Personal.kdbx', active: true, modified: false },
    { id: 2, name: 'Work.kdbx', active: false, modified: true },
  ];
  
  // Update body theme class
  function updateBodyTheme() {
    if (isDarkTheme) {
      document.body.classList.add('th-dark');
      document.body.classList.remove('th-light');
    } else {
      document.body.classList.add('th-light');
      document.body.classList.remove('th-dark');
    }
  }
  
  // Initialize with first entry selected and set theme
  onMount(() => {
    if (entries.length > 0 && !selectedEntry) {
      selectedEntry = entries[0];
    }
    updateBodyTheme();
  });
  
  function selectEntry(entry) {
    selectedEntry = entry;
  }
  
  function toggleTheme() {
    isDarkTheme = !isDarkTheme;
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
    showPassword = !showPassword;
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
</script>

<div class="app">
  <!-- Titlebar drag area (for frameless windows) -->
  <div class="app__titlebar-drag"></div>
  
  <div class="app__body">
    <!-- Left sidebar menu -->
    <div class="app__menu">
      <div class="menu">
        <div class="menu__header">
          <div class="menu__header-title">KeeWeb</div>
          <button class="menu__header-button" on:click={toggleTheme} title="Toggle theme">
            <i class="fa fa-{isDarkTheme ? 'sun' : 'moon'}"></i>
          </button>
        </div>
        
        <div class="menu__sections">
          {#each groups as group}
            <div class="menu__section">
              <button
                class="menu__section-header"
                on:click={() => console.log('Group clicked:', group.name)}
                on:keydown={(e) => handleKeyDown(e, () => console.log('Group clicked:', group.name))}
              >
                  <i class="fa fa-{group.icon} menu__section-icon"></i>
                  <span class="menu__section-name">{group.name}</span>
                  {#if group.count > 0}
                    <span class="menu__section-count">{group.count}</span>
                  {/if}
                </button>
            </div>
          {/each}
        </div>
        
        <div class="menu__footer">
          <button
            class="menu__footer-item"
            on:click={() => console.log('New Group clicked')}
            on:keydown={(e) => handleKeyDown(e, () => console.log('New Group clicked'))}
          >
            <i class="fa fa-plus"></i>
            <span>New Group</span>
          </button>
          <button
            class="menu__footer-item"
            on:click={() => console.log('Search clicked')}
            on:keydown={(e) => handleKeyDown(e, () => console.log('Search clicked'))}
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
                on:click={() => console.log('New Entry clicked')}
                on:keydown={(e) => handleKeyDown(e, () => console.log('New Entry clicked'))}
              >
                <i class="fa fa-plus"></i>
              </button>
              <button 
                class="list__control-button" 
                title="Sort"
                on:click={() => console.log('Sort clicked')}
                on:keydown={(e) => handleKeyDown(e, () => console.log('Sort clicked'))}
              >
                <i class="fa fa-sort"></i>
              </button>
            </div>
          </div>
          
          <div class="list__items">
            {#each entries as entry}
              <div
                class="list__item {selectedEntry?.id === entry.id ? 'list__item--active' : ''}"
                on:click={() => selectEntry(entry)}
                on:keydown={(e) => handleKeyDown(e, () => selectEntry(entry))}
                role="button"
                tabindex="0"
              >
                <div class="list__item-icon">
                  <i class="{getEntryIconClass(entry)}"></i>
                </div>
                <div class="list__item-content">
                  <div class="list__item-title">{entry.title}</div>
                  {#if entry.username}
                    <div class="list__item-subtitle">{entry.username}</div>
                  {/if}
                </div>
                {#if entry.url}
                  <div 
                    class="list__item-action"
                    on:click={() => selectedEntry.url && openUrl(selectedEntry.url)}
                    on:keydown={(e) => handleKeyDown(e, () => selectedEntry.url && openUrl(selectedEntry.url))}
                    title="Open URL"
                  >
                    <i class="fa fa-external-link"></i>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      </div>
      
      <!-- Resize handle -->
      <div class="app__list-drag"></div>
      
      <!-- Entry details -->
      <div class="app__details">
        {#if selectedEntry}
          <div class="details">
            <div class="details__header">
              <i class="details__header-color {getEntryIconClass(selectedEntry)}"></i>
              <h1 class="details__header-title">{selectedEntry.title}</h1>
            </div>
            
            <div class="details__body">
              <div class="scroller">
                <div class="details__body-fields">
                  <div class="details__field">
                      <div class="details__field-label">Username</div>
                      <div class="details__field-value">
                        <input type="text" value={selectedEntry.username || ''} readonly />
                        <button 
                          class="details__field-copy" 
                          title="Copy"
                          on:click={() => copyToClipboard(selectedEntry.username)}
                          on:keydown={(e) => handleKeyDown(e, () => copyToClipboard(selectedEntry.username))}
                        >
                          <i class="fa fa-copy"></i>
                        </button>
                      </div>
                    </div>
                  
                    <div class="details__field">
                      <div class="details__field-label">Password</div>
                      <div class="details__field-value">
                        <input 
                          type={showPassword ? "text" : "password"} 
                          value="password123" 
                          readonly 
                        />
                        <button 
                          class="details__field-copy" 
                          title="Copy"
                          on:click={() => copyToClipboard('password123')}
                          on:keydown={(e) => handleKeyDown(e, () => copyToClipboard('password123'))}
                        >
                          <i class="fa fa-copy"></i>
                        </button>
                        <button 
                          class="details__field-show" 
                          title={showPassword ? "Hide" : "Show"}
                          on:click={togglePasswordVisibility}
                          on:keydown={(e) => handleKeyDown(e, togglePasswordVisibility)}
                        >
                          <i class="fa fa-eye{showPassword ? '-slash' : ''}"></i>
                        </button>
                      </div>
                    </div>
                  
                    {#if selectedEntry.url}
                      <div class="details__field">
                        <div class="details__field-label">URL</div>
                        <div class="details__field-value">
                          <input type="text" value={selectedEntry.url} readonly />
                          <button 
                            class="details__field-copy" 
                            title="Copy"
                            on:click={() => copyToClipboard(selectedEntry.url)}
                            on:keydown={(e) => handleKeyDown(e, () => copyToClipboard(selectedEntry.url))}
                          >
                            <i class="fa fa-copy"></i>
                          </button>
                          <button 
                            class="details__field-open" 
                            title="Open"
                            on:click={() => openUrl(selectedEntry.url)}
                            on:keydown={(e) => handleKeyDown(e, () => openUrl(selectedEntry.url))}
                          >
                            <i class="fa fa-external-link"></i>
                          </button>
                        </div>
                      </div>
                    {/if}
                  
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
                      on:click={() => console.log('Auto-type would be performed here')}
                      on:keydown={(e) => handleKeyDown(e, () => console.log('Auto-type would be performed here'))}
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
                on:click={() => console.log('Edit entry')}
                on:keydown={(e) => handleKeyDown(e, () => console.log('Edit entry'))}
              >
                <i class="fa fa-edit"></i>
              </button>
              <button 
                class="details__button details__button--delete" 
                title="Delete"
                on:click={() => console.log('Delete entry')}
                on:keydown={(e) => handleKeyDown(e, () => console.log('Delete entry'))}
              >
                <i class="fa fa-trash-alt"></i>
              </button>
            </div>
          </div>
        {:else}
          <div class="details details--empty">
            <div class="details__empty-message">
              <i class="fa fa-key details__empty-icon"></i>
              <h2>No Entry Selected</h2>
              <p>Select an entry from the list to view its details</p>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
  
  <!-- Footer -->
  <div class="app__footer">
    <div class="footer">
      {#each openFiles as file}
        <div class="footer__db {file.active ? '' : 'footer__db--dimmed'}" title="{file.name}">
          <i class="fa fa-{file.active ? 'unlock' : 'lock'}"></i>
          <span class="footer__db-name">{file.name}</span>
          {#if file.modified}
            <i class="fa fa-circle footer__db-sign"></i>
          {/if}
        </div>
      {/each}
      
      <div class="footer__db footer__db--dimmed footer__db-open">
        <i class="fa fa-plus"></i>
        <span class="footer__db-text">Open Database</span>
      </div>
      
      <button
        class="footer__btn"
        title="Help"
        on:click={() => console.log('Help clicked')}
        on:keydown={(e) => handleKeyDown(e, () => console.log('Help clicked'))}
      >
        <i class="fa fa-question"></i>
      </button>

      <button
        class="footer__btn"
        title="Settings"
        on:click={() => console.log('Settings clicked')}
        on:keydown={(e) => handleKeyDown(e, () => console.log('Settings clicked'))}
      >
        <i class="fa fa-cog"></i>
      </button>

      <button
        class="footer__btn"
        title="Generate Password"
        on:click={() => console.log('Generate password clicked')}
        on:keydown={(e) => handleKeyDown(e, () => console.log('Generate password clicked'))}
      >
        <i class="fa fa-bolt"></i>
      </button>

      <button
        class="footer__btn"
        title="Lock Database"
        on:click={() => console.log('Lock database clicked')}
        on:keydown={(e) => handleKeyDown(e, () => console.log('Lock database clicked'))}
      >
        <i class="fa fa-sign-out-alt"></i>
      </button>
    </div>
  </div>
</div>

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
    margin-bottom: 2px;
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
</style>
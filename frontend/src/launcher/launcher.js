// Launcher for Wails v2
// Minimal implementation without KeeWeb dependencies

// Simple logger using console
const logger = {
    debug: (...args) => console.debug('[Launcher]', ...args),
    info: (...args) => console.info('[Launcher]', ...args),
    warn: (...args) => console.warn('[Launcher]', ...args),
    error: (...args) => console.error('[Launcher]', ...args),
    ts: () => `t=${Date.now()}`
};

// Helper to call Go methods
function callGo(methodName, ...args) {
    const goApp = window.go?.main?.App;
    if (!goApp || typeof goApp[methodName] !== 'function') {
        logger.warn(`Go method ${methodName} not available`);
        return Promise.reject(new Error(`Go method ${methodName} not available`));
    }
    try {
        return goApp[methodName](...args);
    } catch (err) {
        logger.error(`Error calling ${methodName}:`, err);
        return Promise.reject(err);
    }
}

// Simple event emitter
const EventEmitter = {
    _listeners: {},

    on(event, callback) {
        if (!this._listeners[event]) this._listeners[event] = [];
        this._listeners[event].push(callback);
    },

    emit(event, data) {
        const listeners = this._listeners[event];
        if (listeners) {
            listeners.forEach(cb => cb(data));
        }
        // Also dispatch global event
        window.dispatchEvent(new CustomEvent(event, { detail: data }));
    },

    off(event, callback) {
        const listeners = this._listeners[event];
        if (listeners) {
            const index = listeners.indexOf(callback);
            if (index > -1) listeners.splice(index, 1);
        }
    }
};

const Launcher = {
    name: 'keebell',
    version: '2.11.0',
    autoTypeSupported: false,
    thirdPartyStoragesSupported: true,
    clipboardSupported: true,
    devTools: false,

    // Instance state
    exitRequested: false,
    pendingUpdateFile: null,
    pendingFileToOpen: null,
    readyToOpenFiles: false,

    // Platform methods
    platform() {
        return callGo('Platform');
    },

    arch() {
        return callGo('Arch');
    },

    electron() {
        return null; // Not available in Wails
    },

    remoteApp() {
        // Return a proxy object that forwards calls to Go
        const self = this;
        return {
            getPath: (name) => {
                switch (name) {
                    case 'userData': return self.getUserDataPath();
                    case 'temp': return self.getTempPath();
                    case 'documents': return self.getDocumentsPath();
                    case 'app': return self.getAppPath();
                    default: return '';
                }
            },
            loadConfig: (name) => self.loadConfig(name),
            saveConfig: (name, data) => self.saveConfig(name, data),
            getMainWindow: () => ({ webContents: { session: {} } }),
            setHookBeforeQuitEvent: () => {},
            quit: () => self.exit(),
            restartAndUpdate: (file) => self.requestRestartAndUpdate(file),
            minimizeApp: (options) => self.minimizeApp(options?.restore, options?.quit),
            minimizeThenHideIfInTray: () => {
                self.minimizeMainWindow();
                setTimeout(() => self.hideApp(), 100);
            },
            hide: () => self.hideApp(),
            showAndFocusMainWindow: () => self.showMainWindow(),
            setGlobalShortcuts: (settings) => self.setGlobalShortcuts(settings),
            setAboutPanelOptions: () => {},
            on: (event, handler) => {
                logger.debug(`Event registration: ${event}`);
                // Could map to EventEmitter if needed
            }
        };
    },

    remReq(mod) {
        throw new Error(`Native module ${mod} not available in Wails`);
    },

    // Basic I/O
    openLink(href) {
        if (/^(http|https|ftp|sftp|mailto):/i.test(href)) {
            return callGo('OpenLink', href);
        }
    },

    openDevTools() {
        logger.debug('DevTools not available in Wails v2');
    },

    // File dialogs
    getSaveFileName(defaultPath, callback) {
        const title = 'Save File';
        const filterName = 'KeePass Database';
        const filterExt = '*.kdbx';

        callGo('GetSaveFileName', defaultPath || '', title, filterName, filterExt)
            .then(result => callback(result || ''))
            .catch(err => {
                logger.error('Error getting save file name', err);
                callback('');
            });
    },

    // Path utilities
    getUserDataPath(fileName) {
        return callGo('GetUserDataPath', fileName || '');
    },

    getTempPath(fileName) {
        return callGo('GetTempPath', fileName || '');
    },

    getDocumentsPath(fileName) {
        return callGo('GetDocumentsPath', fileName || '');
    },

    getAppPath(fileName) {
        return callGo('GetAppPath', fileName || '');
    },

    getWorkDirPath(fileName) {
        return callGo('GetWorkDirPath', fileName || '');
    },

    joinPath(...parts) {
        return callGo('JoinPath', ...parts);
    },

    // File operations (callback-based)
    writeFile(path, data, callback) {
        const buffer = data instanceof Uint8Array ? Array.from(data) : data;
        callGo('WriteFile', path, buffer)
            .then(() => callback && callback())
            .catch(err => callback && callback(err));
    },

    readFile(path, encoding, callback) {
        callGo('ReadFile', path)
            .then(data => {
                if (data instanceof Array) {
                    data = new Uint8Array(data);
                }
                callback(data, null);
            })
            .catch(err => {
                callback(null, err);
            });
    },

    fileExists(path, callback) {
        callGo('FileExists', path)
            .then(exists => callback(exists))
            .catch(() => callback(false));
    },

    fileExistsSync(path) {
        // Sync version not available
        return false;
    },

    deleteFile(path, callback) {
        callGo('DeleteFile', path)
            .then(() => callback && callback())
            .catch(err => callback && callback(err));
    },

    statFile(path, callback) {
        callGo('StatFile', path)
            .then(stat => {
                if (stat && typeof stat === 'object') {
                    // Convert to expected format
                    callback({
                        uid: stat.uid || 0,
                        gid: stat.gid || 0,
                        size: stat.size || 0,
                        mtime: stat.mtime ? new Date(stat.mtime) : new Date()
                    }, null);
                } else {
                    callback(null, 'Invalid stat object');
                }
            })
            .catch(err => {
                callback(null, err);
            });
    },

    mkdir(dir, callback) {
        callGo('Mkdir', dir)
            .then(() => callback && callback())
            .catch(err => callback && callback(err));
    },

    parsePath(fileName) {
        const result = callGo('ParsePath', fileName);
        if (result && typeof result === 'object') {
            return {
                path: result.path || fileName,
                dir: result.dir || '',
                file: result.file || ''
            };
        }
        // Fallback
        const lastSlash = fileName.lastIndexOf('/');
        const lastBackslash = fileName.lastIndexOf('\\');
        const lastSep = Math.max(lastSlash, lastBackslash);
        return {
            path: fileName,
            dir: lastSep > 0 ? fileName.substring(0, lastSep) : '',
            file: lastSep > 0 ? fileName.substring(lastSep + 1) : fileName
        };
    },

    createFsWatcher(path) {
        logger.warn('FS watcher not implemented for Wails');
        return { close: () => {} };
    },

    // Config
    loadConfig(name) {
        return callGo('LoadConfig', name);
    },

    saveConfig(name, data) {
        return callGo('SaveConfig', name, data);
    },

    // Application lifecycle
    preventExit(e) {
        e.returnValue = false;
        return false;
    },

    exit() {
        this.exitRequested = true;
        this.requestExit();
    },

    requestExit() {
        callGo('RequestExit');
    },

    requestRestartAndUpdate(updateFilePath) {
        this.pendingUpdateFile = updateFilePath;
        callGo('RequestRestartAndUpdate', updateFilePath);
    },

    cancelRestart() {
        this.pendingUpdateFile = null;
        callGo('CancelRestart');
    },

    // Clipboard
    setClipboardText(text) {
        return callGo('SetClipboardText', text);
    },

    getClipboardText() {
        return callGo('GetClipboardText');
    },

    clearClipboardText() {
        // Wails fallback - set empty text
        callGo('SetClipboardText', '');
    },

    quitOnRealQuitEventIfMinimizeOnQuitIsEnabled() {
        return !!this.pendingUpdateFile;
    },

    // Window management
    minimizeApp(options) {
        if (typeof options === 'string') {
            // Backwards compatibility
            return callGo('MinimizeApp', options, '');
        } else if (options && typeof options === 'object') {
            return callGo('MinimizeApp', options.restore || '', options.quit || '');
        } else {
            return callGo('MinimizeMainWindow');
        }
    },

    canDetectOsSleep() {
        return callGo('CanDetectOsSleep');
    },

    updaterEnabled() {
        return callGo('UpdaterEnabled');
    },

    getMainWindow() {
        // Dummy window object
        return {
            webContents: {
                session: {
                    resolveProxy: () => Promise.resolve('')
                }
            },
            minimize: () => this.minimizeMainWindow(),
            maximize: () => this.maximizeMainWindow(),
            restore: () => this.restoreMainWindow(),
            isMaximized: () => this.mainWindowMaximized()
        };
    },

    resolveProxy(url, callback) {
        // TODO: Implement proxy resolution
        callback(null);
    },

    hideApp() {
        return callGo('HideApp');
    },

    isAppFocused() {
        return callGo('IsAppFocused');
    },

    showMainWindow() {
        return callGo('ShowMainWindow');
    },

    spawn(config) {
        const ts = logger.ts();
        let { complete } = config;
        delete config.complete;

        callGo('Spawn', config)
            .then(result => {
                if (result && typeof result === 'object') {
                    const code = result.code || 0;
                    const stdout = result.stdout || '';
                    const stderr = result.stderr || '';
                    const err = result.err;

                    const msg = 'spawn ' + config.cmd + ': ' + code + ', ' + logger.ts(ts);
                    if (code !== 0) {
                        logger.error(msg + '\n' + stdout + '\n' + stderr);
                    } else {
                        logger.info(msg + (stdout && !config.noStdOutLogging ? '\n' + stdout : ''));
                    }

                    if (complete) {
                        complete(err || (code !== 0 ? 'Exit code ' + code : null), stdout, code);
                    }
                } else {
                    if (complete) {
                        complete('Invalid spawn result');
                    }
                }
            })
            .catch(err => {
                if (complete) {
                    complete(err);
                }
            });
    },

    checkOpenFiles() {
        this.readyToOpenFiles = true;
        if (this.pendingFileToOpen) {
            this.openFile(this.pendingFileToOpen);
            delete this.pendingFileToOpen;
        }
    },

    openFile(file) {
        if (this.readyToOpenFiles) {
            EventEmitter.emit('launcher-open-file', file);
        } else {
            this.pendingFileToOpen = file;
        }
    },

    setGlobalShortcuts(appSettings) {
        return callGo('SetGlobalShortcuts', appSettings);
    },

    minimizeMainWindow() {
        return callGo('MinimizeMainWindow');
    },

    maximizeMainWindow() {
        return callGo('MaximizeMainWindow');
    },

    restoreMainWindow() {
        return callGo('RestoreMainWindow');
    },

    mainWindowMaximized() {
        return callGo('MainWindowMaximized');
    }
};

// Wire up global event handlers
EventEmitter.on('launcher-exit-request', () => {
    setTimeout(() => Launcher.exit(), 0);
});

EventEmitter.on('launcher-minimize', () => {
    setTimeout(() => EventEmitter.emit('app-minimized'), 0);
});

EventEmitter.on('launcher-maximize', () => {
    setTimeout(() => EventEmitter.emit('app-maximized'), 0);
});

EventEmitter.on('launcher-unmaximize', () => {
    setTimeout(() => EventEmitter.emit('app-unmaximized'), 0);
});

EventEmitter.on('launcher-started-minimized', () => {
    setTimeout(() => Launcher.minimizeApp(), 0);
});

// Global function for opening files
window.launcherOpen = (file) => Launcher.openFile(file);

// Check for pending file opens
if (window.launcherOpenedFile) {
    logger.info('Open file request', window.launcherOpenedFile);
    Launcher.openFile(window.launcherOpenedFile);
    delete window.launcherOpenedFile;
}

// Setup app-ready handler
EventEmitter.on('app-ready', () => {
    setTimeout(() => {
        Launcher.checkOpenFiles();
    }, 0);
});

// Export Launcher and EventEmitter for use
window.Launcher = Launcher;
export { Launcher, EventEmitter };

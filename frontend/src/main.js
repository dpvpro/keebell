import './style.css';
import './launcher/launcher.js'; // This will set window.Launcher
import {createApp} from 'vue';
import App from './app.vue';

createApp(App).mount('#app');

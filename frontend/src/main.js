import './style.css';
import './launcher/launcher.js'; // This will set window.Launcher
import App from './App.svelte';

const app = new App({
  target: document.getElementById('app')
});

export default app;

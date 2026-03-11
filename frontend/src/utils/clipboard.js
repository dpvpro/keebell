/**
 * Утилиты для работы с буфером обмена
 */
export function copyToClipboard(text) {
  if (!text) return false;

  // Используем Launcher если доступен
  if (window.Launcher && window.Launcher.setClipboardText) {
    window.Launcher.setClipboardText(text);
    return true;
  }

  // Fallback к navigator clipboard API
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).catch(err => {
      console.error('Failed to copy:', err);
    });
    return true;
  }

  // Старый метод через textarea
  try {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    return true;
  } catch (e) {
    console.error('Failed to copy:', e);
    return false;
  }
}

/**
 * Открытие URL в браузере
 */
export function openUrl(url) {
  if (!url) return;

  if (window.Launcher && window.Launcher.openLink) {
    window.Launcher.openLink(url);
  } else {
    window.open(url, '_blank');
  }
}

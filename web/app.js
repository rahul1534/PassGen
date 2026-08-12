function showLoadError(message) {
  const el = document.getElementById('validation-error');
  el.textContent = message;
  el.classList.remove('hidden');
}

async function loadApp() {
  if (!globalThis.crypto || !globalThis.crypto.getRandomValues) {
    showLoadError('Secure random source unavailable. Cannot generate passwords safely.');
    return;
  }

  const go = new Go();
  const response = await fetch('./app.wasm');
  if (!response.ok) {
    throw new Error('WASM fetch failed: ' + response.status);
  }

  const bytes = await response.arrayBuffer();
  const result = await WebAssembly.instantiate(bytes, go.importObject);
  await go.run(result.instance);
}

loadApp().catch((err) => {
  console.error('PassForge failed to start:', err);
  showLoadError('Failed to load application. Please refresh the page.');
});

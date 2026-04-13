import { mount } from 'svelte'
import App from './App.svelte';

// Disable default browser context menu globally
document.addEventListener('contextmenu', event => {
  event.preventDefault();
});

const app = mount(App, {
  target: document.getElementById('app')!,

})

export default app;

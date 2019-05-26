for (const key in window.on) {
  window.on[key] = e => console.log(typeof e);
}

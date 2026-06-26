// Inicialização do Google Analytics (gtag.js) externalizada — evita script inline
// e permite uma CSP sem 'unsafe-inline' em script-src.
window.dataLayer = window.dataLayer || []
function gtag() {
  dataLayer.push(arguments)
}
gtag("js", new Date())
gtag("config", "G-VQ1KTKXVWF")

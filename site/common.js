document.addEventListener('DOMContentLoaded', () => {
  const ver = typeof cheesseVersion === 'function' ? ' ' + cheesseVersion() : ''
  const footer = document.createElement('footer')
  footer.innerHTML = `
    cheesse${ver} &copy; ${new Date().getFullYear()} Mariano Gappa |
    <a href="https://github.com/marianogappa/cheesse">Project site</a> |
    <a href="https://github.com/marianogappa/cheesse/issues">Report an issue</a> |
    <a href="https://github.com/marianogappa/cheesse/blob/master/LICENSE">MIT License</a>
  `
  document.body.appendChild(footer)
})

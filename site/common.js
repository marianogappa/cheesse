document.addEventListener('DOMContentLoaded', () => {
  const footer = document.createElement('footer')
  footer.innerHTML = `
    cheesse &copy; ${new Date().getFullYear()} Mariano Gappa |
    <a href="https://github.com/marianogappa/cheesse">Project site</a> |
    <a href="https://github.com/marianogappa/cheesse/issues">Report an issue</a> |
    <a href="https://github.com/marianogappa/cheesse/blob/master/LICENSE">MIT License</a>
  `
  document.body.appendChild(footer)
})

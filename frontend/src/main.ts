import './style.less'
import './css/dark.css'
import '../node_modules/bootstrap/dist/css/bootstrap-grid.min.css'
// import '../node_modules/bootstrap/dist/js/bootstrap.bundle.min.js'

import App from './App.svelte'

const app = new App({
  target: document.getElementById('app')
})

export default app

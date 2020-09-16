import Vue from 'vue'

import "@fortawesome/fontawesome-free/css/all.min.css";

Vue.config.productionTip = false

new Vue({
  data: {
    currentRoute: window.location.pathname
  },
  computed: {
    ViewComponent () {
      return routes[this.currentRoute] || Login
    }
  },
  render (h) { return h(this.ViewComponent) },
}).$mount('#app')

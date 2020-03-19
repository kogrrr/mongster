import Vue from 'vue'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

// Install BootstrapVue
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)

window.onload = () => {
  var vm = new Vue({
    el: '#app',
    data() {
      return {
        name: 'BootstrapVue',
        show: true
      }
    },
    watch: {
      show(newVal) {
        console.log('Alert is now ' + (newVal ? 'visible' : 'hidden'))
      }
    },
    methods: {
      toggle() {
        console.log('Toggle button clicked')
        this.show = !this.show
      },
      dismissed() {
        console.log('Alert dismissed')
      }
    }
  });
  fetch('https://api.ipify.org?format=json')
    .then((resp) => {
      return resp.json()
    })
    .then((body) => {
      vm.name = body.ip
    });
}

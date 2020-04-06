import Vue from 'vue'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './registerServiceWorker'

Vue.config.productionTip = false

// Install BootstrapVue
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import * as firebase from 'firebase'
if (process.env.NODE_ENV === 'production') {
  // Your web app's Firebase configuration
  var firebaseConfig = {
    apiKey: process.env.VUE_APP_API_KEY,
    authDomain: 'tag-mng-app.firebaseapp.com',
    databaseURL: 'https://tag-mng-app.firebaseio.com',
    projectId: 'tag-mng-app',
    storageBucket: 'tag-mng-app.appspot.com',
    messagingSenderId: process.env.VUE_APP_MESSAGING_SENDER_ID,
    appId: process.env.VUE_APP_APP_ID,
    measurementId: process.env.VUE_APP_MEASUREMENT_ID
  }

  // Initialize Firebase
  firebase.initializeApp(firebaseConfig)
  firebase.analytics()
  
  const messaging = firebase.messaging()
  
  const publicKey = process.env.VUE_APP_PUBLIC_KEY
  messaging.usePublicVapidKey(publicKey) // 1で取得した鍵ペア
  
  // 通知の受信許可
  messaging.requestPermission().then(() => {
    console.log('Notification permission granted.')
  
    // トークン取得
    messaging.getToken().then((token) => {
      console.log(token)
    })
  }).catch((err) => {
    console.log('Unable to get permission to notify.', err)
  })
}

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')

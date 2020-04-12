<!--
<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <HelloWorld msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'home',
  components: {
    HelloWorld
  }
}
</script>
-->
<template>
  <div class="home">
    <LoginForm title="Welcom to Tag-Mng App"/>
  </div>
</template>

<script>
import LoginForm from '@/components/Login.vue'

export default {
  name: 'home',
  components: {
    LoginForm
  }
}
</script>

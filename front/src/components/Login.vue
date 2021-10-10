<template>
  <div class="login">
    <h1>{{ title }}</h1>
    <LoginForm testValue="tete" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import LoginForm from '@/components/login/Form.vue'

export default defineComponent({
  name: 'Login',
  props: {
    title: String
  },
  components: {
    LoginForm
  }
})
</script>

<style scoped>
#form-usename {
  margin: 0px 0px 10px 0px;
}
#form-password {
  margin: 0px 0px 10px 0px;
}
</style>

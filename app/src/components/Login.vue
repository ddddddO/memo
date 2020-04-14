<template>
  <div class="login">
    <h1>{{ title }}</h1>
    <b-form>
      <label for="text-username">User Name</label>
      <b-input v-model="userName" name="name" type="text" id="text-username"></b-input>
      <label for="text-password">Password</label>
      <b-input v-model="passWord" name="passwd" type="password" id="text-password"></b-input>
      <b-button pill style="margin: 10px" v-on:click="postLoginForm" type="button" size="sm" variant="primary">Login</b-button>
    </b-form>
  </div>
</template>

<script>
export default {
  name: 'LoginForm',
  props: {
    title: String
  },
  data: () => ({
    userName: '',
    passWord: '',
    endpoint: 'http://localhost:8082/auth'
  }),
  methods: {
    postLoginForm: function () {
      try {
        fetch(
          this.endpoint,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded'
            },
            mode: 'cors',
            credentials: 'include',
            body: 'name=' + this.userName + '&' + 'passwd=' + this.passWord
          })
          .then(function (resp) {
            if (!resp.ok) {
              alert('retry login!')
            } else {
              // this.$router.push('/memos') TODO: 遷移させたい
            }
          })
      } catch (err) {
        console.error(err)
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.login {
  margin : 0px 10px 0px 10px;
}
</style>

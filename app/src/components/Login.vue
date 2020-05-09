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
    <b-modal ref="failed-login" ok-only title="Retry Login!">
      <div class="d-block text-center">
        <h3>Wrong UserName or Password.</h3>
      </div>
    </b-modal>
  </div>
</template>

<script>
import router from '../router'

export default {
  name: 'LoginForm',
  props: {
    title: String
  },
  data: () => ({
    userName: '',
    passWord: '',
    endpoint: ''
  }),
  mounted () {
    this.buildEndpoint()
  },
  methods: {
    postLoginForm: function () {
      let own = this
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
              own.$refs['failed-login'].show()
            } else {
              router.push('/memos')
            }
          })
      } catch (err) {
        console.error(err)
      }
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/auth'
      } else {
        this.endpoint = 'http://localhost:8082/auth'
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

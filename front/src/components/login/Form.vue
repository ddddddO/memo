<template>
  <div id="form-usename">
    <label id="label-username" for="text-username">User Name</label>
    <input id="text-username" type="text" v-model="userName">
  </div>
  <div id="form-password">
    <label id="label-password" for="password">Password</label>
    <input id="password" type="password" placeholder="input over 8chars" v-model="password" />
  </div>
  <div id="form-button">
    <button id="button" type="button" @click="ValidateAndPost(userName, password)">Login</button>
  </div>
  {{ testValue }}
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue'
import router from '../../router'

export default defineComponent({
  name: 'LoginForm',
  props: {
    testValue: String
  },
  setup (props) {
    const endpoint = ref('')
    const decideEndpoint = () => {
      if (process.env.NODE_ENV === 'production') {
        endpoint.value = process.env.VUE_APP_API_ENDPOINT + '/auth'
      } else {
        endpoint.value = 'http://localhost:8081/auth'
      }
    }

    onMounted(decideEndpoint)

    const ValidateAndPost = (userName: string, password: string) => {
      console.log('debug', props.testValue)
      console.log('Endpoint', endpoint.value)

      if (password.length < 8) {
        alert('Password is 8 characters or more')
      }

      // NOTE: prepare) python3 -m http.server 8081
      // TODO: using axios
      try {
        fetch(
          endpoint.value,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/x-www-form-urlencoded'
            },
            mode: 'cors',
            credentials: 'include',
            body: 'name=' + userName + '&' + 'passwd=' + password
          })
          .then(function (response) {
            if (!response.ok) {
              // TODO: ちゃんとエラーハンドリングする。4xx or 5xx
              console.log('aaa')
              alert('Wrong User Name or Password')
            } else {
              router.push('/momos')
            }
          })
      } catch (error) {
        console.error(error)
      }
    }

    return {
      ValidateAndPost
    }
  },
  data () {
    return {
      userName: '',
      password: ''
    }
  },
  watch: {
    userName: function (newUserName, oldUserName) {
      // noop
    },
    password: function (newPassword, oldPassword) {
      if (newPassword.length < 8) {
        console.log('debug', newPassword)
      }
    }
  }
})
</script>

<style scoped>
#label-username {
  margin: 1%;
}
#label-password {
  margin: 1%;
}
#button {
  margin: 10px;
  height: 25px;
  color: #FFFFFF;
  background: #66CC66;
  border-color: #99CC66;
  border-radius:30px;
}
</style>

<template>
  <div class="createtag">
    <div class="tag-form">
      <h3 style="text-align:start;font-size: medium;">Tag Name:</h3>
      <b-form-input rows="10" v-model="name"></b-form-input>
      <b-button pill size="sm" variant="primary" v-on:click="createTag(name)">Create</b-button>
    </div>
  </div>
</template>

<style>
body {
  margin : 0px 10px 0px 10px;
}
button {
  margin : 3px;
}
</style>

<script>
export default {
  name: 'createTag',
  data: () => ({
    name: '',
    tagEndpoint: ''
  }),
  async mounted () {
    this.buildEndpoint()
  },
  methods: {
    createTag: function (name) {
      fetch(this.tagEndpoint, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'POST',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          name: name,
          user_id: 1
        })
      })
      setTimeout(() => { this.$router.push('/tags') }, '500')
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.tagEndpoint = process.env.VUE_APP_API_ENDPOINT + '/tags'
      } else {
        this.tagEndpoint = 'http://localhost:8082/tags'
      }
    }
  }
}
</script>

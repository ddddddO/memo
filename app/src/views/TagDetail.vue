<template>
  <div class="tagDetail">
    <h3 style="text-align:start;font-size: medium;">Name:</h3>
    <b-form-input rows="10" v-model="tag.name"></b-form-input>
    <b-button pill size="sm" variant="danger" v-on:click="updateTag(tag.name)">Update</b-button>
    <b-button pill size="sm" variant="danger" v-on:click="$bvModal.show('confirm-delete')">Delete</b-button>

    <b-modal ref="failed-update" ok-only title="Failed to update...">
      <div class="d-block text-center">
        <h3>Sorry, failed to update.</h3>
      </div>
    </b-modal>
    <b-modal id="confirm-delete" hide-footer title="Delete ?">
      <div class="d-block text-center">
        <h3>{{ tag.name }}</h3>
      </div>
      <b-button block variant="outline-danger" @click="deleteTag">OK!</b-button>
    </b-modal>

  </div>
</template>

<style scoped>
body {
  margin : 0px 10px 0px 10px;
}
button {
  margin : 3px;
}
</style>

<script>
export default {
  name: 'tagDetail',
  data: () => ({
    loading: false,
    tagId: 0,
    tag: null,
    tagEndpoint: ''
  }),
  async created () {
    this.loading = true
    this.buildEndpoint()
    this.tagId = this.$route.params.tag_id
    try {
      this.tag = await fetch(this.tagEndpoint + this.tagId + '?userId=1', {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'GET',
        mode: 'cors',
        credentials: 'include'
      })
        .then(function (resp) {
          const tmp1 = resp.json()
          return tmp1
        })
        .then(function (j) {
          const tmp2 = JSON.stringify(j)
          return tmp2
        })
        .then(function (sj) {
          const tmp3 = JSON.parse(sj)
          return tmp3
        })
    } catch (err) {
      console.error(err)
    }
  },
  methods: {
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.tagEndpoint = process.env.VUE_APP_API_ENDPOINT + '/tags/'
      } else {
        this.tagEndpoint = 'http://localhost:8082/tags/'
      }
    },
    updateTag: function (name) {
      let own = this
      try {
        fetch(this.tagEndpoint + this.tagId + '?userId=1', {
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          method: 'PATCH',
          mode: 'cors',
          credentials: 'include',
          body: JSON.stringify({
            name: name
          })
        })
          .then(function (resp) {
            if (!resp.ok) {
              own.$refs['failed-update'].show()
            }
          })
      } catch (err) {
        console.error(err)
      }
      this.$router.push('/tags')
    },
    deleteTag: function () {
      fetch(this.tagEndpoint + this.tagId, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'DELETE',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1
        })
      })
      setTimeout(() => { this.$router.push('/tags') }, '500')
    }
  }
}
</script>

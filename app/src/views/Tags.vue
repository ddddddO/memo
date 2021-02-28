<template>
  <div class="tags">
    <div v-if="loading">
      Loading...
    </div>
    <div v-else>
      <h1>tags</h1>
      <b-card>
        <b-button pill style="margin: 10px" to="/new_tag" size="sm" variant="primary" >New!</b-button>
        <div>
          <b-table
            id="tags"
            :items="tags"
            :fields="fields"
            :small="small"
            sticky-header
            striped
            hover
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'tag-detail', params: { tag_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
    </div>
  </div>
</template>

<script>
export default {
  name: 'tags',
  data: () => ({
    tagEndpoint: '',
    tags: [],
    loading: false,
    fields: ['id', 'name']
  }),
  async created () {
    this.loading = true
    this.buildEndpoint()
    try {
      this.tags = await fetch(this.tagEndpoint + '?userId=1', {
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
          const tagList = tmp3.tags
          // ALLを除外するため
          tagList.shift()
          return tagList
        })
    } catch (err) {
      console.error(err)
    }

    this.loading = false
  },
  methods: {
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

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
            small
            sticky-header
            striped
            hover
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'tag-detail', params: { tag_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
            <template v-slot:cell()="data">
              <router-link :to="{ name:'memos', params: { tag_id: data.item.id }}">
                <a>{{ data.item.name }}</a>
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
    fields: [
      { key: 'id', label: 'Refer to' },
      { key: 'name', label: 'Search for' }
    ]
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
          for (let i = 0; i < tagList.length; i++) {
            if (tagList[i].id === 1) {
              tagList.splice(i, 1)
            }
          }
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

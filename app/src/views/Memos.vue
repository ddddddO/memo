<template>
  <div class="memos">
    <h1>memos</h1>
    <div v-if="loading">
      Loading...
    </div>
    <div v-if="memoList" class="overflow-auto">
      <b-button pill style="margin: 10px" to="/new_memo" size="sm" variant="primary" >New!</b-button>
      <div class="list">
        <b-table
          id="memo-list-table"
          :items="memoList"
          :fields="fields"
          :per-page="perPage"
          :current-page="currentPage"
          small
        >
          <template v-slot:cell(id)="data">
            <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
              <a>{{ data.value }}</a>
            </router-link>
          </template>
        </b-table>
      </div>
      <b-pagination
        pills
        size="sm"
        align="center"
        v-model="currentPage"
        :total-rows="rows"
        :per-page="perPage"
        aria-controls="memo-list-table"
      ></b-pagination>
    </div>
  </div>
</template>

<style>
body {
  margin : 0px 10px 0px 10px;
}

/* PC */
@media only screen and (min-width : 1024px){
  .list {
    width: 50%;
    margin: auto;
  }
}
</style>

<script>
export default {
  name: 'memos',
  data: () => ({
    loading: false,
    memoList: null,
    perPage: 10,
    currentPage: 1,
    fields: ['id', 'subject'],
    endpoint: ''
  }),
  created () {
    this.fetchData()
  },
  computed: {
    rows () {
      return this.memoList.length
    }
  },
  methods: {
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/memos' + '?userId=1'
      } else {
        this.endpoint = 'http://localhost:8082' + '/memos' + '?userId=1'
      }
    },
    fetchData: function () {
      this.loading = true
      this.memoList = null
      this.buildEndpoint()
      let data = null
      const fetchFunc = async () => {
        try {
          data = await fetch(
            this.endpoint,
            {
              mode: 'cors',
              credentials: 'include',
              headers: { 'Accept': 'application/json' }
            })
            .then(function (resp) {
              return resp.json()
            })
            .then(function (json) {
              const tmp = JSON.stringify(json)
              return tmp
            })
            .then(function (sJson) {
              const tmp = JSON.parse(sJson)
              return tmp.memo_list
            })
        } catch (err) {
          console.error(err)
        }
        this.memoList = data
        this.loading = false
      }
      fetchFunc()
    }
  },
  watch: {
    '$route': 'fetchData'
  }
}
</script>

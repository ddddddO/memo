<template>
  <div class="memos">
    <h1>memos</h1>
    <div v-if="loading">
      Loading...
    </div>
    <div v-if="memoList" class="overflow-auto">
      <b-card title="ALL">
        <b-button pill style="margin: 10px" to="/new_memo" size="sm" variant="primary" >New!</b-button>
        <div class="list"> <!--TODO: このb-table近辺をコンポーネント化して以下で使いまわすようにする-->
          <b-table
            id="memo-list-table"
            :items="memoList"
            :fields="fields"
            :per-page="perPage"
            :current-page="currentPage"
            small
            sticky-header
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
      </b-card>
      <b-card title="After 20 days!">
        <div class="list-after-20">
          <b-table
            id="memo-list-table-after-20"
            :items="memoListAfter20days"
            :fields="fields"
            small
            sticky-header
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
      <b-card title="After 15 days!">
        <div class="list-after-15">
          <b-table
            id="memo-list-table-after-15"
            :items="memoListAfter15days"
            :fields="fields"
            small
            sticky-header
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
      <b-card title="After 11 days!">
        <div class="list-after-11">
          <b-table
            id="memo-list-table-after-11"
            :items="memoListAfter11days"
            :fields="fields"
            small
            sticky-header
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
      <b-card title="After 7 days!">
        <div class="list-after-7">
          <b-table
            id="memo-list-table-after-7"
            :items="memoListAfter7days"
            :fields="fields"
            small
            sticky-header
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
      <b-card title="After 4 days!">
        <div class="list-after-4">
          <b-table
            id="memo-list-table-after-4"
            :items="memoListAfter4days"
            :fields="fields"
            small
            sticky-header
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
      <b-card title="After 1 day!">
        <div class="list-after-1">
          <b-table
            id="memo-list-table-after-1"
            :items="memoListAfter1day"
            :fields="fields"
            small
            sticky-header
          >
            <template v-slot:cell(id)="data">
              <router-link :to="{ name:'memo-detail', params: { memo_id: data.value }}">
                <a>{{ data.value }}</a>
              </router-link>
            </template>
          </b-table>
        </div>
      </b-card>
    </div>
  </div>
</template>

<style>
body {
  margin : 0px 10px 0px 10px;
}
.card {
  margin : 5px 0px 5px 0px;
}

/* PC */
@media only screen and (min-width : 1024px){
  .card {
    width: 33%;
    float: left;
  }
  .list {
    width: 100%;
    margin: auto;
  }
  .list-after-20 {
    width: 100%;
    margin: auto;
  }
  .list-after-15 {
    width: 100%;
    margin: auto;
  }
  .list-after-11 {
    width: 100%;
    margin: auto;
  }
  .list-after-7 {
    width: 100%;
    margin: auto;
  }
  .list-after-4 {
    width: 100%;
    margin: auto;
  }
  .list-after-1 {
    width: 100%;
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
    perPage: 50,
    currentPage: 1,
    fields: ['id', 'subject'],
    endpoint: '',
    memoListAfter20days: null,
    memoListAfter15days: null,
    memoListAfter11days: null,
    memoListAfter7days: null,
    memoListAfter4days: null,
    memoListAfter1day: null
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
        this.fetchDataByCondition(data)
        this.loading = false
      }
      fetchFunc()
    },
    fetchDataByCondition: function (data) {
      // condition is in 0~5
      let case5 = []
      let case4 = []
      let case3 = []
      let case2 = []
      let case1 = []
      let case0 = []
      data.forEach(function (datum) {
        switch (datum.notified_cnt) {
          case 5:
            case5.push(datum)
            break
          case 4:
            case4.push(datum)
            break
          case 3:
            case3.push(datum)
            break
          case 2:
            case2.push(datum)
            break
          case 1:
            case1.push(datum)
            break
          case 0:
            case0.push(datum)
            break
          default:
            break
        }
      })
      this.memoListAfter20days = case5
      this.memoListAfter15days = case4
      this.memoListAfter11days = case3
      this.memoListAfter7days = case2
      this.memoListAfter4days = case1
      this.memoListAfter1day = case0
    }
  },
  watch: {
    '$route': 'fetchData'
  }
}
</script>

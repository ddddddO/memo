<template>
  <div class="memos">
    <h1>This is a memos page</h1>
    <div>
      <b-card-group deck v-for="memo in memoList" v-bind:key="memo" >
        <b-card bg-variant="secondary" text-variant="white" header="After 3 days" class="text-center">
          <b-card-text>
            <router-link :to="{ name:'memo-detail', params: { memo_id: memo.id }}">
              <a>{{ memo.subject }}</a>
            </router-link>
          </b-card-text>
        </b-card>
      </b-card-group>
    </div>
  </div>
</template>

<style>
div.card-deck {
  margin : 0px 10px 0px 10px;
}
</style>

<script>
export default {
  name: 'memos',
  data: () => ({
    loaded: false,
    memoList: null
  }),
  async mounted () {
    this.loaded = false
    try {
      this.memoList = await fetch('http://localhost:8082' + '/memos' + '?userId=1',
      {
        mode: 'cors',
        headers: {'Accept': 'application/json'}
      })
      .then(function (resp) {
        return resp.json()
      })
      .then(function (json) {
        const tmp = JSON.stringify(json)
        // NOTE: apiからのレスポンスに含まれるエスケープ文字列をトリムし、かつ、JSONレスポンスの先頭・末尾の「"」をトリム
        return tmp.replace(/\\"/g, '"').slice(1, -1)
      })
      .then(function (sJson) {
        const tmp = JSON.parse(sJson)
        return tmp.memo_list
      })
      console.log('memoList')
      console.log(this.memoList)
    } catch (err) {
      console.error(err)
    }
  }
}
</script>

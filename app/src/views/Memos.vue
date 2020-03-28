<template>
  <div class="memos">
    <h1>This is an memos page</h1>
    <div v-for="memo in memoList" v-bind:key="memo">
      <h3> {{ memo.id }} {{ memo.subject }}</h3>
      <router-link :to="{ name:'memo-detail', params: { memo_id: memo.id }}">メモ詳細へ</router-link>
    </div>
  </div>
</template>

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

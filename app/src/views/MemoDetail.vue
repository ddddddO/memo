<template>
  <div class="memodetail">
    <h1>This is an memo detail page</h1>
    <h2>ID: {{ $route.params.memo_id }}</h2>
    <h2>Subject: {{ memoDetail.subject }}</h2>
    <h3>{{ memoDetail.content }}</h3>
  </div>
</template>

<script>
export default {
  name: 'memoDetail',
  data: () => ({
    memoDetail: null
  }),
  async mounted () {
    const user_id = 1
    const memo_id = this.$route.params.memo_id
    try {
      this.memoDetail = await fetch(
        'http://localhost:8082' + '/memodetail' + '?userId=' + user_id + '&' + 'memoId=' + memo_id,
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
        return tmp
      })
      console.log('memoDetail')
      console.log(this.memoDetail)
    } catch (err) {
      console.error(err)
    }
  }
}
</script>

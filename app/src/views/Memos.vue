<template>
  <div class="memos">
    <h1>This is an memos page</h1>
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
        return JSON.stringify(json)
      })
      .then(function (sJson) {
        return JSON.parse(sJson)
      })
      console.log('memoList')
      console.log(this.memoList)
    } catch (err) {
      console.error(err)
    }
  }
}
</script>

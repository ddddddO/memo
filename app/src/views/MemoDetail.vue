<template>
  <div class="memodetail">
    <h1>This is an memo detail page</h1>
    <h2>ID: {{ $route.params.memo_id }}</h2>
    <h2>Subject: {{ memoDetail.subject }}</h2>
    <div v-if="!activatedEdit">
      <h3 v-html="memoDetail.content"></h3>
      <button v-on:click="activateEditMemo">Edit</button>
    </div>
    <div v-else>
      <textarea name="content" style="width:100%;" rows="20" v-html="contentForTextarea">
      </textarea>
      <button v-on:click="deactivateEditMemo">Cancel</button>
      <button v-on:click="updateMemo">Update</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'memoDetail',
  data: () => ({
    memoDetail: null,
    contentForTextarea: null,
    activatedEdit: false
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
        // NOTE: apiからのレスポンスに含まれるエスケープ文字をトリムし、かつ、JSONレスポンスの先頭・末尾の「"」をトリムし、かつ末尾の改行コード「\n」をトリム
        return tmp.replace(/\\"/g, '"').slice(1, -3)
      })
      .then(function (sJson) {
        const tmp = JSON.parse(sJson)
        return tmp
      })
    } catch (err) {
      console.error(err)
    }

    this.memoDetail.content = this.convertRNtoBR(this.memoDetail.content)
    this.contentForTextarea = this.convertBRto(this.memoDetail.content)
  },
  methods: {
    activateEditMemo: function () {
      this.activatedEdit = true
    },
    deactivateEditMemo: function () {
      this.activatedEdit = false
    },
    updateMemo: function () {
      alert('Update!')
      // TODO: update後、メモ詳細ページへ遷移(更新済みの内容を出力)
    },
    convertRNtoBR: function (content) {
      return content.replace(/(\\r\\n)/g, '<br>') // windows
    },
    convertBRto: function (content) {
      return content.replace(/<br>/g, '&#010;')
    }
  }
}
</script>

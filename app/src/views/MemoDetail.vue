<template>
  <div class="memodetail">
    <div class="memodetail-tags">
      <h3 style="text-align:start;font-size: medium;">Tags:</h3>
      <b style="font-size: medium;" v-for="tagName in memoDetail.tag_names" :key="tagName">{{ tagName }} / </b>
    </div>
    <h3 style="text-align:start;font-size: medium;">Subject:</h3>
    <div v-if="!activatedEdit" class="memodetail-subject">
      <h2 style="font-size: x-large;">{{ memoDetail.subject }}</h2>
    </div>
    <div v-else>
      <b-form-input rows="10" v-model="memoDetail.subject"></b-form-input>
    </div>
    <h3 style="text-align:start;font-size: medium;">Content:</h3>
    <div v-if="!activatedEdit">
      <h3 style="white-space: pre-wrap;font-size: large;text-align:start;" v-html="memoDetail.content"></h3>
      <b-button pill size="sm" v-on:click="activateEditMemo">Edit</b-button>
      <b-button pill size="sm" variant="danger" v-on:click="deleteMemo">Delete</b-button>
    </div>
    <div v-else>
      <b-form-textarea id="textarea" rows="7" v-model="memoDetail.content"></b-form-textarea>
      <b-button pill size="sm" v-on:click="switchPreviewContent">Preview?</b-button>
      <div v-if="activatedPreviewContent">
        <h3 style="text-align:start;font-size: medium;">Preview Content:</h3>
        <h3 style="white-space: pre-wrap;font-size: large;text-align:start;" v-html="memoDetail.content"></h3>
      </div>
      <b-button pill size="sm" v-on:click="deactivateEditMemo">Cancel</b-button>
      <b-button pill size="sm" variant="danger" v-on:click="updateMemo(memoDetail.content)">Update</b-button>
    </div>
  </div>
</template>

<style>
body {
  margin : 0px 10px 0px 10px;
}
button {
  margin : 3px;
}
</style>

<script>
export default {
  name: 'memoDetail',
  data: () => ({
    memoDetail: null,
    activatedEdit: false,
    activatedPreviewContent: false,
    endpoint: ''
  }),
  async created () {
    this.buildEndpoint()
    const userId = 1
    const memoId = this.$route.params.memo_id
    try {
      this.memoDetail = await fetch(
        this.endpoint + '?userId=' + userId + '&' + 'memoId=' + memoId,
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
          return tmp
        })
    } catch (err) {
      console.error(err)
    }

    this.memoDetail.content = this.convertRNtoBR(this.memoDetail.content)
  },
  methods: {
    activateEditMemo: function () {
      this.activatedEdit = true
    },
    deactivateEditMemo: function () {
      this.activatedEdit = false
    },
    switchPreviewContent: function () {
      this.activatedPreviewContent = !this.activatedPreviewContent
    },
    updateMemo: function (content) {
      fetch(this.endpoint, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'PATCH',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1,
          memo_id: this.$route.params.memo_id,
          memo_subject: this.memoDetail.subject,
          memo_content: content
        })
      })
      this.reloadMemos()
    },
    convertRNtoBR: function (content) {
      return content.replace(/(\\r\\n)/g, '<br>').replace(/(\\n)/g, '<br>') // windows+
    },
    deleteMemo: function () {
      if (!window.confirm('Delete ?')) { return }

      fetch(this.endpoint, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'DELETE',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1,
          memo_id: this.$route.params.memo_id
        })
      })
      this.reloadMemos()
    },
    reloadMemos: function () {
      this.$router.replace('/memos')
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/memodetail'
      } else {
        this.endpoint = 'http://localhost:8082/memodetail'
      }
    }
  }
}
</script>

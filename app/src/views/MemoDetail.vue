<template>
  <div class="memodetail">
    <b-modal ref="failed-update" ok-only title="Failed to update...">
      <div class="d-block text-center">
        <h3>Sorry, failed to update.</h3>
      </div>
    </b-modal>
    <div v-if="loading">
      Loading...
    </div>
    <div class="left">
      <div class="memodetail-tags">
        <h3 style="text-align:start;font-size: medium;">Tags:</h3>
        <b style="font-size: medium;" v-for="tagName in memoDetail.tag_names" :key="tagName">{{ tagName }} / </b>
      </div>
      <div v-if="!activatedEdit" class="memodetail-subject">
        <h3 style="text-align:start;font-size: medium;">Subject:</h3>
        <h2 style="font-size: x-large;">{{ memoDetail.subject }}</h2>
        <h3 style="text-align:start;font-size: medium;">Content:</h3>
        <h3 style="font-size: large;text-align:start;" v-html="compiledMarkdownContent"></h3>
        <b-button pill size="sm" v-on:click="activateEditMemo">Edit</b-button>
        <b-button pill size="sm" variant="danger" v-on:click="$bvModal.show('confirm-delete')">Delete</b-button>
      </div>
      <div v-else>
        <h3 style="text-align:start;font-size: medium;">Subject:</h3>
        <b-form-input rows="10" v-model="memoDetail.subject"></b-form-input>
        <h3 style="text-align:start;font-size: medium;">Content:</h3>
        <b-form-textarea id="textarea" rows="7" v-model="memoDetail.content"></b-form-textarea>
        <b-button pill size="sm" v-on:click="switchPreviewContent">Preview</b-button>
        <b-button pill size="sm" v-on:click="deactivateEditMemo">Cancel</b-button>
        <b-button pill size="sm" variant="danger" v-on:click="updateMemo(memoDetail.subject, memoDetail.content)">Update</b-button>
      </div>
    </div>
    <div class="right" v-if="activatedPreviewContent">
      <h3 style="text-align:start;font-size: medium;">Preview Content:</h3>
      <b-card>
        <!--<b-card-text style="font-size: medium;text-align:start;position:relative; height:550px; overflow-y:scroll;" v-html="compiledMarkdownContent">-->
        <b-card-text style="text-align:start; height:550px; overflow-y:scroll;" v-html="compiledMarkdownContent">
        </b-card-text>
      </b-card>
    </div>
    <b-modal id="confirm-delete" hide-footer title="Delete ?">
      <div class="d-block text-center">
        <h3>{{ memoDetail.subject }}</h3>
      </div>
      <b-button block variant="outline-danger" @click="deleteMemo">OK!</b-button>
    </b-modal>
  </div>
</template>

<style>
body {
  margin : 0px 10px 0px 10px;
}
button {
  margin : 3px;
}

/* PC */
@media only screen and (min-width : 1024px){
  .memodetail {
    overflow: auto;
  }
  .left {
    float: left;
    width: 49%;
    margin-right: 2%;
  }
  .right {
    margin-left: 50%;
  }
  .card {
    width: 570px;
  }
  .card-text {
    width: 550px;
  }
}
</style>

<script>
import marked from 'marked'

export default {
  name: 'memoDetail',
  data: () => ({
    loading: false,
    memoDetail: null,
    activatedEdit: false,
    activatedPreviewContent: false,
    endpoint: ''
  }),
  async created () {
    this.loading = true
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
    this.loading = false
  },
  methods: {
    activateEditMemo: function () {
      this.activatedEdit = true
    },
    deactivateEditMemo: function () {
      this.activatedEdit = false
      this.activatedPreviewContent = false
    },
    switchPreviewContent: function () {
      this.activatedPreviewContent = !this.activatedPreviewContent
    },
    updateMemo: function (subject, content) {
      let own = this
      try {
        fetch(this.endpoint, {
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          method: 'PATCH',
          mode: 'cors',
          credentials: 'include',
          body: JSON.stringify({
            user_id: 1,
            memo_id: this.$route.params.memo_id,
            memo_subject: subject,
            memo_content: content
          })
        })
          .then(function (resp) {
            if (!resp.ok) {
              own.$refs['failed-update'].show()
            }
          })
      } catch (err) {
        console.error(err)
      }
      this.$router.push('/memos')
    },
    convertRNtoBR: function (content) {
      return content.replace(/(\\r\\n)/g, '<br>').replace(/(\\n)/g, '<br>') // windows+
    },
    deleteMemo: function () {
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
      setTimeout(() => { this.$router.push('/memos') }, '500')
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/memodetail'
      } else {
        this.endpoint = 'http://localhost:8082/memodetail'
      }
    }
  },
  computed: {
    compiledMarkdownContent: function () {
      return marked(this.memoDetail.content)
    }
  }
}
</script>

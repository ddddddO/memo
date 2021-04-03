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
    <div v-else class="left">
      <div v-if="!activatedEdit" class="memodetail-subject">
        <div class="memodetail-tags">
          <h3 style="text-align:start;font-size: medium;">Tags:</h3>
          <b style="font-size: medium;" v-for="tag in memoDetail.tags" :key="tag">#{{ tag.name | trimWquote }} </b>
        </div>

        <h3 style="text-align:start;font-size: medium;">Subject:</h3>
        <h2 style="font-size: x-large;">{{ memoDetail.subject }}</h2>
        <h3 style="text-align:start;font-size: medium;">Content:</h3>
        <h3 style="font-size: large;text-align:start;" v-html="compiledMarkdownContent"></h3>
        <b-button pill size="sm" v-on:click="activateEditMemo">Edit</b-button>
        <b-button pill size="sm" variant="danger" v-on:click="$bvModal.show('confirm-delete')">Delete</b-button>
      </div>
      <div v-else>
        <div class="memodetail-tags">
          <b-form-group label="Tags:" style="text-align:start;">
            <b-form-checkbox-group
              id="checkbox-group-1"
              v-model="selectedTagIDs"
              name="tags"
            >
              <!-- TODO: タグの選択は、別にモーダルを表示してそこで選択したい。タグが多すぎる -->
              <b-form-checkbox v-for="tag in tags" :key=tag.name :value=tag.id>{{ tag.name }}</b-form-checkbox>
            </b-form-checkbox-group>
          </b-form-group>
        </div>
        <b-form-checkbox v-model="memoDetail.is_exposed">Expose?</b-form-checkbox>
        <h3 style="text-align:start;font-size: medium;">Subject:</h3>
        <b-form-input rows="10" v-model="memoDetail.subject"></b-form-input>
        <h3 style="text-align:start;font-size: medium;">Content:</h3>
        <b-form-textarea id="textarea" rows="7" v-model="memoDetail.content"></b-form-textarea>
        <b-button pill size="sm" v-on:click="switchPreviewContent">Preview</b-button>
        <b-button pill size="sm" v-on:click="deactivateEditMemo">Cancel</b-button>
        <b-button pill size="sm" variant="danger" v-on:click="updateMemo(memoDetail.subject, memoDetail.content, memoDetail.is_exposed)">Update</b-button>
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
    endpoint: '',
    tagEndpoint: '',
    tags: [],
    selectedTagIDs: []
  }),
  async created () {
    this.loading = true
    this.buildEndpoint()
    const userId = 1
    const memoId = this.$route.params.memo_id
    try {
      this.memoDetail = await fetch(
        this.endpoint + memoId + '?userId=' + userId,
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
          // ALLを除外するため
          for (let i = 0; i < tmp.tags.length; i++) {
            if (tmp.tags[i].id === 1) {
              tmp.tags.splice(i, 1)
            }
          }
          return tmp
        })
    } catch (err) {
      console.error(err)
    }

    this.memoDetail.content = this.convertRNtoBR(this.memoDetail.content)

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
          tagList.shift()
          return tagList
        })
    } catch (err) {
      console.error(err)
    }

    let ids = []
    for (const tag of this.memoDetail.tags) {
      ids.push(tag.id)
    }
    this.selectedTagIDs = ids

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
    updateMemo: function (subject, content, isExposed) {
      let tags = []
      for (const id of this.selectedTagIDs) {
        let tag = { id: id }
        tags.push(tag)
      }

      let own = this
      try {
        let memoID = this.$route.params.memo_id
        fetch(this.endpoint + memoID, {
          headers: { 'Content-Type': 'application/json; charset=utf-8' },
          method: 'PATCH',
          mode: 'cors',
          credentials: 'include',
          body: JSON.stringify({
            user_id: 1,
            subject: subject,
            content: content,
            is_exposed: isExposed,
            tags: tags
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
      let memoID = this.$route.params.memo_id
      fetch(this.endpoint + memoID, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'DELETE',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1
        })
      })
      setTimeout(() => { this.$router.push('/memos') }, '500')
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/memos/'
        this.tagEndpoint = process.env.VUE_APP_API_ENDPOINT + '/tags'
      } else {
        this.endpoint = 'http://localhost:8082/memos/'
        this.tagEndpoint = 'http://localhost:8082/tags'
      }
    }
  },
  computed: {
    compiledMarkdownContent: function () {
      return marked(this.memoDetail.content)
    }
  },
  filters: {
    trimWquote: function (value) {
      // FIXME: タグの先頭2文字が"であるのはなぜか
      return value.substring(2, value.length - 1)
    }
  }
}
</script>

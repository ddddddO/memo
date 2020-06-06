<template>
  <div class="creatememo">
    <div class="memodetail-tags" v-if="loaded">
      <!--<h3 style="text-align:start;font-size: medium;">Tags: FIXME: checkbox</h3>-->
      <b-form-group label="Tags:" style="text-align:start;">
      <b-form-checkbox-group
        id="checkbox-group-1"
        v-model="tagsSelected"
        name="tags"
      >
        <b-form-checkbox v-for="tag in tags" :key=tag.name :value=tag.id>{{ tag.name }}</b-form-checkbox>
      </b-form-checkbox-group>
    </b-form-group>
    </div>
    <div class="memodetail-subject">
      <h3 style="text-align:start;font-size: medium;">Subject:</h3>
      <b-form-input rows="10" v-model="subject"></b-form-input>
    </div>
    <div>
      <h3 style="text-align:start;font-size: medium;">Content:</h3>
      <b-form-textarea id="textarea" rows="15" v-model="content"></b-form-textarea>
      <b-button pill size="sm" variant="primary" v-on:click="createMemo()">Create</b-button>
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
  name: 'createMemo',
  data: () => ({
    loaded: false,
    subject: '',
    content: '',
    tags: null,
    tagsSelected: [],
    endpoint: '',
    tagEndpoint: ''
  }),
  async mounted () {
    this.loaded = false
    this.buildEndpoint()
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
          // NOTE: apiからのレスポンスに含まれるエスケープ文字列をトリムし、かつ、JSONレスポンスの先頭・末尾の「"」をトリム
          return tmp2.replace(/\\"/g, '"').slice(1, -1)
        })
        .then(function (sj) {
          const tmp3 = JSON.parse(sj)
          const tagList = tmp3.tag_list
          return tagList
        })
    } catch (err) {
      console.error(err)
    }
    console.log(this.tags)
    this.loaded = true
  },
  methods: {
    createMemo: function () {
      fetch(this.endpoint, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'POST',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1,
          tag_ids: this.tagsSelected,
          memo_subject: this.subject,
          memo_content: this.content
        })
      })
      this.$router.push('/memos')
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/memodetail'
        this.tagEndpoint = process.env.VUE_APP_API_ENDPOINT + '/tags'
      } else {
        this.endpoint = 'http://localhost:8082/memodetail'
        this.tagEndpoint = 'http://localhost:8082/tags'
      }
    }
  }
}
</script>

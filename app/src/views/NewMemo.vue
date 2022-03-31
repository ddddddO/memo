<template>
  <div class="creatememo">
    <div class="memodetail-tags" v-if="loaded">
      <b-form-group label="Tags:" style="text-align:start;">
        <b-form-checkbox-group
          id="checkbox-group-1"
          v-model="selectedTagIDs"
          name="tags"
        >
          <!-- TODO: タグの選択は、別にモーダルを表示してそこで選択したい。タグが多すぎる -->
          <b-form-checkbox v-for="tag in tags" :key="tag.name" :value="tag.id">{{ tag.name }}</b-form-checkbox>
        </b-form-checkbox-group>
      </b-form-group>
    </div>
    <div>
      <b-form-checkbox v-model="isExposed">Expose?</b-form-checkbox>
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
    isExposed: false,
    subject: '',
    content: '',
    tags: [],
    selectedTagIDs: [],
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
          return tmp2
        })
        .then(function (sj) {
          const tmp3 = JSON.parse(sj)
          const tagList = tmp3.tags
          // ALLを除外するため
          for (let i = 0; i < tagList.length; i++) {
            if (tagList[i].id === 1) {
              tagList.splice(i, 1)
            }
          }
          return tagList
        })
    } catch (err) {
      console.error(err)
    }
    this.loaded = true
  },
  methods: {
    createMemo: function () {
      const tagIDAll = { 'id': 1 }
      let selectedTags = [tagIDAll]
      for (const id of this.selectedTagIDs) {
        let tag = { id: id }
        selectedTags.push(tag)
      }

      fetch(this.endpoint, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'POST',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1,
          is_exposed: this.isExposed,
          tags: selectedTags,
          subject: this.subject,
          content: this.content
        })
      })
      setTimeout(() => { this.$router.push('/memos') }, '500')
    },
    buildEndpoint: function () {
      if (process.env.NODE_ENV === 'production') {
        this.endpoint = process.env.VUE_APP_API_ENDPOINT + '/memos'
        this.tagEndpoint = process.env.VUE_APP_API_ENDPOINT + '/tags'
      } else {
        this.endpoint = 'http://localhost:8082/memos'
        this.tagEndpoint = 'http://localhost:8082/tags'
      }
    }
  }
}
</script>

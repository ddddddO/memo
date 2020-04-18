<template>
  <div class="creatememo">
    <div class="memodetail-tags">
      <h3 style="text-align:start;font-size: medium;">Tags: FIXME: checkbox</h3>
      <!--<h2 style="font-size: x-large;">{{ memoDetail.tag_names }}</h2>-->
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
    subject: '',
    content: '',
    endpoint: 'http://localhost:8082/memodetail'
  }),
  methods: {
    createMemo: function () {
      fetch(this.endpoint, {
        headers: { 'Content-Type': 'application/json; charset=utf-8' },
        method: 'POST',
        mode: 'cors',
        credentials: 'include',
        body: JSON.stringify({
          user_id: 1,
          memo_subject: this.subject,
          memo_content: this.content
        })
      })
      this.$router.push('/memos')
    }
  }
}
</script>

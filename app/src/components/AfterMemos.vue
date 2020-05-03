<template>
  <div class="after-memos">
    <h1>{{ msg }}</h1>
    <div>
      <b-card-group deck v-for="memo in memoList" v-bind:key="memo" >
        <b-card bg-variant="secondary" text-variant="white" header="After 3 days" class="text-center">
          <b-card-text>
            <router-link :to="{ name:'memo-detail', params: { memo_id: memo.id }}">
              <a>{{ memo.subject }}</a>
            </router-link>
          </b-card-text>
        </b-card>
      </b-card-group>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AfterMemos',
  props: {
    msg: String
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>

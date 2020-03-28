<template>
  <div class="memodetail">
    <h1>This is an memo detail page</h1>
    <h2>Memo ID: {{ $route.params.memo_id }}</h2>
  </div>
</template>

<script>
export default {
  name: 'memoDetail'
}
</script>

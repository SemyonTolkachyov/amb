<template>
  <div class="container py-5">
    <div class="row mb-4">
      <h1 class="col-12">Anonymous message board</h1>
    </div>
    <div class="row">
      <Timeline class="col" />
      <Search class="col" />
    </div>
  </div>
</template>

<script>
import Timeline from '@/components/Timeline.vue';
import Search from '@/components/Search.vue';

export default {
  components: {
    Timeline,
    Search,
  },
};
</script>

<style lang="scss" scoped>
.container {
  max-width: 768px;
}
</style>

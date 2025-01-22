<template>
  <div>
    <input @keyup="searchMessages" v-model.trim="query" class="form-control" placeholder="Search...">
    <div class="mt-4">
      <Message v-for="item in messages" :key="item.id" :message="item" />
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import Message from '@/components/Message.vue';

export default {
  data() {
    return {
      query: '',
    };
  },
  computed: mapState({
    messages: (state) => state.searchResults,
  }),
  methods: {
    searchMessages() {
      if (this.query !== this.lastQuery) {
        this.$store.dispatch('searchMessages', this.query);
        this.lastQuery = this.query;
      }
    },
  },
  components: {
    Message,
  },
};
</script>

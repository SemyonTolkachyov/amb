<template>
  <div>
    <form v-on:submit.prevent="createMessage">
      <div class="input-group">
        <input v-model.trim="messageBody" type="text" class="form-control" placeholder="What's happening?">
        <div class="input-group-append">
          <button class="btn btn-primary" type="submit">Message</button>
        </div>
      </div>
    </form>

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
      messageBody: '',
    };
  },
  computed: mapState({
    messages: (state) => state.messages,
  }),
  methods: {
    createMessage() {
      if (this.messageBody.length !== 0) {
        this.$store.dispatch('createMessage', { body: this.messageBody });
        this.messageBody = '';
      }
    },
  },
  components: {
    Message,
  },
};
</script>

import Vuex from 'vuex';
import axios from 'axios';

const BACKEND_URL = 'http://localhost:8080';
const PUSHER_URL = 'ws://localhost:8080/pusher';

const SET_MESSAGES = 'SET_MESSAGES';
const CREATE_MESSAGE = 'CREATE_MESSAGE';
const SEARCH_SUCCESS = 'SEARCH_SUCCESS';
const SEARCH_ERROR = 'SEARCH_ERROR';

function createWebSocketPlugin (socket) {
    return (store) => {
        socket.onmessage = ('MESSAGE_CREATED', event => {
            store.commit(CREATE_MESSAGE, JSON.parse(event.data))
        })
    }
}

const socketPlugin = createWebSocketPlugin (new WebSocket(PUSHER_URL));

const store = new Vuex.Store({
    state: {
        messages: [],
        searchResults: [],
    },
    plugins: [socketPlugin],
    mutations: {
        [SET_MESSAGES](state, messages) {
            state.messages = messages;
        },
        [CREATE_MESSAGE](state, message) {
            state.messages = [message, ...state.messages];
        },
        [SEARCH_SUCCESS](state, messages) {
            state.searchResults = messages;
        },
        [SEARCH_ERROR](state) {
            state.searchResults = [];
        },
    },
    actions: {
        getMessages({ commit }) {
            axios
                .get(`${BACKEND_URL}/messages`)
                .then(({ data }) => {
                    commit(SET_MESSAGES, data ?? []);
                })
                .catch((err) => console.error(err));
        },
        async createMessage({ commit }, message) {
            await axios.post(`${BACKEND_URL}/messages`, null, {
                params: {
                    body: message.body,
                },
            });
        },
        async searchMessages({ commit }, query) {
            if (query.length === 0) {
                commit(SEARCH_SUCCESS, []);
                return;
            }
            axios
                .get(`${BACKEND_URL}/search`, {
                    params: { query },
                })
                .then(({ data }) => commit(SEARCH_SUCCESS, data))
                .catch((err) => {
                    console.error(err);
                    commit(SEARCH_ERROR);
                });
        },
    },
});

store.dispatch('getMessages').then();

export default store;

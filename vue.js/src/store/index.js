import Vuex from "vuex";

export default new Vuex.Store({
    state: {
        user: {
            token: null,
            loggedIn: false
        },
    },
    mutations: {
        LOGIN(state, val) {
            state.user.loggedIn = true;
            state.user.token = val.accessToken;
            window.sessionStorage.setItem('token', val.accessToken);
        },

        LOGOUT(state) {
            state.user.token = null;
            state.user.loggedIn = false;
            window.sessionStorage.removeItem('token');
        },
    },
    actions: {
        login({ commit }, val) {
            commit('LOGIN', val);
        },

        logout({ commit }) {
            commit('LOGOUT');
        },
    }
});
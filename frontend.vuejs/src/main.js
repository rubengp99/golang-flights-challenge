import * as VueRouter from 'vue-router';
import { createApp } from 'vue'
import App from './App.vue'
import vueheader from 'vue-head';
import loader from '@/components/Loading.vue';
import store from './store'
import router from './router'
import Vuex from "vuex";
import '@mdi/font/css/materialdesignicons.css';
import { aliases, mdi } from 'vuetify/iconsets/mdi';
// Vuetify
import 'vuetify/styles';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import Toastify from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';
import DayJsAdapter from '@date-io/dayjs'


const vuetify = createVuetify({
  components,
  directives,
  date: {
    adapter: DayJsAdapter,
  },
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi,
    },
  },
});

const app = createApp(App);

let options = {
    progressBar: true
};

app.config.productionTip = true;
let token = window.sessionStorage.getItem('token');
// this allow skipping auth when closing tab if there's a valid token stored
if (token) {
    store.state.user.loggedIn = true;
    store.state.user.token = token;
}

app.use(Toastify, {
    autoClose: 3000,
    position: 'top-center',
  });
app.use(VueRouter)
app.use(router);
app.use(Vuex);
app.use(store);
app.use(vuetify);
app.component('loader', loader);
app.use(vueheader);
app.use(options);
app.mount('#app');

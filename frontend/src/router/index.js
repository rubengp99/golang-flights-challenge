import * as VueRouter from 'vue-router';
import store from '@/store';

const routes = [
    {
        path: '/',
        name: 'dashboard',
        component: () => import('../views/dashboard/Home.vue'),
        meta:{
            auth:true,
        },
    },
    {
        path: '/login',
        name: 'login',
        component: () => import('../views/auth/Login.vue'),
        meta:{
            auth:false,
            transitionName: 'slide'
        }
    },
]
const router = VueRouter.createRouter({
    routes,
    history: VueRouter.createWebHashHistory(),
    mode: "history",
    base: '/',
    linkActiveClass: 'router-link-active',
    linkExactActiveClass: 'router-link-exact-active',
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { x: 0, y: 0 }
        }
    },
    parseQuery: q => q,
    fallback: true,
})

router.beforeEach((to,from,next) => {
    let user = store.state.user.loggedIn;

    if(to.meta.auth){
        if(user){
            next();
        }else{
            next({name: 'login'});
        }
    }else{
        next();
    }
});

export default router
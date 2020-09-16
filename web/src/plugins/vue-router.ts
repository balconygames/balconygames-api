import Vue from 'vue';
import Router from 'vue-router';

import Login from '../login/Login.vue';
import Logout from '../logout/Logout.vue';

// import Info from '../Info.vue';
// import User from '../User.vue';
// import Admin from '../Admin.vue';
// --------------------------------------------------------------------
// vue-router CONFIGURATION
// --------------------------------------------------------------------

Vue.use(Router);

const router = new Router({
  mode: 'history',
  base: '/',
  routes: [
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
    // {
    //   path: '/',
    //   name: 'main-user',
    //   component: User,
    //   children: [
    //     {
    //       path: '',
    //       name: 'user1',
    //       meta: {
    //         auth: true,
    //       },
    //       component: Info,
    //     },
    //     {
    //       path: 'user',
    //       name: 'user2',
    //       meta: {
    //         auth: 'ROLE_ADMIN',
    //       },
    //       component: Info,
    //     },
    //   ],
    // },
    // {
    //   path: '/admin',
    //   meta: {
    //     auth: true,
    //   },
    //   component: Admin,
    //   children: [
    //     {
    //       path: '',
    //       name: 'admin1',
    //       meta: {
    //         auth: ['ROLE_ADMIN'],
    //       },
    //       component: Info,
    //     },
    //     {
    //       path: 'auth',
    //       name: 'admin2',
    //       meta: {
    //         auth: ['ROLE_UNKNOWN'],
    //       },
    //       component: Info,
    //     },
    //   ],
    // },
    {
      path: '/logout',
      component: Logout,
    },
    {
      path: '*',
      redirect: '/',
    },
  ],
});

(Vue as any).router = router;

export default router;
import Vue from 'vue';
import Router from 'vue-router';
import Dashboard from './view/Dashboard';

Vue.use(Router);

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: Dashboard,
    },

    {
      path: '/import',
      name: 'import',
      component: () => import('./view/Import'),
    },

    {
      path: '/export',
      name: 'export',
      component: () => import('./view/Export'),
    },

    {
      path: '/wine/:id',
      component: () => import('./view/Wine'),
    },

    {
      path: '/catalog/new',
      name: 'create-catalog',
      component: () => import('./view/CreateCatalog'),
    },
  ],
});

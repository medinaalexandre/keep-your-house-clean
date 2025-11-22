import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/Login.vue';
import Dashboard from '../components/Dashboard.vue';
import ComplimentsHistory from '../components/ComplimentsHistory.vue';
import UserTasksHistory from '../components/UserTasksHistory.vue';

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/compliments-history',
    name: 'ComplimentsHistory',
    component: ComplimentsHistory,
    meta: { requiresAuth: true }
  },
  {
    path: '/user/:userId/tasks',
    name: 'UserTasksHistory',
    component: UserTasksHistory,
    meta: { requiresAuth: true }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token');
  const isAuthenticated = !!token;

  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: 'Login' });
  } else if (to.name === 'Login' && isAuthenticated) {
    next({ name: 'Dashboard' });
  } else {
    next();
  }
});

export default router;


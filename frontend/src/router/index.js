import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/Search.vue')
  },
  {
    path: '/destination/:id',
    name: 'Destination',
    component: () => import('../views/Destination.vue')
  },
  {
    path: '/product/:id',
    name: 'Product',
    component: () => import('../views/Product.vue')
  },
  {
    path: '/category/:category',
    name: 'Category',
    component: () => import('../views/Category.vue')
  },
  {
    path: '/city/:city',
    name: 'City',
    component: () => import('../views/City.vue')
  },
  {
    path: '/trips',
    name: 'Trips',
    component: () => import('../views/Trips.vue')
  },
  {
    path: '/account',
    name: 'Account',
    component: () => import('../views/Account.vue')
  },
  {
    path: '/platform',
    name: 'Platform',
    component: () => import('../views/Platform.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

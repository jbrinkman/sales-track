import { createRouter, createWebHashHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import SalesDetails from '../views/SalesDetails.vue'
import Reports from '../views/Reports.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/details',
    name: 'SalesDetails',
    component: SalesDetails
  },
  {
    path: '/reports',
    name: 'Reports',
    component: Reports
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router

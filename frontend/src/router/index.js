import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import AdminUsersView from '../views/AdminUsersView.vue'
import BlogListView from '../views/BlogListView.vue'
import CreateBlogView from '../views/CreateBlogView.vue'
import BlogDetailsView from '../views/BlogDetailsView.vue'
import ProfileView from '../views/ProfileView.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/admin/users',
      name: 'adminUsers',
      component: AdminUsersView
    },
    {
      path: '/blogs',
      name: 'blogs',
      component: BlogListView
    },
    {
      path: '/create-blog',
      name: 'createBlog',
      component: CreateBlogView
    },
    {
      path: '/blogs/:id',
      name: 'blogDetails',
      component: BlogDetailsView
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView
    },
  ]
})

export default router
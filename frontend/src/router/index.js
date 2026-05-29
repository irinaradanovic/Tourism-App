import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import AdminUsersView from '../views/AdminUsersView.vue'
import BlogListView from '../views/BlogListView.vue'
import CreateBlogView from '../views/CreateBlogView.vue'
import BlogDetailsView from '../views/BlogDetailsView.vue'
import MyProfileView from '../views/MyProfileView.vue'
import UserProfileView from '../views/UserProfileView.vue'
import PositionSimulatorView from '../views/PositionSimulatorView.vue'
import CreateTourView from "@/views/CreateTourView.vue";
import MyToursView from "@/views/MyToursView.vue";
import TourDetailsView from "@/views/TourDetailsView.vue";
import TourExploreView from "@/views/TourExploreView.vue";
import CartView from "@/views/CartView.vue";
import TourDetailsTouristView from "@/views/TourDetailsTouristView.vue";

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
      component: MyProfileView
    },
    {
      path: '/:id/profile',
      name: 'user-profile',
      component: UserProfileView
    },
    { 
      path: '/simulator', 
      name: 'simulator', 
      component: PositionSimulatorView 
    },
    {
      path: '/create-tour',
      name: 'create-tour',
      component: CreateTourView
    },
    {
      path: '/my-tours',
      name: 'my-tours',
      component: MyToursView
    },
    {
      path: '/tours/:id',
      name: 'tourDetails',
      component: TourDetailsView
    },
    {
      path: '/tours',
      name: 'exploreTours',
      component: TourExploreView
    },
    {
      path: '/cart',
      name: 'cart',
      component: CartView
    },
    {
      path: '/tours-tourist/:id',
      name: 'tourDetailsTourist',
      component: TourDetailsTouristView
    }

  ]
})

export default router
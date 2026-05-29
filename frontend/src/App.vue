<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-logo">
        <router-link to="/">TourismApp</router-link>
      </div>
      <div class="nav-links">
        <router-link to="/blogs" v-if="user && user.role !== 'ADMIN'">Blogs</router-link>

        <router-link to="/tours" v-if="user && user.role === 'TOURIST'">Explore Tours</router-link>
        
        <template v-if="user">
          <router-link to="/profile">My Profile</router-link>

          <router-link to="/simulator" v-if="user.role !== 'ADMIN'">Simulator</router-link>
          
          <template v-if="user.role === 'GUIDE'">
            <router-link to="/create-tour">Create Tour</router-link>
            <router-link to="/my-tours">My Tours</router-link>
          </template>

          <router-link v-if="user.role === 'ADMIN'" to="/admin/users">Admin Panel</router-link>

          <router-link v-if="user.role === 'TOURIST'" to="/cart" class="nav-cart-link" title="View Cart">
            <div class="cart-icon-wrapper">
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="cart-svg">
                <circle cx="9" cy="21" r="1"></circle>
                <circle cx="20" cy="21" r="1"></circle>
                <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path>
              </svg>
              <span class="header-badge" v-if="cartItemsCount > 0">{{ cartItemsCount }}</span>
            </div>
          </router-link>

          <span class="user-greeting">Hi, {{ user.username }}</span>
          <button @click="handleLogout" class="btn-logout">Logout</button>
        </template>

        <template v-else>
          <router-link to="/login">Login</router-link>
          <router-link to="/register" class="btn-signup">Sign Up</router-link>
        </template>
      </div>
    </nav>

    <main class="main-content">
      <router-view @auth-change="handleAuthChange" />
    </main>
  </div>
</template>

<script>
import { logout } from './services/authService'
import { purchaseService } from './services/purchaseService'

export default {
  data() {
    return { 
      user: null,
      cartItemsCount: 0
    }
  },
  created() {
    this.updateUser();
    window.addEventListener('cart-updated', this.fetchCartCount);
  },
  unmounted() {
    window.removeEventListener('cart-updated', this.fetchCartCount);
  },
  methods: {
    updateUser() {
      const stored = localStorage.getItem('user');
      this.user = stored ? JSON.parse(stored) : null;
      
      if (this.user && this.user.role === 'TOURIST') {
        this.fetchCartCount();
      } else {
        this.cartItemsCount = 0;
      }
    },
    async fetchCartCount() {
      if (!this.user || this.user.role !== 'TOURIST') return;
      try {
        const res = await purchaseService.getCart();
        this.cartItemsCount = res.data.items ? res.data.items.length : 0;
      } catch (err) {
        console.error("Error fetching cart count for header:", err);
      }
    },
    handleAuthChange() {
      this.updateUser();
    },
    async handleLogout() {
      localStorage.removeItem('user');
      localStorage.removeItem('token');
      this.user = null;
      this.cartItemsCount = 0;
      try { await logout(); } catch (e) { console.error(e); }
      this.$router.push('/login');
    }
  }
}
</script>

<style>
body {
  margin: 0;
  font-family: 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  background-color: #f4f7f6;
  color: #333;
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 40px;
  height: 70px;
  background: white;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.nav-logo a {
  font-size: 1.5rem;
  font-weight: bold;
  color: #28a745;
  text-decoration: none;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-links a {
  text-decoration: none;
  color: #555;
  font-weight: 500;
  transition: color 0.3s;
}

.nav-links a:hover {
  color: #28a745;
}

.nav-cart-link {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  border-radius: 50%;
  transition: background 0.2s;
}

.nav-cart-link:hover {
  background-color: #f4faf5;
}

.cart-icon-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.cart-svg {
  stroke: #555;
  transition: stroke 0.3s;
}

.nav-cart-link:hover .cart-svg {
  stroke: #28a745;
}

.header-badge {
  position: absolute;
  top: -7px;
  right: -9px;
  background-color: #dc3545;
  color: white;
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 0.7rem;
  font-weight: bold;
  line-height: 1;
}

.user-greeting {
  font-size: 0.9rem;
  color: #888;
  border-left: 1px solid #ddd;
  padding-left: 20px;
}

.btn-logout {
  background: none;
  border: 1px solid #ff4d4d;
  color: #ff4d4d;
  padding: 6px 15px;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-logout:hover {
  background: #ff4d4d;
  color: white;
}

.btn-signup {
  background: #28a745;
  color: white !important;
  padding: 8px 20px;
  border-radius: 5px;
}

.main-content {
  padding: 40px 20px;
}
</style>
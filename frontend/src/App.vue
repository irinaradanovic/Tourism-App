<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-logo">
        <router-link to="/">TourismApp</router-link>
      </div>
      <div class="nav-links">
        <!-- Always visible -->
        <router-link to="/blogs" v-if="user && user.role !== 'ADMIN'">Blogs</router-link>


        
        <!-- Only if logged in -->
        <template v-if="user">
          <router-link to="/profile">My Profile</router-link>

          <router-link to="/simulator" v-if="user && user.role !== 'ADMIN'">Simulator</router-link>
          <!-- GUIDE ONLY -->
          <template v-if="user.role === 'GUIDE'">
            <router-link to="/create-tour">Create Tour</router-link>
            <router-link to="/my-tours">My Tours</router-link>
          </template>

          <!-- ADMIN ONLY -->
          <router-link v-if="user.role === 'ADMIN'" to="/admin/users">Admin Panel</router-link>

          <span class="user-greeting">Hi, {{ user.username }}</span>
          <button @click="handleLogout" class="btn-logout">Logout</button>
        </template>

        <!-- Only if not logged in -->
        <template v-else>
          <router-link to="/login">Login</router-link>
          <router-link to="/register" class="btn-signup">Sign Up</router-link>
        </template>
      </div>
    </nav>

    <main class="main-content">
      <router-view @auth-change="updateUser" />
    </main>
  </div>
</template>

<script>
import { logout } from './services/authService'

export default {
  data() {
    return { user: null }
  },
  created() {
    this.updateUser();
  },
  methods: {
    updateUser() {
      const stored = localStorage.getItem('user');
      this.user = stored ? JSON.parse(stored) : null;
    },
    async handleLogout() {
      localStorage.removeItem('user');
      localStorage.removeItem('token');
      this.user = null;
      try { await logout(); } catch (e) { console.error(e); }
      this.$router.push('/login');
    }
  }
}
</script>

<style>
/* Globalni stilovi */
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
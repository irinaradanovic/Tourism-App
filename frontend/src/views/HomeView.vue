<template>
  <div class="hero-container">
    <div class="hero-content">
      <h1>Welcome to TourismApp</h1>
      <p>Discover beautiful destinations and share your experiences.</p>
      
      <div v-if="!user" class="hero-actions">
        <router-link to="/register" class="btn-primary">Get Started</router-link>
        <router-link to="/login" class="btn-secondary">Sign In</router-link>
      </div>

      <div v-else class="user-status">
        <p>Currently active as: <strong>{{ user.username }}</strong></p>
        <router-link v-if="user.role === 'ADMIN'" to="/admin/users" class="btn-primary">Admin Panel</router-link>
        <router-link v-else to="/blogs" class="btn-primary">Explore Blogs</router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return { user: JSON.parse(localStorage.getItem('user')) }
  }
}
</script>

<style scoped>
.hero-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60vh;
  text-align: center;
}

.hero-content h1 {
  font-size: 3.5rem;
  color: #2c3e50;
  margin-bottom: 20px;
}

.hero-content p {
  font-size: 1.2rem;
  color: #7f8c8d;
  margin-bottom: 40px;
}

.hero-actions, .user-status {
  display: flex;
  gap: 15px;
  justify-content: center;
}

.btn-primary {
  background: #28a745;
  color: white;
  padding: 12px 30px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
}

.btn-secondary {
  border: 2px solid #28a745;
  color: #28a745;
  padding: 10px 30px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
}
</style>
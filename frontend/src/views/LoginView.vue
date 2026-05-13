<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2>Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="input-field">
          <input v-model="form.username" type="text" placeholder="Username" required />
        </div>
        <div class="input-field">
          <input v-model="form.password" type="password" placeholder="Password" required />
        </div>
        <button type="submit" class="btn-auth">Sign In</button>
      </form>
      <p v-if="message" class="error-msg">{{ message }}</p>
      <p class="auth-footer">
        Don't have an account? <router-link to="/register">Create one</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { login } from '../services/authService'

export default {
  data() {
    return { form: { username: '', password: '' }, message: '' }
  },
  methods: {
    async handleLogin() {
      try {
        const response = await login(this.form)
        localStorage.setItem('token', response.data.token)
        localStorage.setItem('user', JSON.stringify(response.data.user))
        
        // Emitujemo događaj da se Header osveži
        this.$emit('auth-change');

        if (response.data.user.role === 'ADMIN') {
          this.$router.push('/admin/users')
        } else {
          this.$router.push('/')
        }
      } catch (e) {
        this.message = 'Invalid username or password!'
      }
    }
  }
}
</script>

<style scoped>
.auth-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 70vh;
}

.auth-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.05);
  width: 100%;
  max-width: 400px;
  text-align: center;
}

h2 { margin-bottom: 30px; color: #2c3e50; }

.input-field { margin-bottom: 20px; }

input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-sizing: border-box;
}

.btn-auth {
  width: 100%;
  background: #28a745;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.error-msg { color: #e74c3c; margin-top: 15px; font-size: 0.9rem; }
.auth-footer { margin-top: 20px; font-size: 0.9rem; color: #7f8c8d; }
</style>
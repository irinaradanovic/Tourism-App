<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h2>Create Account</h2>
      <p class="auth-subtitle">Join our community and start exploring</p>
      
      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="input-field">
          <label>Username</label>
          <input 
            v-model="form.username" 
            type="text" 
            placeholder="Choose a username" 
            required 
          />
        </div>

        <div class="input-field">
          <label>Email Address</label>
          <input 
            v-model="form.email" 
            type="email" 
            placeholder="your@email.com" 
            required 
          />
        </div>

        <div class="input-field">
          <label>Password</label>
          <input 
            v-model="form.password" 
            type="password" 
            placeholder="Create a password" 
            required 
          />
        </div>

        <div class="input-field">
          <label>Role</label>
          <select v-model="form.role" required class="role-select">
            <option value="" disabled selected>Select your role</option>
            <option value="TOURIST">Tourist</option>
            <option value="GUIDE">Guide</option>
          </select>
        </div>

        <button type="submit" class="btn-auth">Sign Up</button>
      </form>

      <p v-if="message" :class="['message', messageClass]">{{ message }}</p>
      
      <p class="auth-footer">
        Already have an account? <router-link to="/login">Sign In</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { register } from '../services/authService'

export default {
  data() {
    return {
      form: { username: '', email: '', password: '', role: '' },
      message: '',
      messageClass: ''
    }
  },
  methods: {
    async handleRegister() {
      try {
        await register(this.form)
        this.message = 'Registration successful! Redirecting to login...'
        this.messageClass = 'success'
        setTimeout(() => {
          this.$router.push('/login')
        }, 2000)
      } catch (e) {
        this.message = 'Error during registration. Please try again.'
        this.messageClass = 'error'
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
  min-height: 80vh;
  padding: 20px;
}

.auth-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.05);
  width: 100%;
  max-width: 450px;
  text-align: center;
}

h2 {
  margin-bottom: 5px;
  color: #2c3e50;
  font-size: 2rem;
}

.auth-subtitle {
  color: #7f8c8d;
  margin-bottom: 30px;
  font-size: 0.95rem;
}

.auth-form {
  text-align: left;
}

.input-field {
  margin-bottom: 20px;
}

.input-field label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  color: #34495e;
}

input, .role-select {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-sizing: border-box;
  font-size: 1rem;
  transition: border-color 0.3s;
}

input:focus, .role-select:focus {
  outline: none;
  border-color: #28a745;
}

.btn-auth {
  width: 100%;
  background: #28a745;
  color: white;
  border: none;
  padding: 14px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.3s;
  margin-top: 10px;
}

.btn-auth:hover {
  background: #218838;
}

.message {
  margin-top: 20px;
  font-size: 0.9rem;
  padding: 10px;
  border-radius: 5px;
}

.success { background-color: #d4edda; color: #155724; }
.error { background-color: #f8d7da; color: #721c24; }

.auth-footer {
  margin-top: 25px;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.auth-footer a {
  color: #28a745;
  text-decoration: none;
  font-weight: bold;
}
</style>
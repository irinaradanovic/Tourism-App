<template>
  <div class="container">
    <h2>Prijava</h2>
    <form @submit.prevent="handleLogin">
      <input v-model="form.username" placeholder="Korisničko ime" required />
      <input v-model="form.password" type="password" placeholder="Lozinka" required />
      <button type="submit">Prijavi se</button>
    </form>
    <p v-if="message">{{ message }}</p>
    <p>Nemaš nalog? <router-link to="/register">Registruj se</router-link></p>
  </div>
</template>

<script>
import { login } from '../services/authService'

export default {
  data() {
    return {
      form: { username: '', password: '' },
      message: ''
    }
  },
  methods: {
    async handleLogin() {
      try {
        const response = await login(this.form)
        localStorage.setItem('user', JSON.stringify({
          ...response.data,
          password: this.form.password
        }))
        if (response.data.role === 'ADMIN') {
          this.$router.push('/admin/users')
        } else {
          this.$router.push('/')
        }
      } catch (e) {
        this.message = 'Pogrešno korisničko ime ili lozinka!'
      }
    }
  }
}
</script>
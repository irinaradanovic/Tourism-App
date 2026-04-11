<template>
  <div class="container">
    <h2>Registracija</h2>
    <form @submit.prevent="handleRegister">
      <input v-model="form.username" placeholder="Korisničko ime" required />
      <input v-model="form.email" type="email" placeholder="Email" required />
      <input v-model="form.password" type="password" placeholder="Lozinka" required />
      <select v-model="form.role" required>
        <option value="">Izaberi ulogu</option>
        <option value="TOURIST">Turista</option>
        <option value="GUIDE">Vodič</option>
      </select>
      <button type="submit">Registruj se</button>
    </form>
    <p v-if="message">{{ message }}</p>
    <p>Već imaš nalog? <router-link to="/login">Prijavi se</router-link></p>
  </div>
</template>

<script>
import { register } from '../services/authService'

export default {
  data() {
    return {
      form: { username: '', email: '', password: '', role: '' },
      message: ''
    }
  },
  methods: {
    async handleRegister() {
      try {
        await register(this.form)
        this.message = 'Registracija uspešna! Možeš se prijaviti.'
        this.$router.push('/login')
      } catch (e) {
        this.message = 'Greška pri registraciji!'
      }
    }
  }
}
</script>
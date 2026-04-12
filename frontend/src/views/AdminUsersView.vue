<template>
  <div class="container">
    <h2>Svi korisnici</h2>
    <button @click="logout">Odjavi se</button>
    <table v-if="users.length > 0">
      <thead>
        <tr>
          <th>ID</th>
          <th>Korisničko ime</th>
          <th>Email</th>
          <th>Uloga</th>
          <th>Blokiran</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.id }}</td>
          <td>{{ user.username }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.role }}</td>
          <td>{{ user.blocked ? 'Da' : 'Ne' }}</td>
        </tr>
      </tbody>
    </table>
    <p v-else>Nema korisnika.</p>
    <p v-if="message">{{ message }}</p>
  </div>
</template>

<script>
import { getAllUsers } from '../services/authService'

export default {
  data() {
    return {
      users: [],
      message: ''
    }
  },
  async mounted() {
    try {
        const user = JSON.parse(localStorage.getItem('user'))
        if (!user || user.role !== 'ADMIN') {
            this.$router.push('/login')
            return
        }
        const response = await getAllUsers()
        this.users = response.data
    } catch (e) {
        this.message = 'Greška pri učitavanju korisnika!'
    }
  },
  methods: {
    logout() {
      localStorage.removeItem('user')
      this.$router.push('/login')
    }
  }
}
</script>
<template>
  <div class="container">
    <h1>Tourism App</h1>
    <p>Dobrodošli u našu turističku aplikaciju!</p>

    <div class="buttons" v-if="!user">
      <router-link to="/register">
        <button>Registruj se</button>
      </router-link>
      <router-link to="/login">
        <button>Prijavi se</button>
      </router-link>
    </div>

    <div v-if="user">
      <p>Prijavljeni ste kao: <strong>{{ user.username }}</strong> ({{ user.role }})</p>
      <router-link to="/admin/users" v-if="user.role === 'ADMIN'">
        <button>Pregled korisnika</button>
      </router-link>
      <button @click="logout">Odjavi se</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      user: null
    }
  },
  mounted() {
    const stored = localStorage.getItem('user')
    if (stored) {
      this.user = JSON.parse(stored)
    }
  },
  methods: {
    logout() {
      localStorage.removeItem('user')
      this.user = null
      this.$router.push('/')
    }
  }
}
</script>
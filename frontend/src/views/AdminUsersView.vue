<template>
  <div class="admin-container">
    <header class="admin-header">
      <h2>User Management</h2>
      <span class="badge">Admin Only</span>
    </header>

    <div class="table-responsive">
      <table v-if="users.length > 0">
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Email</th>
            <th>Role</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>#{{ user.id }}</td>
            <td class="bold">{{ user.username }}</td>
            <td>{{ user.email }}</td>
            <td><span class="role-tag">{{ user.role }}</span></td>
            <td>
              <span :class="user.blocked ? 'status-blocked' : 'status-active'">
                {{ user.blocked ? 'Blocked' : 'Active' }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state">No users found.</div>
    </div>
    <p v-if="message" class="error-msg">{{ message }}</p>
  </div>
</template>

<script>
import { getAllUsers } from '../services/authService'

export default {
  data() { return { users: [], message: '' } },
  async mounted() {
    try {
      const response = await getAllUsers();
      this.users = response.data;
    } catch (e) {
      this.message = 'Error loading users list!';
    }
  }
}
</script>

<style scoped>
.admin-container {
  max-width: 1000px;
  margin: 0 auto;
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.badge {
  background: #fff3cd;
  color: #856404;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  text-align: left;
  padding: 15px;
  background: #f8f9fa;
  color: #7f8c8d;
  font-size: 0.9rem;
}

td {
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.bold { font-weight: 600; }

.role-tag {
  background: #e8f5e9;
  color: #2e7d32;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.8rem;
}

.status-active { color: #28a745; font-weight: 500; }
.status-blocked { color: #e74c3c; font-weight: 500; }
</style>
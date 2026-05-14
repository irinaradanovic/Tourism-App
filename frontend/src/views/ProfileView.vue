<template>
  <div class="profile-container" v-if="profile">

    <div class="profile-card">

      <div class="profile-header">

        <img
            :src="profile.profileImage || defaultImage"
            alt="Profile image"
            class="profile-image"
        />

        <div class="profile-main-info">

          <h1>
            {{ profile.firstName || 'No first name' }}
            {{ profile.lastName || '' }}
          </h1>

          <p class="username">
            @{{ profile.username }}
          </p>

          <p class="role">
            {{ profile.role }}
          </p>

          <p class="motto" v-if="profile.motto">
            "{{ profile.motto }}"
          </p>

        </div>

      </div>

      <div class="profile-section">

        <h2>Biography</h2>

        <p v-if="profile.biography">
          {{ profile.biography }}
        </p>

        <p v-else class="empty-text">
          No biography added yet.
        </p>

      </div>

      <div class="profile-section">

        <h2>Email</h2>

        <p>{{ profile.email }}</p>

      </div>

    </div>

  </div>
</template>

<script>
import { userService } from '@/services/userService'

export default {

  data() {

    return {
      profile: null,
      defaultImage: 'https://placehold.co/200x200'
    }
  },

  async created() {

    try {

      const response = await userService.getMyProfile()

      this.profile = response.data

    } catch (err) {

      console.error('Error fetching profile:', err)
    }
  }
}
</script>

<style scoped>

.profile-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  background: white;
  border-radius: 10px;
  padding: 40px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.08);
}

.profile-header {
  display: flex;
  gap: 30px;
  align-items: center;
  margin-bottom: 40px;
}

.profile-image {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #28a745;
}

.profile-main-info h1 {
  margin: 0;
  color: #333;
  font-size: 2rem;
}

.username {
  color: #888;
  margin-top: 5px;
}

.role {
  color: #28a745;
  font-weight: bold;
  margin-top: 10px;
}

.motto {
  margin-top: 20px;
  font-style: italic;
  color: #28a745;
  font-size: 1.1rem;
}

.profile-section {
  margin-top: 30px;
}

.profile-section h2 {
  color: #28a745;
  margin-bottom: 10px;
}

.profile-section p {
  line-height: 1.6;
  color: #555;
}

.empty-text {
  color: #999;
  font-style: italic;
}

@media (max-width: 700px) {

  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .profile-card {
    padding: 25px;
  }
}
</style>
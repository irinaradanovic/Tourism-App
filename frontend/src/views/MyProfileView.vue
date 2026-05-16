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
            {{ profile.firstName || 'Name' }}
            {{ profile.lastName || '' }}
          </h1>

          <p class="username">
            @{{ profile.username }}
          </p>

          <p class="role">
            {{ profile.role }}
          </p>

        </div>

        <div class="profile-stats-buttons">
          <button class="stats-btn" @click="openModal('followers')">
            Followers
          </button>
          <button class="stats-btn" @click="openModal('following')">
            Following
          </button>
        </div>

      </div>

      <div class="profile-section">

        <h2>Motto</h2>

        <p class="motto" v-if="profile.motto">
          "{{ profile.motto }}"
        </p>

        <p v-else class="empty-text">
          No motto added yet.
        </p>

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

    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ modalTitle }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        
        <div class="modal-body">
          <div v-if="loadingUsers" class="modal-loading">Loading...</div>
          <ul v-else-if="modalUsers.length > 0" class="user-list">
            <li v-for="user in modalUsers" :key="user.userId" class="user-item">
              @{{ user.username }}
            </li>
          </ul>
          <p v-else class="empty-text text-center">No users found.</p>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { userService } from '@/services/userService'

import { followerService } from '@/services/followerService'

export default {
  data() {
    return {
      profile: null,
      defaultImage: 'https://placehold.co/200x200',
      
      isModalOpen: false,
      modalTitle: '',
      modalUsers: [],
      loadingUsers: false
    }
  },

  async created() {
    try {
      const response = await userService.getMyProfile()
      this.profile = response.data
    } catch (err) {
      console.error('Error fetching profile:', err)
    }
  },

  methods: {
    async openModal(type) {
      this.isModalOpen = true;
      this.loadingUsers = true;
      this.modalUsers = [];
      
      if (type === 'followers') {
        this.modalTitle = 'Followers';
        try {
          const response = await followerService.getFollowers(); 
          this.modalUsers = response.data; 
        } catch (err) {
          console.error('Error fetching followers:', err);
        }
      } else if (type === 'following') {
        this.modalTitle = 'Following';
        try {
          const response = await followerService.getFollowing();
          this.modalUsers = response.data;
        } catch (err) {
          console.error('Error fetching following:', err);
        }
      }
      this.loadingUsers = false;
    },

    closeModal() {
      this.isModalOpen = false;
      this.modalUsers = [];
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
  position: relative;
}

.profile-image {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #28a745;
}

.profile-main-info {
  flex-grow: 1; 
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

/* Stilovi za dugmad */
.profile-stats-buttons {
  display: flex;
  gap: 15px;
}

.stats-btn {
  background-color: transparent;
  border: 2px solid #28a745;
  color: #28a745;
  padding: 10px 20px;
  border-radius: 20px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
}

.stats-btn:hover {
  background-color: #28a745;
  color: white;
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

.text-center {
  text-align: center;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 25px;
  border-radius: 10px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.15);
  animation: fadeIn 0.3s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
  margin-bottom: 15px;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.8rem;
  cursor: pointer;
  color: #aaa;
  line-height: 1;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  max-height: 300px;
  overflow-y: auto;
}

.modal-loading {
  text-align: center;
  color: #666;
  font-style: italic;
  padding: 20px;
}

.user-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.user-item {
  padding: 12px 10px;
  border-bottom: 1px solid #f5f5f5;
  color: #333;
  font-weight: 500;
}

.user-item:last-child {
  border-bottom: none;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 700px) {

  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .profile-stats-buttons {
    margin-top: 15px;
    justify-content: center;
  }

  .profile-card {
    padding: 25px;
  }
}
</style>
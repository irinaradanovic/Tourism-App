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
          <h1>{{ profile.firstName || 'Name' }} {{ profile.lastName || '' }}</h1>
          <p class="username">@{{ profile.username }}</p>
          <p class="role">{{ profile.role }}</p>
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

    <!-- Modal za Followers/Following -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ modalTitle }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        
        <div class="modal-body">
          <div v-if="loadingUsers" class="modal-loading">Loading...</div>
          <ul v-else-if="modalUsers.length > 0" class="user-list">
            <li
              v-for="user in modalUsers"
              :key="user.userId"
              class="user-item"
              @click="goToProfile(user.userId)"
            >
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
    
    goToProfile(userId) {
      this.closeModal()
      this.$router.push(`/${userId}/profile`)
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
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.06);
}


.profile-header {
  display: flex;
  gap: 30px;
  align-items: center;
  margin-bottom: 40px;
}

.profile-image {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #28a745;
  box-shadow: 0 4px 10px rgba(40, 167, 69, 0.15);
}

.profile-main-info {
  flex-grow: 1; 
}

.profile-main-info h1 {
  margin: 0;
  color: #2d3748;
  font-size: 2.2rem;
  font-weight: 700;
}

.username {
  color: #718096;
  font-size: 1.1rem;
  margin-top: 4px;
  margin-bottom: 0;
}

.role {
  display: inline-block;
  background-color: #e6fffa;
  color: #234e52;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: bold;
  margin-top: 10px;
}


.profile-stats-buttons {
  display: flex;
  gap: 12px;
}

.stats-btn {
  background-color: transparent;
  border: 2px solid #e2e8f0;
  color: #4a5568;
  padding: 10px 20px;
  border-radius: 20px;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.stats-btn:hover {
  border-color: #28a745;
  color: #28a745;
  background-color: #f0fff4;
}


.motto {
  margin-top: 20px;
  font-style: italic;
  color: #28a745;
  font-size: 1.1rem;
  font-weight: 500;
}

.profile-section {
  margin-top: 35px;
  border-top: 1px solid #edf2f7;
  padding-top: 25px;
}

.profile-section h2 {
  color: #2d3748;
  font-size: 1.3rem;
  margin-bottom: 12px;
}

.profile-section p {
  line-height: 1.6;
  color: #4a5568;
}

.empty-text {
  color: #a0aec0;
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
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 25px;
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.15);
  animation: fadeIn 0.25s ease;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #edf2f7;
  padding-bottom: 12px;
  margin-bottom: 15px;
}

.modal-header h3 {
  margin: 0;
  color: #2d3748;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.8rem;
  cursor: pointer;
  color: #a0aec0;
  line-height: 1;
}

.close-btn:hover {
  color: #2d3748;
}

.modal-body {
  max-height: 300px;
  overflow-y: auto;
}

.modal-loading {
  text-align: center;
  color: #718096;
  font-style: italic;
  padding: 20px;
}

.user-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.user-item {
  padding: 12px 15px;
  border-bottom: 1px solid #f7fafc;
  color: #2d3748;
  font-weight: 500;
  cursor: pointer;             
  transition: all 0.2s ease;
  border-radius: 6px;
}

.user-item:hover {
  background-color: #f0fff4;   
  color: #28a745;             
  padding-left: 20px;          
}

.user-item:last-child {
  border-bottom: none;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-15px); }
  to { opacity: 1; transform: translateY(0); }
}


@media (max-width: 768px) {
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
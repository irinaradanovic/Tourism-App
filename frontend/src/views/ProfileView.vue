<template>
  <div class="page-wrapper" v-if="profile">
    <div class="detail-container">

      <!-- Profile Header -->
      <header class="profile-header">
        <div class="image-upload-wrapper">
          <img
              :src="profile.profileImage || defaultImage"
              alt="Profile image"
              class="profile-image"
          />
          <!-- Skriveni file input i overlay dugme -->
          <label class="upload-overlay">
            <span class="upload-icon">📷</span>
            <span class="upload-text">Change Photo</span>
            <input type="file" @change="handleImageUpload" accept="image/*" class="file-input" />
          </label>
        </div>

        <div class="profile-main-info">
          <h1>
            {{ profile.firstName || 'Name' }}
            {{ profile.lastName || '' }}
          </h1>
          <p class="username">@{{ profile.username }}</p>
          <span :class="['role-badge', profile.role ? profile.role.toLowerCase() : '']">
            {{ profile.role }}
          </span>
        </div>
      </header>

      <!-- Profile Info Sections -->
      <main class="profile-content">
        <section class="profile-section">
          <h3>Motto</h3>
          <p class="motto-text" v-if="profile.motto">
            “{{ profile.motto }}”
          </p>
          <p v-else class="empty-text">
            No motto added yet.
          </p>
        </section>

        <section class="profile-section">
          <h3>Biography</h3>
          <p class="bio-text" v-if="profile.biography">
            {{ profile.biography }}
          </p>
          <p v-else class="empty-text">
            No biography added yet.
          </p>
        </section>

        <section class="profile-section font-email">
          <h3>Email Address</h3>
          <div class="email-box">
            <span class="email-icon">✉️</span>
            <span class="email-text">{{ profile.email }}</span>
          </div>
        </section>
      </main>

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
    this.fetchProfile()
  },

  methods: {
    async fetchProfile() {
      try {
        const response = await userService.getMyProfile()
        this.profile = response.data
      } catch (err) {
        console.error('Error fetching profile:', err)
      }
    },

    async handleImageUpload(e) {
      const file = e.target.files[0]
      if (!file) return

      const formData = new FormData()
      formData.append('image', file)

      try {
        await userService.uploadProfileImage(formData)

        this.fetchProfile()
      } catch (err) {
        console.error('Error uploading profile image:', err)
        alert('Failed to upload image. Please try again.')
      }
    }
  }
}
</script>

<style scoped>
.page-wrapper {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding: 40px 20px;
}

.detail-container {
  max-width: 750px;
  margin: 0 auto;
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
}

.profile-header {
  display: flex;
  gap: 35px;
  align-items: center;
  border-bottom: 2px solid #f0f0f0;
  padding-bottom: 35px;
  margin-bottom: 30px;
}

/* Stilovi za sliku i skriveni upload na hover */
.image-upload-wrapper {
  position: relative;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  border: 4px solid white;
  outline: 2px solid #28a745;
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  opacity: 0;
  cursor: pointer;
  transition: opacity 0.3s ease;
}

.image-upload-wrapper:hover .upload-overlay {
  opacity: 1;
}

.upload-icon {
  font-size: 1.4rem;
  margin-bottom: 4px;
}

.upload-text {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.file-input {
  display: none;
}

/* Informacije o korisniku */
.profile-main-info h1 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 2.2rem;
}

.username {
  color: #7f8c8d;
  margin: 0 0 15px 0;
  font-size: 1.05rem;
}

.role-badge {
  display: inline-block;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
  background-color: #e9ecef;
  color: #495057;
  letter-spacing: 0.5px;
}

/* Menjanje boje bedža u zavisnosti od uloge */
.role-badge.admin {
  background-color: #f8d7da;
  color: #721c24;
}

.role-badge.guide {
  background-color: #d4edda;
  color: #155724;
}

/* Sekcije profila */
.profile-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.profile-section h3 {
  font-size: 1.2rem;
  color: #2c3e50;
  margin: 0 0 10px 0;
  border-left: 3px solid #28a745;
  padding-left: 10px;
}

.motto-text {
  font-size: 1.15rem;
  font-style: italic;
  color: #28a745;
  margin: 0;
  line-height: 1.6;
}

.bio-text {
  font-size: 1.05rem;
  color: #34495e;
  line-height: 1.7;
  margin: 0;
  white-space: pre-wrap;
}

.empty-text {
  color: #95a5a6;
  font-style: italic;
  margin: 0;
  font-size: 0.95rem;
}

/* Email boks */
.email-box {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  background: #f8f9fa;
  padding: 10px 15px;
  border-radius: 6px;
  border: 1px solid #eef0f2;
}

.email-icon {
  font-size: 1.1rem;
}

.email-text {
  font-size: 1rem;
  color: #34495e;
  font-weight: 500;
}

/* Responsive */
@media (max-width: 700px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }

  .detail-container {
    padding: 25px;
  }
}
</style>
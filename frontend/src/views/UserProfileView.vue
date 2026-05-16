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

        <div class="profile-actions">
          <div class="profile-stats-buttons">
            <button class="stats-btn" @click="openModal('followers')">Followers</button>
            <button class="stats-btn" @click="openModal('following')">Following</button>
          </div>

          <div class="follow-action" v-if="!isOwnProfile">
            <button
              v-if="isFollowing"
              class="unfollow-btn"
              @click="handleUnfollow"
              :disabled="followLoading"
            >
              {{ followLoading ? 'Loading...' : 'Unfollow' }}
            </button>
            <button
              v-else
              class="follow-btn"
              @click="handleFollow"
              :disabled="followLoading"
            >
              {{ followLoading ? 'Loading...' : 'Follow' }}
            </button>
          </div>
        </div>
      </div>

      <div class="profile-section">
        <h2>Motto</h2>
        <p class="motto" v-if="profile.motto">"{{ profile.motto }}"</p>
        <p v-else class="empty-text">No motto added yet.</p>

        <h2>Biography</h2>
        <p v-if="profile.biography">{{ profile.biography }}</p>
        <p v-else class="empty-text">No biography added yet.</p>
      </div>

      <div class="profile-section">
        <h2>Email</h2>
        <p>{{ profile.email }}</p>
      </div>
    </div>

    <div class="blogs-section" v-if="isFollowing">
      <h2 class="blogs-title">Blogs by @{{ profile.username }}</h2>

      <div v-if="loadingBlogs" class="blogs-loading">Loading blogs...</div>

      <div v-else-if="blogs.length > 0" class="blogs-grid">
        <div
          v-for="blog in blogs"
          :key="blog.id"
          class="blog-card"
          @click="$router.push(`/blogs/${blog.id}`)"
        >
          <img
            v-if="blog.images && blog.images.length > 0"
            :src="`http://localhost:8081/${blog.images[0]}`"
            class="blog-thumbnail"
            alt="Blog image"
          />
          <div class="blog-info">
            <h3>{{ blog.title }}</h3>
            <p class="blog-desc">{{ blog.description }}</p>
            <span class="blog-date">{{ formatDate(blog.created_at) }}</span>
          </div>
        </div>
      </div>

      <p v-else class="empty-text">This user has no blogs yet.</p>
    </div>

    <div v-else-if="!isOwnProfile" class="follow-to-see">
      <p>Follow this user to see their blogs.</p>
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
import { blogService } from '@/services/blogService'

export default {
  data() {
    return {
      userId: null,
      profile: null,
      defaultImage: 'https://placehold.co/200x200',
      isFollowing: false,
      followLoading: false,
      isOwnProfile: false,
      blogs: [],
      loadingBlogs: false,
      isModalOpen: false,
      modalTitle: '',
      modalUsers: [],
      loadingUsers: false
    }
  },

  async created() {
    this.userId = this.$route.params.id
    const token = localStorage.getItem('token')
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1]))
      this.isOwnProfile = String(payload.sub) === String(this.userId)
    }

    try {
      const response = await userService.getUserProfile(this.userId)
      this.profile = response.data
    } catch (err) {
      console.error('Error fetching user profile:', err)
    }

    await this.checkFollowStatus()

    if (this.isFollowing) {
      await this.loadBlogs()
    }
  },

  methods: {
    async checkFollowStatus() {
      if (this.isOwnProfile) return
      try {
        const response = await followerService.getFollowing()
        const following = response.data
        this.isFollowing = following.some(u => String(u.userId) === String(this.userId))
      } catch (err) {
        console.error('Error checking follow status:', err)
      }
    },

    async handleFollow() {
      this.followLoading = true
      try {
        await followerService.follow(this.userId)
        this.isFollowing = true
        await this.loadBlogs()
      } catch (err) {
        console.error('Error following user:', err)
      } finally {
        this.followLoading = false
      }
    },

    async handleUnfollow() {
      this.followLoading = true
      try {
        await followerService.unfollow(this.userId)
        this.isFollowing = false
        this.blogs = []
      } catch (err) {
        console.error('Error unfollowing user:', err)
      } finally {
        this.followLoading = false
      }
    },

    async loadBlogs() {
      this.loadingBlogs = true
      try {
        const response = await blogService.getBlogsByAuthor(this.userId)
        this.blogs = response.data || []
      } catch (err) {
        console.error('Error fetching blogs:', err)
        this.blogs = []
      } finally {
        this.loadingBlogs = false
      }
    },

    async openModal(type) {
      this.isModalOpen = true
      this.loadingUsers = true
      this.modalUsers = []

      try {
        let response
        if (type === 'followers') {
          this.modalTitle = 'Followers'
          response = await followerService.getUserFollowers(this.userId)
        } else {
          this.modalTitle = 'Following'
          response = await followerService.getUserFollowing(this.userId)
        }
        this.modalUsers = response.data
      } catch (err) {
        console.error('Error fetching modal users:', err)
      } finally {
        this.loadingUsers = false
      }
    },

    closeModal() {
      this.isModalOpen = false
      this.modalUsers = []
    },

    goToProfile(userId) {
      this.closeModal()
      this.$router.push(`/${userId}/profile`)
    },

    formatDate(dateStr) {
      if (!dateStr) return ''
      return new Date(dateStr).toLocaleDateString('sr-RS', {
        year: 'numeric', month: 'long', day: 'numeric'
      })
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
  box-shadow: 0 4px 10px rgba(40, 167, 69, 0.2);
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

.profile-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 15px;
}

.profile-stats-buttons {
  display: flex;
  gap: 12px;
}

.stats-btn {
  background-color: transparent;
  border: 2px solid #e2e8f0;
  color: #4a5568;
  padding: 8px 18px;
  border-radius: 20px;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.stats-btn:hover {
  border-color: #28a745;
  color: #28a745;
  background-color: #f6fedb;
}


.follow-btn {
  background-color: #28a745;
  color: white;
  border: none;
  padding: 10px 30px;
  border-radius: 25px;
  font-weight: bold;
  font-size: 1rem;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(40, 167, 69, 0.3);
  transition: all 0.2s ease;
  width: 140px;
}

.follow-btn:hover {
  background-color: #218838;
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(40, 167, 69, 0.4);
}

.follow-btn:active {
  transform: translateY(0);
}

.unfollow-btn {
  background-color: #fff;
  color: #dc3545;
  border: 2px solid #dc3545;
  padding: 8px 28px;
  border-radius: 25px;
  font-weight: bold;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 140px;
}

.unfollow-btn:hover {
  background-color: #dc3545;
  color: white;
  box-shadow: 0 4px 12px rgba(220, 53, 69, 0.3);
  transform: translateY(-2px);
}

.unfollow-btn:active {
  transform: translateY(0);
}

.follow-btn:disabled,
.unfollow-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
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
  margin-bottom: 20px;
}

.motto {
  font-style: italic;
  color: #28a745;
  font-size: 1.1rem;
  font-weight: 500;
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


.blogs-section {
  margin-top: 40px;
}

.blogs-title {
  font-size: 1.5rem;
  color: #2d3748;
  margin-bottom: 20px;
}

.blogs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.blog-card {
  background: white;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.blog-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.1);
}

.blog-thumbnail {
  width: 100%;
  height: 160px;
  object-fit: cover;
}

.blog-info {
  padding: 15px;
}

.blog-info h3 {
  margin: 0 0 8px;
  color: #2d3748;
  font-size: 1.1rem;
}

.blog-desc {
  color: #718096;
  font-size: 0.9rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.blog-date {
  font-size: 0.8rem;
  color: #a0aec0;
  margin-top: 12px;
  display: block;
}

.follow-to-see {
  margin-top: 30px;
  text-align: center;
  color: #718096;
  font-style: italic;
  padding: 40px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
}

.blogs-loading {
  color: #718096;
  font-style: italic;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .profile-actions {
    align-items: center;
    width: 100%;
  }

  .profile-stats-buttons {
    justify-content: center;
  }

  .profile-card {
    padding: 25px;
  }
}
</style>
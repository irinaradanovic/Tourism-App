<template>
  <div class="page-wrapper" v-if="blog">
    <div class="detail-container">
      <header class="blog-header">
        <h1>{{ blog.title }}</h1>
        <div class="blog-meta">
          <span class="author">👤 {{ blog.authorUsername || 'Unknown author' }}</span>
          <span class="date">📅 {{ formatDate(blog.created_at) }}</span>
        </div>
      </header>

      <article class="blog-content">
        {{ blog.description }}
      </article>
      <div class="image-gallery" v-if="blog.images && blog.images.length > 0">
        <h3>Galerija slika</h3>
        <div class="gallery-grid">
          <div v-for="img in blog.images" :key="img" class="img-container">
            <img :src="'http://localhost:8081/' + img" alt="Blog image" class="thumbnail" @click="openFullImage(img)"/>
          </div>
        </div>
      </div>

      <footer class="blog-footer">
        <div class="like-section">
          <button @click="handleLike" :class="['like-button', { 'is-liked': isLiked }]">
            <svg class="heart-icon" viewBox="0 0 24 24">
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
            </svg>
            <span class="like-count">{{ blog.likes }}</span>
          </button>
          <p class="like-text">Like this blog? Leave a like!</p>
        </div>
      </footer>
    </div>
  </div>
</template>

<script>
import { blogService } from '@/services/blogService';

export default {
  data() {
    return { 
      blog: null,
      isLiked: false 
    };
  },
  async created() {
    this.fetchBlog();
  },
  methods: {
    async fetchBlog() {
      const id = this.$route.params.id;
      const res = await blogService.getBlogById(id);
      this.blog = res.data;
    },
    async handleLike() {
      try {
        await blogService.toggleLike(this.blog.id);
        this.isLiked = !this.isLiked; 
        this.fetchBlog(); // refresh when liked
      } catch (err) {
        alert("You have to log in in order to like this blog.");
      }
    },
    formatDate(date) {
      return new Date(date).toLocaleDateString('sr-RS');
    },
    openFullImage(img) {
      window.open('http://localhost:8081/' + img, '_blank');
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
  max-width: 800px;
  margin: 0 auto;
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
}

.blog-header {
  border-bottom: 2px solid #f0f0f0;
  margin-bottom: 30px;
  padding-bottom: 20px;
}

.blog-header h1 {
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 10px;
}

.blog-meta {
  color: #7f8c8d;
  display: flex;
  gap: 20px;
  font-size: 0.9rem;
}

.blog-content {
  font-size: 1.1rem;
  line-height: 1.8;
  color: #34495e;
  white-space: pre-wrap; 
  margin-bottom: 40px;
}


.image-gallery h3 {
  margin-bottom: 15px;
  color: #2c3e50;
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); 
  gap: 20px;
  margin-bottom: 40px;
}


.thumbnail {
  width: 100%;
  height: 200px; 
  object-fit: cover;
  border-radius: 10px; 
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border: 1px solid #eee;
}

.thumbnail:hover {
  transform: scale(1.02);
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.blog-footer {
  border-top: 2px solid #f0f0f0;
  padding-top: 30px;
  text-align: center;
}

.like-button {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  background: none;
  border: 2px solid #e74c3c;
  padding: 10px 25px;
  border-radius: 30px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #e74c3c;
}

.like-button:hover {
  background-color: #fdf2f2;
  transform: scale(1.05);
}

.like-button.is-liked {
  background-color: #e74c3c;
  color: white;
}

.heart-icon {
  width: 24px;
  height: 24px;
  fill: currentColor;
}

.like-count {
  font-weight: bold;
  font-size: 1.2rem;
}

.like-text {
  margin-top: 10px;
  color: #95a5a6;
  font-size: 0.9rem;
}
</style>
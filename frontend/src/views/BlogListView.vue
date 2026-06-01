<template>
  <div class="container">
    <div class="header-section">
      <h1>All Blogs</h1>
      <router-link to="/create-blog" class="btn-main">Create New Blog</router-link>
    </div>
    
    <div class="blog-grid">
      <div v-for="blog in blogs" :key="blog.id" class="blog-card">
        <div class="blog-image-wrapper">
          <img 
            v-if="blog.images && blog.images.length > 0" 
            :src="'http://localhost:8081/' + blog.images[0]" 
            alt="Blog cover"
            class="blog-card-img"
          />
          <div v-else class="no-image">No image</div>
        </div>

        <div class="blog-content">
          <h2>{{ blog.title }}</h2>
          <p>{{ truncate(blog.description) }}</p>
          
          <div class="card-footer">
            <router-link 
              :to="{ name: 'blogDetails', params: { id: blog.id } }" 
              class="btn-view-more"
            >
              View more
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { blogService } from '@/services/blogService';

export default {
  data() { return { blogs: [] } },
  async created() {
    try {
      const res = await blogService.getAllBlogs();
      this.blogs = res.data;
    } catch (err) {
      console.error("Error while fetching blogs:", err);
    }
  },
  methods: {
    truncate(text) { 
      if (!text) return "";
      return text.length > 120 ? text.substring(0, 120) + '...' : text; 
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}


.btn-main {
  background-color: #28a745;
  color: white;
  padding: 10px 20px;
  text-decoration: none;
  border-radius: 5px;
  font-weight: bold;
  transition: background 0.3s;
}

.btn-main:hover {
  background-color: #218838;
}

.blog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.blog-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  transition: transform 0.2s, box-shadow 0.2s;
}

.blog-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
}


.blog-image-wrapper {
  width: 100%;
  height: 180px;
  background-color: #f0f0f0;
}

.blog-card-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #888;
}


.blog-content {
  padding: 15px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.blog-content h2 {
  margin: 0 0 10px 0;
  font-size: 1.4rem;
  color: #333;
}

.blog-content p {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.4;
  margin-bottom: 20px;
  flex-grow: 1;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
}

.btn-view-more {
  background-color: #28a745;
  color: white;
  padding: 8px 16px;
  text-decoration: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  transition: background 0.3s;
}

.btn-view-more:hover {
  background-color: #218838;
}
</style>
<template>
  <div class="page-wrapper">
    <div class="form-container">
      <header class="form-header">
        <h1>New Blog Post</h1>
        <p>Share Your Trips With Others</p>
      </header>

      <form @submit.prevent="submitForm" class="blog-form">
        <div class="input-group">
          <label>Title</label>
          <input 
            v-model="title" 
            type="text" 
            placeholder="Enter a Title..." 
            required 
          />
        </div>

        <div class="input-group">
          <label>Blog content</label>
          <textarea 
            v-model="description" 
            placeholder="Write a blog description" 
            rows="10" 
            required
          ></textarea>
        </div>

        <div class="input-group">
          <label>Slike</label>
          <div class="upload-area" @click="$refs.fileInput.click()">
            <input 
              type="file" 
              multiple 
              ref="fileInput" 
              @change="onFileChange" 
              style="display: none" 
            />
            <div class="upload-placeholder">
              <span>📷 Click to add images</span>
            </div>
          </div>

          <div class="preview-grid" v-if="previews.length > 0">
            <div v-for="(url, index) in previews" :key="index" class="preview-item">
              <img :src="url" class="preview-img" />
              <button type="button" class="remove-btn" @click="removeImage(index)">×</button>
            </div>
          </div>
        </div>

        <button type="submit" class="btn-submit" :disabled="loading">
          {{ loading ? 'Publishing...' : 'Publish' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { blogService } from '@/services/blogService';

export default {
  data() {
    return { 
      title: '', 
      description: '', 
      files: [], 
      previews: [],
      loading: false
    };
  },
  methods: {
    onFileChange(e) {
      const selectedFiles = Array.from(e.target.files);
      
      selectedFiles.forEach(file => {
        this.files.push(file);
        this.previews.push(URL.createObjectURL(file));
      });
    },
    removeImage(index) {
      this.files.splice(index, 1);
      this.previews.splice(index, 1);
    },
    async submitForm() {
      this.loading = true;
      try {
        const formData = new FormData();
        formData.append('title', this.title);
        formData.append('description', this.description);
        
        this.files.forEach(file => {
          formData.append('images', file);
        });

        await blogService.createBlog(formData);
      } catch (err) {
        alert("Error while creating blog.");
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>

<style scoped>
.page-wrapper {
  background-color: #f4f7f6;
  min-height: 100vh;
  padding: 40px 20px;
}

.form-container {
  max-width: 700px;
  margin: 0 auto;
  background: white;
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
}

.form-header {
  text-align: center;
  margin-bottom: 30px;
}

.form-header h1 {
  color: #2c3e50;
  margin-bottom: 5px;
}

.form-header p {
  color: #7f8c8d;
}

.blog-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-group label {
  font-weight: 600;
  color: #34495e;
}

input[type="text"], textarea {
  padding: 12px;
  border: 2px solid #edf2f7;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

input[type="text"]:focus, textarea:focus {
  outline: none;
  border-color: #28a745;
}

.upload-area {
  border: 2px dashed #cbd5e0;
  padding: 20px;
  text-align: center;
  border-radius: 8px;
  cursor: pointer;
  background-color: #f8fafc;
  transition: all 0.3s;
}

.upload-area:hover {
  background-color: #edf2f7;
  border-color: #28a745;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 10px;
  margin-top: 15px;
}

.preview-item {
  position: relative;
  height: 100px;
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 6px;
}

.remove-btn {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #e74c3c;
  color: white;
  border: none;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  cursor: pointer;
  font-weight: bold;
  line-height: 1;
}

.btn-submit {
  background-color: #28a745;
  color: white;
  padding: 15px;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s;
  margin-top: 10px;
}

.btn-submit:hover {
  background-color: #218838;
}

.btn-submit:disabled {
  background-color: #94d3a2;
  cursor: not-allowed;
}
</style>
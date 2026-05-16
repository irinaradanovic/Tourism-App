<template>
  <div class="page-wrapper">
    <div class="detail-container">
      <!-- Tour Header -->
      <header class="tour-header">
        <h2>{{ tour.title }}</h2>
        <p class="tour-description">{{ tour.description }}</p>
      </header>

      <!-- Key Points Section -->
      <section class="key-points-section">
        <h3>📍 Key Points</h3>

        <div class="key-points-list">
          <div v-for="kp in tour.keyPoints" :key="kp.name" class="kp-card">
            <div class="kp-image-wrapper" v-if="kp.image">
              <img :src="'http://localhost:8083/uploads/' + kp.image" alt="Key point image" class="kp-thumb" />
            </div>
            <div class="kp-info">
              <h4>{{ kp.name }}</h4>
              <p>{{ kp.description }}</p>
              <div class="kp-coordinates" v-if="kp.latitude && kp.longitude">
                🌐 {{ kp.latitude }}, {{ kp.longitude }}
              </div>
            </div>
          </div>
        </div>
      </section>

      <hr class="section-divider" />

      <!-- Add Key Point Section -->
      <section class="add-kp-section">
        <h3>➕ Add New Key Point</h3>

        <div class="form-card">
          <div class="form-grid">
            <div class="form-group full-width">
              <input v-model="kp.name" placeholder="Key Point Name" />
            </div>
            <div class="form-group full-width">
              <input v-model="kp.description" placeholder="Short Description" />
            </div>
            <div class="form-group">
              <input v-model="kp.latitude" placeholder="Latitude (e.g. 45.25)" />
            </div>
            <div class="form-group">
              <input v-model="kp.longitude" placeholder="Longitude (e.g. 19.84)" />
            </div>
            <div class="form-group full-width file-group">
              <label class="file-label">Choose Image File</label>
              <input type="file" @change="handleFile" class="file-input" />
            </div>
          </div>

          <button @click="addKeyPoint" class="btn-add">Add Key Point</button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import { tourService } from '@/services/tourService'

export default {
  data() {
    return {
      tour: {},
      kp: {},
      file: null
    }
  },
  async created() {
    const res = await tourService.getTourById(this.$route.params.id)
    this.tour = res.data
  },
  methods: {
    handleFile(e) {
      this.file = e.target.files[0]
    },
    async addKeyPoint() {
      const formData = new FormData()
      formData.append(
          "data",
          new Blob([JSON.stringify(this.kp)], { type: "application/json" })
      )
      formData.append("image", this.file)
      await tourService.addKeyPoint(this.tour.id, formData)
      location.reload()
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

.tour-header {
  border-bottom: 2px solid #f0f0f0;
  padding-bottom: 25px;
  margin-bottom: 30px;
}

.tour-header h2 {
  font-size: 2.4rem;
  color: #2c3e50;
  margin: 0 0 15px 0;
}

.tour-description {
  font-size: 1.1rem;
  line-height: 1.7;
  color: #4a5568;
  margin: 0;
}

h3 {
  font-size: 1.4rem;
  color: #2c3e50;
  margin-top: 0;
  margin-bottom: 20px;
}

.key-points-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.kp-card {
  display: flex;
  gap: 20px;
  background: #fdfdfd;
  border: 1px solid #f0f0f0;
  padding: 20px;
  border-radius: 8px;
  align-items: center;
}

.kp-image-wrapper {
  flex-shrink: 0;
}

.kp-thumb {
  width: 140px;
  height: 100px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid #eee;
}

.kp-info h4 {
  margin: 0 0 8px 0;
  font-size: 1.15rem;
  color: #2c3e50;
}

.kp-info p {
  margin: 0 0 8px 0;
  color: #555;
  font-size: 0.95rem;
  line-height: 1.5;
}

.kp-coordinates {
  font-size: 0.85rem;
  color: #888;
}

.section-divider {
  border: 0;
  height: 2px;
  background: #f0f0f0;
  margin: 40px 0;
}

.form-card {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 8px;
  border: 1px solid #eef0f2;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-bottom: 20px;
}

.full-width {
  grid-column: span 2;
}

.form-grid input:not([type="file"]) {
  width: 100%;
  box-sizing: border-box;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.95rem;
}

.form-grid input:focus {
  outline: none;
  border-color: #28a745;
  box-shadow: 0 0 0 3px rgba(40, 167, 69, 0.1);
}

.file-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #555;
}

.file-input {
  font-size: 0.9rem;
}

.btn-add {
  background: #28a745;
  color: white;
  border: none;
  padding: 12px 25px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s;
}

.btn-add:hover {
  background: #218838;
}
</style>
<template>
  <div class="page-wrapper">
    <div class="detail-container">
      <header class="form-header">
        <h2>Create Tour</h2>
        <p class="subtitle">Fill in the details below to create a new exciting tour.</p>
      </header>

      <form @submit.prevent="createTour" class="tour-form">
        <div class="form-group">
          <label>Tour Title</label>
          <input v-model="tour.title" placeholder="Enter an catchy title for the tour" />
        </div>

        <div class="form-group">
          <label>Description</label>
          <textarea v-model="tour.description" placeholder="Describe the journey, destinations, and details..."></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Difficulty Level</label>
            <select v-model="tour.difficulty">
              <option>EASY</option>
              <option>MEDIUM</option>
              <option>HARD</option>
            </select>
          </div>

          <div class="form-group">
            <label>Tags</label>
            <input v-model="tagsInput" placeholder="nature, hiking, summer (comma separated)" />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Transport</label>
            <select v-model="duration.transportType">
              <option>WALKING</option>
              <option>BICYCLE</option>
              <option>CAR</option>
            </select>
          </div>

          <div class="form-group">
            <label>Duration (minutes)</label>
            <input v-model.number="duration.minutes" type="number" min="1" placeholder="120" />
          </div>
        </div>

        <button type="submit" class="btn-submit">Create Tour</button>
      </form>

      <p v-if="message" :class="['message-text', { 'error-msg': message.includes('Error') }]">
        {{ message }}
      </p>
    </div>
  </div>
</template>

<script>
import {tourService} from '@/services/tourService'

export default {
  data() {
    return {
      tour: {
        title: '',
        description: '',
        difficulty: 'EASY'
      },
      duration: {
        transportType: 'WALKING',
        minutes: null
      },
      tagsInput: '',
      message: ''
    }
  },
  methods: {
    async createTour() {
      try {
        const payload = {
          ...this.tour,
          tags: this.tagsInput.split(',').map(t => t.trim()).filter(Boolean),
          durations: this.duration.minutes ? [this.duration] : []
        }
        await tourService.createTour(payload)
        this.message = "Tour created!"
        this.$router.push('/my-tours')
      } catch (e) {
        console.error(e)
        this.message = "Error creating tour"
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
  max-width: 650px;
  margin: 0 auto;
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.form-header {
  border-bottom: 2px solid #f0f0f0;
  margin-bottom: 30px;
  padding-bottom: 15px;
}

.form-header h2 {
  font-size: 2rem;
  color: #2c3e50;
  margin: 0 0 5px 0;
}

.subtitle {
  color: #7f8c8d;
  font-size: 0.95rem;
  margin: 0;
}

.tour-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

label {
  font-weight: 600;
  color: #34495e;
  font-size: 0.9rem;
}

input, textarea, select {
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  font-family: inherit;
  color: #333;
  transition: border-color 0.3s, box-shadow 0.3s;
  background-color: #fff;
}

input:focus, textarea:focus, select:focus {
  outline: none;
  border-color: #28a745;
  box-shadow: 0 0 0 3px rgba(40, 167, 69, 0.12);
}

textarea {
  min-height: 120px;
  resize: vertical;
}

.btn-submit {
  background: #28a745;
  color: white;
  border: none;
  padding: 14px;
  border-radius: 6px;
  font-size: 1.05rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s, transform 0.2s;
  margin-top: 10px;
}

.btn-submit:hover {
  background: #218838;
}

.btn-submit:active {
  transform: scale(0.98);
}

.message-text {
  margin-top: 20px;
  padding: 12px;
  background-color: #d4edda;
  color: #155724;
  border-radius: 6px;
  text-align: center;
  font-weight: 500;
}

.message-text.error-msg {
  background-color: #f8d7da;
  color: #721c24;
}
</style>

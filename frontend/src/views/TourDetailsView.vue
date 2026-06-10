<template>
  <!-- v-if="tour.id" osigurava da se stranica renderuje TEK kada backend vrati podatke -->
  <div class="page-wrapper" v-if="tour && tour.id">
    <div class="detail-container">
      <!-- Tour Header -->
      <header class="tour-header">
        <h2>{{ tour.title }}</h2>
        <p class="tour-description">{{ tour.description }}</p>
        <div class="tour-meta">
          <span>Status: {{ tour.status }}</span>
          <span>Distance: {{ tour.distanceKm || 0 }} km</span>
          <span v-for="duration in tour.durations || []" :key="duration.transportType">
            {{ duration.transportType }}: {{ duration.minutes }} min
          </span>
        </div>
        <div class="lifecycle-actions" v-if="canManageKeyPoints">
          <button v-if="tour.status === 'DRAFT'" class="btn-add" @click="publishTour">Publish</button>
          <button v-if="tour.status === 'PUBLISHED'" class="btn-delete" @click="archiveTour">Archive</button>
          <button v-if="tour.status === 'ARCHIVED'" class="btn-add" @click="reactivateTour">Reactivate</button>
        </div>
      </header>

      <!-- Key Points Section -->
      <section class="key-points-section">
        <h3>📍 Key Points</h3>

        <div class="key-points-list" v-if="tour.keyPoints && tour.keyPoints.length > 0">
          <div v-for="(kp,index) in tour.keyPoints" :key="kp.name + '-' + index" class="kp-card">
            <div class="kp-image-wrapper" v-if="kp.image">
              <img :src="'http://localhost:8083/uploads/' + kp.image" alt="Key point image" class="kp-thumb" />
            </div>
            <div class="kp-info">
              <h4>{{ kp.name }}</h4>
              <p>{{ kp.description }}</p>
              <div class="kp-coordinates" v-if="kp.latitude && kp.longitude">
                🌐 {{ kp.latitude }}, {{ kp.longitude }}
              </div>
              <div class="kp-actions" v-if="canManageKeyPoints">
                <button class="btn-edit" @click="startEditKeyPoint(index)">Edit</button>
                <button class="btn-delete" @click="deleteKeyPoint(index)">Delete</button>
              </div>
            </div>
          </div>
        </div>
        <p v-else class="empty-text">No key points added to this tour yet.</p>
      </section>

      <hr class="section-divider" />

      <!-- Add Key Point Section -->
      <section class="add-kp-section">
        <h3>➕ Add New Key Point</h3>

        <div class="form-card">
          <!-- INTERAKTIVNA MAPA -->
          <div class="map-wrapper">
            <label class="file-label">Click on the map to pin the location:</label>
            <div id="map" class="leaflet-map-container"></div>
          </div>

          <div class="form-grid">
            <div class="form-group full-width">
              <input v-model="kp.name" placeholder="Key Point Name (e.g. Museum, Square)" />
            </div>
            <div class="form-group full-width">
              <input v-model="kp.description" placeholder="Short Description" />
            </div>
            <div class="form-group">
              <input v-model="kp.latitude" placeholder="Latitude (Click map)" readonly class="readonly-input" />
            </div>
            <div class="form-group">
              <input v-model="kp.longitude" placeholder="Longitude (Click map)" readonly class="readonly-input" />
            </div>
            <div class="form-group full-width file-group">
              <label class="file-label">Choose Image File</label>
              <input type="file" @change="handleFile" class="file-input" />
            </div>
          </div>
          <button
            @click="editingIndex === null ? addKeyPoint() : updateKeyPoint()"
            class="btn-add"
            :disabled="!kp.latitude"
          >
            {{ editingIndex === null ? 'Add Key Point' : 'Save Changes' }}
          </button>
          <button v-if="editingIndex !== null" @click="cancelEdit" class="btn-cancel">Cancel</button>
        </div>
      </section>

      <hr class="section-divider" />

      <section class="add-kp-section" v-if="canManageKeyPoints">
        <h3>Tour Durations</h3>
        <div class="form-card">
          <div class="form-grid">
            <div class="form-group">
              <select v-model="duration.transportType">
                <option>WALKING</option>
                <option>BICYCLE</option>
                <option>CAR</option>
              </select>
            </div>
            <div class="form-group">
              <input v-model.number="duration.minutes" type="number" min="1" placeholder="Minutes" />
            </div>
          </div>
          <button class="btn-add" @click="saveDuration">Save Duration</button>
        </div>
      </section>
    </div>
  </div>
  <div v-else class="loading-box">
    <p>Fetching tour data...</p>
  </div>
</template>

<script>
import { tourService } from '@/services/tourService'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

export default {
  data() {
    return {
      tour: {},
      kp: {
        name: '',
        description: '',
        latitude: null,
        longitude: null,
      },
      file: null,
      map: null,
      marker: null,
      polyline: null,
      keyPointMarkers: [],
      editingIndex:null,
      currentUserId: null,
      duration: {
        transportType: 'WALKING',
        minutes: null
      }
    }
  },
  async created() {
    const tourId = this.$route.params.id
    this.currentUserId = this.getCurrentUserId()
    try {
      const res = await tourService.getTourById(tourId)
      this.tour = res.data

      // Pokrećemo mapu čim se HTML elementi generišu
      this.$nextTick(() => {
        this.initMap()
      })
    } catch (err) {
      console.error("Greška pri dobavljanju ture:", err)
    }
  },
  computed: {
    canManageKeyPoints() {
      return (
        this.currentUserId !== null &&
        this.tour &&
        String(this.tour.authorId) === String(this.currentUserId)
      )
    }
  },
  methods: {
     initMap() {
      this.map = L.map('map').setView([45.25, 19.84], 13)

      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors'
      }).addTo(this.map)

      this.renderKeyPoints()

      // Klik na mapu postavlja marker i popunjava input polja
      this.map.on('click', (e) => {
        const { lat, lng } = e.latlng
        this.kp.latitude = lat.toFixed(6)
        this.kp.longitude = lng.toFixed(6)

        if (this.marker) {
          this.marker.setLatLng(e.latlng)
        } else {
          this.marker = L.marker(e.latlng).addTo(this.map)
        }
      })
    },
    handleFile(e) {
      this.file = e.target.files[0]
    },
    async addKeyPoint() {
      const tourId = this.tour.id || this.$route.params.id

      if (!this.kp.latitude || !this.kp.longitude) {
        alert("Please select a location on the map first.")
        return
      }

      const formData = new FormData()
      formData.append(
          "data",
          new Blob([JSON.stringify(this.kp)], {type: "application/json"})
      )
      formData.append("image", this.file)

      try {
        await tourService.addKeyPoint(tourId, formData)
        location.reload()
      } catch (err) {
        console.error("Error sending key point:", err)
        alert("An error occurred while saving on the backend.")
      }
    },
    renderKeyPoints() {
      if (!this.map) return

      // Ocisti prethodne markere i liniju
      this.keyPointMarkers.forEach(m => this.map.removeLayer(m))
      this.keyPointMarkers = []
      if (this.polyline) {
        this.map.removeLayer(this.polyline)
        this.polyline = null
      }

      const latLngs = []

      if (this.tour.keyPoints && this.tour.keyPoints.length > 0) {
        this.tour.keyPoints.forEach(point => {
          if (point.latitude && point.longitude) {
            const latLng = [point.latitude, point.longitude]
            latLngs.push(latLng)

            const marker = L.marker(latLng)
              .addTo(this.map)
              .bindPopup(`<b>${point.name}</b><br>${point.description}`)

            this.keyPointMarkers.push(marker)
          }
        })
      }

      if (latLngs.length >= 2) {
        this.polyline = L.polyline(latLngs, {
          color: '#2c7be5',
          weight: 4,
          opacity: 0.85
        }).addTo(this.map)

        this.map.fitBounds(this.polyline.getBounds(), { padding: [20, 20] })
      }
    },
    getCurrentUserId() {
      const token = localStorage.getItem('token')
      if (!token) return null
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        return payload.sub || payload.userId || payload.id || null
      } catch (_) {
        return null
      }
    },
    startEditKeyPoint(index) {
      const point = this.tour.keyPoints[index]
      this.kp = {
        name: point.name,
        description: point.description,
        latitude: point.latitude,
        longitude: point.longitude
      }
      this.editingIndex = index
    },
    async updateKeyPoint() {
      const tourId = this.tour.id || this.$route.params.id
      if (!this.kp.latitude || !this.kp.longitude) {
        alert("Please select a location on the map first.")
        return
      }

      const formData = new FormData()
      formData.append("data", new Blob([JSON.stringify(this.kp)], { type: "application/json" }))
      if (this.file) {
        formData.append("image", this.file)
      }

      try {
        await tourService.updateKeyPoint(tourId, this.editingIndex, formData)
        await this.fetchTour()
        this.resetForm()
      } catch (err) {
        console.error("Error updating key point:", err)
        alert("An error occurred while updating.")
      }
    },
    cancelEdit() {
      this.resetForm()
    },
    resetForm() {
      this.kp = { name: '', description: '', latitude: null, longitude: null }
      this.file = null
      this.editingIndex = null
    },
    async deleteKeyPoint(index) {
      const tourId = this.tour.id || this.$route.params.id
      if (!confirm('Delete this key point?')) return

      try {
        await tourService.deleteKeyPoint(tourId, index)
        await this.fetchTour()
        if (this.editingIndex === index) {
          this.resetForm()
        }
      } catch (err) {
        console.error("Error deleting key point:", err)
        alert("An error occurred while deleting.")
      }
    },
    async fetchTour() {
      const tourId = this.$route.params.id
      const res = await tourService.getTourById(tourId)
      this.tour = res.data
      if (this.map) this.renderKeyPoints()
    },
    async saveDuration() {
      if (!this.duration.minutes || this.duration.minutes <= 0) {
        alert('Duration must be greater than zero.')
        return
      }
      const existing = (this.tour.durations || []).filter(d => d.transportType !== this.duration.transportType)
      await tourService.updateDurations(this.tour.id, [...existing, { ...this.duration }])
      await this.fetchTour()
    },
    async publishTour() {
      try {
        await tourService.publishTour(this.tour.id)
        alert('Publish saga started. Refresh in a moment to see the final status.')
        await this.fetchTour()
      } catch (err) {
        alert(err.response?.data || 'Publish failed.')
      }
    },
    async archiveTour() {
      try {
        await tourService.archiveTour(this.tour.id)
        alert('Archive saga started. Refresh in a moment to see the final status.')
        await this.fetchTour()
      } catch (err) {
        alert(err.response?.data || 'Archive failed.')
      }
    },
    async reactivateTour() {
      try {
        await tourService.reactivateTour(this.tour.id)
        await this.fetchTour()
      } catch (err) {
        alert(err.response?.data || 'Reactivation failed.')
      }
    },
    
  }
}
</script>

<style scoped>
.page-wrapper {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding: 40px 20px;
}

.loading-box {
  text-align: center;
  padding: 50px;
  font-size: 1.2rem;
  color: #666;
}

.detail-container {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
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

.tour-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 16px;
  color: #4a5568;
  font-weight: 600;
}

.tour-meta span {
  background: #f4f7f9;
  border: 1px solid #e5eaef;
  border-radius: 6px;
  padding: 6px 10px;
}

.lifecycle-actions {
  display: flex;
  gap: 10px;
  margin-top: 18px;
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

.empty-text {
  color: #a0aec0;
  font-style: italic;
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

/* STILOVI ZA MAPU */
.map-wrapper {
  margin-bottom: 20px;
}

.leaflet-map-container {
  height: 320px;
  width: 100%;
  border-radius: 8px;
  border: 1px solid #ddd;
  margin-top: 8px;
  z-index: 1;
}

.readonly-input {
  background-color: #e9ecef;
  color: #495057;
  cursor: not-allowed;
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

.btn-add:disabled {
  background: #cbd5e0;
  cursor: not-allowed;
}

.kp-actions {
  margin-top: 10px;
  display: flex;
  gap: 8px;
}

.btn-edit,
.btn-delete,
.btn-cancel {
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  font-size: 0.85rem;
}

.btn-edit {
  background: #e8f0fe;
  color: #1a73e8;
}

.btn-delete {
  background: #fdecea;
  color: #d93025;
}

.btn-cancel {
  margin-left: 10px;
  background: #e2e8f0;
  color: #4a5568;
}
</style>

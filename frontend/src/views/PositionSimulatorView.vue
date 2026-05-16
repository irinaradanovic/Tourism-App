<template>
  <div class="simulator-wrapper">
    <h2>Position Simulator</h2>
    <p class="info-text">Klikni na mapu da postaviš svoju trenutnu poziciju.</p>

    <div v-if="position" class="position-badge">
      Tvoja pozicija: <b>{{ position.lat.toFixed(5) }}, {{ position.lon.toFixed(5) }}</b>
      <span class="updated-at">Ažurirano: {{ formatTime(position.updatedAt) }}</span>
    </div>
    <div v-else class="position-badge empty">
      Pozicija još nije postavljena. Klikni na mapu.
    </div>

    <div id="simulator-map"></div>
  </div>
</template>

<script>
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { tourService } from '@/services/tourService'

delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png'
})

export default {
  name: 'PositionSimulatorView',
  data() {
    return {
      map: null,
      marker: null,
      position: null
    }
  },
  async mounted() {
    this.map = L.map('simulator-map').setView([44.8178, 20.4568], 13)

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors'
    }).addTo(this.map)

    try {
      const res = await tourService.getMyPosition()
      this.position = res.data
      this.placeMarker(res.data.lat, res.data.lon)
      this.map.setView([res.data.lat, res.data.lon], 13)
    } catch (_) {
    }

    this.map.on('click', (e) => this.handleMapClick(e.latlng))
  },
  methods: {
    placeMarker(lat, lon) {
      if (this.marker) {
        this.map.removeLayer(this.marker)
      }
      this.marker = L.marker([lat, lon])
        .addTo(this.map)
        .bindPopup('Tvoja pozicija')
        .openPopup()
    },
    async handleMapClick({ lat, lng }) {
      try {
        const res = await tourService.savePosition(lat, lng)
        this.position = res.data
        this.placeMarker(lat, lng)
      } catch (err) {
        alert('Greška pri čuvanju pozicije.')
        console.error(err)
      }
    },
    formatTime(dateStr) {
      if (!dateStr) return ''
      return new Date(dateStr).toLocaleString('sr-RS')
    }
  }
}
</script>

<style scoped>
.simulator-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

h2 {
  font-size: 1.8rem;
  color: #2c3e50;
  margin-bottom: 8px;
}

.info-text {
  color: #666;
  margin-bottom: 16px;
}

.position-badge {
  background: #eafaf1;
  border: 1px solid #28a745;
  border-radius: 8px;
  padding: 10px 16px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
  color: #2c3e50;
}

.position-badge.empty {
  background: #f8f9fa;
  border-color: #ccc;
  color: #888;
}

.updated-at {
  font-size: 0.85rem;
  color: #888;
  margin-left: auto;
}

#simulator-map {
  height: 500px;
  border-radius: 12px;
  border: 1px solid #ddd;
}
</style>
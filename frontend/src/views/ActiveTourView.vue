<template>
  <div class="container">
    <h1>Active Tour</h1>

    <div v-if="!execution">
      <p>Loading session...</p>
    </div>

    <div v-else>
      <div class="info-card">
        <p><strong>Status:</strong> {{ execution.status }}</p>
        <p><strong>Started:</strong> {{ execution.startedAt }}</p>
        <p><strong>Last Activity:</strong> {{ execution.lastActivity }}</p>
        <p v-if="currentLat"><strong>Your position:</strong> {{ currentLat.toFixed(5) }}, {{ currentLon.toFixed(5) }}</p>
      </div>

      <div id="active-tour-map"></div>

      <div v-if="proximityAlert" class="proximity-alert">
        📍 {{ proximityAlert }}
      </div>

      <div class="keypoints" v-if="tour">
        <h3>Key Points</h3>
        <div v-for="(kp, index) in tour.keyPoints" :key="index" class="kp-item">
          <span>{{ kp.name }}</span>
          <span v-if="isCompleted(index)" class="completed">✓ Reached</span>
          <span v-else class="pending">Not reached</span>
        </div>
      </div>

      <p class="proximity-note">Proximity check runs automatically every 5 seconds.</p>

      <div class="actions">
        <button @click="handleAbandon" class="btn-danger">Abandon Tour</button>
        <button @click="handleComplete" class="btn-success">Complete Tour</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png'
});

const API = 'http://localhost:80';
const PROXIMITY_RADIUS_METERS = 200;

const getAuthHeader = () => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
});

function haversineDistance(lat1, lon1, lat2, lon2) {
  const R = 6371000;
  const toRad = (deg) => (deg * Math.PI) / 180;
  const dLat = toRad(lat2 - lat1);
  const dLon = toRad(lon2 - lon1);
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) ** 2;
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

export default {
  name: 'ActiveTourView',

  data() {
    return {
      execution: null,
      tour: null,
      pollIntervalId: null,
      currentLat: null,
      currentLon: null,
      map: null,
      playerMarker: null,
      keyPointMarkers: [],
      proximityAlert: null
    };
  },

  async created() {
    const executionId = this.$route.params.executionId;
    await this.loadData(executionId);
  },

  beforeUnmount() {
    if (this.pollIntervalId) clearInterval(this.pollIntervalId);
    if (this.map) {
      this.map.remove();
      this.map = null;
    }
  },

  methods: {
    async loadData(executionId) {
      try {
        const execRes = await axios.get(`${API}/api/executions/my`, { headers: getAuthHeader() });
        this.execution = execRes.data.find((e) => e.id === executionId);

        if (this.execution) {
          const tourRes = await axios.get(`${API}/api/tours/${this.execution.tourId}`, {
            headers: getAuthHeader()
          });
          this.tour = tourRes.data;
        }

        if (this.execution && this.execution.status === 'ACTIVE') {
          await this.fetchPositionAndCheck();
          this.startPositionPolling();
        }
      } catch (err) {
        console.error('Failed to load execution', err);
      }
    },

    initMap() {
      if (this.map) return;

      const mapEl = document.getElementById('active-tour-map');
      if (!mapEl) return;

      this.map = L.map('active-tour-map').setView([45.2671, 19.8335], 14);

      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors'
      }).addTo(this.map);

      if (this.tour && this.tour.keyPoints) {
        this.renderKeyPointMarkers();
      }

      if (this.currentLat !== null) {
        this.updatePlayerMarker(this.currentLat, this.currentLon);
      }
    },

    renderKeyPointMarkers() {
      if (!this.map) return;

      this.keyPointMarkers.forEach((m) => this.map.removeLayer(m));
      this.keyPointMarkers = [];

      if (!this.tour || !this.tour.keyPoints) return;

      const pendingIcon = L.divIcon({
        className: '',
        html: `<div style="background:#e74c3c;border:3px solid #fff;border-radius:50%;width:22px;height:22px;box-shadow:0 2px 6px rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;font-size:12px;color:#fff;">📍</div>`,
        iconSize: [26, 26],
        iconAnchor: [13, 13]
      });

      this.tour.keyPoints.forEach((kp, index) => {
        if (kp.latitude == null || kp.longitude == null) return;

        const completed = this.isCompleted(index);
        const icon = completed ? this.completedIcon() : pendingIcon;

        const marker = L.marker([kp.latitude, kp.longitude], { icon })
          .addTo(this.map)
          .bindPopup(
            `<b>${kp.name}</b><br>${kp.description || ''}<br>${completed ? '✅ Reached' : '⏳ Not yet reached'}`
          );

        this.keyPointMarkers.push(marker);
      });

      if (this.keyPointMarkers.length > 0) {
        const group = L.featureGroup(this.keyPointMarkers);
        this.map.fitBounds(group.getBounds().pad(0.3));
      }
    },

    completedIcon() {
      return L.divIcon({
        className: '',
        html: `<div style="background:#27ae60;border:3px solid #fff;border-radius:50%;width:22px;height:22px;box-shadow:0 2px 6px rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;font-size:14px;">✓</div>`,
        iconSize: [26, 26],
        iconAnchor: [13, 13]
      });
    },

    updatePlayerMarker(lat, lon) {
      if (!this.map) return;

      const playerIcon = L.divIcon({
        className: '',
        html: `<div style="background:#3498db;border:3px solid #fff;border-radius:50%;width:24px;height:24px;box-shadow:0 0 0 5px rgba(52,152,219,0.3);display:flex;align-items:center;justify-content:center;font-size:14px;">🧍</div>`,
        iconSize: [28, 28],
        iconAnchor: [14, 14]
      });

      if (this.playerMarker) {
        this.playerMarker.setLatLng([lat, lon]);
      } else {
        this.playerMarker = L.marker([lat, lon], { icon: playerIcon })
          .addTo(this.map)
          .bindPopup('Your position');
      }
    },

    startPositionPolling() {
      this.pollIntervalId = setInterval(async () => {
        await this.fetchPositionAndCheck();
      }, 5000);
    },

    async fetchPositionAndCheck() {
      try {
        const posRes = await axios.get(`${API}/api/position`, { headers: getAuthHeader() });
        const { lat, lon } = posRes.data;

        const positionChanged = lat !== this.currentLat || lon !== this.currentLon;

        this.currentLat = lat;
        this.currentLon = lon;

        this.updatePlayerMarker(lat, lon);

        if (this.tour && this.tour.keyPoints) {
          this.checkLocalProximity(lat, lon);
        }

        if (positionChanged) {
          await this.doServerProximityCheck(lat, lon);
        }
        if (this.execution.status === 'COMPLETED') {
          if (this.pollIntervalId) clearInterval(this.pollIntervalId);
          alert('Tour completed!');
          this.$router.push('/tours');
        }
      } catch (err) {
        console.error('Position fetch failed', err);
      }
    },

    checkLocalProximity(lat, lon) {
      for (let i = 0; i < this.tour.keyPoints.length; i++) {
        const kp = this.tour.keyPoints[i];
        if (kp.latitude == null || kp.longitude == null) continue;
        if (this.isCompleted(i)) continue;

        const dist = haversineDistance(lat, lon, kp.latitude, kp.longitude);
        if (dist <= PROXIMITY_RADIUS_METERS) {
          this.proximityAlert = `You are near "${kp.name}" (${Math.round(dist)}m away)!`;
          return;
        }
      }
      this.proximityAlert = null;
    },

    async doServerProximityCheck(lat, lon) {
      try {
        const res = await axios.post(
          `${API}/api/executions/${this.execution.id}/proximity`,
          { lat, lon },
          { headers: getAuthHeader() }
        );
        const updated = res.data;

        if (
          JSON.stringify(updated.completedKeyPoints) !==
          JSON.stringify(this.execution.completedKeyPoints)
        ) {
          this.execution = updated;
          this.$nextTick(() => this.renderKeyPointMarkers());
        } else {
          this.execution = updated;
        }
      } catch (err) {
        console.error('Server proximity check failed', err);
      }
    },

    isCompleted(index) {
      return (
        this.execution.completedKeyPoints &&
        this.execution.completedKeyPoints[index] != null
      );
    },

    async handleAbandon() {
      try {
        await axios.put(
          `${API}/api/executions/${this.execution.id}/abandon`,
          {},
          { headers: getAuthHeader() }
        );
        if (this.pollIntervalId) clearInterval(this.pollIntervalId);
        this.$router.push('/tours');
      } catch (err) {
        console.error(err);
      }
    },

    async handleComplete() {
      try {
        await axios.put(
          `${API}/api/executions/${this.execution.id}/complete`,
          {},
          { headers: getAuthHeader() }
        );
        if (this.pollIntervalId) clearInterval(this.pollIntervalId);
        alert('Tour completed!');
        this.$router.push('/tours');
      } catch (err) {
        console.error(err);
      }
    }
  },

  watch: {
    execution(newVal) {
      if (newVal) {
        this.$nextTick(() => {
          this.initMap();
        });
      }
    },

    tour(newVal) {
      if (newVal && this.map) {
        this.renderKeyPointMarkers();
      }
    }
  }
};
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.info-card {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}

#active-tour-map {
  height: 420px;
  border-radius: 12px;
  border: 1px solid #ddd;
  margin-bottom: 16px;
}

.proximity-alert {
  background: #fff3cd;
  border: 1px solid #ffc107;
  color: #856404;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 0.95rem;
}

.keypoints {
  margin-bottom: 16px;
}

.kp-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.completed {
  color: #28a745;
  font-weight: bold;
}

.pending {
  color: #999;
}

.proximity-note {
  color: #666;
  font-size: 0.85rem;
  margin: 10px 0;
}

.actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.btn-danger {
  background: #dc3545;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.btn-success {
  background: #28a745;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.btn-danger:hover { background: #c82333; }
.btn-success:hover { background: #218838; }
</style>
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
        <p v-if="currentLat"><strong>Your position:</strong> {{ currentLat }}, {{ currentLon }}</p>
      </div>

      <div class="keypoints" v-if="tour">
        <h3>Key Points</h3>
        <div v-for="(kp, index) in tour.keyPoints" :key="index" class="kp-item">
          <span>{{ kp.name }}</span>
          <span v-if="isCompleted(index)" class="completed">✓ Reached</span>
          <span v-else class="pending">Not reached</span>
        </div>
      </div>

      <p class="proximity-note">Proximity check runs automatically every 10 seconds.</p>

      <div class="actions">
        <button @click="handleAbandon" class="btn-danger">Abandon Tour</button>
        <button @click="handleComplete" class="btn-success">Complete Tour</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

const API = 'http://localhost:80';

const getAuthHeader = () => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`
});

export default {
  data() {
    return {
      execution: null,
      tour: null,
      intervalId: null,
      currentLat: null,
      currentLon: null
    };
  },
  async created() {
    const executionId = this.$route.params.executionId;
    await this.loadData(executionId);
    if (this.execution && this.execution.status === 'ACTIVE') {
      this.startProximityCheck();
    }
  },
  beforeUnmount() {
    if (this.intervalId) clearInterval(this.intervalId);
  },
  methods: {
    async loadData(executionId) {
      try {
        const execRes = await axios.get(`${API}/api/executions/my`, { headers: getAuthHeader() });
        this.execution = execRes.data.find(e => e.id === executionId);
        if (this.execution) {
          const tourRes = await axios.get(`${API}/api/tours/${this.execution.tourId}`, { headers: getAuthHeader() });
          this.tour = tourRes.data;
        }
      } catch (err) {
        console.error('Failed to load execution', err);
      }
    },
    startProximityCheck() {
      this.intervalId = setInterval(async () => {
        await this.doProximityCheck();
      }, 10000);
    },
    async doProximityCheck() {
      try {
        // Simulator integracija: prvo pitaj position simulator za trenutnu lokaciju
        const posRes = await axios.get(`${API}/api/position`, { headers: getAuthHeader() });
        const { lat, lon } = posRes.data;
        this.currentLat = lat;
        this.currentLon = lon;

        // Zatim posalji koordinate na proximity check (ide gRPC putem kroz gateway)
        const res = await axios.post(
          `${API}/api/executions/${this.execution.id}/proximity`,
          { lat, lon },
          { headers: getAuthHeader() }
        );
        this.execution = res.data;
      } catch (err) {
        console.error('Proximity check failed', err);
      }
    },
    isCompleted(index) {
      return this.execution.completedKeyPoints && this.execution.completedKeyPoints[index] != null;
    },
    async handleAbandon() {
      try {
        await axios.put(`${API}/api/executions/${this.execution.id}/abandon`, {}, { headers: getAuthHeader() });
        if (this.intervalId) clearInterval(this.intervalId);
        this.$router.push('/tours');
      } catch (err) {
        console.error(err);
      }
    },
    async handleComplete() {
      try {
        await axios.put(`${API}/api/executions/${this.execution.id}/complete`, {}, { headers: getAuthHeader() });
        if (this.intervalId) clearInterval(this.intervalId);
        alert('Tour completed!');
        this.$router.push('/tours');
      } catch (err) {
        console.error(err);
      }
    }
  }
};
</script>

<style scoped>
.container { max-width: 700px; margin: 0 auto; padding: 20px; }
.info-card { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
.kp-item { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #eee; }
.completed { color: #28a745; font-weight: bold; }
.pending { color: #999; }
.proximity-note { color: #666; font-size: 0.9rem; margin: 15px 0; }
.actions { display: flex; gap: 10px; margin-top: 20px; }
.btn-danger { background: #dc3545; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; }
.btn-success { background: #28a745; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; }
</style>
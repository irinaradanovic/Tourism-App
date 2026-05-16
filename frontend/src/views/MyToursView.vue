<template>
  <div class="page-wrapper">
    <div class="list-container">
      <header class="list-header">
        <h2>My Tours</h2>
        <p class="subtitle">Manage and review the tours you have created.</p>
      </header>

      <div class="tours-grid" v-if="tours.length > 0">
        <div v-for="t in tours" :key="t.id" class="tour-card">
          <div class="card-body">
            <h3>{{ t.title }}</h3>
            <p class="description">{{ t.description }}</p>
          </div>

          <div class="card-meta">
            <span class="badge status-badge">{{ t.status }}</span>
            <span :class="['badge', 'diff-badge', t.difficulty.toLowerCase()]">
              ⚡ {{ t.difficulty }}
            </span>
          </div>

          <div class="card-footer">
            <router-link :to="`/tours/${t.id}`" class="btn-open">
              Open Tour Details
            </router-link>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <p>You haven't created any tours yet.</p>
      </div>
    </div>
  </div>
</template>

<script>
import { tourService } from '@/services/tourService'

export default {
  data() {
    return {
      tours: []
    }
  },
  async created() {
    const res = await tourService.getMyTours()
    this.tours = res.data
  }
}
</script>

<style scoped>
.page-wrapper {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding: 40px 20px;
}

.list-container {
  max-width: 1000px;
  margin: 0 auto;
}

.list-header {
  margin-bottom: 30px;
}

.list-header h2 {
  font-size: 2.2rem;
  color: #2c3e50;
  margin: 0 0 5px 0;
}

.subtitle {
  color: #7f8c8d;
  font-size: 1rem;
  margin: 0;
}

.tours-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 25px;
}

.tour-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
  display: flex;
  flex-direction: column;
  transition: transform 0.3s, box-shadow 0.3s;
  overflow: hidden;
  border: 1px solid #eee;
}

.tour-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.1);
}

.card-body {
  padding: 25px;
  flex-grow: 1;
}

.card-body h3 {
  margin: 0 0 12px 0;
  color: #2c3e50;
  font-size: 1.3rem;
  line-height: 1.4;
}

.description {
  color: #555;
  font-size: 0.95rem;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  padding: 0 25px 15px 25px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-badge {
  background-color: #e9ecef;
  color: #495057;
}

.diff-badge.easy { background-color: #d4edda; color: #155724; }
.diff-badge.medium { background-color: #fff3cd; color: #856404; }
.diff-badge.hard { background-color: #f8d7da; color: #721c24; }

.card-footer {
  padding: 15px 25px 25px 25px;
  border-top: 1px solid #f8f9fa;
}

.btn-open {
  display: block;
  text-align: center;
  background: transparent;
  border: 2px solid #28a745;
  color: #28a745;
  padding: 10px;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
}

.btn-open:hover {
  background: #28a745;
  color: white;
}

.empty-state {
  text-align: center;
  padding: 40px;
  background: white;
  border-radius: 12px;
  color: #7f8c8d;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}
</style>
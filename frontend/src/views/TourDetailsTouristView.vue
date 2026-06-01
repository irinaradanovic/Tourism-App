<template>
  <div class="container" v-if="tour">
    <div class="tour-header">
        <h1>{{ tour.title }}</h1>
        <span class="difficulty-tag" :class="tour.difficulty?.toLowerCase()">
            {{ tour.difficulty }}
        </span>
    </div>

    <div class="tour-layout">
      <div class="main-info">
        <section class="info-section">
          <h3>Description</h3>
          <p class="description-text">{{ tour.description }}</p>
        </section>

        <section class="info-section">
          <h3>Key Points (Route)</h3>
          
          <div v-if="isPurchased" class="purchased-access">
            <p class="access-success">✓ You have purchased this tour! All key points are unlocked.</p>
            <div class="keypoints-list">
              <div v-for="(kp, index) in tour.keyPoints" :key="index" class="kp-card">
                <span class="kp-number">{{ index + 1 }}</span>
                <div>
                  <h4>{{ kp.name }}</h4>
                  <p>{{ kp.description }}</p>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="restricted-access">
            <p class="access-locked">🔒 Detailed information is available only for purchased tours.</p>
            <p v-if="tour.keyPoints && tour.keyPoints.length > 0">
              <strong>Starting point:</strong> {{ tour.keyPoints[0].name }}
            </p>
          </div>
        </section>

        <section class="info-section">
          <h3>Reviews ({{ tour.reviews ? tour.reviews.length : 0 }})</h3>
          <div v-if="!tour.reviews || tour.reviews.length === 0" class="no-reviews">
            Be the first to leave a review after you experience the tour!
          </div>
          <div v-else class="reviews-list">
            <div v-for="review in tour.reviews" :key="review.id" class="review-card">
              <div class="review-meta">
                <span class="review-rating">Rating: {{ review.rating }}/5</span>
                <span class="review-user">User #{{ review.touristId }}</span>
              </div>
              <p>{{ review.comment }}</p>
            </div>
          </div>
        </section>
      </div>

      <div class="sidebar-info">
        <div class="sticky-card">
          <div class="price-box">
            <span class="label">Price:</span>
            <span class="value">{{ tour.price }} RSD</span>
          </div>

          <div class="status-box" v-if="isPurchased">
            <span class="badge-purchased">Purchased</span>
            <button @click="handleStartTour" class="btn-start">
              ▶ Start Tour
            </button>
          </div>
          <div class="action-box" v-else>
            <button v-if="isInCart(tour.id)" @click="handleRemoveFromCart(tour.id)" class="btn-remove">
              Remove from Cart
            </button>
            <button v-else @click="handleAddToCart" class="btn-add">
              Add to Cart
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { tourPublicService, purchaseService } from '@/services/purchaseService';
import axios from 'axios';

const API = 'http://localhost:80';
const getAuthHeader = () => ({ Authorization: `Bearer ${localStorage.getItem('token')}` });

export default {
  name: 'TourDetailsTouristView',
  data() {
    return {
      tour: null,
      cart: { items: [] },
      isPurchased: false
    }
  },
  async created() {
    const tourId = this.$route.params.id;
    await this.fetchTourDetails(tourId);
    await this.fetchCart();
    await this.checkIfPurchased(tourId);
  },
  methods: {
    async fetchTourDetails(id) {
      try {
        const res = await tourPublicService.getPublishedTours();
        this.tour = res.data.find(t => t.id == id);
      } catch (err) {
        console.error("Error fetching tour details:", err);
      }
    },
    async fetchCart() {
      try {
        const res = await purchaseService.getCart();
        this.cart = res.data;
      } catch (err) {
        console.error("Error fetching cart:", err);
      }
    },
    async checkIfPurchased(tourId) {
      try {
        const touristId = this.getCurrentUserId();
        if (!touristId) {
          this.isPurchased = false;
          return;
        }

        const res = await purchaseService.checkPurchase(tourId, touristId);
        this.isPurchased = res.data.purchased;
      } catch (err) {
        console.error("Error checking purchase status:", err);
        this.isPurchased = false;
      }
    },
    getCurrentUserId() {
      const token = localStorage.getItem('token');
      if (!token) return null;

      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.sub || payload.userId || payload.id || null;
      } catch (_) {
        return null;
      }
    },
    isInCart(tourId) {
      if (!this.cart || !this.cart.items) return false;
      return this.cart.items.some(item => item.tour_id === tourId);
    },
    async handleAddToCart() {
      try {
        const payload = {
          tour_id: this.tour.id,
          tour_name: this.tour.title,
          price: this.tour.price
        };
        await purchaseService.addItemToCart(payload);
        await this.fetchCart();
        alert(`Successfully added tour "${this.tour.title}" to the cart!`);
        window.dispatchEvent(new CustomEvent('cart-updated'));
      } catch (err) {
        alert(err.response?.data || "Failed to add item to cart.");
      }
    },
    async handleRemoveFromCart(tourId) {
      try {
        const itemInCart = this.cart.items.find(item => item.tour_id === tourId);
        if (itemInCart) {
          await purchaseService.removeItemFromCart(itemInCart.id);
          await this.fetchCart();
          window.dispatchEvent(new CustomEvent('cart-updated'));
        }
      } catch (err) {
        console.error("Error removing item from cart:", err);
      }
    },
    async handleStartTour() {
      try {
        // Simulator integracija: prvo uzmi poziciju turiste
        const posRes = await axios.get(`${API}/api/position`, { headers: getAuthHeader() });
        const { lat, lon } = posRes.data;

        // SAGA: pozovi startTour koji interno proverava da li je tura kupljena
        const res = await axios.post(
          `${API}/api/executions/start/${this.tour.id}`,
          { lat, lon },
          { headers: getAuthHeader() }
        );

        const execution = res.data;
        this.$router.push(`/active-tour/${execution.id}`);
      } catch (err) {
        if (err.response?.status === 403) {
          alert('You must purchase this tour before starting it.');
        } else if (err.response?.status === 404) {
          alert('Position not found. Please set your position in the simulator first.');
        } else {
          alert('Failed to start tour: ' + (err.response?.data || err.message));
        }
      }
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 1100px;
  margin: 0 auto;
  padding: 20px;
  font-family: 'Segoe UI', sans-serif;
}

.tour-header {
  margin-bottom: 25px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.btn-back {
  text-decoration: none;
  color: #28a745;
  font-weight: 500;
  display: inline-block;
  margin-bottom: 10px;
}

.tour-header h1 {
  margin: 5px 0;
  font-size: 2.2rem;
  color: #2c3e50;
}

.tour-layout {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 40px;
}

.info-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 25px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
}

.info-section h3 {
  margin-top: 0;
  color: #2c3e50;
  border-left: 4px solid #28a745;
  padding-left: 10px;
}

.description-text {
  color: #555;
  line-height: 1.6;
}

.restricted-access {
  background: #fff9e6;
  border: 1px solid #ffeeba;
  padding: 15px;
  border-radius: 6px;
  color: #856404;
}

.purchased-access {
  background: #e2f0d9;
  border: 1px solid #c5e1a5;
  padding: 15px;
  border-radius: 6px;
}

.access-success {
  color: #388e3c;
  font-weight: bold;
  margin-top: 0;
}

.keypoints-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-top: 15px;
}

.kp-card {
  display: flex;
  gap: 15px;
  background: white;
  padding: 12px;
  border-radius: 6px;
  border-left: 3px solid #28a745;
}

.kp-number {
  background: #28a745;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
}

.sticky-card {
  position: sticky;
  top: 90px;
  background: white;
  border: 1px solid #eee;
  border-radius: 10px;
  padding: 25px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
  text-align: center;
}

.price-box {
  margin-bottom: 20px;
}

.price-box .value {
  display: block;
  font-size: 2rem;
  font-weight: bold;
  color: #2c3e50;
}

.btn-add, .btn-remove, .btn-start {
  width: 100%;
  padding: 12px;
  border-radius: 6px;
  font-weight: bold;
  font-size: 1rem;
  cursor: pointer;
  border: none;
}

.btn-add { background: #28a745; color: white; }
.btn-add:hover { background: #218838; }
.btn-remove { background: #dc3545; color: white; }
.btn-remove:hover { background: #c82333; }
.btn-start { background: #007bff; color: white; margin-top: 10px; }
.btn-start:hover { background: #0069d9; }

.badge-purchased {
  background: #28a745;
  color: white;
  padding: 8px 20px;
  border-radius: 20px;
  font-weight: bold;
  display: inline-block;
}

.review-card {
  border-bottom: 1px solid #eee;
  padding: 15px 0;
}
.review-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  color: #7f8c8d;
}
.review-rating {
  font-weight: bold;
  color: #ffc107;
}

.difficulty-tag {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
  color: white;
  text-transform: uppercase;
}
.difficulty-tag.easy { background: #28a745; }
.difficulty-tag.medium { background: #ffc107; color: #333; }
.difficulty-tag.hard { background: #dc3545; }
</style>
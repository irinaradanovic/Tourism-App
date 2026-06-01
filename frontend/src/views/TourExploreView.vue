<template>
  <div class="container">
    <div v-if="!user || user.role !== 'TOURIST'" class="access-denied-box">
      <h2>Access Restricted</h2>
      <p>Only registered tourists are allowed to view available tours and manage a shopping cart.</p>
      <router-link to="/" class="btn-info" style="display:inline-block; max-width:200px; margin-top:15px;">Back to Home</router-link>
    </div>

    <div v-else>
      <div class="header-section">
        <h1>Explore Available Tours</h1>
        <p class="subtitle">Book your next adventure built by our certified guides.</p>
      </div>

      <div v-if="tours.length === 0" class="no-tours">
        No published tours available at the moment.
      </div>

      <div v-else class="tour-grid">
        <div v-for="tour in tours" :key="tour.id" class="tour-card">
          <div class="tour-image-wrapper">
            <div class="tour-placeholder-img">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#bbb" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
            </div>
            <span class="difficulty-tag" :class="tour.difficulty.toLowerCase()">{{ tour.difficulty }}</span>
          </div>

          <div class="tour-content">
            <h2>{{ tour.title }}</h2>
            <p class="tour-desc">{{ truncate(tour.description) }}</p>
            
            <div class="tags-wrapper">
              <span v-for="tag in tour.tags" :key="tag" class="tour-tag">#{{ tag }}</span>
            </div>

            <div class="price-section">
              <span class="price-label">Price:</span>
              <span class="price-value">{{ tour.price }} RSD</span>
            </div>
            
            <div class="card-footer">
              <router-link :to="{ name: 'tourDetailsTourist', params: { id: tour.id } }" class="btn-info">
                Details
              </router-link>
              <button
                v-if="purchasedTourIds.has(tour.id)"
                class="btn-purchased"
                disabled
              >
                Purchased
              </button>

              <button 
                v-else-if="isInCart(tour.id)" 
                @click="handleRemoveFromCart(tour.id)" 
                class="btn-remove"
              >
                Remove
              </button>
              <button 
                v-else 
                @click="handleAddToCart(tour)" 
                class="btn-add"
              >
                Add to Cart
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { purchaseService, tourPublicService } from '@/services/purchaseService';

export default {
  data() {
    return {
      tours: [],
      cart: { items: [] },
      user: JSON.parse(localStorage.getItem('user')),
      purchasedTourIds: new Set()
    }
  },
  async created() {
    if (this.user && this.user.role === 'TOURIST') {
      await this.fetchTours();
      await this.fetchPurchasedTours();
      await this.fetchCart();
    }
  },
  methods: {
    async fetchTours() {
      try {
        const res = await tourPublicService.getPublishedTours();
        this.tours = res.data;
      } catch (err) {
        console.error("Error fetching published tours:", err);
      }
    },
    async fetchPurchasedTours() {
      const touristId = this.getCurrentUserId();
      if (!touristId) return;

      const checks = await Promise.all(
        this.tours.map(async (tour) => {
          try {
            const res = await purchaseService.checkPurchase(tour.id, touristId);
            return res.data.purchased ? tour.id : null;
          } catch (_) {
            return null;
          }
        })
      );

      this.purchasedTourIds = new Set(checks.filter(Boolean));
    },
    async fetchCart() {
      try {
        const res = await purchaseService.getCart();
        this.cart = res.data;
      } catch (err) {
        console.error("Error fetching cart:", err);
      }
    },
    isInCart(tourId) {
      if (!this.cart || !this.cart.items) return false;
      return this.cart.items.some(item => item.tour_id === tourId);
    },
    async handleAddToCart(tour) {
      try {
        const payload = {
          tour_id: tour.id,
          tour_name: tour.title,
          price: tour.price
        };
        await purchaseService.addItemToCart(payload);
        await this.fetchCart();
        
        alert(`You have added "${tour.title}" to the cart!`);
        
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
        alert(err.response?.data || "Failed to remove item from cart.");
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
    truncate(text) {
      if (!text) return "";
      return text.length > 90 ? text.substring(0, 90) + '...' : text;
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

.access-denied-box {
  text-align: center;
  padding: 50px 20px;
  background: #fdf2f2;
  border: 1px solid #f5c6cb;
  border-radius: 8px;
  color: #721c24;
}

.header-section {
  margin-bottom: 35px;
}

.header-section h1 {
  font-size: 2.2rem;
  color: #2c3e50;
  margin: 0;
}

.subtitle {
  color: #7f8c8d;
  font-size: 1.1rem;
  margin-top: 5px;
}

.tour-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 25px;
}

.tour-card {
  border: 1px solid #eee;
  border-radius: 10px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  transition: transform 0.2s, box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.tour-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 18px rgba(0,0,0,0.1);
}

.tour-image-wrapper {
  width: 100%;
  height: 150px;
  background-color: #f7f9fa;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tour-placeholder-img { opacity: 0.6; }

.difficulty-tag {
  position: absolute;
  bottom: 10px;
  left: 10px;
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: bold;
  border-radius: 20px;
  text-transform: uppercase;
  color: white;
}
.difficulty-tag.easy { background-color: #28a745; }
.difficulty-tag.medium { background-color: #ffc107; color: #333; }
.difficulty-tag.hard { background-color: #dc3545; }

.tour-content {
  padding: 18px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.tour-content h2 {
  margin: 0 0 8px 0;
  font-size: 1.3rem;
  color: #2c3e50;
}

.tour-desc {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.4;
  margin-bottom: 12px;
  flex-grow: 1;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 15px;
}

.tour-tag {
  font-size: 0.8rem;
  color: #28a745;
  background: #f0f9f1;
  padding: 3px 8px;
  border-radius: 4px;
}

.price-section {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 15px;
  font-size: 1.05rem;
}

.price-label { color: #7f8c8d; }
.price-value { font-weight: bold; color: #2c3e50; }

.card-footer { display: flex; gap: 10px; }

.btn-info {
  background-color: transparent;
  border: 2px solid #28a745;
  color: #28a745;
  padding: 8px 16px;
  text-decoration: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  text-align: center;
  flex: 1;
  transition: all 0.2s;
}
.btn-info:hover { background-color: #f4faf5; }

.btn-add {
  background-color: #28a745;
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  flex: 2;
}
.btn-add:hover { background-color: #218838; }

.btn-remove {
  background-color: #dc3545;
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  flex: 2;
}
.btn-remove:hover { background-color: #c82333; }

.no-tours {
  text-align: center;
  color: #7f8c8d;
  font-size: 1.1rem;
  margin-top: 40px;
}
</style>
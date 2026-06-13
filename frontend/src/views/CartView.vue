<template>
  <div class="container">
    <div class="header-section">
      <h1>Your Shopping Cart</h1>
      <p class="subtitle">Review your selected tours before completing the purchase.</p>
    </div>

    <div v-if="!cart || !cart.items || cart.items.length === 0" class="empty-cart-message">
      <p>Your cart is currently empty.</p>
      <router-link to="/tours" class="btn-main">Browse Tours</router-link>
    </div>

    <div v-else class="cart-content">
      <div class="cart-table-wrapper">
        <table class="cart-table">
          <thead>
            <tr>
              <th>Tour Name</th>
              <th class="text-right">Price</th>
              <th class="text-center">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in cart.items" :key="item.id">
              <td class="tour-title-cell">
                <strong>{{ item.tour_name }}</strong>
              </td>
              <td class="text-right price-cell">{{ item.price }} RSD</td>
              <td class="text-center">
                <button @click="handleRemoveItem(item.id)" class="btn-action-delete">
                  Remove
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="checkout-summary-card">
        <div class="summary-row">
          <span>Total Price:</span>
          <strong class="total-price-value">{{ cart.total_price }} RSD</strong>
        </div>
        <hr class="divider" />
        <button @click="handleCheckout" class="btn-checkout" :disabled="loading">
          {{ loading ? 'Processing...' : 'Proceed to Purchase' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { purchaseService } from '@/services/purchaseService';

export default {
  data() {
    return {
      cart: { items: [], total_price: 0, status: 'ACTIVE' },
      loading: false
    }
  },
  async created() {
    await this.fetchCart();
  },
  methods: {
    async fetchCart() {
      try {
        const res = await purchaseService.getCart();
        this.cart = res.data;
        this.notifyCartUpdate();
      } catch (err) {
        console.error("Error loading cart:", err);
      }
    },
    async handleRemoveItem(itemId) {
      try {
        await purchaseService.removeItemFromCart(itemId);
        await this.fetchCart(); 
      } catch (err) {
        console.error("Error removing item:", err);
      }
    },
    notifyCartUpdate() {
      window.dispatchEvent(new CustomEvent('cart-updated'));
    },
    async handleCheckout() {
      if (this.loading) {
        return;
      }

      try {
        this.loading = true; 
        await purchaseService.checkoutCart();

        await new Promise(resolve => setTimeout(resolve, 2500));
        await this.fetchCart();

        const itemsCount = this.cart.items ? this.cart.items.length : 0;
        const currentStatus = this.cart.status || "ACTIVE";

        if (currentStatus === "ACTIVE" && itemsCount === 0) {
          alert("Checkout successful!");
          this.notifyCartUpdate(); 
        } else if (currentStatus === "ACTIVE" && itemsCount > 0) {
          alert("Checkout failed: One or more tours are unavailable (Status: DRAFT/ARCHIVED).");
        } else {
          await new Promise(resolve => setTimeout(resolve, 1500));
          await this.fetchCart();
          
          const finalCount = this.cart.items ? this.cart.items.length : 0;
          if (finalCount === 0) {
            alert("Checkout successful!");
            this.notifyCartUpdate();
          } else {
            alert("Checkout is taking longer than expected. Please check your tokens later.");
          }
        }

      } catch (err) {
        console.error("Checkout failed:", err);
        alert("Failed to initiate checkout.");
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header-section {
  margin-bottom: 30px;
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

.empty-cart-message {
  text-align: center;
  padding: 40px;
  background: #f9f9f9;
  border-radius: 8px;
  color: #7f8c8d;
  font-size: 1.2rem;
}

.btn-main {
  display: inline-block;
  background-color: #28a745;
  color: white;
  padding: 10px 20px;
  text-decoration: none;
  border-radius: 5px;
  font-weight: bold;
  margin-top: 15px;
}

/* Tabela */
.cart-content {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.cart-table-wrapper {
  background: white;
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 5px rgba(0,0,0,0.02);
}

.cart-table {
  width: 100%;
  border-collapse: collapse;
}

.cart-table th, .cart-table td {
  padding: 15px;
  text-align: left;
}

.cart-table th {
  background-color: #f8f9fa;
  color: #7f8c8d;
  font-weight: 600;
  border-bottom: 2px solid #eee;
}

.cart-table tbody tr {
  border-bottom: 1px solid #eee;
}

.tour-title-cell {
  color: #2c3e50;
  font-size: 1.05rem;
}

.price-cell {
  font-weight: 600;
  color: #2c3e50;
}

.text-right { text-align: right; }
.text-center { text-align: center; }

.btn-action-delete {
  background-color: transparent;
  border: 1px solid #dc3545;
  color: #dc3545;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-action-delete:hover {
  background-color: #dc3545;
  color: white;
}

.checkout-summary-card {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 20px;
  align-self: flex-end;
  width: 100%;
  max-width: 350px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
}

.summary-row {
  display: flex;
  justify-content: space-between;
  font-size: 1.2rem;
  color: #2c3e50;
  margin-bottom: 15px;
}

.total-price-value {
  color: #28a745;
}

.divider {
  border: 0;
  border-top: 1px solid #eee;
  margin-bottom: 15px;
}

.btn-checkout {
  width: 100%;
  background-color: #28a745;
  border: none;
  color: white;
  padding: 12px;
  border-radius: 6px;
  font-size: 1.1rem;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-checkout:hover {
  background-color: #218838;
}

.btn-checkout:disabled {
  background-color: #8bc99a;
  cursor: not-allowed;
}
</style>

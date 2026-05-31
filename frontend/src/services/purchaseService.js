import axios from 'axios';

const API_URL = 'http://localhost:80/api/purchase';


const getAuthHeader = () => {
  const token = localStorage.getItem('token');

  if (token) {
    return {
      Authorization: `Bearer ${token}`
    };
  }

  return {};
};

export const purchaseService = {
  getCart() {
    return axios.get(`${API_URL}/cart`, { headers: getAuthHeader() });
  },
  addItemToCart(payload) {
    return axios.post(`${API_URL}/cart/items`, payload, { headers: getAuthHeader() });
  },
  removeItemFromCart(itemId) {
    return axios.delete(`${API_URL}/cart/items/${itemId}`, { headers: getAuthHeader() });
  },
  checkoutCart() {
  return axios.post(`${API_URL}/checkout`, {}, { headers: getAuthHeader() });
}
};

export const tourPublicService = {
  getPublishedTours() {
    return axios.get('http://localhost:8083/api/tours/published'); 
  }
};
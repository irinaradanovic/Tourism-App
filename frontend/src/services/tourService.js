import axios from 'axios'

const API_URL = 'http://localhost:80/api'

function getHeaders() {
  const token = localStorage.getItem('token')
  return { Authorization: `Bearer ${token}` }
}

export const tourService = {

  getMyPosition() {
    return axios.get(`${API_URL}/position`, { headers: getHeaders() })
  },

  savePosition(lat, lon) {
    return axios.post(`${API_URL}/position`, { lat, lon }, { headers: getHeaders() })
  },

  createTour(data) {
    return axios.post(`${API_URL}/tours`, data, { headers: getHeaders() })
  },

  getMyTours() {
    return axios.get(`${API_URL}/tours/my`, { headers: getHeaders() })
  },

  getPublishedTours() {
    return axios.get(`${API_URL}/tours`, { headers: getHeaders() })
  },

  getTourById(id) {
    return axios.get(`${API_URL}/tours/${id}`, { headers: getHeaders() })
  },

  addKeyPoint(tourId, keyPoint) {
    return axios.post(`${API_URL}/tours/${tourId}/keypoints`, keyPoint, { headers: getHeaders() })
  },

  updateKeyPoint(tourId, index, keyPoint) {
    return axios.put(`${API_URL}/tours/${tourId}/keypoints/${index}`, keyPoint, { headers: getHeaders() })
  },

  deleteKeyPoint(tourId, index) {
    return axios.delete(`${API_URL}/tours/${tourId}/keypoints/${index}`, { headers: getHeaders() })
  },

  addReview(tourId, review) {
    return axios.post(`${API_URL}/tours/${tourId}/reviews`, review, { headers: getHeaders() })
  }
}
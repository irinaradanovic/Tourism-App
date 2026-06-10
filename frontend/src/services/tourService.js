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
    return axios.get(`${API_URL}/tours/published`, { headers: getHeaders() })
  },

  getTourById(id) {
    return axios.get(`${API_URL}/tours/${id}`, { headers: getHeaders() })
  },

  addKeyPoint(tourId, keyPoint) {
    return axios.post(`${API_URL}/tours/${tourId}/key-points`, keyPoint, { headers: getHeaders() })
  },

  updateKeyPoint(tourId, index, keyPoint) {
    return axios.put(`${API_URL}/tours/${tourId}/key-points/${index}`, keyPoint, { headers: getHeaders() })
  },

  deleteKeyPoint(tourId, index) {
    return axios.delete(`${API_URL}/tours/${tourId}/key-points/${index}`, { headers: getHeaders() })
  },

  updateDurations(tourId, durations) {
    return axios.put(`${API_URL}/tours/${tourId}/durations`, durations, { headers: getHeaders() })
  },

  publishTour(tourId) {
    return axios.post(`${API_URL}/tours/${tourId}/publish`, {}, { headers: getHeaders() })
  },

  archiveTour(tourId) {
    return axios.post(`${API_URL}/tours/${tourId}/archive`, {}, { headers: getHeaders() })
  },

  reactivateTour(tourId) {
    return axios.post(`${API_URL}/tours/${tourId}/reactivate`, {}, { headers: getHeaders() })
  },

addReview(tourId, reviewFormData) {
  return axios.post(`${API_URL}/tours/${tourId}/reviews`, reviewFormData, {
    headers: { ...getHeaders(), 'Content-Type': 'multipart/form-data' }
  })
}
}

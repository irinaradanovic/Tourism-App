import axios from 'axios'

const API_URL = 'http://localhost:8083/api/tours'

function authHeader() {
    const token = localStorage.getItem('token')
    return {
        headers: {
            Authorization: `Bearer ${token}`
        }
    }
}

export const tourService = {

    createTour(data) {
        return axios.post(API_URL, data, authHeader())
    },

    getMyTours() {
        return axios.get(`${API_URL}/my`, authHeader())
    },

    getTourById(id) {
        return axios.get(`${API_URL}/${id}`, authHeader())
    },

    addKeyPoint(tourId, formData) {
        return axios.post(
            `${API_URL}/${tourId}/key-points`,
            formData,
            {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                    'Content-Type': 'multipart/form-data'
                }
            }
        )
    }
}
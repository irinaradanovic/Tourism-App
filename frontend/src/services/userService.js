import axios from 'axios'

const API_URL = 'http://localhost:8082/api/users'

export const userService = {

    getMyProfile() {

        const token = localStorage.getItem('token')

        return axios.get(`${API_URL}/profile`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
    },

    getUserProfile(userId) {
        const token = localStorage.getItem('token')
        return axios.get(`${API_URL}/${userId}`, {
            headers: { Authorization: `Bearer ${token}` }
        })
    },
    
    toggleBlockUser(userId) {
        const token = localStorage.getItem('token')

        return axios.put(`${API_URL}/${userId}/toggle-block`, {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
    }

}
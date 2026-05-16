import axios from 'axios'

const API_URL = 'http://localhost:8000/api/followers' 

export const followerService = {
    
    getFollowers() {
        const token = localStorage.getItem('token')
        return axios.get(`${API_URL}/my-followers`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
    },

    getFollowing() {
        const token = localStorage.getItem('token')
        return axios.get(`${API_URL}/my-followings`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
    },

    getUserFollowers(userId) {
        const token = localStorage.getItem('token')
        return axios.get(`${API_URL}/${userId}/followers`, {
            headers: { Authorization: `Bearer ${token}` }
        })
    },

    getUserFollowing(userId) {
        const token = localStorage.getItem('token')
        return axios.get(`${API_URL}/${userId}/followings`, {
            headers: { Authorization: `Bearer ${token}` }
        })
    },

    follow(userId) {
        const token = localStorage.getItem('token')
        return axios.post(`${API_URL}/follow?followedId=${userId}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
        })
    },

    unfollow(userId) {
    const token = localStorage.getItem('token')
    return axios.delete(`${API_URL}/unfollow/${userId}`, {
        headers: { Authorization: `Bearer ${token}` }
    })
    },

    getRecommendations() {
        const token = localStorage.getItem('token')
        return axios.get(`${API_URL}/recommendations`, {
            headers: { Authorization: `Bearer ${token}` }
        })
    }
}
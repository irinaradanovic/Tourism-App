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
    }
}
import axios from 'axios'

const API_URL = 'http://localhost:8082/api'

export const register = (userData) => {
    return axios.post(`${API_URL}/auth/register`, userData)
}

export const login = (credentials) => {
    return axios.post(`${API_URL}/auth/login`, credentials)
}

export const getAllUsers = () => {
    const token = localStorage.getItem('token')
    return axios.get(`${API_URL}/users`, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })
}
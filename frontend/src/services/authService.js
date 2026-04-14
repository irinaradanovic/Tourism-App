import axios from 'axios'
import router from '../router'

const API_URL = 'http://localhost:8082/api'

axios.interceptors.response.use(
    response => response,
    error => {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
            localStorage.removeItem('user')
            localStorage.removeItem('token')
            router.push('/login')
        }
        return Promise.reject(error)
    }
)

export const register = (userData) => {
    return axios.post(`${API_URL}/auth/register`, userData)
}

export const login = (credentials) => {
    return axios.post(`${API_URL}/auth/login`, credentials)
}

export const logout = () => {
    return axios.post(`${API_URL}/auth/logout`)
}

export const getAllUsers = () => {
    const token = localStorage.getItem('token')
    return axios.get(`${API_URL}/users`, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })
}
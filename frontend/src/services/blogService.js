import axios from 'axios';
import router from '@/router'; 

const API_URL = 'http://localhost:80';

export const blogService = {

  getAllBlogs() {
    const token = localStorage.getItem('token');
    return axios.get(`${API_URL}/blogs`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
  },

  getBlogById(id) {
    const token = localStorage.getItem('token');
    return axios.get(`${API_URL}/blogs/${id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
  },

  async createBlog(formData) {
    const token = localStorage.getItem('token');
    try {
      const res = await axios.post(`${API_URL}/blogs`, formData, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'multipart/form-data'
        }
      });
      router.push({ name: 'blogDetails', params: { id: res.data.id } });
      return res.data;
    } catch (err) {
      console.error("Error while creating blog:", err);
      throw err;
    }
  },

  toggleLike(blogId) {
    const token = localStorage.getItem('token');
    return axios.post(`${API_URL}/blogs/${blogId}/like`, {}, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
  },

  getBlogsByAuthor(authorId) {
    const token = localStorage.getItem('token');
    return axios.get(`${API_URL}/blogs/author/${authorId}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
  },

  getComments(blogId) {
    const token = localStorage.getItem('token');
    return axios.get(`${API_URL}/blogs/${blogId}/comments`, {
      headers: { Authorization: `Bearer ${token}` }
    });
  },

  addComment(blogId, data) {
    const token = localStorage.getItem('token');
    return axios.post(`${API_URL}/blogs/${blogId}/comments`, data, {
      headers: { Authorization: `Bearer ${token}` }
    });
  }
};
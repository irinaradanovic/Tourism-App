import axios from 'axios';
import router from '@/router'; 

const API_URL = 'http://localhost:8081'; 

export const blogService = {

  getAllBlogs() {
    return axios.get(`${API_URL}/blogs`);
  },

  getBlogById(id) {
    return axios.get(`${API_URL}/blogs/${id}`);
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
  }
};
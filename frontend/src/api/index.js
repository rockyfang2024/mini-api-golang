import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Attach JWT from localStorage on every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

// Auth
export const register = (username, email, password) =>
  api.post('/auth/register', { username, email, password })

export const login = (username, password) =>
  api.post('/auth/login', { username, password })

export const getMe = () => api.get('/me')

// Posts
export const getPosts = () => api.get('/posts')
export const createPost = (content, visibility) =>
  api.post('/posts', { content, visibility })

export const getUserPosts = (userId) => api.get(`/users/${userId}/posts`)

export default api

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
export const createPost = (content, visibility, images = []) => {
  const formData = new FormData()
  formData.append('content', content)
  formData.append('visibility', visibility)
  images.forEach((file) => formData.append('images', file))
  return api.post('/posts', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const getUserPosts = (userId) => api.get(`/users/${userId}/posts`)
export const getUserProfile = (userId) => api.get(`/users/${userId}`)

// Avatar
export const uploadAvatar = (file) => {
  const formData = new FormData()
  formData.append('avatar', file)
  return api.post('/me/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// Comments
export const getComments = (postId) => api.get(`/posts/${postId}/comments`)
export const createComment = (postId, content) =>
  api.post(`/posts/${postId}/comments`, { content })
export const replyComment = (commentId, content) =>
  api.post(`/comments/${commentId}/replies`, { content })

// Likes
export const likePost = (postId) => api.post(`/posts/${postId}/like`)
export const unlikePost = (postId) => api.delete(`/posts/${postId}/like`)

// Reposts
export const repostPost = (postId) => api.post(`/posts/${postId}/repost`)

// Follows
export const followUser = (userId) => api.post(`/users/${userId}/follow`)
export const unfollowUser = (userId) => api.delete(`/users/${userId}/follow`)
export const getFollowers = (userId, page = 1, pageSize = 20) =>
  api.get(`/users/${userId}/followers`, { params: { page, page_size: pageSize } })
export const getFollowing = (userId, page = 1, pageSize = 20) =>
  api.get(`/users/${userId}/following`, { params: { page, page_size: pageSize } })

// Notifications
export const getNotifications = (page = 1, pageSize = 20) =>
  api.get('/notifications', { params: { page, page_size: pageSize } })
export const markNotificationRead = (notificationId) =>
  api.put(`/notifications/${notificationId}/read`)
export const markAllNotificationsRead = () => api.put('/notifications/read-all')

// Settings
export const getSettings = () => api.get('/settings')
export const updateSettings = (settings) => api.put('/settings', settings)

export default api

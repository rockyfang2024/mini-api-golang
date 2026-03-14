<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../api/index.js'
import { authStore } from '../stores/auth'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  if (!username.value || !password.value) {
    error.value = '请填写用户名和密码'
    return
  }
  loading.value = true
  try {
    const res = await login(username.value, password.value)
    authStore.setAuth(res.data.data.token, res.data.data.user)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.message || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-wrapper">
    <div class="card auth-card">
      <h2 class="auth-title">登录</h2>
      <div class="form-group">
        <label>用户名</label>
        <input v-model="username" type="text" placeholder="请输入用户名" @keyup.enter="handleLogin" />
      </div>
      <div class="form-group">
        <label>密码</label>
        <input v-model="password" type="password" placeholder="请输入密码" @keyup.enter="handleLogin" />
      </div>
      <p v-if="error" class="error-msg">{{ error }}</p>
      <button class="btn btn-primary" style="width:100%" :disabled="loading" @click="handleLogin">
        {{ loading ? '登录中…' : '登录' }}
      </button>
      <p class="auth-footer">还没有账号？<router-link to="/register">立即注册</router-link></p>
    </div>
  </div>
</template>

<style scoped>
.auth-wrapper { display: flex; justify-content: center; padding-top: 40px; }
.auth-card { width: 100%; max-width: 400px; }
.auth-title { font-size: 1.5rem; font-weight: 700; margin-bottom: 24px; text-align: center; }
.auth-footer { text-align: center; margin-top: 16px; font-size: .9rem; color: #666; }
.auth-footer a { color: #1da1f2; text-decoration: none; }
</style>

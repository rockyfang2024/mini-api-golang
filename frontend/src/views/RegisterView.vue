<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api/index.js'

const router = useRouter()
const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''
  success.value = ''
  if (!username.value || !email.value || !password.value) {
    error.value = '请填写所有字段'
    return
  }
  if (password.value.length < 6) {
    error.value = '密码至少 6 位'
    return
  }
  loading.value = true
  try {
    await register(username.value, email.value, password.value)
    success.value = '注册成功！正在跳转到登录页…'
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    error.value = e.response?.data?.message || '注册失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-wrapper">
    <div class="card auth-card">
      <h2 class="auth-title">注册</h2>
      <div class="form-group">
        <label>用户名</label>
        <input v-model="username" type="text" placeholder="请输入用户名" />
      </div>
      <div class="form-group">
        <label>邮箱</label>
        <input v-model="email" type="email" placeholder="请输入邮箱" />
      </div>
      <div class="form-group">
        <label>密码（至少 6 位）</label>
        <input v-model="password" type="password" placeholder="请输入密码" @keyup.enter="handleRegister" />
      </div>
      <p v-if="error" class="error-msg">{{ error }}</p>
      <p v-if="success" class="success-msg">{{ success }}</p>
      <button class="btn btn-primary" style="width:100%" :disabled="loading" @click="handleRegister">
        {{ loading ? '注册中…' : '注册' }}
      </button>
      <p class="auth-footer">已有账号？<router-link to="/login">立即登录</router-link></p>
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

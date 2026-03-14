<script setup>
import { useRouter } from 'vue-router'
import { authStore } from './stores/auth'

const router = useRouter()

function logout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-brand">
        <router-link to="/">🐦 简易微博</router-link>
      </div>
      <div class="nav-links">
        <router-link to="/">首页</router-link>
        <template v-if="authStore.isLoggedIn">
          <router-link to="/my-posts">我的动态</router-link>
          <span class="nav-user" @click="router.push('/my-posts')">
            <span class="avatar">{{ authStore.user?.username?.[0]?.toUpperCase() }}</span>
            {{ authStore.user?.username }}
          </span>
          <button class="btn-logout" @click="logout">退出</button>
        </template>
        <template v-else>
          <router-link to="/login">登录</router-link>
          <router-link to="/register">注册</router-link>
        </template>
      </div>
    </nav>
    <main class="container">
      <router-view />
    </main>
  </div>
</template>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f0f2f5; color: #333; }
#app { min-height: 100vh; }

.navbar {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 4px rgba(0,0,0,.08);
}
.nav-brand a { font-size: 1.2rem; font-weight: 700; color: #1da1f2; text-decoration: none; }
.nav-links { display: flex; align-items: center; gap: 16px; }
.nav-links a { color: #555; text-decoration: none; font-size: .95rem; }
.nav-links a:hover, .nav-links a.router-link-active { color: #1da1f2; }
.nav-user { cursor: pointer; display: flex; align-items: center; gap: 6px; font-size: .95rem; color: #333; }
.avatar { width: 30px; height: 30px; border-radius: 50%; background: #1da1f2; color: #fff; display: inline-flex; align-items: center; justify-content: center; font-weight: 700; font-size: .85rem; }
.btn-logout { background: none; border: 1px solid #ddd; border-radius: 6px; padding: 4px 12px; cursor: pointer; color: #666; font-size: .9rem; }
.btn-logout:hover { background: #f5f5f5; }
.container { max-width: 640px; margin: 0 auto; padding: 24px 16px; }

/* Shared form styles */
.card { background: #fff; border-radius: 12px; padding: 24px; box-shadow: 0 1px 4px rgba(0,0,0,.08); margin-bottom: 16px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-size: .9rem; color: #555; font-weight: 500; }
.form-group input, .form-group textarea, .form-group select {
  width: 100%; padding: 10px 12px; border: 1px solid #ddd; border-radius: 8px;
  font-size: .95rem; outline: none; transition: border-color .2s;
}
.form-group input:focus, .form-group textarea:focus, .form-group select:focus { border-color: #1da1f2; }
.form-group textarea { resize: vertical; min-height: 80px; }
.btn { padding: 10px 20px; border: none; border-radius: 8px; cursor: pointer; font-size: .95rem; font-weight: 600; transition: opacity .2s; }
.btn:disabled { opacity: .6; cursor: not-allowed; }
.btn-primary { background: #1da1f2; color: #fff; }
.btn-primary:hover:not(:disabled) { background: #0c8de0; }
.error-msg { color: #e0245e; font-size: .9rem; margin-top: 8px; }
.success-msg { color: #17bf63; font-size: .9rem; margin-top: 8px; }

/* Post card */
.post-card { background: #fff; border-radius: 12px; padding: 16px; box-shadow: 0 1px 4px rgba(0,0,0,.08); margin-bottom: 12px; }
.post-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.post-author { font-weight: 600; font-size: .95rem; }
.post-time { font-size: .82rem; color: #888; }
.post-badge { font-size: .78rem; padding: 2px 8px; border-radius: 20px; margin-left: auto; }
.post-badge.public { background: #e8f5fd; color: #1da1f2; }
.post-badge.private { background: #fff3e0; color: #f57c00; }
.post-content { font-size: .97rem; line-height: 1.6; white-space: pre-wrap; word-break: break-word; }
.empty-state { text-align: center; color: #888; padding: 48px 0; font-size: 1rem; }
</style>

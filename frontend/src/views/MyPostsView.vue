<script setup>
import { ref, onMounted } from 'vue'
import { getUserPosts } from '../api/index.js'
import { authStore } from '../stores/auth'
import PostCard from '../components/PostCard.vue'

const posts = ref([])
const loading = ref(false)
const error = ref('')

async function fetchMyPosts() {
  loading.value = true
  error.value = ''
  try {
    const res = await getUserPosts(authStore.user.id)
    posts.value = res.data.data || []
  } catch (e) {
    error.value = '加载失败，请刷新重试'
  } finally {
    loading.value = false
  }
}

onMounted(fetchMyPosts)
</script>

<template>
  <div>
    <div class="card" style="margin-bottom:20px">
      <div style="display:flex;align-items:center;gap:14px">
        <span class="avatar avatar-lg">{{ authStore.user?.username?.[0]?.toUpperCase() }}</span>
        <div>
          <div style="font-size:1.2rem;font-weight:700">{{ authStore.user?.username }}</div>
          <div style="font-size:.85rem;color:#888">{{ authStore.user?.email }}</div>
        </div>
      </div>
    </div>

    <h3 style="margin-bottom:12px;font-size:1rem;color:#555">我的全部动态（包含私密）</h3>

    <div v-if="loading" class="empty-state">加载中…</div>
    <div v-else-if="error" class="empty-state" style="color:#e0245e">{{ error }}</div>
    <div v-else-if="posts.length === 0" class="empty-state">还没有动态，去首页发一条吧 ✏️</div>
    <template v-else>
      <PostCard v-for="post in posts" :key="post.id" :post="post" />
    </template>
  </div>
</template>

<style scoped>
.avatar-lg { width: 48px; height: 48px; font-size: 1.2rem; border-radius: 50%; background: #1da1f2; color: #fff; display: inline-flex; align-items: center; justify-content: center; font-weight: 700; }
</style>

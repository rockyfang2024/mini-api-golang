<script setup>
import { ref, onMounted, watch } from 'vue'
import { getPosts, createPost, getFollowing } from '../api/index.js'
import { authStore } from '../stores/auth'
import PostCard from '../components/PostCard.vue'

const posts = ref([])
const content = ref('')
const visibility = ref('public')
const error = ref('')
const postError = ref('')
const loading = ref(false)
const posting = ref(false)
const followingIds = ref(new Set())

async function fetchPosts() {
  loading.value = true
  error.value = ''
  try {
    const res = await getPosts()
    posts.value = res.data.data || []
  } catch (e) {
    error.value = '加载动态失败，请刷新重试'
  } finally {
    loading.value = false
  }
}

async function fetchFollowing() {
  if (!authStore.isLoggedIn || !authStore.user?.id) {
    followingIds.value = new Set()
    return
  }
  try {
    const res = await getFollowing(authStore.user.id, 1, 100)
    const ids = new Set(res.data.data.following?.map((f) => f.following_id))
    followingIds.value = ids
  } catch (e) {
    followingIds.value = new Set()
  }
}

function isFollowingAuthor(post) {
  const authorId = post.author_id || post.author?.id
  return authorId ? followingIds.value.has(authorId) : false
}

function handleFollowChange({ userId, isFollowing }) {
  const next = new Set(followingIds.value)
  if (isFollowing) {
    next.add(userId)
  } else {
    next.delete(userId)
  }
  followingIds.value = next
}

async function handlePost() {
  postError.value = ''
  if (!content.value.trim()) {
    postError.value = '动态内容不能为空'
    return
  }
  posting.value = true
  try {
    const res = await createPost(content.value, visibility.value)
    posts.value.unshift(res.data.data)
    content.value = ''
    visibility.value = 'public'
  } catch (e) {
    postError.value = e.response?.data?.message || '发布失败，请重试'
  } finally {
    posting.value = false
  }
}

onMounted(fetchPosts)

watch(
  () => authStore.isLoggedIn,
  (loggedIn) => {
    if (loggedIn) {
      fetchFollowing()
    } else {
      followingIds.value = new Set()
    }
  },
  { immediate: true }
)
</script>

<template>
  <div>
    <!-- Post compose box (only for logged-in users) -->
    <div v-if="authStore.isLoggedIn" class="card">
      <div class="compose-header">
        <span class="avatar">{{ authStore.user?.username?.[0]?.toUpperCase() }}</span>
        <span class="compose-hint">有什么新鲜事？</span>
      </div>
      <div class="form-group" style="margin-top:12px">
        <textarea v-model="content" placeholder="写下你的动态…" rows="3"></textarea>
      </div>
      <div class="compose-footer">
        <select v-model="visibility" class="visibility-select">
          <option value="public">🌐 所有人可见</option>
          <option value="private">🔒 仅自己可见</option>
        </select>
        <button class="btn btn-primary" :disabled="posting" @click="handlePost">
          {{ posting ? '发布中…' : '发布' }}
        </button>
      </div>
      <p v-if="postError" class="error-msg" style="margin-top:8px">{{ postError }}</p>
    </div>

    <!-- Feed -->
    <div v-if="loading" class="empty-state">加载中…</div>
    <div v-else-if="error" class="empty-state" style="color:#e0245e">{{ error }}</div>
    <div v-else-if="posts.length === 0" class="empty-state">暂无动态，快来发第一条吧 🎉</div>
    <template v-else>
      <PostCard
        v-for="post in posts"
        :key="post.id"
        :post="post"
        :show-follow="authStore.isLoggedIn"
        :is-following="isFollowingAuthor(post)"
        @follow-change="handleFollowChange"
      />
    </template>
  </div>
</template>

<style scoped>
.compose-header { display: flex; align-items: center; gap: 10px; }
.compose-hint { color: #888; font-size: .95rem; }
.compose-footer { display: flex; align-items: center; justify-content: space-between; }
.visibility-select { padding: 8px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: .9rem; outline: none; }
</style>

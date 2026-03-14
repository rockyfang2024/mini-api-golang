<script setup>
import { computed, ref, watch } from 'vue'
import { authStore } from '../stores/auth'
import { likePost, unlikePost, repostPost, followUser, unfollowUser } from '../api/index.js'

const props = defineProps({
  post: {
    type: Object,
    required: true,
  },
  showFollow: {
    type: Boolean,
    default: false,
  },
  isFollowing: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['follow-change'])

const likeCount = ref(0)
const repostCount = ref(0)
const isLiked = ref(false)
const isReposted = ref(false)
const following = ref(false)
const actionError = ref('')

const authorId = computed(() => props.post.author_id || props.post.author?.id)
const isOwnPost = computed(() => authStore.user?.id && authStore.user.id === authorId.value)
const canFollow = computed(() => props.showFollow && authStore.isLoggedIn && authorId.value && !isOwnPost.value)
const avatarUrl = computed(() => props.post.author?.avatar_url || '')

watch(
  () => props.post,
  (post) => {
    likeCount.value = post?.like_count ?? 0
    repostCount.value = post?.repost_count ?? 0
    isLiked.value = !!post?.is_liked
    isReposted.value = !!post?.is_reposted
  },
  { immediate: true }
)

watch(
  () => props.isFollowing,
  (value) => {
    following.value = value
  },
  { immediate: true }
)

function formatTime(ts) {
  return new Date(ts).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}

async function toggleLike() {
  if (!authStore.isLoggedIn) return
  actionError.value = ''
  try {
    if (isLiked.value) {
      await unlikePost(props.post.id)
      isLiked.value = false
      likeCount.value = Math.max(0, likeCount.value - 1)
      return
    }
    await likePost(props.post.id)
    isLiked.value = true
    likeCount.value += 1
  } catch (e) {
    actionError.value = e.response?.data?.message || '操作失败，请重试'
  }
}

async function handleRepost() {
  if (!authStore.isLoggedIn || isReposted.value) return
  actionError.value = ''
  try {
    await repostPost(props.post.id)
    isReposted.value = true
    repostCount.value += 1
  } catch (e) {
    actionError.value = e.response?.data?.message || '转发失败，请稍后再试'
  }
}

async function toggleFollow() {
  if (!canFollow.value) return
  actionError.value = ''
  try {
    if (following.value) {
      await unfollowUser(authorId.value)
      following.value = false
    } else {
      await followUser(authorId.value)
      following.value = true
    }
    emit('follow-change', { userId: authorId.value, isFollowing: following.value })
  } catch (e) {
    actionError.value = e.response?.data?.message || '操作失败，请重试'
  }
}
</script>

<template>
  <div class="post-card">
    <div class="post-header">
      <span class="avatar">
        <img v-if="avatarUrl" :src="avatarUrl" alt="avatar" class="avatar-img" />
        <span v-else>{{ post.author?.username?.[0]?.toUpperCase() || '?' }}</span>
      </span>
      <div>
        <div class="post-author">{{ post.author?.username || '未知用户' }}</div>
        <div class="post-time">{{ formatTime(post.created_at) }}</div>
      </div>
      <span class="post-badge" :class="post.visibility">
        {{ post.visibility === 'private' ? '🔒 私密' : '🌐 公开' }}
      </span>
    </div>
    <div class="post-content">{{ post.content }}</div>
    <div class="post-actions">
      <button class="action-btn" :class="{ active: isLiked }" :disabled="!authStore.isLoggedIn" @click="toggleLike">
        👍 {{ likeCount }}
      </button>
      <button class="action-btn" :class="{ active: isReposted }" :disabled="!authStore.isLoggedIn || isReposted" @click="handleRepost">
        🔁 {{ repostCount }}
      </button>
      <button v-if="canFollow" class="action-btn follow-btn" @click="toggleFollow">
        {{ following ? '已关注' : '关注' }}
      </button>
    </div>
    <p v-if="actionError" class="error-msg" style="margin-top:8px">{{ actionError }}</p>
  </div>
</template>

<style scoped>
.post-actions { display: flex; align-items: center; gap: 12px; margin-top: 12px; }
.action-btn {
  border: 1px solid #e0e0e0;
  background: #fff;
  border-radius: 18px;
  padding: 6px 12px;
  font-size: .85rem;
  cursor: pointer;
  color: #555;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.action-btn.active { border-color: #1da1f2; color: #1da1f2; }
.action-btn:disabled { opacity: .6; cursor: not-allowed; }
.follow-btn { margin-left: auto; }
.avatar-img { width: 100%; height: 100%; border-radius: 50%; object-fit: cover; }
</style>

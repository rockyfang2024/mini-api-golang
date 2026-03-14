<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getUserProfile, getUserPosts, getFollowers, getFollowing, followUser, unfollowUser } from '../api/index.js'
import { authStore } from '../stores/auth'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const profile = ref(null)
const posts = ref([])
const loading = ref(false)
const profileError = ref('')
const postsError = ref('')
const followLoading = ref(false)
const following = ref(false)
const followersTotal = ref(0)
const followingTotal = ref(0)

const profileId = computed(() => Number(route.params.id))
const isOwnProfile = computed(() => authStore.user?.id === profileId.value)
const canFollow = computed(() => authStore.isLoggedIn && !isOwnProfile.value)

async function fetchProfile() {
  profileError.value = ''
  try {
    const res = await getUserProfile(profileId.value)
    profile.value = res.data.data
  } catch (e) {
    profileError.value = '用户信息加载失败'
  }
}

async function fetchPosts() {
  loading.value = true
  postsError.value = ''
  try {
    const res = await getUserPosts(profileId.value)
    posts.value = res.data.data || []
  } catch (e) {
    if (e.response?.status === 403) {
      postsError.value = e.response?.data?.message || '无权限查看该用户动态'
    } else {
      postsError.value = '动态加载失败'
    }
    posts.value = []
  } finally {
    loading.value = false
  }
}

async function fetchStats() {
  try {
    const [followersRes, followingRes] = await Promise.all([
      getFollowers(profileId.value, 1, 1),
      getFollowing(profileId.value, 1, 1),
    ])
    followersTotal.value = followersRes.data.data.total || 0
    followingTotal.value = followingRes.data.data.total || 0
  } catch (e) {
    followersTotal.value = 0
    followingTotal.value = 0
  }
}

async function fetchFollowingState() {
  if (!authStore.isLoggedIn || !authStore.user?.id || isOwnProfile.value) {
    following.value = false
    return
  }
  try {
    const res = await getFollowing(authStore.user.id, 1, 100)
    const ids = res.data.data.following?.map((item) => item.following_id) || []
    following.value = ids.includes(profileId.value)
  } catch (e) {
    following.value = false
  }
}

async function toggleFollow() {
  if (!canFollow.value) return
  followLoading.value = true
  try {
    if (following.value) {
      await unfollowUser(profileId.value)
      following.value = false
    } else {
      await followUser(profileId.value)
      following.value = true
    }
    fetchStats()
  } catch (e) {
    profileError.value = e.response?.data?.message || '关注操作失败'
  } finally {
    followLoading.value = false
  }
}

async function loadProfile() {
  if (!profileId.value) return
  await Promise.all([fetchProfile(), fetchPosts(), fetchStats(), fetchFollowingState()])
}

onMounted(loadProfile)

watch(
  () => route.params.id,
  () => {
    loadProfile()
  }
)
</script>

<template>
  <div>
    <div class="card profile-card" style="margin-bottom:20px">
      <div class="profile-header">
        <span class="avatar avatar-lg">
          <img v-if="profile?.avatar_url" :src="profile.avatar_url" alt="avatar" class="avatar-img" />
          <span v-else>{{ profile?.username?.[0]?.toUpperCase() || '?' }}</span>
        </span>
        <div class="profile-info">
          <div class="profile-name">{{ profile?.username || '未知用户' }}</div>
          <div class="profile-email">{{ profile?.email }}</div>
          <div class="profile-stats">
            <span>关注 {{ followingTotal }}</span>
            <span>粉丝 {{ followersTotal }}</span>
          </div>
        </div>
        <button v-if="canFollow" class="btn btn-primary" :disabled="followLoading" @click="toggleFollow">
          {{ following ? '已关注' : '关注' }}
        </button>
      </div>
      <p v-if="profileError" class="error-msg" style="margin-top:12px">{{ profileError }}</p>
    </div>

    <h3 style="margin-bottom:12px;font-size:1rem;color:#555">动态</h3>

    <div v-if="loading" class="empty-state">加载中…</div>
    <div v-else-if="postsError" class="empty-state" style="color:#e0245e">{{ postsError }}</div>
    <div v-else-if="posts.length === 0" class="empty-state">暂无动态</div>
    <template v-else>
      <PostCard v-for="post in posts" :key="post.id" :post="post" :show-follow="authStore.isLoggedIn" />
    </template>
  </div>
</template>

<style scoped>
.profile-card { display: flex; flex-direction: column; gap: 12px; }
.profile-header { display: flex; align-items: center; gap: 14px; flex-wrap: wrap; }
.profile-info { display: flex; flex-direction: column; gap: 6px; flex: 1; }
.profile-name { font-size: 1.2rem; font-weight: 700; }
.profile-email { font-size: .85rem; color: #888; }
.profile-stats { display: flex; gap: 12px; font-size: .85rem; color: #666; }
.avatar-lg { width: 52px; height: 52px; font-size: 1.2rem; border-radius: 50%; background: #1da1f2; color: #fff; display: inline-flex; align-items: center; justify-content: center; font-weight: 700; }
.avatar-img { width: 100%; height: 100%; object-fit: cover; border-radius: 50%; }
</style>

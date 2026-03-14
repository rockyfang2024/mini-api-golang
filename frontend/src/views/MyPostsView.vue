<script setup>
import { computed, ref, onMounted } from 'vue'
import { getUserPosts, getFollowers, getFollowing, uploadAvatar, getMe } from '../api/index.js'
import { authStore } from '../stores/auth'
import PostCard from '../components/PostCard.vue'

const posts = ref([])
const loading = ref(false)
const error = ref('')
const followers = ref([])
const following = ref([])
const followersTotal = ref(0)
const followingTotal = ref(0)
const followersError = ref('')
const followingError = ref('')
const avatarError = ref('')
const avatarSuccess = ref('')
const avatarUploading = ref(false)

const avatarUrl = computed(() => authStore.user?.avatar_url || '')

async function ensureUser() {
  if (authStore.user?.id) {
    return authStore.user
  }
  try {
    const res = await getMe()
    authStore.setUser(res.data.data)
    return res.data.data
  } catch (e) {
    error.value = '用户信息获取失败，请重新登录'
    return null
  }
}

async function fetchMyPosts(userId) {
  loading.value = true
  error.value = ''
  try {
    const res = await getUserPosts(userId)
    posts.value = res.data.data || []
  } catch (e) {
    error.value = '加载失败，请刷新重试'
  } finally {
    loading.value = false
  }
}

async function fetchFollowersList(userId) {
  followersError.value = ''
  try {
    const res = await getFollowers(userId, 1, 20)
    followers.value = res.data.data.followers || []
    followersTotal.value = res.data.data.total || 0
  } catch (e) {
    followersError.value = '加载粉丝列表失败'
  }
}

async function fetchFollowingList(userId) {
  followingError.value = ''
  try {
    const res = await getFollowing(userId, 1, 20)
    following.value = res.data.data.following || []
    followingTotal.value = res.data.data.total || 0
  } catch (e) {
    followingError.value = '加载关注列表失败'
  }
}

async function handleAvatarUpload(event) {
  const file = event.target.files?.[0]
  if (!file) {
    return
  }
  avatarError.value = ''
  avatarSuccess.value = ''
  avatarUploading.value = true
  try {
    const res = await uploadAvatar(file)
    authStore.setUser({ ...authStore.user, avatar_url: res.data.data.avatar_url })
    avatarSuccess.value = '头像已更新'
  } catch (e) {
    avatarError.value = e.response?.data?.message || '头像上传失败'
  } finally {
    avatarUploading.value = false
    event.target.value = ''
  }
}

onMounted(async () => {
  const user = await ensureUser()
  if (!user?.id) {
    return
  }
  fetchMyPosts(user.id)
  fetchFollowersList(user.id)
  fetchFollowingList(user.id)
})
</script>

<template>
  <div>
    <div class="card profile-card" style="margin-bottom:20px">
      <div class="profile-header">
        <span class="avatar avatar-lg">
          <img v-if="avatarUrl" :src="avatarUrl" alt="avatar" class="avatar-img" />
          <span v-else>{{ authStore.user?.username?.[0]?.toUpperCase() }}</span>
        </span>
        <div>
          <div style="font-size:1.2rem;font-weight:700">{{ authStore.user?.username }}</div>
          <div style="font-size:.85rem;color:#888">{{ authStore.user?.email }}</div>
        </div>
      </div>
      <div class="avatar-upload">
        <label class="btn btn-outline">
          {{ avatarUploading ? '上传中…' : '上传头像' }}
          <input class="avatar-input" type="file" accept="image/*" @change="handleAvatarUpload" />
        </label>
      </div>
      <p v-if="avatarError" class="error-msg" style="margin-top:8px">{{ avatarError }}</p>
      <p v-if="avatarSuccess" class="success-msg" style="margin-top:8px">{{ avatarSuccess }}</p>
    </div>

    <div class="follow-grid">
      <div class="card follow-card">
        <div class="follow-title">关注 <span class="count">{{ followingTotal }}</span></div>
        <div v-if="followingError" class="empty-state" style="color:#e0245e">{{ followingError }}</div>
        <div v-else-if="following.length === 0" class="empty-state">暂无关注</div>
        <ul v-else class="follow-list">
          <li v-for="item in following" :key="item.id">
            <router-link :to="`/users/${item.following_id}`" class="follow-user">
              <span class="avatar avatar-sm">
                {{ item.following?.username?.[0]?.toUpperCase() || '?' }}
              </span>
              <span>{{ item.following?.username || '未知用户' }}</span>
            </router-link>
          </li>
        </ul>
      </div>
      <div class="card follow-card">
        <div class="follow-title">粉丝 <span class="count">{{ followersTotal }}</span></div>
        <div v-if="followersError" class="empty-state" style="color:#e0245e">{{ followersError }}</div>
        <div v-else-if="followers.length === 0" class="empty-state">暂无粉丝</div>
        <ul v-else class="follow-list">
          <li v-for="item in followers" :key="item.id">
            <router-link :to="`/users/${item.follower_id}`" class="follow-user">
              <span class="avatar avatar-sm">
                {{ item.follower?.username?.[0]?.toUpperCase() || '?' }}
              </span>
              <span>{{ item.follower?.username || '未知用户' }}</span>
            </router-link>
          </li>
        </ul>
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
.profile-card { display: flex; flex-direction: column; gap: 12px; }
.profile-header { display: flex; align-items: center; gap: 14px; }
.avatar-lg { width: 48px; height: 48px; font-size: 1.2rem; border-radius: 50%; background: #1da1f2; color: #fff; display: inline-flex; align-items: center; justify-content: center; font-weight: 700; }
.avatar-img { width: 100%; height: 100%; object-fit: cover; border-radius: 50%; }
.avatar-upload { display: flex; align-items: center; gap: 10px; }
.btn-outline { border: 1px solid #1da1f2; color: #1da1f2; background: transparent; }
.btn-outline:hover { background: #e8f5fd; }
.avatar-input { display: none; }
.follow-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; margin-bottom: 20px; }
.follow-card { padding: 16px; }
.follow-title { font-weight: 600; margin-bottom: 12px; display: flex; align-items: center; gap: 8px; }
.count { font-size: .85rem; color: #888; }
.follow-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 10px; }
.follow-list li { display: flex; align-items: center; gap: 8px; font-size: .9rem; }
.follow-user { display: inline-flex; align-items: center; gap: 8px; text-decoration: none; color: #333; }
.follow-user:hover { color: #1da1f2; }
.avatar-sm { width: 26px; height: 26px; font-size: .75rem; border-radius: 50%; background: #1da1f2; color: #fff; display: inline-flex; align-items: center; justify-content: center; font-weight: 700; }
</style>

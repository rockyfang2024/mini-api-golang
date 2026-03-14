<script setup>
import { ref, onMounted } from 'vue'
import { getNotifications, markNotificationRead, markAllNotificationsRead } from '../api/index.js'

const notifications = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const error = ref('')
const actionError = ref('')

async function fetchNotifications({ reset = false } = {}) {
  if (loading.value) return
  loading.value = true
  error.value = ''
  actionError.value = ''
  if (reset) {
    page.value = 1
    notifications.value = []
  }
  try {
    const res = await getNotifications(page.value, pageSize)
    const items = res.data.data.notifications || []
    notifications.value = page.value === 1 ? items : [...notifications.value, ...items]
    total.value = res.data.data.total || 0
  } catch (e) {
    error.value = '加载通知失败，请刷新重试'
  } finally {
    loading.value = false
  }
}

function formatTime(ts) {
  return new Date(ts).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}

function notificationText(notification) {
  const actor = notification.actor?.username || '有人'
  switch (notification.type) {
    case 'like':
      return `${actor} 点赞了你的动态`
    case 'repost':
      return `${actor} 转发了你的动态`
    case 'follow':
      return `${actor} 关注了你`
    case 'new_post':
      return `${actor} 发布了新动态`
    default:
      return `${actor} 有新动态`
  }
}

async function handleMarkRead(notificationId) {
  actionError.value = ''
  try {
    await markNotificationRead(notificationId)
    notifications.value = notifications.value.map((item) =>
      item.id === notificationId ? { ...item, is_read: true } : item
    )
  } catch (e) {
    actionError.value = e.response?.data?.message || '标记已读失败'
  }
}

async function handleMarkAllRead() {
  actionError.value = ''
  try {
    await markAllNotificationsRead()
    notifications.value = notifications.value.map((item) => ({ ...item, is_read: true }))
  } catch (e) {
    actionError.value = e.response?.data?.message || '操作失败，请重试'
  }
}

function loadMore() {
  if (notifications.value.length >= total.value) {
    return
  }
  page.value += 1
  fetchNotifications()
}

onMounted(() => fetchNotifications({ reset: true }))
</script>

<template>
  <div>
    <div class="card notifications-header">
      <div>
        <div style="font-size:1.2rem;font-weight:700">通知中心</div>
        <div style="font-size:.85rem;color:#888">共 {{ total }} 条通知</div>
      </div>
      <button class="btn btn-primary" :disabled="loading || notifications.length === 0" @click="handleMarkAllRead">
        全部已读
      </button>
    </div>

    <p v-if="actionError" class="error-msg" style="margin-bottom:12px">{{ actionError }}</p>

    <div v-if="loading && notifications.length === 0" class="empty-state">加载中…</div>
    <div v-else-if="error" class="empty-state" style="color:#e0245e">{{ error }}</div>
    <div v-else-if="notifications.length === 0" class="empty-state">暂无通知</div>

    <div v-else class="notification-list">
      <div v-for="notification in notifications" :key="notification.id" class="card notification-item">
        <div class="notification-main">
          <div class="notification-text">
            <span class="status-dot" :class="{ unread: !notification.is_read }"></span>
            <span>{{ notificationText(notification) }}</span>
          </div>
          <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
        </div>
        <div v-if="notification.post?.content" class="notification-post">
          {{ notification.post.content }}
        </div>
        <button
          v-if="!notification.is_read"
          class="btn btn-outline"
          @click="handleMarkRead(notification.id)"
        >
          标记已读
        </button>
      </div>
    </div>

    <div v-if="notifications.length < total" class="load-more">
      <button class="btn btn-primary" :disabled="loading" @click="loadMore">
        {{ loading ? '加载中…' : '加载更多' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.notifications-header { display: flex; align-items: center; justify-content: space-between; }
.notification-list { display: flex; flex-direction: column; gap: 12px; }
.notification-item { display: flex; flex-direction: column; gap: 8px; }
.notification-main { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.notification-text { display: flex; align-items: center; gap: 8px; font-weight: 600; }
.notification-time { font-size: .82rem; color: #888; }
.notification-post { font-size: .9rem; color: #555; padding: 8px 12px; background: #f6f8fa; border-radius: 8px; }
.status-dot { width: 8px; height: 8px; border-radius: 50%; background: #ccc; display: inline-block; }
.status-dot.unread { background: #1da1f2; }
.btn-outline { border: 1px solid #1da1f2; color: #1da1f2; background: transparent; width: fit-content; }
.btn-outline:hover { background: #e8f5fd; }
.load-more { display: flex; justify-content: center; margin-top: 12px; }
</style>

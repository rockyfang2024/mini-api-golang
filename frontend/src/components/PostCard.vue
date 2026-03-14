<script setup>
const props = defineProps({
  post: {
    type: Object,
    required: true,
  },
})

function formatTime(ts) {
  return new Date(ts).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}
</script>

<template>
  <div class="post-card">
    <div class="post-header">
      <span class="avatar">{{ post.author?.username?.[0]?.toUpperCase() || '?' }}</span>
      <div>
        <div class="post-author">{{ post.author?.username || '未知用户' }}</div>
        <div class="post-time">{{ formatTime(post.created_at) }}</div>
      </div>
      <span class="post-badge" :class="post.visibility">
        {{ post.visibility === 'private' ? '🔒 私密' : '🌐 公开' }}
      </span>
    </div>
    <div class="post-content">{{ post.content }}</div>
  </div>
</template>

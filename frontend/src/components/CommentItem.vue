<script setup>
import { computed } from 'vue'
import { authStore } from '../stores/auth'

defineOptions({ name: 'CommentItem' })

const props = defineProps({
  comment: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['reply'])

const authorId = computed(() => props.comment.author_id || props.comment.author?.id)
const authorLink = computed(() => (authorId.value ? `/users/${authorId.value}` : ''))
const avatarUrl = computed(() => props.comment.author?.avatar_url || '')

function formatTime(ts) {
  return new Date(ts).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}

function handleReply() {
  emit('reply', props.comment)
}
</script>

<template>
  <div class="comment-item">
    <div class="comment-header">
      <router-link v-if="authorLink" :to="authorLink" class="comment-author">
        <span class="avatar avatar-xs">
          <img v-if="avatarUrl" :src="avatarUrl" alt="avatar" />
          <span v-else>{{ comment.author?.username?.[0]?.toUpperCase() || '?' }}</span>
        </span>
        {{ comment.author?.username || '未知用户' }}
      </router-link>
      <span v-else class="comment-author">
        <span class="avatar avatar-xs">
          <span>{{ comment.author?.username?.[0]?.toUpperCase() || '?' }}</span>
        </span>
        {{ comment.author?.username || '未知用户' }}
      </span>
      <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
    </div>
    <div class="comment-content">{{ comment.content }}</div>
    <button v-if="authStore.isLoggedIn" class="comment-reply" type="button" @click="handleReply">
      回复
    </button>
    <div v-if="comment.children?.length" class="comment-children">
      <CommentItem
        v-for="child in comment.children"
        :key="child.id"
        :comment="child"
        @reply="emit('reply', $event)"
      />
    </div>
  </div>
</template>

<style scoped>
.comment-item { padding: 12px 0; border-bottom: 1px solid #f0f0f0; }
.comment-header { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.comment-author { font-weight: 600; font-size: .9rem; color: #333; text-decoration: none; display: inline-flex; align-items: center; gap: 6px; }
.comment-time { font-size: .78rem; color: #999; }
.comment-content { margin: 6px 0; font-size: .9rem; color: #444; white-space: pre-wrap; }
.comment-reply { background: none; border: none; color: #1da1f2; cursor: pointer; font-size: .82rem; padding: 0; }
.comment-children { margin-left: 32px; }
.avatar-xs { width: 22px; height: 22px; border-radius: 50%; background: #1da1f2; color: #fff; display: inline-flex; align-items: center; justify-content: center; font-size: .7rem; font-weight: 700; overflow: hidden; }
.avatar-xs img { width: 100%; height: 100%; object-fit: cover; }
</style>

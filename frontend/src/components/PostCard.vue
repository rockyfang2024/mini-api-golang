<script setup>
import { computed, ref, watch } from 'vue'
import { authStore } from '../stores/auth'
import { likePost, unlikePost, repostPost, followUser, unfollowUser, getComments, createComment, replyComment } from '../api/index.js'
import CommentItem from './CommentItem.vue'

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
const commentsOpen = ref(false)
const commentsLoading = ref(false)
const commentError = ref('')
const comments = ref([])
const commentsLoaded = ref(false)
const commentContent = ref('')
const commentSubmitting = ref(false)
const replyTarget = ref(null)

const authorId = computed(() => props.post.author_id || props.post.author?.id)
const isOwnPost = computed(() => authStore.user?.id && authStore.user.id === authorId.value)
const canFollow = computed(() => props.showFollow && authStore.isLoggedIn && authorId.value && !isOwnPost.value)
const avatarUrl = computed(() => props.post.author?.avatar_url || '')
const authorLink = computed(() => (authorId.value ? `/users/${authorId.value}` : ''))
const commentCount = computed(() => comments.value.length)
const commentLabel = computed(() => (commentsLoaded.value ? commentCount.value : '评论'))
const commentTree = computed(() => buildCommentTree(comments.value))

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

function toggleComments() {
  commentsOpen.value = !commentsOpen.value
  if (commentsOpen.value && comments.value.length === 0) {
    fetchComments()
  }
}

async function fetchComments() {
  commentsLoading.value = true
  commentError.value = ''
  try {
    const res = await getComments(props.post.id)
    comments.value = res.data.data || []
    commentsLoaded.value = true
  } catch (e) {
    commentError.value = e.response?.data?.message || '加载评论失败'
  } finally {
    commentsLoading.value = false
  }
}

async function submitComment() {
  if (!authStore.isLoggedIn) return
  commentError.value = ''
  const content = commentContent.value.trim()
  if (!content) {
    commentError.value = '评论内容不能为空'
    return
  }
  commentSubmitting.value = true
  try {
    if (replyTarget.value) {
      await replyComment(replyTarget.value.id, content)
    } else {
      await createComment(props.post.id, content)
    }
    commentContent.value = ''
    replyTarget.value = null
    await fetchComments()
  } catch (e) {
    commentError.value = e.response?.data?.message || '评论失败，请重试'
  } finally {
    commentSubmitting.value = false
  }
}

function setReplyTarget(comment) {
  replyTarget.value = comment
}

function clearReplyTarget() {
  replyTarget.value = null
}

function buildCommentTree(list) {
  const map = new Map()
  list.forEach((comment) => {
    map.set(comment.id, { ...comment, children: [] })
  })
  const roots = []
  map.forEach((comment) => {
    if (comment.parent_comment_id) {
      const parent = map.get(comment.parent_comment_id)
      if (parent) {
        parent.children.push(comment)
        return
      }
    }
    roots.push(comment)
  })
  return roots
}
</script>

<template>
  <div class="post-card">
    <div class="post-header">
      <router-link v-if="authorLink" :to="authorLink" class="avatar">
        <img v-if="avatarUrl" :src="avatarUrl" alt="avatar" class="avatar-img" />
        <span v-else>{{ post.author?.username?.[0]?.toUpperCase() || '?' }}</span>
      </router-link>
      <span v-else class="avatar">
        <span>{{ post.author?.username?.[0]?.toUpperCase() || '?' }}</span>
      </span>
      <div>
        <router-link v-if="authorLink" :to="authorLink" class="post-author">
          {{ post.author?.username || '未知用户' }}
        </router-link>
        <div v-else class="post-author">{{ post.author?.username || '未知用户' }}</div>
        <div class="post-time">{{ formatTime(post.created_at) }}</div>
      </div>
      <span class="post-badge" :class="post.visibility">
        {{ post.visibility === 'private' ? '🔒 私密' : '🌐 公开' }}
      </span>
    </div>
    <div class="post-content">{{ post.content }}</div>
    <div v-if="post.images?.length" class="post-images">
      <img v-for="image in post.images" :key="image.id || image.url" :src="image.url" alt="post image" />
    </div>
    <div class="post-actions">
      <button class="action-btn" :class="{ active: isLiked }" :disabled="!authStore.isLoggedIn" @click="toggleLike">
        👍 {{ likeCount }}
      </button>
      <button class="action-btn" :class="{ active: isReposted }" :disabled="!authStore.isLoggedIn || isReposted" @click="handleRepost">
        🔁 {{ repostCount }}
      </button>
      <button class="action-btn" @click="toggleComments">
        💬 {{ commentLabel }}
      </button>
      <button v-if="canFollow" class="action-btn follow-btn" @click="toggleFollow">
        {{ following ? '已关注' : '关注' }}
      </button>
    </div>
    <p v-if="actionError" class="error-msg" style="margin-top:8px">{{ actionError }}</p>
    <div v-if="commentsOpen" class="comment-section">
      <div v-if="commentsLoading" class="comment-empty">加载评论中…</div>
      <div v-else-if="commentError" class="comment-empty" style="color:#e0245e">{{ commentError }}</div>
      <div v-else-if="commentTree.length === 0" class="comment-empty">暂无评论</div>
      <div v-else class="comment-list">
        <CommentItem v-for="comment in commentTree" :key="comment.id" :comment="comment" @reply="setReplyTarget" />
      </div>
      <div class="comment-form">
        <textarea
          v-model="commentContent"
          rows="2"
          placeholder="写下你的评论…"
          :disabled="!authStore.isLoggedIn || commentSubmitting"
        ></textarea>
        <div class="comment-form-footer">
          <div v-if="replyTarget" class="replying">
            回复 {{ replyTarget.author?.username || '用户' }}
            <button type="button" @click="clearReplyTarget">取消</button>
          </div>
          <button class="btn btn-primary" :disabled="!authStore.isLoggedIn || commentSubmitting" @click="submitComment">
            {{ commentSubmitting ? '发送中…' : '发送' }}
          </button>
        </div>
      </div>
    </div>
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
.post-author { font-weight: 600; font-size: .95rem; color: #333; text-decoration: none; }
.post-images { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 8px; margin-top: 10px; }
.post-images img { width: 100%; height: 110px; object-fit: cover; border-radius: 8px; }
.comment-section { margin-top: 12px; border-top: 1px solid #f0f0f0; padding-top: 12px; }
.comment-empty { text-align: center; color: #888; font-size: .85rem; padding: 12px 0; }
.comment-form textarea { width: 100%; border: 1px solid #ddd; border-radius: 8px; padding: 8px 10px; font-size: .9rem; resize: vertical; }
.comment-form-footer { display: flex; align-items: center; justify-content: space-between; margin-top: 8px; gap: 8px; }
.replying { font-size: .8rem; color: #555; display: flex; align-items: center; gap: 6px; }
.replying button { border: none; background: none; color: #1da1f2; cursor: pointer; padding: 0; font-size: .8rem; }
</style>

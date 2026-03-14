<script setup>
import { ref, onMounted } from 'vue'
import { getSettings, updateSettings } from '../api/index.js'

const loading = ref(false)
const saving = ref(false)
const error = ref('')
const success = ref('')
const form = ref({
  allow_comments: true,
  allow_follow: true,
  only_followers_can_view: false,
  only_following_can_view: false,
})

async function fetchSettings() {
  loading.value = true
  error.value = ''
  try {
    const res = await getSettings()
    form.value = {
      allow_comments: res.data.data.allow_comments,
      allow_follow: res.data.data.allow_follow,
      only_followers_can_view: res.data.data.only_followers_can_view,
      only_following_can_view: res.data.data.only_following_can_view,
    }
  } catch (e) {
    error.value = '设置加载失败'
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  saving.value = true
  error.value = ''
  success.value = ''
  try {
    const res = await updateSettings(form.value)
    form.value = {
      allow_comments: res.data.data.allow_comments,
      allow_follow: res.data.data.allow_follow,
      only_followers_can_view: res.data.data.only_followers_can_view,
      only_following_can_view: res.data.data.only_following_can_view,
    }
    success.value = '设置已更新'
  } catch (e) {
    error.value = e.response?.data?.message || '设置更新失败'
  } finally {
    saving.value = false
  }
}

onMounted(fetchSettings)
</script>

<template>
  <div class="card settings-card">
    <h3 class="settings-title">个人设置</h3>
    <div v-if="loading" class="empty-state">加载中…</div>
    <div v-else>
      <div class="setting-item">
        <label>
          <input v-model="form.allow_comments" type="checkbox" />
          允许他人评论我的动态
        </label>
      </div>
      <div class="setting-item">
        <label>
          <input v-model="form.allow_follow" type="checkbox" />
          允许他人关注我
        </label>
      </div>
      <div class="setting-item">
        <label>
          <input v-model="form.only_followers_can_view" type="checkbox" />
          只有关注我的人可以查看动态
        </label>
      </div>
      <div class="setting-item">
        <label>
          <input v-model="form.only_following_can_view" type="checkbox" />
          只有我关注的人可以查看动态
        </label>
      </div>
      <p class="setting-hint">如同时开启两项可见性限制，将只允许双向关注的用户查看。</p>
      <button class="btn btn-primary" :disabled="saving" @click="saveSettings">
        {{ saving ? '保存中…' : '保存设置' }}
      </button>
      <p v-if="error" class="error-msg" style="margin-top:10px">{{ error }}</p>
      <p v-if="success" class="success-msg" style="margin-top:10px">{{ success }}</p>
    </div>
  </div>
</template>

<style scoped>
.settings-card { display: flex; flex-direction: column; gap: 12px; }
.settings-title { font-size: 1.1rem; margin-bottom: 4px; }
.setting-item { font-size: .95rem; color: #333; display: flex; align-items: center; gap: 8px; }
.setting-item input { margin-right: 6px; }
.setting-hint { font-size: .82rem; color: #888; }
</style>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <header class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-3xl font-black uppercase tracking-tight text-white sm:text-4xl">
          Настройки сайта
        </h1>
        <p class="mt-3 text-sm text-zinc-400 max-w-2xl leading-relaxed">
          Управляйте глобальной конфигурацией, структурой отображаемых разделов и составом игровых категорий.
        </p>
      </div>
    </header>

    <!-- Feedback Message -->
    <div
      v-if="feedback"
      class="rounded-[0.95rem] border px-5 py-4 text-sm font-medium shadow-xl transition-all duration-300"
      :class="feedbackTone === 'error'
        ? 'border-orange-500/25 bg-orange-500/10 text-orange-200'
        : 'border-emerald-500/25 bg-emerald-500/10 text-emerald-200'"
    >
      <div class="flex items-center gap-3">
        <span class="text-lg">{{ feedbackTone === 'error' ? '⚠️' : '✨' }}</span>
        <span>{{ feedback }}</span>
      </div>
    </div>

    <!-- Main settings card -->
    <div class="rounded-3xl border border-cyan-400/10 bg-[#07141d]/40 p-8 backdrop-blur-md shadow-2xl">
      <div class="flex items-start gap-4 mb-8">
        <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-cyan-400/10 border border-cyan-400/20 text-xl text-cyan-300">
          ⚙️
        </div>
        <div>
          <h2 class="text-lg font-bold text-white uppercase tracking-wider">
            Отображение вкладок на странице «Игровая»
          </h2>
          <p class="text-xs text-cyan-100/50 mt-1">
            Настройте категории (теги зон), которые отображаются на странице Игровой в качестве главных табов и определяют доступные конфигурации.
          </p>
        </div>
      </div>

      <div class="grid gap-8 lg:grid-cols-12 mt-6">
        <!-- Available tags list (Left) -->
        <div class="lg:col-span-5 space-y-4">
          <div class="rounded-2xl border border-white/5 bg-[#030d14]/60 p-5">
            <h3 class="text-xs font-black uppercase tracking-widest text-zinc-400 mb-4 flex items-center justify-between">
              <span>Все доступные теги</span>
              <span class="rounded bg-zinc-800 px-2 py-0.5 text-[10px] text-zinc-300 font-mono">{{ availableTags.length }}</span>
            </h3>
            
            <div
              v-if="availableTags.length === 0"
              class="text-xs text-zinc-500 py-6 text-center"
            >
              Нет доступных тегов в базе данных.
            </div>
            
            <div v-else class="space-y-2 max-h-[350px] overflow-y-auto pr-1.5 scrollbar-thin">
              <div
                v-for="tag in availableTags"
                :key="tag.id"
                class="flex items-center justify-between rounded-xl border border-white/5 bg-white/[0.02] p-3 transition hover:border-cyan-400/20 hover:bg-white/[0.04]"
              >
                <div>
                  <p class="text-sm font-bold text-white">{{ tag.name }}</p>
                  <p class="text-[10px] text-zinc-500 mt-0.5">ID: {{ tag.id }}</p>
                </div>
                <button
                  type="button"
                  class="rounded-lg bg-cyan-400/10 border border-cyan-400/20 px-2.5 py-1.5 text-[10px] font-black uppercase text-cyan-300 tracking-wider hover:bg-cyan-300 hover:text-black transition-all"
                  :disabled="isActiveTag(tag.id)"
                  :class="isActiveTag(tag.id) ? 'opacity-30 cursor-not-allowed border-transparent bg-zinc-800' : ''"
                  @click="addTag(tag.id)"
                >
                  {{ isActiveTag(tag.id) ? 'Добавлен' : '+ Добавить' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Reorderable active list (Right) -->
        <div class="lg:col-span-7 space-y-4">
          <div class="rounded-2xl border border-white/5 bg-[#030d14]/60 p-5">
            <h3 class="text-xs font-black uppercase tracking-widest text-zinc-400 mb-4 flex items-center justify-between">
              <span>Отображаемые табы (порядок вывода)</span>
              <span class="rounded bg-cyan-300/10 border border-cyan-300/20 px-2 py-0.5 text-[10px] text-cyan-300 font-mono">{{ activeTagIds.length }}</span>
            </h3>

            <div
              v-if="activeTagIds.length === 0"
              class="text-xs text-zinc-500 py-12 text-center border border-dashed border-white/10 rounded-xl"
            >
              Список пуст. Выберите теги слева, чтобы они отображались на странице «Игровая» в заданном порядке.
            </div>

            <div v-else class="space-y-2 max-h-[350px] overflow-y-auto pr-1.5 scrollbar-thin">
              <div
                v-for="(tagId, index) in activeTagIds"
                :key="tagId"
                class="flex items-center justify-between rounded-xl border border-cyan-300/10 bg-cyan-300/[0.02] p-3 shadow-sm hover:border-cyan-300/25 transition-all"
              >
                <div class="flex items-center gap-3">
                  <span class="text-xs font-black text-cyan-300/50 font-mono w-4">#{{ index + 1 }}</span>
                  <div>
                    <p class="text-sm font-bold text-white">{{ getTagName(tagId) }}</p>
                    <p class="text-[10px] text-zinc-500 mt-0.5">ID тега: {{ tagId }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-1.5">
                  <!-- Reorder buttons -->
                  <button
                    type="button"
                    class="rounded-lg p-1.5 hover:bg-white/5 text-zinc-400 hover:text-white transition disabled:opacity-20"
                    :disabled="index === 0"
                    @click="moveUp(index)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M5 15l7-7 7 7" />
                    </svg>
                  </button>
                  <button
                    type="button"
                    class="rounded-lg p-1.5 hover:bg-white/5 text-zinc-400 hover:text-white transition disabled:opacity-20"
                    :disabled="index === activeTagIds.length - 1"
                    @click="moveDown(index)"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                    </svg>
                  </button>

                  <div class="h-6 w-[1px] bg-white/5 mx-1" />

                  <!-- Delete button -->
                  <button
                    type="button"
                    class="rounded-lg p-1.5 hover:bg-red-500/10 text-zinc-400 hover:text-red-400 transition"
                    @click="removeTag(index)"
                    title="Убрать из списка"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Action buttons -->
      <div class="mt-8 flex justify-end border-t border-white/5 pt-6">
        <button
          type="button"
          class="rounded-[0.95rem] bg-cyan-300 px-8 py-3.5 text-xs font-black uppercase tracking-wider text-[#020c13] transition hover:bg-cyan-200 hover:shadow-[0_0_24px_rgba(34,211,238,0.45)] active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="saving"
          @click="saveSettings"
        >
          {{ saving ? 'Сохранение...' : 'Сохранить конфигурацию' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getAdminZoneTags, getAdminSiteSettings, updateAdminSiteSettings } from '~/api/admin'
import type { AdminZoneTag } from '~/api/types'

useHead({
  title: 'Настройки сайта - PlayGround'
})

const availableTags = ref<AdminZoneTag[]>([])
const activeTagIds = ref<number[]>([])

const saving = ref(false)
const feedback = ref('')
const feedbackTone = ref<'success' | 'error'>('success')

function showFeedback(message: string, tone: 'success' | 'error' = 'success') {
  feedback.value = message
  feedbackTone.value = tone
  setTimeout(() => {
    feedback.value = ''
  }, 4000)
}

async function loadData() {
  try {
    const tagsResponse = await getAdminZoneTags()
    availableTags.value = tagsResponse.zone_tags ?? []

    const settingsResponse = await getAdminSiteSettings()
    // Load config from settings_json
    const loadedIds = settingsResponse.settings?.settings_json
    activeTagIds.value = Array.isArray(loadedIds) ? loadedIds : []
  } catch (err: any) {
    showFeedback('Не удалось загрузить текущие настройки сайта.', 'error')
  }
}

function getTagName(tagId: number): string {
  const found = availableTags.value.find(t => t.id === tagId)
  return found ? found.name : `Тег #${tagId}`
}

function isActiveTag(tagId: number): boolean {
  return activeTagIds.value.includes(tagId)
}

function addTag(tagId: number) {
  if (!isActiveTag(tagId)) {
    activeTagIds.value.push(tagId)
  }
}

function removeTag(index: number) {
  activeTagIds.value.splice(index, 1)
}

function moveUp(index: number) {
  if (index > 0) {
    const temp = activeTagIds.value[index]!
    activeTagIds.value[index] = activeTagIds.value[index - 1]!
    activeTagIds.value[index - 1] = temp
  }
}

function moveDown(index: number) {
  if (index < activeTagIds.value.length - 1) {
    const temp = activeTagIds.value[index]!
    activeTagIds.value[index] = activeTagIds.value[index + 1]!
    activeTagIds.value[index + 1] = temp
  }
}

async function saveSettings() {
  saving.value = true
  try {
    await updateAdminSiteSettings({
      settings_json: activeTagIds.value,
      gallery_items_json: []
    })
    showFeedback('Глобальная конфигурация сайта успешно сохранена!', 'success')
  } catch (err: any) {
    showFeedback('Произошла ошибка при сохранении настроек.', 'error')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

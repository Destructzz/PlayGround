<template>
  <div class="min-h-screen bg-[#020c13] pb-12 pt-20 text-white">
    <div class="mx-auto max-w-[1400px] px-4 sm:px-6 lg:px-8">

      <!-- Header -->
      <div class="mb-8 flex items-center justify-between">
        <div>
          <NuxtLink to="/admin" class="text-[11px] font-bold uppercase tracking-[0.3em] text-fuchsia-100/40 hover:text-fuchsia-100/70 transition">
            &larr; Admin Panel
          </NuxtLink>
          <h1 class="mt-2 text-3xl font-black tracking-tight text-white">Events</h1>
          <p class="mt-1 text-sm text-zinc-400">Управление событиями, форматами и спикерами.</p>
        </div>
        <button
          type="button"
          class="rounded-[0.9rem] bg-fuchsia-300 px-6 py-3 text-sm font-black uppercase tracking-widest text-[#020c13] transition hover:bg-fuchsia-200 hover:shadow-[0_0_20px_rgba(244,114,182,0.4)]"
          @click="showCreate = !showCreate"
        >
          {{ showCreate ? 'Скрыть форму' : '+ Добавить событие' }}
        </button>
      </div>

      <!-- Feedback -->
      <div
        v-if="feedbackMessage"
        class="mb-6 rounded-[0.8rem] border px-4 py-3 text-sm"
        :class="feedbackTone === 'error' ? 'border-orange-300/30 bg-orange-500/10 text-orange-100' : 'border-emerald-300/30 bg-emerald-500/10 text-emerald-100'"
      >
        {{ feedbackMessage }}
      </div>

      <!-- Create form -->
      <div v-if="showCreate" class="mb-8 rounded-[1rem] border border-fuchsia-400/20 bg-[#050f17] p-6 shadow-2xl">
        <p class="mb-5 text-xs font-black uppercase tracking-[0.3em] text-fuchsia-100/50">Новое событие</p>
        <form class="space-y-5" @submit.prevent="submitCreate">
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Название</label>
              <input v-model="createForm.name" type="text" placeholder="Night Bracket Qualifier" required
                class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-fuchsia-100/20 focus:border-fuchsia-300 focus:outline-none">
            </div>
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Вместимость</label>
              <input v-model="createForm.capacity" type="number" min="1" required
                class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
            </div>
          </div>

          <div>
            <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Описание</label>
            <textarea v-model="createForm.description" rows="3" placeholder="Коротко об атмосфере..."
              class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-fuchsia-100/20 focus:border-fuchsia-300 focus:outline-none" />
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Формат</label>
              <select v-model="createForm.format"
                class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
                <option value="">Выбрать формат</option>
                <option value="Tournament">Tournament</option>
                <option value="Meetup">Meetup</option>
                <option value="Community event">Community event</option>
                <option value="LAN Party">LAN Party</option>
                <option value="Workshop">Workshop</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Цвет акцента</label>
              <div class="flex items-center gap-3">
                <input v-model="createForm.accentColor" type="color"
                  class="h-[46px] w-14 cursor-pointer rounded-[0.6rem] border border-fuchsia-400/18 bg-[#06131c] p-1">
                <span class="text-sm font-mono text-zinc-400">{{ createForm.accentColor }}</span>
              </div>
            </div>
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Начало события</label>
              <input v-model="createForm.startTime" type="datetime-local"
                class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
            </div>
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Конец события</label>
              <input v-model="createForm.endTime" type="datetime-local"
                class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
            </div>
          </div>

          <div>
            <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Спикеры / хосты (каждый с новой строки)</label>
            <textarea v-model="createForm.speakersRaw" rows="3" placeholder="Host: Raven, Caster: Miko"
              class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-fuchsia-100/20 focus:border-fuchsia-300 focus:outline-none" />
          </div>

          <div class="flex items-center gap-4">
            <label class="flex items-center gap-2 text-sm text-zinc-300">
              <input v-model="createForm.isActive" type="checkbox" class="h-4 w-4 accent-fuchsia-300">
              Активно (видно в афише)
            </label>
          </div>

          <div>
            <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Zone Tag</label>
            <select v-model="createForm.zoneTagId"
              class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
              <option value="">Выбрать tag</option>
              <option v-for="tag in zoneTags" :key="tag.id" :value="String(tag.id)">{{ tag.name }}</option>
            </select>
          </div>

          <button type="submit" :disabled="isMutating"
            class="w-full rounded-[0.9rem] bg-fuchsia-300 py-3.5 text-sm font-black uppercase tracking-widest text-[#020c13] transition hover:bg-fuchsia-200 disabled:cursor-not-allowed disabled:bg-fuchsia-300/50">
            Создать событие
          </button>
        </form>
      </div>

      <!-- Loading -->
      <div v-if="isLoading" class="rounded-[1rem] border border-white/5 bg-[#050f17] p-12 text-center text-sm text-zinc-400">
        Загружаем события...
      </div>

      <!-- Events list -->
      <div v-else class="space-y-4">
        <div v-if="!eventZones.length" class="rounded-[1rem] border border-white/8 bg-[#050f17] p-8 text-center text-sm text-zinc-400">
          Событий ещё нет. Нажми «Добавить событие» выше.
        </div>

        <div
          v-for="zone in eventZones"
          :key="zone.id"
          class="overflow-hidden rounded-[1rem] border border-white/8 bg-[#050f17] shadow-xl"
        >
          <div class="flex items-start justify-between gap-4 px-6 py-5">
            <div class="flex-1 min-w-0">
              <div class="flex flex-wrap items-center gap-2 mb-2">
                <span
                  class="rounded-full border px-3 py-1 text-[10px] font-bold uppercase tracking-[0.25em]"
                  :style="{ borderColor: (zoneDetails(zone).accent ?? '#f472b6') + '40', backgroundColor: (zoneDetails(zone).accent ?? '#f472b6') + '18', color: zoneDetails(zone).accent ?? '#f472b6' }"
                >
                  {{ zoneDetails(zone).format || 'Event' }}
                </span>
                <span v-if="zone.is_active" class="rounded-full border border-emerald-300/20 bg-emerald-400/10 px-2 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-emerald-200">active</span>
                <span v-else class="rounded-full border border-zinc-300/15 bg-zinc-500/10 px-2 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-400">inactive</span>
              </div>
              <h3 class="text-xl font-black text-white">{{ zone.name }}</h3>
              <p class="mt-1 text-sm text-zinc-400">{{ typeof zone.description === 'string' ? zone.description : '' }}</p>
              <!-- Speakers -->
              <div v-if="zoneDetails(zone).speakers?.length" class="mt-3 flex flex-wrap gap-1.5">
                <span v-for="speaker in zoneDetails(zone).speakers" :key="speaker"
                  class="rounded-full border border-white/8 bg-white/4 px-2.5 py-1 text-[10px] text-zinc-400">
                  {{ speaker }}
                </span>
              </div>
            </div>

            <div class="flex items-center gap-2 flex-shrink-0">
              <div class="rounded-[0.7rem] border border-white/8 bg-white/4 px-4 py-2 text-right">
                <p class="text-[10px] uppercase tracking-[0.2em] text-zinc-500">Capacity</p>
                <p class="text-lg font-black text-white">{{ zone.capacity }}</p>
              </div>
              <button type="button"
                class="rounded-[0.8rem] border border-fuchsia-400/18 px-4 py-2.5 text-sm font-bold text-white transition hover:border-fuchsia-300/40"
                @click="toggleEdit(zone.id)">
                {{ editingId === zone.id ? 'Скрыть' : 'Изменить' }}
              </button>
              <button type="button" :disabled="isMutating"
                class="rounded-[0.8rem] border border-orange-300/18 px-4 py-2.5 text-sm font-bold text-orange-100 transition hover:border-orange-200/40 disabled:cursor-not-allowed"
                @click="removeZone(zone.id, zone.name)">
                Удалить
              </button>
            </div>
          </div>

          <!-- Edit form -->
          <div v-if="editingId === zone.id && currentEditForm" class="border-t border-white/8 bg-[#06131c] px-6 py-5">
            <form class="space-y-4" @submit.prevent="submitEdit(zone.id)">
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Название</label>
                  <input v-model="currentEditForm.name" type="text"
                    class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
                </div>
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Вместимость</label>
                  <input v-model="currentEditForm.capacity" type="number"
                    class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
                </div>
              </div>
              <div>
                <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Описание</label>
                <textarea v-model="currentEditForm.description" rows="2"
                  class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none" />
              </div>
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Формат</label>
                  <select v-model="currentEditForm.format"
                    class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
                    <option value="Tournament">Tournament</option>
                    <option value="Meetup">Meetup</option>
                    <option value="Community event">Community event</option>
                    <option value="LAN Party">LAN Party</option>
                    <option value="Workshop">Workshop</option>
                  </select>
                </div>
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Акцент</label>
                  <div class="flex items-center gap-3">
                    <input v-model="currentEditForm.accentColor" type="color"
                      class="h-[46px] w-14 cursor-pointer rounded-[0.6rem] border border-fuchsia-400/18 bg-[#071926] p-1">
                    <span class="text-sm font-mono text-zinc-400">{{ currentEditForm.accentColor }}</span>
                  </div>
                </div>
              </div>
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Начало</label>
                  <input v-model="currentEditForm.startTime" type="datetime-local"
                    class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
                </div>
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Конец</label>
                  <input v-model="currentEditForm.endTime" type="datetime-local"
                    class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none">
                </div>
              </div>
              <div>
                <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-fuchsia-100/45">Спикеры (каждый с новой строки)</label>
                <textarea v-model="currentEditForm.speakersRaw" rows="3"
                  class="w-full rounded-[0.8rem] border border-fuchsia-400/18 bg-[#071926] px-4 py-3 text-white focus:border-fuchsia-300 focus:outline-none" />
              </div>
              <div class="flex items-center gap-4">
                <label class="flex items-center gap-2 text-sm text-zinc-300">
                  <input v-model="currentEditForm.isActive" type="checkbox" class="h-4 w-4 accent-fuchsia-300">
                  Активно
                </label>
              </div>
              <div class="flex gap-3">
                <button type="submit" :disabled="isMutating"
                  class="rounded-[0.85rem] bg-fuchsia-300 px-5 py-3 text-sm font-black text-[#020c13] transition hover:bg-fuchsia-200 disabled:cursor-not-allowed">
                  Сохранить
                </button>
                <button type="button"
                  class="rounded-[0.85rem] border border-white/10 px-5 py-3 text-sm font-bold text-white transition hover:border-white/25"
                  @click="editingId = null">
                  Отмена
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getAdminZones, getAdminZoneTags, patchAdminZone, deleteAdminZone, createAdminZone } from '~/api/admin'
import type { AdminZone, AdminZoneTag } from '~/api/types'

definePageMeta({ middleware: 'admin' })
useHead({ title: 'Admin · Events — PlayGround' })

const isLoading = ref(true)
const isMutating = ref(false)
const feedbackTone = ref<'success' | 'error'>('success')
const feedbackMessage = ref('')
const showCreate = ref(false)
const editingId = ref<number | null>(null)

const allZones = ref<AdminZone[]>([])
const zoneTags = ref<AdminZoneTag[]>([])

interface EventForm {
  name: string
  capacity: string
  description: string
  format: string
  accentColor: string
  speakersRaw: string
  startTime: string
  endTime: string
  isActive: boolean
  zoneTagId: string
}

const createForm = ref<EventForm>(emptyForm())
const currentEditForm = ref<EventForm | null>(null)

const eventZones = computed(() => allZones.value.filter(z => z.zone_type === 'event'))

onMounted(async () => {
  try {
    const [zonesResp, tagsResp] = await Promise.all([getAdminZones(), getAdminZoneTags()])
    allZones.value = zonesResp.zones ?? []
    zoneTags.value = tagsResp.zone_tags ?? []
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка загрузки')
  } finally {
    isLoading.value = false
  }
})

function emptyForm(): EventForm {
  return { name: '', capacity: '32', description: '', format: 'Tournament', accentColor: '#f472b6', speakersRaw: '', startTime: '', endTime: '', isActive: true, zoneTagId: '' }
}

function zoneDetails(zone: AdminZone) {
  try {
    const d = zone.details_json as any
    if (typeof d === 'object' && d !== null) return d
    if (typeof d === 'string') return JSON.parse(d)
  } catch {}
  return {}
}

function buildDetailsJson(form: EventForm): string {
  const speakers = form.speakersRaw.split('\n').map(s => s.trim()).filter(Boolean)
  const obj: Record<string, unknown> = {
    format: form.format,
    accent: form.accentColor,
    speakers
  }
  if (form.startTime) obj.start_time = new Date(form.startTime).toISOString()
  if (form.endTime) obj.end_time = new Date(form.endTime).toISOString()
  return JSON.stringify(obj)
}

function toggleEdit(id: number) {
  if (editingId.value === id) {
    editingId.value = null
    currentEditForm.value = null
    return
  }
  editingId.value = id
  const zone = allZones.value.find(z => z.id === id)
  if (!zone) return
  const d = zoneDetails(zone)
  currentEditForm.value = {
    name: zone.name,
    capacity: String(zone.capacity),
    description: typeof zone.description === 'string' ? zone.description : '',
    format: d.format ?? 'Tournament',
    accentColor: d.accent ?? '#f472b6',
    speakersRaw: Array.isArray(d.speakers) ? d.speakers.join('\n') : '',
    startTime: d.start_time ? new Date(d.start_time).toISOString().slice(0, 16) : '',
    endTime: d.end_time ? new Date(d.end_time).toISOString().slice(0, 16) : '',
    isActive: zone.is_active,
    zoneTagId: String(zone.zone_tag_id)
  }
}

function showFeedback(tone: 'success' | 'error', msg: string) {
  feedbackTone.value = tone
  feedbackMessage.value = msg
  if (tone === 'success') setTimeout(() => { feedbackMessage.value = '' }, 3000)
}

async function submitCreate() {
  if (!createForm.value.zoneTagId) { showFeedback('error', 'Выбери zone tag'); return }
  isMutating.value = true
  try {
    const resp = await createAdminZone({
      name: createForm.value.name,
      type: 'event',
      zone_tag_id: Number(createForm.value.zoneTagId),
      capacity: Number(createForm.value.capacity),
      description: createForm.value.description,
      is_active: createForm.value.isActive,
      details_json: buildDetailsJson(createForm.value)
    })
    allZones.value.push(resp.zone)
    createForm.value = emptyForm()
    showCreate.value = false
    showFeedback('success', `Событие «${resp.zone.name}» создано.`)
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка создания')
  } finally {
    isMutating.value = false
  }
}

async function submitEdit(id: number) {
  const form = currentEditForm.value
  if (!form) return
  isMutating.value = true
  try {
    const resp = await patchAdminZone(id, {
      name: form.name,
      capacity: Number(form.capacity),
      description: form.description,
      is_active: form.isActive,
      details_json: buildDetailsJson(form)
    })
    const idx = allZones.value.findIndex(z => z.id === id)
    if (idx !== -1) allZones.value[idx] = resp.zone
    editingId.value = null
    currentEditForm.value = null
    showFeedback('success', 'Событие обновлено.')
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка обновления')
  } finally {
    isMutating.value = false
  }
}

async function removeZone(id: number, name: string) {
  if (!confirm(`Удалить событие «${name}»?`)) return
  isMutating.value = true
  try {
    await deleteAdminZone(id)
    allZones.value = allZones.value.filter(z => z.id !== id)
    showFeedback('success', `Событие «${name}» удалено.`)
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка удаления')
  } finally {
    isMutating.value = false
  }
}
</script>

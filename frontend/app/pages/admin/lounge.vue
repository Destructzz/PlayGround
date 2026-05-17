<template>
  <div class="space-y-6">
    <!-- Header Card -->
    <header class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between rounded-[1.25rem] border border-white/5 bg-[#050f17] p-8 shadow-2xl">
      <div>
        <p class="text-[10px] font-black uppercase tracking-[0.4em] text-cyan-300/50">Lounge Management</p>
        <h1 class="mt-2 text-3xl font-black tracking-tight text-white">Lounge Zones</h1>
        <p class="mt-3 text-sm text-zinc-400 max-w-2xl leading-relaxed">Управление лаунж-зонами, перками и настройками атмосферы.</p>
      </div>
      <div class="flex items-center gap-4 flex-shrink-0">
        <button
          type="button"
          class="rounded-[0.9rem] bg-cyan-300 px-6 py-3 text-xs font-black uppercase tracking-widest text-[#020c13] transition hover:bg-cyan-200 hover:shadow-[0_0_20px_rgba(34,211,238,0.4)] active:scale-95"
          @click="showCreate = !showCreate"
        >
          {{ showCreate ? 'Скрыть форму' : '+ Добавить зону' }}
        </button>
      </div>
    </header>

    <!-- Feedback -->
    <div
      v-if="feedbackMessage"
      class="rounded-[0.8rem] border px-4 py-3 text-sm shadow-lg"
      :class="feedbackTone === 'error' ? 'border-orange-300/30 bg-orange-500/10 text-orange-100' : 'border-emerald-300/30 bg-emerald-500/10 text-emerald-100'"
    >
      {{ feedbackMessage }}
    </div>

      <!-- Create form -->
      <div v-if="showCreate" class="mb-8 rounded-[1rem] border border-cyan-400/20 bg-[#050f17] p-6 shadow-2xl">
        <p class="mb-5 text-xs font-black uppercase tracking-[0.3em] text-cyan-100/50">Новая Lounge-зона</p>
        <form class="space-y-5" @submit.prevent="submitCreate">
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Название</label>
              <input v-model="createForm.name" type="text" placeholder="Aurora Corner" required
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
            </div>
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Вместимость</label>
              <input v-model="createForm.capacity" type="number" min="1" required
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none">
            </div>
          </div>

          <div>
            <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Описание</label>
            <textarea v-model="createForm.description" rows="3" placeholder="Атмосфера, особенности, сценарий зоны..."
              class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none" />
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Настроение (mood)</label>
              <input v-model="createForm.mood" type="text" placeholder="Quiet social / Premium private"
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
            </div>
            <div>
              <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Цвет акцента</label>
              <div class="flex items-center gap-3">
                <input v-model="createForm.accentColor" type="color"
                  class="h-[46px] w-14 cursor-pointer rounded-[0.6rem] border border-cyan-400/18 bg-[#06131c] p-1">
                <span class="text-sm font-mono text-zinc-400">{{ createForm.accentColor }}</span>
              </div>
            </div>
          </div>

          <div>
            <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Перки (каждый с новой строки)</label>
            <textarea v-model="createForm.perksRaw" rows="3" placeholder="Signature drinks&#10;Ambient lighting&#10;PS5 side station"
              class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none" />
          </div>

          <div class="flex items-center gap-4">
            <label class="flex items-center gap-2 text-sm text-zinc-300">
              <input v-model="createForm.premium" type="checkbox" class="h-4 w-4 accent-cyan-300">
              Premium зона
            </label>
            <label class="flex items-center gap-2 text-sm text-zinc-300">
              <input v-model="createForm.isActive" type="checkbox" class="h-4 w-4 accent-cyan-300">
              Активна
            </label>
          </div>

          <!-- Zone tag -->
          <div>
            <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Zone Tag</label>
            <select v-model="createForm.zoneTagId"
              class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none">
              <option value="">Выбери tag</option>
              <option v-for="tag in zoneTags" :key="tag.id" :value="String(tag.id)">{{ tag.name }}</option>
            </select>
          </div>

          <!-- Управление тарифом -->
          <div class="rounded-[0.8rem] border border-cyan-400/10 bg-[#061018] p-4 space-y-4">
            <label class="flex items-center gap-2 text-sm font-bold text-white cursor-pointer select-none">
              <input v-model="createForm.hasTariff" type="checkbox" class="h-4 w-4 accent-cyan-300">
              <span>Требуется платный тариф</span>
            </label>
            
            <div v-if="createForm.hasTariff" class="grid gap-4 sm:grid-cols-3 pt-4 border-t border-white/5">
              <div>
                <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.25em] text-cyan-100/45">Стоимость (₽/час)</label>
                <input v-model="createForm.tariffPrice" type="number" min="0" required
                  class="w-full rounded-[0.6rem] border border-cyan-400/18 bg-[#06131c] px-3 py-2 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
              </div>
              <div>
                <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.25em] text-cyan-100/45">Длительность (мин)</label>
                <input v-model="createForm.tariffDuration" type="number" min="1" required
                  class="w-full rounded-[0.6rem] border border-cyan-400/18 bg-[#06131c] px-3 py-2 text-white focus:border-cyan-300 focus:outline-none">
              </div>
              <div>
                <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.25em] text-cyan-100/45">Описание тарифа</label>
                <input v-model="createForm.tariffDesc" type="text" placeholder="Почасовая аренда"
                  class="w-full rounded-[0.6rem] border border-cyan-400/18 bg-[#06131c] px-3 py-2 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
              </div>
            </div>
          </div>

          <button type="submit" :disabled="isMutating"
            class="w-full rounded-[0.9rem] bg-cyan-300 py-3.5 text-sm font-black uppercase tracking-widest text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50">
            Создать lounge-зону
          </button>
        </form>
      </div>

      <!-- Loading -->
      <div v-if="isLoading" class="rounded-[1rem] border border-white/5 bg-[#050f17] p-12 text-center text-sm text-zinc-400">
        Загружаем lounge-зоны...
      </div>

      <!-- Zones list -->
      <div v-else class="space-y-4">
        <div v-if="!loungeZones.length" class="rounded-[1rem] border border-white/8 bg-[#050f17] p-8 text-center text-sm text-zinc-400">
          Lounge-зоны ещё не созданы. Нажми «Добавить зону» выше.
        </div>

        <div
          v-for="zone in loungeZones"
          :key="zone.id"
          class="overflow-hidden rounded-[1rem] border border-white/8 bg-[#050f17] shadow-xl"
        >
          <!-- Zone header -->
          <div class="flex items-start justify-between gap-4 px-6 py-5">
            <div class="flex-1 min-w-0">
              <div class="flex flex-wrap items-center gap-2 mb-2">
                <span class="rounded-full border border-cyan-300/20 bg-cyan-400/10 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.25em] text-cyan-100">
                  Lounge
                </span>
                <span v-if="zoneDetails(zone).premium" class="rounded-full border border-amber-300/20 bg-amber-400/10 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.25em] text-amber-100">
                  Premium
                </span>
                <span v-if="zone.isActive" class="rounded-full border border-emerald-300/20 bg-emerald-400/10 px-2 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-emerald-200">active</span>
                <span v-else class="rounded-full border border-zinc-300/15 bg-zinc-500/10 px-2 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-400">inactive</span>
                
                <span v-if="getZoneService(zone.id)" class="rounded-full border border-cyan-400/20 bg-cyan-400/5 px-2.5 py-1 text-[10px] font-bold text-cyan-200">
                  Тариф: {{ getZoneService(zone.id)?.price }} ₽/ч
                </span>
                <span v-else class="rounded-full border border-emerald-400/20 bg-emerald-400/5 px-2.5 py-1 text-[10px] font-bold text-emerald-300">
                  Бесплатная бронь
                </span>
              </div>
              <h3 class="text-xl font-black text-white">{{ zone.name }}</h3>
              <p class="mt-1 text-sm text-zinc-400">{{ zone.description || 'Описание не заполнено.' }}</p>

              <!-- Perks chips -->
              <div v-if="zoneDetails(zone).perks?.length" class="mt-3 flex flex-wrap gap-1.5">
                <span
                  v-for="perk in zoneDetails(zone).perks"
                  :key="perk"
                  class="rounded-full border border-white/8 bg-white/4 px-2.5 py-1 text-[10px] text-zinc-400"
                >
                  {{ perk }}
                </span>
              </div>
            </div>

            <div class="flex items-center gap-2 flex-shrink-0">
              <div class="rounded-[0.7rem] border border-white/8 bg-white/4 px-4 py-2 text-right">
                <p class="text-[10px] uppercase tracking-[0.2em] text-zinc-500">Capacity</p>
                <p class="text-lg font-black text-white">{{ zone.capacity }}</p>
              </div>
              <button type="button"
                class="rounded-[0.8rem] border border-cyan-400/18 px-4 py-2.5 text-sm font-bold text-white transition hover:border-cyan-300/40"
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
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Название</label>
                  <input v-model="currentEditForm.name" type="text"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none">
                </div>
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Вместимость</label>
                  <input v-model="currentEditForm.capacity" type="number" min="1"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none">
                </div>
              </div>

              <div>
                <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Описание</label>
                <textarea v-model="currentEditForm.description" rows="2"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none" />
              </div>

              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Mood</label>
                  <input v-model="currentEditForm.mood" type="text"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none">
                </div>
                <div>
                  <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Цвет акцента</label>
                  <div class="flex items-center gap-3">
                    <input v-model="currentEditForm.accentColor" type="color"
                      class="h-[46px] w-14 cursor-pointer rounded-[0.6rem] border border-cyan-400/18 bg-[#071926] p-1">
                    <span class="text-sm font-mono text-zinc-400">{{ currentEditForm.accentColor }}</span>
                  </div>
                </div>
              </div>

              <div>
                <label class="mb-2 block text-xs font-bold uppercase tracking-[0.25em] text-cyan-100/45">Перки (каждый с новой строки)</label>
                <textarea v-model="currentEditForm.perksRaw" rows="3"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none" />
              </div>

              <div class="flex items-center gap-4">
                <label class="flex items-center gap-2 text-sm text-zinc-300">
                  <input v-model="currentEditForm.premium" type="checkbox" class="h-4 w-4 accent-cyan-300">
                  Premium
                </label>
                <label class="flex items-center gap-2 text-sm text-zinc-300">
                  <input v-model="currentEditForm.isActive" type="checkbox" class="h-4 w-4 accent-cyan-300">
                  Активна
                </label>
              </div>

              <!-- Управление тарифом -->
              <div class="rounded-[0.8rem] border border-cyan-400/10 bg-[#061018] p-4 space-y-4">
                <label class="flex items-center gap-2 text-sm font-bold text-white cursor-pointer select-none">
                  <input v-model="currentEditForm.hasTariff" type="checkbox" class="h-4 w-4 accent-cyan-300">
                  <span>Требуется платный тариф</span>
                </label>
                
                <div v-if="currentEditForm.hasTariff" class="grid gap-4 sm:grid-cols-3 pt-4 border-t border-white/5">
                  <div>
                    <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.25em] text-cyan-100/45">Стоимость (₽/час)</label>
                    <input v-model="currentEditForm.tariffPrice" type="number" min="0" required
                      class="w-full rounded-[0.6rem] border border-cyan-400/18 bg-[#071926] px-3 py-2 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
                  </div>
                  <div>
                    <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.25em] text-cyan-100/45">Длительность (мин)</label>
                    <input v-model="currentEditForm.tariffDuration" type="number" min="1" required
                      class="w-full rounded-[0.6rem] border border-cyan-400/18 bg-[#071926] px-3 py-2 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
                  </div>
                  <div>
                    <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.25em] text-cyan-100/45">Описание тарифа</label>
                    <input v-model="currentEditForm.tariffDesc" type="text" placeholder="Почасовая аренда"
                      class="w-full rounded-[0.6rem] border border-cyan-400/18 bg-[#071926] px-3 py-2 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none">
                  </div>
                </div>
              </div>

              <div class="flex gap-3">
                <button type="submit" :disabled="isMutating"
                  class="rounded-[0.85rem] bg-cyan-300 px-5 py-3 text-sm font-black text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50">
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
</template>

<script setup lang="ts">
import { getAdminZones, getAdminZoneTags, patchAdminZone, deleteAdminZone, createAdminZone, getAdminServices, createAdminService, patchAdminService, deleteAdminService } from '~/api/admin'
import type { AdminZone, AdminZoneTag, AdminService } from '~/api/types'

definePageMeta({ middleware: 'admin' })
useHead({ title: 'Admin · Lounge — PlayGround' })

// ── State ─────────────────────────────────────────────────
const isLoading = ref(true)
const isMutating = ref(false)
const feedbackTone = ref<'success' | 'error'>('success')
const feedbackMessage = ref('')
const showCreate = ref(false)
const editingId = ref<number | null>(null)

const allZones = ref<AdminZone[]>([])
const zoneTags = ref<AdminZoneTag[]>([])
const allServices = ref<AdminService[]>([])

interface LoungeForm {
  name: string
  capacity: string
  description: string
  mood: string
  accentColor: string
  perksRaw: string
  premium: boolean
  isActive: boolean
  zoneTagId: string
  // Tariff integration fields
  hasTariff: boolean
  tariffPrice: string
  tariffDuration: string
  tariffDesc: string
}

const createForm = ref<LoungeForm>(emptyForm())
const currentEditForm = ref<LoungeForm | null>(null)

// ── Computed ──────────────────────────────────────────────
const loungeZones = computed(() =>
  allZones.value
    .filter(z => z.zone_type === 'lounge')
    .map(z => ({
      ...z,
      isActive: z.is_active,
      description: typeof z.description === 'string' ? z.description : ''
    }))
)

// ── Lifecycle ─────────────────────────────────────────────
onMounted(async () => {
  try {
    const [zonesResp, tagsResp, servicesResp] = await Promise.all([
      getAdminZones(),
      getAdminZoneTags(),
      getAdminServices()
    ])
    allZones.value = zonesResp.zones ?? []
    zoneTags.value = tagsResp.zone_tags ?? []
    allServices.value = servicesResp.services ?? []
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка загрузки')
  } finally {
    isLoading.value = false
  }
})

// ── Helpers ───────────────────────────────────────────────
function emptyForm(): LoungeForm {
  return {
    name: '',
    capacity: '6',
    description: '',
    mood: '',
    accentColor: '#22d3ee',
    perksRaw: '',
    premium: false,
    isActive: true,
    zoneTagId: '',
    hasTariff: false,
    tariffPrice: '1000',
    tariffDuration: '60',
    tariffDesc: 'Почасовая аренда lounge-зоны'
  }
}

function zoneDetails(zone: AdminZone & { isActive: boolean; description: string }) {
  try {
    const d = zone.details_json as any
    if (typeof d === 'object' && d !== null) return d
    if (typeof d === 'string') return JSON.parse(d)
  } catch {}
  return {}
}

function buildDetailsJson(form: LoungeForm): string {
  const perks = form.perksRaw.split('\n').map(s => s.trim()).filter(Boolean)
  return JSON.stringify({
    mood: form.mood,
    accent: form.accentColor,
    premium: form.premium,
    perks
  })
}

function getZoneService(zoneId: number) {
  return allServices.value.find(s => s.zone_id === zoneId)
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
  const d = zoneDetails(zone as any)
  const svc = getZoneService(id)
  currentEditForm.value = {
    name: zone.name,
    capacity: String(zone.capacity),
    description: typeof zone.description === 'string' ? zone.description : '',
    mood: d.mood ?? '',
    accentColor: d.accent ?? '#22d3ee',
    perksRaw: Array.isArray(d.perks) ? d.perks.join('\n') : '',
    premium: !!d.premium,
    isActive: zone.is_active,
    zoneTagId: String(zone.zone_tag_id),
    hasTariff: !!svc,
    tariffPrice: svc ? String(svc.price) : '1000',
    tariffDuration: svc ? String(svc.duration) : '60',
    tariffDesc: svc ? String(svc.description || '') : 'Почасовая аренда lounge-зоны'
  }
}

function showFeedback(tone: 'success' | 'error', msg: string) {
  feedbackTone.value = tone
  feedbackMessage.value = msg
  if (tone === 'success') setTimeout(() => { feedbackMessage.value = '' }, 3000)
}

// ── Actions ───────────────────────────────────────────────
async function submitCreate() {
  if (!createForm.value.zoneTagId) { showFeedback('error', 'Выбери zone tag'); return }
  isMutating.value = true
  try {
    const resp = await createAdminZone({
      name: createForm.value.name,
      type: 'lounge',
      zone_tag_id: Number(createForm.value.zoneTagId),
      capacity: Number(createForm.value.capacity),
      description: createForm.value.description,
      is_active: createForm.value.isActive,
      details_json: buildDetailsJson(createForm.value)
    })
    
    allZones.value.push(resp.zone)

    // Handle seamless service creation if requested
    if (createForm.value.hasTariff) {
      try {
        const svcResp = await createAdminService({
          name: resp.zone.name + " Tarif",
          zone_id: resp.zone.id,
          duration: Number(createForm.value.tariffDuration),
          price: Number(createForm.value.tariffPrice),
          currency: 'RUB',
          description: createForm.value.tariffDesc,
          is_active: true,
          details_json: '{}'
        })
        allServices.value.push(svcResp.service)
      } catch (svcErr: any) {
        showFeedback('error', `Зона создана, но не удалось создать тариф: ${svcErr.message}`)
      }
    }

    createForm.value = emptyForm()
    showCreate.value = false
    showFeedback('success', `Зона «${resp.zone.name}» создана.`)
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

    // Handle service update/creation/deletion
    const existingSvc = getZoneService(id)
    if (form.hasTariff) {
      if (existingSvc) {
        // Update existing service
        const svcResp = await patchAdminService(existingSvc.id, {
          name: form.name + " Tarif",
          duration: Number(form.tariffDuration),
          price: Number(form.tariffPrice),
          description: form.tariffDesc
        })
        const sIdx = allServices.value.findIndex(s => s.id === existingSvc.id)
        if (sIdx !== -1) allServices.value[sIdx] = svcResp.service
      } else {
        // Create new service
        const svcResp = await createAdminService({
          name: form.name + " Tarif",
          zone_id: id,
          duration: Number(form.tariffDuration),
          price: Number(form.tariffPrice),
          currency: 'RUB',
          description: form.tariffDesc,
          is_active: true,
          details_json: '{}'
        })
        allServices.value.push(svcResp.service)
      }
    } else if (existingSvc) {
      // Delete service if checkbox was unchecked
      await deleteAdminService(existingSvc.id)
      allServices.value = allServices.value.filter(s => s.id !== existingSvc.id)
    }

    editingId.value = null
    currentEditForm.value = null
    showFeedback('success', 'Зона и тариф обновлены.')
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка обновления')
  } finally {
    isMutating.value = false
  }
}

async function removeZone(id: number, name: string) {
  if (!confirm(`Удалить зону «${name}»?`)) return
  isMutating.value = true
  try {
    // Delete any associated services first (backend will cascade delete but let's clean local state)
    const existingSvc = getZoneService(id)
    await deleteAdminZone(id)
    allZones.value = allZones.value.filter(z => z.id !== id)
    if (existingSvc) {
      allServices.value = allServices.value.filter(s => s.id !== existingSvc.id)
    }
    showFeedback('success', `Зона «${name}» удалена.`)
  } catch (e: any) {
    showFeedback('error', e?.message ?? 'Ошибка удаления')
  } finally {
    isMutating.value = false
  }
}
</script>

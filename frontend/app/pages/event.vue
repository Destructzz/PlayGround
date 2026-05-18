<template>
  <div class="min-h-screen bg-[#020c13] pb-24 pt-24 text-white">
    <div class="mx-auto flex max-w-[1400px] flex-col gap-10 px-6 sm:px-10 lg:px-12">
      <ExperienceHero
        eyebrow="Event zone"
        title="Афиша событий"
        description="Турниры, митапы, community-вечера — регистрация в один клик для гостей."
        accent="#f472b6"
      />

      <!-- Loading -->
      <ExperienceStatePanel
        v-if="catalogLoading"
        title="Загружаем афишу"
        description="Подтягиваем актуальные события и доступность мест."
        hint="loading..."
      />

      <!-- Error -->
      <div
        v-else-if="catalogError"
        class="rounded-lg border border-orange-500/30 bg-orange-500/10 px-5 py-4 text-xs text-orange-200"
      >
        Не удалось загрузить события: {{ catalogError }}
      </div>

      <!-- Empty -->
      <ExperienceStatePanel
        v-else-if="!zones.length"
        title="Пока ничего не запланировано"
        description="Следи за обновлениями — скоро появятся новые события."
        hint="empty"
        accent="#f472b6"
      />

      <template v-else>
        <div data-testid="event-client-ready" class="hidden" />
        <!-- Event cards — browser-style category tabs -->
        <div>
          <!-- Format filter tabs -->
          <div class="flex gap-1 overflow-x-auto rounded-t-lg border border-b-0 border-white/10 bg-[#020a10] px-2 pt-2">
            <button
              type="button"
              class="flex-shrink-0 rounded-t-md px-4 py-2 text-xs font-bold transition-all"
              :class="selectedFormat === ''
                ? 'bg-[#071724]/75 text-white border border-b-0 border-white/20 -mb-px'
                : 'text-zinc-400 hover:text-zinc-200'"
              @click="selectedFormat = ''"
            >
              Все
            </button>
            <button
              v-for="fmt in availableFormats"
              :key="fmt"
              type="button"
              class="flex-shrink-0 rounded-t-md px-4 py-2 text-xs font-bold transition-all"
              :class="selectedFormat === fmt
                ? 'bg-[#071724]/75 text-white border border-b-0 border-white/20 -mb-px'
                : 'text-zinc-400 hover:text-zinc-200'"
              @click="selectedFormat = fmt"
            >
              {{ fmt }}
            </button>
          </div>
          <div class="rounded-b-lg border border-white/10 bg-[#071724]/75 p-0.5" />
        </div>

        <!-- Events grid -->
        <section class="grid gap-6 lg:grid-cols-3">
          <ExperienceCard
            v-for="zone in filteredZones"
            :key="zone.id"
            eyebrow="Event"
            :title="zone.name"
            :description="zone.description"
            :details="eventDetails(zone)"
            :badge="zone.capacity > 0 ? `${zone.capacity} мест` : '—'"
            badge-tone="success"
            footer-label="Вместимость"
            :footer-value="String(zone.capacity)"
            accent="#f472b6"
          >
            <!-- Beautiful event schedule and format header, styled completely separate from the details list -->
            <div class="flex flex-wrap items-center gap-2.5 text-[11px] font-bold mt-1 mb-2">
              <span class="flex items-center gap-1.5 text-fuchsia-400">
                <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                {{ formatEventDate((zone.details_json as any)?.start_time) }}
              </span>
              <span v-if="(zone.details_json as any)?.format" class="rounded-md border border-white/10 bg-white/5 px-2 py-0.5 text-zinc-400">
                {{ (zone.details_json as any)?.format }}
              </span>
            </div>
            <template #cta>
              <button
                :data-testid="`event-card-${zone.id}`"
                class="rounded-lg px-4 py-2 text-xs font-bold transition"
                :class="selectedZoneId === zone.id
                  ? 'bg-fuchsia-400 text-black shadow-[0_0_12px_rgba(244,114,182,0.2)] hover:bg-fuchsia-300'
                  : 'border border-white/10 bg-white/5 text-white hover:border-white/20 hover:bg-white/10'"
                @click="selectZone(zone.id)"
              >
                {{ selectedZoneId === zone.id ? 'Выбрано' : 'Участвовать' }}
              </button>
            </template>
          </ExperienceCard>
        </section>

        <!-- Registration flow -->
        <ExperienceFlowPanel
          eyebrow="Registration flow"
          title="Регистрация на событие"
          description="Выбери событие и подтверди участие — место фиксируется автоматически на твой аккаунт."
          accent="#f472b6"
        >
          <div class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
            <div class="space-y-6">

              <!-- Selected event card -->
              <div
                data-testid="event-selected-card"
                class="rounded-xl border border-white/10 bg-[#071724]/75 backdrop-blur-md p-5 shadow-[0_8px_30px_rgba(0,0,0,0.4)]"
              >
                <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div class="flex-1 min-w-0">
                    <p class="text-[10px] font-bold uppercase tracking-[0.25em] text-zinc-500">Выбранное событие</p>
                    <h3 class="mt-2 text-xl font-black text-white uppercase truncate">{{ selectedZone?.name ?? 'Сначала выбери событие' }}</h3>
                    <p class="mt-1 text-xs leading-relaxed text-zinc-300">{{ selectedZone?.description ?? 'Нажми «Участвовать» на одной из карточек выше.' }}</p>
                    
                    <div v-if="selectedZone" class="mt-3.5 flex flex-wrap gap-4 text-xs">
                      <div class="flex items-center gap-1.5 text-fuchsia-400 font-medium">
                        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        <span>{{ formatEventDate((selectedZone.details_json as any)?.start_time) }}</span>
                      </div>
                      <div v-if="(selectedZone.details_json as any)?.format" class="flex items-center gap-1.5 text-zinc-400 font-medium">
                        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 113.536 0V21h2v-2.757" />
                        </svg>
                        <span>{{ (selectedZone.details_json as any)?.format }}</span>
                      </div>
                    </div>
                  </div>
                  <div
                    v-if="selectedZone"
                    class="rounded-lg border border-white/10 bg-[#020a10] px-4 py-3 text-right"
                  >
                    <p class="text-[9px] uppercase tracking-[0.25em] text-zinc-500">Вместимость</p>
                    <p class="mt-1 text-xl font-black text-white">{{ selectedZone.capacity }}</p>
                  </div>
                </div>
              </div>

              <!-- Service selector (ticket type) -->
              <div v-if="selectedZone && selectedZone.services.length > 0">
                <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.25em] text-zinc-500">Тип билета</p>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="svc in selectedZone.services"
                    :key="svc.id"
                    type="button"
                    class="rounded-lg border px-4 py-2 text-xs font-bold transition"
                    :class="selectedServiceId === svc.id
                      ? 'border-fuchsia-400 bg-fuchsia-400 text-black shadow-[0_0_16px_rgba(244,114,182,0.15)]'
                      : 'border-white/10 bg-[#020a10] text-zinc-300 hover:border-fuchsia-400/30 hover:text-white'"
                    @click="selectedServiceId = svc.id"
                  >
                    {{ svc.name }}
                    <span class="ml-1.5 text-[10px] opacity-60">{{ svc.price }} ₽</span>
                  </button>
                </div>
              </div>

              <!-- Auth notice / user drawer trigger -->
              <div v-if="!authStore.isAuthenticated" class="rounded-lg border border-amber-500/30 bg-amber-500/10 px-5 py-4 text-xs text-amber-200">
                <p class="font-bold">Требуется вход</p>
                <p class="mt-1 text-amber-200/70">Войди через Google для регистрации на событие. Данные профиля подтянутся автоматически.</p>
              </div>
              <button
                v-else
                type="button"
                class="flex w-full items-center gap-4 rounded-xl border border-white/10 bg-[#020a10] px-5 py-4 text-left transition hover:border-fuchsia-400/30"
                @click="drawerOpen = true"
              >
                <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-fuchsia-400/30 bg-fuchsia-400/10 text-sm font-black text-fuchsia-200">
                  {{ authStore.user?.name?.charAt(0).toUpperCase() ?? '?' }}
                </div>
                <div class="min-w-0">
                  <p class="text-xs font-bold text-white">{{ authStore.user?.name }}</p>
                  <p class="truncate text-[10px] text-zinc-400">{{ authStore.user?.email }} · Нажми чтобы просмотреть</p>
                </div>
              </button>
            </div>

            <!-- Summary sidebar -->
            <div class="space-y-4 rounded-xl border border-white/10 bg-[#071724]/75 backdrop-blur-md p-5 shadow-[0_8px_30px_rgba(0,0,0,0.4)]">
              <div>
                <p class="text-[9px] font-bold uppercase tracking-[0.25em] text-zinc-500">Регистрация</p>
                <h3 class="mt-1 text-xl font-black uppercase tracking-tight text-white">Участие</h3>
              </div>

              <div class="space-y-3 text-xs text-zinc-300">
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Событие</span>
                  <span class="font-bold text-white uppercase">{{ selectedZone?.name ?? '—' }}</span>
                </div>
                <div v-if="selectedZone" class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Когда</span>
                  <span class="font-bold text-fuchsia-400">{{ formatEventDate((selectedZone.details_json as any)?.start_time) }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Тариф</span>
                  <span class="font-bold text-white">{{ selectedServiceLabel }}</span>
                </div>
                <div class="space-y-3 pt-3 border-t border-white/5">
                  <div class="flex items-center justify-between gap-4">
                    <span class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Контакты брони</span>
                    <label class="flex items-center gap-2 cursor-pointer">
                      <input type="checkbox" v-model="useCustomContacts" class="rounded border-white/10 bg-black/50 text-fuchsia-400 focus:ring-0 focus:ring-offset-0 h-3.5 w-3.5" />
                      <span class="text-[10px] text-zinc-400">Указать другие</span>
                    </label>
                  </div>
                  
                  <div v-if="!useCustomContacts" class="flex flex-col gap-1 text-[11px] text-white">
                    <span class="font-bold">{{ authStore.user?.name || '—' }}</span>
                    <span class="text-zinc-400">{{ authStore.user?.email || '—' }}</span>
                  </div>
                  
                  <div v-else class="space-y-2 pt-2">
                    <input v-model="customContacts.name" type="text" placeholder="Имя" class="w-full rounded-md border border-white/10 bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-fuchsia-500/50" />
                    <input v-model="customContacts.email" type="email" placeholder="Email" class="w-full rounded-md border border-white/10 bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-fuchsia-500/50" />
                    <input v-model="customContacts.phone" type="tel" placeholder="Телефон" class="w-full rounded-md border border-white/10 bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-fuchsia-500/50" />
                  </div>
                </div>
              </div>

              <div
                v-if="validationMessage"
                data-testid="event-validation-message"
                class="rounded-lg border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-xs text-rose-200"
              >
                {{ validationMessage }}
              </div>

              <div
                v-if="submitError"
                data-testid="event-submit-error"
                class="rounded-lg border border-orange-500/30 bg-orange-500/10 px-4 py-3 text-xs text-orange-200"
              >
                {{ submitError }}
              </div>

              <div
                v-if="registrationSuccess"
                data-testid="event-success-message"
                class="rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-4 py-3 text-xs text-emerald-200"
              >
                Регистрация подтверждена. Твоё место на «{{ selectedZone?.name }}» зафиксировано.
              </div>

              <button
                data-testid="event-submit"
                type="button"
                class="w-full rounded-lg bg-fuchsia-400 px-5 py-3.5 text-xs font-bold uppercase tracking-widest text-black transition hover:bg-fuchsia-300 hover:shadow-[0_0_16px_rgba(244,114,182,0.15)] disabled:cursor-not-allowed disabled:bg-fuchsia-400/40"
                :disabled="pending || !authStore.isAuthenticated"
                @click="submitRegistration"
              >
                {{ pending ? 'Фиксируем...' : 'Подтвердить участие' }}
              </button>

              <button
                data-testid="event-reset"
                type="button"
                class="w-full rounded-lg border border-white/10 bg-white/5 px-5 py-2.5 text-xs font-bold text-zinc-300 transition hover:bg-white/10 hover:text-white"
                @click="resetForm"
              >
                Очистить форму
              </button>
            </div>
          </div>
        </ExperienceFlowPanel>
      </template>
    </div>

    <!-- User Drawer -->
    <Teleport to="body">
      <Transition name="drawer">
        <div
          v-if="drawerOpen"
          class="fixed inset-0 z-50 flex justify-end"
        >
          <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="drawerOpen = false" />
          <div class="relative z-10 flex h-full w-full max-w-sm flex-col bg-[#050f17] shadow-2xl">
            <div class="flex items-center justify-between border-b border-white/8 px-6 py-5">
              <p class="text-xs font-black uppercase tracking-[0.3em] text-fuchsia-100/50">Профиль участника</p>
              <button
                type="button"
                class="rounded-full border border-white/10 p-2 text-zinc-400 transition hover:border-white/25 hover:text-white"
                @click="drawerOpen = false"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div class="flex-1 overflow-y-auto px-6 py-8">
              <div class="flex flex-col items-center text-center">
                <div class="flex h-20 w-20 items-center justify-center rounded-full border-2 border-fuchsia-300/30 bg-fuchsia-400/10 text-3xl font-black text-fuchsia-200">
                  {{ authStore.user?.name?.charAt(0).toUpperCase() ?? '?' }}
                </div>
                <h2 class="mt-5 text-xl font-black text-white">{{ authStore.user?.name }}</h2>
                <p class="mt-1 text-sm text-fuchsia-100/50">{{ authStore.user?.email }}</p>
                <span class="mt-3 rounded-full border border-fuchsia-300/20 bg-fuchsia-400/10 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-fuchsia-100">
                  {{ authStore.user?.role }}
                </span>
              </div>
              <div class="mt-8 space-y-3 text-sm">
                <div class="rounded-[0.8rem] border border-white/8 bg-white/4 px-4 py-3">
                  <p class="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Имя участника</p>
                  <p class="mt-1 font-medium text-white">{{ authStore.user?.name ?? '—' }}</p>
                </div>
                <div class="rounded-[0.8rem] border border-white/8 bg-white/4 px-4 py-3">
                  <p class="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Email</p>
                  <p class="mt-1 font-medium text-white">{{ authStore.user?.email ?? '—' }}</p>
                </div>
              </div>
              <p class="mt-6 text-center text-[11px] text-zinc-600">
                Данные профиля автоматически прикрепятся к регистрации.
              </p>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import type { CatalogZoneWithServices } from '../api/types'
import { getEventCatalog } from '~/api/catalog'
import { createBooking } from '~/api/booking'
import { useAuthStore } from '../stores/auth'

useHead({ title: 'Events - PlayGround' })

const authStore = useAuthStore()

// ── State ────────────────────────────────────────────────
const zones = ref<CatalogZoneWithServices[]>([])
const catalogLoading = ref(true)
const catalogError = ref('')

const selectedZoneId = ref<number | null>(null)
const selectedServiceId = ref<number | null>(null)
const selectedFormat = ref('')

const drawerOpen = ref(false)
const pending = ref(false)
const registrationSuccess = ref(false)
const submitError = ref('')
const validationMessage = ref('')

const useCustomContacts = ref(false)
const customContacts = ref({ name: '', email: '', phone: '' })

// ── Computed ──────────────────────────────────────────────


const availableFormats = computed(() => {
  const fmts = new Set<string>()
  for (const z of zones.value) {
    const d = z.details_json as any
    if (d?.format) fmts.add(d.format)
  }
  return Array.from(fmts)
})

const filteredZones = computed(() =>
  selectedFormat.value
    ? zones.value.filter(z => (z.details_json as any)?.format === selectedFormat.value)
    : zones.value
)

const selectedZone = computed(() =>
  zones.value.find(z => z.id === selectedZoneId.value) ?? null
)

const selectedServiceLabel = computed(() => {
  if (!selectedZone.value) return '—'
  if (!selectedZone.value.services || selectedZone.value.services.length === 0) {
    return 'Вход свободный'
  }
  if (!selectedServiceId.value) return '—'
  const svc = selectedZone.value.services.find(s => s.id === selectedServiceId.value)
  return svc ? `${svc.name} · ${svc.price} ₽` : '—'
})

// ── Lifecycle ─────────────────────────────────────────────
onMounted(async () => {
  try {
    const data = await getEventCatalog()
    zones.value = data.zones ?? []
    if (zones.value.length > 0) {
      selectedZoneId.value = zones.value[0]!.id
      selectedServiceId.value = zones.value[0]!.services?.[0]?.id ?? null
    }
  } catch (e: any) {
    catalogError.value = e?.message ?? 'Неизвестная ошибка'
  } finally {
    catalogLoading.value = false
  }
})

// ── Methods ───────────────────────────────────────────────
function selectZone(id: number) {
  selectedZoneId.value = id
  const zone = zones.value.find(z => z.id === id)
  selectedServiceId.value = zone?.services?.[0]?.id ?? null
  registrationSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
}

function resetForm() {
  selectedZoneId.value = null
  selectedServiceId.value = null
  registrationSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
  useCustomContacts.value = false
  customContacts.value = { name: '', email: '', phone: '' }
}

function validate(): string {
  if (!selectedZone.value) return 'Сначала выбери событие из афиши.'
  if (selectedZone.value.services && selectedZone.value.services.length > 0 && !selectedServiceId.value) {
    return 'Выбери тип билета.'
  }
  if (!authStore.isAuthenticated) return 'Для регистрации необходимо войти в аккаунт.'
  return ''
}

async function submitRegistration() {
  validationMessage.value = validate()
  submitError.value = ''
  registrationSuccess.value = false
  if (validationMessage.value) return

  const zone = selectedZone.value!
  const service = zone.services.find(s => s.id === selectedServiceId.value)
  const duration = service?.duration ?? 60

  // For events the start/end time comes from the event details_json or we use a placeholder
  const d = zone.details_json as any
  const startTime = d?.start_time ?? new Date().toISOString()
  const endTime = d?.end_time ?? new Date(Date.now() + duration * 60000).toISOString()

  pending.value = true
  
  const cName = useCustomContacts.value ? customContacts.value.name : (authStore.user?.name || '')
  const cEmail = useCustomContacts.value ? customContacts.value.email : (authStore.user?.email || '')
  const cPhone = useCustomContacts.value ? customContacts.value.phone : (localStorage.getItem('playground_phone') || '')

  try {
    await createBooking({
      zone_id: zone.id,
      service_id: service?.id || 0,
      start_time: startTime,
      end_time: endTime,
      participants: 1,
      status: 'created',
      contact_name: cName,
      contact_email: cEmail,
      contact_phone: cPhone
    })
    registrationSuccess.value = true
  } catch (e: any) {
    submitError.value = e?.data?.error?.message ?? e?.message ?? 'Не удалось зарегистрироваться. Попробуй ещё раз.'
  } finally {
    pending.value = false
  }
}

function formatEventDate(isoStr?: string): string {
  if (!isoStr) return '—'
  try {
    const d = new Date(isoStr)
    return d.toLocaleString('ru-RU', {
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      weekday: 'short'
    })
  } catch {
    return isoStr
  }
}

function eventDetails(zone: CatalogZoneWithServices): string[] {
  const d = zone.details_json as any
  const parts: string[] = []
  if (Array.isArray(d?.speakers)) parts.push(...d.speakers)
  if (zone.services?.length) {
    parts.push(...zone.services.map(s => `Билет: ${s.name} · ${s.price} ₽`))
  } else {
    parts.push('Вход свободный')
  }
  return parts
}
</script>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.25s ease;
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}
.drawer-enter-from .relative,
.drawer-leave-to .relative {
  transform: translateX(100%);
}
</style>

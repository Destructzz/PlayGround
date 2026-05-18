<template>
  <div class="min-h-screen bg-[#020c13] pb-24 pt-24 text-white">
    <div class="mx-auto flex max-w-[1400px] flex-col gap-10 px-6 sm:px-10 lg:px-12">
      <ExperienceHero
        eyebrow="Event zone"
        title="Афиша событий"
        description="Турниры, митапы, community-вечера — регистрация в один клик для авторизованных гостей."
        :stats="heroStats"
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
        class="rounded-[0.9rem] border border-orange-300/35 bg-orange-500/12 px-5 py-4 text-sm text-orange-100"
      >
        Не удалось загрузить события: {{ catalogError }}
      </div>

      <!-- Empty -->
      <ExperienceStatePanel
        v-else-if="!zones.length"
        title="Пока ничего не запланировано"
        description="Следи за обновлениями — скоро появятся новые события."
        hint="empty"
      />

      <template v-else>
        <div data-testid="event-client-ready" class="hidden" />
        <!-- Event cards — browser-style category tabs -->
        <div>
          <!-- Format filter tabs -->
          <div class="flex gap-1 overflow-x-auto rounded-t-[0.6rem] border border-b-0 border-fuchsia-400/15 bg-[#07141d] px-2 pt-2">
            <button
              type="button"
              class="flex-shrink-0 rounded-t-[0.5rem] px-4 py-2 text-xs font-bold transition-all"
              :class="selectedFormat === ''
                ? 'bg-[#091924] text-white border border-b-0 border-fuchsia-400/25 -mb-px'
                : 'text-fuchsia-100/40 hover:text-fuchsia-100/70'"
              @click="selectedFormat = ''"
            >
              Все
            </button>
            <button
              v-for="fmt in availableFormats"
              :key="fmt"
              type="button"
              class="flex-shrink-0 rounded-t-[0.5rem] px-4 py-2 text-xs font-bold transition-all"
              :class="selectedFormat === fmt
                ? 'bg-[#091924] text-white border border-b-0 border-fuchsia-400/25 -mb-px'
                : 'text-fuchsia-100/40 hover:text-fuchsia-100/70'"
              @click="selectedFormat = fmt"
            >
              {{ fmt }}
            </button>
          </div>
          <div class="rounded-b-[0.6rem] rounded-tr-[0.6rem] border border-fuchsia-400/15 bg-[#091924] p-0.5" />
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
            <template #cta>
              <button
                :data-testid="`event-card-${zone.id}`"
                class="rounded-full px-5 py-2 text-sm font-bold transition"
                :class="selectedZoneId === zone.id
                  ? 'bg-white text-black'
                  : 'border border-white/15 bg-white/5 text-white hover:border-white/30'"
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
        >
          <div class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
            <div class="space-y-6">

              <!-- Selected event card -->
              <div
                data-testid="event-selected-card"
                class="rounded-[0.85rem] border border-fuchsia-400/20 bg-[#09131d] p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]"
              >
                <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <p class="text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Выбранное событие</p>
                    <h3 class="mt-2 text-2xl font-black text-white">{{ selectedZone?.name ?? 'Сначала выбери событие' }}</h3>
                    <p class="mt-2 text-sm leading-7 text-zinc-300">{{ selectedZone?.description ?? 'Нажми «Участвовать» на одной из карточек выше.' }}</p>
                  </div>
                  <div
                    v-if="selectedZone"
                    class="rounded-[0.7rem] border border-fuchsia-300/25 bg-[#081019] px-4 py-3 text-right"
                  >
                    <p class="text-[11px] uppercase tracking-[0.35em] text-cyan-100/45">Вместимость</p>
                    <p class="mt-2 text-2xl font-black text-white">{{ selectedZone.capacity }}</p>
                  </div>
                </div>
              </div>

              <!-- Service selector (ticket type) -->
              <div v-if="selectedZone && selectedZone.services.length > 0">
                <p class="mb-3 text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Тип билета</p>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="svc in selectedZone.services"
                    :key="svc.id"
                    type="button"
                    class="rounded-[0.55rem] border px-4 py-2 text-sm font-bold transition"
                    :class="selectedServiceId === svc.id
                      ? 'border-fuchsia-200 bg-fuchsia-300 text-[#020c13] shadow-[0_0_24px_rgba(244,114,182,0.25)]'
                      : 'border-fuchsia-400/20 bg-[#07141d] text-white hover:border-fuchsia-300/45'"
                    @click="selectedServiceId = svc.id"
                  >
                    {{ svc.name }}
                    <span class="ml-1.5 text-[10px] opacity-60">{{ svc.price }} ₽</span>
                  </button>
                </div>
              </div>

              <!-- Auth notice / user drawer trigger -->
              <div v-if="!authStore.isAuthenticated" class="rounded-[0.85rem] border border-amber-300/25 bg-amber-500/10 px-5 py-4 text-sm text-amber-100">
                <p class="font-bold">Требуется вход</p>
                <p class="mt-1 text-amber-100/70">Войди через Google для регистрации на событие. Данные профиля подтянутся автоматически.</p>
              </div>
              <button
                v-else
                type="button"
                class="flex w-full items-center gap-4 rounded-[0.85rem] border border-fuchsia-400/15 bg-[#07141d] px-5 py-4 text-left transition hover:border-fuchsia-300/30"
                @click="drawerOpen = true"
              >
                <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-fuchsia-400/15 text-sm font-black text-fuchsia-200">
                  {{ authStore.user?.name?.charAt(0).toUpperCase() ?? '?' }}
                </div>
                <div class="min-w-0">
                  <p class="text-sm font-bold text-white">{{ authStore.user?.name }}</p>
                  <p class="truncate text-xs text-fuchsia-100/50">{{ authStore.user?.email }} · Нажми чтобы просмотреть</p>
                </div>
              </button>
            </div>

            <!-- Summary sidebar -->
            <div class="space-y-4 rounded-[0.85rem] border border-fuchsia-400/18 bg-[#09131d] p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]">
              <div>
                <p class="text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Регистрация</p>
                <h3 class="mt-2 text-2xl font-black uppercase tracking-[-0.04em] text-white">Участие</h3>
              </div>

              <div class="space-y-3 text-sm text-zinc-300">
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Событие</span>
                  <span class="font-bold text-white">{{ selectedZone?.name ?? '—' }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Тариф</span>
                  <span class="font-bold text-white">{{ selectedServiceLabel }}</span>
                </div>
                <div class="space-y-3 pt-3 border-t border-white/8">
                  <div class="flex items-center justify-between gap-4">
                    <span class="text-xs uppercase tracking-wider text-cyan-100/45">Контакты брони</span>
                    <label class="flex items-center gap-2 cursor-pointer">
                      <input type="checkbox" v-model="useCustomContacts" class="rounded border-fuchsia-500/30 bg-black/50 text-fuchsia-400 focus:ring-fuchsia-400 focus:ring-offset-0" />
                      <span class="text-xs">Указать другие</span>
                    </label>
                  </div>
                  
                  <div v-if="!useCustomContacts" class="flex flex-col gap-1 text-xs text-white">
                    <span class="font-bold">{{ authStore.user?.name || '—' }}</span>
                    <span class="text-white/60">{{ authStore.user?.email || '—' }}</span>
                  </div>
                  
                  <div v-else class="space-y-2 pt-2">
                    <input v-model="customContacts.name" type="text" placeholder="Имя" class="w-full rounded bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-fuchsia-500/50" />
                    <input v-model="customContacts.email" type="email" placeholder="Email" class="w-full rounded bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-fuchsia-500/50" />
                    <input v-model="customContacts.phone" type="tel" placeholder="Телефон" class="w-full rounded bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-fuchsia-500/50" />
                  </div>
                </div>
              </div>

              <div
                v-if="validationMessage"
                data-testid="event-validation-message"
                class="rounded-[0.7rem] border border-rose-300/35 bg-rose-500/18 px-4 py-3 text-sm text-rose-100"
              >
                {{ validationMessage }}
              </div>

              <div
                v-if="submitError"
                data-testid="event-submit-error"
                class="rounded-[0.7rem] border border-orange-300/35 bg-orange-500/18 px-4 py-3 text-sm text-orange-100"
              >
                {{ submitError }}
              </div>

              <div
                v-if="registrationSuccess"
                data-testid="event-success-message"
                class="rounded-[0.7rem] border border-emerald-300/35 bg-emerald-500/18 px-4 py-3 text-sm text-emerald-100"
              >
                Регистрация подтверждена. Твоё место на «{{ selectedZone?.name }}» зафиксировано.
              </div>

              <button
                data-testid="event-submit"
                type="button"
                class="w-full rounded-[0.75rem] bg-fuchsia-300 px-5 py-4 text-sm font-black uppercase tracking-[0.25em] text-[#020c13] transition hover:bg-fuchsia-200 hover:shadow-[0_0_28px_rgba(244,114,182,0.28)] disabled:cursor-not-allowed disabled:bg-fuchsia-300/40"
                :disabled="pending || !authStore.isAuthenticated"
                @click="submitRegistration"
              >
                {{ pending ? 'Фиксируем...' : 'Подтвердить участие' }}
              </button>

              <button
                data-testid="event-reset"
                type="button"
                class="w-full rounded-[0.75rem] border border-fuchsia-400/20 bg-[#07141d] px-5 py-3 text-sm font-bold text-white transition hover:border-fuchsia-300/45"
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
const heroStats = computed(() => [
  { label: 'Events', value: String(zones.value.length), hint: 'в текущей афише' },
  { label: 'Регистрация', value: 'Бесплатно', hint: 'без скрытых условий' },
  { label: 'Доступ', value: 'Для всех', hint: 'по аккаунту' }
])

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

function eventDetails(zone: CatalogZoneWithServices): string[] {
  const d = zone.details_json as any
  const parts: string[] = []
  if (d?.date) parts.push(d.date)
  if (d?.format) parts.push(d.format)
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

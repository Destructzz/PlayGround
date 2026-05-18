<template>
  <div class="min-h-screen bg-[#020c13] pb-24 pt-24 text-white">
    <div class="mx-auto flex max-w-[1400px] flex-col gap-10 px-6 sm:px-10 lg:px-12">
      <ExperienceHero
        eyebrow="Lounge zone"
        title="Приватный отдых без потери ритма"
        description="Выбери зону, дату и часовой слот — мы оставим место для твоей компании."
      />

      <!-- Loading state -->
      <ExperienceStatePanel
        v-if="catalogLoading"
        title="Загружаем лаунж-зоны"
        description="Подтягиваем актуальную расстановку и доступность мест."
        hint="loading..."
      />

      <!-- Error state -->
      <div
        v-else-if="catalogError"
        class="rounded-lg border border-orange-500/30 bg-orange-500/10 px-5 py-4 text-xs text-orange-200"
      >
        Не удалось загрузить зоны: {{ catalogError }}
      </div>

      <template v-else>
        <!-- Zones grid -->
        <section class="grid gap-6 lg:grid-cols-3">
          <ExperienceCard
            v-for="zone in zones"
            :key="zone.id"
            eyebrow="Lounge spot"
            :title="zone.name"
            :description="zone.description"
            :details="zonePerks(zone)"
            :badge="zone.capacity > 0 ? `${zone.capacity} мест` : 'Нет данных'"
            badge-tone="success"
            footer-label="Зона"
            :footer-value="zone.zone_type"
            accent="#22d3ee"
          >
            <template #cta>
              <button
                :data-testid="`lounge-zone-${zone.id}`"
                class="rounded-lg px-4 py-2 text-xs font-bold transition"
                :class="selectedZoneId === zone.id
                  ? 'bg-cyan-400 text-black shadow-[0_0_12px_rgba(34,211,238,0.2)] hover:bg-cyan-300'
                  : 'border border-white/10 bg-white/5 text-white hover:border-white/20 hover:bg-white/10'"
                @click="selectZone(zone.id)"
              >
                {{ selectedZoneId === zone.id ? 'Выбрано' : 'Выбрать' }}
              </button>
            </template>
          </ExperienceCard>
        </section>

        <!-- Booking flow -->
        <ExperienceFlowPanel
          eyebrow="Booking flow"
          title="Забронировать lounge"
          description="Выбери дату и доступный слот — место зафиксируется за твоей компанией."
          accent="#22d3ee"
        >
          <div class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
            <div class="space-y-6">

              <!-- Selected zone summary -->
              <div>
                <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.25em] text-zinc-500">Выбранная зона</p>
                <div
                  data-testid="lounge-selected-zone"
                  class="rounded-xl border border-white/10 bg-[#071724]/75 backdrop-blur-md p-5 shadow-[0_8px_30px_rgba(0,0,0,0.4)]"
                >
                  <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                    <div>
                      <h3 class="text-xl font-black text-white uppercase">{{ selectedZone?.name ?? 'Ничего не выбрано' }}</h3>
                      <p class="mt-2 text-xs leading-relaxed text-zinc-300">
                        {{ selectedZone?.description ?? 'Сначала выбери lounge-зону из каталога выше.' }}
                      </p>
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
              </div>

              <!-- Date selector — sticky browser-style tabs -->
              <div v-if="selectedZone">
                <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.25em] text-zinc-500">Дата</p>
                <div class="flex gap-1 overflow-x-auto rounded-t-lg border border-b-0 border-white/10 bg-[#020a10] px-2 pt-2">
                  <button
                    v-for="d in availableDates"
                    :key="d.value"
                    type="button"
                    class="flex-shrink-0 rounded-t-md px-4 py-2 text-xs font-bold transition-all"
                    :class="selectedDate === d.value
                      ? 'bg-[#071724]/75 text-white border border-b-0 border-white/20 -mb-px'
                      : 'text-zinc-400 hover:text-zinc-200'"
                    @click="selectDate(d.value)"
                  >
                    {{ d.label }}
                  </button>
                </div>
                <div class="rounded-b-xl border border-white/10 bg-[#071724]/75 p-4">
                  <p class="text-[10px] uppercase tracking-wider text-zinc-500">{{ selectedDate }}</p>
                </div>
              </div>

              <!-- Hourly slots -->
              <div v-if="selectedZone && selectedDate">
                <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.25em] text-zinc-500">Часовой слот</p>
                <div v-if="availabilityLoading" class="text-xs text-zinc-400">Загружаем слоты...</div>
                <div v-else class="flex flex-wrap gap-2">
                  <button
                    v-for="slot in filteredSlots"
                    :key="slot.hour"
                    :data-testid="`lounge-slot-${slot.hour}`"
                    type="button"
                    class="rounded-lg border px-4 py-2 text-xs font-bold transition"
                    :class="selectedHour === slot.hour
                      ? 'border-cyan-400 bg-cyan-400 text-black shadow-[0_0_16px_rgba(34,211,238,0.15)]'
                      : !slot.available
                        ? 'cursor-not-allowed border-white/10 bg-white/5 text-zinc-600'
                        : 'border-white/10 bg-[#020a10] text-zinc-300 hover:border-cyan-400/30 hover:text-white'"
                    :disabled="!slot.available"
                    @click="selectedHour = slot.hour"
                  >
                    <span>{{ slot.label }}</span>
                    <span class="ml-1.5 text-[10px] opacity-60">{{ slot.remaining }}м</span>
                  </button>
                </div>
              </div>

              <!-- Party size -->
              <div v-if="selectedZone">
                <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.25em] text-zinc-500">Размер компании</p>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="n in partySizeOptions"
                    :key="n"
                    :data-testid="`lounge-party-${n}`"
                    type="button"
                    class="rounded-lg border px-4 py-2 text-xs font-bold transition"
                    :class="partySize === n
                      ? 'border-cyan-400 bg-cyan-400 text-black'
                      : n > selectedSlotRemaining
                        ? 'cursor-not-allowed border-white/10 bg-white/5 text-zinc-600'
                        : 'border-white/10 bg-[#020a10] text-zinc-300 hover:border-cyan-400/30 hover:text-white'"
                    :disabled="n > selectedSlotRemaining"
                    @click="partySize = n"
                  >
                    {{ n }} {{ n === 1 ? 'гость' : n < 5 ? 'гостя' : 'гостей' }}
                  </button>
                </div>
              </div>

              <!-- Auth notice / user drawer trigger -->
              <div v-if="!authStore.isAuthenticated" class="rounded-lg border border-amber-500/30 bg-amber-500/10 px-5 py-4 text-xs text-amber-200">
                <p class="font-bold">Требуется вход</p>
                <p class="mt-1 text-amber-200/70">Войди через Google, чтобы забронировать lounge-зону. Контактные данные подтянутся автоматически.</p>
              </div>
              <button
                v-else
                type="button"
                class="flex w-full items-center gap-4 rounded-xl border border-white/10 bg-[#020a10] px-5 py-4 text-left transition hover:border-cyan-400/30"
                @click="drawerOpen = true"
              >
                <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-cyan-400/30 bg-cyan-400/10 text-sm font-black text-cyan-200">
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
                <p class="text-[9px] font-bold uppercase tracking-[0.25em] text-zinc-500">Итог</p>
                <h3 class="mt-1 text-xl font-black uppercase tracking-tight text-white">Бронь</h3>
              </div>

              <div class="space-y-3 text-xs text-zinc-300">
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Зона</span>
                  <span class="font-bold text-white uppercase">{{ selectedZone?.name ?? '—' }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Дата</span>
                  <span class="font-bold text-white">{{ selectedDate || '—' }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Слот</span>
                  <span class="font-bold text-white">{{ selectedSlotLabel }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Гости</span>
                  <span class="font-bold text-white">{{ partySize > 0 ? partySize : '—' }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/5 pb-3">
                  <span>Стоимость</span>
                  <span
                    :class="totalBookingPrice === 'Бесплатно' ? 'text-emerald-400' : 'text-cyan-400'"
                    class="font-bold"
                  >
                    {{ totalBookingPrice }}
                  </span>
                </div>
              </div>
              
              <div class="space-y-3 pt-3 border-t border-white/5">
                <div class="flex items-center justify-between gap-4">
                  <span class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Контакты брони</span>
                  <label class="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" v-model="useCustomContacts" class="rounded border-white/10 bg-black/50 text-cyan-400 focus:ring-0 focus:ring-offset-0 h-3.5 w-3.5" />
                    <span class="text-[10px] text-zinc-400">Указать другие</span>
                  </label>
                </div>
                
                <div v-if="!useCustomContacts" class="flex flex-col gap-1 text-[11px] text-white">
                  <span class="font-bold">{{ authStore.user?.name || '—' }}</span>
                  <span class="text-zinc-400">{{ authStore.user?.email || '—' }}</span>
                </div>
                
                <div v-else class="space-y-2 pt-2">
                  <input v-model="customContacts.name" type="text" placeholder="Имя" class="w-full rounded-md border border-white/10 bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-cyan-500/50" />
                  <input v-model="customContacts.email" type="email" placeholder="Email" class="w-full rounded-md border border-white/10 bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-cyan-500/50" />
                  <input v-model="customContacts.phone" type="tel" placeholder="Телефон" class="w-full rounded-md border border-white/10 bg-white/5 px-3 py-2 text-xs text-white outline-none focus:border-cyan-500/50" />
                </div>
              </div>

              <div
                v-if="validationMessage"
                data-testid="lounge-validation-message"
                class="rounded-lg border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-xs text-rose-200"
              >
                {{ validationMessage }}
              </div>

              <div
                v-if="submitError"
                data-testid="lounge-submit-error"
                class="rounded-lg border border-orange-500/30 bg-orange-500/10 px-4 py-3 text-xs text-orange-200"
              >
                {{ submitError }}
              </div>

              <div
                v-if="bookingSuccess"
                data-testid="lounge-success-message"
                class="rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-4 py-3 text-xs text-emerald-200"
              >
                Бронь зафиксирована. Ждём тебя в {{ selectedZone?.name ?? 'lounge-зоне' }} в {{ selectedSlotLabel }}.
              </div>

              <button
                data-testid="lounge-submit"
                type="button"
                class="w-full rounded-lg bg-cyan-400 px-5 py-3.5 text-xs font-bold uppercase tracking-widest text-black transition hover:bg-cyan-300 hover:shadow-[0_0_16px_rgba(34,211,238,0.15)] disabled:cursor-not-allowed disabled:bg-cyan-400/40"
                :disabled="pending || !authStore.isAuthenticated"
                @click="submitBooking"
              >
                {{ pending ? 'Подтверждаем...' : 'Подтвердить бронь' }}
              </button>

              <button
                data-testid="lounge-reset"
                type="button"
                class="w-full rounded-lg border border-white/10 bg-white/5 px-5 py-2.5 text-xs font-bold text-zinc-300 transition hover:bg-white/10 hover:text-white"
                @click="resetForm"
              >
                Очистить выбор
              </button>
            </div>
          </div>
        </ExperienceFlowPanel>
      </template>
    </div>

    <!-- User info Drawer -->
    <Teleport to="body">
      <Transition name="drawer">
        <div
          v-if="drawerOpen"
          class="fixed inset-0 z-50 flex justify-end"
          @click.self="drawerOpen = false"
        >
          <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="drawerOpen = false" />
          <div class="relative z-10 flex h-full w-full max-w-sm flex-col bg-[#050f17] shadow-2xl">
            <div class="flex items-center justify-between border-b border-white/8 px-6 py-5">
              <p class="text-xs font-black uppercase tracking-[0.3em] text-cyan-100/50">Профиль гостя</p>
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
                <div class="flex h-20 w-20 items-center justify-center rounded-full border-2 border-cyan-300/30 bg-cyan-400/10 text-3xl font-black text-cyan-200">
                  {{ authStore.user?.name?.charAt(0).toUpperCase() ?? '?' }}
                </div>
                <h2 class="mt-5 text-xl font-black text-white">{{ authStore.user?.name }}</h2>
                <p class="mt-1 text-sm text-cyan-100/50">{{ authStore.user?.email }}</p>
                <span class="mt-3 rounded-full border border-cyan-300/20 bg-cyan-400/10 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-100">
                  {{ authStore.user?.role }}
                </span>
              </div>
              <div class="mt-8 space-y-3 text-sm">
                <div class="rounded-[0.8rem] border border-white/8 bg-white/4 px-4 py-3">
                  <p class="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Имя в брони</p>
                  <p class="mt-1 font-medium text-white">{{ authStore.user?.name ?? '—' }}</p>
                </div>
                <div class="rounded-[0.8rem] border border-white/8 bg-white/4 px-4 py-3">
                  <p class="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Email</p>
                  <p class="mt-1 font-medium text-white">{{ authStore.user?.email ?? '—' }}</p>
                </div>
              </div>
              <p class="mt-6 text-center text-[11px] text-zinc-600">
                Эти данные будут автоматически прикреплены к брони. Изменить их можно в профиле.
              </p>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
    <div data-testid="lounge-client-ready" class="hidden" />
    <div v-if="route.query.state === 'error'" data-testid="lounge-error-mode" class="hidden" />
  </div>
</template>

<script setup lang="ts">
import type { CatalogZoneWithServices, LoungeSlot } from '../api/types'
import { getLoungeCatalog, getLoungeAvailability } from '~/api/catalog'
import { createBooking } from '~/api/booking'
import { useAuthStore } from '../stores/auth'

useHead({ title: 'Lounge Zone - PlayGround' })

const authStore = useAuthStore()
const route = useRoute()

// ── State ────────────────────────────────────────────────
const zones = ref<CatalogZoneWithServices[]>([])
const catalogLoading = ref(true)
const catalogError = ref('')

const selectedZoneId = ref<number | null>(null)
const selectedDate = ref('')
const selectedHour = ref<number | null>(null)
const partySize = ref(0)

const slots = ref<LoungeSlot[]>([])
const availabilityLoading = ref(false)

const drawerOpen = ref(false)
const pending = ref(false)
const bookingSuccess = ref(false)
const submitError = ref('')
const validationMessage = ref('')

const useCustomContacts = ref(false)
const customContacts = ref({ name: '', email: '', phone: '' })

// ── Computed ─────────────────────────────────────────────


const selectedZone = computed(() =>
  zones.value.find(z => z.id === selectedZoneId.value) ?? null
)

const availableDates = computed(() => {
  const result = []
  const now = new Date()
  const days = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']
  for (let i = 0; i < 7; i++) {
    const d = new Date(now)
    d.setDate(d.getDate() + i)
    const iso = d.toISOString().slice(0, 10)
    const label = i === 0 ? 'Сегодня' : i === 1 ? 'Завтра' : `${days[d.getDay()]} ${d.getDate()}`
    result.push({ value: iso, label })
  }
  return result
})

const filteredSlots = computed(() => {
  const now = new Date()
  return slots.value.filter(slot => {
    const slotTime = new Date(`${selectedDate.value}T${String(slot.hour).padStart(2, '0')}:00:00+03:00`)
    return slotTime.getTime() > now.getTime()
  })
})

const selectedSlotRemaining = computed(() => {
  if (selectedHour.value === null) return 0
  const slot = slots.value.find(s => s.hour === selectedHour.value)
  return slot?.remaining ?? 0
})

const selectedSlotLabel = computed(() => {
  if (selectedHour.value === null) return '—'
  const slot = slots.value.find(s => s.hour === selectedHour.value)
  return slot?.label ?? '—'
})

const totalBookingPrice = computed(() => {
  const svc = selectedZone.value?.services?.[0]
  if (!svc) return 'Бесплатно'
  
  const priceNum = parseFloat(svc.price)
  if (isNaN(priceNum) || priceNum <= 0) return 'Бесплатно'
  
  const durationInMinutes = 60
  const participants = partySize.value || 1
  const total = (priceNum / svc.duration) * durationInMinutes * participants
  
  const currencySymbol = svc.currency === 'RUB' ? '₽' : svc.currency
  return `${Math.floor(total)} ${currencySymbol}`
})

const partySizeOptions = computed(() => {
  const max = Math.min(selectedZone.value?.capacity ?? 8, 8)
  return Array.from({ length: max }, (_, i) => i + 1)
})

// ── Lifecycle ─────────────────────────────────────────────
onMounted(async () => {
  try {
    const data = await getLoungeCatalog()
    zones.value = data.zones ?? []
    if (zones.value.length > 0) {
      selectedZoneId.value = zones.value[0]!.id
    }
    selectedDate.value = availableDates.value[0]!.value
  } catch (e: any) {
    catalogError.value = e?.message ?? 'Неизвестная ошибка'
  } finally {
    catalogLoading.value = false
  }

  // Fetch initial availability
  if (selectedZoneId.value && selectedDate.value) {
    await fetchAvailability()
  }
})

// ── Watchers ──────────────────────────────────────────────
watch([selectedZoneId, selectedDate], async ([zoneId, date]) => {
  if (zoneId && date) {
    selectedHour.value = null
    await fetchAvailability()
  }
})

// ── Methods ───────────────────────────────────────────────
async function fetchAvailability() {
  if (!selectedZoneId.value || !selectedDate.value) return
  availabilityLoading.value = true
  try {
    const data = await getLoungeAvailability(selectedZoneId.value, selectedDate.value)
    slots.value = data.slots ?? []
  } catch {
    slots.value = []
  } finally {
    availabilityLoading.value = false
  }
}

function selectZone(id: number) {
  selectedZoneId.value = id
  selectedHour.value = null
  partySize.value = 0
  bookingSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
}

function selectDate(date: string) {
  selectedDate.value = date
  selectedHour.value = null
}

function resetForm() {
  selectedHour.value = null
  partySize.value = 0
  bookingSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
  useCustomContacts.value = false
  customContacts.value = { name: '', email: '', phone: '' }
}

function validate(): string {
  if (!selectedZone.value) return 'Выбери lounge-зону.'
  if (!selectedDate.value) return 'Выбери дату.'
  if (selectedHour.value === null) return 'Выбери часовой слот.'
  if (partySize.value < 1) return 'Укажи размер компании.'
  if (partySize.value > selectedSlotRemaining.value) return 'Размер компании превышает доступное число мест в слоте.'
  if (!authStore.isAuthenticated) return 'Для бронирования необходимо войти в аккаунт.'
  return ''
}

async function submitBooking() {
  validationMessage.value = validate()
  submitError.value = ''
  bookingSuccess.value = false
  if (validationMessage.value) return

  if (route.query.state === 'error') {
    submitError.value = 'Симуляция ошибки'
    return
  }

  const zone = selectedZone.value!
  const service = zone.services?.[0]

  // Build RFC3339 start/end from selected date + hour (MSK UTC+3)
  const startISO = `${selectedDate.value}T${String(selectedHour.value).padStart(2, '0')}:00:00+03:00`
  const endISO = `${selectedDate.value}T${String(selectedHour.value! + 1).padStart(2, '0')}:00:00+03:00`

  pending.value = true
  
  const cName = useCustomContacts.value ? customContacts.value.name : (authStore.user?.name || '')
  const cEmail = useCustomContacts.value ? customContacts.value.email : (authStore.user?.email || '')
  const cPhone = useCustomContacts.value ? customContacts.value.phone : (localStorage.getItem('playground_phone') || '')

  try {
    await createBooking({
      zone_id: zone.id,
      service_id: service?.id || 0,
      start_time: startISO,
      end_time: endISO,
      participants: partySize.value,
      status: 'created',
      contact_name: cName,
      contact_email: cEmail,
      contact_phone: cPhone
    })
    bookingSuccess.value = true
    await fetchAvailability()
  } catch (e: any) {
    submitError.value = e?.data?.error?.message ?? e?.message ?? 'Не удалось создать бронь. Попробуй ещё раз.'
  } finally {
    pending.value = false
  }
}

function zonePerks(zone: CatalogZoneWithServices): string[] {
  const perksList: string[] = []
  try {
    const d = zone.details_json as any
    if (Array.isArray(d?.perks)) {
      perksList.push(...d.perks)
    }
  } catch {}

  if (zone.services?.length) {
    const svc = zone.services[0]
    if (svc) {
      const priceNum = parseFloat(svc.price)
      if (!isNaN(priceNum) && priceNum > 0) {
        const currencySymbol = svc.currency === 'RUB' ? '₽' : svc.currency
        perksList.push(`Тариф: ${Math.floor(priceNum)} ${currencySymbol}/час (на гостя)`)
      } else {
        perksList.push('Бронь: Бесплатно')
      }
    } else {
      perksList.push('Бронь: Бесплатно')
    }
  } else {
    perksList.push('Бронь: Бесплатно')
  }
  
  return perksList
}
</script>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.25s ease;
}
.drawer-enter-active .relative,
.drawer-leave-active .relative {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}
.drawer-enter-from .relative {
  transform: translateX(100%);
}
.drawer-leave-to .relative {
  transform: translateX(100%);
}
</style>

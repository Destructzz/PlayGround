<template>
  <div class="min-h-screen bg-[#020c13] pb-24 pt-24 text-white transition-colors duration-500">
    <div class="mx-auto flex max-w-[1400px] flex-col gap-10 px-6 sm:px-10 lg:px-12">
      <div
        v-if="isClientReady"
        data-testid="event-client-ready"
        class="hidden"
      />
      <ExperienceHero
        eyebrow="Event zone"
        title="Собранный event flow до запуска backend"
        description="Отдельный публичный маршрут для афиши и регистрации: карточки событий, sold-out состояния, локальная форма участия и deterministic preview всех UX-состояний."
        :stats="heroStats"
        accent="#f472b6"
      />

      <ExperienceStatePanel
        v-if="isLoading"
        :title="mockStateLabels.loading.title"
        :description="mockStateLabels.loading.description"
        hint="state=loading"
      />

      <ExperienceStatePanel
        v-else-if="isSuccessPreview"
        :title="mockStateLabels.success.title"
        :description="mockStateLabels.success.description"
        hint="state=success"
      >
        <template #actions>
          <NuxtLink
            to="/event"
            class="rounded-full bg-white px-6 py-3 text-sm font-bold text-black transition hover:bg-zinc-200"
          >Открыть афишу</NuxtLink>
        </template>
      </ExperienceStatePanel>

      <ExperienceStatePanel
        v-else-if="isEmpty"
        :title="mockStateLabels.empty.title"
        :description="mockStateLabels.empty.description"
        hint="state=empty"
      >
        <template #actions>
          <NuxtLink
            to="/event"
            class="rounded-full bg-white px-6 py-3 text-sm font-bold text-black transition hover:bg-zinc-200"
          >Вернуть афишу</NuxtLink>
        </template>
      </ExperienceStatePanel>

      <template v-else>
        <div
          v-if="isError"
          data-testid="event-error-mode"
          class="rounded-[0.9rem] border border-orange-300/35 bg-orange-500/12 px-5 py-4 text-sm text-orange-100 shadow-[0_16px_48px_rgba(0,0,0,0.35)]"
        >
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="text-[11px] font-bold uppercase tracking-[0.35em] text-orange-100/70">
                state=error
              </p>
              <p class="mt-2 leading-7">
                {{ mockStateLabels.error.description }} В этом режиме submit намеренно падает, чтобы проверить failure-path прямо внутри event flow.
              </p>
            </div>
            <button
              type="button"
              class="rounded-[0.7rem] border border-white/15 px-4 py-2 text-sm font-bold text-white transition hover:border-white/30"
              @click="clearTransientStates"
            >
              Вернуть default-state
            </button>
          </div>
        </div>

        <section class="grid gap-6 lg:grid-cols-3">
          <ExperienceCard
            v-for="event in visibleEvents"
            :key="event.id"
            eyebrow="Event slot"
            :title="event.title"
            :description="event.description"
            :details="[`${event.dateLabel} • ${event.timeLabel}`, event.format, ...(event.speakers ?? [])]"
            :badge="event.soldOut ? 'Sold out' : `${event.remaining} мест`"
            :badge-tone="event.soldOut ? 'warning' : 'success'"
            footer-label="Формат"
            :footer-value="event.format"
            :accent="event.accent"
          >
            <template #cta>
              <button
                :data-testid="`event-card-${event.id}`"
                class="rounded-full px-5 py-2 text-sm font-bold transition"
                :class="selectedEventId === event.id ? 'bg-white text-black' : 'border border-white/15 bg-white/5 text-white hover:border-white/30'"
                :disabled="Boolean(event.soldOut)"
                @click="selectEvent(event.id)"
              >
                {{ selectedEventId === event.id ? 'Выбрано' : event.soldOut ? 'Закрыто' : 'Участвовать' }}
              </button>
            </template>
          </ExperienceCard>
        </section>

        <ExperienceFlowPanel
          eyebrow="Registration flow"
          title="Регистрация на событие"
          description="Тот же продуктовый ритм: статус выбранного события, форма участника, блокировка sold-out, pending во время mock-submit и явный success/error output."
        >
          <div class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
            <div class="space-y-6">
              <div
                data-testid="event-selected-card"
                class="rounded-[0.85rem] border border-fuchsia-400/20 bg-[#09131d] p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]"
              >
                <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <p class="text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                      Выбранное событие
                    </p>
                    <h3 class="mt-2 text-2xl font-black text-white">
                      {{ selectedEvent?.title ?? 'Сначала выбери event' }}
                    </h3>
                    <p class="mt-2 text-sm leading-7 text-zinc-300">
                      {{ selectedEvent?.description ?? 'Нужен один доступный event из афиши выше.' }}
                    </p>
                  </div>
                  <div
                    v-if="selectedEvent"
                    class="rounded-[0.7rem] border border-fuchsia-300/25 bg-[#081019] px-4 py-3 text-right"
                  >
                    <p class="text-[11px] uppercase tracking-[0.35em] text-cyan-100/45">
                      Осталось мест
                    </p>
                    <p class="mt-2 text-2xl font-black text-white">
                      {{ selectedEvent.remaining }}/{{ selectedEvent.capacity }}
                    </p>
                  </div>
                </div>
              </div>

              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Имя</label>
                  <input
                    v-model="attendeeName"
                    data-testid="event-attendee-name"
                    type="text"
                    placeholder="Например, Anna"
                    class="w-full rounded-[0.7rem] border border-fuchsia-400/20 bg-[#06131c] px-4 py-3.5 text-white placeholder:text-fuchsia-200/25 focus:border-fuchsia-300 focus:outline-none"
                  >
                </div>
                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Email</label>
                  <input
                    v-model="attendeeEmail"
                    data-testid="event-attendee-email"
                    type="email"
                    placeholder="anna@example.com"
                    class="w-full rounded-[0.7rem] border border-fuchsia-400/20 bg-[#06131c] px-4 py-3.5 text-white placeholder:text-fuchsia-200/25 focus:border-fuchsia-300 focus:outline-none"
                  >
                </div>
              </div>

              <div>
                <p class="mb-3 text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                  Формат участия
                </p>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="option in eventGuestOptions"
                    :key="option.id"
                    :data-testid="`event-attendance-${option.id}`"
                    type="button"
                    class="rounded-[0.55rem] border px-4 py-2 text-sm font-bold transition"
                    :class="attendanceMode === option.id ? 'border-fuchsia-200 bg-fuchsia-300 text-[#020c13] shadow-[0_0_24px_rgba(244,114,182,0.25)]' : 'border-fuchsia-400/20 bg-[#07141d] text-white hover:border-fuchsia-300/45'"
                    @click="attendanceMode = option.id"
                  >
                    {{ option.label }}
                  </button>
                </div>
              </div>
            </div>

            <div class="space-y-4 rounded-[0.85rem] border border-fuchsia-400/18 bg-[#09131d] p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]">
              <div>
                <p class="text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                  Регистрация
                </p>
                <h3 class="mt-2 text-2xl font-black uppercase tracking-[-0.04em] text-white">
                  Участие подтверждается локально
                </h3>
              </div>

              <div class="space-y-3 text-sm text-zinc-300">
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Событие</span>
                  <span class="font-bold text-white">{{ selectedEvent?.title ?? '—' }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Формат</span>
                  <span class="font-bold text-white">{{ selectedAttendanceLabel }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Дата</span>
                  <span class="font-bold text-white">{{ selectedEvent ? `${selectedEvent.dateLabel} • ${selectedEvent.timeLabel}` : '—' }}</span>
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
                Регистрация готова. Твоё место на {{ selectedEvent?.title ?? 'выбранный event' }} удержано в mock-режиме.
              </div>

              <button
                data-testid="event-submit"
                type="button"
                class="w-full rounded-[0.75rem] bg-fuchsia-300 px-5 py-4 text-sm font-black uppercase tracking-[0.25em] text-[#020c13] transition hover:bg-fuchsia-200 hover:shadow-[0_0_28px_rgba(244,114,182,0.28)] disabled:cursor-not-allowed disabled:bg-fuchsia-300/40"
                :disabled="pending"
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
  </div>
</template>

<script setup lang="ts">
import type { EventItem } from '../utils/experienceData'
import { eventGuestOptions, eventItems, mockStateLabels } from '../utils/experienceData'

const route = useRoute()
const router = useRouter()
const { uiState, isLoading, isEmpty, isError, isSuccessPreview } = useExperienceMockState(route)
const { pending, submit } = useMockSubmission()

useHead({
  title: 'Event Zone - PlayGround'
})

const heroStats = [
  { label: 'Events', value: `${eventItems.length}`, hint: 'в текущей mock-афише' },
  { label: 'Next drop', value: '19:30', hint: 'вечерний прайм' },
  { label: 'Access', value: 'Public', hint: 'без auth-gating' }
]

const selectedEventId = ref<string>(eventItems.find(event => !event.soldOut)?.id ?? eventItems[0]?.id ?? '')
const attendeeName = ref('')
const attendeeEmail = ref('')
const attendanceMode = ref('')
const isClientReady = ref(false)
const registrationSuccess = ref(false)
const submitError = ref('')
const validationMessage = ref('')

const visibleEvents = computed(() => (uiState.value === 'empty' ? [] : eventItems))
const selectedEvent = computed<EventItem | undefined>(() => eventItems.find(event => event.id === selectedEventId.value))
const selectedAttendanceLabel = computed(() => eventGuestOptions.find(option => option.id === attendanceMode.value)?.label ?? '—')

onMounted(() => {
  isClientReady.value = true
})

function selectEvent(eventId: string) {
  selectedEventId.value = eventId
  registrationSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
}

function clearTransientStates() {
  submitError.value = ''
  registrationSuccess.value = false
  validationMessage.value = ''
  if (route.query.state) {
    void router.push({ query: {} })
  }
}

function resetForm() {
  attendeeName.value = ''
  attendeeEmail.value = ''
  attendanceMode.value = ''
  registrationSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
}

function validateRegistration() {
  if (!selectedEvent.value) {
    return 'Сначала выбери событие из афиши.'
  }

  if (selectedEvent.value.soldOut || selectedEvent.value.remaining === 0) {
    return 'На это событие мест больше нет. Выбери другой event.'
  }

  if (!attendeeName.value.trim()) {
    return 'Укажи имя участника.'
  }

  if (!attendeeEmail.value.trim() || !attendeeEmail.value.includes('@')) {
    return 'Нужен валидный email для mock-подтверждения.'
  }

  if (!attendanceMode.value) {
    return 'Выбери формат участия.'
  }

  return ''
}

async function submitRegistration() {
  validationMessage.value = validateRegistration()
  submitError.value = ''

  if (validationMessage.value) {
    registrationSuccess.value = false
    return
  }

  await submit(() => {
    registrationSuccess.value = true
  }, uiState.value === 'error')
    .catch(() => {
      submitError.value = 'Симуляция ошибки: регистрация не сохранилась. Повтори попытку или сбрось query-state.'
      registrationSuccess.value = false
    })
}
</script>

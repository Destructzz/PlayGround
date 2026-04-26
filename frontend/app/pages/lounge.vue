<template>
  <div class="min-h-screen bg-[#020c13] pb-24 pt-24 text-white transition-colors duration-500">
    <div class="mx-auto flex max-w-[1400px] flex-col gap-10 px-6 sm:px-10 lg:px-12">
      <div
        v-if="isClientReady"
        data-testid="lounge-client-ready"
        class="hidden"
      />
      <ExperienceHero
        eyebrow="Lounge zone"
        title="Приватный отдых без потери ритма"
        description="Собрали lounge-маршрут как отдельную публичную страницу: каталог зон, проверка доступности, локальная бронь и готовые UI-состояния до подключения backend."
        :stats="heroStats"
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
            to="/lounge"
            class="rounded-full bg-white px-6 py-3 text-sm font-bold text-black transition hover:bg-zinc-200"
          >Новый сценарий</NuxtLink>
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
            to="/lounge"
            class="rounded-full bg-white px-6 py-3 text-sm font-bold text-black transition hover:bg-zinc-200"
          >Вернуть каталог</NuxtLink>
        </template>
      </ExperienceStatePanel>

      <template v-else>
        <div
          v-if="isError"
          data-testid="lounge-error-mode"
          class="rounded-[0.9rem] border border-orange-300/35 bg-orange-500/12 px-5 py-4 text-sm text-orange-100 shadow-[0_16px_48px_rgba(0,0,0,0.35)]"
        >
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="text-[11px] font-bold uppercase tracking-[0.35em] text-orange-100/70">
                state=error
              </p>
              <p class="mt-2 leading-7">
                {{ mockStateLabels.error.description }} В этом режиме submit намеренно падает, чтобы проверить failure-path прямо внутри booking flow.
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
            v-for="zone in visibleZones"
            :key="zone.id"
            eyebrow="Lounge spot"
            :title="zone.name"
            :description="zone.description"
            :details="zone.perks"
            :badge="zone.remaining > 0 ? `${zone.remaining} мест` : 'Sold out'"
            :badge-tone="zone.remaining > 0 ? (zone.premium ? 'warning' : 'success') : 'warning'"
            footer-label="Атмосфера"
            :footer-value="zone.mood"
            :accent="zone.accent"
          >
            <template #cta>
              <button
                :data-testid="`lounge-zone-${zone.id}`"
                class="rounded-full px-5 py-2 text-sm font-bold transition"
                :class="selectedZoneId === zone.id ? 'bg-white text-black' : 'border border-white/15 bg-white/5 text-white hover:border-white/30'"
                :disabled="zone.remaining === 0"
                @click="selectZone(zone.id)"
              >
                {{ selectedZoneId === zone.id ? 'Выбрано' : zone.remaining === 0 ? 'Недоступно' : 'Выбрать' }}
              </button>
            </template>
          </ExperienceCard>
        </section>

        <ExperienceFlowPanel
          eyebrow="Booking flow"
          title="Забронировать lounge"
          description="Сценарий остаётся локальным, но уже показывает disabled, validation, pending и success/error так, как это должно ощущаться в боевом интерфейсе."
        >
          <div class="grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
            <div class="space-y-6">
              <div>
                <p class="mb-3 text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                  Выбранная зона
                </p>
                <div
                  data-testid="lounge-selected-zone"
                  class="rounded-[0.85rem] border border-cyan-400/20 bg-[#091924] p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]"
                >
                  <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                    <div>
                      <h3 class="text-2xl font-black text-white">
                        {{ selectedZone?.name ?? 'Ничего не выбрано' }}
                      </h3>
                      <p class="mt-2 text-sm leading-7 text-zinc-300">
                        {{ selectedZone?.description ?? 'Сначала выбери lounge-зону из каталога выше.' }}
                      </p>
                    </div>
                    <div
                      v-if="selectedZone"
                      class="rounded-[0.7rem] border border-cyan-300/25 bg-[#061018] px-4 py-3 text-right"
                    >
                      <p class="text-[11px] uppercase tracking-[0.35em] text-cyan-100/45">
                        Доступно сейчас
                      </p>
                      <p class="mt-2 text-2xl font-black text-white">
                        {{ selectedZone.remaining }}/{{ selectedZone.capacity }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Имя гостя</label>
                  <input
                    v-model="guestName"
                    data-testid="lounge-guest-name"
                    type="text"
                    placeholder="Например, Alex"
                    class="w-full rounded-[0.7rem] border border-cyan-400/20 bg-[#06131c] px-4 py-3.5 text-white placeholder:text-cyan-200/25 focus:border-cyan-300 focus:outline-none"
                  >
                </div>
                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">Контакт</label>
                  <input
                    v-model="contactHandle"
                    data-testid="lounge-contact"
                    type="text"
                    placeholder="Telegram / phone"
                    class="w-full rounded-[0.7rem] border border-cyan-400/20 bg-[#06131c] px-4 py-3.5 text-white placeholder:text-cyan-200/25 focus:border-cyan-300 focus:outline-none"
                  >
                </div>
              </div>

              <div class="grid gap-4 xl:grid-cols-2">
                <div>
                  <p class="mb-3 text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                    Размер компании
                  </p>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="option in availablePartyOptions"
                      :key="option.id"
                      :data-testid="`lounge-party-${option.id}`"
                      type="button"
                      class="rounded-[0.55rem] border px-4 py-2 text-sm font-bold transition"
                      :class="partySize === option.id ? 'border-cyan-200 bg-cyan-300 text-[#020c13] shadow-[0_0_24px_rgba(34,211,238,0.25)]' : option.disabled ? 'cursor-not-allowed border-white/5 bg-white/5 text-zinc-500' : 'border-cyan-400/20 bg-[#07141d] text-white hover:border-cyan-300/45'"
                      :disabled="option.disabled"
                      @click="partySize = option.id"
                    >
                      {{ option.label }}
                    </button>
                  </div>
                </div>

                <div>
                  <p class="mb-3 text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                    Слот
                  </p>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="slot in loungeTimeSlots"
                      :key="slot.id"
                      :data-testid="`lounge-slot-${slot.id}`"
                      type="button"
                      class="rounded-[0.55rem] border px-4 py-2 text-sm font-bold transition"
                      :class="selectedSlot === slot.id ? 'border-cyan-200 bg-cyan-300 text-[#020c13] shadow-[0_0_24px_rgba(34,211,238,0.25)]' : slot.disabled ? 'cursor-not-allowed border-white/5 bg-white/5 text-zinc-500' : 'border-cyan-400/20 bg-[#07141d] text-white hover:border-cyan-300/45'"
                      :disabled="slot.disabled"
                      @click="selectedSlot = slot.id"
                    >
                      {{ slot.label }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="space-y-4 rounded-[0.85rem] border border-cyan-400/18 bg-[#091924] p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]">
              <div>
                <p class="text-xs font-bold uppercase tracking-[0.35em] text-cyan-100/45">
                  Итог
                </p>
                <h3 class="mt-2 text-2xl font-black uppercase tracking-[-0.04em] text-white">
                  Почти production flow
                </h3>
              </div>

              <div class="space-y-3 text-sm text-zinc-300">
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Зона</span>
                  <span class="font-bold text-white">{{ selectedZone?.name ?? '—' }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Гости</span>
                  <span class="font-bold text-white">{{ selectedPartyLabel }}</span>
                </div>
                <div class="flex items-center justify-between gap-4 border-b border-white/8 pb-3">
                  <span>Слот</span>
                  <span class="font-bold text-white">{{ selectedSlot || '—' }}</span>
                </div>
              </div>

              <div
                v-if="validationMessage"
                data-testid="lounge-validation-message"
                class="rounded-[0.7rem] border border-rose-300/35 bg-rose-500/18 px-4 py-3 text-sm text-rose-100"
              >
                {{ validationMessage }}
              </div>

              <div
                v-if="submitError"
                data-testid="lounge-submit-error"
                class="rounded-[0.7rem] border border-orange-300/35 bg-orange-500/18 px-4 py-3 text-sm text-orange-100"
              >
                {{ submitError }}
              </div>

              <div
                v-if="bookingSuccess"
                data-testid="lounge-success-message"
                class="rounded-[0.7rem] border border-emerald-300/35 bg-emerald-500/18 px-4 py-3 text-sm text-emerald-100"
              >
                Бронь зафиксирована. Мы оставили тебе {{ selectedZone?.name ?? 'lounge-зону' }} на слот {{ selectedSlot || '—' }}.
              </div>

              <button
                data-testid="lounge-submit"
                type="button"
                class="w-full rounded-[0.75rem] bg-cyan-300 px-5 py-4 text-sm font-black uppercase tracking-[0.25em] text-[#020c13] transition hover:bg-cyan-200 hover:shadow-[0_0_28px_rgba(34,211,238,0.28)] disabled:cursor-not-allowed disabled:bg-cyan-300/40"
                :disabled="pending"
                @click="submitBooking"
              >
                {{ pending ? 'Подтверждаем...' : 'Подтвердить бронь' }}
              </button>

              <button
                data-testid="lounge-reset"
                type="button"
                class="w-full rounded-[0.75rem] border border-cyan-400/20 bg-[#07141d] px-5 py-3 text-sm font-bold text-white transition hover:border-cyan-300/45"
                @click="resetForm"
              >
                Очистить выбор
              </button>
            </div>
          </div>
        </ExperienceFlowPanel>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LoungeZone } from '../utils/experienceData'
import { loungePartyOptions, loungeTimeSlots, loungeZones, mockStateLabels } from '../utils/experienceData'

const route = useRoute()
const router = useRouter()
const { uiState, isLoading, isEmpty, isError, isSuccessPreview } = useExperienceMockState(route)
const { pending, submit } = useMockSubmission()

useHead({
  title: 'Lounge Zone - PlayGround'
})

const heroStats = [
  { label: 'Zones', value: `${loungeZones.length}`, hint: 'открытые и приватные' },
  { label: 'Peak slot', value: '20:00', hint: 'самый перегруженный час' },
  { label: 'Flow mode', value: 'Mock', hint: 'без backend' }
]

const selectedZoneId = ref<string>(loungeZones.find(zone => zone.remaining > 0)?.id ?? loungeZones[0]?.id ?? '')
const selectedSlot = ref('')
const guestName = ref('')
const contactHandle = ref('')
const partySize = ref('')
const isClientReady = ref(false)
const bookingSuccess = ref(false)
const submitError = ref('')
const validationMessage = ref('')

const visibleZones = computed(() => (uiState.value === 'empty' ? [] : loungeZones))
const selectedZone = computed<LoungeZone | undefined>(() => loungeZones.find(zone => zone.id === selectedZoneId.value))
const availablePartyOptions = computed(() => loungePartyOptions.map((option) => {
  const remaining = selectedZone.value?.remaining ?? 0
  const exceedsCapacity = Number(option.id) > remaining

  return {
    ...option,
    disabled: option.disabled || exceedsCapacity
  }
}))
const selectedPartyLabel = computed(() => loungePartyOptions.find(option => option.id === partySize.value)?.label ?? '—')

onMounted(() => {
  isClientReady.value = true
})

function selectZone(zoneId: string) {
  selectedZoneId.value = zoneId
  bookingSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''

  if (selectedZone.value && Number(partySize.value) > selectedZone.value.remaining) {
    partySize.value = ''
  }
}

function clearTransientStates() {
  submitError.value = ''
  bookingSuccess.value = false
  validationMessage.value = ''
  if (route.query.state) {
    void router.push({ query: {} })
  }
}

function resetForm() {
  guestName.value = ''
  contactHandle.value = ''
  partySize.value = ''
  selectedSlot.value = ''
  bookingSuccess.value = false
  submitError.value = ''
  validationMessage.value = ''
}

function validateBooking() {
  if (!selectedZone.value) {
    return 'Сначала выбери lounge-зону.'
  }

  if (selectedZone.value.remaining === 0) {
    return 'Эта lounge-зона уже заполнена. Выбери другую секцию.'
  }

  if (!guestName.value.trim()) {
    return 'Укажи имя гостя, чтобы зафиксировать бронь.'
  }

  if (!contactHandle.value.trim()) {
    return 'Добавь контакт для подтверждения mock-брони.'
  }

  if (!partySize.value) {
    return 'Выбери размер компании.'
  }

  if (Number(partySize.value) > selectedZone.value.remaining) {
    return 'Размер компании превышает оставшуюся вместимость выбранной зоны.'
  }

  if (!selectedSlot.value) {
    return 'Выбери доступный слот времени.'
  }

  if (loungeTimeSlots.some(slot => slot.id === selectedSlot.value && slot.disabled)) {
    return 'Этот слот недоступен. Нужен другой слот.'
  }

  return ''
}

async function submitBooking() {
  validationMessage.value = validateBooking()
  submitError.value = ''

  if (validationMessage.value) {
    bookingSuccess.value = false
    return
  }

  await submit(() => {
    bookingSuccess.value = true
  }, uiState.value === 'error')
    .catch(() => {
      submitError.value = 'Симуляция ошибки: подтверждение не прошло. Попробуй ещё раз или смени сценарий query-state.'
      bookingSuccess.value = false
    })
}
</script>

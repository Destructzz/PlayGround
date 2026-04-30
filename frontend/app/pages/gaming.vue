<template>
  <div class="min-h-screen bg-[#020c13] text-white pt-24 pb-32 transition-colors duration-1000">
    <div class="max-w-[1400px] mx-auto px-6 sm:px-12 flex flex-col gap-12">
      <!-- Навигация обратно -->
      <div>
        <NuxtLink
          to="/"
          class="dynamic-text hover:opacity-75 transition-opacity flex items-center gap-2 w-fit"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          ><path d="m15 18-6-6 6-6" /></svg>
          На главную
        </NuxtLink>
      </div>

      <div class="text-center mb-4">
        <h1 class="text-5xl md:text-7xl font-black uppercase tracking-tighter text-transparent bg-clip-text bg-gradient-to-br from-white to-zinc-400 relative inline-block">
          Игровая
          <span class="absolute inset-0 bg-clip-text text-transparent opacity-50 dynamic-gradient blur-[20px] transition-all duration-700">Игровая</span>
        </h1>
      </div>

      <div
        v-if="catalogError"
        class="rounded-2xl border border-red-400/25 bg-red-500/10 px-5 py-4 text-sm text-red-100"
      >
        {{ catalogError }}
      </div>
      <div
        v-else-if="isCatalogEmpty"
        class="rounded-2xl border border-zinc-700 bg-zinc-900/40 px-5 py-4 text-sm text-zinc-300"
      >
        В backend пока нет игровых зон для этой страницы.
      </div>
    </div>

    <!-- Бегущая лента технологий -->
    <div class="w-full flex flex-col gap-4 overflow-hidden py-8 bg-zinc-950/40 border-y border-white/5 my-8 mb-16 pointer-events-none select-none">
      <div class="w-full overflow-hidden whitespace-nowrap opacity-60">
        <div class="inline-flex animate-[scroll_20s_linear_infinite]">
          <span class="text-2xl font-black tracking-widest text-zinc-500 uppercase pr-8">PREMIUM GAMING LOUNGE • HIGH REFRESH RATE • ULTIMATE PERFORMANCE • PREMIUM GAMING LOUNGE • HIGH REFRESH RATE • ULTIMATE PERFORMANCE •</span>
          <span class="text-2xl font-black tracking-widest text-zinc-500 uppercase pr-8">PREMIUM GAMING LOUNGE • HIGH REFRESH RATE • ULTIMATE PERFORMANCE • PREMIUM GAMING LOUNGE • HIGH REFRESH RATE • ULTIMATE PERFORMANCE •</span>
        </div>
      </div>
      <div class="w-full overflow-hidden whitespace-nowrap opacity-80">
        <div class="inline-flex animate-[scroll_25s_linear_infinite_reverse]">
          <span class="text-3xl font-black tracking-widest text-transparent bg-clip-text dynamic-gradient uppercase pr-8">E-SPORTS ARENA • NO DELAY • MAXIMUM FPS • E-SPORTS ARENA • NO DELAY • MAXIMUM FPS •</span>
          <span class="text-3xl font-black tracking-widest text-transparent bg-clip-text dynamic-gradient uppercase pr-8">E-SPORTS ARENA • NO DELAY • MAXIMUM FPS • E-SPORTS ARENA • NO DELAY • MAXIMUM FPS •</span>
        </div>
      </div>
      <div class="w-full overflow-hidden whitespace-nowrap opacity-60">
        <div class="inline-flex animate-[scroll_15s_linear_infinite]">
          <span class="text-2xl font-black tracking-widest text-zinc-500 uppercase pr-8">VIP ROOMS • BOOTCAMP ZONES • 24/7 OPEN • TOURNAMENTS • VIP ROOMS • BOOTCAMP ZONES • 24/7 OPEN • TOURNAMENTS •</span>
          <span class="text-2xl font-black tracking-widest text-zinc-500 uppercase pr-8">VIP ROOMS • BOOTCAMP ZONES • 24/7 OPEN • TOURNAMENTS • VIP ROOMS • BOOTCAMP ZONES • 24/7 OPEN • TOURNAMENTS •</span>
        </div>
      </div>
    </div>

    <div class="max-w-[1400px] mx-auto px-6 sm:px-12 flex flex-col">
      <!-- Селектор Вкладок (С плавным перетеканием ползунка) -->
      <section class="mb-12 flex justify-center w-full px-4">
        <div class="bg-zinc-900/60 backdrop-blur-md p-2 rounded-full border border-white/10 shadow-xl w-full max-w-2xl">
          <div
            class="relative grid w-full"
            :style="{ gridTemplateColumns: `repeat(${tabCount}, minmax(0, 1fr))` }"
          >
            <!-- Плывущий активный задний фон (ползунок) -->
            <div
              class="absolute top-0 bottom-0 rounded-full z-0 dynamic-bg dynamic-glow transition-all duration-500 ease-[cubic-bezier(0.25,1,0.5,1)]"
              :style="{ width: `${100 / tabCount}%`, transform: `translateX(${activeClass * 100}%)` }"
            />

            <!-- Кнопки -->
            <button
              v-for="(cls, index) in classDataList"
              :key="'btn-'+index"
              :class="['relative z-10 py-3.5 text-center rounded-full font-bold text-sm sm:text-base transition-colors duration-500',
                       activeClass === index ? 'text-black' : 'text-zinc-400 hover:text-white']"
              @click="activeClass = index"
            >
              {{ cls.name }}
            </button>
          </div>
        </div>
      </section>

      <!-- Плавная анимация смены контента -->
      <div class="min-h-[500px]">
        <transition
          name="smooth-slide"
          mode="out-in"
        >
          <!-- Используем :key чтобы Vue знал, когда нужно перерисовывать элемент с анимацией -->
          <div
            v-if="currentClassData"
            :key="activeClass"
            class="flex flex-col lg:flex-row gap-12 lg:gap-24 items-center"
          >
            <!-- Левый бокс -->
            <div class="w-full lg:w-1/2 bg-[#0a192f] border-l-4 dynamic-border p-10 flex flex-col justify-center relative overflow-hidden dynamic-shadow-box rounded-r-2xl transform transition-all duration-700">
              <div class="absolute -right-20 -top-20 w-64 h-64 rounded-full blur-[80px] pointer-events-none dynamic-bg opacity-30 transition-all duration-700" />
              <h2 class="text-4xl sm:text-5xl font-black uppercase mb-6 tracking-tighter text-white">
                <span class="block">{{ currentClassData.prefix }}</span>
                <span class="block">{{ currentClassData.title }}</span>
              </h2>
              <p class="text-zinc-300 text-lg relative z-10 leading-relaxed font-medium">
                {{ currentClassData.desc }}
              </p>
            </div>

            <!-- Правые линии текста (Specs) -->
            <div class="w-full lg:w-1/2 flex flex-col gap-4">
              <div
                v-for="(spec, i) in currentClassData.specs"
                :key="tempKey + i"
                class="flex justify-between items-center py-4 border-b border-zinc-800 transition-all duration-300 group cursor-default dynamic-border-hover relative overflow-hidden"
              >
                <div class="absolute inset-0 opacity-0 group-hover:opacity-10 dynamic-bg transition-opacity duration-300 -z-10" />
                <span class="text-xl font-medium text-zinc-300 group-hover:text-white transition-colors">{{ spec.title }}</span>
                <span class="font-mono tracking-wider dynamic-text transition-colors duration-700">{{ spec.value }}</span>
              </div>
            </div>
          </div>
        </transition>
      </div>

      <!-- Pricing comparison bar -->
      <section class="mt-20">
        <div
          v-if="pricingTiers.length"
          class="bg-zinc-950/60 border border-white/5 rounded-3xl p-8 flex flex-col md:flex-row items-center justify-between gap-6"
        >
          <div
            v-for="(tier, i) in pricingTiers"
            :key="'tier'+i"
            class="flex-1 text-center"
          >
            <p class="text-zinc-500 text-xs uppercase tracking-widest mb-2">
              {{ tier.name }}
            </p>
            <p class="text-3xl font-black text-white">
              {{ tier.price }}<span class="text-zinc-500 text-lg font-medium">/час</span>
            </p>
            <p class="text-zinc-600 text-xs mt-1">
              {{ tier.note }}
            </p>
          </div>
          <div
            v-for="i in Math.max(pricingTiers.length - 1, 0)"
            :key="'div'+i"
            class="hidden md:block w-[1px] h-16 bg-zinc-800"
          />
        </div>
        <div
          v-else
          class="rounded-3xl border border-zinc-800 bg-zinc-950/40 p-6 text-sm text-zinc-400"
        >
          Backend пока не вернул ни одного активного игрового тарифа.
        </div>
      </section>

      <!-- ====== BOOKING PANEL ====== -->
      <div class="mt-16">
        <div class="text-center mb-6">
          <h2 class="text-3xl font-bold text-white">
            Бронирование
          </h2>
          <p class="text-zinc-500 mt-2">
            Выбери зону, место и время
          </p>
        </div>
      </div>

      <section class="relative mt-16">
        <div class="flex justify-center mb-6">
          <button
            class="px-10 py-4 rounded-2xl font-bold text-lg border-2 transition-all duration-500 relative overflow-hidden group"
            :class="panels.v1 ? 'bg-white text-black border-white' : 'border-zinc-700 text-white hover:border-zinc-400'"
            @click="panels.v1 = !panels.v1"
          >
            <span class="relative z-10">{{ panels.v1 ? 'Закрыть' : 'Забронировать место' }}</span>
            <div class="absolute inset-0 bg-gradient-to-r dynamic-bg opacity-0 group-hover:opacity-20 transition-opacity" />
          </button>
        </div>

        <transition name="slide-down">
          <div
            v-if="panels.v1"
            class="bg-zinc-950/70 backdrop-blur-2xl border border-white/10 rounded-3xl p-8 shadow-[0_20px_60px_rgba(0,0,0,0.6)]"
          >
            <!-- Date Selector -->
            <p class="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-3">
              Выберите дату
            </p>
            <div class="flex gap-2 overflow-x-auto pb-4 mb-6 -mx-4 px-4 sm:mx-0 sm:px-0">
              <button
                v-for="(date, d) in dates"
                :key="'date'+d"
                :class="['flex-shrink-0 px-4 py-3 rounded-xl border flex flex-col items-center justify-center transition-all duration-300 min-w-[76px]',
                         selectedDate === d ? 'dynamic-bg text-black dynamic-glow border-transparent' : 'bg-zinc-900 border-zinc-800 text-zinc-400 hover:border-zinc-600 hover:text-white']"
                @click="selectedDate = d"
              >
                <span class="text-[10px] uppercase font-bold opacity-80 mb-0.5">{{ date.label }}</span>
                <span class="font-black text-sm">{{ date.date }}</span>
              </button>
            </div>

            <!-- Area Selector -->
            <p class="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-3">
              Выберите зону
            </p>
            <div class="flex gap-3 mb-8 flex-wrap">
              <button
                v-for="(zone, i) in zones"
                :key="'v1z'+i"
                :class="['px-6 py-3 rounded-xl font-bold text-sm border transition-all duration-300',
                         selectedZone === i ? 'dynamic-bg text-black dynamic-glow border-transparent' : 'border-zinc-800 text-zinc-400 hover:border-zinc-600 hover:text-white']"
                @click="selectedZone = i"
                >
                  {{ zone.name }}
                </button>
              </div>
            <p
              v-if="!zones.length"
              class="mb-8 text-sm text-zinc-500"
            >
              В backend пока нет зон с местами.
            </p>
            <!-- Places Grid -->
            <p class="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-3">
              Выберите место
            </p>
            <div
              v-if="zones[selectedZone]?.places?.length"
              class="grid grid-cols-5 sm:grid-cols-10 gap-2 mb-8"
            >
              <button
                v-for="(pc, j) in zones[selectedZone]?.places"
                :key="'v1p'+j"
                :class="['h-12 rounded-lg text-xs font-bold transition-all duration-200 border flex items-center justify-center',
                         pc.booked ? 'bg-red-950/40 border-red-900/50 text-red-400 cursor-not-allowed'
                         : selectedPlace === j ? 'dynamic-bg text-black border-transparent scale-110 dynamic-glow' : 'bg-zinc-900 border-zinc-800 text-zinc-400 hover:border-zinc-600']"
                @click="!pc.booked && (selectedPlace = j)"
                >
                  {{ pc.name }}
                </button>
            </div>
            <p
              v-else
              class="mb-8 text-sm text-zinc-500"
            >
              Для выбранной зоны backend пока не вернул ни одного игрового места.
            </p>
            <!-- Timeline -->
            <p class="text-xs font-bold text-zinc-500 uppercase tracking-widest mb-3">
              Выберите время
            </p>
            <div class="flex overflow-x-auto pb-6 pt-4 px-4 -mx-4 -mt-2">
              <div
                v-for="(hour, k) in hours"
                :key="'v1h'+k"
                :class="getHourClass(k, hour)"
                @click="toggleHour(k)"
              >
                <svg
                  v-if="hour.taken"
                  xmlns="http://www.w3.org/2000/svg"
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                ><rect
                  width="18"
                  height="11"
                  x="3"
                  y="11"
                  rx="2"
                  ry="2"
                /><path d="M7 11V7a5 5 0 0 1 10 0v4" /></svg>
                <span>{{ hour.time }}</span>
              </div>
            </div>
            <p class="text-xs text-zinc-500">
              Слоты времени пока локальные: `public/gaming` ещё не отдаёт публичную занятость мест по часам.
            </p>
          </div>
        </transition>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watchEffect } from 'vue'

type GamingService = {
  id: number
  name: string
  duration: number
  price: string
  currency: string
  description: string
  details_json?: unknown
}

type GamingPlace = {
  id: number
  label: string
  configuration_id?: number | null
  sort_order: number
  is_active: boolean
  specs?: unknown
}

type GamingZone = {
  id: number
  name: string
  zone_type: string
  zone_tag_id: number
  capacity: number
  description: string
  is_active: boolean
  details_json?: unknown
  places: GamingPlace[]
  services: GamingService[]
}

type GamingConfiguration = {
  id: number
  zone_tag_id: number
  specs_json?: unknown
}

type GamingZoneTag = {
  id: number
  name: string
}

type GamingCatalogResponse = {
  zone_tags?: GamingZoneTag[] | null
  zones?: GamingZone[] | null
  configurations?: GamingConfiguration[] | null
}

type ClassSpec = {
  title: string
  value: string
}

type ClassData = {
  name: string
  prefix: string
  title: string
  desc: string
  hexColor: string
  specs: ClassSpec[]
}

type ZoneView = {
  name: string
  places: Array<{
    id: number
    name: string
    booked: boolean
  }>
}

type PricingTier = {
  name: string
  price: string
  note: string
}

useHead({
  title: 'Игровая зона - PlayGround'
})

const { data: catalog, error } = await useAsyncData('gaming-catalog', () => $fetch<GamingCatalogResponse>('/api/v1/public/gaming'))

const zoneColors = ['#10b981', '#22d3ee', '#f97316', '#a855f7', '#f43f5e']

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function prettifyKey(value: string) {
  return value
    .replace(/[_-]+/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
    .replace(/^./, (letter) => letter.toUpperCase())
}

function stringifyValue(value: unknown): string {
  if (typeof value === 'string') {
    return value
  }

  if (typeof value === 'number' || typeof value === 'boolean') {
    return String(value)
  }

  if (Array.isArray(value)) {
    return value.map((item) => stringifyValue(item)).filter(Boolean).join(', ')
  }

  if (isRecord(value)) {
    return Object.entries(value)
      .map(([key, item]) => `${prettifyKey(key)}: ${stringifyValue(item)}`)
      .join(', ')
  }

  return ''
}

function extractSpecs(input: unknown): ClassSpec[] {
  if (Array.isArray(input)) {
    return input.flatMap((item, index) => {
      if (isRecord(item)) {
        if (typeof item.title === 'string' && item.value !== undefined) {
          return [{ title: item.title, value: stringifyValue(item.value) }]
        }

        const firstPair = Object.entries(item)[0]
        if (!firstPair) {
          return []
        }

        return [{ title: prettifyKey(firstPair[0]), value: stringifyValue(firstPair[1]) }]
      }

      const value = stringifyValue(item)
      return value ? [{ title: `Параметр ${index + 1}`, value }] : []
    }).filter((spec) => spec.value)
  }

  if (isRecord(input)) {
    return Object.entries(input)
      .map(([key, value]) => ({ title: prettifyKey(key), value: stringifyValue(value) }))
      .filter((spec) => spec.value)
  }

  return []
}

function formatPrice(price: string, currency: string) {
  const suffix = currency === 'RUB' ? '₽' : ` ${currency}`
  return `${price}${suffix}`
}

const rawZones = computed(() => catalog.value?.zones ?? [])
const zoneTags = computed(() => catalog.value?.zone_tags ?? [])
const rawConfigurations = computed(() => catalog.value?.configurations ?? [])
const catalogError = computed(() => error.value ? 'Не удалось загрузить игровые зоны из backend.' : '')
const isCatalogEmpty = computed(() => !catalogError.value && zoneTags.value.length === 0)
const tabCount = computed(() => Math.max(zoneTags.value.length, 1))

const tagOrder = computed(() => new Map(zoneTags.value.map((tag, index) => [tag.id, index])))
const userZones = computed(() => [...rawZones.value].sort((left, right) => {
  const leftOrder = tagOrder.value.get(left.zone_tag_id) ?? Number.MAX_SAFE_INTEGER
  const rightOrder = tagOrder.value.get(right.zone_tag_id) ?? Number.MAX_SAFE_INTEGER

  if (leftOrder !== rightOrder) {
    return leftOrder - rightOrder
  }

  return left.id - right.id
}))

const classDataList = computed<ClassData[]>(() => zoneTags.value.map((tag, index) => {
  const tagZones = userZones.value.filter((zone) => zone.zone_tag_id === tag.id)
  const representativeZone = tagZones[0]
  const details = isRecord(representativeZone?.details_json) ? representativeZone.details_json : {}
  const config = rawConfigurations.value.find((item) => item.zone_tag_id === tag.id)
  const configSpecs = extractSpecs(config?.specs_json)
  const placeSpecs = extractSpecs(representativeZone?.places.find((place) => place.specs)?.specs)
  const detailsSpecs = extractSpecs(details.specs)
  const service = representativeZone?.services[0]
  const specs = (configSpecs.length ? configSpecs : placeSpecs.length ? placeSpecs : detailsSpecs).slice(0, 5)

  return {
    name: tag.name,
    prefix: typeof details.prefix === 'string' ? details.prefix : `ТЕГ ${index + 1}`,
    title: typeof details.title === 'string' ? details.title : representativeZone?.name || tag.name,
    desc: representativeZone?.description || service?.description || 'Игровая зона доступна для бронирования через backend-каталог.',
    hexColor: typeof details.hexColor === 'string' ? details.hexColor : zoneColors[index % zoneColors.length]!,
    specs: specs.length ? specs : [
      { title: 'Вместимость', value: String(representativeZone?.capacity ?? 0) },
      { title: 'Игровых мест', value: String(representativeZone?.places.length ?? 0) },
      { title: 'Тарифов', value: String(representativeZone?.services.length ?? 0) },
      { title: 'Конфигураций', value: String(rawConfigurations.value.filter((item) => item.zone_tag_id === tag.id).length) },
      { title: 'Статус', value: representativeZone?.is_active ? 'Активна' : 'Скрыта' }
    ]
  }
}))

const activeClass = ref(0)
const currentClassData = computed(() => classDataList.value[activeClass.value])
const activeColor = computed(() => currentClassData.value?.hexColor || '#10b981')
const tempKey = computed(() => `spec-list-${activeClass.value}-`)

const pricingTiers = computed<PricingTier[]>(() => zoneTags.value.flatMap((tag) => {
  const zone = userZones.value.find((item) => item.zone_tag_id === tag.id)
  const service = zone?.services[0]

  if (!zone || !service) {
    return []
  }

  return [{
    name: tag.name,
    price: formatPrice(service.price, service.currency),
    note: service.description || `${zone.name} • ${service.duration} мин`
  }]
}))

const panels = reactive({ v1: false, v2: false, v3: false, v4: false, v5: false })
const selectedDate = ref(0)
const selectedZone = ref(0)
const selectedPlace = ref(-1)
const selectedHours = ref<number[]>([])

const zones = computed<ZoneView[]>(() => userZones.value.map((zone) => ({
  name: zone.name,
  places: [...zone.places]
    .sort((left, right) => left.sort_order - right.sort_order)
    .map((place) => ({
      id: place.id,
      name: place.label,
      booked: false
    }))
})))

const hours = Array.from({ length: 12 }, (_, i) => ({
  time: `${10 + i}:00`,
  label: i < 6 ? 'День' : 'Вечер',
  taken: [2, 3, 8].includes(i)
}))

const dates = Array.from({ length: 7 }, (_, i) => {
  const d = new Date()
  d.setDate(d.getDate() + i)

  let label = ''
  if (i === 0) label = 'Сегодня'
  else if (i === 1) label = 'Завтра'
  else {
    const days = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']
    label = days[d.getDay()] || ''
  }

  return {
    label,
    date: d.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
  }
})

function toggleHour(index: number) {
  if (hours[index]?.taken) return
  const pos = selectedHours.value.indexOf(index)
  if (pos >= 0) {
    selectedHours.value.splice(pos, 1)
  } else {
    selectedHours.value.push(index)
  }
}

function getHourClass(k: number, hour: (typeof hours)[number]) {
  const isSelected = selectedHours.value.includes(k)
  const isPrevSelected = isSelected && selectedHours.value.includes(k - 1)
  const isNextSelected = isSelected && selectedHours.value.includes(k + 1)

  const classes = ['flex-shrink-0 w-16 h-20 flex flex-col items-center justify-center text-xs font-bold border cursor-pointer transition-all duration-200 gap-0.5']

  if (isSelected) {
    classes.push('dynamic-bg text-black border-y-transparent dynamic-glow relative z-10')

    if (isPrevSelected) {
      classes.push('border-l-transparent rounded-l-none')
    } else {
      classes.push('border-l-transparent rounded-l-xl')
    }

    if (isNextSelected) {
      classes.push('border-r-transparent rounded-r-none mr-0')
    } else {
      classes.push('border-r-transparent rounded-r-xl mr-1.5')
    }
  } else {
    classes.push('rounded-xl mr-1.5 relative z-0')
    if (hour.taken) {
      classes.push('bg-blue-950/40 border-blue-500/30 text-blue-400 cursor-not-allowed opacity-80')
    } else {
      classes.push('bg-zinc-900 border-zinc-800 text-zinc-400 hover:border-zinc-500 hover:bg-zinc-800')
    }
  }

  return classes.join(' ')
}

watchEffect(() => {
  if (classDataList.value.length === 0) {
    activeClass.value = 0
    return
  }

  if (activeClass.value >= classDataList.value.length) {
    activeClass.value = 0
  }
})

watchEffect(() => {
  if (zones.value.length === 0) {
    selectedZone.value = 0
    selectedPlace.value = -1
    return
  }

  if (selectedZone.value >= zones.value.length) {
    selectedZone.value = 0
  }

  const places = zones.value[selectedZone.value]?.places ?? []
  if (selectedPlace.value >= places.length) {
    selectedPlace.value = -1
  }
})
</script>

<style scoped>
/* Динамические переменные цвета на основе выбранного класса */
.dynamic-text { color: v-bind(activeColor); }
.dynamic-bg { background-color: v-bind(activeColor); }
.dynamic-border { border-color: v-bind(activeColor); }
.dynamic-glow { box-shadow: 0 0 20px v-bind(activeColor); }
.dynamic-shadow-box { box-shadow: 0 0 40px v-bind('activeColor + "1a"'); }
.dynamic-border-hover:hover { border-color: v-bind(activeColor); }
.dynamic-gradient { background-image: linear-gradient(to right, v-bind(activeColor), #ffffff); }

/* Анимация смены контента (Smooth Slide & Fade) */
.smooth-slide-enter-active,
.smooth-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.25, 1, 0.5, 1);
}
.smooth-slide-enter-from {
  opacity: 0;
  transform: translateY(15px);
}
.smooth-slide-leave-to {
  opacity: 0;
  transform: translateY(-15px);
}

/* Slide-down panel */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.5s cubic-bezier(0.25, 1, 0.5, 1);
}
.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* Fade overlay */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Side drawer */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.4s cubic-bezier(0.25, 1, 0.5, 1);
}
.slide-right-enter-from,
.slide-right-leave-to {
  transform: translateX(100%);
}

/* Bottom sheet */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.4s cubic-bezier(0.25, 1, 0.5, 1);
}
.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
}

/* Анимация строк с задержкой */
.smooth-slide-enter-active :deep(.dynamic-text),
.smooth-slide-enter-active :deep(.group) {
  animation: slideInUp 0.5s both cubic-bezier(0.25, 1, 0.5, 1);
}
.smooth-slide-enter-active :deep(.group:nth-child(1)) { animation-delay: 0.05s; }
.smooth-slide-enter-active :deep(.group:nth-child(2)) { animation-delay: 0.1s; }
.smooth-slide-enter-active :deep(.group:nth-child(3)) { animation-delay: 0.15s; }
.smooth-slide-enter-active :deep(.group:nth-child(4)) { animation-delay: 0.2s; }
.smooth-slide-enter-active :deep(.group:nth-child(5)) { animation-delay: 0.25s; }

@keyframes slideInUp {
  0% { opacity: 0; transform: translateY(15px); }
  100% { opacity: 1; transform: translateY(0); }
}
</style>

<!-- Unscoped: keyframes & animation utilities must NOT be scoped or Vue breaks the escaped selectors -->
<style>
@keyframes scroll {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}
@keyframes scroll_reverse {
  0% { transform: translateX(-50%); }
  100% { transform: translateX(0); }
}
.animate-\[scroll_20s_linear_infinite\] {
  animation: scroll 20s linear infinite;
}
.animate-\[scroll_25s_linear_infinite_reverse\] {
  animation: scroll_reverse 25s linear infinite;
}
.animate-\[scroll_15s_linear_infinite\] {
  animation: scroll 15s linear infinite;
}
</style>

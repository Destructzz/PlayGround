<script setup lang="ts">
import { useApi } from '~/composables/useApi'

const { get } = useApi()

const stats = ref({ zones: 0, services: 0, bookings: 0, staff: 0, payments: 0 })

async function fetchStats() {
  try {
    const [z, s, b, st, p] = await Promise.allSettled([
      get('/zone'),
      get('/service'),
      get('/booking'),
      get('/staff'),
      get('/payment')
    ])
    stats.value.zones = z.status === 'fulfilled' ? (z.value?.zones?.length || 0) : 0
    stats.value.services = s.status === 'fulfilled' ? (s.value?.services?.length || 0) : 0
    stats.value.bookings = b.status === 'fulfilled' ? (b.value?.bookings?.length || 0) : 0
    stats.value.staff = st.status === 'fulfilled' ? (st.value?.staff?.length || 0) : 0
    stats.value.payments = p.status === 'fulfilled' ? (p.value?.payments?.length || 0) : 0
  } catch (e) { console.error(e) }
}

const highlights = computed(() => [
  { title: 'Зоны', value: stats.value.zones, to: '/zones' },
  { title: 'Услуги', value: stats.value.services, to: '/services' },
  { title: 'Бронирования', value: stats.value.bookings, to: '/bookings' },
  { title: 'Сотрудники', value: stats.value.staff, to: '/staff' },
  { title: 'Платежи', value: stats.value.payments, to: '/payments' }
])

const actions = [
  { title: 'Зоны', description: 'Создание, редактирование и удаление игровых зон.', to: '/zones' },
  { title: 'Услуги', description: 'Каталог услуг с ценами и привязкой к зонам.', to: '/services' },
  { title: 'Бронирования', description: 'Управление бронями и статусами.', to: '/bookings' },
  { title: 'Сотрудники', description: 'Кадровый учёт: должности, найм, контакты.', to: '/staff' },
  { title: 'Платежи', description: 'Финансы: суммы, методы оплаты, чеки.', to: '/payments' },
  { title: 'Админ-панель', description: 'Быстрый доступ к административным функциям.', to: '/admin' }
]

onMounted(fetchStats)
</script>

<template>
  <section class="space-y-10">
    <div class="grid items-start gap-6 lg:grid-cols-[1.15fr,0.85fr]">
      <div class="space-y-5">
        <p class="section-title">PlayGround · АИС досугового центра</p>
        <h1 class="text-3xl font-bold text-white sm:text-4xl">
          Управление бронированием, персоналом и финансами
        </h1>
        <p class="text-gray-400">
          Полноценный CRUD по всем сущностям: зоны, услуги, бронирования, сотрудники, платежи. Данные загружаются с бэкенда в реальном времени.
        </p>
        <div class="flex flex-wrap gap-3">
          <NuxtLink to="/bookings" class="solid-button">Бронирования</NuxtLink>
          <NuxtLink to="/zones" class="ghost-button">Зоны</NuxtLink>
          <NuxtLink to="/payments" class="ghost-button">Платежи</NuxtLink>
        </div>
      </div>
      <div class="grid gap-3 sm:grid-cols-3 lg:grid-cols-3">
        <NuxtLink
          v-for="item in highlights"
          :key="item.title"
          :to="item.to"
          class="glass p-4 hover:bg-white/8 transition cursor-pointer"
        >
          <div class="text-sm text-gray-400">{{ item.title }}</div>
          <div class="mt-2 text-2xl font-bold text-white">{{ item.value }}</div>
        </NuxtLink>
      </div>
    </div>

    <div class="card-grid">
      <NuxtLink
        v-for="action in actions"
        :key="action.title"
        :to="action.to"
        class="glass border-white/10 bg-white/5 p-4 hover:bg-white/8 transition block"
      >
        <div class="text-sm font-semibold text-white flex items-center gap-2">
          {{ action.title }}
        </div>
        <p class="mt-2 text-sm text-gray-400">{{ action.description }}</p>
        <span class="mt-3 inline-flex text-sm font-semibold text-gray-200 hover:text-white">
          Перейти →
        </span>
      </NuxtLink>
    </div>
  </section>
</template>

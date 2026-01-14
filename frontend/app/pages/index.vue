<script setup lang="ts">
const highlights = [
  { title: 'Зоны', value: '12', note: 'VR, консоли, настолки, event' },
  { title: 'Мероприятия', value: 'Ежедневно', note: 'ДР, турниры, корпоративы' },
  { title: 'Онлайн-бронирование', value: '24/7', note: 'Самостоятельно или через сотрудника' }
]

const actions = [
  {
    title: 'Новая бронь',
    description: 'Выберите услугу, зону и слот. Проверка доступности позже.',
    to: '/bookings/new'
  },
  {
    title: 'Каталог услуг',
    description: 'Статический список VR, консоли, настолки, мероприятия.',
    to: '/services'
  },
  {
    title: 'Кабинет бронирований',
    description: 'Управление текущими бронями и статусами.',
    to: '/bookings'
  }
]
</script>

<template>
  <section class="space-y-10">
    <div class="grid items-start gap-6 lg:grid-cols-[1.15fr,0.85fr]">
      <div class="space-y-5">
        <p class="section-title">PlayGround · АИС досугового центра</p>
        <h1 class="text-3xl font-bold text-white sm:text-4xl">
          Тёмный премиальный каркас интерфейса с бронированием и расписанием
        </h1>
        <p class="text-gray-400">
          Статические экраны для клиента, сотрудника и администратора. Фокус на визуальной системе:
          стеклянные панели, мягкие ховеры, капс-заголовки и лаконичные списки.
        </p>
        <div class="flex flex-wrap gap-3">
          <NuxtLink to="/bookings/new" class="solid-button">Создать бронь</NuxtLink>
          <NuxtLink to="/services" class="ghost-button">Каталог услуг</NuxtLink>
        </div>
      </div>
      <div class="grid gap-3 sm:grid-cols-3 sm:gap-4">
        <div
          v-for="item in highlights"
          :key="item.title"
          class="glass p-4"
        >
          <div class="text-sm text-gray-400">{{ item.title }}</div>
          <div class="mt-2 text-2xl font-bold text-white">{{ item.value }}</div>
          <div class="mt-1 text-sm text-gray-500">{{ item.note }}</div>
        </div>
      </div>
    </div>

    <div class="card-grid">
      <div
        v-for="action in actions"
        :key="action.title"
        class="glass border-white/10 bg-white/5 p-4 hover:bg-white/8 transition"
      >
        <div class="text-sm font-semibold text-white flex items-center gap-2">
          {{ action.title }}
          <span class="badge">NEW</span>
        </div>
        <p class="mt-2 text-sm text-gray-400">{{ action.description }}</p>
        <NuxtLink :to="action.to" class="mt-3 inline-flex text-sm font-semibold text-gray-200 hover:text-white">
          Перейти →
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const panels = [
  { title: 'Зоны', desc: 'Создание, редактирование и удаление зон.', to: '/zones' },
  { title: 'Услуги', desc: 'Каталог услуг, цены, длительность.', to: '/services' },
  { title: 'Бронирования', desc: 'Контроль статусов и управление.', to: '/bookings' },
  { title: 'Сотрудники', desc: 'Кадровый учёт персонала.', to: '/staff' },
  { title: 'Платежи', desc: 'Финансы и статусы платежей.', to: '/payments' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2">
      <p class="section-title">Админ</p>
      <h1 class="text-2xl font-bold text-white">Административная панель</h1>
      <p class="text-gray-400">Быстрый доступ ко всем разделам управления.</p>
    </div>

    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <NuxtLink
        v-for="panel in panels"
        :key="panel.title"
        :to="panel.to"
        class="glass p-5 hover:bg-white/8 transition block"
      >
        <h2 class="text-lg font-semibold text-white">{{ panel.title }}</h2>
        <p class="mt-1 text-sm text-gray-400">{{ panel.desc }}</p>
        <span class="mt-3 inline-flex text-sm font-semibold text-gray-200">Открыть →</span>
      </NuxtLink>
    </div>
  </section>
</template>

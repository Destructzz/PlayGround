<script setup lang="ts">
const panels = [
  { title: 'Зоны и помещения', desc: 'Редактирование названий, вместимости и статуса доступности.', to: '/admin/zones' },
  { title: 'Услуги', desc: 'Каталог услуг, цены, длительность.', to: '/admin/services' },
  { title: 'Пользователи', desc: 'Роли, активность, блокировки.', to: '/admin/users' },
  { title: 'Бронирования', desc: 'Контроль статусов и отмен.', to: '/bookings' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2">
      <p class="section-title">Админ</p>
      <h1 class="text-2xl font-bold text-white">Административная панель</h1>
      <p class="text-gray-400">Демо-экраны для настройки зон, услуг и пользователей.</p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <article
        v-for="panel in panels"
        :key="panel.title"
        class="glass p-4"
      >
        <h2 class="text-lg font-semibold text-white flex items-center gap-2">
          {{ panel.title }}
          <span class="badge" v-if="panel.title === 'Пользователи'">Admin</span>
        </h2>
        <p class="mt-1 text-sm text-gray-400">{{ panel.desc }}</p>
        <NuxtLink :to="panel.to" class="mt-3 inline-flex text-sm font-semibold text-gray-200 hover:text-white">Открыть →</NuxtLink>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
const events = [
  { title: 'VR-сессия', zone: 'VR-зона', time: 'Сегодня, 12:00–13:00', status: 'Идёт' },
  { title: 'День рождения', zone: 'Праздничный зал', time: 'Сегодня, 14:00–17:00', status: 'Запланировано' },
  { title: 'Кибер-турнир', zone: 'Киберзона', time: 'Завтра, 11:00–15:00', status: 'Запланировано' },
  { title: 'Настолки', zone: 'Комната настолок', time: 'Завтра, 18:00–20:00', status: 'Окно' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Расписание</p>
        <h1 class="text-2xl font-bold text-white">Слоты по зонам</h1>
        <p class="text-gray-400">Данные статические. Формат карточек — стекло с мягким контрастом.</p>
      </div>
      <NuxtLink to="/bookings/new" class="solid-button">Создать бронь</NuxtLink>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <article
        v-for="event in events"
        :key="event.title + event.time"
        class="glass p-4"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-white">{{ event.title }}</h2>
            <p class="text-sm text-gray-400">{{ event.zone }}</p>
          </div>
          <span class="badge">{{ event.status }}</span>
        </div>
        <p class="mt-2 text-sm text-gray-300">{{ event.time }}</p>
      </article>
    </div>
  </section>
</template>

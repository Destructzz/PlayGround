<script setup lang="ts">
const zones = [
  { name: 'VR-зона', type: 'game', capacity: 4, status: 'Доступна', description: 'VR-станции с контроллерами, затемнённое помещение.' },
  { name: 'Зал консолей №1', type: 'game', capacity: 6, status: 'Доступна', description: 'PS5 / Xbox, кооперативные игры.' },
  { name: 'Зал консолей №2', type: 'game', capacity: 6, status: 'Занята', description: 'Параллельные сессии и мини-турниры.' },
  { name: 'Комната настолок', type: 'lounge', capacity: 8, status: 'Доступна', description: 'Настольные и ролевые игры.' },
  { name: 'Киберспортивный зал', type: 'event', capacity: 10, status: 'Доступна', description: 'Командные тренировки, LAN.' },
  { name: 'Праздничный зал', type: 'event', capacity: 15, status: 'Бронь до 18:00', description: 'ДР и закрытые мероприятия.' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2">
      <p class="section-title">Ресурсы</p>
      <h1 class="text-2xl font-bold text-white">Игровые зоны и помещения</h1>
      <p class="text-gray-400">Вместимость, тип и статус доступности. Демо-данные.</p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <article
        v-for="zone in zones"
        :key="zone.name"
        class="glass p-4"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-white">{{ zone.name }}</h2>
            <p class="text-sm text-gray-400">Тип: {{ zone.type }} · Вместимость: {{ zone.capacity }}</p>
          </div>
          <span class="badge">{{ zone.status }}</span>
        </div>
        <p class="mt-2 text-sm text-gray-400">{{ zone.description }}</p>
        <div class="mt-3 flex items-center gap-3 text-sm">
          <NuxtLink to="/bookings/new" class="solid-button">Забронировать</NuxtLink>
          <NuxtLink to="/services" class="nav-link">Услуги в зоне</NuxtLink>
        </div>
      </article>
    </div>
  </section>
</template>

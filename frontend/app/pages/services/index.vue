<script setup lang="ts">
const services = [
  { name: 'VR-сессия', zone: 'VR-зона', duration: '60 мин', price: '1 500 ₽', description: 'Иммерсивные VR-игры, индивидуальные и парные.' },
  { name: 'Консоли', zone: 'Зал консолей', duration: '90 мин', price: '1 200 ₽', description: 'PS5 / Xbox, кооперативные и PvP сессии.' },
  { name: 'Настольные игры', zone: 'Комната настолок', duration: '120 мин', price: '900 ₽', description: 'Каталог настолок, аренда по времени.' },
  { name: 'Киберспорт', zone: 'Киберзона', duration: '120 мин', price: '1 800 ₽', description: 'Командные турниры и тренировки.' },
  { name: 'День рождения', zone: 'Праздничный зал', duration: '180 мин', price: '6 500 ₽', description: 'Пакет для праздника, аниматор по запросу.' },
  { name: 'VIP-комната', zone: 'VIP-зал', duration: '120 мин', price: '4 000 ₽', description: 'Закрытая комната, премиум-сервис.' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2">
      <p class="section-title">Каталог</p>
      <h1 class="text-2xl font-bold text-white">Услуги PlayGround</h1>
      <p class="text-gray-400">Статический перечень предложений. Цены и параметры — демо.</p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <article
        v-for="service in services"
        :key="service.name"
        class="glass p-4"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-white">{{ service.name }}</h2>
            <p class="text-sm text-gray-400">{{ service.zone }} · {{ service.duration }}</p>
          </div>
          <div class="text-base font-semibold text-gray-100">{{ service.price }}</div>
        </div>
        <p class="mt-2 text-sm text-gray-400">{{ service.description }}</p>
        <div class="mt-3 flex items-center gap-3 text-sm">
          <NuxtLink to="/bookings/new" class="solid-button">Забронировать</NuxtLink>
          <NuxtLink to="/zones" class="nav-link">Посмотреть зону</NuxtLink>
        </div>
      </article>
    </div>
  </section>
</template>

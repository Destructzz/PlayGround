<script setup lang="ts">
const services = [
  { name: 'VR-сессия', price: '1 500 ₽', duration: '60 мин', zone: 'VR-зона', status: 'Активна' },
  { name: 'День рождения', price: '6 500 ₽', duration: '180 мин', zone: 'Праздничный зал', status: 'Активна' },
  { name: 'VIP-комната', price: '4 000 ₽', duration: '120 мин', zone: 'VIP-зал', status: 'Скрыта' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Админ · Услуги</p>
        <h1 class="text-2xl font-bold text-white">Каталог услуг</h1>
        <p class="text-gray-400">Редактирование длительности, цены и статуса. Демо-значения.</p>
      </div>
      <button class="ghost-button">Добавить услугу</button>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/10 bg-white/5 shadow-[0_20px_60px_rgba(0,0,0,0.45)]">
      <table class="min-w-full divide-y divide-white/5 text-sm text-gray-100">
        <thead class="bg-white/5 text-xs uppercase text-gray-400 tracking-[0.14em]">
          <tr>
            <th class="px-4 py-3 text-left">Название</th>
            <th class="px-4 py-3 text-left">Зона</th>
            <th class="px-4 py-3 text-left">Длительность</th>
            <th class="px-4 py-3 text-left">Цена</th>
            <th class="px-4 py-3 text-left">Статус</th>
            <th class="px-4 py-3 text-left">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="service in services" :key="service.name" class="hover:bg-white/5 transition">
            <td class="px-4 py-3 font-medium text-white">{{ service.name }}</td>
            <td class="px-4 py-3 text-gray-300">{{ service.zone }}</td>
            <td class="px-4 py-3 text-gray-300">{{ service.duration }}</td>
            <td class="px-4 py-3 text-gray-100 font-semibold">{{ service.price }}</td>
            <td class="px-4 py-3 text-gray-300">{{ service.status }}</td>
            <td class="px-4 py-3 text-gray-300">
              <button class="nav-link">Редактировать</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

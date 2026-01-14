<script setup lang="ts">
const zones = [
  { name: 'VR-зона', capacity: 4, type: 'game', status: 'Активна' },
  { name: 'Зал консолей №1', capacity: 6, type: 'game', status: 'Активна' },
  { name: 'Праздничный зал', capacity: 15, type: 'event', status: 'На ремонте' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Админ · Зоны</p>
        <h1 class="text-2xl font-bold text-white">Зоны и помещения</h1>
        <p class="text-gray-400">Список зон для администрирования. Данные статичны.</p>
      </div>
      <button class="ghost-button">Добавить зону</button>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/10 bg-white/5 shadow-[0_20px_60px_rgba(0,0,0,0.45)]">
      <table class="min-w-full divide-y divide-white/5 text-sm text-gray-100">
        <thead class="bg-white/5 text-xs uppercase text-gray-400 tracking-[0.14em]">
          <tr>
            <th class="px-4 py-3 text-left">Название</th>
            <th class="px-4 py-3 text-left">Тип</th>
            <th class="px-4 py-3 text-left">Вместимость</th>
            <th class="px-4 py-3 text-left">Статус</th>
            <th class="px-4 py-3 text-left">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="zone in zones" :key="zone.name" class="hover:bg-white/5 transition">
            <td class="px-4 py-3 font-medium text-white">{{ zone.name }}</td>
            <td class="px-4 py-3 text-gray-300">{{ zone.type }}</td>
            <td class="px-4 py-3 text-gray-300">{{ zone.capacity }}</td>
            <td class="px-4 py-3 text-gray-300">{{ zone.status }}</td>
            <td class="px-4 py-3 text-gray-300">
              <button class="nav-link">Редактировать</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

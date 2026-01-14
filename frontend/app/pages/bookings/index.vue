<script setup lang="ts">
const bookings = [
  { code: '#B-1021', user: 'Иван Петров', service: 'VR-сессия', zone: 'VR-зона', start: '12 фев, 14:00', end: '15:00', participants: 2, status: 'confirmed', price: '1 500 ₽' },
  { code: '#B-1022', user: 'Мария Смирнова', service: 'День рождения', zone: 'Праздничный зал', start: '13 фев, 11:00', end: '14:00', participants: 10, status: 'created', price: '6 500 ₽' },
  { code: '#B-1023', user: 'Алексей Ким', service: 'Консоли', zone: 'Зал консолей №1', start: '12 фев, 16:00', end: '17:30', participants: 3, status: 'canceled', price: '—' }
]

const statusMap: Record<string, string> = {
  confirmed: 'text-green-700 bg-green-50',
  created: 'text-amber-700 bg-amber-50',
  canceled: 'text-gray-600 bg-gray-100',
  completed: 'text-blue-700 bg-blue-50'
}
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Бронирования</p>
        <h1 class="text-2xl font-bold text-white">Список броней (демо)</h1>
        <p class="text-gray-400">Управление статусами пока статичное. Формат строк — стеклянные панели.</p>
      </div>
      <NuxtLink to="/bookings/new" class="solid-button">Новое бронирование</NuxtLink>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/10 bg-white/5 shadow-[0_20px_60px_rgba(0,0,0,0.45)]">
      <table class="min-w-full divide-y divide-white/5 text-sm text-gray-100">
        <thead class="bg-white/5 text-xs uppercase text-gray-400 tracking-[0.14em]">
          <tr>
            <th class="px-4 py-3 text-left">Код</th>
            <th class="px-4 py-3 text-left">Клиент</th>
            <th class="px-4 py-3 text-left">Услуга</th>
            <th class="px-4 py-3 text-left">Зона</th>
            <th class="px-4 py-3 text-left">Время</th>
            <th class="px-4 py-3 text-left">Участники</th>
            <th class="px-4 py-3 text-left">Статус</th>
            <th class="px-4 py-3 text-left">Стоимость</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="item in bookings" :key="item.code" class="hover:bg-white/5 transition">
            <td class="px-4 py-3 font-semibold text-white">{{ item.code }}</td>
            <td class="px-4 py-3 text-gray-200">{{ item.user }}</td>
            <td class="px-4 py-3 text-gray-300">{{ item.service }}</td>
            <td class="px-4 py-3 text-gray-300">{{ item.zone }}</td>
            <td class="px-4 py-3 text-gray-300">{{ item.start }} — {{ item.end }}</td>
            <td class="px-4 py-3 text-gray-300">{{ item.participants }}</td>
            <td class="px-4 py-3">
              <span :class="`rounded-full px-3 py-1 text-xs font-semibold ${statusMap[item.status]}`">{{ item.status }}</span>
            </td>
            <td class="px-4 py-3 font-semibold text-white">{{ item.price }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

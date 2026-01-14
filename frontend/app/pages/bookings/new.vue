<script setup lang="ts">
const services = ['VR-сессия', 'Консоли', 'Настольные игры', 'День рождения', 'VIP-комната']
const zones = ['VR-зона', 'Зал консолей №1', 'Комната настолок', 'Праздничный зал', 'VIP-зал']
const timeSlots = ['10:00', '12:00', '14:00', '16:00', '18:00']
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2">
      <p class="section-title">Новая бронь</p>
      <h1 class="text-2xl font-bold text-white">Форма-заглушка</h1>
      <p class="text-gray-400">Выбор услуги, зоны, времени и участников. Валидации и логика подключатся позже.</p>
    </div>

    <form class="grid gap-4 rounded-2xl border border-white/10 bg-white/5 p-5 shadow-[0_20px_60px_rgba(0,0,0,0.45)] backdrop-blur-xl md:grid-cols-2">
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Услуга
        <select class="input-dark">
          <option v-for="service in services" :key="service">{{ service }}</option>
        </select>
      </label>

      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Зона / помещение
        <select class="input-dark">
          <option v-for="zone in zones" :key="zone">{{ zone }}</option>
        </select>
      </label>

      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Дата
        <input type="date" class="input-dark" />
      </label>

      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Время
        <select class="input-dark">
          <option v-for="slot in timeSlots" :key="slot">{{ slot }}</option>
        </select>
      </label>

      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Количество участников
        <input type="number" min="1" value="2" class="input-dark" />
      </label>

      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Контактный телефон
        <input type="tel" placeholder="+7 (___) ___-__-__" class="input-dark" />
      </label>

      <label class="md:col-span-2 flex flex-col gap-2 text-sm font-medium text-gray-200">
        Комментарий
        <textarea rows="3" placeholder="Пожелания к мероприятию" class="input-dark"></textarea>
      </label>

      <div class="md:col-span-2 flex flex-wrap items-center gap-3">
        <button type="button" class="solid-button">Создать бронь</button>
        <NuxtLink to="/bookings" class="nav-link">Отмена</NuxtLink>
        <span class="text-sm text-gray-400">Стоимость рассчитана условно: 1 500 ₽</span>
      </div>
    </form>
  </section>
</template>

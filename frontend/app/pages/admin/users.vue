<script setup lang="ts">
const users = [
  { name: 'Анна Румянцева', email: 'anna@example.com', role: 'admin', status: 'Активен' },
  { name: 'Илья Орлов', email: 'ilya@example.com', role: 'seller', status: 'Активен' },
  { name: 'Мария Смирнова', email: 'maria@example.com', role: 'client', status: 'Заблокирован' }
]
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Админ · Пользователи</p>
        <h1 class="text-2xl font-bold text-white">Роли и статусы</h1>
        <p class="text-gray-400">Демо-управление ролями и активностью пользователей.</p>
      </div>
      <button class="ghost-button">Добавить пользователя</button>
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/10 bg-white/5 shadow-[0_20px_60px_rgba(0,0,0,0.45)]">
      <table class="min-w-full divide-y divide-white/5 text-sm text-gray-100">
        <thead class="bg-white/5 text-xs uppercase text-gray-400 tracking-[0.14em]">
          <tr>
            <th class="px-4 py-3 text-left">Имя</th>
            <th class="px-4 py-3 text-left">Email</th>
            <th class="px-4 py-3 text-left">Роль</th>
            <th class="px-4 py-3 text-left">Статус</th>
            <th class="px-4 py-3 text-left">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="user in users" :key="user.email" class="hover:bg-white/5 transition">
            <td class="px-4 py-3 font-medium text-white">{{ user.name }}</td>
            <td class="px-4 py-3 text-gray-300">{{ user.email }}</td>
            <td class="px-4 py-3 text-gray-300 uppercase">
              {{ user.role }}
            </td>
            <td class="px-4 py-3 text-gray-300">{{ user.status }}</td>
            <td class="px-4 py-3 text-gray-300">
              <button class="nav-link">Редактировать</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

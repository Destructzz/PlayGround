<template>
  <header class="sticky top-0 z-50 w-full border-b border-cyan-950 bg-[#020c13]/80 backdrop-blur-md">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <!-- Логотип -->
      <div class="flex items-center">
        <NuxtLink
          to="/"
          class="flex items-center font-black text-[26px] tracking-tighter transition-transform hover:scale-105"
        >
          <span class="text-white">Play</span>
          <span class="text-cyan-300">Ground</span>
        </NuxtLink>
      </div>

      <!-- Навигация -->
      <nav class="hidden md:flex items-center gap-8 text-sm font-medium text-cyan-100/60">
        <NuxtLink
          to="/"
          class="hover:text-white hover:drop-shadow-[0_0_8px_rgba(255,255,255,0.6)] transition-all"
        >Главная</NuxtLink>
        <NuxtLink
          to="/gaming"
          class="hover:text-white hover:drop-shadow-[0_0_8px_rgba(255,255,255,0.6)] transition-all"
        >Gaming</NuxtLink>
        <NuxtLink
          to="/lounge"
          class="hover:text-white hover:drop-shadow-[0_0_8px_rgba(255,255,255,0.6)] transition-all"
        >Lounge</NuxtLink>
        <NuxtLink
          to="/event"
          class="hover:text-white hover:drop-shadow-[0_0_8px_rgba(255,255,255,0.6)] transition-all"
        >Event</NuxtLink>
      </nav>

      <!-- Кнопки действий -->
      <div class="flex items-center gap-4">
        <NuxtLink
          to="/login"
          class="text-sm font-semibold text-cyan-100/60 hover:text-white hover:drop-shadow-[0_0_5px_rgba(255,255,255,0.4)] transition-all hidden sm:block"
        >Войти</NuxtLink>
        <NuxtLink
          to="/register"
          class="text-sm font-semibold text-cyan-300 hover:text-white hover:drop-shadow-[0_0_5px_rgba(34,211,238,0.5)] transition-all hidden sm:block"
        >Регистрация</NuxtLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
</script>

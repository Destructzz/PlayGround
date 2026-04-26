<template>
  <footer class="border-t border-cyan-950 bg-[#020c13]/90 backdrop-blur-md">
    <div class="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8 md:flex-row md:items-center md:justify-between">
      <div>
        <div class="flex items-center font-black text-[22px] tracking-tighter select-none">
          <span class="text-white">Play</span>
          <span class="text-cyan-300">Ground</span>
        </div>
        <p class="mt-2 max-w-md text-sm text-cyan-100/50">
          Lounge, gaming и event-маршруты в едином интерфейсе PlayGround.
        </p>
      </div>

      <div class="flex flex-col gap-2 text-sm text-cyan-100/55 md:items-end">
        <div class="flex flex-wrap gap-4">
          <NuxtLink
            to="/gaming"
            class="transition-colors hover:text-white"
          >Gaming</NuxtLink>
          <NuxtLink
            to="/lounge"
            class="transition-colors hover:text-white"
          >Lounge</NuxtLink>
          <NuxtLink
            to="/event"
            class="transition-colors hover:text-white"
          >Event</NuxtLink>
          <NuxtLink
            to="/login"
            class="transition-colors hover:text-white"
          >Login</NuxtLink>
          <NuxtLink
            to="/register"
            class="transition-colors hover:text-white"
          >Register</NuxtLink>
        </div>
        <p>© 2026 PlayGround</p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
</script>

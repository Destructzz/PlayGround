<template>
  <div class="min-h-screen bg-[#020c13] text-white flex flex-col font-sans selection:bg-cyan-400 selection:text-[#020c13]">
    <AppHeader v-if="route.path !== '/login'" />
    <main class="flex-1 flex flex-col">
      <NuxtPage />
    </main>
  </div>
</template>

<script setup>
const route = useRoute()
</script>

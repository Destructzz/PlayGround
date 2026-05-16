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
        <NuxtLink
          v-if="authStore.user?.role === 'admin'"
          to="/admin"
          class="hover:text-white hover:drop-shadow-[0_0_8px_rgba(103,232,249,0.8)] transition-all"
        >Admin</NuxtLink>
      </nav>

      <!-- Кнопки действий -->
      <div class="flex items-center gap-4">
        <template v-if="authStore.isLoading">
          <span class="text-sm text-cyan-100/50">Загрузка...</span>
        </template>
        <template v-else-if="authStore.isAuthenticated">
          <div class="flex items-center gap-3 border border-cyan-400/20 bg-cyan-950/40 px-3 py-2.5">
            <NuxtLink
              to="/profile"
              class="flex items-center gap-3 group transition-opacity hover:opacity-80"
            >
              <img
                v-if="authStore.user?.avatar_url"
                :src="authStore.user.avatar_url"
                :alt="authStore.user.name"
                class="h-9 w-9 rounded-full border border-cyan-300/30 object-cover"
              >
              <div
                v-else
                class="flex h-9 w-9 items-center justify-center rounded-full border border-cyan-300/30 bg-cyan-400/10 text-sm font-bold text-cyan-100"
              >
                {{ accountInitial }}
              </div>

              <div class="hidden min-w-0 sm:block">
                <p class="max-w-48 truncate text-sm font-medium text-cyan-100 group-hover:text-white transition-colors">
                  {{ displayName }}
                </p>
              </div>
            </NuxtLink>
            <button
              class="text-sm font-semibold text-red-400 hover:text-red-300 transition-all"
              @click="authStore.logout"
            >
              Выйти
            </button>
          </div>
        </template>
        <template v-else>
          <NuxtLink
            to="/login"
            class="text-sm font-semibold text-cyan-100/60 hover:text-white hover:drop-shadow-[0_0_5px_rgba(255,255,255,0.4)] transition-all hidden sm:block"
          >Войти</NuxtLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
const displayName = computed(() => authStore.user?.name || 'Пользователь')
const accountInitial = computed(() => displayName.value.charAt(0).toUpperCase())
</script>

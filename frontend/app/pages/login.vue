<template>
  <div class="flex-1 flex flex-col items-center justify-center p-8 relative min-h-full">
    <div class="mb-12 text-center flex flex-col items-center gap-1">
      <span class="text-cyan-100/60 text-lg font-medium tracking-widest uppercase mb-2">Sign in to</span>
      <div class="flex items-center font-black text-[2.75rem] tracking-tighter select-none">
        <span class="text-white drop-shadow-[0_0_12px_rgba(255,255,255,0.9)]">Play</span>
        <span class="text-cyan-300 drop-shadow-[0_0_18px_rgba(103,232,249,0.9)]">Ground</span>
      </div>
    </div>

    <div class="relative w-full max-w-md mx-auto group">
      <div class="absolute -inset-1 bg-cyan-500/20 rounded-[1.25rem] blur-xl opacity-50 group-hover:opacity-100 transition duration-700" />

      <div class="absolute inset-0 rounded-2xl overflow-hidden bg-cyan-950/30">
        <div class="absolute top-1/2 left-1/2 w-[200%] h-[200%] -translate-x-1/2 -translate-y-1/2 bg-[conic-gradient(from_0deg_at_50%_50%,transparent_60%,rgba(34,211,238,0.2)_85%,#22d3ee_98%,#ffffff_100%)] animate-[spin_6s_linear_infinite]" />
      </div>

      <div class="absolute inset-0 rounded-2xl overflow-hidden mix-blend-screen">
        <div class="absolute top-1/2 left-1/2 w-[200%] h-[200%] -translate-x-1/2 -translate-y-1/2 bg-[conic-gradient(from_180deg_at_50%_50%,transparent_60%,rgba(249,115,22,0.2)_85%,#f97316_98%,#ffffff_100%)] animate-[spin_6s_linear_infinite]" />
      </div>

      <div class="relative z-10 bg-[#020c13] m-[2px] rounded-[14px] p-8 sm:p-10 shadow-[0_0_40px_rgba(2,12,19,0.8)] flex flex-col gap-7">
        <template v-if="authStore.isAuthenticated">
          <div class="rounded-2xl border border-cyan-400/20 bg-cyan-950/40 p-5">
            <div class="flex items-center gap-4">
              <img
                v-if="authStore.user?.avatar_url"
                :src="authStore.user.avatar_url"
                :alt="authStore.user.name"
                class="h-12 w-12 rounded-full border border-cyan-300/30 object-cover"
              >
              <div
                v-else
                class="flex h-12 w-12 items-center justify-center rounded-full border border-cyan-300/30 bg-cyan-400/10 text-lg font-bold text-cyan-100"
              >
                {{ accountInitial }}
              </div>

              <div class="min-w-0">
                <p class="truncate text-base font-medium text-white">{{ displayName }}</p>
              </div>
            </div>

            <div class="mt-4 flex flex-wrap items-center gap-3">
              <NuxtLink
                to="/"
                class="text-sm font-semibold text-cyan-300 transition-colors hover:text-white"
              >
                На главную
              </NuxtLink>
              <button
                type="button"
                class="text-sm font-semibold text-red-300 transition-colors hover:text-red-200"
                @click="authStore.logout"
              >
                Выйти
              </button>
            </div>
          </div>
        </template>

        <template v-else>
          <a
            :href="googleAuthUrl"
            class="w-full bg-white hover:bg-gray-100 text-[#020c13] flex items-center justify-center gap-3 text-lg font-bold px-6 py-4 rounded-xl transition-all shadow-[0_0_15px_rgba(255,255,255,0.2)] hover:shadow-[0_0_20px_rgba(255,255,255,0.4)] active:scale-[0.98]"
          >
            <UIcon name="i-simple-icons-google" class="w-5 h-5" />
            Continue with Google
          </a>

          <div class="text-center text-sm text-cyan-100/55">
            Первый вход через Google автоматически создаст аккаунт.
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '~/stores/auth'

useHead({
  title: 'Sign In - PlayGround'
})

const authStore = useAuthStore()
const route = useRoute()
const config = useRuntimeConfig()
const displayName = computed(() => authStore.user?.name || 'Пользователь')

const returnTo = computed(() => {
  const value = route.query.return_to

  if (typeof value === 'string' && value.startsWith('/')) {
    return value
  }

  return '/'
})

const googleAuthUrl = computed(() => `${config.public.backendUrl.replace(/\/$/, '')}/api/v1/auth/google?${new URLSearchParams({ return_to: returnTo.value }).toString()}`)
const accountInitial = computed(() => displayName.value.charAt(0).toUpperCase())
</script>

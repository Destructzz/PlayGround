<template>
  <div class="admin-layout min-h-screen bg-[#020c13] pb-24 pt-24 text-white">
    <div class="mx-auto flex max-w-[1400px] flex-col gap-8 px-4 sm:px-6 lg:px-8">
      
      <!-- Admin Header (Title + Minimalist Profile) -->
      <header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b border-white/8 pb-6">
        <div>
          <span class="text-[10px] font-black uppercase tracking-[0.4em] text-cyan-400/60">System Operations</span>
          <h1 class="mt-1 text-3xl font-black tracking-tight text-white uppercase">Control Center</h1>
        </div>

        <!-- Minimalist Profile Button (Identical to Lounge) -->
        <button
          type="button"
          class="flex items-center gap-4 rounded-[0.85rem] border border-cyan-400/15 bg-[#07141d] px-5 py-3 text-left transition hover:border-cyan-300/30"
          @click="drawerOpen = true"
        >
          <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-cyan-400/15 text-sm font-black text-cyan-200">
            {{ authStore.user?.name?.charAt(0).toUpperCase() ?? 'A' }}
          </div>
          <div class="min-w-0">
            <p class="text-sm font-bold text-white leading-tight">{{ authStore.user?.name || 'Admin user' }}</p>
            <p class="truncate text-xs text-cyan-100/50 leading-tight mt-0.5">{{ authStore.user?.email || 'admin@playground.local' }}</p>
          </div>
        </button>
      </header>

      <!-- Horizontal Tabs (Equally ranking all administrative actions) -->
      <div class="flex gap-1 overflow-x-auto rounded-[0.85rem] border border-cyan-400/15 bg-[#07141d] p-1.5 shadow-lg scrollbar-thin">
        <button
          v-for="tab in adminTabs"
          :key="tab.id"
          type="button"
          class="flex items-center gap-2.5 rounded-[0.7rem] px-5 py-3 text-xs font-black uppercase tracking-wider transition-all duration-200 whitespace-nowrap"
          :class="activeTab === tab.id
            ? 'bg-cyan-300 text-[#020c13] shadow-[0_4px_20px_rgba(34,211,238,0.25)]'
            : 'text-cyan-100/50 hover:bg-white/[0.03] hover:text-white'"
          @click="selectTab(tab.id)"
        >
          <span>{{ tab.label }}</span>
        </button>
      </div>

      <!-- Main Nested Page Render -->
      <main class="relative">
        <NuxtPage />
      </main>

      <!-- Teleported User Profile Drawer (Exactly like Lounge) -->
      <Teleport to="body">
        <Transition name="drawer">
          <div
            v-if="drawerOpen"
            class="fixed inset-0 z-50 flex justify-end"
            @click.self="drawerOpen = false"
          >
            <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="drawerOpen = false" />
            <div class="relative z-10 flex h-full w-full max-w-sm flex-col bg-[#050f17] shadow-2xl">
              <div class="flex items-center justify-between border-b border-white/8 px-6 py-5">
                <p class="text-xs font-black uppercase tracking-[0.3em] text-cyan-100/50">Профиль администратора</p>
                <button
                  type="button"
                  class="rounded-full border border-white/10 p-2 text-zinc-400 transition hover:border-white/25 hover:text-white"
                  @click="drawerOpen = false"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="flex-1 overflow-y-auto px-6 py-8">
                <div class="flex flex-col items-center text-center">
                  <div class="flex h-20 w-20 items-center justify-center rounded-full border-2 border-cyan-300/30 bg-cyan-400/10 text-3xl font-black text-cyan-200">
                    {{ authStore.user?.name?.charAt(0).toUpperCase() ?? '?' }}
                  </div>
                  <h2 class="mt-5 text-xl font-black text-white">{{ authStore.user?.name }}</h2>
                  <p class="mt-1 text-sm text-cyan-100/50">{{ authStore.user?.email }}</p>
                  <span class="mt-3 rounded-full border border-cyan-300/20 bg-cyan-400/10 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-100">
                    {{ authStore.user?.role }}
                  </span>
                </div>
                <div class="mt-8 space-y-3 text-sm">
                  <div class="rounded-[0.8rem] border border-white/8 bg-white/4 px-4 py-3">
                    <p class="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Роль в системе</p>
                    <p class="mt-1 font-medium text-white">{{ authStore.user?.role ?? '—' }}</p>
                  </div>
                  <div class="rounded-[0.8rem] border border-white/8 bg-white/4 px-4 py-3">
                    <p class="text-[10px] font-bold uppercase tracking-[0.2em] text-zinc-500">Email адрес</p>
                    <p class="mt-1 font-medium text-white">{{ authStore.user?.email ?? '—' }}</p>
                  </div>
                </div>
                
                <button
                  type="button"
                  class="mt-8 w-full rounded-[0.85rem] border border-red-500/20 bg-red-500/10 py-3.5 text-xs font-black uppercase tracking-widest text-red-400 transition hover:bg-red-500/20 active:scale-95"
                  @click="logoutSession"
                >
                  Выйти из сессии
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </Teleport>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useRoute, useRouter } from 'vue-router'

definePageMeta({
  middleware: 'admin'
})

useHead({
  title: 'Control Center - PlayGround'
})

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

const drawerOpen = ref(false)

const adminTabs = [
  { id: 'zones', label: 'Зоны', icon: '🗺️' },
  { id: 'zone-tags', label: 'Теги зон', icon: '🏷️' },
  { id: 'bookings', label: 'Бронирования', icon: '📅' },
  { id: 'shifts', label: 'Смены', icon: '⏳' },
  { id: 'lounge', label: 'Lounge зоны', icon: '🛋️' },
  { id: 'event', label: 'События', icon: '🎪' },
  { id: 'settings', label: 'Настройки', icon: '⚙️' }
] as const

type TabId = typeof adminTabs[number]['id']

const activeTab = computed<TabId>(() => {
  if (route.path.endsWith('/lounge')) return 'lounge'
  if (route.path.endsWith('/event')) return 'event'
  if (route.path.endsWith('/settings')) return 'settings'
  return (route.query.tab as TabId) || 'zones'
})

function selectTab(tabId: TabId) {
  if (tabId === 'lounge') {
    router.push('/admin/lounge')
  } else if (tabId === 'event') {
    router.push('/admin/event')
  } else if (tabId === 'settings') {
    router.push('/admin/settings')
  } else {
    router.push({ path: '/admin', query: { tab: tabId } })
  }
}

function logoutSession() {
  drawerOpen.value = false
  authStore.logout()
}
</script>

<style scoped>
.scrollbar-thin::-webkit-scrollbar {
  height: 4px;
}
.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}
.scrollbar-thin::-webkit-scrollbar-thumb {
  background: rgba(34, 211, 238, 0.2);
  border-radius: 2px;
}
.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background: rgba(34, 211, 238, 0.4);
}

.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.25s ease;
}
.drawer-enter-active .relative,
.drawer-leave-active .relative {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}
.drawer-enter-from .relative {
  transform: translateX(100%);
}
.drawer-leave-to .relative {
  transform: translateX(100%);
}
</style>

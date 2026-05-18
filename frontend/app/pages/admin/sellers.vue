<template>
  <div class="space-y-6 pt-6">
    <!-- Local sub-tabs -->
    <div class="flex gap-1 rounded-[0.75rem] border border-white/8 bg-black/20 p-1">
      <button
        v-for="tab in localTabs"
        :key="tab.id"
        type="button"
        class="flex-1 rounded-[0.6rem] px-4 py-2.5 text-xs font-black uppercase tracking-wider transition-all duration-200"
        :class="activeLocal === tab.id
          ? 'bg-cyan-300 text-[#020c13] shadow-[0_2px_12px_rgba(34,211,238,0.3)]'
          : 'text-cyan-100/50 hover:text-white'"
        @click="activeLocal = tab.id"
      >
        {{ tab.label }}
        <span v-if="tab.id === 'sellers'" class="ml-1.5 rounded-full bg-white/15 px-1.5 py-0.5 text-[9px]">
          {{ sellers.length }}
        </span>
      </button>
    </div>

    <!-- ── Sellers Tab ── -->
    <template v-if="activeLocal === 'sellers'">
      <div v-if="loadingSellers" class="py-12 text-center text-sm text-zinc-500">
        Загружаем продавцов...
      </div>
      <div v-else-if="sellers.length === 0" class="rounded-[0.85rem] border border-white/8 bg-white/[0.02] px-6 py-12 text-center text-sm text-zinc-500">
        Нет пользователей с ролью seller.
      </div>
      <div v-else class="overflow-hidden rounded-[0.85rem] border border-white/8 bg-[#07141d]">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-white/8">
              <th class="px-5 py-3.5 text-left text-[10px] font-black uppercase tracking-widest text-zinc-500">Имя</th>
              <th class="px-5 py-3.5 text-left text-[10px] font-black uppercase tracking-widest text-zinc-500">Email</th>
              <th class="px-5 py-3.5 text-left text-[10px] font-black uppercase tracking-widest text-zinc-500">Телефон</th>
              <th class="px-5 py-3.5 text-right text-[10px] font-black uppercase tracking-widest text-zinc-500">Действие</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="seller in sellers"
              :key="seller.id"
              class="border-b border-white/[0.04] last:border-0 hover:bg-white/[0.02]"
            >
              <td class="px-5 py-3.5 font-medium text-white">{{ seller.full_name }}</td>
              <td class="px-5 py-3.5 text-zinc-400">{{ seller.email }}</td>
              <td class="px-5 py-3.5 text-zinc-500">{{ seller.phone || '—' }}</td>
              <td class="px-5 py-3.5 text-right">
                <button
                  class="rounded-lg border border-red-500/25 bg-red-500/10 px-3 py-1.5 text-[10px] font-black uppercase tracking-wider text-red-400 transition hover:bg-red-500/20"
                  :disabled="processing === seller.id"
                  @click="revokeRole(seller)"
                >
                  {{ processing === seller.id ? '...' : 'Снять роль' }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- ── Users Tab ── -->
    <template v-if="activeLocal === 'users'">
      <div class="relative">
        <input
          v-model="searchQuery"
          type="search"
          placeholder="Поиск по email..."
          class="w-full rounded-[0.75rem] border border-white/10 bg-[#07141d] px-5 py-3.5 pr-12 text-sm text-white placeholder-zinc-600 outline-none transition focus:border-cyan-400/40"
          @input="debouncedSearch"
        />
        <div class="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2 text-zinc-600">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>

      <div v-if="loadingUsers" class="py-12 text-center text-sm text-zinc-500">
        Ищем пользователей...
      </div>
      <div v-else-if="users.length === 0 && searchQuery" class="rounded-[0.85rem] border border-white/8 bg-white/[0.02] px-6 py-12 text-center text-sm text-zinc-500">
        Ничего не найдено по запросу «{{ searchQuery }}».
      </div>
      <div v-else-if="users.length > 0" class="overflow-hidden rounded-[0.85rem] border border-white/8 bg-[#07141d]">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-white/8">
              <th class="px-5 py-3.5 text-left text-[10px] font-black uppercase tracking-widest text-zinc-500">Имя</th>
              <th class="px-5 py-3.5 text-left text-[10px] font-black uppercase tracking-widest text-zinc-500">Email</th>
              <th class="px-5 py-3.5 text-left text-[10px] font-black uppercase tracking-widest text-zinc-500">Роль</th>
              <th class="px-5 py-3.5 text-right text-[10px] font-black uppercase tracking-widest text-zinc-500">Действие</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="user in users"
              :key="user.id"
              class="border-b border-white/[0.04] last:border-0 hover:bg-white/[0.02]"
            >
              <td class="px-5 py-3.5 font-medium text-white">{{ user.full_name }}</td>
              <td class="px-5 py-3.5 text-zinc-400">{{ user.email }}</td>
              <td class="px-5 py-3.5">
                <span
                  class="rounded-full px-2.5 py-0.5 text-[9px] font-black uppercase tracking-wider"
                  :class="{
                    'bg-cyan-400/15 text-cyan-300': user.role === 'seller',
                    'bg-red-400/15 text-red-300': user.role === 'admin',
                    'bg-zinc-700/50 text-zinc-400': user.role === 'client',
                  }"
                >{{ user.role }}</span>
              </td>
              <td class="px-5 py-3.5 text-right">
                <button
                  v-if="user.role === 'client'"
                  class="rounded-lg border border-cyan-400/25 bg-cyan-400/10 px-3 py-1.5 text-[10px] font-black uppercase tracking-wider text-cyan-300 transition hover:bg-cyan-400/20"
                  :disabled="processing === user.id"
                  @click="assignSeller(user)"
                >
                  {{ processing === user.id ? '...' : 'Назначить продавцом' }}
                </button>
                <span v-else class="text-[10px] text-zinc-600">—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- Toast notification -->
    <Transition name="fade">
      <div
        v-if="toast"
        class="fixed bottom-6 right-6 z-50 rounded-[0.85rem] border px-5 py-3.5 text-sm font-bold shadow-2xl"
        :class="toast.type === 'success'
          ? 'border-emerald-400/30 bg-emerald-500/15 text-emerald-300'
          : 'border-red-400/30 bg-red-500/15 text-red-300'"
      >
        {{ toast.message }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { searchAdminUsers, listSellers, setUserRole, type UserDoc } from '~/api/users'

const localTabs = [
  { id: 'sellers', label: 'Продавцы' },
  { id: 'users', label: 'Все пользователи' }
] as const

type LocalTab = typeof localTabs[number]['id']

const activeLocal = ref<LocalTab>('sellers')

const sellers = ref<UserDoc[]>([])
const users = ref<UserDoc[]>([])
const searchQuery = ref('')
const loadingSellers = ref(true)
const loadingUsers = ref(false)
const processing = ref<string | null>(null)
const toast = ref<{ type: 'success' | 'error', message: string } | null>(null)

let searchTimer: ReturnType<typeof setTimeout> | null = null

function showToast(type: 'success' | 'error', message: string) {
  toast.value = { type, message }
  setTimeout(() => { toast.value = null }, 3000)
}

async function fetchSellers() {
  loadingSellers.value = true
  try {
    const resp = await listSellers()
    sellers.value = resp.sellers ?? []
  } catch {
    sellers.value = []
  } finally {
    loadingSellers.value = false
  }
}

async function fetchUsers(q = '') {
  loadingUsers.value = true
  try {
    const resp = await searchAdminUsers(q || undefined)
    users.value = resp.users ?? []
  } catch {
    users.value = []
  } finally {
    loadingUsers.value = false
  }
}

function debouncedSearch() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    fetchUsers(searchQuery.value)
  }, 300)
}

async function assignSeller(user: UserDoc) {
  processing.value = user.id
  try {
    const resp = await setUserRole(user.id, 'seller')
    const idx = users.value.findIndex(u => u.id === user.id)
    if (idx !== -1 && resp.user) {
      users.value[idx] = resp.user
    }
    await fetchSellers()
    showToast('success', `${user.full_name} назначен продавцом`)
  } catch {
    showToast('error', 'Не удалось назначить роль')
  } finally {
    processing.value = null
  }
}

async function revokeRole(seller: UserDoc) {
  processing.value = seller.id
  try {
    const resp = await setUserRole(seller.id, 'client')
    sellers.value = sellers.value.filter(s => s.id !== seller.id)
    const idx = users.value.findIndex(u => u.id === seller.id)
    if (idx !== -1 && resp.user) {
      users.value[idx] = resp.user
    }
    showToast('success', `Роль seller снята с ${seller.full_name}`)
  } catch {
    showToast('error', 'Не удалось снять роль')
  } finally {
    processing.value = null
  }
}

onMounted(async () => {
  await fetchSellers()
  await fetchUsers()
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>

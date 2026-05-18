<template>
  <div class="min-h-screen bg-[#020c13] pt-24 pb-20 text-white">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="grid gap-8 lg:grid-cols-[380px_1fr]">
        <!-- Левая колонка: Профиль (Sidebar) -->
        <aside class="space-y-6">
          <div class="overflow-hidden rounded-[1rem] border border-cyan-400/20 bg-[#050f17] shadow-[0_20px_48px_rgba(0,0,0,0.4)]">
            <!-- Header Section -->
            <div class="border-b border-white/8 bg-[#08131b] px-6 py-10 flex flex-col items-center text-center">
              <div class="relative">
                <img
                  v-if="user?.avatar_url && !avatarError"
                  :src="user.avatar_url"
                  :alt="user.name"
                  class="h-32 w-32 rounded-full border-4 border-cyan-300/30 object-cover shadow-[0_0_30px_rgba(34,211,238,0.2)]"
                  @error="avatarError = true"
                >
                <div
                  v-else
                  class="flex h-32 w-32 items-center justify-center rounded-full border-4 border-cyan-300/30 bg-cyan-400/10 text-4xl font-black text-cyan-100 shadow-[0_0_30px_rgba(34,211,238,0.2)]"
                >
                  {{ accountInitial }}
                </div>
                <div class="absolute bottom-1 right-2 h-6 w-6 rounded-full border-4 border-[#08131b] bg-cyan-400"></div>
              </div>
              
              <h2 class="mt-6 text-2xl font-black tracking-tight text-white">{{ user?.name }}</h2>
              <p class="text-sm font-medium text-cyan-100/60">{{ user?.email }}</p>
              
              <div class="mt-4 inline-flex items-center rounded-full border border-cyan-300/20 bg-cyan-400/10 px-4 py-1.5 text-[11px] font-bold uppercase tracking-widest text-cyan-100">
                {{ user?.role === 'admin' ? 'Administrator' : 'Verified Member' }}
              </div>
            </div>

            <!-- Content Section -->
            <div class="p-6 space-y-8">
              <div>
                <h3 class="text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-100/40 mb-4">Contact Information</h3>
                
                <div class="space-y-5">
                  <div>
                    <label class="mb-1.5 block text-[10px] font-bold uppercase tracking-widest text-zinc-500">Email Address</label>
                    <div class="flex items-center gap-3 rounded-[0.7rem] border border-white/5 bg-white/5 px-4 py-3 text-sm text-zinc-400">
                      <span>📧</span>
                      {{ user?.email }}
                    </div>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-[10px] font-bold uppercase tracking-widest text-cyan-100/40">Full Name</label>
                    <div class="flex gap-2 mb-3">
                      <input
                        v-model="fullName"
                        type="text"
                        placeholder="Ivan Ivanov"
                        class="flex-1 rounded-[0.7rem] border border-cyan-400/20 bg-[#081824] px-4 py-3 text-sm text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none focus:ring-1 focus:ring-cyan-300 transition-all"
                      >
                    </div>
                    <label class="mb-1.5 block text-[10px] font-bold uppercase tracking-widest text-cyan-100/40">Phone Number</label>
                    <div class="flex gap-2">
                      <input
                        v-model="phone"
                        type="tel"
                        placeholder="+7 (___) ___-__-__"
                        class="flex-1 rounded-[0.7rem] border border-cyan-400/20 bg-[#081824] px-4 py-3 text-sm text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none focus:ring-1 focus:ring-cyan-300 transition-all"
                      >
                    </div>
                    <button 
                      class="mt-4 w-full rounded-[0.7rem] bg-cyan-300 py-3 text-xs font-black uppercase tracking-widest text-black hover:bg-cyan-200 transition-all shadow-lg active:scale-95"
                      @click="saveProfile"
                    >
                      Update Profile
                    </button>
                    <p v-if="saveMessage" class="mt-3 text-center text-xs font-bold flex items-center justify-center gap-2" :class="saveMessage.includes('Ошибка') ? 'text-red-400' : 'text-emerald-400'">
                      <span>{{ saveMessage.includes('Ошибка') ? '❌' : '✅' }}</span> {{ saveMessage }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- Mini Stats -->
              <div class="pt-6 border-t border-white/10">
                <div class="grid grid-cols-2 gap-4">
                  <div class="rounded-[0.8rem] bg-white/5 p-4 text-center">
                    <p class="text-[9px] font-bold uppercase tracking-widest text-zinc-500">Bookings</p>
                    <p class="mt-1 text-xl font-black">{{ allBookings.length }}</p>
                  </div>
                  <div class="rounded-[0.8rem] bg-white/5 p-4 text-center">
                    <p class="text-[9px] font-bold uppercase tracking-widest text-zinc-500">Member Since</p>
                    <p class="mt-1 text-xs font-black text-cyan-100/80">May 2026</p>
                  </div>
                </div>
              </div>
              
              <div class="pt-4">
                 <button
                   @click="authStore.logout"
                   class="w-full rounded-[0.8rem] border border-red-500/20 bg-red-500/10 py-3.5 text-[11px] font-black uppercase tracking-widest text-red-400 transition hover:bg-red-500/20 hover:text-red-300"
                 >
                   Sign Out
                 </button>
              </div>
            </div>
          </div>
        </aside>

        <!-- Правая колонка: Заказы (Main Content) -->
        <div class="space-y-8">
          <header class="flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold uppercase tracking-[0.3em] text-cyan-300/60">Booking History</p>
              <h1 class="mt-1 text-3xl font-black tracking-tight text-white">My Reservations</h1>
              <p class="mt-2 text-sm text-zinc-400">Track your gaming sessions and lounge bookings below.</p>
            </div>
          </header>

          <!-- Tab Switcher -->
          <div class="mt-8 flex gap-8 border-b border-white/5">
            <button 
              @click="activeBookingTab = 'current'"
              class="relative pb-4 text-[11px] font-black uppercase tracking-[0.2em] transition-all"
              :class="activeBookingTab === 'current' ? 'text-cyan-400' : 'text-zinc-500 hover:text-white'"
            >
              Current
              <span v-if="activeBookingTab === 'current'" class="absolute bottom-0 left-0 h-0.5 w-full bg-cyan-400 shadow-[0_0_10px_rgba(34,211,238,0.5)]"></span>
              <span class="ml-2 rounded-md bg-white/5 px-2 py-0.5 text-[9px] text-zinc-400">{{ currentBookings.length }}</span>
            </button>
            <button 
              @click="activeBookingTab = 'archive'"
              class="relative pb-4 text-[11px] font-black uppercase tracking-[0.2em] transition-all"
              :class="activeBookingTab === 'archive' ? 'text-cyan-400' : 'text-zinc-500 hover:text-white'"
            >
              Archive
              <span v-if="activeBookingTab === 'archive'" class="absolute bottom-0 left-0 h-0.5 w-full bg-cyan-400 shadow-[0_0_10px_rgba(34,211,238,0.5)]"></span>
              <span class="ml-2 rounded-md bg-white/5 px-2 py-0.5 text-[9px] text-zinc-400">{{ archiveBookings.length }}</span>
            </button>
          </div>

          <div v-if="isLoadingBookings" class="rounded-[1.5rem] border border-white/5 bg-[#050f17] p-20 text-center shadow-2xl">
            <div class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-cyan-300 border-t-transparent"></div>
            <p class="mt-4 text-sm font-bold text-zinc-500 tracking-widest uppercase">Fetching Records...</p>
          </div>
          
          <div v-else-if="!displayBookings.length" class="rounded-[1.5rem] border border-white/5 bg-[#050f17] p-20 text-center shadow-2xl">
            <div class="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-white/5">
              <span class="text-4xl opacity-50">🕹️</span>
            </div>
            <h3 class="text-2xl font-black text-white">No {{ activeBookingTab }} bookings found</h3>
            <p class="mt-3 text-sm text-zinc-400 max-w-sm mx-auto leading-relaxed">It looks like you haven't made any reservations yet. Start your journey in the gaming or lounge zones.</p>
            <div class="mt-8 flex flex-wrap justify-center gap-4">
              <NuxtLink to="/gaming" class="rounded-full bg-cyan-300 px-8 py-3.5 text-xs font-black uppercase tracking-widest text-black transition hover:bg-cyan-200 hover:shadow-[0_0_20px_rgba(34,211,238,0.4)]">
                Gaming Zone
              </NuxtLink>
              <NuxtLink to="/lounge" class="rounded-full border border-white/20 px-8 py-3.5 text-xs font-black uppercase tracking-widest text-white transition hover:bg-white/10">
                Lounge Zone
              </NuxtLink>
            </div>
          </div>

          <div v-else class="space-y-4">
            <div 
              v-for="booking in displayBookings" 
              :key="booking.id"
              class="overflow-hidden rounded-[1rem] border border-cyan-400/15 bg-[#050f17] transition-all hover:border-cyan-400/30"
            >
              <!-- Краткая информация (всегда видна) -->
              <div 
                class="flex cursor-pointer items-center justify-between p-5"
                @click="toggleBooking(booking.id)"
              >
                <div class="flex items-center gap-5">
                  <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-[0.8rem] bg-[#0a1824] shadow-[inset_0_1px_0_rgba(255,255,255,0.05)]">
                    <span class="text-lg" v-if="booking.zone_id === 1">🎮</span>
                    <span class="text-lg" v-else-if="booking.zone_id === 2">🛋️</span>
                    <span class="text-lg" v-else>🎟️</span>
                  </div>
                  
                  <div>
                    <div class="flex items-center gap-2">
                      <h3 class="text-lg font-black text-white">
                        {{ zoneNamesMap[booking.zone_id] || `Бронь #${booking.id}` }}
                      </h3>
                      <span 
                        class="rounded-full border px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider"
                        :class="statusClasses(booking.status)"
                      >
                        {{ statusLabel(booking.status) }}
                      </span>
                    </div>
                    <p class="mt-1 text-sm font-medium text-cyan-100/60">
                      {{ formatDateTime(booking.start_time) }}
                    </p>
                  </div>
                </div>
                
                <div class="flex items-center gap-4">
                  <div class="text-right hidden sm:block">
                    <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-500">Сумма</p>
                    <p class="text-lg font-black text-white">{{ booking.total_price }} ₽</p>
                  </div>
                  
                  <div class="flex h-8 w-8 items-center justify-center rounded-full border border-white/10 bg-white/5 transition-colors" :class="{'rotate-180 bg-white/10': expandedBooking === booking.id}">
                    <svg width="12" height="8" viewBox="0 0 12 8" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M1 1.5L6 6.5L11 1.5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                  </div>
                </div>
              </div>
              
              <!-- Развернутая информация -->
              <div 
                v-if="expandedBooking === booking.id"
                class="border-t border-white/5 bg-[#030a10] p-5"
              >
                <div class="grid gap-6 sm:grid-cols-2">
                  <div class="space-y-4">
                    <div>
                      <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-500">Детали бронирования</p>
                      <ul class="mt-2 space-y-2 text-sm text-zinc-300">
                        <li class="flex justify-between">
                          <span class="text-zinc-500">Начало:</span>
                          <span class="font-medium text-white">{{ formatDateTime(booking.start_time) }}</span>
                        </li>
                        <li class="flex justify-between">
                          <span class="text-zinc-500">Окончание:</span>
                          <span class="font-medium text-white">{{ formatDateTime(booking.end_time) }}</span>
                        </li>
                        <li class="flex justify-between">
                          <span class="text-zinc-500">Количество гостей:</span>
                          <span class="font-medium text-white">{{ booking.participants }} чел.</span>
                        </li>
                        <li class="flex justify-between sm:hidden">
                          <span class="text-zinc-500">Сумма:</span>
                          <span class="font-medium text-white">{{ booking.total_price }} ₽</span>
                        </li>
                      </ul>
                    </div>
                  </div>
                  
                  <div class="space-y-4">
                      <div class="flex items-center justify-between">
                        <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-500">Контактные данные</p>
                        <button 
                          v-if="activeBookingTab === 'current' && editBookingId !== booking.id"
                          @click.stop="startEditBooking(booking)"
                          class="text-[10px] font-bold uppercase text-cyan-400 hover:text-cyan-300"
                        >
                          Редактировать
                        </button>
                      </div>

                      <div v-if="editBookingId === booking.id" class="mt-3 space-y-3 rounded-lg border border-cyan-500/20 bg-black/20 p-3">
                        <div>
                          <label class="mb-1 text-[9px] font-bold uppercase text-zinc-500">Имя</label>
                          <input v-model="editBookingContacts.name" type="text" class="w-full rounded bg-white/5 px-3 py-1.5 text-xs text-white outline-none focus:border-cyan-500/50" />
                        </div>
                        <div>
                          <label class="mb-1 text-[9px] font-bold uppercase text-zinc-500">Email</label>
                          <input v-model="editBookingContacts.email" type="email" class="w-full rounded bg-white/5 px-3 py-1.5 text-xs text-white outline-none focus:border-cyan-500/50" />
                        </div>
                        <div>
                          <label class="mb-1 text-[9px] font-bold uppercase text-zinc-500">Телефон</label>
                          <input v-model="editBookingContacts.phone" type="tel" class="w-full rounded bg-white/5 px-3 py-1.5 text-xs text-white outline-none focus:border-cyan-500/50" />
                        </div>
                        <div class="flex gap-2 pt-1">
                          <button @click.stop="saveBookingContacts(booking.id)" class="flex-1 rounded bg-cyan-400 py-1.5 text-[10px] font-bold uppercase text-black hover:bg-cyan-300">
                            {{ isSavingBooking ? '...' : 'Сохранить' }}
                          </button>
                          <button @click.stop="cancelEditBooking" class="flex-1 rounded bg-white/10 py-1.5 text-[10px] font-bold uppercase text-white hover:bg-white/20">
                            Отмена
                          </button>
                        </div>
                      </div>

                      <ul v-else class="mt-2 space-y-2 text-sm text-zinc-300">
                        <li class="flex justify-between">
                          <span class="text-zinc-500">Имя:</span>
                          <span class="font-medium text-white">{{ booking.contact_name || user?.name || 'Не указано' }}</span>
                        </li>
                        <li class="flex justify-between">
                          <span class="text-zinc-500">Email:</span>
                          <span class="font-medium text-white">{{ booking.contact_email || user?.email || 'Не указан' }}</span>
                        </li>
                        <li class="flex justify-between">
                          <span class="text-zinc-500">Телефон:</span>
                          <span class="font-medium text-white">{{ booking.contact_phone || phone || 'Не указан' }}</span>
                        </li>
                      </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { updateUserProfile } from '~/api/auth'
import { getMyBookingsCategorized, patchBooking } from '~/api/booking'
import { getHomeCatalog } from '~/api/catalog'

definePageMeta({
  middleware: ['auth']
})

useHead({
  title: 'Мой профиль - PlayGround'
})

const authStore = useAuthStore()
const user = computed(() => authStore.user)
const accountInitial = computed(() => user.value?.name?.charAt(0).toUpperCase() || 'P')

const avatarError = ref(false)
watch(() => user.value, () => {
  avatarError.value = false
})

const fullName = ref('')
const phone = ref('')
const saveMessage = ref('')
const zoneNamesMap = ref<Record<number, string>>({})
const activeBookingTab = ref<'current' | 'archive'>('current')
const currentBookings = ref<any[]>([])
const archiveBookings = ref<any[]>([])
const allBookings = computed(() => [...currentBookings.value, ...archiveBookings.value])
const displayBookings = computed(() => activeBookingTab.value === 'current' ? currentBookings.value : archiveBookings.value)
const isLoadingBookings = ref(true)
const expandedBooking = ref<number | null>(null)

// Booking Edit State
const editBookingId = ref<number | null>(null)
const editBookingContacts = ref({ name: '', email: '', phone: '' })
const isSavingBooking = ref(false)

function startEditBooking(booking: any) {
  editBookingId.value = booking.id
  editBookingContacts.value = {
    name: booking.contact_name || user.value?.name || '',
    email: booking.contact_email || user.value?.email || '',
    phone: booking.contact_phone || phone.value || ''
  }
}

function cancelEditBooking() {
  editBookingId.value = null
}

async function saveBookingContacts(bookingId: number) {
  isSavingBooking.value = true
  try {
    const resp = await patchBooking(bookingId, {
      contact_name: editBookingContacts.value.name,
      contact_email: editBookingContacts.value.email,
      contact_phone: editBookingContacts.value.phone
    })
    
    // Update local state
    const index = currentBookings.value.findIndex(b => b.id === bookingId)
    if (index !== -1 && resp.booking) {
      currentBookings.value[index].contact_name = resp.booking.contact_name
      currentBookings.value[index].contact_email = resp.booking.contact_email
      currentBookings.value[index].contact_phone = resp.booking.contact_phone
    }
    
    editBookingId.value = null
  } catch (e) {
    console.error('Failed to update booking contacts:', e)
  } finally {
    isSavingBooking.value = false
  }
}


onMounted(async () => {
  if (user.value) {
    fullName.value = user.value.name
  }
  const savedPhone = localStorage.getItem('playground_phone')
  if (savedPhone) {
    phone.value = savedPhone
  }
  
  try {
    const [bookingResp, catalogResp] = await Promise.all([
      getMyBookingsCategorized().catch(() => null),
      getHomeCatalog().catch(() => null)
    ])

    if (bookingResp) {
      currentBookings.value = bookingResp.current || []
      archiveBookings.value = bookingResp.archive || []
    }
    
    if (catalogResp) {
      const allZones = [
        ...(catalogResp.gaming || []),
        ...(catalogResp.lounge || []),
        ...(catalogResp.event || [])
      ]
      for (const z of allZones) {
        zoneNamesMap.value[z.id] = z.name
      }
    }
  } catch (e) {
    console.error('Failed to load profile data:', e)
  } finally {
    isLoadingBookings.value = false
  }
})

async function saveProfile() {
  localStorage.setItem('playground_phone', phone.value)
  try {
    const resp = await updateUserProfile({
      full_name: fullName.value,
      phone: phone.value
    })
    if (resp.user) {
      await authStore.fetchSession()
    }
    saveMessage.value = 'Сохранено'
  } catch (e) {
    console.error(e)
    saveMessage.value = 'Ошибка сохранения'
  }
  setTimeout(() => {
    saveMessage.value = ''
  }, 2000)
}

function toggleBooking(id: number) {
  expandedBooking.value = expandedBooking.value === id ? null : id
}

function formatDateTime(dateString: string) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

function statusLabel(status: string) {
  switch (status) {
    case 'created': return 'Ожидает'
    case 'confirmed': return 'Подтвержден'
    case 'completed': return 'Завершен'
    case 'canceled': return 'Отменен'
    default: return status
  }
}

function statusClasses(status: string) {
  switch (status) {
    case 'created': return 'border-cyan-300/20 bg-cyan-400/10 text-cyan-100'
    case 'confirmed': return 'border-emerald-300/20 bg-emerald-500/10 text-emerald-300'
    case 'completed': return 'border-white/10 bg-white/5 text-zinc-300'
    case 'canceled': return 'border-red-300/20 bg-red-500/10 text-red-300'
    default: return 'border-white/10 bg-white/5 text-zinc-300'
  }
}
</script>

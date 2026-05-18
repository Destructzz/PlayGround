<template>
  <div class="min-h-screen bg-[#020c13] pb-24 pt-24 text-white">
    <div class="mx-auto flex max-w-[900px] flex-col gap-8 px-4 sm:px-6 lg:px-8">
      
      <!-- Header -->
      <header class="border-b border-white/8 pb-6">
        <span class="text-[10px] font-black uppercase tracking-[0.4em] text-cyan-400/60">Seller Panel</span>
        <h1 class="mt-1 text-3xl font-black uppercase tracking-tight text-white">Панель продавца</h1>
        <p class="mt-2 text-sm text-zinc-500">Введи номер заказа или выбери из списка, чтобы проверить бронирование клиента</p>
      </header>

      <!-- Search form -->
      <div class="space-y-4">
        <label class="block text-[10px] font-black uppercase tracking-widest text-zinc-500">Номер брони (ID)</label>
        <div class="flex gap-3">
          <input
            v-model="bookingIdInput"
            type="number"
            placeholder="Например: 42"
            min="1"
            class="flex-1 rounded-[0.75rem] border border-white/10 bg-[#07141d] px-5 py-4 text-base text-white placeholder-zinc-600 outline-none transition focus:border-cyan-400/40 focus:shadow-[0_0_0_3px_rgba(34,211,238,0.08)]"
            @keydown.enter="searchBooking"
          />
          <button
            class="rounded-[0.75rem] bg-cyan-300 px-8 py-4 text-sm font-black uppercase tracking-widest text-[#020c13] transition hover:bg-cyan-200 disabled:opacity-50"
            :disabled="!bookingIdInput || loading"
            @click="searchBooking"
          >
            {{ loading ? 'Ищем...' : 'Найти' }}
          </button>
        </div>
        <p v-if="error" class="text-sm text-red-400">{{ error }}</p>
      </div>

      <!-- Booking card -->
      <Transition name="slide-up">
        <div
          v-if="booking"
          class="overflow-hidden rounded-[0.9rem] border border-cyan-400/20 bg-[#07141d] shadow-[0_8px_40px_rgba(0,0,0,0.4)]"
        >
          <!-- Card header -->
          <div class="flex items-center justify-between border-b border-white/8 px-6 py-5">
            <div>
              <p class="text-[10px] font-black uppercase tracking-widest text-cyan-400/60">Бронирование</p>
              <h2 class="mt-1 text-2xl font-black text-white">#{{ booking.id }}</h2>
            </div>
            <span
              class="rounded-full px-3 py-1.5 text-[10px] font-black uppercase tracking-wider"
              :class="{
                'bg-emerald-400/15 text-emerald-300': booking.status === 'confirmed',
                'bg-cyan-400/15 text-cyan-300': booking.status === 'created',
                'bg-red-400/15 text-red-300': booking.status === 'canceled',
                'bg-zinc-700/50 text-zinc-400': booking.status === 'completed',
              }"
            >{{ booking.status }}</span>
          </div>

          <div class="grid gap-6 p-6 sm:grid-cols-2">
            <!-- Booking info -->
            <div class="space-y-4">
              <p class="text-[10px] font-black uppercase tracking-widest text-zinc-500">Детали брони</p>
              <div class="space-y-2.5 text-sm">
                <div class="flex justify-between">
                  <span class="text-zinc-500">Зона</span>
                  <span class="font-bold text-white">{{ booking.zone_name || `Зона #${booking.zone_id}` }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-zinc-500">Услуга / Тариф</span>
                  <span class="font-bold text-white">{{ booking.service_name || `Сервис #${booking.service_id}` }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-zinc-500">Начало</span>
                  <span class="font-bold text-white">{{ formatDateTime(booking.start_time) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-zinc-500">Конец</span>
                  <span class="font-bold text-white">{{ formatDateTime(booking.end_time) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-zinc-500">Участников</span>
                  <span class="font-bold text-white">{{ booking.participants }}</span>
                </div>
                <div class="flex justify-between border-t border-white/8 pt-2.5">
                  <span class="text-zinc-500">Итого</span>
                  <span class="text-lg font-black text-cyan-300">{{ booking.total_price }} ₽</span>
                </div>
              </div>
            </div>

            <!-- Contact info -->
            <div class="space-y-4">
              <p class="text-[10px] font-black uppercase tracking-widest text-zinc-500">Контакты клиента</p>
              <div class="space-y-3">
                <div class="rounded-[0.7rem] border border-white/8 bg-white/[0.03] px-4 py-3">
                  <p class="text-[9px] font-black uppercase tracking-wider text-zinc-600">Имя</p>
                  <p class="mt-1 font-bold text-white">{{ booking.contact_name || '—' }}</p>
                </div>
                <div class="rounded-[0.7rem] border border-white/8 bg-white/[0.03] px-4 py-3">
                  <p class="text-[9px] font-black uppercase tracking-wider text-zinc-600">Email</p>
                  <p class="mt-1 font-medium text-white">{{ booking.contact_email || '—' }}</p>
                </div>
                <div class="rounded-[0.7rem] border border-white/8 bg-white/[0.03] px-4 py-3">
                  <p class="text-[9px] font-black uppercase tracking-wider text-zinc-600">Телефон</p>
                  <p class="mt-1 font-medium text-white">{{ booking.contact_phone || '—' }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Bookings List Section -->
      <section class="mt-8 space-y-4">
        <div class="flex items-center justify-between border-b border-white/8 pb-4">
          <div>
            <h2 class="text-xl font-black uppercase tracking-tight text-white">Все заявки</h2>
            <p class="text-xs text-zinc-500">Отсортированы по дате начала (сначала новые/ближайшие)</p>
          </div>
          <button 
            @click="loadAllBookings" 
            class="rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs font-bold text-white hover:bg-white/10 transition"
            :disabled="listLoading"
          >
            {{ listLoading ? 'Обновление...' : 'Обновить список' }}
          </button>
        </div>

        <div v-if="listLoading && bookingsList.length === 0" class="flex flex-col items-center justify-center py-12 text-zinc-500">
          <svg class="h-8 w-8 animate-spin text-cyan-400" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          <span class="mt-3 text-xs uppercase tracking-wider font-bold">Загрузка заявок...</span>
        </div>

        <div v-else-if="bookingsList.length === 0" class="text-center py-12 text-zinc-500 border border-dashed border-white/8 rounded-[0.9rem] bg-white/[0.01]">
          Нет зарегистрированных заявок.
        </div>

        <div v-else class="grid gap-4">
          <div 
            v-for="b in bookingsList" 
            :key="b.id"
            class="group relative overflow-hidden rounded-[0.8rem] border border-white/5 bg-[#07141d]/50 p-5 transition-all duration-300 hover:border-cyan-400/30 hover:bg-[#07141d]"
          >
            <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <div class="flex items-center gap-2.5">
                  <span class="text-sm font-black text-white">#{{ b.id }}</span>
                  <span class="text-xs font-bold text-cyan-300">{{ b.contact_name }}</span>
                  <span 
                    class="rounded-full px-2 py-0.5 text-[9px] font-bold uppercase tracking-wider"
                    :class="{
                      'bg-emerald-400/10 text-emerald-400 border border-emerald-400/20': b.status === 'confirmed',
                      'bg-cyan-400/10 text-cyan-400 border border-cyan-400/20': b.status === 'created',
                      'bg-red-400/10 text-red-400 border border-red-400/20': b.status === 'canceled',
                      'bg-zinc-700/30 text-zinc-400 border border-white/5': b.status === 'completed',
                    }"
                  >
                    {{ b.status }}
                  </span>
                </div>
                
                <div class="mt-2.5 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-xs text-zinc-400">
                  <span class="flex items-center gap-1">
                    <span class="text-zinc-600 font-medium">Зона:</span>
                    <strong class="text-zinc-300 font-semibold">{{ b.zone_name || `Зона #${b.zone_id}` }}</strong>
                  </span>
                  <span class="flex items-center gap-1">
                    <span class="text-zinc-600 font-medium">Тариф:</span>
                    <strong class="text-zinc-300 font-semibold">{{ b.service_name || `Сервис #${b.service_id}` }}</strong>
                  </span>
                  <span class="flex items-center gap-1">
                    <span class="text-zinc-600 font-medium">Участников:</span>
                    <strong class="text-zinc-300 font-semibold">{{ b.participants }}</strong>
                  </span>
                </div>

                <div class="mt-1.5 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-zinc-500">
                  <span>📅 {{ formatDateTime(b.start_time) }} — {{ formatDateTime(b.end_time) }}</span>
                  <span v-if="b.contact_phone">📞 {{ b.contact_phone }}</span>
                </div>
              </div>

              <div class="flex items-center justify-between border-t border-white/5 pt-3 sm:border-t-0 sm:pt-0 gap-4">
                <div class="text-left sm:text-right">
                  <p class="text-[9px] uppercase tracking-wider text-zinc-600 font-bold">Итого</p>
                  <p class="text-base font-black text-cyan-300">{{ b.total_price }} ₽</p>
                </div>
                <button 
                  @click="selectBookingFromList(b.id)"
                  class="rounded-lg bg-cyan-300/10 hover:bg-cyan-300/20 px-3.5 py-2 text-xs font-bold uppercase tracking-wider text-cyan-300 transition"
                >
                  Выбрать
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getBookingForSeller, listBookingsForSeller } from '~/api/users'

definePageMeta({
  middleware: 'seller'
})

useHead({ title: 'Seller Panel - PlayGround' })

const bookingIdInput = ref<number | null>(null)
const booking = ref<any>(null)
const loading = ref(false)
const error = ref('')

const bookingsList = ref<any[]>([])
const listLoading = ref(false)

async function searchBooking() {
  if (!bookingIdInput.value) return
  loading.value = true
  error.value = ''
  booking.value = null
  try {
    const resp = await getBookingForSeller(bookingIdInput.value)
    booking.value = resp.booking
  } catch (e: any) {
    error.value = e?.data?.error?.message ?? 'Бронирование не найдено. Проверь номер.'
  } finally {
    loading.value = false
  }
}

async function loadAllBookings() {
  listLoading.value = true
  try {
    const resp = await listBookingsForSeller()
    bookingsList.value = resp.bookings || []
  } catch (e) {
    console.error('Failed to load bookings list:', e)
  } finally {
    listLoading.value = false
  }
}

function selectBookingFromList(id: number) {
  bookingIdInput.value = id
  searchBooking()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadAllBookings()
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(16px);
}
</style>

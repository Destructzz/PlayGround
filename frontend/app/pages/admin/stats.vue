<template>
  <div class="space-y-8">
    <!-- Summary Cards -->
    <div v-if="stats" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg transition hover:border-cyan-400/20">
        <p class="text-[10px] font-black uppercase tracking-widest text-zinc-500">Всего бронирований</p>
        <h3 class="mt-2 text-3xl font-black text-white">{{ stats.summary.total_bookings }}</h3>
        <p class="mt-1 text-xs text-cyan-400/60">+12% с прошлого месяца</p>
      </div>
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg transition hover:border-cyan-400/20">
        <p class="text-[10px] font-black uppercase tracking-widest text-zinc-500">Общая выручка</p>
        <h3 class="mt-2 text-3xl font-black text-cyan-300">{{ formatCurrency(stats.summary.total_revenue) }}</h3>
        <p class="mt-1 text-xs text-emerald-400/60">+8% с прошлого месяца</p>
      </div>
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg transition hover:border-cyan-400/20">
        <p class="text-[10px] font-black uppercase tracking-widest text-zinc-500">Активные пользователи</p>
        <h3 class="mt-2 text-3xl font-black text-white">{{ stats.summary.active_users }}</h3>
        <p class="mt-1 text-xs text-zinc-500">Зарегистрировано в системе</p>
      </div>
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg transition hover:border-cyan-400/20">
        <p class="text-[10px] font-black uppercase tracking-widest text-zinc-500">Ожидают подтверждения</p>
        <h3 class="mt-2 text-3xl font-black text-amber-400">{{ stats.summary.pending_bookings }}</h3>
        <p class="mt-1 text-xs text-zinc-500">Новые бронирования</p>
      </div>
    </div>

    <div v-if="loading" class="flex h-64 items-center justify-center">
      <div class="h-12 w-12 animate-spin rounded-full border-4 border-cyan-400/20 border-t-cyan-400"></div>
    </div>

    <div v-else-if="stats" class="grid gap-8 lg:grid-cols-2">
      <!-- Revenue Chart -->
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg">
        <div class="mb-6 flex items-center justify-between">
          <h4 class="text-sm font-black uppercase tracking-widest text-white">Выручка (30 дней)</h4>
          <span class="text-[10px] font-bold text-zinc-500">RUB / День</span>
        </div>
        <div class="h-[300px]">
          <Line :data="revenueChartData" :options="chartOptions" />
        </div>
      </div>

      <!-- Bookings Chart -->
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg">
        <div class="mb-6 flex items-center justify-between">
          <h4 class="text-sm font-black uppercase tracking-widest text-white">Бронирования (30 дней)</h4>
          <span class="text-[10px] font-bold text-zinc-500">Кол-во / День</span>
        </div>
        <div class="h-[300px]">
          <Bar :data="bookingsChartData" :options="chartOptions" />
        </div>
      </div>

      <!-- By Zone Chart -->
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg">
        <div class="mb-6">
          <h4 class="text-sm font-black uppercase tracking-widest text-white">Популярность зон</h4>
        </div>
        <div class="flex h-[300px] items-center justify-center">
          <Doughnut :data="zoneChartData" :options="doughnutOptions" />
        </div>
      </div>

      <!-- Top Services / Recent activity placeholder -->
      <div class="rounded-[0.9rem] border border-white/8 bg-[#07141d] p-6 shadow-lg">
        <div class="mb-6">
          <h4 class="text-sm font-black uppercase tracking-widest text-white">Последние события</h4>
        </div>
        <div class="space-y-4">
          <div v-for="i in 5" :key="i" class="flex items-center gap-4 border-b border-white/4 pb-3 last:border-0">
            <div class="h-2 w-2 rounded-full bg-cyan-400 shadow-[0_0_8px_rgba(34,211,238,0.6)]"></div>
            <div class="flex-1">
              <p class="text-sm font-medium text-white">Новое бронирование #{{ 100 + i }}</p>
              <p class="text-[10px] text-zinc-500">2 часа назад • Клиент: Иван Иванов</p>
            </div>
            <span class="text-[10px] font-bold text-cyan-300">1 500 ₽</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getAdminStats } from '~/api/admin'
import type { AdminStatsResponse } from '~/api/types'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { Line, Bar, Doughnut } from 'vue-chartjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const stats = ref<AdminStatsResponse | null>(null)
const loading = ref(true)

async function fetchStats() {
  try {
    stats.value = await getAdminStats()
  } catch (e) {
    console.error('Failed to fetch stats', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchStats)

const formatCurrency = (val: string) => {
  const num = parseFloat(val) || 0
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0
  }).format(num)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: 'short'
  })
}

const revenueChartData = computed(() => {
  if (!stats.value) return { labels: [], datasets: [] }
  return {
    labels: stats.value.revenue.map(d => formatDate(d.date)),
    datasets: [
      {
        label: 'Выручка',
        backgroundColor: 'rgba(34, 211, 238, 0.1)',
        borderColor: '#22d3ee',
        borderWidth: 3,
        pointBackgroundColor: '#22d3ee',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: '#22d3ee',
        fill: true,
        tension: 0.4,
        data: stats.value.revenue.map(d => d.value)
      }
    ]
  }
})

const bookingsChartData = computed(() => {
  if (!stats.value) return { labels: [], datasets: [] }
  return {
    labels: stats.value.bookings.map(d => formatDate(d.date)),
    datasets: [
      {
        label: 'Бронирования',
        backgroundColor: '#22d3ee',
        borderRadius: 4,
        data: stats.value.bookings.map(d => d.value)
      }
    ]
  }
})

const zoneChartData = computed(() => {
  if (!stats.value) return { labels: [], datasets: [] }
  const colors = ['#22d3ee', '#818cf8', '#f472b6', '#fbbf24', '#34d399']
  return {
    labels: stats.value.by_zone.map(z => z.zone_type),
    datasets: [
      {
        backgroundColor: colors,
        borderColor: '#07141d',
        borderWidth: 4,
        data: stats.value.by_zone.map(z => z.count)
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#050f17',
      titleFont: { size: 12, weight: 'bold' },
      bodyFont: { size: 12 },
      padding: 12,
      borderColor: 'rgba(255,255,255,0.1)',
      borderWidth: 1,
      displayColors: false
    }
  },
  scales: {
    y: {
      grid: { color: 'rgba(255,255,255,0.05)' },
      ticks: { color: 'rgba(255,255,255,0.4)', font: { size: 10 } }
    },
    x: {
      grid: { display: false },
      ticks: { color: 'rgba(255,255,255,0.4)', font: { size: 10 } }
    }
  }
}

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right' as const,
      labels: {
        color: 'rgba(255,255,255,0.6)',
        font: { size: 11, weight: 'bold' },
        padding: 20,
        usePointStyle: true,
        pointStyle: 'circle'
      }
    }
  },
  cutout: '70%'
}
</script>

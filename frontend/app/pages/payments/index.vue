<script setup lang="ts">
const { get, post, patch, del } = useApi()

const payments = ref<any[]>([])
const bookings = ref<any[]>([])
const zones = ref<any[]>([])
const services = ref<any[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingItem = ref<any>(null)

const form = ref({
  booking_id: 0,
  amount: '',
  currency: 'RUB',
  payment_method: 'card' as string,
  status: 'pending' as string,
  receipt_number: '',
  paid_at: ''
})

const methods = [
  { value: 'cash', label: 'Наличные' },
  { value: 'card', label: 'Карта' },
  { value: 'online', label: 'Онлайн' }
]

const statuses = [
  { value: 'pending', label: 'Ожидание' },
  { value: 'paid', label: 'Оплачен' },
  { value: 'failed', label: 'Ошибка' },
  { value: 'refunded', label: 'Возврат' }
]

const statusBadge: Record<string, string> = {
  pending: 'badge-amber',
  paid: 'badge-green',
  failed: 'badge-red',
  refunded: 'badge-blue'
}

const statusLabel: Record<string, string> = {
  pending: 'Ожидание', paid: 'Оплачен', failed: 'Ошибка', refunded: 'Возврат'
}

const methodLabel: Record<string, string> = {
  cash: 'Наличные', card: 'Карта', online: 'Онлайн'
}

function zoneName(zid: number) { return zones.value.find((z: any) => z.id === zid)?.name || `#${zid}` }
function serviceName(sid: number) { return services.value.find((s: any) => s.id === sid)?.name || `#${sid}` }

function bookingLabel(bid: number) {
  const b = bookings.value.find((b: any) => b.id === bid)
  if (!b) return `Бронь #${bid}`
  return `#${b.id} — ${zoneName(b.zone_id)}, ${serviceName(b.service_id)}`
}

async function fetchAll() {
  loading.value = true
  try {
    const [pData, bData, zData, sData] = await Promise.all([
      get('/payment'), get('/booking'), get('/zone'), get('/service')
    ])
    payments.value = pData.payments || []
    bookings.value = bData.bookings || []
    zones.value = zData.zones || []
    services.value = sData.services || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

function openCreate() {
  editingItem.value = null
  form.value = {
    booking_id: bookings.value[0]?.id || 0,
    amount: '', currency: 'RUB', payment_method: 'card',
    status: 'pending', receipt_number: '', paid_at: ''
  }
  showModal.value = true
}

function openEdit(p: any) {
  editingItem.value = p
  form.value = {
    booking_id: p.booking_id,
    amount: p.amount?.String || String(p.amount) || '',
    currency: p.currency,
    payment_method: p.payment_method,
    status: p.status,
    receipt_number: p.receipt_number?.String || p.receipt_number || '',
    paid_at: p.paid_at?.Time || p.paid_at || ''
  }
  showModal.value = true
}

async function save() {
  try {
    const body: any = {
      booking_id: form.value.booking_id,
      amount: form.value.amount,
      currency: form.value.currency,
      payment_method: form.value.payment_method,
      status: form.value.status
    }
    if (form.value.receipt_number) body.receipt_number = form.value.receipt_number
    if (form.value.paid_at) body.paid_at = form.value.paid_at

    if (editingItem.value) {
      await patch(`/payment/${editingItem.value.id}`, body)
    } else {
      await post('/payment', body)
    }
    showModal.value = false
    await fetchAll()
  } catch (e: any) { alert(e?.data?.message || 'Ошибка сохранения') }
}

async function remove(id: number) {
  if (!confirm('Удалить платёж?')) return
  try { await del(`/payment/${id}`); await fetchAll() }
  catch (e: any) { alert(e?.data?.message || 'Ошибка удаления') }
}

function formatAmount(a: any) {
  if (!a) return '—'
  if (typeof a === 'string') return a
  if (a.String) return a.String
  return String(a)
}

function fmtTime(t: any) {
  if (!t) return '—'
  const raw = t?.Time || t
  if (!raw || raw === '0001-01-01T00:00:00Z') return '—'
  try { return new Date(raw).toLocaleString('ru-RU', { dateStyle: 'short', timeStyle: 'short' }) }
  catch { return String(raw) }
}

onMounted(fetchAll)
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Финансы</p>
        <h1 class="text-2xl font-bold text-white">Платежи</h1>
        <p class="text-gray-400">Приём, возврат и контроль оплат.</p>
      </div>
      <button class="ghost-button" @click="openCreate">+ Новый платёж</button>
    </div>

    <div v-if="loading" class="text-center text-gray-400 py-12">Загрузка…</div>
    <div v-else-if="payments.length === 0" class="glass p-8 text-center text-gray-400">Платежей пока нет.</div>

    <div v-else class="table-wrap overflow-x-auto">
      <table class="table-dark">
        <thead>
          <tr>
            <th>ID</th>
            <th>Бронирование</th>
            <th>Сумма</th>
            <th>Валюта</th>
            <th>Способ оплаты</th>
            <th>Статус</th>
            <th>Чек</th>
            <th>Дата оплаты</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in payments" :key="p.id">
            <td class="font-semibold text-white">#{{ p.id }}</td>
            <td class="text-gray-200">{{ bookingLabel(p.booking_id) }}</td>
            <td class="font-semibold text-gray-100">{{ formatAmount(p.amount) }}</td>
            <td class="text-gray-300">{{ p.currency }}</td>
            <td class="text-gray-300">{{ methodLabel[p.payment_method] || p.payment_method }}</td>
            <td><span :class="statusBadge[p.status] || 'badge'">{{ statusLabel[p.status] || p.status }}</span></td>
            <td class="text-gray-300">{{ p.receipt_number?.String || p.receipt_number || '—' }}</td>
            <td class="text-gray-300">{{ fmtTime(p.paid_at) }}</td>
            <td>
              <div class="flex gap-2">
                <button class="nav-link text-xs" @click="openEdit(p)">Изменить</button>
                <button class="text-red-400 text-xs hover:text-red-300" @click="remove(p.id)">Удалить</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay">
        <div class="modal-box">
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-lg font-bold text-white">{{ editingItem ? 'Изменить платёж' : 'Новый платёж' }}</h2>
            <button class="text-gray-400 hover:text-white transition" @click="showModal = false">✕</button>
          </div>
          <form class="grid gap-4" @submit.prevent="save">
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Бронирование
              <select v-model="form.booking_id" class="input-dark" required>
                <option value="" disabled>Выберите бронирование…</option>
                <option v-for="b in bookings" :key="b.id" :value="b.id">
                  #{{ b.id }} — {{ zoneName(b.zone_id) }}, {{ serviceName(b.service_id) }}
                </option>
              </select>
              <span class="text-xs text-gray-500">К какой брони привязать платёж</span>
            </label>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Сумма
                <input v-model="form.amount" class="input-dark" placeholder="1500.00" required />
                <span class="text-xs text-gray-500">Сколько списать / принять</span>
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Валюта
                <input v-model="form.currency" class="input-dark" />
              </label>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Способ оплаты
                <select v-model="form.payment_method" class="input-dark">
                  <option v-for="m in methods" :key="m.value" :value="m.value">{{ m.label }}</option>
                </select>
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Статус
                <select v-model="form.status" class="input-dark">
                  <option v-for="s in statuses" :key="s.value" :value="s.value">{{ s.label }}</option>
                </select>
              </label>
            </div>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Номер чека <span class="text-xs text-gray-500 font-normal">(необязательно)</span>
              <input v-model="form.receipt_number" class="input-dark" placeholder="RCP-001" />
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Дата оплаты <span class="text-xs text-gray-500 font-normal">(необязательно)</span>
              <input v-model="form.paid_at" type="datetime-local" lang="ru-RU" class="input-dark" />
            </label>
            <div class="flex gap-3 justify-end pt-2">
              <button type="button" class="ghost-button" @click="showModal = false">Отмена</button>
              <button type="submit" class="solid-button">{{ editingItem ? 'Сохранить' : 'Принять оплату' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<script setup lang="ts">
const { get, post, patch, del } = useApi()

const bookings = ref<any[]>([])
const users = ref<any[]>([])
const zones = ref<any[]>([])
const services = ref<any[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingItem = ref<any>(null)

const form = ref({
  user_id: '',
  zone_id: 0,
  service_id: 0,
  start_time: '',
  end_time: '',
  participants: 2,
  status: 'created' as string
})

const statuses = [
  { value: 'created', label: 'Создана' },
  { value: 'confirmed', label: 'Подтверждена' },
  { value: 'canceled', label: 'Отменена' },
  { value: 'completed', label: 'Завершена' }
]

const statusBadge: Record<string, string> = {
  created: 'badge-amber',
  confirmed: 'badge-green',
  canceled: 'badge-red',
  completed: 'badge-blue'
}

const statusLabel: Record<string, string> = {
  created: 'Создана',
  confirmed: 'Подтверждена',
  canceled: 'Отменена',
  completed: 'Завершена'
}

async function fetchAll() {
  loading.value = true
  try {
    const [bData, uData, zData, sData] = await Promise.all([
      get('/booking'), get('/user'), get('/zone'), get('/service')
    ])
    bookings.value = bData.bookings || []
    users.value = uData.users || []
    zones.value = zData.zones || []
    services.value = sData.services || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

function userName(uid: any) {
  const raw = uid?.Bytes ? bytesToUUID(uid.Bytes) : uid
  const u = users.value.find((u: any) => {
    const id = u.id?.Bytes ? bytesToUUID(u.id.Bytes) : u.id
    return id === raw
  })
  return u ? u.full_name : String(raw || '—').substring(0, 8) + '…'
}

function zoneName(zid: number) {
  return zones.value.find((z: any) => z.id === zid)?.name || `#${zid}`
}

function serviceName(sid: number) {
  return services.value.find((s: any) => s.id === sid)?.name || `#${sid}`
}

function bytesToUUID(bytes: number[]): string {
  if (!bytes || bytes.length !== 16) return ''
  const hex = bytes.map((b: number) => b.toString(16).padStart(2, '0')).join('')
  return `${hex.slice(0,8)}-${hex.slice(8,12)}-${hex.slice(12,16)}-${hex.slice(16,20)}-${hex.slice(20)}`
}

function userUUID(u: any): string {
  if (!u?.id) return ''
  return u.id.Bytes ? bytesToUUID(u.id.Bytes) : String(u.id)
}

function openCreate() {
  editingItem.value = null
  form.value = {
    user_id: users.value.length ? userUUID(users.value[0]) : '',
    zone_id: zones.value[0]?.id || 0,
    service_id: services.value[0]?.id || 0,
    start_time: '',
    end_time: '',
    participants: 2,
    status: 'created'
  }
  showModal.value = true
}

function openEdit(b: any) {
  editingItem.value = b
  const uid = b.user_id?.Bytes ? bytesToUUID(b.user_id.Bytes) : b.user_id
  form.value = {
    user_id: uid || '',
    zone_id: b.zone_id,
    service_id: b.service_id,
    start_time: b.start_time?.Time || b.start_time || '',
    end_time: b.end_time?.Time || b.end_time || '',
    participants: b.participants,
    status: b.status
  }
  showModal.value = true
}

async function save() {
  try {
    const body: any = {
      user_id: form.value.user_id || undefined,
      zone_id: form.value.zone_id,
      service_id: form.value.service_id,
      start_time: form.value.start_time,
      end_time: form.value.end_time,
      participants: form.value.participants,
      status: form.value.status
    }
    if (editingItem.value) {
      await patch(`/booking/${editingItem.value.id}`, body)
    } else {
      await post('/booking', body)
    }
    showModal.value = false
    await fetchAll()
  } catch (e: any) {
    alert(e?.data?.message || 'Ошибка сохранения')
  }
}

async function remove(id: number) {
  if (!confirm('Удалить бронирование?')) return
  try {
    await del(`/booking/${id}`)
    await fetchAll()
  } catch (e: any) { alert(e?.data?.message || 'Ошибка удаления') }
}

function fmtTime(t: any) {
  if (!t) return '—'
  const raw = t?.Time || t
  if (!raw) return '—'
  try { return new Date(raw).toLocaleString('ru-RU', { dateStyle: 'short', timeStyle: 'short' }) }
  catch { return String(raw) }
}

onMounted(fetchAll)
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Бронирования</p>
        <h1 class="text-2xl font-bold text-white">Управление бронями</h1>
        <p class="text-gray-400">Создание и редактирование бронирований.</p>
      </div>
      <button class="ghost-button" @click="openCreate">+ Новое бронирование</button>
    </div>

    <div v-if="loading" class="text-center text-gray-400 py-12">Загрузка…</div>
    <div v-else-if="bookings.length === 0" class="glass p-8 text-center text-gray-400">Бронирований пока нет.</div>

    <div v-else class="table-wrap overflow-x-auto">
      <table class="table-dark">
        <thead>
          <tr>
            <th>ID</th>
            <th>Клиент</th>
            <th>Зона</th>
            <th>Услуга</th>
            <th>Начало</th>
            <th>Конец</th>
            <th>Участники</th>
            <th>Статус</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in bookings" :key="b.id">
            <td class="font-semibold text-white">#{{ b.id }}</td>
            <td class="text-gray-200">{{ userName(b.user_id) }}</td>
            <td class="text-gray-300">{{ zoneName(b.zone_id) }}</td>
            <td class="text-gray-300">{{ serviceName(b.service_id) }}</td>
            <td class="text-gray-300">{{ fmtTime(b.start_time) }}</td>
            <td class="text-gray-300">{{ fmtTime(b.end_time) }}</td>
            <td class="text-gray-300">{{ b.participants }}</td>
            <td><span :class="statusBadge[b.status] || 'badge'">{{ statusLabel[b.status] || b.status }}</span></td>
            <td>
              <div class="flex gap-2">
                <button class="nav-link text-xs" @click="openEdit(b)">Изменить</button>
                <button class="text-red-400 text-xs hover:text-red-300" @click="remove(b.id)">Удалить</button>
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
            <h2 class="text-lg font-bold text-white">{{ editingItem ? 'Изменить бронирование' : 'Новое бронирование' }}</h2>
            <button class="text-gray-400 hover:text-white transition" @click="showModal = false">✕</button>
          </div>
          <form class="grid gap-4" @submit.prevent="save">
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Клиент
              <select v-model="form.user_id" class="input-dark" required>
                <option value="" disabled>Выберите клиента…</option>
                <option v-for="u in users" :key="userUUID(u)" :value="userUUID(u)">
                  {{ u.full_name }} ({{ u.email }})
                </option>
              </select>
              <span class="text-xs text-gray-500">Кто оформляет бронь</span>
            </label>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Зона
                <select v-model="form.zone_id" class="input-dark" required>
                  <option value="" disabled>Выберите зону…</option>
                  <option v-for="z in zones" :key="z.id" :value="z.id">{{ z.name }}</option>
                </select>
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Услуга
                <select v-model="form.service_id" class="input-dark" required>
                  <option value="" disabled>Выберите услугу…</option>
                  <option v-for="s in services" :key="s.id" :value="s.id">{{ s.name }}</option>
                </select>
              </label>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Начало
                <input v-model="form.start_time" type="datetime-local" lang="ru-RU" class="input-dark" required />
                <span class="text-xs text-gray-500">Когда начинается сеанс</span>
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Конец
                <input v-model="form.end_time" type="datetime-local" lang="ru-RU" class="input-dark" required />
              </label>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Участники
                <input v-model.number="form.participants" type="number" min="1" class="input-dark" required />
                <span class="text-xs text-gray-500">Сколько человек</span>
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Статус
                <select v-model="form.status" class="input-dark">
                  <option v-for="s in statuses" :key="s.value" :value="s.value">{{ s.label }}</option>
                </select>
              </label>
            </div>
            <div class="flex gap-3 justify-end pt-2">
              <button type="button" class="ghost-button" @click="showModal = false">Отмена</button>
              <button type="submit" class="solid-button">{{ editingItem ? 'Сохранить' : 'Создать бронь' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </section>
</template>

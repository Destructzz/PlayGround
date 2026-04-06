<script setup lang="ts">
const { get, post, patch, del } = useApi()

const staff = ref<any[]>([])
const users = ref<any[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingItem = ref<any>(null)

const form = ref({
  user_id: '',
  position: 'operator' as string,
  hire_date: '',
  phone: '',
  email: '',
  is_active: true
})

const positions = [
  { value: 'admin', label: 'Администратор' },
  { value: 'seller', label: 'Продавец' },
  { value: 'operator', label: 'Оператор' },
  { value: 'tech', label: 'Техник' }
]

function bytesToUUID(bytes: number[]): string {
  if (!bytes || bytes.length !== 16) return ''
  const hex = bytes.map((b: number) => b.toString(16).padStart(2, '0')).join('')
  return `${hex.slice(0,8)}-${hex.slice(8,12)}-${hex.slice(12,16)}-${hex.slice(16,20)}-${hex.slice(20)}`
}

function userUUID(u: any): string {
  if (!u?.id) return ''
  return u.id.Bytes ? bytesToUUID(u.id.Bytes) : String(u.id)
}

function userName(uid: any) {
  const raw = uid?.Bytes ? bytesToUUID(uid.Bytes) : uid
  const u = users.value.find((u: any) => userUUID(u) === raw)
  return u ? u.full_name : '—'
}

function positionLabel(pos: string) {
  return positions.find(p => p.value === pos)?.label || pos
}

async function fetchAll() {
  loading.value = true
  try {
    const [sData, uData] = await Promise.all([get('/staff'), get('/user')])
    staff.value = sData.staff || []
    users.value = uData.users || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

function openCreate() {
  editingItem.value = null
  form.value = {
    user_id: users.value.length ? userUUID(users.value[0]) : '',
    position: 'operator', hire_date: '', phone: '', email: '', is_active: true
  }
  showModal.value = true
}

function openEdit(s: any) {
  editingItem.value = s
  const uid = s.user_id?.Bytes ? bytesToUUID(s.user_id.Bytes) : s.user_id
  form.value = {
    user_id: uid || '',
    position: s.position,
    hire_date: s.hire_date?.Time ? new Date(s.hire_date.Time).toISOString().split('T')[0] : s.hire_date || '',
    phone: s.phone?.String || s.phone || '',
    email: s.email?.String || s.email || '',
    is_active: s.is_active
  }
  showModal.value = true
}

async function save() {
  try {
    if (editingItem.value) {
      await patch(`/staff/${editingItem.value.id}`, {
        position: form.value.position,
        hire_date: form.value.hire_date,
        phone: form.value.phone,
        email: form.value.email,
        is_active: form.value.is_active
      })
    } else {
      await post('/staff', {
        user_id: form.value.user_id,
        position: form.value.position,
        hire_date: form.value.hire_date,
        phone: form.value.phone,
        email: form.value.email,
        is_active: form.value.is_active
      })
    }
    showModal.value = false
    await fetchAll()
  } catch (e: any) { alert(e?.data?.message || 'Ошибка сохранения') }
}

async function remove(id: number) {
  if (!confirm('Удалить сотрудника?')) return
  try { await del(`/staff/${id}`); await fetchAll() }
  catch (e: any) { alert(e?.data?.message || 'Ошибка удаления') }
}

function fmtDate(d: any) {
  if (!d) return '—'
  const raw = d?.Time || d
  if (!raw) return '—'
  try { return new Date(raw).toLocaleDateString('ru-RU') }
  catch { return String(raw) }
}

onMounted(fetchAll)
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Команда</p>
        <h1 class="text-2xl font-bold text-white">Сотрудники</h1>
        <p class="text-gray-400">Управление персоналом досугового центра.</p>
      </div>
      <button class="ghost-button" @click="openCreate">+ Добавить сотрудника</button>
    </div>

    <div v-if="loading" class="text-center text-gray-400 py-12">Загрузка…</div>
    <div v-else-if="staff.length === 0" class="glass p-8 text-center text-gray-400">Сотрудников пока нет.</div>

    <div v-else class="table-wrap overflow-x-auto">
      <table class="table-dark">
        <thead>
          <tr>
            <th>ID</th>
            <th>Пользователь</th>
            <th>Должность</th>
            <th>Дата найма</th>
            <th>Телефон</th>
            <th>Email</th>
            <th>Статус</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in staff" :key="s.id">
            <td class="font-semibold text-white">#{{ s.id }}</td>
            <td class="text-gray-200">{{ userName(s.user_id) }}</td>
            <td class="text-gray-300">{{ positionLabel(s.position) }}</td>
            <td class="text-gray-300">{{ fmtDate(s.hire_date) }}</td>
            <td class="text-gray-300">{{ s.phone?.String || s.phone || '—' }}</td>
            <td class="text-gray-300">{{ s.email?.String || s.email || '—' }}</td>
            <td>
              <span :class="s.is_active ? 'badge-green' : 'badge-red'">
                {{ s.is_active ? 'Работает' : 'Уволен' }}
              </span>
            </td>
            <td>
              <div class="flex gap-2">
                <button class="nav-link text-xs" @click="openEdit(s)">Изменить</button>
                <button class="text-red-400 text-xs hover:text-red-300" @click="remove(s.id)">Удалить</button>
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
            <h2 class="text-lg font-bold text-white">{{ editingItem ? 'Изменить сотрудника' : 'Новый сотрудник' }}</h2>
            <button class="text-gray-400 hover:text-white transition" @click="showModal = false">✕</button>
          </div>
          <form class="grid gap-4" @submit.prevent="save">
            <label v-if="!editingItem" class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Пользователь
              <select v-model="form.user_id" class="input-dark" required>
                <option value="" disabled>Выберите пользователя…</option>
                <option v-for="u in users" :key="userUUID(u)" :value="userUUID(u)">
                  {{ u.full_name }} ({{ u.email }})
                </option>
              </select>
              <span class="text-xs text-gray-500">Привязка к аккаунту</span>
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Должность
              <select v-model="form.position" class="input-dark">
                <option v-for="p in positions" :key="p.value" :value="p.value">{{ p.label }}</option>
              </select>
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Дата найма
              <input v-model="form.hire_date" type="date" lang="ru-RU" class="input-dark" required />
              <span class="text-xs text-gray-500">Когда сотрудник вышел на работу</span>
            </label>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Телефон
                <input v-model="form.phone" class="input-dark" placeholder="+7 900 123-45-67" />
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Email
                <input v-model="form.email" type="email" class="input-dark" placeholder="staff@example.com" />
              </label>
            </div>
            <label class="flex items-center gap-2 text-sm text-gray-300">
              <input v-model="form.is_active" type="checkbox" class="rounded" />
              Активный (работает)
            </label>
            <div class="flex gap-3 justify-end pt-2">
              <button type="button" class="ghost-button" @click="showModal = false">Отмена</button>
              <button type="submit" class="solid-button">{{ editingItem ? 'Сохранить' : 'Добавить' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </section>
</template>

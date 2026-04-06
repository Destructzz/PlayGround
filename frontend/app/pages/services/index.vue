<script setup lang="ts">
const { get, post, patch, del } = useApi()

const services = ref<any[]>([])
const zones = ref<any[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingItem = ref<any>(null)

const form = ref({
  name: '',
  zone_id: 0,
  duration: 60,
  price: '',
  currency: 'RUB',
  description: '',
  is_active: true
})

async function fetchData() {
  loading.value = true
  try {
    const [sData, zData] = await Promise.all([get('/service'), get('/zone')])
    services.value = sData.services || []
    zones.value = zData.zones || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

function zoneName(zid: number) {
  return zones.value.find((z: any) => z.id === zid)?.name || `Зона #${zid}`
}

function openCreate() {
  editingItem.value = null
  form.value = { name: '', zone_id: zones.value[0]?.id || 0, duration: 60, price: '', currency: 'RUB', description: '', is_active: true }
  showModal.value = true
}

function openEdit(svc: any) {
  editingItem.value = svc
  form.value = {
    name: svc.name,
    zone_id: svc.zone_id,
    duration: svc.duration,
    price: svc.price?.String || svc.price?.toString() || '',
    currency: svc.currency,
    description: svc.description?.String || svc.description || '',
    is_active: svc.is_active
  }
  showModal.value = true
}

async function save() {
  try {
    const body = {
      name: form.value.name,
      zone_id: form.value.zone_id,
      duration: form.value.duration,
      price: form.value.price,
      currency: form.value.currency,
      description: form.value.description,
      is_active: form.value.is_active
    }
    if (editingItem.value) {
      await patch(`/service/${editingItem.value.id}`, body)
    } else {
      await post('/service', body)
    }
    showModal.value = false
    await fetchData()
  } catch (e: any) { alert(e?.data?.message || 'Ошибка сохранения') }
}

async function remove(id: number) {
  if (!confirm('Удалить услугу?')) return
  try { await del(`/service/${id}`); await fetchData() }
  catch (e: any) { alert(e?.data?.message || 'Ошибка удаления') }
}

function formatPrice(p: any) {
  if (!p) return '—'
  if (typeof p === 'string') return p
  if (p.String) return p.String
  return String(p)
}

onMounted(fetchData)
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Каталог</p>
        <h1 class="text-2xl font-bold text-white">Услуги PlayGround</h1>
        <p class="text-gray-400">Управление услугами, ценами и зонами.</p>
      </div>
      <button class="ghost-button" @click="openCreate">+ Добавить услугу</button>
    </div>

    <div v-if="loading" class="text-center text-gray-400 py-12">Загрузка…</div>
    <div v-else-if="services.length === 0" class="glass p-8 text-center text-gray-400">Услуг пока нет. Создайте первую!</div>

    <div v-else class="grid gap-4 md:grid-cols-2">
      <article v-for="svc in services" :key="svc.id" class="glass p-4">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-white">{{ svc.name }}</h2>
            <p class="text-sm text-gray-400">{{ zoneName(svc.zone_id) }} · {{ svc.duration }} мин</p>
          </div>
          <div class="text-right">
            <div class="text-base font-semibold text-gray-100">{{ formatPrice(svc.price) }} {{ svc.currency }}</div>
            <span :class="svc.is_active ? 'badge-green' : 'badge-red'" class="mt-1">
              {{ svc.is_active ? 'Активна' : 'Скрыта' }}
            </span>
          </div>
        </div>
        <p v-if="svc.description?.String || svc.description" class="mt-2 text-sm text-gray-400">
          {{ svc.description?.String || svc.description }}
        </p>
        <div class="mt-3 flex items-center gap-3 text-sm">
          <button class="ghost-button text-xs !px-3 !py-1.5" @click="openEdit(svc)">Изменить</button>
          <button class="danger-button" @click="remove(svc.id)">Удалить</button>
        </div>
      </article>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay">
        <div class="modal-box">
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-lg font-bold text-white">{{ editingItem ? 'Изменить услугу' : 'Новая услуга' }}</h2>
            <button class="text-gray-400 hover:text-white transition" @click="showModal = false">✕</button>
          </div>
          <form class="grid gap-4" @submit.prevent="save">
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Название услуги
              <input v-model="form.name" class="input-dark" placeholder="VR-сессия, Консоли…" required />
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Зона / помещение
              <select v-model="form.zone_id" class="input-dark" required>
                <option value="" disabled>Выберите зону…</option>
                <option v-for="z in zones" :key="z.id" :value="z.id">{{ z.name }}</option>
              </select>
              <span class="text-xs text-gray-500">Где проводится услуга</span>
            </label>
            <div class="grid grid-cols-2 gap-4">
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Длительность (мин)
                <input v-model.number="form.duration" type="number" min="1" class="input-dark" required />
              </label>
              <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
                Цена
                <input v-model="form.price" class="input-dark" placeholder="1500.00" required />
              </label>
            </div>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Валюта
              <input v-model="form.currency" class="input-dark" />
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Описание <span class="text-xs text-gray-500 font-normal">(необязательно)</span>
              <textarea v-model="form.description" rows="2" class="input-dark" placeholder="Подробнее об услуге…"></textarea>
            </label>
            <label class="flex items-center gap-2 text-sm text-gray-300">
              <input v-model="form.is_active" type="checkbox" class="rounded" />
              Показывать в каталоге
            </label>
            <div class="flex gap-3 justify-end pt-2">
              <button type="button" class="ghost-button" @click="showModal = false">Отмена</button>
              <button type="submit" class="solid-button">{{ editingItem ? 'Сохранить' : 'Создать услугу' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </section>
</template>

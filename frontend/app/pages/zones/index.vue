<script setup lang="ts">
const { get, post, patch, del } = useApi()

const zones = ref<any[]>([])
const loading = ref(true)
const showModal = ref(false)
const editingZone = ref<any>(null)

const form = ref({
  name: '',
  type: 'game' as string,
  capacity: 1,
  description: '',
  is_active: true
})

const zoneTypes = [
  { value: 'game', label: 'Игровая' },
  { value: 'event', label: 'Мероприятия' },
  { value: 'vip', label: 'VIP' },
  { value: 'lounge', label: 'Лаунж' },
  { value: 'sys', label: 'Сервисная' }
]

const typeLabel: Record<string, string> = {
  game: 'Игровая', event: 'Мероприятия', vip: 'VIP', lounge: 'Лаунж', sys: 'Сервисная'
}

async function fetchZones() {
  loading.value = true
  try {
    const data = await get('/zone')
    zones.value = data.zones || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

function openCreate() {
  editingZone.value = null
  form.value = { name: '', type: 'game', capacity: 1, description: '', is_active: true }
  showModal.value = true
}

function openEdit(zone: any) {
  editingZone.value = zone
  form.value = {
    name: zone.name,
    type: zone.zone_type,
    capacity: zone.capacity,
    description: zone.description?.String || zone.description || '',
    is_active: zone.is_active
  }
  showModal.value = true
}

async function save() {
  try {
    const body = {
      name: form.value.name, type: form.value.type,
      capacity: form.value.capacity, description: form.value.description,
      is_active: form.value.is_active
    }
    if (editingZone.value) {
      await patch(`/zone/${editingZone.value.id}`, body)
    } else {
      await post('/zone', body)
    }
    showModal.value = false
    await fetchZones()
  } catch (e: any) { alert(e?.data?.message || 'Ошибка сохранения') }
}

async function remove(id: number) {
  if (!confirm('Удалить зону?')) return
  try { await del(`/zone/${id}`); await fetchZones() }
  catch (e: any) { alert(e?.data?.message || 'Ошибка удаления') }
}

onMounted(fetchZones)
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p class="section-title">Ресурсы</p>
        <h1 class="text-2xl font-bold text-white">Зоны и помещения</h1>
        <p class="text-gray-400">Создание, редактирование и удаление зон досугового центра.</p>
      </div>
      <button class="ghost-button" @click="openCreate">+ Добавить зону</button>
    </div>

    <div v-if="loading" class="text-center text-gray-400 py-12">Загрузка…</div>
    <div v-else-if="zones.length === 0" class="glass p-8 text-center text-gray-400">Зон пока нет. Создайте первую!</div>

    <div v-else class="grid gap-4 md:grid-cols-2">
      <article v-for="zone in zones" :key="zone.id" class="glass p-4">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-white">{{ zone.name }}</h2>
            <p class="text-sm text-gray-400">{{ typeLabel[zone.zone_type] || zone.zone_type }} · до {{ zone.capacity }} чел.</p>
          </div>
          <span :class="zone.is_active ? 'badge-green' : 'badge-red'">
            {{ zone.is_active ? 'Открыта' : 'Закрыта' }}
          </span>
        </div>
        <p v-if="zone.description?.String || zone.description" class="mt-2 text-sm text-gray-400">
          {{ zone.description?.String || zone.description }}
        </p>
        <div class="mt-3 flex items-center gap-3 text-sm">
          <button class="ghost-button text-xs !px-3 !py-1.5" @click="openEdit(zone)">Изменить</button>
          <button class="danger-button" @click="remove(zone.id)">Удалить</button>
        </div>
      </article>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay">
        <div class="modal-box">
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-lg font-bold text-white">{{ editingZone ? 'Изменить зону' : 'Новая зона' }}</h2>
            <button class="text-gray-400 hover:text-white transition" @click="showModal = false">✕</button>
          </div>
          <form class="grid gap-4" @submit.prevent="save">
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Название
              <input v-model="form.name" class="input-dark" placeholder="VR-арена, Лаунж-бар…" required />
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Тип зоны
              <select v-model="form.type" class="input-dark">
                <option v-for="t in zoneTypes" :key="t.value" :value="t.value">{{ t.label }}</option>
              </select>
              <span class="text-xs text-gray-500">Определяет категорию зоны</span>
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Вместимость (чел.)
              <input v-model.number="form.capacity" type="number" min="1" class="input-dark" required />
              <span class="text-xs text-gray-500">Максимальное количество гостей</span>
            </label>
            <label class="flex flex-col gap-1.5 text-sm font-medium text-gray-300">
              Описание <span class="text-xs text-gray-500 font-normal">(необязательно)</span>
              <textarea v-model="form.description" rows="2" class="input-dark" placeholder="Подробнее о зоне…"></textarea>
            </label>
            <label class="flex items-center gap-2 text-sm text-gray-300">
              <input v-model="form.is_active" type="checkbox" class="rounded" />
              Зона открыта для посетителей
            </label>
            <div class="flex gap-3 justify-end pt-2">
              <button type="button" class="ghost-button" @click="showModal = false">Отмена</button>
              <button type="submit" class="solid-button">{{ editingZone ? 'Сохранить' : 'Создать зону' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </section>
</template>

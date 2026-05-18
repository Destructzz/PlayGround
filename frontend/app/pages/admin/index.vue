<template>
  <div class="space-y-6">
    <!-- Quick Stats Grid (Horizontal widescreen view) -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div
        v-for="stat in dashboardStats"
        :key="stat.label"
        class="rounded-[1.25rem] border border-white/5 bg-[#050f17] p-5 text-center shadow-xl transition-transform hover:scale-[1.02]"
      >
        <p class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-500">
          {{ stat.label }}
        </p>
        <p class="mt-2 text-2xl font-black text-white leading-none">
          {{ stat.value }}
        </p>
      </div>
    </div>

    <!-- Header Card with responsive action button -->
    <header class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between rounded-[1.25rem] border border-white/5 bg-[#050f17] p-8 shadow-2xl">
      <div>
        <p class="text-[10px] font-black uppercase tracking-[0.4em] text-cyan-300/50">{{ activeTabItem.eyebrow }}</p>
        <h1 class="mt-2 text-3xl font-black tracking-tight text-white">{{ activeTabItem.title }}</h1>
        <p class="mt-3 text-sm text-zinc-400 max-w-2xl leading-relaxed">{{ activeTabItem.description }}</p>
      </div>
      <div class="flex items-center gap-4 flex-shrink-0">
        <button
          v-if="showCreateAction"
          type="button"
          class="rounded-[0.9rem] bg-cyan-300 px-6 py-3 text-xs font-black uppercase tracking-widest text-[#020c13] transition hover:bg-cyan-200 hover:shadow-[0_0_20px_rgba(34,211,238,0.4)] active:scale-95"
          @click="openCreateModal"
        >
          {{ createModalMeta.buttonLabel }}
        </button>
      </div>
    </header>

    <!-- Feedback Message -->
    <div
      v-if="feedbackMessage"
      class="rounded-[0.8rem] border px-4 py-3 text-sm shadow-lg"
      :class="feedbackTone === 'error' ? 'border-orange-300/30 bg-orange-500/10 text-orange-100' : 'border-emerald-300/30 bg-emerald-500/10 text-emerald-100'"
    >
      {{ feedbackMessage }}
    </div>

    <!-- Data Loading State -->
    <section
      v-if="isLoading"
      class="rounded-[1rem] border border-white/5 bg-[#050f17] p-12 text-center text-sm font-medium text-zinc-400"
    >
      Загрузка системных данных...
    </section>

    <!-- Error State -->
    <section
      v-else-if="pageError"
      class="rounded-[1rem] border border-orange-300/20 bg-orange-500/10 p-8"
    >
      <p class="text-sm font-bold text-orange-200">Ошибка загрузки: {{ pageError }}</p>
      <button
        class="mt-4 rounded-full border border-white/20 px-6 py-2 text-xs font-bold text-white transition hover:bg-white/10"
        @click="loadAdminData()"
      >
        Повторить попытку
      </button>
    </section>

    <!-- Content Sections -->
    <template v-else>

      <!-- Zones Tab -->
      <section
        v-if="activeTab === 'zones'"
        class="space-y-4"
      >
        <article class="space-y-4">
          <div
            v-if="!sortedZones.length"
            class="border border-white/8 bg-[#07141d] px-4 py-5 text-sm text-zinc-300"
          >
            Пока нет зон. Начни с создания первой записи через кнопку выше.
          </div>

          <div
            v-for="zone in sortedZones"
            :key="zone.id"
            class="overflow-hidden border border-white/8 bg-[#07141d] shadow-[0_12px_30px_rgba(0,0,0,0.24)]"
          >
            <div class="border-b border-white/8 px-4 py-4">
              <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <span
                      class="rounded-full border px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em]"
                      :class="zoneBadgeClass(zone.zoneType)"
                    >
                      {{ zoneTypeLabel(zone.zoneType) }}
                    </span>
                    <span class="rounded-full border border-white/10 bg-white/4 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-cyan-100/55">
                      {{ zoneTagName(zone.zoneTagId) }}
                    </span>
                    <span
                      v-if="zone.isActive"
                      class="rounded-full border border-emerald-300/20 bg-emerald-500/12 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-emerald-100/80"
                    >
                      active
                    </span>
                    <span
                      v-else
                      class="rounded-full border border-zinc-300/15 bg-zinc-500/10 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-zinc-300/70"
                    >
                      inactive
                    </span>
                  </div>
                  <h3 class="mt-4 text-2xl font-black text-white">
                    {{ zone.name }}
                  </h3>
                  <p class="mt-2 text-sm leading-7 text-zinc-300">
                    {{ zone.description || 'Описание не заполнено.' }}
                  </p>
                </div>

                <div class="flex flex-wrap items-center gap-3">
                  <div class="border border-white/8 bg-white/4 px-3 py-2 text-right">
                    <p class="text-[11px] uppercase tracking-[0.28em] text-cyan-100/45">
                      Capacity
                    </p>
                    <p class="mt-2 text-2xl font-black text-white">
                      {{ zone.capacity }}
                    </p>
                  </div>
                  <button
                    type="button"
                    class="rounded-[0.8rem] border border-cyan-400/18 px-4 py-2.5 text-sm font-bold text-white transition hover:border-cyan-300/40"
                    @click="toggleZoneEditor(zone.id)"
                  >
                    {{ expandedZoneId === zone.id ? 'Скрыть форму' : 'Редактировать' }}
                  </button>
                  <button
                    type="button"
                    class="rounded-[0.8rem] border border-orange-300/18 px-4 py-2.5 text-sm font-bold text-orange-100 transition hover:border-orange-200/40"
                    :disabled="isMutating"
                    @click="removeZone(zone.id, zone.name)"
                  >
                    Удалить
                  </button>
                </div>
              </div>

              <div
                v-if="hasZoneDetails(zone.detailsJson)"
                class="mt-4 border border-white/8 bg-[#06131c] p-3"
              >
                <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
                  Details JSON
                </p>
                <pre class="mt-3 overflow-x-auto text-xs leading-6 text-zinc-300">{{ zone.detailsJson }}</pre>
              </div>
            </div>

            <div
              v-if="expandedZoneId === zone.id"
              class="border-b border-white/8 px-4 py-4"
            >
              <form
                class="space-y-4"
                @submit.prevent="updateZone(zone.id)"
              >
                <div class="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Название</label>
                    <input
                      v-model="ensureZoneDraft(zone.id).name"
                      type="text"
                      class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                    >
                  </div>
                  <div>
                    <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Тип зоны</label>
                    <select
                      v-model="ensureZoneDraft(zone.id).type"
                      class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                    >
                      <option
                        v-for="option in zoneTypeOptions"
                        :key="option.value"
                        :value="option.value"
                      >
                        {{ option.label }}
                      </option>
                    </select>
                  </div>
                </div>

                <div class="grid gap-4 lg:grid-cols-[1fr_auto]">
                  <div>
                    <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Zone tag</label>
                    <select
                      v-model="ensureZoneDraft(zone.id).zoneTagId"
                      class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                    >
                      <option
                        disabled
                        value=""
                      >
                        Выбери tag
                      </option>
                      <option
                        v-for="tag in sortedZoneTags"
                        :key="tag.id"
                        :value="String(tag.id)"
                      >
                        {{ tag.name }}
                      </option>
                    </select>
                  </div>

                  <div class="border border-white/8 bg-white/4 p-3 lg:min-w-64">
                    <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
                      Создать tag здесь
                    </p>
                    <label class="mb-2 mt-3 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Название тега</label>
                    <div class="flex gap-2">
                      <input
                        v-model="ensureZoneEditTagNames(zone.id)[zone.id]"
                        type="text"
                        placeholder="Новый tag"
                        class="min-w-0 flex-1 rounded-[0.75rem] border border-cyan-400/18 bg-[#06131c] px-3 py-2.5 text-sm text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
                      >
                      <button
                        type="button"
                        class="rounded-[0.75rem] bg-cyan-300 px-4 py-2.5 text-sm font-black text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
                        :disabled="isMutating"
                        @click="createZoneTagFromZoneContext(ensureZoneEditTagNames(zone.id)[zone.id] ?? '', zone.id)"
                      >
                        Добавить
                      </button>
                    </div>
                  </div>
                </div>

                <div class="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Вместимость</label>
                    <input
                      v-model="ensureZoneDraft(zone.id).capacity"
                      type="number"
                      min="1"
                      class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                    >
                  </div>
                  <label class="flex items-center gap-3 rounded-[0.9rem] border border-white/8 bg-white/4 px-4 py-3 text-sm text-zinc-200">
                    <input
                      v-model="ensureZoneDraft(zone.id).isActive"
                      type="checkbox"
                      class="h-4 w-4 accent-cyan-300"
                    >
                    Зона активна
                  </label>
                </div>

                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Описание</label>
                  <textarea
                    v-model="ensureZoneDraft(zone.id).description"
                    rows="4"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                  />
                </div>

                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Details JSON</label>
                  <textarea
                    v-model="ensureZoneDraft(zone.id).detailsJson"
                    rows="6"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 font-mono text-sm text-white focus:border-cyan-300 focus:outline-none"
                  />
                </div>

                <div class="flex flex-wrap gap-3">
                  <button
                    type="submit"
                    class="rounded-[0.85rem] bg-cyan-300 px-5 py-3 text-sm font-black uppercase tracking-[0.28em] text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
                    :disabled="isMutating"
                  >
                    Сохранить изменения
                  </button>
                  <button
                    type="button"
                    class="rounded-[0.85rem] border border-white/10 px-5 py-3 text-sm font-bold text-white transition hover:border-white/25"
                    :disabled="isMutating"
                    @click="resetZoneDraft(zone.id)"
                  >
                    Сбросить
                  </button>
                </div>
              </form>
            </div>

            <div class="px-4 py-4">
              <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
                <div>
                  <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
                    Places inside this zone
                  </p>
                  <p class="mt-2 text-sm leading-7 text-zinc-300">
                    Все места уже привязаны к текущей зоне. Здесь не нужно думать про `zone_id` и архитектуру отдельно.
                  </p>
                </div>
                <span class="rounded-full border border-white/10 bg-white/4 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-cyan-100/55">
                  {{ placesForZone(zone.id).length }} мест
                </span>
              </div>

              <form
                class="mt-4 grid gap-3 border border-white/8 bg-[#06131c] p-3 lg:grid-cols-[1.15fr_0.75fr_0.75fr_auto]"
                @submit.prevent="submitPlace(zone.id)"
              >
                <div>
                  <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Label</label>
                  <input
                    v-model="ensureNewPlaceForm(zone.id).label"
                    type="text"
                    placeholder="Например, PC-07"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
                  >
                </div>
                <div>
                  <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Sort order</label>
                  <input
                    v-model="ensureNewPlaceForm(zone.id).sortOrder"
                    type="number"
                    min="0"
                    placeholder="0"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                  >
                </div>
                <div>
                  <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Configuration ID</label>
                  <input
                    v-model="ensureNewPlaceForm(zone.id).configurationId"
                    type="number"
                    min="1"
                    placeholder="Если нужен"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
                  >
                </div>
                <div class="lg:self-end">
                  <button
                    type="submit"
                    class="w-full rounded-[0.8rem] bg-white px-4 py-3 text-sm font-black text-[#020c13] transition hover:bg-zinc-200 disabled:cursor-not-allowed disabled:bg-white/50"
                    :disabled="isMutating"
                  >
                    Добавить место
                  </button>
                </div>
                <label class="lg:col-span-4 flex items-center gap-3 text-sm text-zinc-300">
                  <input
                    v-model="ensureNewPlaceForm(zone.id).isActive"
                    type="checkbox"
                    class="h-4 w-4 accent-cyan-300"
                  >
                  Место активно. `Configuration ID` можно оставить пустым, если для зоны он не нужен.
                </label>
              </form>

              <div
                v-if="!placesForZone(zone.id).length"
                class="mt-4 border border-white/8 bg-white/4 px-4 py-3 text-sm text-zinc-300"
              >
                У этой зоны пока нет мест.
              </div>

              <div
                v-else
                class="mt-4 grid gap-3"
              >
                <div
                  v-for="place in placesForZone(zone.id)"
                  :key="place.id"
                  class="border border-white/8 bg-white/4 p-3"
                >
                  <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                    <div>
                      <div class="flex flex-wrap items-center gap-2">
                        <span class="rounded-full border border-cyan-300/18 bg-cyan-400/10 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-cyan-100/70">
                          {{ place.label }}
                        </span>
                        <span class="rounded-full border border-white/10 bg-white/4 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-zinc-300/75">
                          sort {{ place.sortOrder }}
                        </span>
                        <span
                          v-if="place.isActive"
                          class="rounded-full border border-emerald-300/20 bg-emerald-500/12 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-emerald-100/80"
                        >
                          active
                        </span>
                        <span
                          v-else
                          class="rounded-full border border-zinc-300/15 bg-zinc-500/10 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-zinc-300/70"
                        >
                          inactive
                        </span>
                      </div>
                      <p class="mt-3 text-sm text-zinc-300">
                        Configuration ID: <span class="font-bold text-white">{{ place.configurationId ?? 'не задан' }}</span>
                      </p>
                    </div>

                    <div class="flex flex-wrap gap-3">
                      <button
                        type="button"
                        class="rounded-[0.75rem] border border-cyan-400/18 px-4 py-2.5 text-sm font-bold text-white transition hover:border-cyan-300/40"
                        @click="togglePlaceEditor(place.id)"
                      >
                        {{ expandedPlaceId === place.id ? 'Скрыть форму' : 'Редактировать' }}
                      </button>
                      <button
                        type="button"
                        class="rounded-[0.75rem] border border-orange-300/18 px-4 py-2.5 text-sm font-bold text-orange-100 transition hover:border-orange-200/40"
                        :disabled="isMutating"
                        @click="removePlace(place.id, place.label)"
                      >
                        Удалить
                      </button>
                    </div>
                  </div>

                  <form
                    v-if="expandedPlaceId === place.id"
                    class="mt-3 grid gap-3 lg:grid-cols-[1fr_0.7fr_0.7fr_auto]"
                    @submit.prevent="updatePlace(place.id)"
                  >
                    <div>
                      <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Label</label>
                      <input
                        v-model="ensurePlaceDraft(place.id).label"
                        type="text"
                        class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                      >
                    </div>
                    <div>
                      <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Sort order</label>
                      <input
                        v-model="ensurePlaceDraft(place.id).sortOrder"
                        type="number"
                        min="0"
                        class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                      >
                    </div>
                    <div>
                      <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Configuration ID</label>
                      <input
                        v-model="ensurePlaceDraft(place.id).configurationId"
                        type="number"
                        min="1"
                        placeholder="Если нужен"
                        class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#071926] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
                      >
                    </div>
                    <div class="lg:self-end">
                      <button
                        type="submit"
                        class="w-full rounded-[0.8rem] bg-cyan-300 px-4 py-3 text-sm font-black text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
                        :disabled="isMutating"
                      >
                        Сохранить
                      </button>
                    </div>
                    <label class="lg:col-span-4 flex items-center gap-3 text-sm text-zinc-300">
                      <input
                        v-model="ensurePlaceDraft(place.id).isActive"
                        type="checkbox"
                        class="h-4 w-4 accent-cyan-300"
                      >
                      Место активно
                    </label>
                  </form>
                </div>
              </div>
            </div>
          </div>
        </article>
      </section>

      <!-- Zone Tags Tab -->
      <section
        v-if="activeTab === 'zone-tags'"
        class="space-y-4"
      >
        <article class="space-y-4">
          <div
            v-if="!sortedZoneTags.length"
            class="border border-white/8 bg-[#07141d] px-4 py-5 text-sm text-zinc-300"
          >
            Теги ещё не созданы.
          </div>

          <div
            v-for="tag in sortedZoneTags"
            :key="tag.id"
            class="border border-white/8 bg-[#07141d] p-4 shadow-[0_12px_30px_rgba(0,0,0,0.24)]"
          >
            <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
                  Tag #{{ tag.id }}
                </p>
                <h3 class="mt-2 text-2xl font-black text-white">
                  {{ tag.name }}
                </h3>
              </div>

              <div class="flex flex-wrap gap-3">
                <button
                  type="button"
                  class="rounded-[0.75rem] border border-orange-300/18 px-4 py-2.5 text-sm font-bold text-white transition hover:border-orange-200/40"
                  @click="toggleZoneTagEditor(tag.id)"
                >
                  {{ expandedZoneTagId === tag.id ? 'Скрыть форму' : 'Редактировать' }}
                </button>
                <button
                  type="button"
                  class="rounded-[0.75rem] border border-orange-300/18 px-4 py-2.5 text-sm font-bold text-orange-100 transition hover:border-orange-200/40"
                  :disabled="isMutating"
                  @click="removeZoneTag(tag.id, tag.name)"
                >
                  Удалить
                </button>
              </div>
            </div>

            <form
              v-if="expandedZoneTagId === tag.id"
              class="mt-3 flex flex-col gap-3 sm:flex-row"
              @submit.prevent="updateZoneTag(tag.id)"
            >
              <div class="min-w-0 flex-1">
                <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Название тега</label>
                <input
                  v-model="zoneTagDrafts[tag.id]"
                  type="text"
                  class="min-w-0 w-full rounded-[0.8rem] border border-orange-300/18 bg-[#06131c] px-4 py-3 text-white focus:border-orange-300 focus:outline-none"
                />
              </div>
              <button
                type="submit"
                class="rounded-[0.8rem] bg-orange-300 px-5 py-3 text-sm font-black text-[#020c13] transition hover:bg-orange-200 disabled:cursor-not-allowed disabled:bg-orange-300/50"
                :disabled="isMutating"
              >
                Сохранить
              </button>
            </form>
          </div>
        </article>
      </section>

      <!-- Bookings Tab -->
      <section
        v-if="activeTab === 'bookings'"
        class="space-y-4"
      >
        <article class="space-y-4">
          <div
            v-if="!sortedBookings.length"
            class="border border-white/8 bg-[#07141d] px-4 py-5 text-sm text-zinc-300"
          >
            Брони пока не появились.
          </div>

          <div
            v-for="booking in sortedBookings"
            :key="booking.id"
            class="border border-white/8 bg-[#07141d] p-4 shadow-[0_12px_30px_rgba(0,0,0,0.24)]"
          >
            <div class="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
              <div>
                <div class="flex flex-wrap items-center gap-2">
                  <span
                    class="rounded-full border px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em]"
                    :class="bookingStatusClass(booking.status)"
                  >
                    {{ booking.status }}
                  </span>
                  <span class="rounded-full border border-white/10 bg-white/4 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-zinc-300/75">
                    Booking #{{ booking.id }}
                  </span>
                </div>

                <h3 class="mt-4 text-2xl font-black text-white">
                  {{ zoneNameById(booking.zoneId) }} / {{ placeLabelById(booking.placeId) }}
                </h3>
                <p class="mt-2 text-sm leading-7 text-zinc-300">
                  {{ serviceNameById(booking.serviceId) }} • {{ booking.totalPrice }} RUB • {{ booking.participants }} участ.
                </p>
                <p class="mt-2 text-sm text-zinc-400">
                  {{ formatDateRange(booking.startTime, booking.endTime) }}
                </p>
              </div>

              <form
                class="flex flex-col gap-3 sm:min-w-[260px]"
                @submit.prevent="updateBookingStatus(booking.id)"
              >
                <div>
                  <label class="mb-2 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Статус</label>
                  <select
                    v-model="bookingStatusDrafts[booking.id]"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                  >
                    <option
                      v-for="status in bookingStatusOptions"
                      :key="status"
                      :value="status"
                    >
                      {{ status }}
                    </option>
                  </select>
                </div>
                <div class="flex gap-2">
                  <button
                    type="submit"
                    class="flex-1 rounded-[0.8rem] bg-cyan-300 px-4 py-3 text-sm font-black text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
                    :disabled="isMutating"
                  >
                    Сохранить
                  </button>
                  <button
                    type="button"
                    class="flex-1 rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-sm font-bold text-white transition hover:bg-cyan-900/50"
                    @click="toggleBookingDetails(booking.id)"
                  >
                    {{ expandedBookingDetails === booking.id ? 'Скрыть' : 'Подробнее' }}
                  </button>
                </div>
              </form>
            </div>

            <template v-if="expandedBookingDetails === booking.id">
              <div class="mt-4 grid gap-3 lg:grid-cols-2">
                <div class="border border-white/8 bg-white/4 p-3 text-sm text-zinc-300">
                  <p>
                    <span class="text-zinc-500">Контакт:</span> <span class="text-white">{{ booking.contactName || 'Не указан' }}</span>
                  </p>
                  <p class="mt-2">
                    <span class="text-zinc-500">Email:</span> <span class="text-white">{{ booking.contactEmail || 'Не указан' }}</span>
                  </p>
                  <p class="mt-2">
                    <span class="text-zinc-500">Телефон:</span> <span class="text-white">{{ booking.contactPhone || 'Не указан' }}</span>
                  </p>
                </div>
                <div class="border border-white/8 bg-white/4 p-3 text-sm text-zinc-300">
                  <p>
                    <span class="text-zinc-500">User ID:</span> <span class="break-all text-white">{{ booking.userId }}</span>
                  </p>
                  <p class="mt-2">
                    <span class="text-zinc-500">Создана:</span> <span class="text-white">{{ formatAdminDateTime(booking.createdAt) }}</span>
                  </p>
                  <p class="mt-2">
                    <span class="text-zinc-500">Обновлена:</span> <span class="text-white">{{ formatAdminDateTime(booking.updatedAt) }}</span>
                  </p>
                </div>
              </div>

              <div
                v-if="hasZoneDetails(booking.detailsJson)"
                class="mt-4 border border-white/8 bg-[#06131c] p-3"
              >
                <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
                  Details JSON
                </p>
                <pre class="mt-3 overflow-x-auto text-xs leading-6 text-zinc-300">{{ booking.detailsJson }}</pre>
              </div>
            </template>
          </div>
        </article>
      </section>

      <!-- Shifts Tab -->
      <section
        v-if="activeTab === 'shifts'"
        class="space-y-4"
      >
        <article class="space-y-4">
          <div
            v-if="!sortedShifts.length"
            class="border border-white/8 bg-[#07141d] px-4 py-5 text-sm text-zinc-300"
          >
            Смен пока нет.
          </div>

          <div
            v-for="shift in sortedShifts"
            :key="shift.id"
            class="border border-white/8 bg-[#07141d] p-4 shadow-[0_12px_30px_rgba(0,0,0,0.24)]"
          >
            <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
              <div>
                <div class="flex flex-wrap items-center gap-2">
                  <span class="rounded-full border border-cyan-300/18 bg-cyan-400/10 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-cyan-100/70">
                    {{ zoneTagName(shift.zone_tag_id ?? undefined) }}
                  </span>
                  <span class="rounded-full border border-white/10 bg-white/4 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.28em] text-zinc-300/75">
                    Shift #{{ shift.id }}
                  </span>
                </div>
                <h3 class="mt-4 text-2xl font-black text-white">
                  {{ formatDateRange(shift.start_time, shift.end_time) }}
                </h3>

                <template v-if="expandedShiftDetails === shift.id">
                  <p class="mt-2 text-sm leading-7 text-zinc-300">
                    {{ shift.note || 'Без заметки.' }}
                  </p>
                  <p class="mt-3 text-sm text-zinc-400">
                    Создал: <span class="font-semibold text-white">{{ shift.user.full_name }}</span>
                    <span class="text-zinc-500">{{ shift.user.email }}</span>
                  </p>
                </template>
              </div>

              <div class="flex flex-wrap gap-3">
                <button
                  type="button"
                  class="rounded-[0.75rem] border border-cyan-400/18 bg-cyan-900/20 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-cyan-900/40"
                  @click="toggleShiftDetails(shift.id)"
                >
                  {{ expandedShiftDetails === shift.id ? 'Скрыть детали' : 'Подробнее' }}
                </button>
                <button
                  type="button"
                  class="rounded-[0.75rem] border border-cyan-400/18 px-4 py-2.5 text-sm font-bold text-white transition hover:border-cyan-300/40"
                  @click="toggleShiftEditor(shift.id)"
                >
                  {{ expandedShiftId === shift.id ? 'Скрыть форму' : 'Редактировать' }}
                </button>
                <button
                  type="button"
                  class="rounded-[0.75rem] border border-orange-300/18 px-4 py-2.5 text-sm font-bold text-orange-100 transition hover:border-orange-200/40"
                  :disabled="isMutating"
                  @click="removeShift(shift.id)"
                >
                  Удалить
                </button>
              </div>
            </div>

            <form
              v-if="expandedShiftId === shift.id"
              class="mt-4 space-y-4"
              @submit.prevent="updateShift(shift.id)"
            >
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Zone tag</label>
                <select
                  v-model="ensureShiftDraft(shift.id).zoneTagId"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                >
                  <option value="">
                    Без привязки
                  </option>
                  <option
                    v-for="tag in shiftAvailableZoneTags"
                    :key="tag.id"
                    :value="String(tag.id)"
                  >
                    {{ tag.name }}
                  </option>
                </select>
              </div>

              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Начало</label>
                  <input
                    v-model="ensureShiftDraft(shift.id).startTime"
                    type="datetime-local"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                  >
                </div>
                <div>
                  <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Конец</label>
                  <input
                    v-model="ensureShiftDraft(shift.id).endTime"
                    type="datetime-local"
                    class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                  >
                </div>
              </div>

              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Заметка</label>
                <textarea
                  v-model="ensureShiftDraft(shift.id).note"
                  rows="4"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                />
              </div>

              <button
                type="submit"
                class="rounded-[0.8rem] bg-cyan-300 px-5 py-3 text-sm font-black text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
                :disabled="isMutating"
              >
                Сохранить смену
              </button>
            </form>
          </div>
        </article>
      </section>
    </template>

    <!-- Modal for creation -->
    <div
      v-if="isCreateModalOpen"
      class="fixed inset-0 z-[60] flex items-center justify-center bg-[#020c13]/78 px-3 py-6 backdrop-blur-sm"
      @click.self="closeCreateModal"
    >
      <div class="w-full max-w-2xl overflow-hidden border border-white/10 bg-[#07141d] shadow-[0_32px_80px_rgba(0,0,0,0.55)]">
        <div class="flex items-start justify-between gap-4 border-b border-white/8 bg-[#08131b] px-5 py-4">
          <div>
            <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
              {{ createModalMeta.eyebrow }}
            </p>
            <h2 class="mt-2 text-2xl font-black text-white">
              {{ createModalMeta.title }}
            </h2>
            <p class="mt-2 text-sm text-zinc-400">
              {{ createModalMeta.description }}
            </p>
          </div>

          <button
            type="button"
            class="rounded-[0.75rem] border border-white/10 px-3 py-2 text-xs font-bold uppercase tracking-[0.24em] text-zinc-300 transition hover:border-white/25 hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="isMutating"
            @click="closeCreateModal"
          >
            Закрыть
          </button>
        </div>

        <div class="max-h-[calc(100vh-8rem)] overflow-y-auto px-5 py-5">
          <form
            v-if="activeTab === 'zones'"
            class="space-y-4"
            @submit.prevent="submitZone"
          >
            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Название</label>
                <input
                  v-model="newZoneForm.name"
                  type="text"
                  placeholder="Например, Lounge Neon"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
                >
              </div>
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Тип зоны</label>
                <select
                  v-model="newZoneForm.type"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                >
                  <option
                    v-for="option in zoneTypeOptions"
                    :key="option.value"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>

            <div class="grid gap-4 lg:grid-cols-[1fr_auto]">
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Zone tag</label>
                <select
                  v-model="newZoneForm.zoneTagId"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                >
                  <option
                    disabled
                    value=""
                  >
                    Выбери tag
                  </option>
                  <option
                    v-for="tag in sortedZoneTags"
                    :key="tag.id"
                    :value="String(tag.id)"
                  >
                    {{ tag.name }}
                  </option>
                </select>
              </div>

              <div class="border border-white/8 bg-white/4 p-3 lg:min-w-64">
                <p class="text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">
                  Быстрый tag
                </p>
                <label class="mb-2 mt-3 block text-[11px] font-bold uppercase tracking-[0.24em] text-cyan-100/45">Название тега</label>
                <div class="flex gap-2">
                  <input
                    v-model="zoneCreateTagName"
                    type="text"
                    placeholder="Например, lounge-vip"
                    class="min-w-0 flex-1 rounded-[0.75rem] border border-cyan-400/18 bg-[#06131c] px-3 py-2.5 text-sm text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
                  >
                  <button
                    type="button"
                    class="rounded-[0.75rem] bg-cyan-300 px-4 py-2.5 text-sm font-black text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
                    :disabled="isMutating"
                    @click="createZoneTagFromZoneContext(zoneCreateTagName, 'create')"
                  >
                    Добавить
                  </button>
                </div>
              </div>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Вместимость</label>
                <input
                  v-model="newZoneForm.capacity"
                  type="number"
                  min="1"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                >
              </div>
              <label class="flex items-center gap-3 rounded-[0.9rem] border border-white/8 bg-white/4 px-4 py-3 text-sm text-zinc-200">
                <input
                  v-model="newZoneForm.isActive"
                  type="checkbox"
                  class="h-4 w-4 accent-cyan-300"
                >
                Зона активна и видна в рабочем состоянии
              </label>
            </div>

            <div>
              <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Описание</label>
              <textarea
                v-model="newZoneForm.description"
                rows="4"
                placeholder="Коротко опиши сценарий, атмосферу или ограничения зоны"
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
              />
            </div>

            <div>
              <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Details JSON</label>
              <textarea
                v-model="newZoneForm.detailsJson"
                rows="6"
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 font-mono text-sm text-white focus:border-cyan-300 focus:outline-none"
              />
            </div>
            <button
              type="submit"
              class="w-full rounded-[0.85rem] bg-cyan-300 px-5 py-3.5 text-sm font-black uppercase tracking-[0.28em] text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
              :disabled="isMutating"
            >
              Создать зону
            </button>
          </form>

          <form
            v-else-if="activeTab === 'zone-tags'"
            class="space-y-4"
            @submit.prevent="submitZoneTag"
          >
            <div>
              <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Название тега</label>
              <input
                v-model="newZoneTagName"
                type="text"
                placeholder="Например, gaming-premium"
                class="w-full rounded-[0.8rem] border border-orange-300/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-orange-100/20 focus:border-orange-300 focus:outline-none"
              >
            </div>

            <button
              type="submit"
              class="w-full rounded-[0.85rem] bg-orange-300 px-5 py-3.5 text-sm font-black uppercase tracking-[0.28em] text-[#020c13] transition hover:bg-orange-200 disabled:cursor-not-allowed disabled:bg-orange-300/50"
              :disabled="isMutating"
            >
              Создать тег
            </button>
          </form>

          <form
            v-else
            class="space-y-4"
            @submit.prevent="submitShift"
          >
            <div>
              <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Zone tag</label>
              <select
                v-model="newShiftForm.zoneTagId"
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
              >
                <option value="">
                  Без привязки
                </option>
                <option
                  v-for="tag in shiftAvailableZoneTags"
                  :key="tag.id"
                  :value="String(tag.id)"
                >
                  {{ tag.name }}
                </option>
              </select>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Начало</label>
                <input
                  v-model="newShiftForm.startTime"
                  type="datetime-local"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                >
              </div>
              <div>
                <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Конец</label>
                <input
                  v-model="newShiftForm.endTime"
                  type="datetime-local"
                  class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white focus:border-cyan-300 focus:outline-none"
                >
              </div>
            </div>

            <div>
              <label class="mb-3 block text-xs font-bold uppercase tracking-[0.28em] text-cyan-100/45">Заметка</label>
              <textarea
                v-model="newShiftForm.note"
                rows="4"
                placeholder="Например, вечерняя смена для front desk и lounge"
                class="w-full rounded-[0.8rem] border border-cyan-400/18 bg-[#06131c] px-4 py-3 text-white placeholder:text-cyan-100/20 focus:border-cyan-300 focus:outline-none"
              />
            </div>

            <button
              type="submit"
              class="w-full rounded-[0.85rem] bg-cyan-300 px-5 py-3.5 text-sm font-black uppercase tracking-[0.28em] text-[#020c13] transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-300/50"
              :disabled="isMutating"
            >
              Создать смену
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import {
  createAdminPlace,
  getAdminBookings,
  createAdminShift,
  createAdminZone,
  createAdminZoneTag,
  deleteAdminPlace,
  deleteAdminShift,
  getAdminServices,
  deleteAdminZone,
  deleteAdminZoneTag,
  getAdminPlaces,
  getAdminShifts,
  getAdminZones,
  getAdminZoneTags,
  getAdminSiteSettings,
  patchAdminBooking,
  patchAdminPlace,
  patchAdminShift,
  patchAdminZone,
  patchAdminZoneTag
} from '~/api/admin'
import type { AdminPlace, AdminService, AdminZone, AdminZoneTag, BookingRecord, BookingStatus, ShiftSchedule } from '~/api/types'
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'admin'
})

useHead({
  title: 'Admin - PlayGround'
})

const authStore = useAuthStore()

type AdminTab = 'zones' | 'zone-tags' | 'bookings' | 'shifts'

interface ZoneTagView {
  id: number
  name: string
  createdAt: string
  updatedAt: string
}

interface ZoneView {
  id: number
  name: string
  zoneType: string
  zoneTagId: number
  capacity: number
  description: string
  detailsJson: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

interface PlaceView {
  id: number
  zoneId: number
  label: string
  configurationId?: number
  sortOrder: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

interface ServiceView {
  id: number
  name: string
  zoneId: number
  duration: number
  price: string
  currency: string
  description: string
  isActive: boolean
}

interface BookingView {
  id: number
  userId: string
  zoneId: number
  serviceId: number
  placeId?: number
  startTime: string
  endTime: string
  participants: number
  totalPrice: string
  status: BookingStatus
  contactName: string
  contactEmail: string
  contactPhone: string
  detailsJson: string
  createdAt: string
  updatedAt: string
}

interface ZoneForm {
  name: string
  type: string
  zoneTagId: string
  capacity: string
  description: string
  detailsJson: string
  isActive: boolean
}

interface PlaceForm {
  label: string
  configurationId: string
  sortOrder: string
  isActive: boolean
}

interface ShiftForm {
  zoneTagId: string
  startTime: string
  endTime: string
  note: string
}

const tabs = [
  {
    id: 'zones' as const,
    eyebrow: 'Layout + places',
    label: 'Zones',
    title: 'Zone Management',
    description: 'Карточки зон и вложенные места прямо внутри каждой зоны.'
  },
  {
    id: 'zone-tags' as const,
    eyebrow: 'Reference data',
    label: 'Zone tags',
    title: 'Zone Tag Directory',
    description: 'Отдельный справочник тегов для связки смен и зон.'
  },
  {
    id: 'bookings' as const,
    eyebrow: 'Reservations',
    label: 'Bookings',
    title: 'Booking Registry',
    description: 'Все клиентские брони со статусами и контактной информацией.'
  },
  {
    id: 'shifts' as const,
    eyebrow: 'Operations',
    label: 'Shifts',
    title: 'Shift Control',
    description: 'Планирование смен текущего admin-пользователя.'
  }
]

const zoneTypeOptions = [
  { value: 'game', label: 'Gaming' },
  { value: 'lounge', label: 'Lounge' },
  { value: 'event', label: 'Event' },
  { value: 'sys', label: 'System' }
]

const route = useRoute()
const activeTab = ref<AdminTab>((route.query.tab as AdminTab) || 'zones')

watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab && ['zones', 'zone-tags', 'bookings', 'shifts'].includes(newTab as string)) {
      activeTab.value = newTab as AdminTab
    }
  }
)
const isLoading = ref(true)
const isMutating = ref(false)
const pageError = ref('')
const feedbackTone = ref<'success' | 'error'>('success')
const feedbackMessage = ref('')
const isCreateModalOpen = ref(false)

const zoneTags = ref<ZoneTagView[]>([])
const activeSettingsTagIds = ref<number[]>([])
const zones = ref<ZoneView[]>([])
const places = ref<PlaceView[]>([])
const services = ref<ServiceView[]>([])
const bookings = ref<BookingView[]>([])
const shifts = ref<ShiftSchedule[]>([])

const newZoneTagName = ref('')
const zoneCreateTagName = ref('')
const zoneEditTagNames = ref<Record<number, string>>({})
const newZoneForm = ref<ZoneForm>(emptyZoneForm())
const zoneDrafts = ref<Record<number, ZoneForm>>({})
const newPlaceForms = ref<Record<number, PlaceForm>>({})
const placeDrafts = ref<Record<number, PlaceForm>>({})
const newShiftForm = ref<ShiftForm>(emptyShiftForm())
const shiftDrafts = ref<Record<number, ShiftForm>>({})
const zoneTagDrafts = ref<Record<number, string>>({})
const bookingStatusDrafts = ref<Record<number, BookingStatus>>({})

const bookingStatusOptions: BookingStatus[] = ['created', 'confirmed', 'canceled', 'completed']

const expandedZoneId = ref<number | null>(null)
const expandedPlaceId = ref<number | null>(null)
const expandedShiftId = ref<number | null>(null)
const expandedZoneTagId = ref<number | null>(null)
const expandedBookingDetails = ref<number | null>(null)
const expandedShiftDetails = ref<number | null>(null)

const dashboardStats = computed(() => [
  {
    label: 'Zones',
    value: zones.value.length,
    shellClass: 'border-cyan-400/18 bg-cyan-400/10 shadow-[inset_4px_0_0_rgba(34,211,238,0.8)]',
    headerClass: 'border-cyan-400/18 bg-cyan-400/10'
  },
  {
    label: 'Places',
    value: places.value.length,
    shellClass: 'border-orange-300/18 bg-orange-400/10 shadow-[inset_4px_0_0_rgba(253,186,116,0.75)]',
    headerClass: 'border-orange-300/18 bg-orange-400/10'
  },
  {
    label: 'Shifts',
    value: shifts.value.length,
    shellClass: 'border-white/10 bg-white/[0.05] shadow-[inset_4px_0_0_rgba(255,255,255,0.35)]',
    headerClass: 'border-white/10 bg-white/[0.04]'
  },
  {
    label: 'Bookings',
    value: bookings.value.length,
    shellClass: 'border-emerald-300/18 bg-emerald-400/10 shadow-[inset_4px_0_0_rgba(52,211,153,0.6)]',
    headerClass: 'border-emerald-300/18 bg-emerald-400/10'
  }
])

const activeTabItem = computed(() => tabs.find(tab => tab.id === activeTab.value) ?? tabs[0]!)
const adminDisplayName = computed(() => authStore.user?.name || 'Admin user')
const adminEmail = computed(() => authStore.user?.email || 'session@playground.local')
const showCreateAction = computed(() => activeTab.value !== 'bookings')
const createModalMeta = computed(() => {
  switch (activeTab.value) {
    case 'zones':
      return {
        eyebrow: 'Create zone',
        title: 'Новая зона',
        description: 'Создай новую зону во всплывающем окне, не сжимая карточки текущего списка.',
        buttonLabel: 'Добавить зону'
      }
    case 'zone-tags':
      return {
        eyebrow: 'Create tag',
        title: 'Новый тег зоны',
        description: 'Теги остаются отдельным справочником и сразу доступны в формах зон и смен.',
        buttonLabel: 'Добавить тег'
      }
    case 'shifts':
      return {
        eyebrow: 'Create shift',
        title: 'Новая смена',
        description: 'Планируй смены в модальном окне, не теряя контекст текущего списка.',
        buttonLabel: 'Добавить смену'
      }
    case 'bookings':
      return {
        eyebrow: 'Booking log',
        title: 'Журнал бронирований',
        description: 'Все записи из gaming появляются здесь, вместе со статусами и контактными данными.',
        buttonLabel: ''
      }
  }

  return {
    eyebrow: 'Create entity',
    title: 'Новая сущность',
    description: 'Создание новой записи в отдельном модальном окне.',
    buttonLabel: 'Добавить'
  }
})

const sortedZoneTags = computed(() => [...zoneTags.value].sort((left, right) => left.name.localeCompare(right.name)))

const shiftAvailableZoneTags = computed(() => {
  if (activeSettingsTagIds.value.length === 0) {
    return sortedZoneTags.value
  }
  return sortedZoneTags.value.filter(tag => activeSettingsTagIds.value.includes(tag.id))
})

const sortedZones = computed(() => [...zones.value].sort((left, right) => left.name.localeCompare(right.name)))
const sortedBookings = computed(() => [...bookings.value].sort((left, right) => new Date(right.createdAt).getTime() - new Date(left.createdAt).getTime()))
const sortedShifts = computed(() => [...shifts.value].sort((left, right) => new Date(right.start_time).getTime() - new Date(left.start_time).getTime()))

let feedbackTimer: ReturnType<typeof setTimeout> | null = null

function emptyZoneForm(): ZoneForm {
  return {
    name: '',
    type: 'game',
    zoneTagId: '',
    capacity: '1',
    description: '',
    detailsJson: '{}',
    isActive: true
  }
}

function emptyPlaceForm(): PlaceForm {
  return {
    label: '',
    configurationId: '',
    sortOrder: '0',
    isActive: true
  }
}

function emptyShiftForm(): ShiftForm {
  return {
    zoneTagId: '',
    startTime: '',
    endTime: '',
    note: ''
  }
}

function resetFeedbackTimer() {
  if (feedbackTimer) {
    clearTimeout(feedbackTimer)
    feedbackTimer = null
  }
}

function showFeedback(tone: 'success' | 'error', message: string) {
  resetFeedbackTimer()
  feedbackTone.value = tone
  feedbackMessage.value = message

  if (tone === 'success') {
    feedbackTimer = setTimeout(() => {
      feedbackMessage.value = ''
      feedbackTimer = null
    }, 3000)
  }
}

function clearFeedback() {
  resetFeedbackTimer()
  feedbackMessage.value = ''
}

function openCreateModal() {
  isCreateModalOpen.value = true
}

function closeCreateModal() {
  if (isMutating.value) {
    return
  }

  isCreateModalOpen.value = false
}

function extractErrorMessage(error: unknown) {
  if (error && typeof error === 'object') {
    const maybeError = error as {
      data?: { message?: string }
      message?: string
    }

    return maybeError.data?.message || maybeError.message || 'Произошла ошибка'
  }

  return 'Произошла ошибка'
}

function readNullableText(value: unknown) {
  if (typeof value === 'string') {
    return value
  }

  if (value && typeof value === 'object') {
    const candidate = value as { String?: string, Valid?: boolean }
    if (candidate.Valid && typeof candidate.String === 'string') {
      return candidate.String
    }
  }

  return ''
}

function decodeMaybeBase64(value: string) {
  try {
    if (import.meta.server) {
      return Buffer.from(value, 'base64').toString('utf-8')
    }

    return window.atob(value)
  } catch {
    return value
  }
}

function normalizeDetailsJson(value: unknown) {
  if (!value) {
    return '{}'
  }

  if (typeof value === 'string') {
    const direct = value.trim()
    if (direct.startsWith('{') || direct.startsWith('[')) {
      return formatJsonString(direct)
    }

    const decoded = decodeMaybeBase64(direct).trim()
    if (decoded.startsWith('{') || decoded.startsWith('[')) {
      return formatJsonString(decoded)
    }

    return direct
  }

  if (Array.isArray(value)) {
    return formatJsonString(JSON.stringify(value))
  }

  return formatJsonString(JSON.stringify(value))
}

function formatJsonString(value: string) {
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}

function normalizeOptionalNumber(value: unknown) {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value
  }

  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : undefined
  }

  if (value && typeof value === 'object') {
    const candidate = value as { Int64?: number, Int32?: number, Valid?: boolean }
    if (candidate.Valid) {
      if (typeof candidate.Int64 === 'number') {
        return candidate.Int64
      }
      if (typeof candidate.Int32 === 'number') {
        return candidate.Int32
      }
    }
  }

  return undefined
}

function normalizeZoneTag(raw: AdminZoneTag): ZoneTagView {
  return {
    id: raw.id,
    name: raw.name,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at
  }
}

function normalizeZone(raw: AdminZone): ZoneView {
  return {
    id: raw.id,
    name: raw.name,
    zoneType: raw.zone_type,
    zoneTagId: raw.zone_tag_id,
    capacity: raw.capacity,
    description: readNullableText(raw.description),
    detailsJson: normalizeDetailsJson(raw.details_json),
    isActive: raw.is_active,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at
  }
}

function normalizePlace(raw: AdminPlace): PlaceView {
  return {
    id: raw.id,
    zoneId: raw.zone_id,
    label: raw.label,
    configurationId: normalizeOptionalNumber(raw.configuration_id),
    sortOrder: raw.sort_order,
    isActive: raw.is_active,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at
  }
}

function normalizeService(raw: AdminService): ServiceView {
  return {
    id: raw.id,
    name: raw.name,
    zoneId: raw.zone_id,
    duration: raw.duration,
    price: raw.price,
    currency: raw.currency,
    description: readNullableText(raw.description),
    isActive: raw.is_active
  }
}

function normalizeBooking(raw: BookingRecord): BookingView {
  return {
    id: raw.id,
    userId: raw.user_id,
    zoneId: raw.zone_id,
    serviceId: raw.service_id,
    placeId: normalizeOptionalNumber(raw.place_id),
    startTime: raw.start_time,
    endTime: raw.end_time,
    participants: raw.participants,
    totalPrice: raw.total_price,
    status: raw.status,
    contactName: raw.contact_name,
    contactEmail: raw.contact_email,
    contactPhone: raw.contact_phone,
    detailsJson: normalizeDetailsJson(raw.details_json),
    createdAt: raw.created_at,
    updatedAt: raw.updated_at
  }
}

function makeZoneForm(zone?: ZoneView): ZoneForm {
  if (!zone) {
    return emptyZoneForm()
  }

  return {
    name: zone.name,
    type: zone.zoneType,
    zoneTagId: String(zone.zoneTagId),
    capacity: String(zone.capacity),
    description: zone.description,
    detailsJson: zone.detailsJson,
    isActive: zone.isActive
  }
}

function makePlaceForm(place?: PlaceView): PlaceForm {
  if (!place) {
    return emptyPlaceForm()
  }

  return {
    label: place.label,
    configurationId: place.configurationId ? String(place.configurationId) : '',
    sortOrder: String(place.sortOrder),
    isActive: place.isActive
  }
}

function makeShiftForm(shift?: ShiftSchedule): ShiftForm {
  if (!shift) {
    return emptyShiftForm()
  }

  return {
    zoneTagId: shift.zone_tag_id ? String(shift.zone_tag_id) : '',
    startTime: toLocalInputValue(shift.start_time),
    endTime: toLocalInputValue(shift.end_time),
    note: shift.note || ''
  }
}

function hydrateDrafts() {
  zoneDrafts.value = Object.fromEntries(zones.value.map(zone => [zone.id, makeZoneForm(zone)]))
  newPlaceForms.value = Object.fromEntries(zones.value.map(zone => [zone.id, emptyPlaceForm()]))
  placeDrafts.value = Object.fromEntries(places.value.map(place => [place.id, makePlaceForm(place)]))
  shiftDrafts.value = Object.fromEntries(shifts.value.map(shift => [shift.id, makeShiftForm(shift)]))
  zoneTagDrafts.value = Object.fromEntries(zoneTags.value.map(tag => [tag.id, tag.name]))
  zoneEditTagNames.value = Object.fromEntries(zones.value.map(zone => [zone.id, '']))
  bookingStatusDrafts.value = Object.fromEntries(bookings.value.map(booking => [booking.id, booking.status]))
}

function ensureZoneDraft(zoneId: number) {
  const currentDraft = zoneDrafts.value[zoneId]
  if (currentDraft) {
    return currentDraft
  }

  const createdDraft = emptyZoneForm()
  zoneDrafts.value[zoneId] = createdDraft
  return createdDraft
}

function ensureNewPlaceForm(zoneId: number) {
  const currentDraft = newPlaceForms.value[zoneId]
  if (currentDraft) {
    return currentDraft
  }

  const createdDraft = emptyPlaceForm()
  newPlaceForms.value[zoneId] = createdDraft
  return createdDraft
}

function ensurePlaceDraft(placeId: number) {
  const currentDraft = placeDrafts.value[placeId]
  if (currentDraft) {
    return currentDraft
  }

  const createdDraft = emptyPlaceForm()
  placeDrafts.value[placeId] = createdDraft
  return createdDraft
}

function ensureShiftDraft(shiftId: number) {
  const currentDraft = shiftDrafts.value[shiftId]
  if (currentDraft) {
    return currentDraft
  }

  const createdDraft = emptyShiftForm()
  shiftDrafts.value[shiftId] = createdDraft
  return createdDraft
}

function ensureZoneEditTagNames(zoneId: number) {
  if (zoneEditTagNames.value[zoneId] === undefined) {
    zoneEditTagNames.value[zoneId] = ''
  }

  return zoneEditTagNames.value
}

async function loadAdminData(showLoader = true) {
  clearFeedback()
  pageError.value = ''

  if (showLoader) {
    isLoading.value = true
  }

  try {
    const [zoneTagsResponse, zonesResponse, placesResponse, servicesResponse, bookingsResponse, shiftsResponse, settingsResponse] = await Promise.all([
      getAdminZoneTags(),
      getAdminZones(),
      getAdminPlaces(),
      getAdminServices(),
      getAdminBookings(),
      getAdminShifts(),
      getAdminSiteSettings()
    ])

    zoneTags.value = (zoneTagsResponse.zone_tags ?? []).map(normalizeZoneTag)
    zones.value = (zonesResponse.zones ?? []).map(normalizeZone)
    places.value = (placesResponse.places ?? []).map(normalizePlace)
    services.value = (servicesResponse.services ?? []).map(normalizeService)
    bookings.value = (bookingsResponse.bookings ?? []).map(normalizeBooking)
    shifts.value = shiftsResponse.shifts ?? []

    const loadedIds = settingsResponse.settings?.settings_json
    activeSettingsTagIds.value = Array.isArray(loadedIds) ? loadedIds : []

    hydrateDrafts()
  } catch (error) {
    pageError.value = extractErrorMessage(error)
  } finally {
    if (showLoader) {
      isLoading.value = false
    }
  }
}

function placesForZone(zoneId: number) {
  return places.value
    .filter(place => place.zoneId === zoneId)
    .sort((left, right) => left.sortOrder - right.sortOrder || left.label.localeCompare(right.label))
}

function zoneTagName(zoneTagId?: number) {
  if (!zoneTagId) {
    return 'Без тега'
  }

  return zoneTags.value.find(tag => tag.id === zoneTagId)?.name || `Tag #${zoneTagId}`
}

function zoneTypeLabel(zoneType: string) {
  return zoneTypeOptions.find(option => option.value === zoneType)?.label || zoneType
}

function tabMetric(tabId: AdminTab) {
  switch (tabId) {
    case 'zones':
      return zones.value.length
    case 'zone-tags':
      return zoneTags.value.length
    case 'bookings':
      return bookings.value.length
    case 'shifts':
      return shifts.value.length
  }
}

function zoneNameById(zoneId: number) {
  return zones.value.find(zone => zone.id === zoneId)?.name || `Zone #${zoneId}`
}

function placeLabelById(placeId?: number) {
  if (!placeId) {
    return 'Без места'
  }

  return places.value.find(place => place.id === placeId)?.label || `Place #${placeId}`
}

function serviceNameById(serviceId: number) {
  return services.value.find(service => service.id === serviceId)?.name || `Service #${serviceId}`
}

function bookingStatusClass(status: BookingStatus) {
  switch (status) {
    case 'confirmed':
      return 'border-emerald-300/20 bg-emerald-500/12 text-emerald-100/80'
    case 'canceled':
      return 'border-orange-300/20 bg-orange-500/12 text-orange-100/80'
    case 'completed':
      return 'border-cyan-300/20 bg-cyan-400/10 text-cyan-100/80'
    default:
      return 'border-white/10 bg-white/4 text-zinc-200'
  }
}

function formatAdminDateTime(value: string) {
  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(new Date(value))
}

function zoneBadgeClass(zoneType: string) {
  switch (zoneType) {
    case 'game':
      return 'border-cyan-300/20 bg-cyan-400/10 text-cyan-100/80'
    case 'lounge':
      return 'border-orange-300/20 bg-orange-400/10 text-orange-100/80'
    case 'event':
      return 'border-amber-300/20 bg-amber-400/10 text-amber-100/80'
    default:
      return 'border-zinc-300/20 bg-zinc-400/10 text-zinc-100/80'
  }
}

function hasZoneDetails(detailsJson: string) {
  return detailsJson.trim() !== '{}' && detailsJson.trim() !== ''
}

function toLocalInputValue(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }

  const pad = (part: number) => String(part).padStart(2, '0')

  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

function toIsoValue(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toISOString()
}

function formatDateRange(startAt: string, endAt: string) {
  const formatter = new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short'
  })

  return `${formatter.format(new Date(startAt))} - ${formatter.format(new Date(endAt))}`
}

function toPositiveNumber(value: string) {
  if (value.trim() === '') {
    return undefined
  }

  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : undefined
}

function toggleZoneEditor(zoneId: number) {
  expandedZoneId.value = expandedZoneId.value === zoneId ? null : zoneId
}

function togglePlaceEditor(placeId: number) {
  expandedPlaceId.value = expandedPlaceId.value === placeId ? null : placeId
}

function toggleShiftEditor(shiftId: number) {
  expandedShiftId.value = expandedShiftId.value === shiftId ? null : shiftId
}

function toggleZoneTagEditor(zoneTagId: number) {
  expandedZoneTagId.value = expandedZoneTagId.value === zoneTagId ? null : zoneTagId
}

function toggleBookingDetails(bookingId: number) {
  expandedBookingDetails.value = expandedBookingDetails.value === bookingId ? null : bookingId
}

function toggleShiftDetails(shiftId: number) {
  expandedShiftDetails.value = expandedShiftDetails.value === shiftId ? null : shiftId
}

function resetZoneDraft(zoneId: number) {
  const zone = zones.value.find(item => item.id === zoneId)
  if (!zone) {
    return
  }

  zoneDrafts.value[zoneId] = makeZoneForm(zone)
  zoneEditTagNames.value[zoneId] = ''
}

async function createZoneTagFromZoneContext(name: string, target: 'create' | number) {
  const trimmedName = name.trim()
  if (!trimmedName) {
    showFeedback('error', 'Для нового zone tag нужно указать название.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    const response = await createAdminZoneTag({ name: trimmedName })
    await loadAdminData(false)

    const createdId = response.zone_tag.id

    if (target === 'create') {
      newZoneForm.value.zoneTagId = String(createdId)
      zoneCreateTagName.value = ''
    } else {
      const draft = zoneDrafts.value[target]
      if (draft) {
        draft.zoneTagId = String(createdId)
      }
      zoneEditTagNames.value[target] = ''
    }

    showFeedback('success', `Тег ${trimmedName} создан и сразу выбран в форме зоны.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function submitZoneTag() {
  const trimmedName = newZoneTagName.value.trim()
  if (!trimmedName) {
    showFeedback('error', 'Название zone tag обязательно.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await createAdminZoneTag({ name: trimmedName })
    newZoneTagName.value = ''
    await loadAdminData(false)
    isCreateModalOpen.value = false
    showFeedback('success', `Тег ${trimmedName} создан.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function updateZoneTag(zoneTagId: number) {
  const trimmedName = zoneTagDrafts.value[zoneTagId]?.trim()
  if (!trimmedName) {
    showFeedback('error', 'Название zone tag не может быть пустым.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await patchAdminZoneTag(zoneTagId, { name: trimmedName })
    await loadAdminData(false)
    expandedZoneTagId.value = null
    showFeedback('success', 'Тег зоны обновлён.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function removeZoneTag(zoneTagId: number, zoneTagNameValue: string) {
  if (import.meta.client && !window.confirm(`Удалить zone tag ${zoneTagNameValue}?`)) {
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await deleteAdminZoneTag(zoneTagId)
    await loadAdminData(false)
    showFeedback('success', `Тег ${zoneTagNameValue} удалён.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function submitZone() {
  const zoneTagId = toPositiveNumber(newZoneForm.value.zoneTagId)
  const capacity = toPositiveNumber(newZoneForm.value.capacity)

  if (!newZoneForm.value.name.trim()) {
    showFeedback('error', 'Для зоны нужно указать название.')
    return
  }

  if (!zoneTagId) {
    showFeedback('error', 'Для зоны нужно выбрать zone tag.')
    return
  }

  if (!capacity || capacity < 1) {
    showFeedback('error', 'Вместимость зоны должна быть больше нуля.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await createAdminZone({
      name: newZoneForm.value.name.trim(),
      type: newZoneForm.value.type,
      zone_tag_id: zoneTagId,
      capacity,
      description: newZoneForm.value.description.trim(),
      is_active: newZoneForm.value.isActive,
      details_json: newZoneForm.value.detailsJson.trim() || '{}'
    })

    newZoneForm.value = emptyZoneForm()
    zoneCreateTagName.value = ''
    await loadAdminData(false)
    isCreateModalOpen.value = false
    showFeedback('success', 'Зона создана.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function updateZone(zoneId: number) {
  const draft = zoneDrafts.value[zoneId]
  if (!draft) {
    return
  }

  const zoneTagId = toPositiveNumber(draft.zoneTagId)
  const capacity = toPositiveNumber(draft.capacity)

  if (!draft.name.trim()) {
    showFeedback('error', 'Название зоны не может быть пустым.')
    return
  }

  if (!zoneTagId) {
    showFeedback('error', 'Для зоны нужно выбрать zone tag.')
    return
  }

  if (!capacity || capacity < 1) {
    showFeedback('error', 'Вместимость зоны должна быть больше нуля.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await patchAdminZone(zoneId, {
      name: draft.name.trim(),
      type: draft.type,
      zone_tag_id: zoneTagId,
      capacity,
      description: draft.description.trim(),
      is_active: draft.isActive,
      details_json: draft.detailsJson.trim() || '{}'
    })

    await loadAdminData(false)
    showFeedback('success', 'Зона обновлена.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function removeZone(zoneId: number, zoneNameValue: string) {
  if (import.meta.client && !window.confirm(`Удалить зону ${zoneNameValue}? Сначала убедись, что связанные места уже не нужны.`)) {
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await deleteAdminZone(zoneId)
    await loadAdminData(false)
    showFeedback('success', `Зона ${zoneNameValue} удалена.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function submitPlace(zoneId: number) {
  const draft = newPlaceForms.value[zoneId]
  if (!draft) {
    return
  }

  const sortOrder = toPositiveNumber(draft.sortOrder)

  if (!draft.label.trim()) {
    showFeedback('error', 'Для места нужен label.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await createAdminPlace({
      zone_id: zoneId,
      label: draft.label.trim(),
      configuration_id: toPositiveNumber(draft.configurationId),
      sort_order: sortOrder,
      is_active: draft.isActive
    })

    newPlaceForms.value[zoneId] = emptyPlaceForm()
    await loadAdminData(false)
    showFeedback('success', 'Место создано внутри зоны.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function updatePlace(placeId: number) {
  const draft = placeDrafts.value[placeId]
  if (!draft) {
    return
  }

  if (!draft.label.trim()) {
    showFeedback('error', 'Label места не может быть пустым.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await patchAdminPlace(placeId, {
      label: draft.label.trim(),
      configuration_id: toPositiveNumber(draft.configurationId),
      sort_order: toPositiveNumber(draft.sortOrder),
      is_active: draft.isActive
    })

    await loadAdminData(false)
    showFeedback('success', 'Место обновлено.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function removePlace(placeId: number, placeLabel: string) {
  if (import.meta.client && !window.confirm(`Удалить место ${placeLabel}?`)) {
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await deleteAdminPlace(placeId)
    await loadAdminData(false)
    showFeedback('success', `Место ${placeLabel} удалено.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function updateBookingStatus(bookingId: number) {
  const nextStatus = bookingStatusDrafts.value[bookingId]
  const booking = bookings.value.find(item => item.id === bookingId)

  if (!booking || !nextStatus) {
    return
  }

  if (booking.status === nextStatus) {
    showFeedback('success', 'Статус уже актуален.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await patchAdminBooking(bookingId, { status: nextStatus })
    await loadAdminData(false)
    showFeedback('success', `Статус брони #${bookingId} обновлён.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function submitShift() {
  if (!newShiftForm.value.startTime || !newShiftForm.value.endTime) {
    showFeedback('error', 'Для смены нужно указать начало и конец.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await createAdminShift({
      start_time: toIsoValue(newShiftForm.value.startTime),
      end_time: toIsoValue(newShiftForm.value.endTime),
      zone_tag_id: toPositiveNumber(newShiftForm.value.zoneTagId),
      note: newShiftForm.value.note.trim() || undefined
    })

    newShiftForm.value = emptyShiftForm()
    await loadAdminData(false)
    isCreateModalOpen.value = false
    showFeedback('success', 'Смена создана.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function updateShift(shiftId: number) {
  const draft = shiftDrafts.value[shiftId]
  if (!draft) {
    return
  }

  if (!draft.startTime || !draft.endTime) {
    showFeedback('error', 'Для смены нужно указать начало и конец.')
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await patchAdminShift(shiftId, {
      start_time: toIsoValue(draft.startTime),
      end_time: toIsoValue(draft.endTime),
      zone_tag_id: toPositiveNumber(draft.zoneTagId),
      note: draft.note.trim() || undefined
    })

    await loadAdminData(false)
    showFeedback('success', 'Смена обновлена.')
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

async function removeShift(shiftId: number) {
  if (import.meta.client && !window.confirm(`Удалить смену #${shiftId}?`)) {
    return
  }

  isMutating.value = true
  clearFeedback()

  try {
    await deleteAdminShift(shiftId)
    await loadAdminData(false)
    showFeedback('success', `Смена #${shiftId} удалена.`)
  } catch (error) {
    showFeedback('error', extractErrorMessage(error))
  } finally {
    isMutating.value = false
  }
}

onBeforeUnmount(() => {
  resetFeedbackTimer()
})

onMounted(async () => {
  await loadAdminData()
})
</script>
<style scoped>
.admin-page input:not([type='checkbox']),
.admin-page select,
.admin-page textarea {
  min-height: 2.45rem;
  border-radius: 0.35rem !important;
  border-color: rgba(103, 232, 249, 0.16) !important;
  background: linear-gradient(180deg, rgba(10, 24, 34, 0.98), rgba(8, 20, 28, 0.98)) !important;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.03);
}

.admin-page textarea {
  min-height: 6.5rem;
}

.admin-page select {
  appearance: none;
  color-scheme: dark;
  background-image:
    linear-gradient(45deg, transparent 50%, rgba(103, 232, 249, 0.82) 50%),
    linear-gradient(135deg, rgba(103, 232, 249, 0.82) 50%, transparent 50%);
  background-position:
    calc(100% - 16px) calc(50% - 2px),
    calc(100% - 10px) calc(50% - 2px);
  background-repeat: no-repeat;
  background-size: 6px 6px, 6px 6px;
  padding-right: 2.25rem !important;
}

.admin-page option,
.admin-page optgroup {
  background: #08131b;
  color: #f8fafc;
}

.admin-page button {
  border-radius: 0.35rem !important;
}

.admin-page article,
.admin-page section {
  border-radius: 0.45rem !important;
}

.admin-page pre {
  border-radius: 0.35rem;
}
</style>

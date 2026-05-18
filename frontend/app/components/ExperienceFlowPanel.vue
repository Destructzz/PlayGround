<template>
  <section class="rounded-xl border border-white/10 bg-[#071724]/75 backdrop-blur-md p-6 shadow-[0_18px_60px_rgba(0,0,0,0.4)] sm:p-8">
    <div class="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-[10px] font-bold uppercase tracking-[0.3em]" :style="{ color: accent || '#22d3ee' }">
          {{ eyebrow }}
        </p>
        <h2 class="mt-2 text-2xl font-black tracking-tight text-white uppercase">
          {{ title }}
        </h2>
      </div>
      <p class="max-w-xl text-xs leading-relaxed text-zinc-200">
        {{ description }}
      </p>
    </div>

    <slot />
  </section>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    eyebrow: string
    title: string
    description: string
    accent?: string
  }>(),
  {
    accent: '#22d3ee'
  }
)
</script>



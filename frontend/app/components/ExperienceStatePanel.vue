<template>
  <div class="rounded-xl border border-white/5 bg-[#030c14]/80 backdrop-blur-md p-8 text-center shadow-[0_16px_50px_rgba(0,0,0,0.55)] sm:p-10">
    <div
      class="mx-auto mb-5 flex h-12 w-12 items-center justify-center rounded-lg border text-lg"
      :style="iconStyle"
    >
      <slot name="icon">
        <span>✦</span>
      </slot>
    </div>
    <h2 class="text-xl font-black uppercase tracking-tight text-white sm:text-2xl">
      {{ title }}
    </h2>
    <p class="mx-auto mt-2 max-w-xl text-xs leading-relaxed text-zinc-300">
      {{ description }}
    </p>
    <div
      v-if="hint"
      class="mt-3 text-[10px] font-bold uppercase tracking-[0.3em]"
      :style="{ color: accent || '#22d3ee' }"
    >
      {{ hint }}
    </div>
    <div
      v-if="$slots.actions"
      class="mt-6 flex flex-wrap items-center justify-center gap-3"
    >
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    title: string
    description: string
    hint?: string
    accent?: string
  }>(),
  {
    accent: '#22d3ee'
  }
)

const iconStyle = computed(() => {
  const acc = props.accent || '#22d3ee'
  return {
    borderColor: `${acc}40`,
    backgroundColor: `${acc}10`,
    color: acc
  }
})
</script>


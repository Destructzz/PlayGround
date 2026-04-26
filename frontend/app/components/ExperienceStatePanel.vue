<template>
  <div class="rounded-[0.9rem] border border-cyan-400/20 bg-[#07141d] p-8 text-center shadow-[0_20px_80px_rgba(0,0,0,0.5)] sm:p-10">
    <div class="mx-auto mb-5 flex h-14 w-14 items-center justify-center rounded-[0.6rem] border border-cyan-300/35 bg-cyan-500/12 text-cyan-200">
      <slot name="icon">
        <span class="text-xl">✦</span>
      </slot>
    </div>
    <h2 class="text-2xl font-black uppercase tracking-[-0.04em] text-white sm:text-3xl">
      {{ title }}
    </h2>
    <p class="mx-auto mt-3 max-w-2xl text-sm leading-7 text-zinc-200 sm:text-base">
      {{ description }}
    </p>
    <div
      v-if="hint"
      class="mt-3 text-xs uppercase tracking-[0.35em] text-cyan-300/80"
    >
      {{ hint }}
    </div>
    <div
      v-if="$slots.actions"
      class="mt-7 flex flex-wrap items-center justify-center gap-3"
    >
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description: string
  hint?: string
}>()
</script>

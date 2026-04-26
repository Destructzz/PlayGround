<template>
  <section class="rounded-[0.9rem] border border-white/10 bg-[#07141d] p-6 shadow-[0_22px_84px_rgba(0,0,0,0.55)] sm:p-8">
    <div class="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-xs font-bold uppercase tracking-[0.35em] text-cyan-300/80">
          {{ eyebrow }}
        </p>
        <h2 class="mt-2 text-3xl font-black uppercase tracking-[-0.04em] text-white">
          {{ title }}
        </h2>
      </div>
      <p class="max-w-xl text-sm leading-7 text-zinc-200">
        {{ description }}
      </p>
    </div>

    <slot />
  </section>
</template>

<script setup lang="ts">
defineProps<{
  eyebrow: string
  title: string
  description: string
}>()
</script>

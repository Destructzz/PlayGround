<template>
  <section class="relative overflow-hidden rounded-[0.9rem] border border-cyan-400/25 bg-[#07141d] px-6 py-12 shadow-[0_24px_90px_rgba(0,0,0,0.55)] sm:px-10 sm:py-14">
    <div
      class="absolute inset-0 opacity-90"
      :style="heroSurfaceStyle"
    />
    <div class="absolute inset-0 bg-[linear-gradient(180deg,rgba(3,12,18,0.08),rgba(3,12,18,0.92)),radial-gradient(circle_at_top_right,rgba(255,255,255,0.04),transparent_28%)]" />
    <div
      class="absolute -right-10 top-0 h-40 w-40 rotate-12 blur-2xl"
      :style="orbStyle"
    />
    <div
      class="absolute inset-x-0 bottom-0 h-px"
      :style="edgeStyle"
    />

    <div class="relative z-10 flex flex-col gap-8 lg:flex-row lg:items-end lg:justify-between">
      <div class="max-w-3xl space-y-5">
        <p class="text-xs font-bold uppercase tracking-[0.5em] text-cyan-300/80">
          {{ eyebrow }}
        </p>
        <div class="space-y-3">
          <h1 class="text-4xl font-black uppercase tracking-[-0.05em] text-white sm:text-6xl">
            {{ title }}
          </h1>
          <p class="max-w-2xl text-sm leading-7 text-zinc-200 sm:text-base">
            {{ description }}
          </p>
        </div>
      </div>

      <div class="grid gap-3 sm:grid-cols-3 lg:w-[28rem] lg:min-w-[28rem]">
        <div
          v-for="stat in stats"
          :key="stat.label"
          class="rounded-[0.75rem] border border-white/10 bg-[#091924] px-4 py-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.03)]"
        >
          <p class="text-[11px] font-bold uppercase tracking-[0.35em] text-cyan-300/80">
            {{ stat.label }}
          </p>
          <p class="mt-2 text-2xl font-black text-white">
            {{ stat.value }}
          </p>
          <p
            v-if="stat.hint"
            class="mt-1 text-xs text-zinc-300/70"
          >
            {{ stat.hint }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const props = defineProps<{
  eyebrow: string
  title: string
  description: string
  accent?: string
  stats: Array<{ label: string, value: string, hint?: string }>
}>()

const heroSurfaceStyle = computed(() => ({
  background: `linear-gradient(135deg, ${props.accent ?? '#22d3ee'}60, rgba(5,19,28,0.92) 38%, rgba(2,12,19,0.98) 100%)`
}))

const orbStyle = computed(() => ({
  background: `${props.accent ?? '#22d3ee'}a0`
}))

const edgeStyle = computed(() => ({
  background: `linear-gradient(90deg, transparent, ${props.accent ?? '#22d3ee'}, transparent)`
}))
</script>

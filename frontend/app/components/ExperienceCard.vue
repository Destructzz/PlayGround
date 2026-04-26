<template>
  <article class="group relative overflow-hidden rounded-[0.85rem] border border-white/10 bg-[#07141d] p-6 shadow-[0_14px_48px_rgba(0,0,0,0.55)] transition-all duration-500 hover:-translate-y-1 hover:border-cyan-300/50">
    <div
      class="absolute inset-0 opacity-0 transition-opacity duration-500 group-hover:opacity-100"
      :style="glowStyle"
    />
    <div class="absolute inset-x-0 top-0 h-px bg-white/10" />
    <div class="relative z-10 flex h-full flex-col gap-5">
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="text-[11px] font-bold uppercase tracking-[0.35em] text-cyan-300/80">
            {{ eyebrow }}
          </p>
          <h3 class="mt-2 text-2xl font-black text-white">
            {{ title }}
          </h3>
        </div>
        <div
          class="rounded-[0.45rem] border px-3 py-1 text-[11px] font-bold uppercase tracking-[0.25em]"
          :class="badgeClass"
        >
          {{ badge }}
        </div>
      </div>

      <p class="text-sm leading-7 text-zinc-200">
        {{ description }}
      </p>

      <div class="grid gap-2 text-sm text-zinc-200">
        <div
          v-for="item in details"
          :key="item"
          class="flex items-center gap-2"
        >
          <span
            class="h-1.5 w-1.5 rounded-full"
            :style="dotStyle"
          />
          <span>{{ item }}</span>
        </div>
      </div>

      <div class="mt-auto flex items-end justify-between gap-4 pt-4">
        <div>
          <p class="text-[11px] uppercase tracking-[0.35em] text-cyan-300/70">
            {{ footerLabel }}
          </p>
          <p class="mt-2 text-lg font-bold text-white">
            {{ footerValue }}
          </p>
        </div>
        <slot name="cta" />
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
const props = defineProps<{
  eyebrow: string
  title: string
  description: string
  details: string[]
  badge: string
  badgeTone?: 'default' | 'warning' | 'success'
  footerLabel: string
  footerValue: string
  accent?: string
}>()

const glowStyle = computed(() => ({
  background: `linear-gradient(145deg, ${props.accent ?? '#22d3ee'}55, transparent 42%, rgba(255,255,255,0.01))`
}))

const dotStyle = computed(() => ({
  backgroundColor: props.accent ?? '#22d3ee'
}))

const badgeClass = computed(() => {
  if (props.badgeTone === 'warning') {
    return 'border-orange-300/55 bg-orange-500/20 text-orange-100'
  }

  if (props.badgeTone === 'success') {
    return 'border-emerald-300/55 bg-emerald-500/20 text-emerald-100'
  }

  return 'border-cyan-300/55 bg-cyan-500/20 text-cyan-50'
})
</script>

<template>
  <article
    :style="{ '--accent': accent || '#22d3ee' }"
    class="group relative overflow-hidden rounded-xl border border-white/5 bg-[#030c14]/80 backdrop-blur-md p-6 shadow-[0_12px_36px_rgba(0,0,0,0.6)] transition-all duration-500 hover:-translate-y-1 hover:border-[var(--accent)]"
  >
    <div
      class="absolute inset-0 opacity-0 transition-opacity duration-500 group-hover:opacity-100 pointer-events-none"
      :style="glowStyle"
    />
    <div class="absolute inset-x-0 top-0 h-px bg-white/5" />
    <div class="relative z-10 flex h-full flex-col gap-5">
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="text-[10px] font-bold uppercase tracking-[0.3em]" :style="{ color: accent || '#22d3ee' }">
            {{ eyebrow }}
          </p>
          <h3 class="mt-2 text-2xl font-black tracking-tight text-white uppercase">
            {{ title }}
          </h3>
        </div>
        <div
          class="rounded px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-[0.2em] border"
          :style="badgeStyle"
        >
          {{ badge }}
        </div>
      </div>

      <p class="text-xs leading-relaxed text-zinc-300">
        {{ description }}
      </p>

      <slot />

      <div class="grid gap-2 text-xs text-zinc-300">
        <div
          v-for="item in details"
          :key="item"
          class="flex items-center gap-2"
        >
          <span
            class="h-1 w-1 rounded-full animate-pulse"
            :style="dotStyle"
          />
          <span>{{ item }}</span>
        </div>
      </div>

      <div class="mt-auto flex items-end justify-between gap-4 pt-4 border-t border-white/5">
        <div>
          <p class="text-[10px] uppercase tracking-[0.3em] opacity-80" :style="{ color: accent || '#22d3ee' }">
            {{ footerLabel }}
          </p>
          <p class="mt-1 text-base font-bold text-white uppercase">
            {{ footerValue }}
          </p>
        </div>
        <slot name="cta" />
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    eyebrow: string
    title: string
    description: string
    details: string[]
    badge: string
    badgeTone?: 'default' | 'warning' | 'success'
    footerLabel: string
    footerValue: string
    accent?: string
  }>(),
  {
    badgeTone: 'default',
    accent: '#22d3ee'
  }
)

const glowStyle = computed(() => ({
  background: `linear-gradient(145deg, ${props.accent}15, transparent 42%, rgba(255,255,255,0.01))`
}))

const dotStyle = computed(() => ({
  backgroundColor: props.accent
}))

const badgeStyle = computed(() => {
  if (props.badgeTone === 'warning') {
    return {
      borderColor: 'rgba(249, 115, 22, 0.4)',
      backgroundColor: 'rgba(249, 115, 22, 0.1)',
      color: '#ffedd5'
    }
  }

  if (props.badgeTone === 'success') {
    return {
      borderColor: 'rgba(16, 185, 129, 0.4)',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      color: '#d1fae5'
    }
  }

  return {
    borderColor: `${props.accent}40`,
    backgroundColor: `${props.accent}10`,
    color: '#ffffff'
  }
})
</script>


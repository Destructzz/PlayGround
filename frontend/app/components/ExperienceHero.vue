<template>
  <section
    :style="{ '--accent-glow': accent || '#22d3ee' }"
    class="relative overflow-hidden rounded-xl border border-white/10 bg-[#071724]/75 backdrop-blur-md px-6 py-8 sm:px-8 sm:py-10 shadow-[0_12px_40px_rgba(0,0,0,0.45)]"
  >
    <!-- Coordinated ambient glowing backdrops -->
    <div
      class="absolute inset-0 opacity-60 pointer-events-none transition-all duration-500"
      :style="heroSurfaceStyle"
    />
    
    <!-- Beautiful colored glow orb on the right side -->
    <div
      class="absolute right-0 top-1/2 -translate-y-1/2 h-32 w-72 rounded-full blur-3xl opacity-35 pointer-events-none"
      :style="orbStyle"
    />

    <div class="relative z-10 max-w-3xl space-y-4">
      <div class="inline-flex items-center gap-2">
        <!-- Tiny pulsing indicator of theme accent color -->
        <span class="h-1.5 w-1.5 rounded-full animate-pulse" :style="{ backgroundColor: accent || '#22d3ee' }" />
        <p class="text-[10px] font-bold uppercase tracking-[0.3em] opacity-90" :style="{ color: accent || '#22d3ee' }">
          {{ eyebrow }}
        </p>
      </div>
      <div class="space-y-2">
        <h1 class="text-3xl font-black tracking-tight text-white sm:text-4xl uppercase">
          {{ title }}
        </h1>
        <p class="max-w-2xl text-xs leading-relaxed text-zinc-200 sm:text-sm">
          {{ description }}
        </p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
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

const heroSurfaceStyle = computed(() => ({
  background: `linear-gradient(135deg, ${props.accent}30, rgba(7,23,36,0.4) 60%, transparent 100%)`
}))

const orbStyle = computed(() => ({
  background: props.accent
}))
</script>


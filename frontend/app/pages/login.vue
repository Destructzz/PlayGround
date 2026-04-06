<template>
  <div class="flex-1 flex flex-col items-center justify-center p-8 relative min-h-full">
    
    <!-- Title -->
    <div class="mb-12 text-center flex flex-col items-center gap-1">
      <span class="text-cyan-100/60 text-lg font-medium tracking-widest uppercase mb-2">Sign in to</span>
      <div class="flex items-center font-black text-[2.75rem] tracking-tighter select-none">
        <span class="text-white drop-shadow-[0_0_12px_rgba(255,255,255,0.9)]">Play</span>
        <span class="text-cyan-300 drop-shadow-[0_0_18px_rgba(103,232,249,0.9)]">Ground</span>
      </div>
    </div>

    <!-- Login Form Wrapper with Animated Trace -->
    <div class="relative w-full max-w-md mx-auto group">
      
      <!-- Outer soft glow that fades in on hover -->
      <div class="absolute -inset-1 bg-cyan-500/20 rounded-[1.25rem] blur-xl opacity-50 group-hover:opacity-100 transition duration-700"></div>

      <!-- First Spark: Neon Cyan -->
      <div class="absolute inset-0 rounded-2xl overflow-hidden bg-cyan-950/30">
        <div class="absolute top-1/2 left-1/2 w-[200%] h-[200%] -translate-x-1/2 -translate-y-1/2 bg-[conic-gradient(from_0deg_at_50%_50%,transparent_60%,rgba(34,211,238,0.2)_85%,#22d3ee_98%,#ffffff_100%)] animate-[spin_6s_linear_infinite]"></div>
      </div>
      
      <!-- Second Spark: Neon Orange (Complementary color) chasing opposite -->
      <div class="absolute inset-0 rounded-2xl overflow-hidden mix-blend-screen">
        <div class="absolute top-1/2 left-1/2 w-[200%] h-[200%] -translate-x-1/2 -translate-y-1/2 bg-[conic-gradient(from_180deg_at_50%_50%,transparent_60%,rgba(249,115,22,0.2)_85%,#f97316_98%,#ffffff_100%)] animate-[spin_6s_linear_infinite]"></div>
      </div>

      <!-- Form Content layer -->
      <!-- Margin 2px exposes the background underneath as a border -->
      <div class="relative z-10 bg-[#020c13] m-[2px] rounded-[14px] p-8 sm:p-10 shadow-[0_0_40px_rgba(2,12,19,0.8)] flex flex-col gap-7">
        
        <form @submit.prevent class="flex flex-col gap-5">
          <!-- Login Field -->
          <div class="flex flex-col gap-2">
            <label for="login" class="text-xs font-bold tracking-wide text-cyan-100/70 uppercase">Login / Username</label>
            <input 
              type="text" 
              id="login" 
              class="w-full bg-[#03141f] border border-cyan-900/40 rounded-xl px-4 py-3.5 text-white placeholder-cyan-100/20 focus:outline-none focus:border-cyan-400 focus:ring-1 focus:ring-cyan-400 focus:bg-[#041a27] transition-all font-medium"
              placeholder="Enter your login"
            >
          </div>

          <!-- Password Field -->
          <div class="flex flex-col gap-2">
            <label for="password" class="text-xs font-bold tracking-wide text-cyan-100/70 uppercase">Password</label>
            <input 
              type="password" 
              id="password" 
              class="w-full bg-[#03141f] border border-cyan-900/40 rounded-xl px-4 py-3.5 text-white placeholder-cyan-100/20 focus:outline-none focus:border-cyan-400 focus:ring-1 focus:ring-cyan-400 focus:bg-[#041a27] transition-all font-medium"
              placeholder="••••••••"
            >
          </div>

          <!-- Submit Button -->
          <button 
            type="submit" 
            class="mt-4 w-full bg-cyan-400 hover:bg-cyan-300 text-[#020c13] text-lg font-bold px-6 py-4 rounded-xl transition-all shadow-[0_0_15px_rgba(34,211,238,0.4)] hover:shadow-[0_0_20px_rgba(34,211,238,0.7)] active:scale-[0.98]"
          >
            Sign In
          </button>
        </form>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Sign In - PlayGround'
})
</script>

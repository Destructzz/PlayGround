<template>
  <div class="flex flex-col w-full min-h-screen pb-16">
    <!-- Video Section (approx 66vh height) -->
    <div class="relative w-full h-[66vh] overflow-hidden flex items-center justify-center">
      <!-- Video Background -->
      <video
        autoplay
        loop
        muted
        playsinline
        class="absolute inset-0 w-full h-full object-cover"
      >
        <source
          src="/background1.mp4"
          type="video/mp4"
        >
      </video>

      <!-- Overlay to make text and content pop -->
      <div class="absolute inset-0 bg-black/60 pointer-events-none" />

      <!-- PlayGround Text -->
      <div class="relative z-10 text-center pointer-events-none group">
        <h1 class="text-6xl sm:text-7xl md:text-8xl lg:text-[8rem] font-extrabold tracking-tighter mt-12 text-transparent bg-clip-text bg-gradient-to-br from-white to-cyan-400 drop-shadow-[0_0_25px_rgba(34,211,238,0.4)] group-hover:drop-shadow-[0_0_40px_rgba(34,211,238,0.6)] transition-all duration-500 uppercase leading-none pb-4">
          PlayGround
        </h1>
      </div>

      <!-- Fade out to black at the bottom of the video so it blends with the dark background -->
      <div class="absolute bottom-0 w-full h-32 bg-gradient-to-t from-black to-transparent pointer-events-none" />
    </div>

    <!-- 3 Forms/Cards overlapping the video by half -->
    <div class="relative z-20 max-w-[1200px] mx-auto px-4 sm:px-6 w-full -mt-24">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Form/Card 1 -->
        <div class="bg-zinc-950/60 backdrop-blur-2xl border border-white/5 p-6 rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.8)] flex flex-col h-64 transition-all duration-500 hover:shadow-[0_0_30px_rgba(34,211,238,0.15)] hover:border-cyan-500/20 group relative overflow-hidden">
          <div class="absolute inset-0 bg-gradient-to-br from-cyan-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />

          <h2 class="text-2xl font-bold text-white mb-1.5">
            Игровая (Game)
          </h2>
          <p class="text-zinc-400 text-xs mb-4 flex-1">
            Мощные ПК, мониторы 240Hz и топовая периферия для идеальной игры.
          </p>

          <form
            class="space-y-3 relative z-10 mt-auto"
            @submit.prevent
          >
            <NuxtLink
              to="/gaming"
              data-testid="home-gaming-cta"
              class="block w-full text-center bg-white text-black font-bold py-2.5 rounded-xl text-sm hover:bg-zinc-200 transition-colors active:scale-[0.98]"
            >Забронировать ПК</NuxtLink>
          </form>
        </div>

        <!-- Form/Card 2 -->
        <div class="bg-zinc-950/60 backdrop-blur-2xl border border-white/5 p-6 rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.8)] flex flex-col h-64 transition-all duration-500 hover:shadow-[0_0_30px_rgba(249,115,22,0.15)] hover:border-orange-500/20 group relative overflow-hidden">
          <div class="absolute inset-0 bg-gradient-to-br from-orange-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />

          <h2 class="text-2xl font-bold text-white mb-1.5">
            Lounge зона
          </h2>
          <p class="text-zinc-400 text-xs mb-4 flex-1">
            Релакс с друзьями: комфортные диваны, кальяны и освежающие напитки.
          </p>

          <form
            class="space-y-3 relative z-10 mt-auto"
            @submit.prevent
          >
            <NuxtLink
              to="/lounge"
              data-testid="home-lounge-cta"
              class="block w-full text-center bg-white text-black font-bold py-2.5 rounded-xl text-sm hover:bg-zinc-200 transition-colors active:scale-[0.98]"
            >Забронировать столик</NuxtLink>
          </form>
        </div>

        <!-- Form/Card 3 -->
        <div class="bg-zinc-950/60 backdrop-blur-2xl border border-white/5 p-6 rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.8)] flex flex-col h-64 transition-all duration-500 hover:shadow-[0_0_30px_rgba(168,85,247,0.15)] hover:border-purple-500/20 group relative overflow-hidden">
          <div class="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />

          <h2 class="text-2xl font-bold text-white mb-1.5">
            Мероприятия (Event)
          </h2>
          <p class="text-zinc-400 text-xs mb-4 flex-1">
            Идеальная площадка для турниров, корпоративов и дней рождения.
          </p>

          <form
            class="space-y-3 relative z-10 mt-auto"
            @submit.prevent
          >
            <NuxtLink
              to="/event"
              data-testid="home-event-cta"
              class="block w-full text-center bg-white text-black font-bold py-2.5 rounded-xl text-sm hover:bg-zinc-200 transition-colors active:scale-[0.98]"
            >Оставить заявку</NuxtLink>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

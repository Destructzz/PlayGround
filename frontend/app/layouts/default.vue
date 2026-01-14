<template>
  <div class="relative min-h-screen text-gray-100">
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute -right-10 top-10 h-96 w-96 rounded-full bg-gradient-to-br from-white/10 via-purple-500/10 to-cyan-400/10 blur-3xl" />
      <div class="absolute left-10 bottom-20 h-80 w-80 rounded-full bg-gradient-to-tr from-indigo-500/10 via-sky-400/10 to-transparent blur-3xl" />
    </div>

    <header class="sticky top-0 z-20 border-b border-white/10 bg-black/40 backdrop-blur-xl">
      <div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-3">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-white/10 text-sm font-semibold text-white">PG</div>
          <div class="leading-tight">
            <div class="text-xs uppercase tracking-[0.18em] text-gray-400">PlayGround</div>
            <div class="text-sm font-semibold text-white">АИС досугового центра</div>
          </div>
        </div>
        <nav class="hidden items-center gap-4 md:flex">
          <NuxtLink to="/" class="nav-link">Главная</NuxtLink>
          <NuxtLink to="/services" class="nav-link">Услуги ▾</NuxtLink>
          <NuxtLink to="/zones" class="nav-link">Зоны ▾</NuxtLink>
          <NuxtLink to="/bookings" class="nav-link">Бронирования</NuxtLink>
          <NuxtLink to="/schedule" class="nav-link">Расписание</NuxtLink>
          <NuxtLink to="/admin" class="nav-link">Админ</NuxtLink>
        </nav>
        <div class="flex items-center gap-3 text-sm">
          <NuxtLink to="/auth/login" class="nav-link">Войти</NuxtLink>
          <NuxtLink to="/auth/register" class="solid-button">Регистрация</NuxtLink>
        </div>
      </div>
    </header>

    <main class="relative z-10 mx-auto max-w-6xl px-4 py-10 space-y-10">
      <slot />
    </main>

    <footer class="relative z-10 border-t border-white/10 bg-black/40 py-6 text-center text-sm text-gray-500 backdrop-blur-xl">
      © PlayGround, концепт интерфейса. Все данные демо.
    </footer>
  </div>
</template>

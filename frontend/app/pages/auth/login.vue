<template>
  <section class="mx-auto max-w-md space-y-6">
    <div class="space-y-2 text-center">
      <p class="section-title">PlayGround</p>
      <h1 class="text-2xl font-bold text-white">Вход</h1>
      <p class="text-gray-400">Демо-форма. OAuth/обычный вход появятся позже.</p>
    </div>

    <form class="space-y-4 rounded-2xl border border-white/10 bg-white/5 p-5 shadow-[0_20px_60px_rgba(0,0,0,0.45)] backdrop-blur-xl">
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Email
        <input type="email" placeholder="you@example.com" class="input-dark" />
      </label>
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Пароль
        <input type="password" placeholder="••••••••" class="input-dark" />
      </label>
      <button type="button" class="solid-button w-full">
        Войти
      </button>
      <NuxtLink to="/auth/register" class="block text-center text-sm font-semibold text-gray-200 hover:text-white">
        Создать аккаунт
      </NuxtLink>
    </form>
  </section>
</template>

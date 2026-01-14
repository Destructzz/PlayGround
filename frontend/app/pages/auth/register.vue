<template>
  <section class="mx-auto max-w-md space-y-6">
    <div class="space-y-2 text-center">
      <p class="section-title">PlayGround</p>
      <h1 class="text-2xl font-bold text-white">Регистрация</h1>
      <p class="text-gray-400">Статическая форма клиента. OAuth подключится позже.</p>
    </div>

    <form class="space-y-4 rounded-2xl border border-white/10 bg-white/5 p-5 shadow-[0_20px_60px_rgba(0,0,0,0.45)] backdrop-blur-xl">
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Имя и фамилия
        <input type="text" placeholder="Иван Петров" class="input-dark" />
      </label>
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Email
        <input type="email" placeholder="you@example.com" class="input-dark" />
      </label>
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Телефон
        <input type="tel" placeholder="+7 (___) ___-__-__" class="input-dark" />
      </label>
      <label class="flex flex-col gap-2 text-sm font-medium text-gray-200">
        Пароль
        <input type="password" placeholder="••••••••" class="input-dark" />
      </label>
      <button type="button" class="solid-button w-full">
        Зарегистрироваться
      </button>
      <NuxtLink to="/auth/login" class="block text-center text-sm font-semibold text-gray-200 hover:text-white">
        Уже есть аккаунт? Войти
      </NuxtLink>
    </form>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { login, register } from '../services/auth';

const router = useRouter();

const { t } = useI18n();

const isRegisterMode = ref(false);

const email = ref('');
const password = ref('');
const error = ref('');
const loading = ref(false);

const tenantName = ref('');
const tenantDomain = ref('');
const userName = ref('');

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  try {
    if (!email.value || !password.value) {
      error.value = t('login.requiredFields');
      loading.value = false;
      return;
    }
    await login({
      email: email.value,
      password: password.value
    });
    router.push({ name: 'Dashboard' });
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
};

const handleRegister = async () => {
  loading.value = true;
  error.value = '';
  try {
    if (!tenantName.value || !tenantDomain.value || !userName.value || !email.value || !password.value) {
      error.value = t('register.requiredFields');
      loading.value = false;
      return;
    }
    await register({
      tenant_name: tenantName.value,
      tenant_domain: tenantDomain.value,
      user_name: userName.value,
      email: email.value,
      password: password.value
    });
    router.push({ name: 'Dashboard' });
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
};

const toggleMode = () => {
  isRegisterMode.value = !isRegisterMode.value;
  error.value = '';
  email.value = '';
  password.value = '';
  tenantName.value = '';
  tenantDomain.value = '';
  userName.value = '';
};
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-slate-950 p-4">
    <div class="w-full max-w-md rounded-lg bg-slate-900 p-8 shadow-xl ring-1 ring-slate-800/50">
      <h2 class="mb-6 text-center text-2xl font-bold text-slate-100">
        {{ isRegisterMode ? t('register.title') : t('login.title') }}
      </h2>
      
      <form @submit.prevent="isRegisterMode ? handleRegister() : handleLogin()" class="space-y-4">
        <div v-if="isRegisterMode">
          <label for="tenant_name" class="block text-sm font-medium text-slate-300">{{ t('register.tenantName') }}</label>
          <input 
            id="tenant_name" 
            v-model="tenantName" 
            type="text" 
            :placeholder="t('register.tenantNamePlaceholder')"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div v-if="isRegisterMode">
          <label for="tenant_domain" class="block text-sm font-medium text-slate-300">{{ t('register.tenantDomain') }}</label>
          <input 
            id="tenant_domain" 
            v-model="tenantDomain" 
            type="text" 
            :placeholder="t('register.tenantDomainPlaceholder')"
            required
            pattern="[a-z0-9-]+"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
          <p class="mt-1 text-xs text-slate-500">Apenas letras minúsculas, números e hífens</p>
        </div>

        <div v-if="isRegisterMode">
          <label for="user_name" class="block text-sm font-medium text-slate-300">{{ t('register.userName') }}</label>
          <input 
            id="user_name" 
            v-model="userName" 
            type="text" 
            :placeholder="t('register.userNamePlaceholder')"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-slate-300">{{ t('login.email') }}</label>
          <input 
            id="email" 
            v-model="email" 
            type="email" 
            :placeholder="t(isRegisterMode ? 'register.emailPlaceholder' : 'login.emailPlaceholder')"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-slate-300">{{ t('login.password') }}</label>
          <input 
            id="password" 
            v-model="password" 
            type="password" 
            :placeholder="t(isRegisterMode ? 'register.passwordPlaceholder' : 'login.passwordPlaceholder')"
            :minlength="isRegisterMode ? 6 : undefined"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div v-if="error" class="rounded bg-red-900/50 p-2 text-center text-sm text-red-200 ring-1 ring-red-900/50">
          {{ error }}
        </div>

        <button 
          type="submit" 
          :disabled="loading"
          class="w-full rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50"
        >
          {{ loading ? (isRegisterMode ? t('register.creating') : t('login.entering')) : (isRegisterMode ? t('register.create') : t('login.enter')) }}
        </button>

        <div class="mt-4 text-center">
          <button
            type="button"
            @click="toggleMode"
            class="text-sm text-blue-400 hover:text-blue-300 transition-colors duration-200"
          >
            {{ isRegisterMode ? t('login.alreadyHaveAccount') : t('login.dontHaveAccount') }}
            <span class="font-medium">{{ isRegisterMode ? t('login.title') : t('login.register') }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

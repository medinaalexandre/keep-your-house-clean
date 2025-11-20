<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { availableLocales } from '../i18n';

const { locale, t } = useI18n();

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
}>();

const changeLanguage = (langCode: string) => {
  locale.value = langCode;
  localStorage.setItem('locale', langCode);
  emit('close');
};

const handleClose = () => {
  emit('close');
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    @click.self="handleClose"
  >
    <div class="w-full max-w-2xl rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
      <div class="mb-6 flex items-center justify-between">
        <h2 class="text-2xl font-bold text-slate-100">{{ t('settings.selectLanguage') }}</h2>
        <button
          @click="handleClose"
          class="rounded-md p-2 text-slate-400 hover:bg-slate-800 hover:text-slate-100"
        >
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <button
          v-for="lang in availableLocales"
          :key="lang.code"
          @click="changeLanguage(lang.code)"
          class="group flex flex-col items-center gap-3 rounded-lg border-2 p-6 transition-all duration-200"
          :class="locale === lang.code 
            ? 'border-blue-500 bg-blue-500/20' 
            : 'border-slate-700 bg-slate-800/50 hover:border-blue-500 hover:bg-slate-800'"
        >
          <span class="text-6xl">{{ lang.flag }}</span>
          <span class="text-lg font-medium text-slate-100">{{ lang.name }}</span>
        </button>
      </div>
    </div>
  </div>
</template>


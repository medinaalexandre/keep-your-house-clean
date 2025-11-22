<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { getUserComplimentsHistory, type ComplimentWithUser } from '../services/compliments';
import { getCurrentUser } from '../services/users';

const { t, locale } = useI18n();
const router = useRouter();

const compliments = ref<ComplimentWithUser[]>([]);
const loading = ref(false);
const error = ref('');

const currentUser = getCurrentUser();

const formatDate = (dateString: string | null): string => {
  if (!dateString) return t('common.noDate');
  const date = new Date(dateString);
  const localeMap: Record<string, string> = {
    'pt-BR': 'pt-BR',
    'es': 'es-ES',
    'en': 'en-US',
  };
  return date.toLocaleDateString(localeMap[locale.value] || 'pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const isReceived = (compliment: ComplimentWithUser): boolean => {
  return currentUser ? compliment.to_user_id === currentUser.id : false;
};

const loadCompliments = async () => {
  loading.value = true;
  error.value = '';
  try {
    const data = await getUserComplimentsHistory();
    compliments.value = data || [];
  } catch (e: any) {
    error.value = e.message;
    compliments.value = [];
  } finally {
    loading.value = false;
  }
};

const goToDashboard = () => {
  router.push({ name: 'Dashboard' });
};

onMounted(() => {
  loadCompliments();
});
</script>

<template>
  <div class="min-h-screen bg-slate-950 p-6">
    <div class="mx-auto max-w-4xl">
      <div class="mb-6 flex items-center justify-between">
        <h1 class="text-3xl font-bold text-slate-100">{{ t('complimentsHistory.title') }}</h1>
        <button
          @click="goToDashboard"
          class="cursor-pointer flex items-center gap-2 rounded-md border border-slate-700 bg-slate-800 px-4 py-2 text-slate-300 transition-colors duration-200 hover:bg-slate-700"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>{{ t('complimentsHistory.backToDashboard') }}</span>
        </button>
      </div>

      <div class="rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
        <div v-if="loading" class="py-8 text-center text-slate-400">
          {{ t('dashboard.loadingMore') }}
        </div>

        <div v-else-if="error" class="rounded bg-red-900/50 p-3 text-sm text-red-200 ring-1 ring-red-900/50">
          {{ error }}
        </div>

        <div v-else-if="compliments.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
          <svg class="mb-4 h-32 w-32 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-lg font-medium text-slate-300">{{ t('complimentsHistory.noCompliments') }}</p>
        </div>

        <ul v-else class="space-y-4">
          <li
            v-for="compliment in compliments"
            :key="compliment.id"
            class="rounded-lg border p-4 transition-all duration-200"
            :class="isReceived(compliment) 
              ? 'border-yellow-500/50 bg-gradient-to-r from-yellow-500/10 to-amber-500/10' 
              : 'border-blue-500/50 bg-gradient-to-r from-blue-500/10 to-cyan-500/10'"
          >
            <div class="flex items-start gap-4">
              <div class="flex-shrink-0">
                <svg 
                  class="h-6 w-6" 
                  :class="isReceived(compliment) ? 'text-yellow-500' : 'text-blue-500'"
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-4">
                  <div class="flex-1">
                    <div class="flex items-center gap-2 mb-1">
                      <span 
                        class="text-xs font-semibold px-2 py-1 rounded"
                        :class="isReceived(compliment) 
                          ? 'bg-yellow-500/20 text-yellow-400' 
                          : 'bg-blue-500/20 text-blue-400'"
                      >
                        {{ isReceived(compliment) ? t('complimentsHistory.received') : t('complimentsHistory.sent') }}
                      </span>
                      <h3 class="text-lg font-semibold text-slate-100">{{ compliment.title }}</h3>
                    </div>
                    <p v-if="compliment.description" class="mt-1 text-sm text-slate-300">{{ compliment.description }}</p>
                    <div class="mt-3 flex flex-wrap items-center gap-3 text-xs text-slate-400">
                      <span>
                        {{ isReceived(compliment) ? t('complimentsHistory.from') : t('complimentsHistory.to') }}
                        <span class="font-medium text-slate-300">{{ compliment.from_user_name || t('dashboard.unknownUser') }}</span>
                      </span>
                      <span>•</span>
                      <span>{{ formatDate(compliment.created_at) }}</span>
                      <span v-if="compliment.points > 0" 
                        class="font-medium"
                        :class="isReceived(compliment) ? 'text-yellow-400' : 'text-blue-400'"
                      >
                        {{ isReceived(compliment) ? '+' : '' }}{{ compliment.points }} {{ t('dashboard.points') }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter, useRoute } from 'vue-router';
import { getCompletedTasksByUser, type TaskWithUser } from '../services/tasks';
import { getUsers } from '../services/users';

const { t, locale } = useI18n();
const router = useRouter();
const route = useRoute();

const userId = ref<number | null>(null);
const userName = ref<string>('');
const tasks = ref<TaskWithUser[]>([]);
const loading = ref(false);
const loadingMore = ref(false);
const error = ref('');
const hasMoreTasks = ref(true);
const tasksContainer = ref<HTMLElement | null>(null);

const TASKS_PER_PAGE = 20;

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

const loadTasks = async (reset: boolean = false) => {
  if (!userId.value) return;

  if (reset) {
    tasks.value = [];
    hasMoreTasks.value = true;
  }

  if (!hasMoreTasks.value || loadingMore.value) return;

  loading.value = reset;
  loadingMore.value = !reset;
  error.value = '';

  try {
    const currentTasks = tasks.value;
    const offset = reset ? 0 : currentTasks.length;
    const data = await getCompletedTasksByUser(userId.value, TASKS_PER_PAGE, offset);
    const tasksArray = Array.isArray(data) ? data : [];

    if (reset) {
      tasks.value = tasksArray;
    } else {
      tasks.value = [...currentTasks, ...tasksArray];
    }

    hasMoreTasks.value = tasksArray.length === TASKS_PER_PAGE;
  } catch (e: any) {
    error.value = e.message;
    tasks.value = reset ? [] : tasks.value;
  } finally {
    loading.value = false;
    loadingMore.value = false;
  }
};

const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement;
  if (!target) return;

  const scrollBottom = target.scrollHeight - target.scrollTop - target.clientHeight;

  if (scrollBottom < 100 && hasMoreTasks.value && !loadingMore.value) {
    loadTasks(false);
  }
};

const goToDashboard = () => {
  router.push({ name: 'Dashboard' });
};

const loadUserInfo = async () => {
  try {
    const users = await getUsers();
    const user = users.find(u => u.id === userId.value);
    if (user) {
      userName.value = user.name;
    }
  } catch (e: any) {
    console.error('Failed to load user info:', e);
  }
};

onMounted(async () => {
  const userIdParam = route.params.userId;
  if (typeof userIdParam === 'string') {
    userId.value = parseInt(userIdParam, 10);
  } else if (Array.isArray(userIdParam) && userIdParam.length > 0 && typeof userIdParam[0] === 'string') {
    userId.value = parseInt(userIdParam[0], 10);
  }

  if (!userId.value || isNaN(userId.value)) {
    error.value = 'ID de usuário inválido';
    return;
  }

  await loadUserInfo();
  await loadTasks(true);

  setTimeout(() => {
    if (tasksContainer.value) {
      tasksContainer.value.addEventListener('scroll', handleScroll);
    }
  }, 100);
});

onUnmounted(() => {
  if (tasksContainer.value) {
    tasksContainer.value.removeEventListener('scroll', handleScroll);
  }
});
</script>

<template>
  <div class="min-h-screen bg-slate-950 p-6">
    <div class="mx-auto max-w-4xl">
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-slate-100">{{ t('userTasksHistory.title') }}</h1>
          <p v-if="userName" class="mt-1 text-sm text-slate-400">{{ t('userTasksHistory.subtitle', { name: userName }) }}</p>
        </div>
        <button
          @click="goToDashboard"
          class="cursor-pointer flex items-center gap-2 rounded-md border border-slate-700 bg-slate-800 px-4 py-2 text-slate-300 transition-colors duration-200 hover:bg-slate-700"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>{{ t('userTasksHistory.backToDashboard') }}</span>
        </button>
      </div>

      <div class="rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
        <div v-if="loading" class="py-8 text-center text-slate-400">
          {{ t('dashboard.loadingMore') }}
        </div>

        <div v-else-if="error" class="rounded bg-red-900/50 p-3 text-sm text-red-200 ring-1 ring-red-900/50">
          {{ error }}
        </div>

        <div v-else-if="tasks.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
          <svg class="mb-4 h-32 w-32 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p class="text-lg font-medium text-slate-300">{{ t('userTasksHistory.noTasks') }}</p>
        </div>

        <div
          v-else
          ref="tasksContainer"
          class="overflow-y-auto"
          style="max-height: 600px; scrollbar-width: none; -ms-overflow-style: none;"
        >
          <ul class="space-y-4">
            <li
              v-for="task in tasks"
              :key="task.id"
              class="rounded-lg border border-slate-700 bg-slate-800 p-4 transition-all duration-200 hover:border-green-500/50"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1">
                  <h3 class="text-lg font-semibold text-slate-100">{{ task.title }}</h3>
                  <p v-if="task.description" class="mt-1 text-sm text-slate-300">{{ task.description }}</p>
                  <div class="mt-3 flex flex-wrap items-center gap-3 text-xs text-slate-400">
                    <span>{{ t('dashboard.completedOn') + ': ' + formatDate(task.updated_at) }}</span>
                    <span v-if="task.points > 0" class="font-medium text-green-400">
                      +{{ task.points }} {{ t('dashboard.points') }}
                    </span>
                  </div>
                </div>
              </div>
            </li>

            <li v-if="loadingMore" class="py-4 text-center text-slate-400">
              {{ t('dashboard.loadingMore') }}
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overflow-y-auto::-webkit-scrollbar {
  display: none;
}
</style>


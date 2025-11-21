<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { getUpcomingTasks, getCompletedTasksHistory, undoCompleteTask, deleteTask, type Task, type TaskWithUser } from '../services/tasks';
import { getTopUsers, type User } from '../services/users';
import { logout } from '../services/auth';
import CreateTaskModal from './CreateTaskModal.vue';
import CreateUserModal from './CreateUserModal.vue';
import CompleteTaskModal from './CompleteTaskModal.vue';
import EditTaskModal from './EditTaskModal.vue';
import LanguageSelectorModal from './LanguageSelectorModal.vue';

const { t, locale } = useI18n();
const router = useRouter();

const upcomingTasks = ref<Task[]>([]);
const completedTasks = ref<TaskWithUser[]>([]);
const topUsers = ref<User[]>([]);

const ensureArray = <T>(value: T[] | null | undefined): T[] => {
  return Array.isArray(value) ? value : [];
};
const loadingUpcoming = ref(false);
const loadingHistory = ref(false);
const loadingRanking = ref(false);
const loadingMore = ref(false);
const errorUpcoming = ref('');
const errorHistory = ref('');
const errorRanking = ref('');
const isTaskModalOpen = ref(false);
const isUserModalOpen = ref(false);
const isCompleteModalOpen = ref(false);
const isEditModalOpen = ref(false);
const selectedTaskForComplete = ref<Task | null>(null);
const selectedTaskForEdit = ref<Task | null>(null);
const hasMoreTasks = ref(true);
const upcomingTasksContainer = ref<HTMLElement | null>(null);
const showFloatingMenu = ref(false);
const showSettingsMenu = ref(false);
const isLanguageModalOpen = ref(false);

const TASKS_PER_PAGE = 5;

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

const loadUpcomingTasks = async (reset: boolean = false) => {
  if (reset) {
    upcomingTasks.value = [];
    hasMoreTasks.value = true;
  }

  if (!hasMoreTasks.value || loadingMore.value) return;

  loadingUpcoming.value = reset;
  loadingMore.value = !reset;
  errorUpcoming.value = '';

  try {
    const currentTasks = ensureArray(upcomingTasks.value);
    const offset = reset ? 0 : currentTasks.length;
    const tasks = await getUpcomingTasks(TASKS_PER_PAGE, offset);
    const tasksArray = ensureArray(tasks);
    
    if (reset) {
      upcomingTasks.value = tasksArray;
    } else {
      upcomingTasks.value = [...currentTasks, ...tasksArray];
    }

    hasMoreTasks.value = tasksArray.length === TASKS_PER_PAGE;
  } catch (e: any) {
    errorUpcoming.value = e.message;
    upcomingTasks.value = reset ? [] : ensureArray(upcomingTasks.value);
  } finally {
    loadingUpcoming.value = false;
    loadingMore.value = false;
  }
};

const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement;
  if (!target) return;

  const scrollBottom = target.scrollHeight - target.scrollTop - target.clientHeight;
  
  if (scrollBottom < 100 && hasMoreTasks.value && !loadingMore.value) {
    loadUpcomingTasks(false);
  }
};

const loadCompletedTasksHistory = async () => {
  loadingHistory.value = true;
  errorHistory.value = '';
  try {
    const tasks = await getCompletedTasksHistory(5);
    completedTasks.value = ensureArray(tasks);
  } catch (e: any) {
    errorHistory.value = e.message;
    completedTasks.value = [];
  } finally {
    loadingHistory.value = false;
    completedTasks.value = ensureArray(completedTasks.value);
  }
};

const loadTopUsers = async () => {
  loadingRanking.value = true;
  errorRanking.value = '';
  try {
    const users = await getTopUsers();
    topUsers.value = ensureArray(users);
  } catch (e: any) {
    errorRanking.value = e.message;
    topUsers.value = [];
  } finally {
    loadingRanking.value = false;
    topUsers.value = ensureArray(topUsers.value);
  }
};

const handleTaskCreated = () => {
  loadUpcomingTasks(true);
  loadCompletedTasksHistory();
  loadTopUsers();
};

const handleUserCreated = () => {
  loadCompletedTasksHistory();
  loadTopUsers();
};

const handleCompleteClick = (task: Task) => {
  selectedTaskForComplete.value = task;
  isCompleteModalOpen.value = true;
};

const handleEditClick = (task: Task) => {
  selectedTaskForEdit.value = task;
  isEditModalOpen.value = true;
};

const handleTaskUpdated = () => {
  loadUpcomingTasks(true);
  loadCompletedTasksHistory();
  loadTopUsers();
};

const handleTaskCompleted = () => {
  loadUpcomingTasks(true);
  loadCompletedTasksHistory();
  loadTopUsers();
};

const undoingTaskId = ref<number | null>(null);

const handleUndoTask = async (taskId: number) => {
  undoingTaskId.value = taskId;
  try {
    await undoCompleteTask(taskId);
    loadUpcomingTasks(true);
    loadCompletedTasksHistory();
    loadTopUsers();
  } catch (e: any) {
    errorHistory.value = e.message;
  } finally {
    undoingTaskId.value = null;
  }
};

const deletingTaskId = ref<number | null>(null);

const handleDeleteTask = async (taskId: number) => {
  if (!confirm(t('dashboard.confirmDelete'))) {
    return;
  }

  deletingTaskId.value = taskId;
  try {
    await deleteTask(taskId);
    loadUpcomingTasks(true);
    loadCompletedTasksHistory();
    loadTopUsers();
  } catch (e: any) {
    errorUpcoming.value = e.message;
  } finally {
    deletingTaskId.value = null;
  }
};

const handleMenuMouseLeave = () => {
  setTimeout(() => {
    const menuContainer = document.querySelector('.floating-menu-container');
    if (menuContainer && !menuContainer.matches(':hover')) {
      showFloatingMenu.value = false;
    }
  }, 200);
};

const handleLogout = () => {
  logout();
  router.push({ name: 'Login' });
};

onMounted(async () => {
  await loadUpcomingTasks(true);
  loadCompletedTasksHistory();
  loadTopUsers();
  
  setTimeout(() => {
    if (upcomingTasksContainer.value) {
      upcomingTasksContainer.value.addEventListener('scroll', handleScroll);
    }
  }, 100);
});

onUnmounted(() => {
  if (upcomingTasksContainer.value) {
    upcomingTasksContainer.value.removeEventListener('scroll', handleScroll);
  }
});
</script>

<template>
  <div class="min-h-screen bg-slate-950 p-6">
    <div class="mx-auto max-w-7xl">
      <div class="mb-6 flex items-center justify-between">
        <h1 class="text-3xl font-bold text-slate-100">{{ t('dashboard.title') }}</h1>
        <div class="relative">
          <button
            @click="showSettingsMenu = !showSettingsMenu"
            class="cursor-pointer rounded-md p-2 text-slate-400 transition-colors duration-200 hover:bg-slate-800 hover:text-slate-100"
            aria-label="Configurações"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
          
          <Transition
            enter-active-class="transition ease-out duration-200"
            enter-from-class="opacity-0 translate-y-2"
            enter-to-class="opacity-100 translate-y-0"
            leave-active-class="transition ease-in duration-150"
            leave-from-class="opacity-100 translate-y-0"
            leave-to-class="opacity-0 translate-y-2"
          >
            <div
              v-if="showSettingsMenu"
              class="absolute right-0 top-full z-50 mt-2 w-56 rounded-lg bg-slate-800/50 backdrop-blur-sm p-1 shadow-xl ring-1 ring-slate-700/30"
            >
              <button
                @click.stop="isLanguageModalOpen = true; showSettingsMenu = false"
                class="cursor-pointer flex w-full items-center gap-3 whitespace-nowrap rounded-md px-4 py-2 text-left text-slate-100 transition-colors duration-200 hover:bg-slate-700/50"
              >
                <svg class="h-5 w-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
                </svg>
                <span>{{ t('settings.changeLanguage') }}</span>
              </button>
              <button
                @click.stop="handleLogout(); showSettingsMenu = false"
                class="cursor-pointer flex w-full items-center gap-3 whitespace-nowrap rounded-md px-4 py-2 text-left text-red-400 transition-colors duration-200 hover:bg-slate-700/50"
              >
                <svg class="h-5 w-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                <span>{{ t('settings.logout') }}</span>
              </button>
            </div>
          </Transition>
        </div>
      </div>
      
      <div
        v-if="showSettingsMenu"
        class="fixed inset-0 z-30"
        @click="showSettingsMenu = false"
      />
      
      <div class="mb-6 flex flex-col rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
        <h2 class="mb-4 text-xl font-semibold text-slate-100">{{ t('dashboard.ranking') }}</h2>
        
        <div v-if="loadingRanking" class="py-8 text-center text-slate-400">
          {{ t('dashboard.loadingMore') }}
        </div>
        
        <div v-else-if="errorRanking" class="rounded bg-red-900/50 p-3 text-sm text-red-200 ring-1 ring-red-900/50">
          {{ errorRanking }}
        </div>
        
        <div v-else-if="topUsers.length === 0" class="py-8 text-center text-slate-400">
          <p>{{ t('dashboard.noUpcomingTasks') }}</p>
        </div>
        
        <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div
            v-for="(user, index) in topUsers"
            :key="user.id"
            class="flex flex-col items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 p-4"
            :class="{ 'border-yellow-500 bg-yellow-500/10': index === 0, 'border-slate-600 bg-slate-800/30': index === 1, 'border-amber-600 bg-amber-600/10': index === 2 }"
          >
            <div class="text-2xl font-bold text-slate-100">{{ user.points }}</div>
            <div class="text-sm font-medium text-slate-300">{{ user.name }}</div>
            <div class="text-xs text-slate-500">{{ t('dashboard.points') }}</div>
          </div>
        </div>
      </div>
      
      <div class="grid gap-6 md:grid-cols-2">
        <div class="flex flex-col rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
          <h2 class="mb-4 text-xl font-semibold text-slate-100">{{ t('dashboard.upcomingTasks') }}</h2>
          
          <div class="relative flex-1">
            <div
              ref="upcomingTasksContainer"
              class="overflow-y-auto upcoming-tasks-scroll"
              style="max-height: 600px; scrollbar-width: none; -ms-overflow-style: none;"
            >
            <div v-if="loadingUpcoming" class="py-8 text-center text-slate-400">
              {{ t('dashboard.loadingMore') }}
            </div>
            
            <div v-else-if="errorUpcoming" class="rounded bg-red-900/50 p-3 text-sm text-red-200 ring-1 ring-red-900/50">
              {{ errorUpcoming }}
            </div>
            
            <div v-else-if="!upcomingTasks || upcomingTasks.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
              <svg class="mb-4 h-32 w-32 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
              <p class="text-lg font-medium text-slate-300">{{ t('dashboard.noUpcomingTasks') }}</p>
              <button
                @click="isTaskModalOpen = true"
                class="mt-2 cursor-pointer text-sm font-medium text-blue-400 transition-colors duration-200 hover:text-blue-300"
              >
                {{ t('dashboard.createFirstTask') }}
              </button>
            </div>
            
            <ul v-else class="space-y-3">
              <li
                v-for="task in upcomingTasks"
                :key="task.id"
                class="group relative rounded-md border border-slate-700 bg-slate-800 p-4 transition-all duration-200 hover:border-green-500/50"
              >
                <div class="flex items-start justify-between">
                  <div class="flex-1">
                    <h3 class="font-medium text-slate-100">{{ task.title }}</h3>
                    <p class="mt-1 text-sm text-slate-400">{{ task.description }}</p>
                    <div class="mt-2 flex items-center gap-4 text-xs text-slate-500">
                      <span>{{ formatDate(task.scheduled_to) }}</span>
                      <span>{{ task.points }} pontos</span>
                    </div>
                  </div>
                  <div v-if="!task.completed" class="relative z-10 ml-4 flex flex-shrink-0 gap-2 md:opacity-0 md:group-hover:opacity-100">
                    <button
                      @click.stop="handleEditClick(task)"
                      class="cursor-pointer rounded-md border border-blue-600 bg-blue-600/10 p-2 text-blue-400 transition-all duration-200 hover:bg-blue-600/20"
                      :title="t('dashboard.edit')"
                    >
                      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <button
                      @click.stop="handleDeleteTask(task.id)"
                      :disabled="deletingTaskId === task.id"
                      class="cursor-pointer rounded-md border border-red-600 bg-red-600/10 p-2 text-red-400 transition-all duration-200 hover:bg-red-600/20 disabled:cursor-not-allowed disabled:opacity-50"
                      :title="deletingTaskId === task.id ? t('dashboard.deleting') : t('dashboard.delete')"
                    >
                      <svg v-if="deletingTaskId !== task.id" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      <svg v-else class="h-5 w-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                      </svg>
                    </button>
                  </div>
                  <div v-else class="relative z-10 ml-4 flex flex-shrink-0 md:opacity-0 md:group-hover:opacity-100">
                    <button
                      @click.stop="handleDeleteTask(task.id)"
                      :disabled="deletingTaskId === task.id"
                      class="cursor-pointer rounded-md border border-red-600 bg-red-600/10 p-2 text-red-400 transition-all duration-200 hover:bg-red-600/20 disabled:cursor-not-allowed disabled:opacity-50"
                      :title="deletingTaskId === task.id ? t('dashboard.deleting') : t('dashboard.delete')"
                    >
                      <svg v-if="deletingTaskId !== task.id" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      <svg v-else class="h-5 w-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <div
                  v-if="!task.completed"
                  class="absolute inset-0 flex cursor-pointer items-center justify-center rounded-md bg-green-600/90 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
                  @click="handleCompleteClick(task)"
                >
                  <span class="font-semibold text-white">{{ t('completeTask.complete') }}</span>
                </div>
              </li>
              
              <li v-if="loadingMore" class="py-4 text-center text-slate-400">
                {{ t('dashboard.loadingMore') }}
              </li>
            </ul>
            </div>
            <div
              v-if="hasMoreTasks && upcomingTasks && upcomingTasks.length > 0 && !loadingUpcoming"
              class="pointer-events-none absolute bottom-0 left-0 right-0 h-20 bg-gradient-to-t from-slate-900 via-slate-900/60 to-transparent"
            />
          </div>
        </div>

        <div class="rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
          <h2 class="mb-4 text-xl font-semibold text-slate-100">{{ t('dashboard.taskHistory') }}</h2>
          
          <div v-if="loadingHistory" class="text-center text-slate-400">
            Carregando...
          </div>
          
          <div v-else-if="errorHistory" class="rounded bg-red-900/50 p-3 text-sm text-red-200 ring-1 ring-red-900/50">
            {{ errorHistory }}
          </div>
          
          <div v-else-if="!completedTasks || completedTasks.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
            <svg class="mb-4 h-32 w-32 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p class="text-lg font-medium text-slate-300">{{ t('dashboard.noCompletedTasks') }}</p>
            <p class="mt-2 text-sm text-slate-500">Complete tarefas pendentes para ver o histórico aqui</p>
          </div>
          
            <ul v-else-if="completedTasks && completedTasks.length > 0" class="space-y-3">
              <li
                v-for="task in completedTasks"
                :key="task.id"
                class="group relative rounded-md border border-slate-700 bg-slate-800 p-4 transition-all duration-200 hover:border-orange-600/50"
              >
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="font-medium text-slate-100">{{ task.title }}</h3>
                  <p class="mt-1 text-sm text-slate-400">{{ task.description }}</p>
                  <div class="mt-2 flex items-center gap-4 text-xs text-slate-500">
                    <span>{{ t('dashboard.completedOn') + ': ' + formatDate(task.updated_at) }}</span>
                    <span v-if="task.completed_by_name" class="text-slate-300">
                      {{ t('dashboard.completedBy') + ' ' + task.completed_by_name }}
                    </span>
                    <span v-else class="text-slate-500">{{ t('dashboard.completedBy') + ' ' + t('dashboard.unknownUser') }}</span>
                  </div>
                </div>
                <button
                  @click="handleUndoTask(task.id)"
                  :disabled="undoingTaskId === task.id"
                  class="ml-4 flex-shrink-0 cursor-pointer rounded-md border border-orange-600 bg-orange-600/10 p-2 text-orange-400 transition-all duration-200 hover:bg-orange-600/20 disabled:cursor-not-allowed disabled:opacity-50 md:opacity-0 md:group-hover:opacity-100 md:px-3 md:py-1.5"
                  :title="undoingTaskId === task.id ? t('dashboard.undoing') : t('dashboard.undo')"
                >
                  <span v-if="undoingTaskId === task.id" class="md:hidden">{{ t('dashboard.undoing') }}</span>
                  <span v-else class="md:hidden">{{ t('dashboard.undo') }}</span>
                  <svg v-if="undoingTaskId !== task.id" class="h-5 w-5 md:block" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                  </svg>
                  <svg v-else class="h-5 w-5 animate-spin md:block" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                </button>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <div class="fixed bottom-8 right-8 z-50 floating-menu-container">
      <div class="relative">
        <Transition
          enter-active-class="transition ease-out duration-200"
          enter-from-class="opacity-0 translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition ease-in duration-150"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 translate-y-2"
        >
          <div
            v-if="showFloatingMenu"
            @mouseenter="showFloatingMenu = true"
            @mouseleave="showFloatingMenu = false"
            class="absolute bottom-full right-0 mb-2 flex flex-col gap-1 rounded-lg bg-slate-800/50 backdrop-blur-sm p-1 shadow-xl ring-1 ring-slate-700/30"
          >
            <button
              @click="isTaskModalOpen = true; showFloatingMenu = false"
              class="cursor-pointer flex items-center justify-between gap-3 whitespace-nowrap rounded-md px-4 py-3 text-left text-slate-100 transition-colors duration-200 hover:bg-slate-700/50"
            >
              <span class="flex-1">{{ t('dashboard.createTask') }}</span>
              <svg class="h-5 w-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </button>
            <button
              @click="isUserModalOpen = true; showFloatingMenu = false"
              class="cursor-pointer flex items-center justify-between gap-3 whitespace-nowrap rounded-md px-4 py-3 text-left text-slate-100 transition-colors duration-200 hover:bg-slate-700/50"
            >
              <span class="flex-1">{{ t('dashboard.createUser') }}</span>
              <svg class="h-5 w-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </button>
          </div>
        </Transition>
        
        <button
          @click="isTaskModalOpen = true; showFloatingMenu = false"
          @mouseenter="showFloatingMenu = true"
          @mouseleave="handleMenuMouseLeave"
          class="cursor-pointer flex h-14 w-14 items-center justify-center rounded-full bg-blue-600 text-white shadow-lg transition-colors duration-200 hover:bg-blue-500 hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-950"
          aria-label="Criar nova tarefa"
        >
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
    </div>
    
    <div
      v-if="showFloatingMenu"
      class="fixed inset-0 z-30"
      @click="showFloatingMenu = false"
    />

    <CreateTaskModal :is-open="isTaskModalOpen" @close="isTaskModalOpen = false" @created="handleTaskCreated" />
    <CreateUserModal :is-open="isUserModalOpen" @close="isUserModalOpen = false" @created="handleUserCreated" />
    
    <CompleteTaskModal
      :is-open="isCompleteModalOpen"
      :task="selectedTaskForComplete"
      @close="isCompleteModalOpen = false; selectedTaskForComplete = null"
      @completed="handleTaskCompleted"
    />
    
    <EditTaskModal
      :is-open="isEditModalOpen"
      :task="selectedTaskForEdit"
      @close="isEditModalOpen = false; selectedTaskForEdit = null"
      @updated="handleTaskUpdated"
    />
    
    <LanguageSelectorModal
      :is-open="isLanguageModalOpen"
      @close="isLanguageModalOpen = false"
    />
  </div>
</template>

<style scoped>
  .upcoming-tasks-scroll::-webkit-scrollbar {
    display: none;
  }
</style>
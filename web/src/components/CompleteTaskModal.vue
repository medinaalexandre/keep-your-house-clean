<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { completeTask, type Task } from '../services/tasks';
import { getUsers, getCurrentUser, type User } from '../services/users';

const { t } = useI18n();

const props = defineProps<{
  isOpen: boolean;
  task: Task | null;
}>();

const emit = defineEmits<{
  close: [];
  completed: [];
}>();

const users = ref<User[]>([]);
const selectedUserId = ref<number | null>(null);
const loading = ref(false);
const loadingUsers = ref(false);
const error = ref('');

watch(() => props.isOpen, async (isOpen) => {
  if (isOpen && props.task) {
    const currentUser = getCurrentUser();
    await loadUsers();
    await nextTick();
    if (currentUser) {
      selectedUserId.value = currentUser.id;
    } else {
      selectedUserId.value = users.value.length > 0 ? users.value[0].id : null;
    }
  } else {
    selectedUserId.value = null;
    error.value = '';
  }
});

const loadUsers = async () => {
  loadingUsers.value = true;
  try {
    users.value = await getUsers();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loadingUsers.value = false;
  }
};

const handleClose = () => {
  selectedUserId.value = null;
  error.value = '';
  emit('close');
};

const handleSubmit = async () => {
  if (!props.task || !selectedUserId.value) return;

  loading.value = true;
  error.value = '';

  try {
    await completeTask(props.task.id, selectedUserId.value);
    emit('completed');
    handleClose();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div
    v-if="isOpen && task"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    @click.self="handleClose"
  >
    <div class="w-full max-w-md rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-2xl font-bold text-slate-100">{{ t('completeTask.title') }}</h2>
        <button
          @click="handleClose"
          class="rounded-md p-2 text-slate-400 hover:bg-slate-800 hover:text-slate-100"
        >
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="mb-4 rounded-md border border-slate-700 bg-slate-800 p-4">
        <h3 class="font-medium text-slate-100">{{ task.title }}</h3>
        <p class="mt-1 text-sm text-slate-400">{{ task.description }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="completed_by" class="block text-sm font-medium text-slate-300">
            {{ t('completeTask.whoCompleted') }}
          </label>
          <div v-if="loadingUsers" class="mt-2 text-sm text-slate-400">
            {{ t('completeTask.loadingUsers') }}
          </div>
          <select
            v-else
            id="completed_by"
            v-model="selectedUserId"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          >
            <option :value="null" disabled>{{ t('completeTask.selectUser') }}</option>
            <option
              v-for="user in users"
              :key="user.id"
              :value="user.id"
            >
              {{ user.name }} {{ user.id === getCurrentUser()?.id ? t('completeTask.you') : '' }}
            </option>
          </select>
        </div>

        <div v-if="error" class="rounded bg-red-900/50 p-3 text-sm text-red-200 ring-1 ring-red-900/50">
          {{ error }}
        </div>

        <div class="flex justify-end gap-3 pt-4">
          <button
            type="button"
            @click="handleClose"
            class="rounded-md border border-slate-700 bg-slate-800 px-4 py-2 text-slate-300 transition-colors duration-200 hover:bg-slate-700"
          >
            {{ t('completeTask.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="loading || !selectedUserId"
            class="rounded-md bg-green-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-green-500 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50"
          >
            {{ loading ? t('completeTask.completing') : t('completeTask.complete') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


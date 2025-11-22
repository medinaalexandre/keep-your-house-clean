<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { createCompliment, type CreateComplimentRequest } from '../services/compliments';
import { getUsers, type User } from '../services/users';
import { getCurrentUser } from '../services/users';

const { t } = useI18n();

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
  created: [];
}>();

const formData = ref<CreateComplimentRequest>({
  title: '',
  description: '',
  points: 0,
  to_user_id: 0,
});

const users = ref<User[]>([]);
const loadingUsers = ref(false);
const loading = ref(false);
const error = ref('');

const loadUsers = async () => {
  loadingUsers.value = true;
  try {
    const allUsers = await getUsers();
    const currentUser = getCurrentUser();
    users.value = allUsers.filter(user => currentUser && user.id !== currentUser.id);
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loadingUsers.value = false;
  }
};

const resetForm = () => {
  formData.value = {
    title: '',
    description: '',
    points: 0,
    to_user_id: 0,
  };
  error.value = '';
};

const handleClose = () => {
  resetForm();
  emit('close');
};

const handleSubmit = async () => {
  if (!formData.value.title.trim()) {
    error.value = t('compliment.titleRequired');
    return;
  }

  if (formData.value.to_user_id === 0) {
    error.value = t('compliment.userRequired');
    return;
  }

  if (formData.value.points < 0 || formData.value.points > 5) {
    error.value = t('compliment.invalidPoints');
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    await createCompliment(formData.value);
    resetForm();
    emit('created');
    handleClose();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadUsers();
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    @click.self="handleClose"
  >
    <div class="w-full max-w-2xl rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-2xl font-bold text-slate-100">{{ t('compliment.title') }}</h2>
        <button
          @click="handleClose"
          class="rounded-md p-2 text-slate-400 hover:bg-slate-800 hover:text-slate-100"
        >
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="to_user_id" class="block text-sm font-medium text-slate-300">{{ t('compliment.toUserLabel') }} *</label>
          <select
            id="to_user_id"
            v-model.number="formData.to_user_id"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          >
            <option value="0">{{ loadingUsers ? t('compliment.loadingUsers') : t('compliment.selectUser') }}</option>
            <option
              v-for="user in users"
              :key="user.id"
              :value="user.id"
            >
              {{ user.name }}
            </option>
          </select>
        </div>

        <div>
          <label for="title" class="block text-sm font-medium text-slate-300">{{ t('compliment.titleLabel') }} *</label>
          <input
            id="title"
            v-model="formData.title"
            type="text"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('compliment.titlePlaceholder')"
          />
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-slate-300">{{ t('compliment.descriptionLabel') }}</label>
          <textarea
            id="description"
            v-model="formData.description"
            rows="3"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('compliment.descriptionPlaceholder')"
          />
        </div>

        <div>
          <label for="points" class="block text-sm font-medium text-slate-300">{{ t('compliment.pointsLabel') }} ({{ t('compliment.maxPoints') }})</label>
          <input
            id="points"
            v-model.number="formData.points"
            type="number"
            min="0"
            max="5"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
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
            {{ t('compliment.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50"
          >
            {{ loading ? t('compliment.creating') : t('compliment.create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


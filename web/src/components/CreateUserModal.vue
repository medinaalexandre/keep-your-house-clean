<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { createUser, getCurrentUser, type CreateUserRequest } from '../services/users';

const { t } = useI18n();

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
  created: [];
}>();

const currentUser = getCurrentUser();

const formData = ref<CreateUserRequest>({
  name: '',
  email: '',
  password: '',
  tenant_id: currentUser?.tenant_id || 0,
  points: 0,
  role: 'user',
  status: 'active',
});

const loading = ref(false);
const error = ref('');

const resetForm = () => {
  formData.value = {
    name: '',
    email: '',
    password: '',
    tenant_id: currentUser?.tenant_id || 0,
    points: 0,
    role: 'user',
    status: 'active',
  };
  error.value = '';
};

const handleClose = () => {
  resetForm();
  emit('close');
};

const handleSubmit = async () => {
  if (!formData.value.name || !formData.value.email || !formData.value.password) {
    error.value = t('user.requiredFields');
    return;
  }

  if (!formData.value.tenant_id) {
    error.value = t('user.tenantNotFound');
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    await createUser(formData.value);
    resetForm();
    emit('created');
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
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    @click.self="handleClose"
  >
    <div class="w-full max-w-md rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-2xl font-bold text-slate-100">{{ t('user.title') }}</h2>
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
          <label for="name" class="block text-sm font-medium text-slate-300">{{ t('user.nameLabel') }} *</label>
          <input
            id="name"
            v-model="formData.name"
            type="text"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            placeholder="João Silva"
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-slate-300">{{ t('user.emailLabel') }} *</label>
          <input
            id="email"
            v-model="formData.email"
            type="email"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            placeholder="joao@example.com"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-slate-300">{{ t('user.passwordLabel') }} *</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            required
            minlength="6"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            placeholder="Mínimo 6 caracteres"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="points" class="block text-sm font-medium text-slate-300">{{ t('user.pointsLabel') }}</label>
            <input
              id="points"
              v-model.number="formData.points"
              type="number"
              min="0"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="role" class="block text-sm font-medium text-slate-300">{{ t('user.roleLabel') }}</label>
            <select
              id="role"
              v-model="formData.role"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            >
              <option value="user">{{ t('user.user') }}</option>
              <option value="admin">{{ t('user.admin') }}</option>
            </select>
          </div>
        </div>

        <div>
          <label for="status" class="block text-sm font-medium text-slate-300">{{ t('user.statusLabel') }}</label>
          <select
            id="status"
            v-model="formData.status"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          >
            <option value="active">{{ t('user.active') }}</option>
            <option value="inactive">{{ t('user.inactive') }}</option>
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
            {{ t('user.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50"
          >
            {{ loading ? t('user.creating') : t('user.create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


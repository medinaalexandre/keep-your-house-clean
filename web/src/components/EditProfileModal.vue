<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { updateUser, getUserById, getCurrentUser } from '../services/users';
import { logout } from '../services/auth';
import { useRouter } from 'vue-router';

const { t } = useI18n();
const router = useRouter();

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
  updated: [];
}>();

const currentUser = getCurrentUser();
const userEmail = ref('');
const formData = ref({
  name: '',
  password: '',
  confirmPassword: '',
});

const loading = ref(false);
const loadingUser = ref(false);
const error = ref('');
const showConfirmPassword = ref(false);

const loadUserData = async () => {
  if (!currentUser) return;
  
  loadingUser.value = true;
  error.value = '';
  
  try {
    const user = await getUserById(currentUser.id);
    formData.value.name = user.name;
    userEmail.value = user.email;
  } catch (e: any) {
    error.value = e.message;
  } finally {
    loadingUser.value = false;
  }
};

watch(() => formData.value.password, (newPassword) => {
  showConfirmPassword.value = newPassword.length > 0;
  if (newPassword.length === 0) {
    formData.value.confirmPassword = '';
  }
});

watch(() => formData.value.confirmPassword, () => {
  if (formData.value.confirmPassword && formData.value.password !== formData.value.confirmPassword) {
    error.value = t('profile.passwordsDoNotMatch');
  } else if (error.value === t('profile.passwordsDoNotMatch')) {
    error.value = '';
  }
});

watch(() => props.isOpen, (isOpen) => {
  if (isOpen && currentUser) {
    loadUserData();
  }
});

onMounted(() => {
  if (currentUser && props.isOpen) {
    loadUserData();
  }
});

const resetForm = () => {
  if (currentUser) {
    loadUserData();
  }
  formData.value.password = '';
  formData.value.confirmPassword = '';
  showConfirmPassword.value = false;
  error.value = '';
};

const handleClose = () => {
  resetForm();
  emit('close');
};

const handleSubmit = async () => {
  if (!formData.value.name) {
    error.value = t('profile.nameRequired');
    return;
  }

  if (formData.value.password && formData.value.password !== formData.value.confirmPassword) {
    error.value = t('profile.passwordsDoNotMatch');
    return;
  }

  if (!currentUser) {
    error.value = t('profile.userNotFound');
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    const updateData: { name: string; password?: string } = {
      name: formData.value.name,
    };

    if (formData.value.password) {
      updateData.password = formData.value.password;
    }

    await updateUser(currentUser.id, updateData);
    
    if (formData.value.password) {
      logout();
      router.push({ name: 'Login' });
    } else {
      localStorage.setItem('user_name', formData.value.name);
      resetForm();
      emit('updated');
      handleClose();
    }
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
        <h2 class="text-2xl font-bold text-slate-100">{{ t('profile.title') }}</h2>
        <button
          @click="handleClose"
          class="rounded-md p-2 text-slate-400 hover:bg-slate-800 hover:text-slate-100"
        >
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div v-if="loadingUser" class="py-8 text-center text-slate-400">
        {{ t('profile.loading') }}
      </div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="email" class="block text-sm font-medium text-slate-300">{{ t('profile.emailLabel') }}</label>
          <input
            id="email"
            :value="userEmail"
            type="email"
            disabled
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800/50 px-3 py-2 text-slate-400 cursor-not-allowed shadow-sm sm:text-sm"
          />
        </div>

        <div>
          <label for="name" class="block text-sm font-medium text-slate-300">{{ t('profile.nameLabel') }} *</label>
          <input
            id="name"
            v-model="formData.name"
            type="text"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('profile.namePlaceholder')"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-slate-300">{{ t('profile.passwordLabel') }}</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('profile.passwordPlaceholder')"
          />
          <p class="mt-1 text-xs text-slate-500">{{ t('profile.passwordHint') }}</p>
        </div>

        <div v-if="showConfirmPassword">
          <label for="confirmPassword" class="block text-sm font-medium text-slate-300">{{ t('profile.confirmPasswordLabel') }} *</label>
          <input
            id="confirmPassword"
            v-model="formData.confirmPassword"
            type="password"
            :required="showConfirmPassword"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('profile.confirmPasswordPlaceholder')"
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
            {{ t('profile.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50"
          >
            {{ loading ? t('profile.updating') : t('profile.update') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


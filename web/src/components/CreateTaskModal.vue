<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { createTask, type CreateTaskRequest } from '../services/tasks';

const { t } = useI18n();

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
  created: [];
}>();

const formData = ref<CreateTaskRequest>({
  title: '',
  description: '',
  points: 0,
  status: 'pending',
  scheduled_to: null,
  scheduled_by_id: null,
  frequency_value: 0,
  frequency_unit: 'days',
});

const loading = ref(false);
const error = ref('');

const resetForm = () => {
  formData.value = {
    title: '',
    description: '',
    points: 0,
    status: 'pending',
    scheduled_to: null,
    scheduled_by_id: null,
    frequency_value: 0,
    frequency_unit: 'days',
  };
  error.value = '';
};

const handleClose = () => {
  resetForm();
  emit('close');
};

const handleSubmit = async () => {
  if (!formData.value.title || !formData.value.description) {
    error.value = t('task.titleRequired');
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    const scheduledTo = formData.value.scheduled_to
      ? new Date(formData.value.scheduled_to + 'T00:00:00').toISOString()
      : null;

    await createTask({
      ...formData.value,
      scheduled_to: scheduledTo,
    });

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
    <div class="w-full max-w-2xl rounded-lg bg-slate-900 p-6 shadow-xl ring-1 ring-slate-800/50">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-2xl font-bold text-slate-100">{{ t('task.title') }}</h2>
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
          <label for="title" class="block text-sm font-medium text-slate-300">{{ t('task.titleLabel') }} *</label>
          <input
            id="title"
            v-model="formData.title"
            type="text"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('task.titlePlaceholder')"
          />
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-slate-300">{{ t('task.descriptionLabel') }} *</label>
          <textarea
            id="description"
            v-model="formData.description"
            required
            rows="3"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('task.descriptionPlaceholder')"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="points" class="block text-sm font-medium text-slate-300">{{ t('task.pointsLabel') }}</label>
            <input
              id="points"
              v-model.number="formData.points"
              type="number"
              min="0"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="status" class="block text-sm font-medium text-slate-300">{{ t('task.statusLabel') }}</label>
            <select
              id="status"
              v-model="formData.status"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            >
              <option value="pending">{{ t('task.pending') }}</option>
              <option value="in_progress">{{ t('task.inProgress') }}</option>
              <option value="completed">{{ t('task.completed') }}</option>
            </select>
          </div>
        </div>

        <div>
          <label for="scheduled_to" class="block text-sm font-medium text-slate-300">{{ t('task.scheduledToLabel') }}</label>
          <input
            id="scheduled_to"
            v-model="formData.scheduled_to"
            type="date"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="frequency_value" class="block text-sm font-medium text-slate-300">{{ t('task.frequencyValueLabel') }}</label>
            <input
              id="frequency_value"
              v-model.number="formData.frequency_value"
              type="number"
              min="0"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="frequency_unit" class="block text-sm font-medium text-slate-300">{{ t('task.frequencyUnitLabel') }}</label>
            <select
              id="frequency_unit"
              v-model="formData.frequency_unit"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            >
              <option value="days">{{ t('task.days') }}</option>
              <option value="weeks">{{ t('task.weeks') }}</option>
              <option value="months">{{ t('task.months') }}</option>
            </select>
          </div>
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
            {{ t('task.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50"
          >
            {{ loading ? t('task.creating') : t('task.create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


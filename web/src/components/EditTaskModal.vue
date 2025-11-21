<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { updateTask, type Task } from '../services/tasks';

const { t } = useI18n();

const props = defineProps<{
  isOpen: boolean;
  task: Task | null;
}>();

const emit = defineEmits<{
  close: [];
  updated: [];
}>();

const formData = ref({
  title: '',
  description: '',
  points: 0,
  status: 'pending',
  scheduled_to: null as string | null,
  frequency_value: 0,
  frequency_unit: 'days' as 'days' | 'weeks' | 'months',
});

const loading = ref(false);
const error = ref('');

watch(() => props.task, (newTask) => {
  if (newTask) {
    formData.value = {
      title: newTask.title,
      description: newTask.description,
      points: newTask.points,
      status: newTask.status,
      scheduled_to: newTask.scheduled_to ? newTask.scheduled_to.split('T')[0] : null,
      frequency_value: newTask.frequency_value,
      frequency_unit: newTask.frequency_unit as 'days' | 'weeks' | 'months',
    };
  }
}, { immediate: true });

const resetForm = () => {
  if (props.task) {
    formData.value = {
      title: props.task.title,
      description: props.task.description,
      points: props.task.points,
      status: props.task.status,
      scheduled_to: props.task.scheduled_to ? props.task.scheduled_to.split('T')[0] : null,
      frequency_value: props.task.frequency_value,
      frequency_unit: props.task.frequency_unit as 'days' | 'weeks' | 'months',
    };
  }
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

  if (!props.task) {
    error.value = t('task.taskNotFound');
    return;
  }

  if (props.task.completed) {
    error.value = t('task.cannotEditCompleted');
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    const scheduledTo = formData.value.scheduled_to
      ? new Date(formData.value.scheduled_to + 'T00:00:00').toISOString()
      : null;

    await updateTask(props.task.id, {
      title: formData.value.title,
      description: formData.value.description,
      points: formData.value.points,
      status: formData.value.status,
      scheduled_to: scheduledTo,
      frequency_value: formData.value.frequency_value,
      frequency_unit: formData.value.frequency_unit,
    });

    resetForm();
    emit('updated');
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
        <h2 class="text-2xl font-bold text-slate-100">{{ t('task.editTitle') }}</h2>
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
          <label for="edit-title" class="block text-sm font-medium text-slate-300">{{ t('task.titleLabel') }} *</label>
          <input
            id="edit-title"
            v-model="formData.title"
            type="text"
            required
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('task.titlePlaceholder')"
          />
        </div>

        <div>
          <label for="edit-description" class="block text-sm font-medium text-slate-300">{{ t('task.descriptionLabel') }} *</label>
          <textarea
            id="edit-description"
            v-model="formData.description"
            required
            rows="3"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            :placeholder="t('task.descriptionPlaceholder')"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="edit-points" class="block text-sm font-medium text-slate-300">{{ t('task.pointsLabel') }}</label>
            <input
              id="edit-points"
              v-model.number="formData.points"
              type="number"
              min="0"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="edit-status" class="block text-sm font-medium text-slate-300">{{ t('task.statusLabel') }}</label>
            <select
              id="edit-status"
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
          <label for="edit-scheduled_to" class="block text-sm font-medium text-slate-300">{{ t('task.scheduledToLabel') }}</label>
          <input
            id="edit-scheduled_to"
            v-model="formData.scheduled_to"
            type="date"
            class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="edit-frequency_value" class="block text-sm font-medium text-slate-300">{{ t('task.frequencyValueLabel') }}</label>
            <input
              id="edit-frequency_value"
              v-model.number="formData.frequency_value"
              type="number"
              min="0"
              class="mt-1 block w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-slate-100 placeholder-slate-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="edit-frequency_unit" class="block text-sm font-medium text-slate-300">{{ t('task.frequencyUnitLabel') }}</label>
            <select
              id="edit-frequency_unit"
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
            {{ loading ? t('task.updating') : t('task.update') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


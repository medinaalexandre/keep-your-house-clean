<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { markComplimentsAsViewed, type ComplimentWithUser } from '../services/compliments';

const { t } = useI18n();

const props = defineProps<{
  isOpen: boolean;
  compliments: ComplimentWithUser[];
}>();

const emit = defineEmits<{
  close: [];
  viewed: [];
}>();

const handleClose = async () => {
  if (props.compliments.length > 0) {
    try {
      const ids = props.compliments.map(c => c.id);
      await markComplimentsAsViewed(ids);
      emit('viewed');
    } catch (error) {
      console.error('Failed to mark compliments as viewed:', error);
    }
  }
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
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-2xl font-bold text-yellow-400">
          {{ t('unviewedCompliments.title', { count: compliments.length }) }}
        </h2>
        <button
          @click="handleClose"
          class="rounded-md p-2 text-slate-400 hover:bg-slate-800 hover:text-slate-100"
        >
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="space-y-4 max-h-[60vh] overflow-y-auto">
        <div
          v-for="(compliment, index) in compliments"
          :key="compliment.id"
          class="rounded-lg border border-yellow-500/50 bg-gradient-to-r from-yellow-500/10 to-amber-500/10 p-4"
        >
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0">
              <div class="flex h-8 w-8 items-center justify-center rounded-full bg-yellow-500/20 text-yellow-400 font-semibold">
                {{ index + 1 }}
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-semibold text-yellow-400">{{ compliment.title }}</h3>
              <p v-if="compliment.description" class="mt-1 text-sm text-slate-300">{{ compliment.description }}</p>
              <div class="mt-2 flex items-center gap-3 text-xs text-slate-400">
                <span>
                  {{ t('unviewedCompliments.from') }}
                  <span class="font-medium text-slate-300">{{ compliment.from_user_name || t('dashboard.unknownUser') }}</span>
                </span>
                <span v-if="compliment.points > 0" class="text-yellow-400 font-medium">
                  +{{ compliment.points }} {{ t('dashboard.points') }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6 flex justify-end">
        <button
          @click="handleClose"
          class="rounded-md bg-yellow-600 px-6 py-2 text-white transition-colors duration-200 hover:bg-yellow-500 focus:outline-none focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2 focus:ring-offset-slate-900"
        >
          {{ t('unviewedCompliments.close') }}
        </button>
      </div>
    </div>
  </div>
</template>


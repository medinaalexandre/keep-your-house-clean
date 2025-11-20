import { createI18n } from 'vue-i18n';
import ptBR from './locales/pt-BR.json';
import es from './locales/es.json';
import en from './locales/en.json';

const savedLocale = localStorage.getItem('locale') || 'pt-BR';

export const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'pt-BR',
  messages: {
    'pt-BR': ptBR,
    'es': es,
    'en': en,
  },
});

export const availableLocales = [
  { code: 'pt-BR', name: 'Português (Brasil)', flag: '🇧🇷' },
  { code: 'es', name: 'Español', flag: '🇪🇸' },
  { code: 'en', name: 'English (US)', flag: '🇺🇸' },
];


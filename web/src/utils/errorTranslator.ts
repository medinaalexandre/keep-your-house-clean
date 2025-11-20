import { i18n } from '../i18n';

const errorMap: Record<string, string> = {
  'invalid credentials': 'errors.invalidCredentials',
  'user account is inactive': 'errors.userAccountInactive',
  'domain already exists': 'errors.domainAlreadyExists',
  'email already exists': 'errors.emailAlreadyExists',
  'failed to hash password': 'errors.failedToHashPassword',
  'invalid request body': 'errors.invalidRequestBody',
  'email and password are required': 'errors.emailAndPasswordRequired',
  'all fields are required': 'errors.allFieldsRequired',
  'login failed': 'errors.loginFailed',
  'registration failed': 'errors.registrationFailed',
  'unknown error': 'errors.unknownError',
};

export const translateError = (errorMessage: string): string => {
  const normalizedMessage = errorMessage.trim().toLowerCase();
  const translationKey = errorMap[normalizedMessage];
  
  if (translationKey) {
    return i18n.global.t(translationKey);
  }
  
  for (const [key, value] of Object.entries(errorMap)) {
    if (normalizedMessage.includes(key.toLowerCase())) {
      return i18n.global.t(value);
    }
  }
  
  return errorMessage;
};


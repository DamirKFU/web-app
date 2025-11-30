import { API_CONFIG } from '@/config/api';

const API_BASE_URL = API_CONFIG.baseURL;

// Получаем CSRF токен из localStorage
export async function getCSRFToken(): Promise<string | undefined> {
  return localStorage.getItem("csrf_token") || undefined;
}

// Функция для добавления CSRF токена в заголовки запроса
export async function addCSRFToken(): Promise<HeadersInit> {
  const csrfToken = localStorage.getItem("csrf_token");

  return {
    "Content-Type": "application/json",
    "X-CSRF-TOKEN": csrfToken || "",   // Django принимает любой вариант регистра
  };
}

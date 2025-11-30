import { API_CONFIG } from '@/config/api';
import { ProductResponse, CategoriesResponse, ColorsResponse } from '@/types/api';

// Получаем CSRF-токен из localStorage
const getCsrfToken = (): string | null => {
  return localStorage.getItem("csrf_token");
};

// Генерация headers с CSRF
const getHeaders = (): HeadersInit => {
  const csrfToken = getCsrfToken();

  return {
    "X-CSRF-TOKEN": csrfToken || "",
    "Content-Type": "application/json",
  };
};

interface FetchProductsParams {
  category?: string;
  color?: string;
  size?: string;
  page?: string;
}

export const fetchProducts = async (params?: FetchProductsParams): Promise<ProductResponse> => {
  const queryParams = new URLSearchParams();
  if (params?.category && params.category !== '_all') queryParams.append('category', params.category);
  if (params?.color && params.color !== '_all') queryParams.append('color', params.color);
  if (params?.size && params.size !== '_all') queryParams.append('size', params.size);
  if (params?.page) queryParams.append('page', params.page);

  const url = `${API_CONFIG.baseURL}${API_CONFIG.endpoints.products}${
    queryParams.toString() ? `?${queryParams.toString()}` : ''
  }`;

  const response = await fetch(url, {
    credentials: "include",
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error("Failed to fetch products");
  }

  return response.json();
};

export const fetchCategories = async (): Promise<CategoriesResponse> => {
  const response = await fetch(`${API_CONFIG.baseURL}${API_CONFIG.endpoints.categories}`, {
    credentials: "include",
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error("Failed to fetch categories");
  }

  return response.json();
};

export const fetchColors = async (): Promise<ColorsResponse> => {
  const response = await fetch(`${API_CONFIG.baseURL}${API_CONFIG.endpoints.colors}`, {
    credentials: "include",
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error("Failed to fetch colors");
  }

  return response.json();
};

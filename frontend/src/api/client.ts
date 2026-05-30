import axios, { type AxiosRequestConfig } from 'axios';
import { z } from 'zod';
import { ApiError, triggerUnauthorized } from '@/lib/apiError';
import { getAccessToken } from '@/lib/authStorage';
import { ApiEnvelopeSchema } from '@/schemas/apiSchema';

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '',
  headers: {
    Accept: 'application/json',
  },
});

apiClient.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      triggerUnauthorized();
    }
    const message =
      error.response?.data?.error ??
      error.response?.data?.message ??
      error.message ??
      'Request failed';
    return Promise.reject(
      new ApiError(message, error.response?.status ?? 500),
    );
  },
);

export async function apiRequest<T>(
  config: AxiosRequestConfig,
  dataSchema: z.ZodType<T>,
): Promise<T> {
  const response = await apiClient.request(config);
  const contentType = response.headers['content-type'] as string | undefined;

  if (contentType?.includes('application/json')) {
    const envelope = ApiEnvelopeSchema.parse(response.data);
    if (!envelope.success) {
      if (response.status === 401) triggerUnauthorized();
      throw new ApiError(envelope.error ?? 'Request failed', response.status);
    }
    return dataSchema.parse(envelope.data);
  }

  return dataSchema.parse(response.data);
}

export async function apiBlobRequest(config: AxiosRequestConfig): Promise<Blob> {
  const response = await apiClient.request({
    ...config,
    responseType: 'blob',
  });
  return response.data as Blob;
}

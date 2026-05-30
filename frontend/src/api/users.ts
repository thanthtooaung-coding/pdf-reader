import { apiRequest } from '@/api/client';
import { UserResponseSchema, type UserResponse } from '@/schemas/apiSchema';

export async function getMe(): Promise<UserResponse> {
  return apiRequest(
    { method: 'GET', url: '/api/v1/users/me' },
    UserResponseSchema,
  );
}

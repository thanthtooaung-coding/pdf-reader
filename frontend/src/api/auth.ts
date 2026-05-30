import { apiRequest } from '@/api/client';
import {
  AuthResponseSchema,
  LoginRequestSchema,
  OTPSentResponseSchema,
  RegisterRequestSchema,
  ResendOTPRequestSchema,
  VerifyOTPRequestSchema,
  type AuthResponse,
  type LoginRequest,
  type OTPSentResponse,
  type RegisterRequest,
  type ResendOTPRequest,
  type VerifyOTPRequest,
} from '@/schemas/apiSchema';

export async function registerRequest(
  body: RegisterRequest,
): Promise<OTPSentResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/auth/register/request',
      data: RegisterRequestSchema.parse(body),
    },
    OTPSentResponseSchema,
  );
}

export async function registerVerify(
  body: VerifyOTPRequest,
): Promise<AuthResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/auth/register/verify',
      data: VerifyOTPRequestSchema.parse(body),
    },
    AuthResponseSchema,
  );
}

export async function registerResend(
  body: ResendOTPRequest,
): Promise<OTPSentResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/auth/register/resend',
      data: ResendOTPRequestSchema.parse(body),
    },
    OTPSentResponseSchema,
  );
}

export async function loginRequest(body: LoginRequest): Promise<OTPSentResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/auth/login',
      data: LoginRequestSchema.parse(body),
    },
    OTPSentResponseSchema,
  );
}

export async function loginVerify(body: VerifyOTPRequest): Promise<AuthResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/auth/login/verify',
      data: VerifyOTPRequestSchema.parse(body),
    },
    AuthResponseSchema,
  );
}

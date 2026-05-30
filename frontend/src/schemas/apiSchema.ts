import { z } from 'zod';

const uuid = z.string().uuid();
const dateString = z.string();

export const ApiEnvelopeSchema = z.object({
  success: z.boolean(),
  message: z.string().optional(),
  data: z.unknown().optional(),
  error: z.string().optional(),
});

export type ApiEnvelope = z.infer<typeof ApiEnvelopeSchema>;

export const UserResponseSchema = z.object({
  id: uuid,
  role_id: uuid,
  role_code: z.string().optional(),
  fullname: z.string().optional(),
  username: z.string(),
  email: z.string().email(),
  is_enable: z.boolean(),
  created_at: dateString,
  updated_at: dateString,
});

export const AuthResponseSchema = z.object({
  access_token: z.string(),
  token_type: z.string(),
  expires_in: z.number(),
  user: UserResponseSchema,
});

export const OTPSentResponseSchema = z.object({
  message: z.string(),
  expires_in_seconds: z.number(),
  otp: z.string().optional(),
});

export const RegisterRequestSchema = z.object({
  fullname: z.string().min(1).max(128),
  username: z.string().min(2).max(64),
  email: z.string().email().max(255),
  password: z.string().min(8).max(72),
});

export const LoginRequestSchema = z.object({
  email: z.string().email().max(255),
  username: z.string().min(2).max(64),
  password: z.string().min(8).max(72),
});

export const VerifyOTPRequestSchema = z.object({
  email: z.string().email().max(255),
  otp: z.string().length(6),
});

export const ResendOTPRequestSchema = z.object({
  email: z.string().email().max(255),
});

export const CreateWorkspaceRequestSchema = z.object({
  name: z.string().min(1).max(255),
});

export const WorkspaceResponseSchema = z.object({
  id: uuid,
  user_id: uuid,
  name: z.string(),
  file_count: z.number().optional(),
  created_at: dateString,
  updated_at: dateString,
});

export const PageMetaSchema = z.object({
  page: z.number(),
  limit: z.number(),
  total: z.number(),
  total_pages: z.number(),
});

export const PagedWorkspacesSchema = z.object({
  items: z.array(WorkspaceResponseSchema),
  meta: PageMetaSchema,
});

export const FileResponseSchema = z.object({
  id: uuid,
  user_id: uuid,
  workspace_id: uuid,
  original_file_name: z.string(),
  name: z.string(),
  url: z.string(),
  type: z.string(),
  created_at: dateString,
  updated_at: dateString,
});

export const CreateCommentRequestSchema = z.object({
  message: z.string().min(1).max(5000),
});

export const CommentResponseSchema = z.object({
  id: uuid,
  file_id: uuid,
  user_id: uuid,
  username: z.string().optional(),
  message: z.string(),
  created_at: dateString,
});

export const AIJobTypeSchema = z.enum(['TRANSLATE', 'SUMMARIZE', 'COMMENT']);
export const AIJobStatusSchema = z.enum([
  'pending',
  'processing',
  'completed',
  'failed',
]);

export const CreateAIJobRequestSchema = z.object({
  file_id: uuid,
  type: AIJobTypeSchema,
  input: z.string().max(10000).optional(),
});

export const AIJobResponseSchema = z.object({
  id: uuid,
  user_id: uuid,
  workspace_id: uuid,
  file_id: uuid,
  type: z.string(),
  status: AIJobStatusSchema,
  input: z.string().optional(),
  output: z.string().optional(),
  duration_ms: z.number(),
  retry_count: z.number(),
  created_at: dateString,
  updated_at: dateString,
});

export const CommentFormSchema = z.object({
  message: z.string().min(1, 'Comment cannot be empty'),
});

export type UserResponse = z.infer<typeof UserResponseSchema>;
export type AuthResponse = z.infer<typeof AuthResponseSchema>;
export type OTPSentResponse = z.infer<typeof OTPSentResponseSchema>;
export type RegisterRequest = z.infer<typeof RegisterRequestSchema>;
export type LoginRequest = z.infer<typeof LoginRequestSchema>;
export type VerifyOTPRequest = z.infer<typeof VerifyOTPRequestSchema>;
export type ResendOTPRequest = z.infer<typeof ResendOTPRequestSchema>;
export type CreateWorkspaceRequest = z.infer<typeof CreateWorkspaceRequestSchema>;
export type WorkspaceResponse = z.infer<typeof WorkspaceResponseSchema>;
export type PagedWorkspaces = z.infer<typeof PagedWorkspacesSchema>;
export type FileResponse = z.infer<typeof FileResponseSchema>;
export type CommentResponse = z.infer<typeof CommentResponseSchema>;
export type CreateCommentRequest = z.infer<typeof CreateCommentRequestSchema>;
export type CreateAIJobRequest = z.infer<typeof CreateAIJobRequestSchema>;
export type AIJobResponse = z.infer<typeof AIJobResponseSchema>;
export type CommentFormValues = z.infer<typeof CommentFormSchema>;

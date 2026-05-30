import { apiRequest } from '@/api/client';
import {
  CommentResponseSchema,
  CreateCommentRequestSchema,
  type CommentResponse,
  type CreateCommentRequest,
} from '@/schemas/apiSchema';
import { z } from 'zod';

const CommentListSchema = z.array(CommentResponseSchema);

export async function listComments(fileId: string): Promise<CommentResponse[]> {
  return apiRequest(
    { method: 'GET', url: `/api/v1/files/${fileId}/comments` },
    CommentListSchema,
  );
}

export async function createComment(
  fileId: string,
  body: CreateCommentRequest,
): Promise<CommentResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: `/api/v1/files/${fileId}/comments`,
      data: CreateCommentRequestSchema.parse(body),
    },
    CommentResponseSchema,
  );
}

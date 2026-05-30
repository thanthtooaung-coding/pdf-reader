import { apiRequest } from '@/api/client';
import {
  AIJobResponseSchema,
  CreateAIJobRequestSchema,
  type AIJobResponse,
  type CreateAIJobRequest,
} from '@/schemas/apiSchema';
import { z } from 'zod';

const AIJobListSchema = z.array(AIJobResponseSchema);

export async function createAIJob(
  body: CreateAIJobRequest,
): Promise<AIJobResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/ai-jobs',
      data: CreateAIJobRequestSchema.parse(body),
    },
    AIJobResponseSchema,
  );
}

export async function getAIJob(jobId: string): Promise<AIJobResponse> {
  return apiRequest(
    { method: 'GET', url: `/api/v1/ai-jobs/${jobId}` },
    AIJobResponseSchema,
  );
}

export async function listAIJobsByWorkspace(
  workspaceId: string,
): Promise<AIJobResponse[]> {
  return apiRequest(
    { method: 'GET', url: `/api/v1/workspaces/${workspaceId}/ai-jobs` },
    AIJobListSchema,
  );
}

export function isAIJobRunning(status: string): boolean {
  return status === 'pending' || status === 'processing';
}

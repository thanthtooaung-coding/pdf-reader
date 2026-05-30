import { apiRequest } from '@/api/client';
import {
  CreateWorkspaceRequestSchema,
  PagedWorkspacesSchema,
  WorkspaceResponseSchema,
  type CreateWorkspaceRequest,
  type PagedWorkspaces,
  type WorkspaceResponse,
} from '@/schemas/apiSchema';
import { z } from 'zod';

export async function listWorkspaces(
  page = 1,
  limit = 20,
): Promise<PagedWorkspaces> {
  return apiRequest(
    {
      method: 'GET',
      url: '/api/v1/workspaces',
      params: { page, limit },
    },
    PagedWorkspacesSchema,
  );
}

export async function createWorkspace(
  body: CreateWorkspaceRequest,
): Promise<WorkspaceResponse> {
  return apiRequest(
    {
      method: 'POST',
      url: '/api/v1/workspaces',
      data: CreateWorkspaceRequestSchema.parse(body),
    },
    WorkspaceResponseSchema,
  );
}

export async function getWorkspace(id: string): Promise<WorkspaceResponse> {
  return apiRequest(
    { method: 'GET', url: `/api/v1/workspaces/${id}` },
    WorkspaceResponseSchema,
  );
}

export { WorkspaceResponseSchema };

export const WorkspaceListSchema = z.array(WorkspaceResponseSchema);

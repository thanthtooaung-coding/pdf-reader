import { apiBlobRequest, apiRequest } from '@/api/client';
import {
  FileResponseSchema,
  type FileResponse,
} from '@/schemas/apiSchema';
import { z } from 'zod';

const FileListSchema = z.array(FileResponseSchema);

export async function listFilesByWorkspace(
  workspaceId: string,
): Promise<FileResponse[]> {
  return apiRequest(
    { method: 'GET', url: `/api/v1/workspaces/${workspaceId}/files` },
    FileListSchema,
  );
}

export async function getFile(fileId: string): Promise<FileResponse> {
  return apiRequest(
    { method: 'GET', url: `/api/v1/files/${fileId}` },
    FileResponseSchema,
  );
}

export async function uploadFile(
  workspaceId: string,
  file: File,
): Promise<FileResponse> {
  const formData = new FormData();
  formData.append('file', file);

  return apiRequest(
    {
      method: 'POST',
      url: `/api/v1/workspaces/${workspaceId}/files`,
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' },
    },
    FileResponseSchema,
  );
}

export async function downloadFile(fileId: string): Promise<Blob> {
  return apiBlobRequest({
    method: 'GET',
    url: `/api/v1/files/${fileId}/download`,
  });
}

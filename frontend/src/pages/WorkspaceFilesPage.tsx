import { useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { listFilesByWorkspace, uploadFile } from '@/api/files';
import { getWorkspace } from '@/api/workspaces';
import AppHeader from '@/components/layout/AppHeader';
import Button from '@/components/ui/Button';
import { ApiError } from '@/lib/apiError';

export default function WorkspaceFilesPage() {
  const { workspaceId = '' } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);

  const workspaceQuery = useQuery({
    queryKey: ['workspace', workspaceId],
    queryFn: () => getWorkspace(workspaceId),
    enabled: Boolean(workspaceId),
  });

  const filesQuery = useQuery({
    queryKey: ['files', workspaceId],
    queryFn: () => listFilesByWorkspace(workspaceId),
    enabled: Boolean(workspaceId),
  });

  const uploadMutation = useMutation({
    mutationFn: (file: File) => uploadFile(workspaceId, file),
    onSuccess: (file) => {
      setUploadError(null);
      queryClient.invalidateQueries({ queryKey: ['files', workspaceId] });
      queryClient.invalidateQueries({ queryKey: ['workspace', workspaceId] });
      navigate(`/workspaces/${workspaceId}/files/${file.id}`);
    },
    onError: (err: ApiError) => setUploadError(err.message),
  });

  const handleFile = (file: File | null) => {
    if (!file || file.type !== 'application/pdf') {
      setUploadError('Please upload a PDF file.');
      return;
    }
    uploadMutation.mutate(file);
  };

  const workspaceName = workspaceQuery.data?.name ?? 'Workspace';

  return (
    <div className="min-h-screen bg-slate-50">
      <AppHeader
        title={workspaceName}
        breadcrumbs={[
          { label: 'Workspaces', to: '/' },
          { label: workspaceName },
        ]}
      />
      <main className="mx-auto max-w-3xl p-6">
        <label
          className={`mb-8 flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed bg-white p-8 transition-all duration-200 ${
            isDragging
              ? 'border-slate-400 bg-slate-100'
              : 'border-slate-300 hover:border-slate-400'
          }`}
          onDragOver={(event) => {
            event.preventDefault();
            setIsDragging(true);
          }}
          onDragLeave={() => setIsDragging(false)}
          onDrop={(event) => {
            event.preventDefault();
            setIsDragging(false);
            handleFile(event.dataTransfer.files[0] ?? null);
          }}
        >
          <input
            type="file"
            accept="application/pdf"
            className="hidden"
            onChange={(event) =>
              handleFile(event.target.files?.[0] ?? null)
            }
            disabled={uploadMutation.isPending}
          />
          <p className="text-base font-medium text-slate-700">
            {uploadMutation.isPending
              ? 'Uploading...'
              : 'Drop a PDF here or click to upload'}
          </p>
          <p className="mt-1 text-sm text-slate-500">Max 25 MB</p>
        </label>

        {uploadError && (
          <p className="mb-4 text-sm text-red-600">{uploadError}</p>
        )}

        {filesQuery.isLoading && (
          <p className="text-sm text-slate-500">Loading files...</p>
        )}

        {filesQuery.data && filesQuery.data.length === 0 && (
          <p className="text-sm text-slate-500">
            No files in this workspace yet.
          </p>
        )}

        <ul className="space-y-3">
          {filesQuery.data?.map((file) => (
            <li key={file.id}>
              <Link
                to={`/workspaces/${workspaceId}/files/${file.id}`}
                className="flex items-center justify-between rounded-lg border border-slate-200 bg-white p-4 transition-all duration-200 hover:border-slate-300 hover:shadow-sm"
              >
                <div>
                  <p className="font-medium text-slate-900">
                    {file.original_file_name}
                  </p>
                  <p className="mt-1 text-xs text-slate-500">
                    Uploaded {new Date(file.created_at).toLocaleString()}
                  </p>
                </div>
                <Button variant="ghost">Open</Button>
              </Link>
            </li>
          ))}
        </ul>
      </main>
    </div>
  );
}

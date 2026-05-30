import { useQuery } from '@tanstack/react-query';
import { Link, useParams } from 'react-router-dom';
import { getFile } from '@/api/files';
import { getWorkspace } from '@/api/workspaces';
import Button from '@/components/ui/Button';
import SplitLayout from '@/components/workspace/SplitLayout';
import { useAuth } from '@/context/AuthProvider';

export default function ReaderPage() {
  const { workspaceId = '', fileId = '' } = useParams();
  const { user, logout } = useAuth();

  const workspaceQuery = useQuery({
    queryKey: ['workspace', workspaceId],
    queryFn: () => getWorkspace(workspaceId),
    enabled: Boolean(workspaceId),
  });

  const fileQuery = useQuery({
    queryKey: ['file', fileId],
    queryFn: () => getFile(fileId),
    enabled: Boolean(fileId),
  });

  const workspaceName = workspaceQuery.data?.name ?? 'Workspace';
  const fileName = fileQuery.data?.original_file_name ?? 'Document';

  if (fileQuery.isError) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4">
        <p className="text-sm text-red-600">File not found or access denied.</p>
        <Link to={`/workspaces/${workspaceId}`}>
          <Button variant="ghost">Back to files</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="flex h-screen flex-col">
      <header className="flex shrink-0 items-center justify-between border-b border-slate-200 bg-white px-4 py-2">
        <nav className="flex flex-wrap items-center gap-1 text-xs text-slate-500">
          <Link to="/" className="hover:text-slate-800">
            Workspaces
          </Link>
          <span>/</span>
          <Link
            to={`/workspaces/${workspaceId}`}
            className="hover:text-slate-800"
          >
            {workspaceName}
          </Link>
          <span>/</span>
          <span className="text-slate-700">{fileName}</span>
        </nav>
        <div className="flex items-center gap-3">
          {user && (
            <span className="text-sm text-slate-600">{user.username}</span>
          )}
          <Link to={`/workspaces/${workspaceId}`}>
            <Button variant="ghost">Back</Button>
          </Link>
          <Button variant="ghost" onClick={logout}>
            Logout
          </Button>
        </div>
      </header>

      <div className="min-h-0 flex-1">
        {fileQuery.data && (
          <SplitLayout
            workspaceId={workspaceId}
            fileId={fileId}
            fileName={fileQuery.data.original_file_name}
          />
        )}
        {fileQuery.isLoading && (
          <div className="flex h-full items-center justify-center text-sm text-slate-500">
            Loading document...
          </div>
        )}
      </div>
    </div>
  );
}

import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createWorkspace, listWorkspaces } from '@/api/workspaces';
import AppHeader from '@/components/layout/AppHeader';
import Button from '@/components/ui/Button';
import { ApiError } from '@/lib/apiError';

export default function WorkspaceListPage() {
  const queryClient = useQueryClient();
  const [name, setName] = useState('');
  const [error, setError] = useState<string | null>(null);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['workspaces'],
    queryFn: () => listWorkspaces(),
  });

  const createMutation = useMutation({
    mutationFn: createWorkspace,
    onSuccess: () => {
      setName('');
      setError(null);
      queryClient.invalidateQueries({ queryKey: ['workspaces'] });
    },
    onError: (err: ApiError) => setError(err.message),
  });

  const handleCreate = (event: React.FormEvent) => {
    event.preventDefault();
    if (!name.trim()) return;
    createMutation.mutate({ name: name.trim() });
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <AppHeader title="Workspaces" />
      <main className="mx-auto max-w-3xl p-6">
        <form
          onSubmit={handleCreate}
          className="mb-8 flex gap-2 rounded-lg border border-slate-200 bg-white p-4"
        >
          <input
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="New workspace name"
            className="flex-1 rounded-md border border-slate-300 px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
          />
          <Button type="submit" disabled={createMutation.isPending}>
            {createMutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </form>

        {error && <p className="mb-4 text-sm text-red-600">{error}</p>}

        {isLoading && (
          <p className="text-sm text-slate-500">Loading workspaces...</p>
        )}
        {isError && (
          <p className="text-sm text-red-600">Failed to load workspaces.</p>
        )}

        {data && data.items.length === 0 && (
          <p className="text-sm text-slate-500">
            No workspaces yet. Create one to upload PDFs.
          </p>
        )}

        <ul className="space-y-3">
          {data?.items.map((workspace) => (
            <li key={workspace.id}>
              <Link
                to={`/workspaces/${workspace.id}`}
                className="block rounded-lg border border-slate-200 bg-white p-4 transition-all duration-200 hover:border-slate-300 hover:shadow-sm"
              >
                <p className="font-medium text-slate-900">{workspace.name}</p>
                <p className="mt-1 text-xs text-slate-500">
                  {workspace.file_count ?? 0} files
                </p>
              </Link>
            </li>
          ))}
        </ul>
      </main>
    </div>
  );
}

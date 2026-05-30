import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  createAIJob,
  getAIJob,
  isAIJobRunning,
  listAIJobsByWorkspace,
} from '@/api/aiJobs';
import type { AIJobResponse, CreateAIJobRequest } from '@/schemas/apiSchema';

export function useCreateAIJob() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: CreateAIJobRequest) => createAIJob(body),
    onSuccess: (job) => {
      queryClient.invalidateQueries({
        queryKey: ['ai-jobs', job.workspace_id],
      });
    },
  });
}

export function useAIJobPoll(jobId: string | null) {
  return useQuery({
    queryKey: ['ai-job', jobId],
    queryFn: () => getAIJob(jobId!),
    enabled: Boolean(jobId),
    refetchInterval: (query) => {
      const status = query.state.data?.status;
      if (!status || !isAIJobRunning(status)) return false;
      return 500;
    },
  });
}

export function useWorkspaceAIJobs(workspaceId: string, fileId: string) {
  return useQuery({
    queryKey: ['ai-jobs', workspaceId, fileId],
    queryFn: async () => {
      const jobs = await listAIJobsByWorkspace(workspaceId);
      return jobs.filter((job) => job.file_id === fileId);
    },
    enabled: Boolean(workspaceId && fileId),
  });
}

export function pickLatestCompletedJob(
  jobs: AIJobResponse[],
  type: string,
): AIJobResponse | null {
  const completed = jobs
    .filter((job) => job.type === type && job.status === 'completed')
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    );
  return completed[0] ?? null;
}

export function pickCompletedSummaries(jobs: AIJobResponse[]): string[] {
  return jobs
    .filter((job) => job.type === 'SUMMARIZE' && job.status === 'completed')
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )
    .map((job) => job.output ?? '')
    .filter(Boolean);
}

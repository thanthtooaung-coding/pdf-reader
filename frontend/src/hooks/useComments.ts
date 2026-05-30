import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { createComment, listComments } from '@/api/comments';
import type { CreateCommentRequest } from '@/schemas/apiSchema';

export function useComments(fileId: string) {
  return useQuery({
    queryKey: ['comments', fileId],
    queryFn: () => listComments(fileId),
    enabled: Boolean(fileId),
  });
}

export function useCreateComment(fileId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: CreateCommentRequest) => createComment(fileId, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', fileId] });
    },
  });
}

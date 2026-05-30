import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { AnimatePresence, motion } from 'framer-motion';
import Button from '@/components/ui/Button';
import { CommentFormSchema, type CommentFormValues } from '@/schemas/apiSchema';
import { useComments, useCreateComment } from '@/hooks/useComments';
import { fadeSlide } from '@/lib/motion';
import { ApiError } from '@/lib/apiError';

interface NotesPanelProps {
  fileId: string;
  draftSourceText: string | null;
  onDraftConsumed: () => void;
}

export default function NotesPanel({
  fileId,
  draftSourceText,
  onDraftConsumed,
}: NotesPanelProps) {
  const { data: comments = [], isLoading, isError } = useComments(fileId);
  const createComment = useCreateComment(fileId);
  const [quotedSource, setQuotedSource] = useState<string | undefined>();
  const [submitError, setSubmitError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<CommentFormValues>({
    resolver: zodResolver(CommentFormSchema),
    defaultValues: { message: '' },
  });

  useEffect(() => {
    if (!draftSourceText) return;
    setValue('message', draftSourceText);
    setQuotedSource(draftSourceText);
    onDraftConsumed();
  }, [draftSourceText, onDraftConsumed, setValue]);

  const onSubmit = handleSubmit(async (values) => {
    setSubmitError(null);
    try {
      await createComment.mutateAsync({ message: values.message });
      setQuotedSource(undefined);
      reset({ message: '' });
    } catch (error) {
      setSubmitError(
        error instanceof ApiError ? error.message : 'Failed to save comment',
      );
    }
  });

  return (
    <article className="flex h-full flex-col">
      <header className="border-b border-slate-200 px-4 py-3">
        <h2 className="text-sm font-semibold text-slate-900">My Notes</h2>
        <p className="mt-1 text-xs text-slate-500">
          Comments synced to this file on the server
        </p>
      </header>

      <motion.form
        onSubmit={onSubmit}
        className="space-y-3 border-b border-slate-200 p-4"
        layout
      >
        {quotedSource && (
          <p className="rounded bg-slate-50 p-2 text-xs italic text-slate-500">
            From selection: &ldquo;{quotedSource}&rdquo;
          </p>
        )}
        <label className="block">
          <span className="mb-1 block text-xs font-medium text-slate-600">
            New comment
          </span>
          <textarea
            {...register('message')}
            rows={4}
            placeholder="Write a note..."
            className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm text-slate-800 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
          />
        </label>
        {errors.message && (
          <motion.p
            className="text-xs text-red-600"
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
          >
            {errors.message.message}
          </motion.p>
        )}
        {submitError && <p className="text-xs text-red-600">{submitError}</p>}
        <Button type="submit" disabled={isSubmitting || createComment.isPending}>
          Save Comment
        </Button>
      </motion.form>

      <div className="flex-1 overflow-auto p-4">
        {isLoading && (
          <p className="text-sm text-slate-500">Loading comments...</p>
        )}
        {isError && (
          <p className="text-sm text-red-600">Failed to load comments.</p>
        )}
        {!isLoading && comments.length === 0 && (
          <p className="text-sm text-slate-400">
            No comments yet. Select text and click Take Note to get started.
          </p>
        )}
        {!isLoading && comments.length > 0 && (
          <motion.ul layout className="space-y-3">
            <AnimatePresence initial={false}>
              {comments.map((comment) => (
                <motion.li
                  key={comment.id}
                  layout
                  {...fadeSlide}
                  className="rounded-md border border-slate-200 bg-white p-3"
                >
                  <div className="mb-2 flex items-start justify-between gap-2">
                    <time
                      dateTime={comment.created_at}
                      className="text-xs text-slate-400"
                    >
                      {new Date(comment.created_at).toLocaleString()}
                    </time>
                    {comment.username && (
                      <span className="text-xs text-slate-500">
                        {comment.username}
                      </span>
                    )}
                  </div>
                  <pre className="whitespace-pre-wrap font-sans text-sm text-slate-800">
                    {comment.message}
                  </pre>
                </motion.li>
              ))}
            </AnimatePresence>
          </motion.ul>
        )}
      </div>
    </article>
  );
}

import { useEffect, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import Button from '@/components/ui/Button';
import {
  pickCompletedSummaries,
  pickLatestCompletedJob,
  useAIJobPoll,
  useCreateAIJob,
  useWorkspaceAIJobs,
} from '@/hooks/useAIJob';
import { isAIJobRunning } from '@/api/aiJobs';

interface InsightsPanelProps {
  translation: string | null;
  summaries: string[];
  isLoading: boolean;
  isError: boolean;
  jobStatus?: string;
  onRetry: () => void;
}

export default function InsightsPanel({
  translation,
  summaries,
  isLoading,
  isError,
  jobStatus,
  onRetry,
}: InsightsPanelProps) {
  const [expandedIndex, setExpandedIndex] = useState<number | null>(0);

  return (
    <article className="flex h-full flex-col">
      <header className="border-b border-slate-200 px-4 py-3">
        <h2 className="text-sm font-semibold text-slate-900">AI Insights</h2>
        <p className="mt-1 text-xs text-slate-500">
          Document-level translations and TL;DR summaries
        </p>
      </header>

      <div className="flex-1 space-y-4 overflow-auto p-4">
        <AnimatePresence mode="wait">
          {isLoading && (
            <motion.p
              key="loading"
              className="text-sm text-slate-500"
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
            >
              {jobStatus === 'processing'
                ? 'Processing document...'
                : 'Starting AI job...'}
            </motion.p>
          )}
        </AnimatePresence>

        <AnimatePresence>
          {isError && (
            <motion.div
              key="error"
              className="rounded-md border border-red-200 bg-red-50 p-3"
              initial={{ opacity: 0, scale: 0.98 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.98 }}
            >
              <p className="text-sm text-red-700">
                AI job failed. Please try again.
              </p>
              <Button variant="ghost" className="mt-2" onClick={onRetry}>
                Retry
              </Button>
            </motion.div>
          )}
        </AnimatePresence>

        <section>
          <h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
            Translation
          </h3>
          <AnimatePresence mode="wait">
            {translation ? (
              <motion.p
                key={translation}
                className="rounded-md border border-slate-200 bg-white p-3 text-sm leading-relaxed text-slate-800"
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -8 }}
              >
                {translation}
              </motion.p>
            ) : (
              <motion.p
                key="empty-translation"
                className="text-sm text-slate-400"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
              >
                Click Translate doc to translate the entire PDF.
              </motion.p>
            )}
          </AnimatePresence>
        </section>

        <section>
          <h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
            Layered TL;DR
          </h3>
          {summaries.length === 0 ? (
            <p className="text-sm text-slate-400">
              Click Summarize doc to build your summary stack.
            </p>
          ) : (
            <motion.ul layout className="space-y-2">
              <AnimatePresence initial={false}>
                {summaries.map((summary, index) => {
                  const isExpanded = expandedIndex === index;
                  return (
                    <motion.li
                      key={`${summary.slice(0, 24)}-${index}`}
                      layout
                      initial={{ opacity: 0, y: -12 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="overflow-hidden rounded-md border border-slate-200 bg-white"
                    >
                      <button
                        type="button"
                        className="flex w-full items-center justify-between px-3 py-2 text-left text-sm font-medium text-slate-800 transition-colors duration-200 hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300"
                        onClick={() =>
                          setExpandedIndex(isExpanded ? null : index)
                        }
                      >
                        <span>Summary {summaries.length - index}</span>
                        <motion.span
                          animate={{ rotate: isExpanded ? 180 : 0 }}
                          transition={{ duration: 0.2 }}
                          className="text-xs text-slate-400"
                        >
                          {isExpanded ? 'Hide' : 'Show'}
                        </motion.span>
                      </button>
                      <AnimatePresence initial={false}>
                        {isExpanded && (
                          <motion.div
                            initial={{ height: 0, opacity: 0 }}
                            animate={{ height: 'auto', opacity: 1 }}
                            exit={{ height: 0, opacity: 0 }}
                            transition={{ duration: 0.2 }}
                            className="overflow-hidden"
                          >
                            <p className="border-t border-slate-100 px-3 py-2 text-sm leading-relaxed text-slate-700">
                              {summary}
                            </p>
                          </motion.div>
                        )}
                      </AnimatePresence>
                    </motion.li>
                  );
                })}
              </AnimatePresence>
            </motion.ul>
          )}
        </section>
      </div>
    </article>
  );
}

export function useInsightsState(workspaceId: string, fileId: string) {
  const createJob = useCreateAIJob();
  const [activeJobId, setActiveJobId] = useState<string | null>(null);
  const [lastAction, setLastAction] = useState<'TRANSLATE' | 'SUMMARIZE' | null>(
    null,
  );

  const pollQuery = useAIJobPoll(activeJobId);
  const jobsQuery = useWorkspaceAIJobs(workspaceId, fileId);

  const jobs = jobsQuery.data ?? [];
  const translation = pickLatestCompletedJob(jobs, 'TRANSLATE')?.output ?? null;
  const summaries = pickCompletedSummaries(jobs);

  useEffect(() => {
    if (!pollQuery.data) return;
    if (!isAIJobRunning(pollQuery.data.status)) {
      setActiveJobId(null);
      void jobsQuery.refetch();
    }
  }, [pollQuery.data, jobsQuery]);

  const runAction = (type: 'TRANSLATE' | 'SUMMARIZE') => {
    setLastAction(type);
    createJob.mutate(
      {
        file_id: fileId,
        type,
        input: type === 'TRANSLATE' ? 'en' : undefined,
      },
      {
        onSuccess: (job) => setActiveJobId(job.id),
      },
    );
  };

  const isPolling =
    Boolean(activeJobId) &&
    (!pollQuery.data || isAIJobRunning(pollQuery.data.status));

  return {
    translation,
    summaries,
    isLoading: createJob.isPending || isPolling,
    isError:
      createJob.isError ||
      pollQuery.data?.status === 'failed' ||
      pollQuery.isError,
    jobStatus: pollQuery.data?.status,
    runTranslate: () => runAction('TRANSLATE'),
    runSummarize: () => runAction('SUMMARIZE'),
    handleRetry: () => {
      if (!lastAction) return;
      runAction(lastAction);
    },
  };
}

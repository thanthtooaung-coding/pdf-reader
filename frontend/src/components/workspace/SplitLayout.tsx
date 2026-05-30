import { useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import PDFWorkspace from '@/components/pdf/PDFWorkspace';
import Tab from '@/components/ui/Tab';
import InsightsPanel, {
  useInsightsState,
} from '@/components/workspace/InsightsPanel';
import NotesPanel from '@/components/workspace/NotesPanel';
import { panelSlide } from '@/lib/motion';

type RightPanelTab = 'insights' | 'notes';

interface SplitLayoutProps {
  workspaceId: string;
  fileId: string;
  fileName: string;
}

export default function SplitLayout({
  workspaceId,
  fileId,
  fileName,
}: SplitLayoutProps) {
  const [activeTab, setActiveTab] = useState<RightPanelTab>('insights');
  const [noteDraft, setNoteDraft] = useState<string | null>(null);
  const insights = useInsightsState(workspaceId, fileId);

  const handleTranslate = () => {
    setActiveTab('insights');
    insights.runTranslate();
  };

  const handleSummarize = () => {
    setActiveTab('insights');
    insights.runSummarize();
  };

  const handleTakeNote = (text: string) => {
    setActiveTab('notes');
    setNoteDraft(text);
  };

  return (
    <main className="flex h-full overflow-hidden bg-white">
      <section className="w-[60%] shrink-0 border-r border-slate-200">
        <PDFWorkspace
          fileId={fileId}
          fileName={fileName}
          onTranslate={handleTranslate}
          onSummarize={handleSummarize}
          onTakeNote={handleTakeNote}
          isAIProcessing={insights.isLoading}
        />
      </section>

      <aside className="flex w-[40%] shrink-0 flex-col">
        <nav
          role="tablist"
          aria-label="Right panel tabs"
          className="flex border-b border-slate-200 bg-white"
        >
          <Tab
            label="AI Insights"
            isActive={activeTab === 'insights'}
            onClick={() => setActiveTab('insights')}
          />
          <Tab
            label="My Notes"
            isActive={activeTab === 'notes'}
            onClick={() => setActiveTab('notes')}
          />
        </nav>

        <div className="relative flex-1 overflow-hidden">
          <AnimatePresence mode="wait">
            {activeTab === 'insights' ? (
              <motion.div
                key="insights"
                className="absolute inset-0"
                {...panelSlide}
                transition={{ duration: 0.2 }}
              >
                <InsightsPanel
                  translation={insights.translation}
                  summaries={insights.summaries}
                  isLoading={insights.isLoading}
                  isError={insights.isError}
                  jobStatus={insights.jobStatus}
                  onRetry={insights.handleRetry}
                />
              </motion.div>
            ) : (
              <motion.div
                key="notes"
                className="absolute inset-0"
                {...panelSlide}
                transition={{ duration: 0.2 }}
              >
                <NotesPanel
                  fileId={fileId}
                  draftSourceText={noteDraft}
                  onDraftConsumed={() => setNoteDraft(null)}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </aside>
    </main>
  );
}

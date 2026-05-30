import { useEffect, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { AnimatePresence } from 'framer-motion';
import { Document, Page } from 'react-pdf';
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';
import { downloadFile } from '@/api/files';
import Button from '@/components/ui/Button';
import FloatingToolbar from '@/components/pdf/FloatingToolbar';
import { getToolbarPosition, useSelection } from '@/hooks/useSelection';

interface PDFWorkspaceProps {
  fileId: string;
  fileName: string;
  onTranslate: () => void;
  onSummarize: () => void;
  onTakeNote: (text: string) => void;
  isAIProcessing?: boolean;
}

export default function PDFWorkspace({
  fileId,
  fileName,
  onTranslate,
  onSummarize,
  onTakeNote,
  isAIProcessing,
}: PDFWorkspaceProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [numPages, setNumPages] = useState(0);
  const [pageNumber, setPageNumber] = useState(1);
  const [blobUrl, setBlobUrl] = useState<string | null>(null);

  const { selectedText, rect, clearSelection } = useSelection(containerRef);

  const downloadQuery = useQuery({
    queryKey: ['file-download', fileId],
    queryFn: () => downloadFile(fileId),
    enabled: Boolean(fileId),
  });

  useEffect(() => {
    if (!downloadQuery.data) return;
    const url = URL.createObjectURL(downloadQuery.data);
    setBlobUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [downloadQuery.data]);

  const toolbarPosition = rect ? getToolbarPosition(rect) : null;

  return (
    <section className="flex h-full flex-col bg-slate-50">
      <header className="flex items-center justify-between border-b border-slate-200 bg-white px-4 py-3">
        <div>
          <h1 className="text-sm font-semibold text-slate-900">{fileName}</h1>
          <p className="text-xs text-slate-500">PDF Workspace</p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            disabled={isAIProcessing}
            onClick={onTranslate}
          >
            Translate doc
          </Button>
          <Button
            variant="ghost"
            disabled={isAIProcessing}
            onClick={onSummarize}
          >
            Summarize doc
          </Button>
          {numPages > 0 && (
            <>
              <Button
                variant="ghost"
                disabled={pageNumber <= 1}
                onClick={() => setPageNumber((page) => Math.max(1, page - 1))}
              >
                Prev
              </Button>
              <span className="text-sm text-slate-600">
                Page {pageNumber} of {numPages}
              </span>
              <Button
                variant="ghost"
                disabled={pageNumber >= numPages}
                onClick={() =>
                  setPageNumber((page) => Math.min(numPages, page + 1))
                }
              >
                Next
              </Button>
            </>
          )}
        </div>
      </header>

      <div ref={containerRef} className="relative flex-1 overflow-auto p-6">
        {downloadQuery.isLoading && (
          <p className="text-sm text-slate-500">Loading PDF...</p>
        )}
        {downloadQuery.isError && (
          <p className="text-sm text-red-600">
            Failed to load PDF. Try again from the file list.
          </p>
        )}
        {blobUrl && (
          <div className="mx-auto w-fit rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
            <Document
              file={blobUrl}
              onLoadSuccess={({ numPages: loadedPages }) =>
                setNumPages(loadedPages)
              }
              loading={
                <p className="p-8 text-sm text-slate-500">Rendering PDF...</p>
              }
              error={
                <p className="p-8 text-sm text-red-600">
                  Failed to render PDF.
                </p>
              }
            >
              <Page
                pageNumber={pageNumber}
                renderTextLayer
                renderAnnotationLayer
              />
            </Document>
          </div>
        )}

        <AnimatePresence>
          {selectedText && toolbarPosition && (
            <FloatingToolbar
              key="selection-toolbar"
              x={toolbarPosition.x}
              y={toolbarPosition.y}
              onTranslate={() => {
                onTranslate();
                clearSelection();
              }}
              onSummarize={() => {
                onSummarize();
                clearSelection();
              }}
              onTakeNote={() => {
                onTakeNote(selectedText);
                clearSelection();
              }}
            />
          )}
        </AnimatePresence>
      </div>
    </section>
  );
}

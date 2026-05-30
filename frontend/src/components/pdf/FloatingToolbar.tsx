import { motion } from 'framer-motion';
import Button from '@/components/ui/Button';
import { springSnappy, toolbarPop } from '@/lib/motion';

interface FloatingToolbarProps {
  x: number;
  y: number;
  onTranslate: () => void;
  onSummarize: () => void;
  onTakeNote: () => void;
}

export default function FloatingToolbar({
  x,
  y,
  onTranslate,
  onSummarize,
  onTakeNote,
}: FloatingToolbarProps) {
  return (
    <motion.div
      data-floating-toolbar
      className="fixed z-50 flex items-center gap-1 rounded-lg border border-slate-200 bg-white p-1 shadow-lg"
      style={{ transform: 'translate(-50%, calc(-100% - 8px))' }}
      initial={toolbarPop.initial}
      animate={{ ...toolbarPop.animate, left: x, top: y }}
      exit={toolbarPop.exit}
      transition={springSnappy}
    >
      <Button variant="toolbar" onClick={onTranslate}>
        Translate
      </Button>
      <Button variant="toolbar" onClick={onSummarize}>
        Summarize
      </Button>
      <Button variant="toolbar" onClick={onTakeNote}>
        Take Note
      </Button>
    </motion.div>
  );
}

import { motion } from 'framer-motion';

interface TabProps {
  label: string;
  isActive: boolean;
  onClick: () => void;
}

export default function Tab({ label, isActive, onClick }: TabProps) {
  return (
    <button
      type="button"
      role="tab"
      aria-selected={isActive}
      onClick={onClick}
      className={`relative px-4 py-3 text-sm font-medium transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300 ${
        isActive ? 'text-slate-900' : 'text-slate-500 hover:text-slate-700'
      }`}
    >
      {label}
      {isActive && (
        <motion.span
          layoutId="active-tab-indicator"
          className="absolute inset-x-0 bottom-0 h-0.5 bg-slate-900"
          transition={{ type: 'spring', stiffness: 400, damping: 30 }}
        />
      )}
    </button>
  );
}

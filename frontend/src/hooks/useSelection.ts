import { useCallback, useEffect, useState } from 'react';

export interface SelectionState {
  selectedText: string | null;
  rect: DOMRect | null;
}

export function useSelection(containerRef: React.RefObject<HTMLElement | null>) {
  const [selection, setSelection] = useState<SelectionState>({
    selectedText: null,
    rect: null,
  });

  const clearSelection = useCallback(() => {
    setSelection({ selectedText: null, rect: null });
    window.getSelection()?.removeAllRanges();
  }, []);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const handleMouseUp = () => {
      const activeSelection = window.getSelection();
      if (!activeSelection || activeSelection.isCollapsed) {
        setSelection({ selectedText: null, rect: null });
        return;
      }

      const text = activeSelection.toString().trim();
      if (!text) {
        setSelection({ selectedText: null, rect: null });
        return;
      }

      const range = activeSelection.getRangeAt(0);
      const anchorNode = activeSelection.anchorNode;
      if (!anchorNode || !container.contains(anchorNode)) {
        setSelection({ selectedText: null, rect: null });
        return;
      }

      const rect = range.getBoundingClientRect();
      setSelection({ selectedText: text, rect });
    };

    const handleMouseDown = (event: MouseEvent) => {
      const target = event.target as Node;
      if (
        target instanceof Element &&
        target.closest('[data-floating-toolbar]')
      ) {
        return;
      }
      setSelection({ selectedText: null, rect: null });
    };

    container.addEventListener('mouseup', handleMouseUp);
    document.addEventListener('mousedown', handleMouseDown);

    return () => {
      container.removeEventListener('mouseup', handleMouseUp);
      document.removeEventListener('mousedown', handleMouseDown);
    };
  }, [containerRef]);

  return {
    selectedText: selection.selectedText,
    rect: selection.rect,
    clearSelection,
  };
}

export function getToolbarPosition(rect: DOMRect) {
  return {
    x: rect.left + rect.width / 2,
    y: rect.top,
  };
}

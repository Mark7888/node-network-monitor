import { ReactNode, useEffect, useRef } from 'react';
import { X } from 'lucide-react';
import Button from './Button';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  actions?: ReactNode;
  size?: 'sm' | 'md' | 'lg';
}

/**
 * Reusable Modal component using daisyUI modal
 */
export default function Modal({ 
  isOpen, 
  onClose, 
  title, 
  children, 
  actions,
  size = 'md' 
}: ModalProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (isOpen) {
      dialog.showModal();
    } else {
      dialog.close();
    }
  }, [isOpen]);

  const sizeClass = {
    sm: 'max-w-sm',
    md: 'max-w-2xl',
    lg: 'max-w-5xl',
  }[size];

  return (
    <dialog ref={dialogRef} className="modal" onClose={onClose}>
      <div className={`modal-box ${sizeClass}`}>
        {/* Header */}
        <div className="flex justify-between items-center mb-4">
          {title && <h3 className="font-bold text-lg">{title}</h3>}
          <button
            onClick={onClose}
            className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
            aria-label="Close modal"
          >
            <X size={20} />
          </button>
        </div>

        {/* Content */}
        <div className="py-4">{children}</div>

        {/* Actions */}
        {actions && (
          <div className="modal-action">
            {actions}
          </div>
        )}
      </div>
      
      {/* Backdrop */}
      <form method="dialog" className="modal-backdrop">
        <button onClick={onClose}>close</button>
      </form>
    </dialog>
  );
}

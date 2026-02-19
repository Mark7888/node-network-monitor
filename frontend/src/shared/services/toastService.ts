import toast from 'react-hot-toast';

/**
 * Track active toasts by message to prevent duplicates
 */
const activeToasts = new Map<string, string>();

/**
 * Custom toast handler that prevents duplicate toasts with the same message
 * from appearing simultaneously. Only one toast per unique message will be shown at a time.
 */
export const showToast = {
  /**
   * Show a success toast (unless the same message is already displayed)
   */
  success: (message: string) => {
    // Check if this message is already being displayed
    if (activeToasts.has(message)) {
      return;
    }
    
    // Create the toast and store its ID
    const toastId = toast.success(message);
    activeToasts.set(message, toastId);
    
    // Clean up after the default duration (3000ms from main.tsx config)
    setTimeout(() => {
      activeToasts.delete(message);
    }, 3000);
  },
  
  /**
   * Show an error toast (unless the same message is already displayed)
   */
  error: (message: string) => {
    // Check if this message is already being displayed
    if (activeToasts.has(message)) {
      return;
    }
    
    // Create the toast and store its ID
    const toastId = toast.error(message);
    activeToasts.set(message, toastId);
    
    // Clean up after the default duration (4000ms from main.tsx config)
    setTimeout(() => {
      activeToasts.delete(message);
    }, 4000);
  },
};

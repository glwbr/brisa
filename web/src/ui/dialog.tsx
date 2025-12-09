import * as React from 'react';
import {
  Dialog as AriaDialog,
  DialogTrigger as AriaDialogTrigger,
  Heading as AriaHeading,
  Modal as AriaModal,
  ModalOverlay as AriaModalOverlay,
} from 'react-aria-components';
import type {
  DialogProps as AriaDialogProps,
  HeadingProps as AriaHeadingProps,
  ModalOverlayProps as AriaModalOverlayProps,
} from 'react-aria-components';
import { cn } from '@/lib/cn';

const DialogTrigger = AriaDialogTrigger;

const DialogOverlay = React.forwardRef<HTMLDivElement, AriaModalOverlayProps>(
  ({ className, isDismissable = true, ...props }, ref) => {
    return (
      <AriaModalOverlay
        ref={ref}
        isDismissable={isDismissable}
        className={(values) =>
          cn(
            'fixed inset-0 z-50 flex min-h-full items-center justify-center overflow-y-auto p-4',
            'bg-black/60 backdrop-blur-sm',
            'data-entering:animate-in data-entering:fade-in data-entering:duration-200',
            'data-exiting:animate-out data-exiting:fade-out data-exiting:duration-150',
            typeof className === 'function' ? className(values) : className,
          )
        }
        {...props}
      />
    );
  },
);

DialogOverlay.displayName = 'DialogOverlay';

interface DialogModalProps extends AriaModalOverlayProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
}

const DialogModal = React.forwardRef<HTMLDivElement, DialogModalProps>(({ className, size = 'md', ...props }, ref) => {
  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
  };

  return (
    <AriaModal
      ref={ref}
      className={(values) =>
        cn(
          'w-full overflow-hidden rounded-2xl p-6 text-left shadow-xl',
          'bg-slate-900 ring-1 ring-white/10',
          'data-entering:animate-in data-entering:zoom-in-95 data-entering:fade-in data-entering:duration-200',
          'data-exiting:animate-out data-exiting:zoom-out-95 data-exiting:fade-out data-exiting:duration-150',
          sizeClasses[size],
          typeof className === 'function' ? className(values) : className,
        )
      }
      {...props}
    />
  );
});

DialogModal.displayName = 'DialogModal';

const DialogContent = React.forwardRef<HTMLDivElement, AriaDialogProps>(({ className, ...props }, ref) => {
  return <AriaDialog ref={ref} className={cn('relative outline-none', className)} {...props} />;
});

DialogContent.displayName = 'DialogContent';

const DialogHeader = ({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  return <div className={cn('mb-4', className)} {...props} />;
};

DialogHeader.displayName = 'DialogHeader';

const DialogTitle = React.forwardRef<HTMLHeadingElement, AriaHeadingProps>(({ className, ...props }, ref) => {
  return (
    <AriaHeading
      ref={ref}
      slot="title"
      className={cn('text-lg leading-6 font-semibold text-white', className)}
      {...props}
    />
  );
});

DialogTitle.displayName = 'DialogTitle';

const DialogDescription = React.forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  ({ className, ...props }, ref) => {
    return <p ref={ref} className={cn('text-sm text-slate-400', className)} {...props} />;
  },
);

DialogDescription.displayName = 'DialogDescription';

const DialogFooter = ({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) => {
  return <div className={cn('mt-6 flex justify-end gap-3', className)} {...props} />;
};

DialogFooter.displayName = 'DialogFooter';

export {
  DialogTrigger,
  DialogOverlay,
  DialogModal,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
};

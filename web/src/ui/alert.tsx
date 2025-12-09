import * as React from 'react';
import { cva } from 'class-variance-authority';
import { AlertCircle, AlertTriangle, CheckCircle2, Info, X } from 'lucide-react';
import type { VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/cn';
import { Button } from './button';

const alertVariants = cva(
  ['relative rounded-xl border p-4 text-sm', 'animate-in fade-in slide-in-from-top-2 duration-300'],
  {
    variants: {
      variant: {
        error: 'border-red-500/20 bg-red-500/10 text-red-200',
        warning: 'border-amber-500/20 bg-amber-500/10 text-amber-200',
        success: 'border-emerald-500/20 bg-emerald-500/10 text-emerald-200',
        info: 'border-sky-500/20 bg-sky-500/10 text-sky-200',
      },
    },
    defaultVariants: {
      variant: 'info',
    },
  },
);

const alertIconMap = {
  error: AlertCircle,
  warning: AlertTriangle,
  success: CheckCircle2,
  info: Info,
};

export interface AlertProps extends React.HTMLAttributes<HTMLDivElement>, VariantProps<typeof alertVariants> {
  title?: string;
  onDismiss?: () => void;
  icon?: React.ReactNode;
}

const Alert = React.forwardRef<HTMLDivElement, AlertProps>(
  ({ className, variant = 'info', title, children, onDismiss, icon, ...props }, ref) => {
    const IconComponent = alertIconMap[variant ?? 'info'];

    return (
      <div ref={ref} role="alert" className={cn(alertVariants({ variant }), className)} {...props}>
        <div className="flex gap-3">
          <div className="shrink-0">{icon ?? <IconComponent className="h-5 w-5" />}</div>
          <div className="flex-1 space-y-1">
            {title && <p className="font-semibold">{title}</p>}
            {children && <div className="opacity-90">{children}</div>}
          </div>
          {onDismiss && (
            <Button
              variant="ghost"
              size="icon"
              onClick={onDismiss}
            >
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>
      </div>
    );
  },
);

Alert.displayName = 'Alert';

export { Alert, alertVariants };

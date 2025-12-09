import * as React from 'react';
import { Button as AriaButton } from 'react-aria-components';
import { cva } from 'class-variance-authority';
import { Loader2 } from 'lucide-react';
import { cn } from '../lib/cn';
import type { VariantProps } from 'class-variance-authority';
import type { ButtonProps as AriaButtonProps, ButtonRenderProps } from 'react-aria-components';
import type { LucideIcon } from 'lucide-react';

const buttonVariants = cva(
  [
    'font-medium transition-all duration-300 flex items-center justify-center gap-2 cursor-pointer',
    'outline-none data-focus-visible:ring-2 data-focus-visible:ring-sky-500/50 data-focus-visible:ring-offset-2 data-focus-visible:ring-offset-slate-950',
    'data-pressed:scale-95',
    'data-disabled:pointer-events-none data-disabled:opacity-50 data-disabled:cursor-not-allowed',
  ],
  {
    variants: {
      variant: {
        primary: [
          'bg-sky-500 text-white shadow-lg shadow-sky-500/20 border-transparent',
          'data-hovered:bg-sky-600 data-hovered:shadow-sky-500/40',
          'data-pressed:bg-sky-700',
        ],
        secondary: [
          'bg-slate-900 dark:bg-white text-white dark:text-slate-900 shadow-md border-transparent',
          'data-hovered:bg-slate-800 dark:data-hovered:bg-slate-100',
          'data-pressed:bg-slate-700 dark:data-pressed:bg-slate-200',
        ],
        outline: [
          'border-2 border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-200 bg-transparent',
          'data-hovered:border-sky-500 data-hovered:text-sky-500 dark:data-hovered:border-sky-400 dark:data-hovered:text-sky-400',
          'data-pressed:border-sky-600 data-pressed:text-sky-600',
        ],
        ghost: [
          'text-slate-600 dark:text-slate-400 border-transparent',
          'data-hovered:text-sky-600 dark:data-hovered:text-sky-400 data-hovered:bg-sky-50 dark:data-hovered:bg-slate-800/50',
          'data-pressed:bg-sky-100 dark:data-pressed:bg-slate-800',
        ],
        destructive: [
          'bg-red-500 text-white shadow-lg shadow-red-500/20 border-transparent',
          'data-hovered:bg-red-600 data-hovered:shadow-red-500/40',
          'data-pressed:bg-red-700',
        ],
      },
      size: {
        sm: 'h-8 px-4 text-xs rounded-full',
        md: 'h-11 px-6 text-sm rounded-full',
        lg: 'h-14 px-8 text-base rounded-full',
        icon: 'h-10 w-10 rounded-full p-0',
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md',
    },
  },
);

export interface ButtonProps extends AriaButtonProps, VariantProps<typeof buttonVariants> {
  className?: string | ((values: ButtonRenderProps) => string);
  icon?: LucideIcon;
  isLoading?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, icon: Icon, isLoading, isDisabled, children, ...props }, ref) => {
    return (
      <AriaButton
        className={(values) =>
          cn(buttonVariants({ variant, size }), typeof className === 'function' ? className(values) : className)
        }
        ref={ref}
        isDisabled={isDisabled || isLoading}
        {...props}
      >
        {(renderProps) => (
          <>
            {isLoading ? (
              <Loader2 className="animate-spin" size={size === 'sm' ? 14 : 18} />
            ) : (
              Icon && <Icon size={size === 'sm' ? 14 : 18} />
            )}
            {typeof children === 'function' ? children(renderProps) : children}
          </>
        )}
      </AriaButton>
    );
  },
);

Button.displayName = 'Button';

export { Button, buttonVariants };

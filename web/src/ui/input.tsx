import * as React from 'react';
import { Input as AriaInput } from 'react-aria-components';
import { cva } from 'class-variance-authority';
import type { VariantProps } from 'class-variance-authority';
import type { InputProps as AriaInputProps } from 'react-aria-components';
import { cn } from '@/lib/cn';

const inputVariants = cva(
  [
    'w-full border rounded-xl',
    'text-slate-900 dark:text-white',
    'placeholder-slate-400',
    'transition-all duration-200',
    'outline-none',

    'data-disabled:opacity-60 data-disabled:cursor-not-allowed',

    'data-focused:ring-2',

    'group-data-invalid:border-red-500',
    'group-data-invalid:text-red-900 group-data-invalid:dark:text-red-200',
    'group-data-invalid:placeholder-red-400',
    'group-data-invalid:data-focused:ring-2 group-data-invalid:data-focused:ring-red-500/20',
    'group-data-invalid:bg-red-50/50 group-data-invalid:dark:bg-red-900/20',
    'group-data-invalid:data-focused:border-red-600',
  ],
  {
    variants: {
      variant: {
        outline: [
          'bg-white dark:bg-slate-900',
          'border-slate-200 dark:border-slate-700',
          'data-focused:border-sky-500 data-focused:ring-sky-500/20',
        ],
        filled: [
          'bg-slate-100 dark:bg-slate-800',
          'border-transparent',
          'data-focused:bg-white dark:data-focused:bg-slate-900',
          'data-focused:border-sky-500 data-focused:ring-sky-500/20',
        ],
        ghost: [
          'bg-transparent border-transparent',
          'data-hovered:bg-slate-100 dark:data-hovered:bg-slate-800',
          'data-focused:bg-transparent',
          'data-focused:border-sky-500',
        ],
      },
      size: {
        sm: 'h-9 px-3 text-xs',
        md: 'h-11 px-4 text-sm',
        lg: 'h-14 px-6 text-base',
      },
    },
    defaultVariants: {
      variant: 'outline',
      size: 'md',
    },
  },
);

export interface InputProps extends Omit<AriaInputProps, 'size'>, VariantProps<typeof inputVariants> {}

const Input = React.forwardRef<HTMLInputElement, InputProps>((props, ref) => {
  const { className, variant, size, ...rest } = props;
  return <AriaInput ref={ref} className={cn(inputVariants({ variant, size }), className)} {...rest} />;
});

Input.displayName = 'Input';

export { Input, inputVariants };

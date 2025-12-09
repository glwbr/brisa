import * as React from 'react';
import {
  FieldError as AriaFieldError,
  Label as AriaLabel,
  Text as AriaText,
  TextField as AriaTextField,
} from 'react-aria-components';
import { AlertCircle, CheckCircle2 } from 'lucide-react';
import { cn } from '../lib/cn';
import { Input } from './input';
import type { TextFieldProps as AriaTextFieldProps, ValidationResult } from 'react-aria-components';
import type { InputProps } from './input';

const TextFieldLabel = (props: React.ComponentProps<typeof AriaLabel>) => {
  const { className, children, ...rest } = props;

  return (
    <AriaLabel
      className={cn(
        'mb-1.5 ml-1 block text-sm font-semibold text-slate-700 dark:text-slate-300',
        'group-data-invalid:text-red-500',
        'group-data-required:after:ml-0.5 group-data-required:after:text-red-500 group-data-required:after:content-["*"]',
        className,
      )}
      {...rest}
    >
      {children}
    </AriaLabel>
  );
};

const TextFieldError = ({ className, ...rest }: React.ComponentProps<typeof AriaFieldError>) => {
  return (
    <AriaFieldError
      className={cn('mt-1.5 ml-1 block text-xs font-medium text-red-500', 'animate-in-fade-slide', className)}
      {...rest}
    />
  );
};

export interface TextFieldProps extends AriaTextFieldProps {
  label?: string;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
  placeholder?: string;
  hideLabel?: boolean;
  isSuccess?: boolean;
  inputProps?: InputProps;
}

const TextField = React.forwardRef<HTMLInputElement, TextFieldProps>((props, ref) => {
  const { className, label, description, errorMessage, placeholder, hideLabel, isSuccess, inputProps, ...rest } = props;

  return (
    <AriaTextField
      {...rest}
      data-success={isSuccess || undefined}
      className={cn('group flex w-full flex-col', className)}
    >
      {({ isInvalid }) => {
        const hasValue = !!props.value || !!props.defaultValue;
        const showIcon = isInvalid || (isSuccess && hasValue);

        return (
          <>
            {label && <TextFieldLabel className={hideLabel ? 'sr-only' : ''}>{label}</TextFieldLabel>}

            <div className="relative">
              <Input ref={ref} placeholder={placeholder} className={showIcon ? 'pr-10' : ''} {...inputProps} />

              {showIcon && (
                <div className="animate-in-zoom pointer-events-none absolute top-1/2 right-3 -translate-y-1/2">
                  {isInvalid ? (
                    <AlertCircle size={16} className="text-red-500" />
                  ) : (
                    <CheckCircle2 size={16} className="text-sky-500" />
                  )}
                </div>
              )}
            </div>

            {description && !isInvalid && (
              <AriaText slot="description" className="mt-1.5 ml-1 block text-xs text-slate-500">
                {description}
              </AriaText>
            )}

            <TextFieldError>{errorMessage}</TextFieldError>
          </>
        );
      }}
    </AriaTextField>
  );
});

TextField.displayName = 'TextField';

export { TextField };

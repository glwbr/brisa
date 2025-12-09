import { useStore } from '@tanstack/react-form';
import { Button } from './button';
import { TextField } from './text-field';
import type { TextFieldProps } from './text-field';
import type { ButtonProps } from './button';
import { useFieldContext, useFormContext } from '@/hooks/use-app-form';

type ErrorLike = string | { message: string };

function getFirstError(errors: Array<ErrorLike>): string | undefined {
  if (!errors.length) return undefined;
  const first = errors[0];
  return typeof first === 'string' ? first : first.message;
}

export interface FormTextFieldProps extends Omit<
  TextFieldProps,
  'name' | 'value' | 'onChange' | 'onBlur' | 'isInvalid' | 'errorMessage' | 'defaultValue' | 'isSuccess'
> {}

export function FormTextField(props: FormTextFieldProps) {
  const field = useFieldContext<string>();

  const errors = useStore(field.store, (s) => s.meta.errors);
  const isTouched = useStore(field.store, (s) => s.meta.isTouched);
  const value = useStore(field.store, (s) => s.value);

  const hasErrors = errors.length > 0 && isTouched;
  const errorMessage = hasErrors ? getFirstError(errors) : undefined;
  const isInvalid = hasErrors && !!errorMessage;
  const isSuccess = isTouched && errors.length === 0 && !!value;
  const isDisabled = props.isDisabled;

  return (
    <TextField
      {...props}
      {...{
        value,
        isDisabled,
        isInvalid,
        isSuccess,
        errorMessage,
        name: field.name,
        onChange: field.handleChange,
        onBlur: field.handleBlur,
      }}
    />
  );
}

export interface SubmitButtonProps extends ButtonProps {}

export function SubmitButton(props: SubmitButtonProps) {
  const form = useFormContext();

  return (
    <form.Subscribe selector={(state) => [state.canSubmit, state.isSubmitting]}>
      {([canSubmit, isSubmitting]) => (
        <Button type="submit" isDisabled={!canSubmit} isLoading={isSubmitting} {...props}>
          {props.children}
        </Button>
      )}
    </form.Subscribe>
  );
}

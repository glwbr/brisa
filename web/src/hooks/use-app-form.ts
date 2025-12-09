import { createFormHook, createFormHookContexts } from '@tanstack/react-form';
import { FormTextField, SubmitButton } from '@/ui';

export const { fieldContext, formContext, useFieldContext, useFormContext } = createFormHookContexts();

export const { useAppForm } = createFormHook({
  fieldComponents: {
    FormTextField,
  },
  formComponents: {
    SubmitButton,
  },
  fieldContext,
  formContext,
});

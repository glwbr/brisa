import { z } from 'zod';
import { useSubmitCaptcha } from '@/api/endpoints/invoice-jobs';
import { useAppForm } from '@/hooks/use-app-form';

const ZCaptchaSolution = z.object({
  solution: z.string().min(1, 'Please enter the captcha solution'),
});

interface UseCaptchaFormOptions {
  jobId: string;
  onSuccess?: () => void;
  onError?: (error: unknown) => void;
}

export function useCaptchaForm({ jobId, onSuccess, onError }: UseCaptchaFormOptions) {
  const {
    mutateAsync: submitCaptcha,
    isPending: isSubmitCaptchaPending,
    isError: isSubmitCaptchaError,
    error: submitCaptchaError,
  } = useSubmitCaptcha({
    mutation: {
      onSuccess: () => {
        onSuccess?.();
      },
      onError: (error) => {
        onError?.(error);
      },
    },
  });

  const form = useAppForm({
    defaultValues: {
      solution: '',
    },
    validators: {
      onChange: ZCaptchaSolution,
    },
    onSubmit: async ({ value }) => {
      await submitCaptcha({ id: jobId, data: value });
    },
  });

  return {
    form,
    isPending: isSubmitCaptchaPending,
    isError: isSubmitCaptchaError,
    error: submitCaptchaError,
  };
}

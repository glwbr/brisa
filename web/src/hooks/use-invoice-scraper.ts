import { z } from 'zod';
import { useQueryClient } from '@tanstack/react-query';
import { useCreateInvoiceJob, useGetInvoiceJob } from '@/api/endpoints/invoice-jobs';
import { useAppForm } from '@/hooks/use-app-form';
import { JobStatus } from '@/api/model/jobStatus';

const ZAccessKey = z.object({
  accessKey: z
    .string()
    .length(44, 'Access key must be exactly 44 digits')
    .regex(/^\d+$/, 'Access key must contain only numbers'),
});

const ACTIVE_JOB_KEY = ['invoice-scraper', 'active-job-id'];

export function useInvoiceScraper() {
  const queryClient = useQueryClient();

  const { mutateAsync: createJob } = useCreateInvoiceJob({
    mutation: {
      onSuccess: (data) => {
        if (data.jobId) {
          queryClient.setQueryData(ACTIVE_JOB_KEY, data.jobId);
        }
      },
    },
  });

  const jobId = queryClient.getQueryData<string>(ACTIVE_JOB_KEY);

  // TODO: replace this polling with SSE for fun.
  const { data: job } = useGetInvoiceJob(jobId!, {
    query: {
      enabled: !!jobId,
      refetchInterval: ({ state }) => {
        const status = state.data?.status;
        if (status === JobStatus.completed || status === JobStatus.failed) return false;
        return 1000;
      },
    },
  });

  const form = useAppForm({
    defaultValues: {
      accessKey: '',
    },
    validators: {
      onChange: ZAccessKey,
    },
    onSubmit: async ({ value }) => {
      await createJob({ data: value });
    },
  });

  const isWaitingCaptcha = job?.status === JobStatus.waiting_captcha;
  const isProcessing = job?.status === JobStatus.running || job?.status === JobStatus.created;
  const isCompleted = job?.status === JobStatus.completed;
  const isFailed = job?.status === JobStatus.failed;

  return {
    form,
    jobId,
    job,
    isWaitingCaptcha,
    isProcessing,
    isCompleted,
    isFailed,
    error: job?.error,
    result: job?.result,
    captcha: job?.captcha,
  };
}

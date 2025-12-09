import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export function getContext() {
  const queryClient = new QueryClient();
  return {
    queryClient,
  };
}

interface ProviderProps {
  children: React.ReactNode;
  queryClient: QueryClient;
}

export function Provider(props: ProviderProps) {
  const { children, queryClient } = props;

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
}

import { Outlet, createRootRouteWithContext } from '@tanstack/react-router';
import { TanStackDevtools } from '@tanstack/react-devtools';

import type { QueryClient } from '@tanstack/react-query';

import TanStackRouterDevtools from '@/integrations/tanstack-router/devtools';
import TanStackQueryDevtools from '@/integrations/tanstack-query/devtools';

interface RouterContext {
  queryClient: QueryClient;
}

function RootLayout() {
  return (
    <div className="flex min-h-dvh flex-col">
      <main className="flex flex-1 flex-col">
        <Outlet />
      </main>
      <TanStackDevtools
        config={{ position: 'bottom-right' }}
        plugins={[TanStackRouterDevtools, TanStackQueryDevtools]}
      />
    </div>
  );
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
});

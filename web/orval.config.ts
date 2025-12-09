import { defineConfig } from 'orval';

export default defineConfig({
  brisa: {
    output: {
      mode: 'single',
      target: 'src/api/endpoints/invoice-jobs.ts',
      schemas: 'src/api/model',
      client: 'react-query',
      mock: false,
      override: {
        mutator: {
          path: './src/lib/custom-fetch.ts',
          name: 'customFetch',
        },
      },
      prettier: true,
      clean: true,
    },
    input: {
      target: '../openapi.yaml',
    },
  },
});

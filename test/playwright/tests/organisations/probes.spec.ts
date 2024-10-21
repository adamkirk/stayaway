import { test, expect } from '@playwright/test';

test('organisations', async ({ request }) => {
  const resp = await request.get(`/api/v1/_probes/startup`);

  expect(resp.status()).toEqual(204)
});

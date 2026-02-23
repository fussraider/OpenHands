import { test, expect } from '@playwright/test';

test('has title and loads successfully from Go backend', async ({ page }) => {
  await page.goto('http://localhost:3000/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/OpenHands/);
});

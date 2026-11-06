import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e-p0',
  outputDir: '.p0-test-results',
  workers: 1,
  timeout: 120000,
  forbidOnly: !!process.env.CI,
  reporter: 'list',
  use: { actionTimeout: 15000, baseURL: 'http://127.0.0.1:18800', trace: 'retain-on-failure', screenshot: 'only-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'], launchOptions: process.env.PLAYWRIGHT_EXECUTABLE_PATH ? { executablePath: process.env.PLAYWRIGHT_EXECUTABLE_PATH } : {} } }],
  webServer: [
    { command: 'node scripts/p0-backend.mjs', url: 'http://127.0.0.1:18808/api/v1/system/health', timeout: 180000, reuseExistingServer: false },
    { command: 'node scripts/p0-preview.mjs', url: 'http://127.0.0.1:18800', timeout: 30000, reuseExistingServer: false },
  ],
});

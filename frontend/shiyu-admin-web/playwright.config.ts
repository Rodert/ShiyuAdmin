import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: {
    command: 'npm run start:dev',
    url: 'http://localhost:8000',
    reuseExistingServer: !process.env.CI,
    timeout: 300 * 1000, // 增加到 5 分钟，给前端服务更多启动时间
    stdout: 'pipe',
    stderr: 'pipe',
  },
});


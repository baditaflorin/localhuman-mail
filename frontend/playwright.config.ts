import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  timeout: 30_000,
  use: {
    baseURL: "http://127.0.0.1:4187/localhuman-mail/",
    trace: "retain-on-failure"
  },
  webServer: {
    command: "npm run preview -- --port 4187 --strictPort",
    url: "http://127.0.0.1:4187/localhuman-mail/",
    reuseExistingServer: false,
    timeout: 30_000
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] }
    }
  ]
});

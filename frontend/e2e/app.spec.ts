import { expect, test } from "@playwright/test";
import path from "node:path";

test("renders the mailbox workbench and local project links", async ({ page }) => {
  await page.goto("/");

  await expect(page.getByRole("heading", { name: "localhuman-mail" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Star" })).toHaveAttribute(
    "href",
    "https://github.com/baditaflorin/localhuman-mail"
  );
  await expect(page.getByRole("link", { name: "PayPal" })).toHaveAttribute(
    "href",
    "https://www.paypal.com/paypalme/florinbadita"
  );

  await page.getByLabel("Search messages").fill("launch");
  await expect(page.getByRole("heading", { name: "Thursday launch checklist" })).toBeVisible();
  await page.getByRole("button", { name: "Draft" }).click();
  await expect(page.getByLabel("Generated reply draft")).toBeVisible();
});

test("imports real EML and exports user work paths when backend is running", async ({
  page,
  request,
  context
}) => {
  const backendHealth = await request.get("http://127.0.0.1:18080/healthz").catch(() => null);
  test.skip(!backendHealth?.ok(), "smoke backend is not running");
  await context.grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.goto("/");
  await page.getByLabel("Backend URL").fill("http://127.0.0.1:18080");
  await expect(page.getByText("Connected")).toBeVisible();

  const chooserPromise = page.waitForEvent("filechooser");
  await page.getByRole("button", { name: "EML files" }).click();
  const chooser = await chooserPromise;
  await chooser.setFiles(path.resolve("../test/fixtures/realdata/06-calendar-invite-rfc5545.eml"));

  await expect(page.getByText("Import complete")).toBeVisible();
  await page.getByLabel("Search messages").fill("Design Review");
  await expect(page.getByRole("heading", { name: "Design Review" })).toBeVisible();
  await expect(page.getByText("calendar_invite · high confidence")).toBeVisible();

  await page.getByRole("button", { name: "Body" }).click();
  await expect(page.getByRole("status")).toContainText("Body copied");

  await page.getByRole("button", { name: "Draft" }).click();
  await expect(page.getByLabel("Generated reply draft")).toBeVisible();
  await page.getByRole("button", { name: "Copy draft" }).click();
  await expect(page.getByRole("status")).toContainText("Draft copied");

  const downloadPromise = page.waitForEvent("download");
  await page.getByRole("button", { name: "Export state" }).click();
  const download = await downloadPromise;
  expect(download.suggestedFilename()).toBe("localhuman-mail-state.json");

  await page.getByRole("button", { name: "Share" }).click();
  await expect(page.getByRole("status")).toContainText("Share URL copied");
});

import { expect, test } from "@playwright/test";

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

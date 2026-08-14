import { expect, test } from "@playwright/test";

test("shows the sign in screen", async ({ page }) => {
  await page.goto("/login");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
});

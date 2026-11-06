import { expect, type Page, type TestInfo } from "@playwright/test";

export async function verifyTeamFlow(page: Page, testInfo: TestInfo) {
  const token = await page.evaluate(() => localStorage.getItem("shiyu_token"));
  const headers = { Authorization: `Bearer ${token}` };
  const user = await page.request.post("/api/v1/system/users", {
    headers,
    data: {
      username: "reviewer",
      nickname: "审批同事",
      password: "Shiyu!Review2026",
      status: 1,
    },
  });
  expect(user.status()).toBe(200);

  await page.goto("/system/settings");
  await page.getByLabel("系统名称", { exact: true }).fill("团队协作管理系统");
  await page.getByLabel("单文件上传上限（MiB）").fill("10");
  await page.getByRole("button", { name: "保存设置" }).click();
  await expect(page.getByText("系统设置已保存", { exact: true })).toBeVisible();
  await expect(
    page.getByRole("cell", {
      name: /系统名称：P0 验证管理系统 → 团队协作管理系统/,
    })
  ).toBeVisible();
  await page.evaluate(() => window.scrollTo(0, 0));
  await page.screenshot({
    path: testInfo.outputPath("settings.png"),
    fullPage: true,
  });

  await page.goto("/system/role");
  await page.getByRole("button", { name: /新建/ }).click();
  const roleDialog = page.getByRole("dialog", { name: "新建角色" });
  await roleDialog.getByRole("combobox", { name: "角色模板" }).click();
  await page.getByText("团队成员：本人数据与文件", { exact: true }).click();
  await roleDialog.getByLabel("角色编码", { exact: true }).fill("TEAM_MEMBER");
  await roleDialog.getByLabel("角色标识", { exact: true }).fill("team_member");
  await expect(roleDialog.getByLabel("角色名称", { exact: true })).toHaveValue(
    "团队成员"
  );
  await roleDialog.getByRole("button", { name: "确 定" }).click();
  await expect(roleDialog).not.toBeVisible();
  const permissions = await page.request.get(
    "/api/v1/system/roles/TEAM_MEMBER/menus",
    { headers }
  );
  expect(permissions.status()).toBe(200);
  expect(
    (await permissions.json()).data.map((m: { perms: string }) => m.perms)
  ).toContain("system:file:upload");

  await page.goto("/work/requests");
  await page.getByRole("button", { name: "新建申请", exact: true }).click();
  const dialog = page.getByRole("dialog", { name: "新建申请" });
  await dialog.getByLabel("事项标题").fill("采购办公设备");
  await dialog
    .getByLabel("申请说明")
    .fill("申请采购办公设备，预算与附件请审批。");
  await dialog.getByLabel("审批人", { exact: true }).click();
  await page.getByText("审批同事（reviewer）", { exact: true }).click();
  await dialog
    .locator("input[type=file]")
    .setInputFiles({
      name: "budget.txt",
      mimeType: "text/plain",
      buffer: Buffer.from("Budget: 500"),
    });
  await expect(dialog.getByText("budget.txt")).toBeVisible();
  await dialog.getByRole("button", { name: "提交申请", exact: true }).click();
  await expect(dialog).not.toBeVisible();
  await expect(
    page.getByRole("heading", { name: "采购办公设备", exact: true })
  ).toBeVisible();
  await page.screenshot({
    path: testInfo.outputPath("request.png"),
    fullPage: true,
  });

  const browser = page.context().browser()!;
  const context = await browser.newContext({
    baseURL: "http://127.0.0.1:18800",
  });
  try {
    const reviewer = await context.newPage();
    const errors: string[] = [];
    reviewer.on("response", (response) => {
      if (response.url().includes("/api/") && response.status() >= 400)
        errors.push(`${response.status()} ${response.url()}`);
    });
    await reviewer.goto("/user/login");
    await reviewer.getByPlaceholder("用户名", { exact: true }).fill("reviewer");
    await reviewer
      .getByPlaceholder("密码", { exact: true })
      .fill("Shiyu!Review2026");
    await reviewer.getByRole("button", { name: "进入控制台" }).click();
    await expect(reviewer).toHaveURL(/\/welcome$/);
    await expect(
      reviewer.getByText("采购办公设备", { exact: true })
    ).toBeVisible();
    await reviewer.screenshot({
      path: testInfo.outputPath("workspace.png"),
      fullPage: true,
    });
    await reviewer.goto("/work/notifications");
    await expect(
      reviewer.getByText("待审批：采购办公设备", { exact: true })
    ).toBeVisible();
    await reviewer.getByRole("button", { name: "查看事项" }).click();
    await expect(
      reviewer.getByRole("heading", { name: "采购办公设备", exact: true })
    ).toBeVisible();
    const download = reviewer.waitForEvent("download");
    await reviewer
      .getByRole("button", { name: "budget.txt", exact: true })
      .click();
    expect((await download).suggestedFilename()).toBe("budget.txt");
    await reviewer
      .getByRole("button", { name: "通过申请", exact: true })
      .click();
    const decision = reviewer.getByRole("dialog", {
      name: "通过申请",
      exact: true,
    });
    await decision.getByLabel("处理意见").fill("预算已确认，同意采购。");
    await decision.getByRole("button", { name: "确认处理" }).click();
    await expect(decision).not.toBeVisible();
    await expect(
      reviewer.getByRole("dialog", { name: "事项详情" }).getByText("预算已确认，同意采购。", { exact: true })
    ).toBeVisible();
    await reviewer.screenshot({
      path: testInfo.outputPath("approved.png"),
      fullPage: true,
    });
    expect(errors).toEqual([]);
  } finally {
    await context.close();
  }
  await page.goto("/work/notifications");
  await expect(
    page.getByText("申请已通过：采购办公设备", { exact: true })
  ).toBeVisible();
  await page.getByRole("button", { name: "查看事项" }).click();
  await expect(
    page.getByText("预算已确认，同意采购。", { exact: true })
  ).toBeVisible();
}

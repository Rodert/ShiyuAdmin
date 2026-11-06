import { expect, type Page, type TestInfo } from "@playwright/test";

export async function verifyMaintenanceFlow(page: Page, testInfo: TestInfo) {
  await page.goto("/work/requests");
  await page.getByRole("button", { name: "新建申请", exact: true }).click();
  let form = page.getByRole("dialog", { name: "新建申请", exact: true });
  await form
    .locator("input[type=file]")
    .setInputFiles({
      name: "unused-invoice.txt",
      mimeType: "text/plain",
      buffer: Buffer.from("invoice to reuse"),
    });
  await expect(form.getByText("unused-invoice.txt")).toBeVisible();
  await form.getByRole("button", { name: /取\s*消/ }).click();
  await page
    .getByRole("button", { name: "我的未使用附件", exact: true })
    .click();
  let drawer = page.getByRole("dialog", {
    name: "我的未使用附件",
    exact: true,
  });
  await expect(
    drawer.getByRole("cell", { name: "unused-invoice.txt", exact: true })
  ).toBeVisible();
  await drawer
    .getByRole("row")
    .filter({ hasText: "unused-invoice.txt" })
    .getByRole("button", { name: "移入回收站" })
    .click();
  await page
    .getByRole("dialog", { name: "附件移入回收站", exact: true })
    .getByRole("button", { name: "移入回收站" })
    .click();
  await expect(
    drawer.getByRole("cell", { name: "unused-invoice.txt", exact: true })
  ).not.toBeVisible();
  await drawer.getByRole("tab", { name: "回收站", exact: true }).click();
  await drawer
    .getByRole("row")
    .filter({ hasText: "unused-invoice.txt" })
    .getByRole("button", { name: "恢复附件" })
    .click();
  await expect(
    drawer.getByRole("cell", { name: "unused-invoice.txt", exact: true })
  ).not.toBeVisible();
  await page.goto("/work/requests");
  await page.getByRole("button", { name: "新建申请", exact: true }).click();
  form = page.getByRole("dialog", { name: "新建申请", exact: true });
  await form.getByLabel("事项标题").fill("设备费用报销");
  await form.getByLabel("申请说明").fill("请审核设备费用。");
  await form.getByLabel("审批人", { exact: true }).click();
  await page.getByText("审批同事（reviewer）", { exact: true }).click();
  await form.getByRole("button", { name: "选择已上传附件" }).click();
  drawer = page.getByRole("dialog", { name: "我的未使用附件", exact: true });
  await drawer
    .getByRole("row")
    .filter({ hasText: "unused-invoice.txt" })
    .getByRole("button", { name: "添加至申请" })
    .click();
  await expect(form.getByText("unused-invoice.txt")).toBeVisible();
  await form.getByRole("button", { name: "提交申请", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "设备费用报销", exact: true })
  ).toBeVisible();
  const originalID = new URL(page.url()).searchParams.get("id");
  const login = await page.request.post("/api/v1/system/auth/login", {
    data: { username: "reviewer", password: "Shiyu!Review2026" },
  });
  expect(login.status()).toBe(200);
  const reviewerToken = (await login.json()).data.token;
  const rejected = await page.request.post(
    `/api/v1/system/work/requests/${originalID}/actions`,
    {
      headers: { Authorization: `Bearer ${reviewerToken}` },
      data: { action: "rejected", comment: "请补充发票明细" },
    }
  );
  expect(rejected.status()).toBe(200);
  await page.reload();
  await page.getByRole("button", { name: "修改后重新提交" }).click();
  form = page.getByRole("dialog", { name: "重新提交申请", exact: true });
  await expect(form.getByLabel("事项标题")).toHaveValue("设备费用报销");
  await expect(form.getByText("unused-invoice.txt")).toBeVisible();
  await form.getByLabel("事项标题").fill("设备费用报销（补充发票）");
  await form.getByLabel("申请说明").fill("已补充发票明细，请再次审核。");
  await form.getByLabel("审批人", { exact: true }).click();
  await page.getByText("审批同事（reviewer）", { exact: true }).click();
  await form.getByRole("button", { name: "提交申请", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "设备费用报销（补充发票）", exact: true })
  ).toBeVisible();
  expect(new URL(page.url()).searchParams.get("id")).not.toBe(originalID);
  await page.getByRole("button", { name: "查看原申请" }).click();
  await expect(
    page.getByRole("heading", { name: "设备费用报销", exact: true })
  ).toBeVisible();
  await expect(page.getByText("请补充发票明细", { exact: true })).toBeVisible();
  await page.goto("/work/requests");
  await page.getByRole("searchbox", { name: "搜索事项" }).fill("补充发票");
  await page.getByRole("button", { name: /搜\s*索/, exact: true }).click();
  await page.getByRole("combobox", { name: "筛选事项状态" }).click();
  await page.getByText("待审批", { exact: true }).last().click();
  await expect(
    page.getByRole("button", { name: "设备费用报销（补充发票）", exact: true })
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "设备费用报销", exact: true })
  ).not.toBeVisible();
  await page.screenshot({
    path: testInfo.outputPath("filtered-requests.png"),
    fullPage: true,
  });
  await page
    .getByRole("button", { name: "设备费用报销（补充发票）", exact: true })
    .click();
  await page
    .getByRole("dialog", { name: "事项详情" })
    .getByRole("button", { name: "关闭", exact: true })
    .click();
  expect(new URL(page.url()).searchParams.get("keyword")).toBe("补充发票");
  expect(new URL(page.url()).searchParams.get("status")).toBe("pending");
  await page.goto("/system/files");
  await expect(page.getByText("有效文件", { exact: true })).toBeVisible();
  await page.getByRole("button", { name: "到期文件清理" }).click();
  await expect(
    page
      .getByRole("dialog", { name: "到期文件清理" })
      .getByRole("button", { name: "清理所选文件" })
  ).toBeDisabled();
  await page.screenshot({
    path: testInfo.outputPath("file-cleanup.png"),
    fullPage: true,
  });
}

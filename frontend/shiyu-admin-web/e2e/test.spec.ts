import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:8000';
const API_BASE = 'http://localhost:8080';

// 登录辅助函数
async function login(page: any) {
  await page.goto(`${BASE_URL}/user/login`, { waitUntil: 'domcontentloaded' });
  
  // 等待页面加载，尝试多种选择器
  await page.waitForLoadState('networkidle', { timeout: 60000 }).catch(() => {});
  
  // 等待任何输入框出现
  await page.waitForSelector('input, .ant-input, [name="username"], [name="password"]', { 
    timeout: 60000,
    state: 'visible'
  });
  
  // 尝试找到用户名输入框
  const usernameSelectors = [
    'input[name="username"]',
    '.ant-input[name="username"]',
    'input[placeholder*="用户名"]',
    'input[placeholder*="admin"]',
    'input[placeholder*="账户"]'
  ];
  
  let usernameInput = null;
  for (const selector of usernameSelectors) {
    try {
      usernameInput = page.locator(selector).first();
      await usernameInput.waitFor({ timeout: 5000, state: 'visible' });
      break;
    } catch {
      continue;
    }
  }
  
  if (!usernameInput) {
    throw new Error('无法找到用户名输入框');
  }
  
  // 尝试找到密码输入框
  const passwordSelectors = [
    'input[name="password"]',
    'input[type="password"]',
    '.ant-input-password input',
    '.ant-input[type="password"]'
  ];
  
  let passwordInput = null;
  for (const selector of passwordSelectors) {
    try {
      passwordInput = page.locator(selector).first();
      await passwordInput.waitFor({ timeout: 5000, state: 'visible' });
      break;
    } catch {
      continue;
    }
  }
  
  if (!passwordInput) {
    throw new Error('无法找到密码输入框');
  }
  
  // 填写表单
  await usernameInput.fill('admin');
  await passwordInput.fill('Admin@123');
  
  // 点击登录按钮
  const submitSelectors = [
    'button.ant-btn-primary',
    'button[type="submit"]',
    'button:has-text("登录")',
    'button:has-text("登 录")'
  ];
  
  let submitButton = null;
  for (const selector of submitSelectors) {
    try {
      submitButton = page.locator(selector).first();
      await submitButton.waitFor({ timeout: 5000, state: 'visible' });
      break;
    } catch {
      continue;
    }
  }
  
  if (!submitButton) {
    throw new Error('无法找到登录按钮');
  }
  
  await submitButton.click();
  
  // 等待跳转
  await page.waitForURL(/\/(welcome|\?redirect=)/, { timeout: 20000 });
}

test.describe('Shiyu Admin E2E Tests', () => {
  // 移除 beforeEach，让每个测试自己处理导航

  test('登录功能', async ({ page }) => {
    await login(page);
    
    // 验证登录成功
    const currentUrl = page.url();
    expect(currentUrl).toMatch(/\/(welcome|\?redirect=)/);
  });

  test('用户管理 - 查看列表', async ({ page }) => {
    // 先登录
    await login(page);
    
    // 导航到用户管理页面
    await page.goto(`${BASE_URL}/system/user`);
    await page.waitForLoadState('networkidle');
    
    // 等待页面加载完成
    await page.waitForTimeout(2000);
    
    // 验证页面标题或表格存在
    const hasTitle = await page.locator('text=用户管理').isVisible().catch(() => false);
    const hasTable = await page.locator('table, .ant-table').isVisible().catch(() => false);
    
    expect(hasTitle || hasTable).toBeTruthy();
  });

  test('用户管理 - 创建用户', async ({ page }) => {
    // 先登录
    await login(page);
    
    // 导航到用户管理页面
    await page.goto(`${BASE_URL}/system/user`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // 点击新建按钮 - 使用更灵活的选择器
    const newButton = page.locator('button:has-text("新建"), button:has-text("新增"), .ant-btn-primary:has-text("新建")').first();
    await newButton.waitFor({ timeout: 10000 });
    await newButton.click();
    
    // 等待表单弹窗
    await page.waitForSelector('.ant-modal, .ant-drawer', { timeout: 10000 });
    await page.waitForTimeout(1000);
    
    // 填写表单
    const timestamp = Date.now();
    const username = `testuser${timestamp}`;
    await page.fill('input[name="user_code"], input[placeholder*="用户编码"]', `TEST-${timestamp}`);
    await page.fill('input[name="username"], input[placeholder*="用户名"]', username);
    await page.fill('input[name="password"], input[placeholder*="密码"]', 'Test123456!');
    await page.fill('input[name="nickname"], input[placeholder*="昵称"]', '测试用户');
    
    // 提交表单
    await page.click('button:has-text("确定"), button:has-text("提交"), .ant-btn-primary:has-text("确定")');
    
    // 在表格中验证新用户出现
    const userRow = page.locator(`.ant-table-row td:has-text("${username}")`).first();
    await expect(userRow).toBeVisible({ timeout: 15000 });
  });

  test('角色管理 - 查看列表', async ({ page }) => {
    // 先登录
    await login(page);
    
    // 导航到角色管理页面
    await page.goto(`${BASE_URL}/system/role`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // 验证页面标题或内容
    const hasTitle = await page.locator('text=角色管理').isVisible().catch(() => false);
    const hasContent = await page.locator('table, .ant-table, .ant-pro-table').isVisible().catch(() => false);
    
    expect(hasTitle || hasContent).toBeTruthy();
  });

  test('菜单管理 - 查看列表', async ({ page }) => {
    // 先登录
    await login(page);
    
    // 导航到菜单管理页面
    await page.goto(`${BASE_URL}/system/menu`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // 验证页面标题或内容
    const hasTitle = await page.locator('text=菜单管理').isVisible().catch(() => false);
    const hasContent = await page.locator('table, .ant-table, .ant-pro-table').isVisible().catch(() => false);
    
    expect(hasTitle || hasContent).toBeTruthy();
  });

  test('部门管理 - 查看列表', async ({ page }) => {
    // 先登录
    await login(page);
    
    // 导航到部门管理页面
    await page.goto(`${BASE_URL}/system/dept`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // 验证页面标题或内容
    const hasTitle = await page.locator('text=部门管理').isVisible().catch(() => false);
    const hasContent = await page.locator('table, .ant-table, .ant-pro-table').isVisible().catch(() => false);
    
    expect(hasTitle || hasContent).toBeTruthy();
  });

  test('用户与角色 - 新建并在编辑表单中交叉验证', async ({ page }) => {
    await login(page);

    const timestamp = Date.now();
    const roleCode = `ROLE-${timestamp}`;
    const roleName = `测试角色${timestamp}`;
    const roleKey = `test_role_${timestamp}`;
    const userCode = `TEST-${timestamp}`;
    const username = `testuser${timestamp}`;
    const password = 'Test123456!';

    // 先新建角色
    await page.goto(`${BASE_URL}/system/role`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const newRoleButton = page
      .locator('button:has-text("新建"), button:has-text("新增"), .ant-btn-primary:has-text("新建")')
      .first();
    await newRoleButton.waitFor({ timeout: 10000 });
    await newRoleButton.click();

    await page.waitForSelector('.ant-modal, .ant-drawer', { timeout: 10000 });
    await page.waitForTimeout(500);

    await page.fill('input[name="role_code"], input[placeholder*="角色编码"]', roleCode);
    await page.fill('input[name="role_name"], input[placeholder*="角色名称"]', roleName);
    await page.fill('input[name="role_key"], input[placeholder*="角色标识"]', roleKey);

    await page.click(
      'button:has-text("确定"), button:has-text("提交"), .ant-btn-primary:has-text("确定")',
    );

    const roleRow = page.locator(`.ant-table-row td:has-text("${roleName}")`).first();
    await expect(roleRow).toBeVisible({ timeout: 15000 });

    // 再新建用户
    await page.goto(`${BASE_URL}/system/user`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const newUserButton = page
      .locator('button:has-text("新建"), button:has-text("新增"), .ant-btn-primary:has-text("新建")')
      .first();
    await newUserButton.waitFor({ timeout: 10000 });
    await newUserButton.click();

    await page.waitForSelector('.ant-modal, .ant-drawer', { timeout: 10000 });
    await page.waitForTimeout(500);

    await page.fill('input[name="user_code"], input[placeholder*="用户编码"]', userCode);
    await page.fill('input[name="username"], input[placeholder*="用户名"]', username);
    await page.fill('input[name="password"], input[placeholder*="密码"]', password);
    await page.fill('input[name="nickname"], input[placeholder*="昵称"]', '测试用户');

    await page.click(
      'button:has-text("确定"), button:has-text("提交"), .ant-btn-primary:has-text("确定")',
    );

    const userRowInTable = page.locator(`.ant-table-row td:has-text("${username}")`).first();
    await expect(userRowInTable).toBeVisible({ timeout: 15000 });

    // 在用户编辑表单中验证新建角色存在
    const editButton = page
      .locator(`tr:has-text("${username}") button:has-text("编辑")`)
      .first();
    await editButton.waitFor({ timeout: 10000 });
    await editButton.click();

    await page.waitForSelector('.ant-modal, .ant-drawer', { timeout: 10000 });
    await page.waitForTimeout(500);

    // 打开角色选择下拉
    const roleSelectTrigger = page
      .locator('.ant-modal .ant-form-item:has-text("角色") .ant-select-selector')
      .first();
    await roleSelectTrigger.click();

    const roleOption = page
      .locator(
        `.ant-select-dropdown .ant-select-item-option[title="${roleName}"], .ant-select-dropdown .ant-select-item-option:has-text("${roleName}")`,
      )
      .first();
    await expect(roleOption).toBeVisible({ timeout: 10000 });
  });
});


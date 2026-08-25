import { test, expect } from '@playwright/test';

async function waitForApplication(page) {
  await page.goto('/');
  await expect(page.locator('#password-output')).not.toHaveText('');
  await expect(page.locator('#validation-error')).toBeHidden();
}

test.describe('PassForge generator', () => {
  test('generates a password and updates its strength estimate', async ({ page }) => {
    await waitForApplication(page);

    const initial = await page.locator('#password-output').textContent();
    expect(initial).toHaveLength(20);
    await expect(page.locator('#strength-label')).not.toHaveText('—');
    await expect(page.locator('#entropy-bits')).toContainText('bits estimated entropy');
    await expect(page.locator('#strength-bar')).toHaveAttribute('data-level', /.+/);

    await page.locator('#input-length').fill('21');
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.locator('#password-output')).toHaveText(/^.{21}$/);
  });

  test('supports strong password, passphrase, and PIN modes', async ({ page }) => {
    await waitForApplication(page);

    await page.locator('#mode-strong').check();
    await expect(page.locator('#panel-strong')).toBeVisible();
    await expect(page.locator('#panel-random')).toBeHidden();
    await page.locator('#input-strong-length').fill('32');
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.locator('#password-output')).toHaveText(/^[\s\S]{32}$/);

    await page.locator('#mode-passphrase').check();
    await expect(page.locator('#panel-passphrase')).toBeVisible();
    await page.locator('#input-words').fill('4');
    await page.locator('#input-separator').fill('.');
    await page.locator('#chk-add-number').uncheck();
    await page.getByRole('button', { name: 'Generate' }).click();
    const passphrase = await page.locator('#password-output').textContent();
    expect(passphrase.split('.')).toHaveLength(4);

    await page.locator('#mode-pin').check();
    await expect(page.locator('#panel-pin')).toBeVisible();
    await page.locator('#input-pin-length').fill('8');
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.locator('#password-output')).toHaveText(/^\d{8}$/);
    await expect(page.locator('#panel-pin .warning')).toContainText('much less entropy');
  });

  test('validates options, exposes advanced controls, and resets defaults', async ({ page }) => {
    await waitForApplication(page);

    await page.getByRole('button', { name: /Advanced Options/ }).click();
    await expect(page.locator('#advanced-panel')).toBeVisible();
    await expect(page.locator('#btn-advanced-toggle')).toHaveAttribute('aria-expanded', 'true');

    await page.locator('#input-excluded').fill('aA');
    await page.locator('#chk-exclude-ambiguous').check();
    await page.locator('#chk-prevent-repeated').check();
    await page.getByRole('button', { name: 'Generate' }).click();
    const configuredPassword = await page.locator('#password-output').textContent();
    expect(configuredPassword).toHaveLength(20);
    expect(configuredPassword).not.toMatch(/[aA0Oo1Il]/);
    expect(new Set(configuredPassword).size).toBe(configuredPassword.length);

    for (const id of ['chk-upper', 'chk-lower', 'chk-numbers', 'chk-symbols']) {
      await page.locator(`#${id}`).uncheck();
    }
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.locator('#validation-error')).toBeVisible();
    await expect(page.locator('#validation-error')).toContainText('character type');
    await expect(page.getByRole('button', { name: 'Generate' })).toBeEnabled();

    for (const id of ['chk-upper', 'chk-lower', 'chk-numbers', 'chk-symbols']) {
      await page.locator(`#${id}`).check();
    }
    await page.locator('#input-length').fill('21');
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect(page.locator('#validation-error')).toBeHidden();
    await expect(page.locator('#password-output')).toHaveText(/^.{21}$/);

    await page.getByRole('button', { name: 'Reset to defaults' }).click();
    await expect(page.locator('#mode-random')).toBeChecked();
    await expect(page.locator('#panel-random')).toBeVisible();
    await expect(page.locator('#validation-error')).toBeHidden();
    await expect(page.locator('#password-output')).toHaveText(/\S/);
    await expect(page.locator('#chk-upper')).toBeChecked();
    await expect(page.locator('#chk-lower')).toBeChecked();
    await expect(page.locator('#chk-numbers')).toBeChecked();
    await expect(page.locator('#chk-symbols')).toBeChecked();
  });

  test('changes theme, copies the generated password, and supports keyboard generation', async ({ page, context }) => {
    await context.grantPermissions(['clipboard-read', 'clipboard-write'], { origin: 'http://localhost:8080' });
    await waitForApplication(page);

    await page.locator('#theme-select').selectOption('dark');
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
    await page.locator('#theme-select').selectOption('light');
    await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');

    const password = await page.locator('#password-output').textContent();
    await page.getByRole('button', { name: 'Copy' }).click();
    await expect(page.getByRole('button', { name: 'Copied!' })).toBeVisible();
    expect(await page.evaluate(() => navigator.clipboard.readText())).toBe(password);

    await page.locator('#input-length').fill('21');
    await page.locator('body').press('Control+Enter');
    await expect(page.locator('#password-output')).toHaveText(/^.{21}$/);
  });
});
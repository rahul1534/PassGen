import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

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

    await page.getByRole('button', { name: /Advanced options/ }).click();
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

test.describe('PassForge privacy guarantees (runtime)', () => {
  test('makes no cross-origin requests and raises no CSP violations', async ({ page, baseURL }) => {
    const origin = new URL(baseURL).origin;
    const crossOrigin = [];
    const problems = [];

    page.on('request', (request) => {
      const url = request.url();
      if (!url.startsWith(origin) && !url.startsWith('data:') && !url.startsWith('blob:')) {
        crossOrigin.push(url);
      }
    });
    // Chromium reports CSP violations (e.g. eval blocked because 'unsafe-eval' is
    // absent) as console errors; uncaught exceptions surface via pageerror.
    page.on('console', (msg) => {
      if (msg.type() === 'error') problems.push(`console: ${msg.text()}`);
    });
    page.on('pageerror', (err) => problems.push(`pageerror: ${err.message}`));

    await waitForApplication(page);
    for (const mode of ['#mode-strong', '#mode-passphrase', '#mode-pin', '#mode-random']) {
      await page.locator(mode).check();
      await page.getByRole('button', { name: 'Generate' }).click();
      await expect(page.locator('#password-output')).not.toHaveText('');
    }
    await page.waitForLoadState('networkidle');

    expect(crossOrigin, 'no request may leave the origin').toEqual([]);
    expect(problems, 'no CSP violations or script errors').toEqual([]);
  });

  test('ships the approved strict CSP', async ({ page }) => {
    await page.goto('/');
    const csp = await page
      .locator('meta[http-equiv="Content-Security-Policy"]')
      .getAttribute('content');

    expect(csp).toContain("connect-src 'self';");
    expect(csp).toContain("script-src 'self' 'wasm-unsafe-eval';");
    // Plain 'unsafe-eval' must be absent ('wasm-unsafe-eval' is fine and required).
    expect(csp).not.toMatch(/(^|\s)'unsafe-eval'/);
    expect(csp).not.toMatch(/https?:/);
  });

  test('rejects negative minimum counts instead of over-long output', async ({ page }) => {
    await waitForApplication(page);
    await page.getByRole('button', { name: /Advanced options/ }).click();

    await page.locator('#input-length').fill('10');
    await page.locator('#input-min-upper').fill('30');
    await page.locator('#input-min-lower').fill('-20');
    await page.getByRole('button', { name: 'Generate' }).click();

    await expect(page.locator('#validation-error')).toBeVisible();
    await expect(page.locator('#validation-error')).toContainText('cannot be negative');
    await expect(page.locator('#password-output')).toHaveText('');
  });
});


test.describe('PassForge interface', () => {
  test('regenerates as soon as an option changes', async ({ page }) => {
    await waitForApplication(page);
    await page.locator('#chk-symbols').uncheck();
    await expect(page.locator('#password-output')).toHaveText(/^[A-Za-z0-9]{20}$/);
    await page.locator('#chk-numbers').uncheck();
    await expect(page.locator('#password-output')).toHaveText(/^[A-Za-z]{20}$/);
  });

  test('keeps the length slider and number field in step', async ({ page }) => {
    await waitForApplication(page);

    await page.locator('#input-length-range').fill('32');
    await expect(page.locator('#input-length')).toHaveValue('32');
    await expect(page.locator('#password-output')).toHaveText(/^.{32}$/);

    await page.locator('#input-length').fill('40');
    await page.locator('#input-length').press('Tab');
    await expect(page.locator('#input-length-range')).toHaveValue('40');
    await expect(page.locator('#password-output')).toHaveText(/^.{40}$/);

    await page.getByRole('button', { name: 'Reset to defaults' }).click();
    await expect(page.locator('#input-length')).toHaveValue('20');
    await expect(page.locator('#input-length-range')).toHaveValue('20');
  });

  test('colours digits and symbols separately without changing the text', async ({ page }) => {
    await waitForApplication(page);
    await page.locator('#input-length').fill('64');
    await page.locator('#input-length').press('Tab');
    await expect(page.locator('#password-output')).toHaveText(/^.{64}$/);

    const parts = await page.locator('#password-output > span').evaluateAll((spans) =>
      spans.map((s) => ({ cls: s.className, text: s.textContent })));
    expect(parts.length).toBeGreaterThan(1);
    for (const { cls, text } of parts) {
      if (cls === 'ch-digit') expect(text).toMatch(/^[0-9]+$/);
      else if (cls === 'ch-symbol') expect(text).toMatch(/^[^A-Za-z0-9]+$/);
      else expect(text).toMatch(/^[A-Za-z]+$/);
    }
    const joined = parts.map((p) => p.text).join('');
    expect(await page.locator('#password-output').textContent()).toBe(joined);
    expect(await page.locator('#password-output').evaluate((el) => el.innerHTML)).not.toContain('<script');
  });

  test('shows an empty strength meter when there is nothing to rate', async ({ page }) => {
    await waitForApplication(page);
    for (const id of ['chk-upper', 'chk-lower', 'chk-numbers', 'chk-symbols']) {
      await page.locator(`#${id}`).uncheck();
    }
    await expect(page.locator('#validation-error')).toBeVisible();
    await expect(page.locator('#strength-label')).toHaveText('—');
    await expect(page.locator('#entropy-bits')).toHaveText('');
    await expect(page.locator('#strength-bar')).toHaveAttribute('data-level', '');
    await expect(page.locator('#strength-track')).toHaveAttribute('aria-valuenow', '0');
  });

  test('exposes the strength meter to assistive technology', async ({ page }) => {
    await waitForApplication(page);
    await expect(page.locator('#strength-track')).toHaveAttribute('aria-valuenow', /^[1-9]\d*$/);
    await expect(page.locator('#strength-track')).toHaveAttribute('aria-valuetext', /\S/);
  });

  test('mode tabs show a visible keyboard focus ring', async ({ page }) => {
    await waitForApplication(page);
    await page.locator('#mode-random').focus();
    await page.keyboard.press('ArrowRight'); // radio groups move focus with arrow keys
    await expect(page.locator('#mode-strong')).toBeFocused();
    await expect(page.locator('.mode:has(#mode-strong)')).toHaveCSS('outline-style', 'solid');
  });

  test('character chips show a visible keyboard focus ring', async ({ page }) => {
    await waitForApplication(page);
    await page.locator('#chk-upper').focus();
    await page.keyboard.press('Shift+Tab');
    await page.keyboard.press('Tab');
    await expect(page.locator('#chk-upper')).toBeFocused();
    await expect(page.locator('.chip:has(#chk-upper)')).toHaveCSS('outline-style', 'solid');
  });

  test('marks the copy button while the copy confirmation is showing', async ({ page, context }) => {
    await context.grantPermissions(['clipboard-read', 'clipboard-write']);
    await waitForApplication(page);
    await page.getByRole('button', { name: 'Copy' }).click();
    await expect(page.locator('#btn-copy')).toHaveAttribute('data-state', 'copied');
    await expect(page.locator('#copy-announcer')).toHaveText('Copied to clipboard.');
    await expect(page.locator('#btn-copy')).not.toHaveAttribute('data-state', 'copied', { timeout: 5000 });
    await expect(page.locator('#btn-copy')).toHaveText('Copy');
  });

  test('has no horizontal overflow on a phone-sized screen', async ({ page }) => {
    await page.setViewportSize({ width: 360, height: 740 });
    await waitForApplication(page);
    await page.locator('#input-length').fill('128');
    await page.locator('#input-length').press('Tab');
    await expect(page.locator('#password-output')).toHaveText(/^.{128}$/);
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);
    expect(overflow).toBeLessThanOrEqual(0);
  });
});


// Automated WCAG 2.x A/AA checks (contrast, names, roles, landmarks) in both
// colour schemes and every mode, plus the transient copied and error states.
for (const scheme of ['light', 'dark']) {
  test.describe(`PassForge accessibility (${scheme})`, () => {
    test.use({ colorScheme: scheme, reducedMotion: 'reduce' });

    const audit = async (page) => {
      const results = await new AxeBuilder({ page })
        .withTags(['wcag2a', 'wcag2aa', 'wcag21aa', 'best-practice'])
        .analyze();
      const summary = results.violations.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(', ')}`);
      expect(summary, 'accessibility violations').toEqual([]);
    };

    for (const mode of ['random', 'strong', 'passphrase', 'pin']) {
      test(`${mode} mode has no violations`, async ({ page }) => {
        await waitForApplication(page);
        await page.locator(`#mode-${mode}`).check();
        if (mode === 'random') await page.getByRole('button', { name: /Advanced options/ }).click();
        await audit(page);
      });
    }

    test('copied and error states have no violations', async ({ page, context }) => {
      await context.grantPermissions(['clipboard-read', 'clipboard-write']);
      await waitForApplication(page);
      await page.getByRole('button', { name: 'Copy' }).click();
      await expect(page.locator('#btn-copy')).toHaveAttribute('data-state', 'copied');
      await audit(page);

      for (const id of ['chk-upper', 'chk-lower', 'chk-numbers', 'chk-symbols']) {
        await page.locator(`#${id}`).uncheck();
      }
      await expect(page.locator('#validation-error')).toBeVisible();
      await audit(page);
    });
  });
}

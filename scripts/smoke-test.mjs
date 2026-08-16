import { chromium } from 'playwright';

const port = process.env.PORT || '8083';
const errors = [];

const browser = await chromium.launch();
const page = await browser.newPage();
page.on('console', (msg) => {
  if (msg.type() === 'error') errors.push(msg.text());
});
page.on('pageerror', (err) => errors.push(String(err)));

await page.goto(`http://localhost:${port}/`, { waitUntil: 'networkidle' });
await page.waitForTimeout(3000);

const initial = await page.locator('#password-output').textContent();
const validation = await page.locator('#validation-error').textContent();
const strength = await page.locator('#strength-label').textContent();
const entropy = await page.locator('#entropy-bits').textContent();

await page.click('#btn-generate');
await page.waitForTimeout(500);
const regenerated = await page.locator('#password-output').textContent();

await page.click('#mode-passphrase');
await page.waitForTimeout(800);
const passphrase = await page.locator('#password-output').textContent();

await page.click('#mode-pin');
await page.waitForTimeout(800);
const pin = await page.locator('#password-output').textContent();
const pinWarning = await page.locator('#panel-pin .warning').textContent();

const result = {
  initial,
  regenerated,
  passphrase,
  pin,
  validation,
  strength,
  entropy,
  pinWarning: (pinWarning || '').trim().slice(0, 80),
  errors,
  ok:
    !!initial?.trim() &&
    !!regenerated?.trim() &&
    regenerated !== initial &&
    !!passphrase?.includes('-') &&
    /^\d{4,}$/.test(pin || '') &&
    !(validation || '').trim() &&
    errors.length === 0,
};

console.log(JSON.stringify(result, null, 2));
await browser.close();
process.exit(result.ok ? 0 : 1);

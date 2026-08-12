import { chromium } from 'playwright';

const port = process.env.PORT || '8082';
const errors = [];

const browser = await chromium.launch();
const page = await browser.newPage();
page.on('console', (msg) => {
  if (msg.type() === 'error') errors.push(msg.text());
});
page.on('pageerror', (err) => errors.push(String(err)));

await page.goto(`http://localhost:${port}/`, { waitUntil: 'networkidle' });
await page.waitForTimeout(4000);

const password = await page.locator('#password-output').textContent();
const validation = await page.locator('#validation-error').textContent();
const strength = await page.locator('#strength-label').textContent();

console.log(JSON.stringify({ password, validation, strength, errors }, null, 2));

await browser.close();

if (!password || password.trim().length === 0) {
  process.exit(1);
}

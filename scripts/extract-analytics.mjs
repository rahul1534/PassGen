#!/usr/bin/env node

/**
 * Extract analytics data from PassForge using Playwright
 * This script loads the live website and extracts localStorage analytics data
 */

import { chromium } from 'playwright';
import * as fs from 'fs';
import * as path from 'path';

const SITE_URL = process.env.SITE_URL || 'https://rahul1534.github.io/PassGen/';
const OUTPUT_DIR = process.env.OUTPUT_DIR || './analytics-data';
const ANALYTICS_DASHBOARD_URL = SITE_URL.replace(/\/$/, '') + '/analytics-dashboard.html';

async function extractAnalytics() {
  const browser = await chromium.launch();
  const context = await browser.createBrowserContext();
  const page = await context.newPage();

  try {
    console.log(`Loading analytics dashboard from: ${ANALYTICS_DASHBOARD_URL}`);
    await page.goto(ANALYTICS_DASHBOARD_URL, { waitUntil: 'networkidle' });

    // Wait for analytics data to be loaded and processed
    await page.waitForFunction(() => {
      const el = document.getElementById('json-output');
      return el && el.textContent && el.textContent !== 'Loading...';
    }, { timeout: 10000 });

    // Extract the JSON data from the page
    const jsonData = await page.evaluate(() => {
      const el = document.getElementById('json-output');
      return el ? el.textContent : null;
    });

    if (!jsonData || jsonData === 'No analytics data found.') {
      console.log('No analytics data found on the website');
      return null;
    }

    const analyticsData = JSON.parse(jsonData);

    // Also extract the summary stats
    const summary = await page.evaluate(() => {
      return {
        total_events: document.getElementById('stat-events')?.textContent || '0',
        page_views: document.getElementById('stat-views')?.textContent || '0',
        passwords_generated: document.getElementById('stat-generated')?.textContent || '0',
        last_activity: document.getElementById('stat-last')?.textContent || '—',
      };
    });

    console.log('Analytics extracted successfully:');
    console.log(`  - Total events: ${summary.total_events}`);
    console.log(`  - Page views: ${summary.page_views}`);
    console.log(`  - Passwords generated: ${summary.passwords_generated}`);
    console.log(`  - Last activity: ${summary.last_activity}`);

    return {
      data: analyticsData,
      summary: summary,
      extracted_at: new Date().toISOString()
    };
  } catch (error) {
    console.error('Error extracting analytics:', error);
    throw error;
  } finally {
    await context.close();
    await browser.close();
  }
}

async function saveAnalytics(analyticsData) {
  // Create output directory if it doesn't exist
  if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  }

  // Save full analytics data
  const dataPath = path.join(OUTPUT_DIR, 'analytics.json');
  fs.writeFileSync(dataPath, JSON.stringify(analyticsData.data, null, 2));
  console.log(`✓ Analytics data saved to ${dataPath}`);

  // Save summary
  const summaryPath = path.join(OUTPUT_DIR, 'summary.json');
  fs.writeFileSync(summaryPath, JSON.stringify({
    summary: analyticsData.summary,
    extracted_at: analyticsData.extracted_at
  }, null, 2));
  console.log(`✓ Summary saved to ${summaryPath}`);

  // Save as CSV for easy import to spreadsheets
  const csvPath = path.join(OUTPUT_DIR, 'events.csv');
  const events = analyticsData.data.events || [];
  const csv = [
    'timestamp,type,mode,referrer,user_agent_short',
    ...events.map(e => `"${e.timestamp}","${e.type}","${e.metadata?.mode || ''}","${e.referrer}","${e.user_agent?.substring(0, 50) || ''}"`)
  ].join('\n');
  fs.writeFileSync(csvPath, csv);
  console.log(`✓ Events exported to CSV: ${csvPath}`);

  // Save index with metadata
  const indexPath = path.join(OUTPUT_DIR, 'index.json');
  fs.writeFileSync(indexPath, JSON.stringify({
    extracted_at: analyticsData.extracted_at,
    files: {
      analytics: 'Full analytics data with all events',
      summary: 'Summary statistics',
      events_csv: 'Events exported as CSV'
    }
  }, null, 2));
  console.log(`✓ Index saved to ${indexPath}`);
}

async function main() {
  try {
    console.log('PassForge Analytics Extractor');
    console.log('=============================\n');

    const analyticsData = await extractAnalytics();

    if (analyticsData) {
      await saveAnalytics(analyticsData);
      console.log('\n✓ Analytics extraction completed successfully');
      process.exit(0);
    } else {
      console.log('\n⚠ No analytics data to save');
      process.exit(0);
    }
  } catch (error) {
    console.error('\n✗ Analytics extraction failed:', error.message);
    process.exit(1);
  }
}

main();

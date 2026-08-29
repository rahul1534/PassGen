# PassForge Analytics Setup Guide

## Overview

PassForge now includes a **privacy-first, completely free analytics system** that:
- ✅ Tracks usage **locally in your browser** (localStorage)
- ✅ Stores no sensitive data
- ✅ Requires no external services
- ✅ Automatically exports data to your private GitHub repo daily
- ✅ Complies with privacy laws (GDPR, CCPA)

## How It Works

```
User visits PassForge
    ↓
JavaScript tracks events in browser localStorage
    ↓
Analytics Dashboard displays live stats (/analytics-dashboard.html)
    ↓
GitHub Actions runs daily
    ↓
Extracts data using Playwright
    ↓
Commits to your private analytics repo
    ↓
You can analyze historical data anytime
```

## Setup Steps

### Step 1: Enable Analytics Collection (Already Done ✓)

Analytics collection is automatically enabled. The website now:
- Tracks page views on load
- Tracks password generations (method and timestamp)
- Stores data in browser localStorage

**No user action needed** — tracking is automatic and transparent.

### Step 2: Create a Private Analytics Repository

Create a new **private** GitHub repository to store analytics data:

1. Go to [github.com/new](https://github.com/new)
2. Repository name: `passgen-analytics` (or any name you prefer)
3. Select **Private**
4. Initialize with a README
5. Create the repository

### Step 3: Create a GitHub Token

GitHub Actions needs a token to upload analytics data:

1. Go to [github.com/settings/tokens](https://github.com/settings/tokens?type=pat)
2. Click **"Generate new token"** → **"Generate new token (classic)"**
3. Name: `PassGen Analytics`
4. Scopes: Check only **`repo`** (full control of private repositories)
5. Expiration: 90 days (auto-renewal recommended)
6. Click **"Generate token"**
7. **Copy the token** (you won't see it again!)

### Step 4: Add Token to PassGen Repository Secrets

1. Go to your **PassGen** repository (not the analytics repo)
2. Settings → Secrets and variables → Actions
3. Click **"New repository secret"**
   - Name: `ANALYTICS_REPO_TOKEN`
   - Value: [Paste your GitHub token from Step 3]
4. Click **"Add secret"**

### Step 5: Configure the Workflow (Optional)

Edit `.github/workflows/extract-analytics.yml` to customize:

**Change extraction schedule:**
```yaml
schedule:
  - cron: '0 2 * * *'  # Daily at 2 AM UTC
  # Alternatives:
  # - cron: '0 0 * * 0'  # Weekly on Sundays at midnight
  # - cron: '0 * * * *'  # Every hour
```

**Change analytics repo name:**
Edit line in workflow if you used a different repo name.

### Step 6: Test the Workflow

1. Go to your PassGen repo
2. Actions tab → **Extract Analytics** workflow
3. Click **"Run workflow"** → **"Run workflow"** (manual trigger)
4. Watch it run and commit analytics to your private repo

## Accessing Analytics

### View Live Analytics Dashboard

Open in your browser:
```
https://rahul1534.github.io/PassGen/analytics-dashboard.html
```

Shows:
- ✓ Total page views
- ✓ Password generations by method
- ✓ Referrer sources
- ✓ User agent info
- ✓ Full JSON export

### Download Analytics Data

From the dashboard, you can:
- **Download JSON** — Export as JSON file
- **Copy to Clipboard** — For pasting into spreadsheets
- **Clear Data** — Reset all analytics

### Access Historical Data in GitHub

Your private `passgen-analytics` repo contains:

```
analytics-repo/
├── latest/
│   ├── analytics.json      # Current snapshot
│   └── summary.json        # Current stats
├── history/
│   ├── 2024-01-15-02-00-00/
│   │   ├── analytics.json
│   │   └── summary.json
│   └── ...
└── README.md
```

## What Gets Tracked

For **each event**:
- ⏱️ Timestamp (ISO 8601)
- 📍 Event type (page_view, password_generated)
- 🔧 Generation method (random, strong, passphrase, pin) — if applicable
- 📊 Referrer (where user came from)
- 📱 User agent (browser info)

**What's NOT tracked:**
- ❌ Generated passwords (they never leave your browser!)
- ❌ IP addresses
- ❌ Personal information
- ❌ Anything sent to external servers
- ❌ Cookies or tracking pixels

## Analytics Dashboard Features

### Real-Time Stats
```
┌─────────────────┬──────────────┐
│  Total Events   │   Page Views  │
│       42        │      15       │
├─────────────────┴──────────────┤
│  Passwords Generated │   27     │
├──────────────────────┴──────────┤
│   Last Activity: Today 2:30 PM  │
└─────────────────────────────────┘
```

### Generation Breakdown
Shows password generations by method:
- Random: 12
- Strong: 8
- Passphrase: 6
- PIN: 1

### Export Options
1. **Download JSON** — Save as `passforge-analytics-2024-01-15.json`
2. **Copy to Clipboard** — Paste into Excel/Google Sheets
3. **View Raw JSON** — Inspect all fields

## Troubleshooting

### "No analytics data found"
- **Cause:** No events tracked yet
- **Fix:** Visit the PassForge homepage and generate a password

### GitHub Actions workflow fails
- **Check:** Is `ANALYTICS_REPO_TOKEN` secret added?
- **Check:** Is the private analytics repo created?
- **Check:** Does the token still have valid scopes?
- **Fix:** View workflow logs: Actions → Extract Analytics → failed run

### Analytics data not appearing in private repo
- **Cause:** Token might not have repo access
- **Fix:** Recreate token with `repo` scope checked

### Want to disable analytics
1. Remove tracking script from `web/index.html`:
   ```html
   <!-- Remove this line: -->
   <script src="./analytics.js"></script>
   ```
2. Disable the workflow: `.github/workflows/extract-analytics.yml` → disable

## Advanced: Local Backup

Manually export analytics anytime:

1. Visit: `https://rahul1534.github.io/PassGen/analytics-dashboard.html`
2. Click **"Download JSON"**
3. Save locally

## Advanced: Analyze in Excel/Google Sheets

1. Visit dashboard → Click **"Copy to Clipboard"**
2. Paste into Google Sheets / Excel
3. Use built-in charts and analysis

## Security Notes

- ✅ All data stored locally in browser (localStorage)
- ✅ No data sent to external analytics services
- ✅ GitHub token stored as encrypted repository secret
- ✅ Private analytics repo is not accessible to the public
- ✅ Analytics data is your own — you own and control it

## Questions?

- 📖 See [docs/security.md](../docs/security.md) for PassForge security details
- 🐛 Report issues: [GitHub Issues](https://github.com/rahul1534/PassGen/issues)

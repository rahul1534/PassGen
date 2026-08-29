/**
 * PassForge Analytics - Privacy-first tracking stored in localStorage
 * Data is stored locally and never sent to external servers
 */

const PassForgeAnalytics = (() => {
  const STORAGE_KEY = 'passforge_analytics';
  const VERSION = '1';
  
  // Initialize analytics on first load or retrieve existing data
  const initStorage = () => {
    if (!localStorage.getItem(STORAGE_KEY)) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        version: VERSION,
        created_at: new Date().toISOString(),
        events: []
      }));
    }
  };

  // Get all stored analytics data
  const getData = () => {
    try {
      const data = localStorage.getItem(STORAGE_KEY);
      return data ? JSON.parse(data) : null;
    } catch (e) {
      console.error('Failed to parse analytics data:', e);
      return null;
    }
  };

  // Add an event to analytics storage
  const trackEvent = (eventType, metadata = {}) => {
    try {
      const data = getData();
      if (!data) {
        initStorage();
        return trackEvent(eventType, metadata);
      }

      data.events.push({
        type: eventType,
        timestamp: new Date().toISOString(),
        metadata: metadata,
        user_agent: navigator.userAgent,
        referrer: document.referrer || 'direct'
      });

      localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
    } catch (e) {
      console.error('Failed to track event:', e);
    }
  };

  // Track page view on load
  const trackPageView = () => {
    trackEvent('page_view', {
      url: window.location.href,
      title: document.title
    });
  };

  // Clear all analytics data
  const clearData = () => {
    localStorage.removeItem(STORAGE_KEY);
    initStorage();
    console.log('Analytics data cleared');
  };

  // Export data as JSON (for manual download)
  const exportData = () => {
    const data = getData();
    return data ? JSON.stringify(data, null, 2) : '{}';
  };

  // Get analytics summary
  const getSummary = () => {
    const data = getData();
    if (!data || !data.events) return null;

    const events = data.events;
    const summary = {
      total_events: events.length,
      page_views: events.filter(e => e.type === 'page_view').length,
      password_generations: events.filter(e => e.type === 'password_generated').length,
      generations_by_mode: {},
      first_visit: data.created_at,
      last_activity: events.length > 0 ? events[events.length - 1].timestamp : null,
      unique_referrers: [...new Set(events.map(e => e.referrer))]
    };

    // Count generations by mode
    events
      .filter(e => e.type === 'password_generated')
      .forEach(e => {
        const mode = e.metadata.mode || 'unknown';
        summary.generations_by_mode[mode] = (summary.generations_by_mode[mode] || 0) + 1;
      });

    return summary;
  };

  // Initialize on load
  initStorage();

  return {
    trackEvent,
    trackPageView,
    getData,
    getSummary,
    exportData,
    clearData
  };
})();

// Track page view when script loads
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    PassForgeAnalytics.trackPageView();
  });
} else {
  PassForgeAnalytics.trackPageView();
}

// Make available globally for tracking from Go WASM
globalThis.PassForgeAnalytics = PassForgeAnalytics;

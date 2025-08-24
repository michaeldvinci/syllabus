import React from 'react';

const Features = () => {
  return (
    <div className="prose">
      <h1>Features</h1>

      <h2>Core Functionality</h2>
      <div className="feature-grid">
        <div className="feature-card">
          <h3>📄 YAML Configuration</h3>
          <p>Parse audiobook series from a simple YAML file with hot-reload support</p>
        </div>
        <div className="feature-card">
          <h3>🌐 Multi-Provider Scraping</h3>
          <p>Fetch data from both Audible and Amazon with intelligent rate limiting</p>
        </div>
        <div className="feature-card">
          <h3>📅 Release Date Tracking</h3>
          <p>Extract latest and next release dates automatically</p>
        </div>
        <div className="feature-card">
          <h3>💾 Database Persistence</h3>
          <p>SQLite database for reliable data storage with automatic migrations</p>
        </div>
        <div className="feature-card">
          <h3>⚙️ Background Processing</h3>
          <p>Multi-threaded background scraper with job queue</p>
        </div>
        <div className="feature-card">
          <h3>📡 Real-time Updates</h3>
          <p>Server-sent events for live UI updates during scraping</p>
        </div>
      </div>

      <h2>User Experience</h2>
      <div className="feature-grid">
        <div className="feature-card">
          <h3>🔐 Authentication System</h3>
          <p>Secure login with role-based access (Admin/User) and session management</p>
        </div>
        <div className="feature-card">
          <h3>📱 Responsive Web UI</h3>
          <p>Clean, mobile-friendly interface with dark mode support</p>
        </div>
        <div className="feature-card">
          <h3>🔄 Auto-refresh</h3>
          <p>Configurable automatic data refresh (2-10 hours) with UI controls</p>
        </div>
        <div className="feature-card">
          <h3>🎯 Manual Refresh</h3>
          <p>On-demand data refresh with progress tracking and live updates</p>
        </div>
        <div className="feature-card">
          <h3>📲 iCal Export</h3>
          <p>Subscribe to release date calendar in your favorite app</p>
        </div>
        <div className="feature-card">
          <h3>⚙️ Settings Panel</h3>
          <p>Manage refresh intervals, view modes, and user preferences</p>
        </div>
      </div>

      <h2>View Modes</h2>
      
      <h3>🖥️ Desktop Views</h3>
      <ul>
        <li><strong>Unified View:</strong> Side-by-side Audible and Amazon data in one table</li>
        <li><strong>Tabbed View:</strong> Separate tabs for Audible and Amazon with configurable default tab</li>
        <li><strong>Search & Filter:</strong> Real-time search with clear button and filter by upcoming releases</li>
        <li><strong>Sticky Headers:</strong> Table headers remain visible while scrolling</li>
      </ul>

      <h3>📱 Mobile Experience</h3>
      <ul>
        <li><strong>Full-screen Layout:</strong> Optimized mobile interface using full screen space</li>
        <li><strong>Mobile Tabs:</strong> Touch-friendly tabbed navigation for Audible/Amazon</li>
        <li><strong>Card Layout:</strong> Easy-to-read card-based design for series information</li>
        <li><strong>Responsive Grid:</strong> Adaptive statistics tiles that work on all screen sizes</li>
      </ul>

      <h2>Technical Features</h2>
      <div className="feature-grid">
        <div className="feature-card">
          <h3>📂 File Watching</h3>
          <p>Auto-reload when YAML configuration changes without restart</p>
        </div>
        <div className="feature-card">
          <h3>🛑 Graceful Shutdown</h3>
          <p>Clean application termination handling with job cleanup</p>
        </div>
        <div className="feature-card">
          <h3>🐳 Docker Support</h3>
          <p>Full containerization with Docker Compose and environment variable overrides</p>
        </div>
        <div className="feature-card">
          <h3>🔌 JSON API</h3>
          <p>Programmatic access to series data with authentication</p>
        </div>
        <div className="feature-card">
          <h3>⏱️ Rate Limiting</h3>
          <p>Intelligent scraping delays to avoid provider restrictions</p>
        </div>
        <div className="feature-card">
          <h3>🔄 Smart Caching</h3>
          <p>Configurable cache timeout to minimize unnecessary requests</p>
        </div>
      </div>

      <h2>Data Sources & Scraping</h2>
      
      <h3>Audible Scraping</h3>
      <ul>
        <li><strong>Series Count:</strong> Number of books from series page HTML</li>
        <li><strong>Latest Release:</strong> Most recent release date from series page</li>
        <li><strong>Next Release:</strong> Extracted from "Coming Soon" or pre-order sections</li>
      </ul>

      <h3>Amazon Scraping</h3>
      <ul>
        <li><strong>Series Count:</strong> Parsed from collection-size elements</li>
        <li><strong>Next Release:</strong> Date from success/bold span elements</li>
        <li><strong>Series Detection:</strong> Automatic ASIN extraction from URLs</li>
      </ul>

      <h3>Background Processing</h3>
      <ul>
        <li><strong>Multi-threaded:</strong> Configurable concurrent workers (default: 4)</li>
        <li><strong>Job Queue:</strong> Persistent SQLite-based job management</li>
        <li><strong>Rate Limiting:</strong> Intelligent delays to respect provider limits</li>
        <li><strong>Error Handling:</strong> Automatic retry logic with exponential backoff</li>
      </ul>

      <h2>Auto-Refresh System</h2>
      
      <h3>Configurable Intervals</h3>
      <ul>
        <li><strong>Options:</strong> 2, 4, 6, 8, 10 hours</li>
        <li><strong>Default:</strong> 6 hours</li>
        <li><strong>UI Control:</strong> Settings panel slider with real-time updates</li>
        <li><strong>Persistence:</strong> Interval survives restarts and is stored in database</li>
      </ul>

      <h3>Refresh Behavior</h3>
      <ul>
        <li><strong>Incremental:</strong> Only updates stale data based on cache timeout</li>
        <li><strong>Background:</strong> Non-blocking operation that doesn't affect UI</li>
        <li><strong>Progress:</strong> Real-time updates via Server-Sent Events</li>
        <li><strong>Manual Override:</strong> Refresh button forces immediate update</li>
      </ul>

      <h2>Search & Filtering</h2>
      <ul>
        <li><strong>Real-time Search:</strong> Instant filtering as you type</li>
        <li><strong>Clear Button:</strong> X button to quickly clear search terms</li>
        <li><strong>Auto-clear:</strong> Search automatically clears when switching tabs/views</li>
        <li><strong>Filter Options:</strong> Filter by series with upcoming releases</li>
        <li><strong>Multi-view Support:</strong> Search works across unified, tabbed, and mobile views</li>
      </ul>

      <h2>User Management</h2>
      <ul>
        <li><strong>Role-based Access:</strong> Admin and User roles with different permissions</li>
        <li><strong>Secure Authentication:</strong> bcrypt password hashing and session management</li>
        <li><strong>User Creation:</strong> Admin users can create additional accounts</li>
        <li><strong>Session Persistence:</strong> Login sessions survive browser restarts</li>
      </ul>

      <h2>Performance Expectations</h2>
      <ul>
        <li><strong>Initial Setup:</strong> First scrape takes 30-90 seconds depending on series count</li>
        <li><strong>Background Updates:</strong> Automatic refreshes are optimized and much faster</li>
        <li><strong>Rate Limiting:</strong> Intentional delays prevent overwhelming provider servers</li>
        <li><strong>Caching:</strong> Smart caching reduces unnecessary requests and improves performance</li>
      </ul>
    </div>
  );
};

export default Features;
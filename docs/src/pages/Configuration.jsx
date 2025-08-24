import React from 'react';

const Configuration = () => {
  return (
    <div className="prose">
      <h1>Configuration</h1>

      <h2>YAML Configuration Schema</h2>
      <p>
        The main configuration file is <code>config/books.yaml</code>. It supports optional 
        application settings and audiobook series configuration.
      </p>

      <h3>Complete Configuration Example</h3>
      <pre><code>{`# Application Settings (optional - defaults shown)
settings:
  auto_refresh_interval: 6  # Hours between automatic data refreshes (default: 6)
  default_workers: 4        # Number of concurrent scraper workers (default: 4)  
  server_port: 8080         # Port for the web server (default: 8080)
  cache_timeout: 6          # Cache timeout in hours (default: 6)
  log_level: "info"         # Logging level: debug, info, warn, error (default: info)
  main_view: "unified"      # Default view mode: unified, tabbed (default: unified)

# Audiobook/Ebook Series Configuration
audiobooks:
  - title: "1% Lifesteal"
    audible: "https://www.audible.com/series/1-Lifesteal-Audiobooks/B0F8QMLV9T"
    amazon: "https://www.amazon.com/dp/B0DGWCJ6JP"
  
  - title: "A Soldier's Life"
    audible: "https://www.audible.com/series/A-Soldiers-Life-Audiobooks/B0D34549LX"
    amazon: "https://www.amazon.com/dp/B0CW18NDBQ"`}</code></pre>

      <h3>Required Fields</h3>
      <p>Only <code>title</code>, <code>audible</code>, and <code>amazon</code> are required for each series.</p>
      
      <h3>Settings Section</h3>
      <p>
        All settings are optional and will use sensible defaults if not specified. The settings 
        section allows you to customize application behavior without modifying code.
      </p>

      <h2>Environment Variables</h2>
      <p>
        For containerized deployments, you can override any setting using environment variables. 
        Environment variables take precedence over YAML configuration.
      </p>

      <h3>Server Configuration</h3>
      <table>
        <thead>
          <tr>
            <th>Variable</th>
            <th>Description</th>
            <th>Default</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>SYLLABUS_SERVER_PORT</code></td>
            <td>Server port (1-65535)</td>
            <td>8080</td>
          </tr>
          <tr>
            <td><code>PORT</code></td>
            <td>Standard port env var (alternative)</td>
            <td>8080</td>
          </tr>
        </tbody>
      </table>

      <h3>Scraping Configuration</h3>
      <table>
        <thead>
          <tr>
            <th>Variable</th>
            <th>Description</th>
            <th>Default</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>SYLLABUS_AUTO_REFRESH_INTERVAL</code></td>
            <td>Hours between auto-refreshes (&gt;0)</td>
            <td>6</td>
          </tr>
          <tr>
            <td><code>SYLLABUS_DEFAULT_WORKERS</code></td>
            <td>Number of concurrent scraper workers (&gt;0)</td>
            <td>4</td>
          </tr>
          <tr>
            <td><code>SYLLABUS_CACHE_TIMEOUT</code></td>
            <td>Cache timeout in hours (&gt;0)</td>
            <td>6</td>
          </tr>
        </tbody>
      </table>

      <h3>UI Configuration</h3>
      <table>
        <thead>
          <tr>
            <th>Variable</th>
            <th>Description</th>
            <th>Options</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>SYLLABUS_MAIN_VIEW</code></td>
            <td>Default view mode</td>
            <td>"unified" or "tabbed"</td>
          </tr>
        </tbody>
      </table>

      <h3>Logging Configuration</h3>
      <table>
        <thead>
          <tr>
            <th>Variable</th>
            <th>Description</th>
            <th>Options</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>SYLLABUS_LOG_LEVEL</code></td>
            <td>Log level</td>
            <td>"debug", "info", "warn", "error"</td>
          </tr>
        </tbody>
      </table>

      <h2>Docker Compose Override</h2>
      <p>
        Use <code>docker-compose.override.yaml</code> to customize your deployment without 
        modifying the main docker-compose.yaml file:
      </p>

      <pre><code>{`# docker-compose.override.yaml
services:
  syllabus:
    environment:
      # Server Configuration
      SYLLABUS_SERVER_PORT: "8080"
      
      # Scraping Configuration  
      SYLLABUS_AUTO_REFRESH_INTERVAL: "4"  # Check every 4 hours
      SYLLABUS_DEFAULT_WORKERS: "2"        # Use 2 workers
      SYLLABUS_CACHE_TIMEOUT: "8"          # 8 hour cache
      
      # UI Configuration
      SYLLABUS_MAIN_VIEW: "tabbed"         # Default to tabbed view
      
      # Logging Configuration
      SYLLABUS_LOG_LEVEL: "debug"          # Enable debug logging`}</code></pre>

      <h2>Configuration Priority</h2>
      <p>Settings are applied in the following order (highest to lowest priority):</p>
      <ol>
        <li><strong>Runtime UI changes</strong> (highest priority) - persisted to database</li>
        <li><strong>Environment Variables</strong> - override YAML at startup</li>
        <li><strong>YAML Configuration</strong> - file-based defaults</li>
        <li><strong>Built-in Defaults</strong> (lowest priority)</li>
      </ol>

      <h2>Database Persistence</h2>
      <p>
        UI changes (like auto-refresh interval) are saved to the database and survive container 
        restarts. The UI will always show the current server state.
      </p>

      <h2>Configuration Watching</h2>
      <ul>
        <li><strong>Auto-reload:</strong> YAML file changes trigger incremental updates</li>
        <li><strong>Smart Updates:</strong> Only new series are scraped, existing data preserved</li>
        <li><strong>Hot Refresh:</strong> No application restart required</li>
      </ul>

      <h2>Example Configurations</h2>
      
      <h3>High-Frequency Monitoring</h3>
      <pre><code>{`settings:
  auto_refresh_interval: 2  # Check every 2 hours
  default_workers: 6        # More workers for faster scraping
  log_level: "debug"        # Detailed logging`}</code></pre>

      <h3>Low-Resource Setup</h3>
      <pre><code>{`settings:
  auto_refresh_interval: 12  # Check twice daily
  default_workers: 2         # Fewer workers
  cache_timeout: 24          # Longer cache`}</code></pre>

      <h3>Development Setup</h3>
      <pre><code>{`settings:
  auto_refresh_interval: 1  # Frequent updates for testing
  log_level: "debug"        # Debug logging
  main_view: "tabbed"       # Test tabbed view`}</code></pre>
    </div>
  );
};

export default Configuration;
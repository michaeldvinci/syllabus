import React from 'react';

const Troubleshooting = () => {
  return (
    <div className="prose">
      <h1>Troubleshooting</h1>

      <h2>Common Issues</h2>

      <h3>Application Won't Start</h3>
      
      <h4>Port 8080 Already in Use</h4>
      <pre><code>{`# Check what's using port 8080
lsof -i :8080

# Change the port in docker-compose.override.yaml
services:
  syllabus:
    ports:
      - "9000:8080"  # Use port 9000 instead`}</code></pre>

      <h4>Docker Compose Issues</h4>
      <pre><code>{`# Check Docker Compose logs
docker compose logs syllabus

# Rebuild containers
docker compose down
docker compose up -d --build

# Check Docker daemon status
docker info`}</code></pre>

      <h3>No Data Appearing</h3>

      <h4>Initial Scrape in Progress</h4>
      <p>The first scrape takes 30-90 seconds. Check the progress in the web UI or logs:</p>
      <pre><code>{`# Watch real-time logs
docker compose logs -f syllabus

# Check for scraping activity
grep -i "scraping\|worker\|job" logs`}</code></pre>

      <h4>Configuration Issues</h4>
      <ul>
        <li>Verify <code>config/books.yaml</code> has valid URLs</li>
        <li>Check that Audible and Amazon URLs are accessible</li>
        <li>Ensure YAML syntax is correct</li>
      </ul>

      <pre><code>{`# Validate YAML syntax
python -c "import yaml; yaml.safe_load(open('config/books.yaml'))"

# Or use a YAML linter
yamllint config/books.yaml`}</code></pre>

      <h3>Scraping Errors</h3>

      <h4>Rate Limiting</h4>
      <p>If you're seeing rate limiting errors, adjust the worker count:</p>
      <pre><code>{`# In docker-compose.override.yaml
services:
  syllabus:
    environment:
      SYLLABUS_DEFAULT_WORKERS: "2"  # Reduce from default 4
      SYLLABUS_CACHE_TIMEOUT: "12"   # Increase cache timeout`}</code></pre>

      <h4>Network Connectivity</h4>
      <pre><code>{`# Test connectivity to providers
curl -I https://www.audible.com
curl -I https://www.amazon.com

# Check Docker network
docker compose exec syllabus ping audible.com`}</code></pre>

      <h3>Authentication Issues</h3>

      <h4>Can't Login</h4>
      <ul>
        <li>Try default credentials: <code>admin</code> / <code>admin</code></li>
        <li>Check if user database is corrupted</li>
        <li>Reset user database if needed</li>
      </ul>

      <pre><code>{`# Reset user database (WARNING: loses all users)
docker compose down
rm data/users.json
docker compose up -d

# Default admin user will be recreated`}</code></pre>

      <h4>Session Expired</h4>
      <p>Clear browser cookies and log in again:</p>
      <ul>
        <li>Open browser developer tools (F12)</li>
        <li>Go to Application/Storage tab</li>
        <li>Clear cookies for your syllabus domain</li>
        <li>Refresh page and log in again</li>
      </ul>

      <h2>Performance Issues</h2>

      <h3>Slow Scraping</h3>
      <p>If scraping is taking too long:</p>
      <pre><code>{`# Increase workers (use carefully to avoid rate limits)
services:
  syllabus:
    environment:
      SYLLABUS_DEFAULT_WORKERS: "6"  # Increase from default 4

# Or reduce series count for testing
# Comment out some series in books.yaml`}</code></pre>

      <h3>High Memory Usage</h3>
      <pre><code>{`# Check container resource usage
docker stats

# Limit container memory if needed
services:
  syllabus:
    deploy:
      resources:
        limits:
          memory: 512M`}</code></pre>

      <h2>Database Issues</h2>

      <h3>Corrupted Database</h3>
      <pre><code>{`# Check database integrity
sqlite3 data/syllabus.db "PRAGMA integrity_check;"

# Reset database if corrupted (WARNING: loses all data)
docker compose down
rm data/syllabus.db
docker compose up -d`}</code></pre>

      <h3>Migration Failures</h3>
      <pre><code>{`# Check logs for migration errors
docker compose logs syllabus | grep -i migration

# If migrations fail, might need to reset database
# Backup data first if possible`}</code></pre>

      <h2>Log Analysis</h2>

      <h3>Log Locations</h3>
      <ul>
        <li><strong>Docker:</strong> <code>docker compose logs syllabus</code></li>
        <li><strong>Local:</strong> Console output</li>
        <li><strong>Detailed:</strong> Set <code>SYLLABUS_LOG_LEVEL=debug</code></li>
      </ul>

      <h3>Enable Debug Logging</h3>
      <pre><code>{`# In docker-compose.override.yaml
services:
  syllabus:
    environment:
      SYLLABUS_LOG_LEVEL: "debug"`}</code></pre>

      <h3>Common Log Messages</h3>
      <table>
        <thead>
          <tr>
            <th>Message</th>
            <th>Meaning</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>"Rate limited"</td>
            <td>Provider is throttling requests</td>
            <td>Reduce workers or increase delays</td>
          </tr>
          <tr>
            <td>"Connection timeout"</td>
            <td>Network connectivity issue</td>
            <td>Check internet connection</td>
          </tr>
          <tr>
            <td>"Parse error"</td>
            <td>HTML structure changed</td>
            <td>May need application update</td>
          </tr>
          <tr>
            <td>"Database locked"</td>
            <td>Concurrent access issue</td>
            <td>Usually resolves automatically</td>
          </tr>
        </tbody>
      </table>

      <h2>Browser Issues</h2>

      <h3>UI Not Loading</h3>
      <ul>
        <li>Hard refresh: Ctrl+F5 (Windows/Linux) or Cmd+Shift+R (Mac)</li>
        <li>Clear browser cache</li>
        <li>Try incognito/private mode</li>
        <li>Check browser console for JavaScript errors</li>
      </ul>

      <h3>Mobile Issues</h3>
      <ul>
        <li>Ensure you're using a modern mobile browser</li>
        <li>Try landscape orientation if layout seems broken</li>
        <li>Clear mobile browser cache</li>
      </ul>

      <h2>Configuration Validation</h2>

      <h3>YAML Syntax</h3>
      <pre><code>{`# Valid YAML structure
audiobooks:
  - title: "Series Name"
    audible: "https://www.audible.com/series/..."
    amazon: "https://www.amazon.com/dp/..."

# Common YAML mistakes:
# - Missing quotes around URLs with special characters
# - Incorrect indentation (must use spaces, not tabs)
# - Missing colons or dashes`}</code></pre>

      <h3>URL Validation</h3>
      <pre><code>{`# Test URLs manually
curl -I "https://www.audible.com/series/Your-Series-Name/B0EXAMPLE"
curl -I "https://www.amazon.com/dp/B0EXAMPLE"`}</code></pre>

      <h2>Getting Help</h2>

      <h3>Information to Collect</h3>
      <p>When reporting issues, please include:</p>
      <ul>
        <li>Docker Compose logs: <code>docker compose logs syllabus</code></li>
        <li>Configuration file (with sensitive URLs redacted)</li>
        <li>Browser console errors (if UI issue)</li>
        <li>System information (OS, Docker version)</li>
        <li>Steps to reproduce the issue</li>
      </ul>

      <h3>Support Channels</h3>
      <ul>
        <li><strong>GitHub Issues:</strong> <a href="https://github.com/michaeldvinci/syllabus/issues" target="_blank">Report bugs and feature requests</a></li>
        <li><strong>Documentation:</strong> Check other sections of this documentation</li>
        <li><strong>Logs:</strong> Enable debug logging for detailed troubleshooting</li>
      </ul>

      <h2>Emergency Recovery</h2>

      <h3>Complete Reset</h3>
      <p>If all else fails, complete reset (WARNING: loses all data):</p>
      <pre><code>{`# Stop everything
docker compose down

# Remove all data
rm -rf data/

# Reset configuration to example
cp config/example-settings.yaml config/books.yaml

# Start fresh
docker compose up -d`}</code></pre>

      <h3>Backup Before Reset</h3>
      <pre><code>{`# Backup important files
cp data/syllabus.db backup-syllabus.db
cp data/users.json backup-users.json
cp config/books.yaml backup-books.yaml`}</code></pre>
    </div>
  );
};

export default Troubleshooting;
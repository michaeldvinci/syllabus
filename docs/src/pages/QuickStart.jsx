import React from 'react';
import { Link } from 'react-router-dom';

const QuickStart = () => {
  return (
    <div className="prose">
      <h1>Quick Start</h1>
      
      <div className="warning-box">
        <h4>Default Credentials</h4>
        <p>
          <strong>Username:</strong> <code>admin</code><br />
          <strong>Password:</strong> <code>admin</code><br />
          <strong>⚠️ Change the default password immediately after first login</strong>
        </p>
      </div>

      <h2>Docker Compose (Recommended)</h2>
      <p>The fastest way to get started:</p>

      <pre><code>{`# Clone the repository
git clone https://github.com/michaeldvinci/syllabus.git
cd syllabus

# Start the application
docker compose up -d

# View logs (optional)
docker compose logs -f`}</code></pre>

      <h2>Access the Application</h2>
      <ul>
        <li><strong>Web UI:</strong> <a href="http://localhost:8080" target="_blank">http://localhost:8080</a></li>
        <li><strong>Login:</strong> Use <code>admin</code> / <code>admin</code> (change immediately!)</li>
        <li><strong>JSON API:</strong> <a href="http://localhost:8080/api/series" target="_blank">http://localhost:8080/api/series</a></li>
        <li><strong>Calendar:</strong> <a href="http://localhost:8080/calendar.ics" target="_blank">http://localhost:8080/calendar.ics</a></li>
      </ul>

      <h2>First Steps</h2>
      <ol>
        <li>
          <strong>Change Default Password</strong>
          <ul>
            <li>Log in with <code>admin</code> / <code>admin</code></li>
            <li>Go to Settings → User Management</li>
            <li>Create a new admin user or change the default password</li>
          </ul>
        </li>
        
        <li>
          <strong>Configure Your Series</strong>
          <ul>
            <li>Edit the <code>config/books.yaml</code> file</li>
            <li>Add your audiobook series with Audible and Amazon URLs</li>
            <li>The app will automatically detect changes and start scraping</li>
          </ul>
        </li>

        <li>
          <strong>Wait for Initial Scrape</strong>
          <ul>
            <li>First scrape takes 30-90 seconds depending on series count</li>
            <li>Watch the progress in the web UI</li>
            <li>Background updates are much faster</li>
          </ul>
        </li>
      </ol>

      <h2>Sample Configuration</h2>
      <p>Here's a basic <code>config/books.yaml</code> to get you started:</p>

      <pre><code>{`# Application Settings (optional)
settings:
  auto_refresh_interval: 6  # Hours between automatic refreshes
  default_workers: 4        # Number of concurrent scraper workers
  server_port: 8080         # Web server port
  main_view: "unified"      # Default view: unified or tabbed

# Audiobook Series
audiobooks:
  - title: "Example Series 1"
    audible: "https://www.audible.com/series/Example-Series-Audiobooks/B0EXAMPLE1"
    amazon: "https://www.amazon.com/dp/B0EXAMPLE1"
  
  - title: "Example Series 2"
    audible: "https://www.audible.com/series/Example-Series-2-Audiobooks/B0EXAMPLE2"
    amazon: "https://www.amazon.com/dp/B0EXAMPLE2"`}</code></pre>

      <h2>Environment Variables (Optional)</h2>
      <p>You can customize settings using environment variables in <code>docker-compose.override.yaml</code>:</p>

      <pre><code>{`# Copy the example override file
cp docker-compose.override.example.yaml docker-compose.override.yaml

# Edit your settings
nano docker-compose.override.yaml

# Restart with new settings
docker compose up -d`}</code></pre>

      <h2>Troubleshooting</h2>
      <ul>
        <li><strong>Container won't start:</strong> Check port 8080 isn't in use</li>
        <li><strong>No data appearing:</strong> Wait for initial scrape to complete (30-90 seconds)</li>
        <li><strong>Scraping errors:</strong> Check the logs with <code>docker compose logs syllabus</code></li>
      </ul>

      <h2>Next Steps</h2>
      <div className="bg-gray-50 rounded-lg p-6 border border-gray-200">
        <h3 className="mt-0">Explore More Features</h3>
        <ul className="mb-0">
          <li><Link to="/configuration">Configuration Guide</Link> - Detailed setup options</li>
          <li><Link to="/features">Features Overview</Link> - Complete feature list</li>
          <li><Link to="/api">API Reference</Link> - JSON API documentation</li>
          <li><Link to="/troubleshooting">Troubleshooting</Link> - Common issues and solutions</li>
        </ul>
      </div>
    </div>
  );
};

export default QuickStart;
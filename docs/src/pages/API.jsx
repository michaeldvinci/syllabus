import React from 'react';

const API = () => {
  return (
    <div className="prose">
      <h1>API Reference</h1>

      <div className="warning-box">
        <h4>🔐 Authentication Required</h4>
        <p>
          All API endpoints require authentication via session cookie. 
          Login at <code>/login</code> to establish a session.
        </p>
      </div>

      <h2>GET /api/series</h2>
      <p>Returns an array of series objects with complete scraping data.</p>

      <h3>Response Format</h3>
      <pre><code>{`[
  {
    "Title": "Example Series",
    "AudibleCount": 5,
    "AudibleLatestTitle": "Book 5: Latest Release",
    "AudibleLatestDate": "2024-01-15T00:00:00Z",
    "AudibleNextTitle": "Book 6: Coming Soon", 
    "AudibleNextDate": "2024-03-20T00:00:00Z",
    "AmazonCount": 5,
    "AmazonLatestTitle": "Book 5: Latest Release",
    "AmazonLatestDate": "2024-01-15T00:00:00Z",
    "AmazonNextTitle": "Book 6: Coming Soon",
    "AmazonNextDate": "2024-03-20T00:00:00Z",
    "AudibleID": "B0EXAMPLE123",
    "AmazonASIN": "B08EXAMPLE456",
    "Err": null
  }
]`}</code></pre>

      <h3>Field Descriptions</h3>
      <table>
        <thead>
          <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>Title</code></td>
            <td>string</td>
            <td>Series name from configuration</td>
          </tr>
          <tr>
            <td><code>AudibleCount</code></td>
            <td>int</td>
            <td>Number of books in Audible series</td>
          </tr>
          <tr>
            <td><code>AudibleLatestTitle</code></td>
            <td>string</td>
            <td>Title of most recent Audible release</td>
          </tr>
          <tr>
            <td><code>AudibleLatestDate</code></td>
            <td>string</td>
            <td>Release date of latest Audible book (ISO 8601)</td>
          </tr>
          <tr>
            <td><code>AudibleNextTitle</code></td>
            <td>string</td>
            <td>Title of upcoming Audible release</td>
          </tr>
          <tr>
            <td><code>AudibleNextDate</code></td>
            <td>string</td>
            <td>Expected release date for next Audible book (ISO 8601)</td>
          </tr>
          <tr>
            <td><code>AmazonCount</code></td>
            <td>int</td>
            <td>Number of books in Amazon series</td>
          </tr>
          <tr>
            <td><code>AmazonLatestTitle</code></td>
            <td>string</td>
            <td>Title of most recent Amazon release</td>
          </tr>
          <tr>
            <td><code>AmazonLatestDate</code></td>
            <td>string</td>
            <td>Release date of latest Amazon book (ISO 8601)</td>
          </tr>
          <tr>
            <td><code>AmazonNextTitle</code></td>
            <td>string</td>
            <td>Title of upcoming Amazon release</td>
          </tr>
          <tr>
            <td><code>AmazonNextDate</code></td>
            <td>string</td>
            <td>Expected release date for next Amazon book (ISO 8601)</td>
          </tr>
          <tr>
            <td><code>AudibleID</code></td>
            <td>string</td>
            <td>Audible series identifier</td>
          </tr>
          <tr>
            <td><code>AmazonASIN</code></td>
            <td>string</td>
            <td>Amazon ASIN identifier</td>
          </tr>
          <tr>
            <td><code>Err</code></td>
            <td>string|null</td>
            <td>Error message if scraping failed, null if successful</td>
          </tr>
        </tbody>
      </table>

      <h3>Example Request</h3>
      <pre><code>{`curl -X GET "http://localhost:8080/api/series" \\
     -H "Cookie: session_token=your_session_token"`}</code></pre>

      <h2>POST /refresh</h2>
      <p>Triggers a manual refresh of all series data.</p>

      <h3>Request</h3>
      <pre><code>{`curl -X POST "http://localhost:8080/refresh" \\
     -H "Cookie: session_token=your_session_token"`}</code></pre>

      <h3>Response</h3>
      <p>Returns HTTP 200 on success. Progress can be monitored via Server-Sent Events.</p>

      <h2>GET /calendar.ics</h2>
      <p>Returns iCal calendar file with all upcoming release dates.</p>

      <h3>Request</h3>
      <pre><code>{`curl -X GET "http://localhost:8080/calendar.ics" \\
     -H "Cookie: session_token=your_session_token"`}</code></pre>

      <h3>Response</h3>
      <p>Returns an iCal formatted calendar file with <code>Content-Type: text/calendar</code>.</p>

      <h2>POST /api/auto-refresh</h2>
      <p>Updates automatic refresh interval (Admin only).</p>

      <h3>Request Body</h3>
      <pre><code>{`{
  "interval": 6
}`}</code></pre>

      <h3>Example Request</h3>
      <pre><code>{`curl -X POST "http://localhost:8080/api/auto-refresh" \\
     -H "Content-Type: application/json" \\
     -H "Cookie: session_token=your_session_token" \\
     -d '{"interval": 6}'`}</code></pre>

      <h3>Parameters</h3>
      <table>
        <thead>
          <tr>
            <th>Parameter</th>
            <th>Type</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>interval</code></td>
            <td>int</td>
            <td>Hours between automatic refreshes (2-10)</td>
          </tr>
        </tbody>
      </table>

      <h2>Date Formatting</h2>
      <p>All dates are returned in ISO 8601 format:</p>
      <ul>
        <li><strong>Format:</strong> <code>YYYY-MM-DDTHH:mm:ssZ</code></li>
        <li><strong>Timezone:</strong> UTC</li>
        <li><strong>Example:</strong> <code>2024-03-20T00:00:00Z</code></li>
      </ul>

      <h2>Error Responses</h2>

      <h3>401 Unauthorized</h3>
      <pre><code>{`{
  "error": "authentication required"
}`}</code></pre>

      <h3>403 Forbidden</h3>
      <pre><code>{`{
  "error": "admin access required"
}`}</code></pre>

      <h3>500 Internal Server Error</h3>
      <pre><code>{`{
  "error": "internal server error"
}`}</code></pre>

      <h2>Authentication</h2>
      
      <h3>Login</h3>
      <p>Authenticate via the web interface at <code>/login</code> to establish a session cookie.</p>

      <h3>Session Management</h3>
      <ul>
        <li><strong>Cookie Name:</strong> <code>session_token</code></li>
        <li><strong>HttpOnly:</strong> Yes (not accessible via JavaScript)</li>
        <li><strong>Secure:</strong> Yes (when using HTTPS)</li>
        <li><strong>SameSite:</strong> Strict</li>
      </ul>

      <h2>Rate Limiting</h2>
      <p>The API respects the same rate limiting as the web scraper:</p>
      <ul>
        <li>Requests are throttled to avoid overwhelming external providers</li>
        <li>Manual refresh operations are queued and processed sequentially</li>
        <li>Concurrent scraping is limited by the <code>default_workers</code> setting</li>
      </ul>

      <h2>Server-Sent Events</h2>
      <p>Real-time updates during scraping operations are available via SSE:</p>

      <h3>Endpoint</h3>
      <pre><code>GET /events</code></pre>

      <h3>Event Types</h3>
      <ul>
        <li><strong>progress:</strong> Scraping progress updates</li>
        <li><strong>complete:</strong> Scraping operation completed</li>
        <li><strong>error:</strong> Scraping error occurred</li>
      </ul>

      <h2>Example Integration</h2>
      
      <h3>Python Script</h3>
      <pre><code>{`import requests
import json

# Login to get session
session = requests.Session()
login_response = session.post('http://localhost:8080/login', data={
    'username': 'admin',
    'password': 'your_password'
})

if login_response.status_code == 200:
    # Get series data
    series_response = session.get('http://localhost:8080/api/series')
    series_data = series_response.json()
    
    for series in series_data:
        print(f"Series: {series['Title']}")
        print(f"Next Release: {series['AudibleNextDate'] or series['AmazonNextDate']}")
        print()
else:
    print("Login failed")`}</code></pre>

      <h3>JavaScript (Browser)</h3>
      <pre><code>{`// Assuming you're already logged in via the web interface
fetch('/api/series')
  .then(response => response.json())
  .then(data => {
    data.forEach(series => {
      console.log(\`\${series.Title}: Next release \${series.AudibleNextDate}\`);
    });
  })
  .catch(error => console.error('Error:', error));`}</code></pre>
    </div>
  );
};

export default API;
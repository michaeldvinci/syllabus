import React from 'react';

const Terms = () => {
  return (
    <div className="prose">
      <h1>Terms of Service & Compliance</h1>

      <div className="warning-box">
        <h4>⚠️ Important Notice</h4>
        <p>
          This application collects publicly available data from Audible and Amazon through 
          automated web scraping for personal use only. Users are responsible for ensuring 
          their usage complies with applicable terms of service.
        </p>
      </div>

      <h2>Data Collection & Usage</h2>

      <h3>What Data is Collected</h3>
      <ul>
        <li><strong>Public Information Only:</strong> Release dates, book titles, series counts</li>
        <li><strong>No Account Data:</strong> Does not access or require user accounts on either platform</li>
        <li><strong>No Personal Information:</strong> Does not collect personal data about users or purchases</li>
        <li><strong>Metadata Only:</strong> Focuses on publicly visible series and release information</li>
      </ul>

      <h3>How Data is Used</h3>
      <ul>
        <li><strong>Personal Tracking:</strong> Intended for individual tracking of series you're interested in</li>
        <li><strong>Local Storage:</strong> All data is stored locally in your own database</li>
        <li><strong>No Redistribution:</strong> Data is not shared, sold, or redistributed to third parties</li>
        <li><strong>Automation Only:</strong> Used solely to automate manual tracking tasks</li>
      </ul>

      <h2>Compliance & Best Practices</h2>

      <h3>Rate Limiting</h3>
      <ul>
        <li><strong>Respectful Scraping:</strong> Implements intelligent delays between requests</li>
        <li><strong>Limited Concurrency:</strong> Configurable worker limits (default: 4 concurrent requests)</li>
        <li><strong>Caching:</strong> Reduces unnecessary requests through intelligent caching</li>
        <li><strong>Reasonable Frequency:</strong> Default 6-hour intervals between automatic updates</li>
      </ul>

      <h3>Technical Safeguards</h3>
      <ul>
        <li><strong>User-Agent Headers:</strong> Identifies requests appropriately</li>
        <li><strong>Request Spacing:</strong> Prevents overwhelming provider servers</li>
        <li><strong>Error Handling:</strong> Graceful handling of rate limits and blocks</li>
        <li><strong>Retry Logic:</strong> Exponential backoff for failed requests</li>
      </ul>

      <h2>Terms of Service Review</h2>

      <h3>Audible Terms of Use</h3>
      <p>
        Users should review <a href="https://www.audible.com/conditions-of-use" target="_blank" rel="noopener noreferrer">
        Audible's Conditions of Use</a> to ensure their personal use case complies with Audible's policies.
      </p>

      <h4>Key Considerations:</h4>
      <ul>
        <li>Automated access restrictions and limitations</li>
        <li>Personal vs. commercial use guidelines</li>
        <li>Data collection and usage policies</li>
        <li>Geographic restrictions and availability</li>
      </ul>

      <h3>Amazon Conditions of Use</h3>
      <p>
        Users should review <a href="https://www.amazon.com/gp/help/customer/display.html?nodeId=508088" target="_blank" rel="noopener noreferrer">
        Amazon's Conditions of Use</a> to ensure their personal use case complies with Amazon's policies.
      </p>

      <h4>Key Considerations:</h4>
      <ul>
        <li>Acceptable use policies for automated tools</li>
        <li>Restrictions on automated data collection</li>
        <li>Personal use vs. commercial use distinctions</li>
        <li>Intellectual property and content usage rights</li>
      </ul>

      <h2>Responsible Usage Guidelines</h2>

      <h3>Recommended Practices</h3>
      <ul>
        <li><strong>Personal Use Only:</strong> Use solely for tracking your own reading interests</li>
        <li><strong>Reasonable Intervals:</strong> Use default refresh intervals (6 hours) or longer</li>
        <li><strong>Limited Series Count:</strong> Track a reasonable number of series (suggest &lt;50)</li>
        <li><strong>Monitor Logs:</strong> Watch for rate limiting and adjust settings accordingly</li>
      </ul>

      <h3>Configuration Recommendations</h3>
      <pre><code>{`# Conservative configuration for compliance
settings:
  auto_refresh_interval: 8    # Check every 8 hours
  default_workers: 2          # Use only 2 concurrent workers
  cache_timeout: 12           # Cache for 12 hours to reduce requests`}</code></pre>

      <h3>What to Avoid</h3>
      <ul>
        <li><strong>Commercial Use:</strong> Do not use for commercial purposes or profit</li>
        <li><strong>Data Redistribution:</strong> Do not share, sell, or republish scraped data</li>
        <li><strong>Excessive Requests:</strong> Avoid high-frequency scraping or large worker counts</li>
        <li><strong>Circumventing Blocks:</strong> Do not attempt to bypass rate limits or IP blocks</li>
      </ul>

      <h2>Legal Disclaimer</h2>

      <blockquote>
        <p>
          <strong>No Legal Advice:</strong> This documentation does not constitute legal advice. 
          Users are responsible for understanding and complying with applicable laws, terms of 
          service, and usage policies in their jurisdiction.
        </p>
      </blockquote>

      <blockquote>
        <p>
          <strong>Terms May Change:</strong> Provider terms of service may change without notice. 
          Users should periodically review the terms of service for Audible and Amazon to ensure 
          continued compliance.
        </p>
      </blockquote>

      <blockquote>
        <p>
          <strong>Use at Own Risk:</strong> This software is provided "as is" without warranty. 
          Users assume all responsibility for their usage and any consequences thereof.
        </p>
      </blockquote>

      <h2>Performance Expectations</h2>

      <h3>Timing Guidelines</h3>
      <ul>
        <li><strong>Initial Setup:</strong> First scrape may take 30-90 seconds, depending on series count</li>
        <li><strong>Background Updates:</strong> Automatic refreshes are optimized and much faster</li>
        <li><strong>Rate Limiting:</strong> Intentional delays prevent overwhelming provider servers</li>
      </ul>

      <h3>Resource Usage</h3>
      <ul>
        <li><strong>Network Impact:</strong> Minimal bandwidth usage for text-based scraping</li>
        <li><strong>Server Load:</strong> Designed to be respectful of provider infrastructure</li>
        <li><strong>Local Resources:</strong> Lightweight local storage and processing</li>
      </ul>

      <h2>Monitoring & Adjustment</h2>

      <h3>Signs to Reduce Usage</h3>
      <ul>
        <li>Frequent rate limiting errors in logs</li>
        <li>Requests being blocked or denied</li>
        <li>Significantly slower response times</li>
        <li>Error messages indicating policy violations</li>
      </ul>

      <h3>Adjustment Strategies</h3>
      <pre><code>{`# If experiencing rate limiting, try:
settings:
  auto_refresh_interval: 12   # Reduce frequency
  default_workers: 1          # Single worker only
  cache_timeout: 24           # Longer cache retention`}</code></pre>

      <h2>Contact & Reporting</h2>

      <h3>Issue Reporting</h3>
      <p>
        If you encounter issues that may be related to terms of service compliance, 
        please report them via <a href="https://github.com/michaeldvinci/syllabus/issues" target="_blank" rel="noopener noreferrer">
        GitHub Issues</a> so they can be addressed in future versions.
      </p>

      <h3>Policy Changes</h3>
      <p>
        The maintainers of this project will make reasonable efforts to update the software 
        to maintain compliance with provider policies, but users are ultimately responsible 
        for their own usage.
      </p>

      <h2>Summary</h2>

      <div className="bg-blue-50 rounded-lg p-6 border border-blue-200">
        <h3 className="text-lg font-semibold text-blue-900 mb-4 mt-0">
          Key Takeaways
        </h3>
        <ul className="text-blue-800 mb-0">
          <li><strong>Personal Use Only:</strong> Keep it personal, not commercial</li>
          <li><strong>Be Respectful:</strong> Use reasonable intervals and worker counts</li>
          <li><strong>Monitor Usage:</strong> Watch logs and adjust if needed</li>
          <li><strong>Review Terms:</strong> Periodically check provider terms of service</li>
          <li><strong>Stay Informed:</strong> Keep software updated for compliance improvements</li>
        </ul>
      </div>

      <p>
        This tool is designed for personal audiobook enthusiasts who want to automate the 
        tedious task of manually checking for new releases. When used responsibly and in 
        compliance with provider terms, it provides a valuable service while respecting 
        the infrastructure and policies of Audible and Amazon.
      </p>
    </div>
  );
};

export default Terms;
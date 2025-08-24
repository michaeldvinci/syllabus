import React from 'react';

const Installation = () => {
  return (
    <div className="prose">
      <h1>Installation</h1>

      <h2>Requirements</h2>
      <ul>
        <li>Go 1.21+ (for local development)</li>
        <li>Docker & Docker Compose (recommended)</li>
        <li>Internet access for scraping</li>
        <li>YAML configuration file</li>
      </ul>

      <h2>Docker Compose (Recommended)</h2>
      <p>The easiest way to deploy syllabus:</p>

      <pre><code>{`git clone https://github.com/michaeldvinci/syllabus.git
cd syllabus

# Basic deployment
docker compose up -d

# View logs
docker compose logs -f`}</code></pre>

      <h3>Customization with Environment Variables</h3>
      <pre><code>{`# Copy example override file
cp docker-compose.override.example.yaml docker-compose.override.yaml

# Edit settings
nano docker-compose.override.yaml

# Deploy with custom settings
docker compose up -d`}</code></pre>

      <p>
        The override file allows you to customize settings without modifying the main 
        docker-compose.yaml file.
      </p>

      <h2>Local Development</h2>
      <p>For development or running directly on your machine:</p>

      <pre><code>{`git clone https://github.com/michaeldvinci/syllabus.git
cd syllabus

# Run directly
go run cmd/syllabus/main.go config/books.yaml

# Or build first
go build -o syllabus cmd/syllabus/main.go
./syllabus config/books.yaml`}</code></pre>

      <h2>Data Persistence</h2>
      
      <h3>SQLite Database</h3>
      <ul>
        <li><strong>Location:</strong> <code>./data/syllabus.db</code></li>
        <li><strong>Schema:</strong> Series, books, and job queue tables</li>
        <li><strong>Persistence:</strong> Survives application restarts</li>
        <li><strong>Migration:</strong> Automatic schema updates on startup</li>
      </ul>

      <h3>User Management</h3>
      <ul>
        <li><strong>Storage:</strong> <code>./data/users.json</code></li>
        <li><strong>Encryption:</strong> bcrypt password hashing</li>
        <li><strong>Roles:</strong> Admin and User access levels</li>
        <li><strong>Default:</strong> Admin user created on first run</li>
      </ul>

      <h2>Directory Structure</h2>
      <pre><code>{`syllabus/
├── config/
│   ├── books.yaml              # Main configuration
│   └── example-settings.yaml   # Example settings
├── data/
│   ├── syllabus.db            # SQLite database
│   └── users.json             # User accounts
├── docker-compose.yaml         # Main compose file
└── docker-compose.override.example.yaml  # Example customizations`}</code></pre>

      <h2>Port Configuration</h2>
      <p>By default, syllabus runs on port 8080. You can change this:</p>

      <h3>Docker Compose</h3>
      <pre><code>{`# In docker-compose.override.yaml
services:
  syllabus:
    ports:
      - "9000:8080"  # Map host port 9000 to container port 8080
    environment:
      PORT: "8080"   # Internal container port`}</code></pre>

      <h3>Local Development</h3>
      <pre><code>{`# Set environment variable
export SYLLABUS_SERVER_PORT=9000
go run cmd/syllabus/main.go config/books.yaml

# Or in your YAML config
settings:
  server_port: 9000`}</code></pre>

      <h2>Initial Setup</h2>
      <ol>
        <li>
          <strong>Start the application</strong>
          <pre><code>docker compose up -d</code></pre>
        </li>
        
        <li>
          <strong>Access the web interface</strong>
          <p>Navigate to <a href="http://localhost:8080">http://localhost:8080</a></p>
        </li>
        
        <li>
          <strong>Login with default credentials</strong>
          <ul>
            <li>Username: <code>admin</code></li>
            <li>Password: <code>admin</code></li>
          </ul>
        </li>
        
        <li>
          <strong>Change the default password</strong>
          <p>Go to Settings → User Management and create a new admin user</p>
        </li>
        
        <li>
          <strong>Configure your series</strong>
          <p>Edit <code>config/books.yaml</code> with your audiobook series</p>
        </li>
        
        <li>
          <strong>Wait for initial scrape</strong>
          <p>The first scrape will take 30-90 seconds depending on series count</p>
        </li>
      </ol>

      <h2>Updating</h2>
      <p>To update to the latest version:</p>

      <pre><code>{`# Pull latest changes
git pull origin main

# Rebuild and restart
docker compose down
docker compose up -d --build

# Or if using pre-built images
docker compose pull
docker compose up -d`}</code></pre>

      <h2>Uninstalling</h2>
      <p>To completely remove syllabus:</p>

      <pre><code>{`# Stop and remove containers
docker compose down

# Remove data (optional - this will delete your database)
rm -rf data/

# Remove configuration (optional)
rm config/books.yaml

# Remove Docker images (optional)
docker image rm syllabus:latest`}</code></pre>

      <div className="warning-box">
        <h4>⚠️ Data Loss Warning</h4>
        <p>
          Removing the <code>data/</code> directory will permanently delete your database, 
          user accounts, and all scraped data. Make sure to backup important data before 
          uninstalling.
        </p>
      </div>
    </div>
  );
};

export default Installation;
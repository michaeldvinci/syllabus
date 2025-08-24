import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowRightIcon, BookOpenIcon, ClockIcon, ServerIcon } from '@heroicons/react/24/outline';

const Home = () => {
  return (
    <div className="prose">
      <div className="text-center mb-12">
        <BookOpenIcon className="h-16 w-16 text-blue-600 mx-auto mb-4" />
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          syllabus
        </h1>
        <p className="text-xl text-gray-600 max-w-3xl mx-auto">
          A Go web application that tracks audiobook series release dates by scraping Audible and Amazon. 
          Features user authentication, database persistence, background scraping, and automatic refresh 
          capabilities with a clean web UI and JSON API.
        </p>
      </div>

      <div className="warning-box">
        <h4>⚠️ Important Notice</h4>
        <p>
          This application collects publicly available data from Audible and Amazon through automated 
          web scraping for personal use only. Please review the Terms & Compliance section before using.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div className="text-center p-6 bg-white rounded-lg border border-gray-200 hover:shadow-md transition-shadow">
          <ClockIcon className="h-12 w-12 text-green-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Quick Setup</h3>
          <p className="text-gray-600 text-sm mb-4">Initial scrape: 30-90 seconds, then automatic updates</p>
        </div>
        <div className="text-center p-6 bg-white rounded-lg border border-gray-200 hover:shadow-md transition-shadow">
          <ServerIcon className="h-12 w-12 text-blue-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Docker Ready</h3>
          <p className="text-gray-600 text-sm mb-4">Perfect for homelab deployment with Docker Compose</p>
        </div>
        <div className="text-center p-6 bg-white rounded-lg border border-gray-200 hover:shadow-md transition-shadow">
          <BookOpenIcon className="h-12 w-12 text-purple-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Multi-Platform</h3>
          <p className="text-gray-600 text-sm mb-4">Tracks both Audible and Amazon release data</p>
        </div>
      </div>

      <h2>Why syllabus?</h2>
      <p>
        Originally created to replace manual maintenance of an Obsidian "database" for tracking audiobook 
        release dates. I call it <code>syllabus</code> because it's a list of things to read.
      </p>

      <h2>Key Features</h2>
      <div className="feature-grid">
        <div className="feature-card">
          <h3>📊 Multi-View Interface</h3>
          <p>Switch between unified and tabbed views on both desktop and mobile</p>
        </div>
        <div className="feature-card">
          <h3>🔄 Background Scraping</h3>
          <p>Automatic data refresh with intelligent rate limiting</p>
        </div>
        <div className="feature-card">
          <h3>📱 Mobile Responsive</h3>
          <p>Full-featured mobile experience with tabbed navigation</p>
        </div>
        <div className="feature-card">
          <h3>🔐 User Authentication</h3>
          <p>Secure login system with admin and user roles</p>
        </div>
        <div className="feature-card">
          <h3>📅 iCal Export</h3>
          <p>Subscribe to release dates in your favorite calendar app</p>
        </div>
        <div className="feature-card">
          <h3>🎛️ Live Updates</h3>
          <p>Real-time progress updates via Server-Sent Events</p>
        </div>
      </div>

      <h2>Screenshots</h2>
      
      <h3>Unified View</h3>
      <div className="mb-6">
        <img 
          src="https://github.com/michaeldvinci/syllabus/raw/main/res/syllabus-unified.png" 
          alt="Syllabus Unified View" 
          className="w-full rounded-lg border border-gray-200 shadow-sm"
        />
      </div>

      <h3>Tabbed View</h3>
      <div className="mb-6">
        <img 
          src="https://github.com/michaeldvinci/syllabus/raw/main/res/syllabus-tabbed.png" 
          alt="Syllabus Tabbed View" 
          className="w-full rounded-lg border border-gray-200 shadow-sm"
        />
      </div>

      <h2>Ready to get started?</h2>
      <div className="bg-blue-50 rounded-lg p-6 border border-blue-200">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-blue-900 mb-2 mt-0">
              Quick Start Guide
            </h3>
            <p className="text-blue-700 mb-0">
              Get up and running with Docker in under 5 minutes
            </p>
          </div>
          <Link 
            to="/quick-start"
            className="inline-flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 no-underline"
          >
            Get Started
            <ArrowRightIcon className="ml-2 h-4 w-4" />
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Home;
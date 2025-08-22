package handlers

const IndexHTML = `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Syllabus</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
:root {
  --bg: #ffffff;
  --text: #111827;
  --muted: #6b7280;
  --line: #e5e7eb;
  --head-bg: #f9fafb;
  --head-shadow: 0 1px 0 rgba(0,0,0,.04);
  --row-hover: #f3f4f6;
  --aud: #0ea5e9; /* cyan-ish */
  --amz: #f59e0b; /* amber-ish */
}

[data-theme="dark"] {
  --bg: #2b2b2b;
  --text: #a9b7c6;
  --muted: #808080;
  --line: #3c3f41;
  --head-bg: #3c3f41;
  --head-shadow: 0 1px 0 rgba(255,255,255,.03);
  --row-hover: #3c3f41;
}

/* Dark mode overrides for table elements */
[data-theme="dark"] table {
  background: #313335;
  border-color: var(--line);
}

[data-theme="dark"] tbody tr {
  background: #313335;
}

[data-theme="dark"] tbody tr:nth-child(even) {
  background: #3c3f41;
}

[data-theme="dark"] tbody tr:hover {
  background: #4c5052;
}

[data-theme="dark"] .settings-btn {
  background: #3c3f41;
  color: var(--text);
  border-color: var(--line);
}

[data-theme="dark"] .settings-btn:hover {
  background: #4c5052;
}

[data-theme="dark"] .settings-panel {
  background: #3c3f41;
  border-color: var(--line);
}
body { 
  font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial; 
  margin: 0; 
  color: var(--text); 
  background: var(--bg); 
  height: 100vh;
  display: flex;
  flex-direction: column;
}
/* Content area with scrolling */
.content-area {
  flex: 1;
  overflow: auto;
  padding: 1rem 1.5rem;
}

table { border-collapse: separate; border-spacing: 0; width: 100%; background: #fff; border: 1px solid var(--line); border-radius: .5rem; overflow: hidden; }
thead th { position: sticky; top: 0; background: var(--head-bg); z-index: 3; }
thead tr:nth-child(2) th { top: 2.5rem; z-index: 2; }
thead th { border-bottom: 1px solid var(--line); padding: .6rem .75rem; text-align: left; font-weight: 600; }
thead tr:first-child th { box-shadow: var(--head-shadow); }
th, td { border-bottom: 1px solid var(--line); padding: .5rem .75rem; text-align: left; vertical-align: middle; }
tbody tr:nth-child(even) { background: #fcfcfd; }
tbody tr:hover { background: var(--row-hover); }
small { color: var(--muted); }
/* Sticky first column for easier scanning */
th:first-child, td:first-child { position: sticky; left: 0; background: inherit; z-index: 1; }
thead th:first-child { z-index: 4; }
/* Series cell with inline source pills */
.series-cell { display: flex; align-items: center; gap: .5rem; }
.series-title { font-weight: 600; }
.links { display: inline-flex; gap: .35rem; min-width: 4rem; }
.pill { display: inline-flex; align-items: center; justify-content: center; font-size: .72rem; line-height: 1; padding: .28rem .45rem; border-radius: 999px; text-decoration: none; border: 1px solid rgba(0,0,0,.06); width: 1.8rem; }
.pill-aud { background: rgba(14,165,233,.08); color: var(--aud); }
.pill-amz { background: rgba(245,158,11,.10); color: var(--amz); }
.badge { display: inline-block; min-width: 1.5em; padding: .15rem .5rem; border-radius: .5rem; background: #eef2ff; font-weight: 600; text-align: center; }
.count-aud { background: rgba(14,165,233,.12); color: #0369a1; }
.count-amz { background: rgba(245,158,11,.16); color: #92400e; }
.date { white-space: nowrap; color: var(--text); }
.sortable { cursor: pointer; user-select: none; }
.sortable::after { content: '\25B4\25BE'; font-size: .7em; opacity: .35; margin-left: .35rem; }
th.sort-asc::after { content: '\25B4'; opacity: .8; }
th.sort-desc::after { content: '\25BE'; opacity: .8; }
/* Top bar and settings panel */
.topbar { 
  display: flex; 
  align-items: center; 
  justify-content: space-between; 
  gap: 1rem; 
  padding: .75rem 1.5rem;
  background: var(--bg);
  border-bottom: 1px solid var(--line);
  position: sticky;
  top: 0;
  z-index: 10;
  flex-shrink: 0;
}
.topbar h1 {
  margin: 0;
  font-size: 1.44rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: .5rem;
}

.logo {
  height: 1.99rem;
  width: auto;
}
.settings-btn { display: inline-flex; align-items: center; justify-content: center; width: 2.25rem; height: 2.25rem; border-radius: .5rem; border: 1px solid var(--line); background: #fff; box-shadow: var(--head-shadow); cursor: pointer; font-size: 1.05rem; }
.settings-btn:hover { background: #f8fafc; }
.settings-btn:focus { outline: 2px solid #93c5fd; outline-offset: 2px; }

.loading-spinner { 
  display: inline-flex; 
  align-items: center; 
  justify-content: center; 
  width: 2.25rem; 
  height: 2.25rem; 
  margin-right: 0.5rem;
  opacity: 0;
  transition: opacity 0.3s ease;
  position: relative;
}
.loading-spinner.active { opacity: 1; }

.loading-spinner .dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: #6b7280;
  position: absolute;
  animation: spin-dots 1.2s linear infinite;
}

.loading-spinner .dot:nth-child(1) { transform: rotate(0deg) translate(8px) rotate(0deg); animation-delay: 0s; }
.loading-spinner .dot:nth-child(2) { transform: rotate(45deg) translate(8px) rotate(-45deg); animation-delay: -0.15s; }
.loading-spinner .dot:nth-child(3) { transform: rotate(90deg) translate(8px) rotate(-90deg); animation-delay: -0.3s; }
.loading-spinner .dot:nth-child(4) { transform: rotate(135deg) translate(8px) rotate(-135deg); animation-delay: -0.45s; }
.loading-spinner .dot:nth-child(5) { transform: rotate(180deg) translate(8px) rotate(-180deg); animation-delay: -0.6s; }
.loading-spinner .dot:nth-child(6) { transform: rotate(225deg) translate(8px) rotate(-225deg); animation-delay: -0.75s; }
.loading-spinner .dot:nth-child(7) { transform: rotate(270deg) translate(8px) rotate(-270deg); animation-delay: -0.9s; }
.loading-spinner .dot:nth-child(8) { transform: rotate(315deg) translate(8px) rotate(-315deg); animation-delay: -1.05s; }

@keyframes spin-dots {
  0%, 20% { opacity: 0.2; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1); }
  80%, 100% { opacity: 0.2; transform: scale(0.8); }
}

[data-theme="dark"] .loading-spinner .dot {
  background: #9ca3af;
}
.settings-wrap { position: relative; }
.settings-panel { position: absolute; right: 0; top: 2.8rem; width: 320px; max-width: calc(100vw - 2rem); background: #fff; border: 1px solid var(--line); border-radius: .5rem; box-shadow: 0 10px 20px rgba(0,0,0,.08), 0 2px 6px rgba(0,0,0,.06); padding: .75rem; z-index: 10; }
.settings-panel .panel-section { padding: .5rem .25rem; }
.settings-panel .panel-heading { font-weight: 700; font-size: .85rem; color: var(--muted); margin-bottom: .25rem; text-transform: uppercase; letter-spacing: .02em; }
.settings-panel code { background: #f3f4f6; padding: .15rem .35rem; border-radius: .35rem; }

/* Theme toggle */
.theme-toggle {
  display: flex;
  align-items: center;
  gap: .75rem;
}

.theme-label {
  font-size: .85rem;
  color: var(--muted);
  font-weight: 500;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 3rem;
  height: 1.5rem;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #e5e7eb;
  transition: .3s;
  border-radius: 1.5rem;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 1.125rem;
  width: 1.125rem;
  left: .1875rem;
  bottom: .1875rem;
  background: white;
  transition: .3s;
  border-radius: 50%;
  box-shadow: 0 1px 3px rgba(0,0,0,.2);
}

input:checked + .toggle-slider {
  background: #3b82f6;
}

input:checked + .toggle-slider:before {
  transform: translateX(1.5rem);
}

/* Dark mode adjustments */
[data-theme="dark"] .settings-panel code {
  background: #2b2b2b;
  color: #e5e7eb;
}

/* Export button styling */
.export-btn {
  display: inline-flex;
  align-items: center;
  gap: .5rem;
  padding: .5rem .75rem;
  background: #3b82f6;
  color: white;
  text-decoration: none;
  border-radius: .375rem;
  font-size: .85rem;
  font-weight: 500;
  transition: background-color .15s ease;
}

.export-btn:hover {
  background: #2563eb;
  color: white;
}

.export-btn:active {
  background: #1d4ed8;
}

[data-theme="dark"] .export-btn {
  background: #2563eb;
}

[data-theme="dark"] .export-btn:hover {
  background: #1d4ed8;
}

/* Subscription info styling */
.subscription-info {
  margin-top: .75rem;
  padding-top: .5rem;
  border-top: 1px solid var(--line);
}

.subscription-info small {
  display: block;
  margin-bottom: .35rem;
  color: var(--muted);
}

.subscription-url {
  display: block;
  padding: .4rem .5rem;
  background: #f8fafc;
  border: 1px solid var(--line);
  border-radius: .25rem;
  font-size: .75rem;
  word-break: break-all;
  user-select: all;
  cursor: pointer;
}

.subscription-url:hover {
  background: #f1f5f9;
}

[data-theme="dark"] .subscription-url {
  background: #2b2b2b;
  color: #e5e7eb;
  border-color: var(--line);
}

[data-theme="dark"] .subscription-url:hover {
  background: #374151;
}

[data-theme="dark"] .toggle-slider {
  background: #4b5563;
}

[data-theme="dark"] .toggle-slider:before {
  background: #f3f4f6;
}

/* Auto refresh slider styling */
.auto-refresh-setting {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.auto-refresh-label {
  font-size: 0.9rem;
  color: var(--text);
  text-align: center;
  font-weight: 500;
}

.auto-refresh-slider {
  width: 100%;
  height: 4px;
  border-radius: 2px;
  background: #e5e7eb;
  outline: none;
  -webkit-appearance: none;
  cursor: pointer;
}

.auto-refresh-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.auto-refresh-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
  border: none;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.auto-refresh-markers {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--muted);
  margin-top: -0.25rem;
}

[data-theme="dark"] .auto-refresh-slider {
  background: #4b5563;
}

[data-theme="dark"] .auto-refresh-slider::-webkit-slider-thumb {
  background: #60a5fa;
}

[data-theme="dark"] .auto-refresh-slider::-moz-range-thumb {
  background: #60a5fa;
}

/* Dark mode count badge styling for better visibility */
[data-theme="dark"] .count-aud {
  background: rgba(14,165,233,.25);
  color: #60a5fa;
}

[data-theme="dark"] .count-amz {
  background: rgba(245,158,11,.25);
  color: #fbbf24;
}

/* Dark mode mobile card adjustments */
@media (max-width: 768px) {
  [data-theme="dark"] tbody tr {
    background: #313335;
    border-color: var(--line);
  }
  
  [data-theme="dark"] tbody tr:nth-child(even) {
    background: #313335;
  }
  
  [data-theme="dark"] tbody tr:hover {
    background: #313335;
  }
  
  [data-theme="dark"] tbody td:first-child {
    background: #3c3f41;
    border-color: var(--line);
  }
}

/* Mobile responsive styles */
@media (max-width: 768px) {
  .topbar { padding: .5rem 1rem; }
  .topbar h1 { font-size: 1.27rem; }
  .content-area { padding: 1rem; }
  
  /* Hide table structure on mobile */
  table, thead, tbody, th, td, tr { 
    display: block; 
  }
  
  /* Hide table headers */
  thead tr { 
    position: absolute;
    top: -9999px;
    left: -9999px;
  }
  
  /* Style each row as a card */
  tbody tr {
    border: 1px solid var(--line);
    border-radius: .5rem;
    margin-bottom: 1rem;
    padding: 0;
    background: #fff;
    box-shadow: 0 1px 3px rgba(0,0,0,.1);
  }
  
  tbody tr:nth-child(even) {
    background: #fff;
  }
  
  tbody tr:hover {
    background: #fff;
    box-shadow: 0 2px 6px rgba(0,0,0,.15);
  }
  
  /* Remove sticky positioning */
  th:first-child, td:first-child { 
    position: static;
    z-index: auto;
  }
  
  /* Series title and links at the top */
  tbody td:first-child {
    padding: .75rem;
    border-bottom: 1px solid var(--line);
    background: var(--head-bg);
    border-radius: .5rem .5rem 0 0;
  }
  
  .series-cell {
    flex-direction: column;
    align-items: flex-start;
    gap: .75rem;
  }
  
  .series-title {
    font-size: 1.1rem;
    line-height: 1.3;
  }
  
  .links {
    gap: .5rem;
  }
  
  .pill {
    font-size: .8rem;
    padding: .35rem .6rem;
  }
  
  /* Platform sections underneath */
  tbody td:nth-child(n+2) {
    border-bottom: none;
    padding: .5rem .75rem;
  }
  
  /* Group Audible data */
  tbody td:nth-child(2) {
    border-top: 1px solid var(--line);
    padding-top: .75rem;
    position: relative;
  }
  
  tbody td:nth-child(2)::before {
    content: "Audible";
    display: block;
    font-weight: 600;
    color: var(--aud);
    margin-bottom: .5rem;
    font-size: .9rem;
    text-transform: uppercase;
    letter-spacing: .03em;
  }
  
  /* Group Amazon data */
  tbody td:nth-child(5) {
    border-top: 1px solid var(--line);
    padding-top: .75rem;
    position: relative;
  }
  
  tbody td:nth-child(5)::before {
    content: "Amazon";
    display: block;
    font-weight: 600;
    color: var(--amz);
    margin-bottom: .5rem;
    font-size: .9rem;
    text-transform: uppercase;
    letter-spacing: .03em;
  }
  
  /* Style data cells with labels */
  tbody td:nth-child(n+2) {
    display: flex;
    justify-content: space-between;
    align-items: center;
    min-height: 2rem;
  }
  
  tbody td:nth-child(2)::after { content: "Count"; }
  tbody td:nth-child(3)::after { content: "Latest"; }
  tbody td:nth-child(4)::after { content: "Next"; }
  tbody td:nth-child(5)::after { content: "Count"; }
  tbody td:nth-child(6)::after { content: "Latest"; }
  tbody td:nth-child(7)::after { content: "Next"; }
  
  tbody td:nth-child(n+2)::after {
    font-weight: 500;
    color: var(--muted);
    font-size: .85rem;
    order: -1;
  }
  
  /* Last cell rounded bottom */
  tbody td:last-child {
    border-radius: 0 0 .5rem .5rem;
    padding-bottom: .75rem;
  }
  
  /* Adjust badges and dates for mobile */
  .badge {
    font-size: .9rem;
    min-width: 2em;
  }
  
  .date {
    font-size: .9rem;
  }
}

/* Tablet adjustments */
@media (max-width: 1024px) and (min-width: 769px) {
  body { margin: 1.5rem; }
  table { font-size: .9rem; }
  th, td { padding: .45rem .6rem; }
  .pill { font-size: .7rem; }
}

/* Alphabetical index */
.alpha-index {
  position: fixed;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: .1rem;
  background: rgba(255,255,255,.9);
  backdrop-filter: blur(10px);
  border-radius: .5rem;
  padding: .5rem .25rem;
  border: 1px solid var(--line);
  box-shadow: 0 4px 12px rgba(0,0,0,.1);
  z-index: 20;
  user-select: none;
  opacity: 0;
  visibility: hidden;
  transition: opacity .3s ease, visibility .3s ease, transform .3s ease;
  transform: translateY(-50%) translateX(20px);
}

.alpha-index.visible {
  opacity: 1;
  visibility: visible;
  transform: translateY(-50%) translateX(0);
}

.alpha-letter {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  font-size: .7rem;
  font-weight: 600;
  color: var(--muted);
  cursor: pointer;
  border-radius: .25rem;
  transition: all .15s ease;
  text-decoration: none;
}

.alpha-letter:hover,
.alpha-letter:focus {
  background: var(--row-hover);
  color: var(--text);
  outline: none;
}

.alpha-letter.active {
  background: #3b82f6;
  color: white;
}

.alpha-letter.disabled {
  color: #d1d5db;
  cursor: default;
  pointer-events: none;
}

/* Hide on desktop by default, show on mobile/tablet */
@media (max-width: 1024px) {
  .alpha-index {
    display: flex;
  }
}

@media (min-width: 1025px) {
  .alpha-index {
    display: none;
  }
}

/* Mobile specific adjustments */
@media (max-width: 768px) {
  .alpha-index {
    right: .5rem;
    padding: .4rem .2rem;
  }
  
  .alpha-letter {
    width: 1.3rem;
    height: 1.3rem;
    font-size: .65rem;
  }
}

/* Modal styles */
.modal {
  position: fixed;
  z-index: 1000;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  backdrop-filter: blur(4px);
}

.modal-content {
  background: var(--bg);
  margin: 5% auto;
  padding: 0;
  border: 1px solid var(--line);
  border-radius: 0.5rem;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--line);
  background: var(--head-bg);
  border-radius: 0.5rem 0.5rem 0 0;
}

.modal-header h2 {
  margin: 0;
  color: var(--text);
  font-size: 1.5rem;
}

.close {
  color: var(--muted);
  font-size: 2rem;
  font-weight: bold;
  cursor: pointer;
  line-height: 1;
}

.close:hover,
.close:focus {
  color: var(--text);
}

.modal-body {
  padding: 1.5rem;
}

.user-section {
  margin-bottom: 2rem;
}

.user-section h3 {
  margin: 0 0 1rem 0;
  color: var(--text);
  font-size: 1.2rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: 500;
  color: var(--text);
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--line);
  border-radius: 0.375rem;
  background: var(--bg);
  color: var(--text);
  font-size: 1rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.message {
  padding: 0.75rem;
  border-radius: 0.375rem;
  margin-top: 1rem;
}

.message.success {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.message.error {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

[data-theme="dark"] .message.success {
  background: #064e3b;
  color: #6ee7b7;
  border-color: #047857;
}

[data-theme="dark"] .message.error {
  background: #7f1d1d;
  color: #fca5a5;
  border-color: #dc2626;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.users-table th,
.users-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid var(--line);
}

.users-table th {
  background: var(--head-bg);
  font-weight: 600;
  color: var(--text);
}

.users-table td {
  color: var(--text);
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
}

.role-admin {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.role-user {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

[data-theme="dark"] .role-admin {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
}

[data-theme="dark"] .role-user {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
}

.reset-btn {
  padding: 0.25rem 0.5rem;
  background: #f59e0b;
  color: white;
  border: none;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.reset-btn:hover {
  background: #d97706;
}

[data-theme="dark"] .reset-btn {
  background: #f59e0b;
}

[data-theme="dark"] .reset-btn:hover {
  background: #d97706;
}
</style>
<link rel="icon" href="/static/favicon.ico" type="image/x-icon">
</head>
<body>
  <div class="topbar">
    <h1>
      <img src="/static/syllabus_logo.png" alt="Syllabus Logo" class="logo">
      Syllabus
    </h1>
    <div class="settings-wrap">
      <div class="loading-spinner" id="loadingSpinner" title="Processing data...">
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
      </div>
      <button class="settings-btn" id="settingsBtn" aria-expanded="false" aria-controls="settingsPanel" title="Settings" aria-label="Settings">⚙️</button>
      <div class="settings-panel" id="settingsPanel" hidden>
        <div class="panel-section">
          <div class="panel-heading">Theme</div>
          <div class="panel-content">
            <div class="theme-toggle">
              <span class="theme-label">Light</span>
              <label class="toggle-switch">
                <input type="checkbox" id="themeToggle">
                <span class="toggle-slider"></span>
              </label>
              <span class="theme-label">Dark</span>
            </div>
          </div>
        </div>
        <div class="panel-section">
          <div class="panel-heading">Authentication</div>
          <div class="panel-content">
            {{ if .Authenticated }}
              <div style="margin-bottom: 0.75rem;">
                <small>Logged in as: <strong>{{ .User.Username }}</strong> ({{ .User.Role }})</small>
              </div>
              {{ if eq .User.Role "admin" }}
                <button onclick="showUserManagement()" class="export-btn" style="background: #059669; margin-right: 0.5rem;">
                  👥 Manage Users
                </button>
              {{ end }}
              <button onclick="logout()" class="export-btn" style="background: #dc2626; color: white;">
                🚪 Sign Out
              </button>
            {{ else }}
              <a href="/login" class="export-btn">
                🔐 Sign In
              </a>
            {{ end }}
          </div>
        </div>
        <div class="panel-section">
          <div class="panel-heading">Data Refresh</div>
          <div class="panel-content">
            <button onclick="refreshData()" class="export-btn" id="refreshBtn">🔄 Refresh All Data</button>
            <div id="refreshStatus" style="margin-top: 0.5rem; font-size: 0.85rem; color: var(--muted);"></div>
          </div>
        </div>
        <div class="panel-section">
          <div class="panel-heading">Auto Refresh</div>
          <div class="panel-content">
            <div class="auto-refresh-setting">
              <label for="autoRefreshSlider" class="auto-refresh-label">Every <span id="autoRefreshValue">6</span> hours</label>
              <input type="range" id="autoRefreshSlider" min="2" max="10" step="2" value="6" class="auto-refresh-slider">
              <div class="auto-refresh-markers">
                <span>2h</span>
                <span>4h</span>
                <span>6h</span>
                <span>8h</span>
                <span>10h</span>
              </div>
            </div>
          </div>
        </div>
        <div class="panel-section">
          <div class="panel-heading">Calendar Subscription</div>
          <div class="panel-content">
            <a href="/calendar.ics" target="_blank" class="export-btn">📅 Subscribe to iCal</a>
            <div class="subscription-info">
              <small>Copy this URL to subscribe in your calendar app:</small>
              <code class="subscription-url">{{ .CalendarURL }}</code>
            </div>
          </div>
        </div>
        <div class="panel-section">
          <div class="panel-heading">Generated at</div>
          <div class="panel-content"><code>{{ .Now }}</code></div>
        </div>
      </div>
    </div>
  </div>
  
  <div class="content-area">
    <table>
    <thead>
      <tr>
        <th rowspan="2" scope="col" class="sortable" data-col="0" data-type="text">Series</th>
        <th colspan="3" scope="colgroup">Audible</th>
        <th colspan="3" scope="colgroup">Amazon</th>
      </tr>
      <tr>
        <th scope="col" title="Number of audiobooks in the series" class="sortable" data-col="1" data-type="number">Count</th>
        <th scope="col" class="sortable" data-col="2" data-type="date">Latest</th>
        <th scope="col" class="sortable" data-col="3" data-type="date">Next</th>
        <th scope="col" title="Number of ebooks in the series on Amazon" class="sortable" data-col="4" data-type="number">Count</th>
        <th scope="col" class="sortable" data-col="5" data-type="date">Latest</th>
        <th scope="col" class="sortable" data-col="6" data-type="date">Next</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Rows }}
      <tr>
        <td>
          <div class="series-cell">
            <span class="series-title">{{ .Title }}</span>
            <span class="links">
              {{ if .AudibleURL }}<a class="pill pill-aud" href="{{ .AudibleURL }}" target="_blank" rel="noopener" aria-label="Open series on Audible">Au</a>{{ end }}
              {{ if .AmazonURL }}<a class="pill pill-amz" href="{{ .AmazonURL }}" target="_blank" rel="noopener" aria-label="Open series on Amazon">Am</a>{{ end }}
            </span>
          </div>
        </td>
        <td><span class="badge count-aud">{{ .AudibleCount }}</span></td>
        <td><span class="date">{{ .AudibleLatest }}</span></td>
        <td><span class="date">{{ .AudibleNext }}</span></td>
        <td><span class="badge count-amz">{{ .AmazonCount }}</span></td>
        <td><span class="date">{{ .AmazonLatest }}</span></td>
        <td><span class="date">{{ .AmazonNext }}</span></td>
      </tr>
      {{ end }}
    </tbody>
    </table>
  </div>
  
  <!-- Alphabetical Index -->
  <div class="alpha-index" id="alphaIndex">
    <div class="alpha-letter" data-letter="A">A</div>
    <div class="alpha-letter" data-letter="B">B</div>
    <div class="alpha-letter" data-letter="C">C</div>
    <div class="alpha-letter" data-letter="D">D</div>
    <div class="alpha-letter" data-letter="E">E</div>
    <div class="alpha-letter" data-letter="F">F</div>
    <div class="alpha-letter" data-letter="G">G</div>
    <div class="alpha-letter" data-letter="H">H</div>
    <div class="alpha-letter" data-letter="I">I</div>
    <div class="alpha-letter" data-letter="J">J</div>
    <div class="alpha-letter" data-letter="K">K</div>
    <div class="alpha-letter" data-letter="L">L</div>
    <div class="alpha-letter" data-letter="M">M</div>
    <div class="alpha-letter" data-letter="N">N</div>
    <div class="alpha-letter" data-letter="O">O</div>
    <div class="alpha-letter" data-letter="P">P</div>
    <div class="alpha-letter" data-letter="Q">Q</div>
    <div class="alpha-letter" data-letter="R">R</div>
    <div class="alpha-letter" data-letter="S">S</div>
    <div class="alpha-letter" data-letter="T">T</div>
    <div class="alpha-letter" data-letter="U">U</div>
    <div class="alpha-letter" data-letter="V">V</div>
    <div class="alpha-letter" data-letter="W">W</div>
    <div class="alpha-letter" data-letter="X">X</div>
    <div class="alpha-letter" data-letter="Y">Y</div>
    <div class="alpha-letter" data-letter="Z">Z</div>
  </div>
  
  <!-- User Management Modal -->
  <div id="userModal" class="modal" style="display: none;">
    <div class="modal-content">
      <div class="modal-header">
        <h2>User Management</h2>
        <span class="close" onclick="closeUserManagement()">&times;</span>
      </div>
      <div class="modal-body">
        <div class="user-section">
          <h3>Create New User</h3>
          <form id="createUserForm">
            <div class="form-group">
              <label for="newUsername">Username:</label>
              <input type="text" id="newUsername" required>
            </div>
            <div class="form-group">
              <label for="newPassword">Password:</label>
              <input type="password" id="newPassword" required>
            </div>
            <div class="form-group">
              <label for="newRole">Role:</label>
              <select id="newRole">
                <option value="user">User</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <button type="submit" class="export-btn">Create User</button>
          </form>
          <div id="createUserMessage" class="message" style="display: none;"></div>
        </div>
        
        <div class="user-section">
          <h3>Existing Users</h3>
          <div id="usersList">Loading...</div>
        </div>
      </div>
    </div>
  </div>
  
  <!-- Password Reset Modal -->
  <div id="resetPasswordModal" class="modal" style="display: none;">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Reset Password</h2>
        <span class="close" onclick="closeResetPasswordModal()">&times;</span>
      </div>
      <div class="modal-body">
        <form id="resetPasswordForm">
          <div class="form-group">
            <label for="resetUsername">Username:</label>
            <input type="text" id="resetUsername" readonly>
          </div>
          <div class="form-group">
            <label for="resetNewPassword">New Password:</label>
            <input type="password" id="resetNewPassword" required>
          </div>
          <div class="form-group">
            <label for="resetConfirmPassword">Confirm New Password:</label>
            <input type="password" id="resetConfirmPassword" required>
          </div>
          <button type="submit" class="export-btn" style="background: #f59e0b;">Reset Password</button>
          <button type="button" onclick="closeResetPasswordModal()" class="export-btn" style="background: #6b7280; margin-left: 0.5rem;">Cancel</button>
        </form>
        <div id="resetPasswordMessage" class="message" style="display: none;"></div>
      </div>
    </div>
  </div>
  
  <script>
  (function(){
    const table = document.querySelector('table');
    if(!table) return;
    const tbody = table.querySelector('tbody');
    const getText = (cell) => (cell.textContent || '').trim();
    const parseNumber = (s) => {
      const m = (s.match(/[-+]?[0-9]*\.?[0-9]+/)||[])[0];
      if(m === undefined || m === '') return NaN;
      return parseFloat(m);
    };
    const parseDate = (s) => {
      // Look for YYYY-MM-DD anywhere in the string
      const m = s.match(/\b(\d{4})-(\d{2})-(\d{2})\b/);
      if(m){
        const t = Date.parse(m[0] + 'T00:00:00Z');
        return isNaN(t) ? null : t;
      }
      // Fallback: try native Date
      const t = Date.parse(s);
      return isNaN(t) ? null : t;
    };
    const comparators = {
      text: (a,b) => a.localeCompare(b, undefined, {numeric:true, sensitivity:'base'}),
      number: (a,b) => (a - b),
      date: (a,b) => (a - b)
    };
    const extractors = {
      text: (cell) => getText(cell).toLowerCase(),
      number: (cell) => parseNumber(getText(cell)),
      date: (cell) => { const t = parseDate(getText(cell)); return t===null? Number.NEGATIVE_INFINITY : t; }
    };
    const clearSortStates = () => table.querySelectorAll('th.sort-asc, th.sort-desc').forEach(th=>{ th.classList.remove('sort-asc','sort-desc'); th.removeAttribute('aria-sort'); });
    const sortBy = (col, type, direction) => {
      const rows = Array.from(tbody.querySelectorAll('tr'));
      const idx = col|0;
      const getVal = (row) => extractors[type](row.children[idx]);
      const cmp = comparators[type];
      rows.sort((r1, r2) => {
        const a = getVal(r1); const b = getVal(r2);
        const c = cmp(a,b);
        return direction === 'desc' ? -c : c;
      });
      // Re-append in sorted order
      rows.forEach(r => tbody.appendChild(r));
    };
    table.querySelectorAll('th.sortable').forEach(th => {
      th.setAttribute('role', 'button');
      th.tabIndex = 0;
      let dir = th.dataset.defaultDir || 'asc';
      th.addEventListener('click', () => {
        const col = parseInt(th.dataset.col,10);
        const type = th.dataset.type || 'text';
        clearSortStates();
        sortBy(col, type, dir);
        th.classList.add(dir==='asc'?'sort-asc':'sort-desc');
        th.setAttribute('aria-sort', dir==='asc'?'ascending':'descending');
        dir = (dir === 'asc') ? 'desc' : 'asc';
      });
      th.addEventListener('keydown', (e) => { if(e.key==='Enter' || e.key===' '){ e.preventDefault(); th.click(); }});
    });
    // Settings panel toggle
    const btn = document.getElementById('settingsBtn');
    const panel = document.getElementById('settingsPanel');
    if (btn && panel) {
      const closePanel = () => { panel.hidden = true; btn.setAttribute('aria-expanded','false'); };
      const openPanel  = () => { panel.hidden = false; btn.setAttribute('aria-expanded','true'); };
      btn.addEventListener('click', (e) => {
        e.stopPropagation();
        if (panel.hidden) openPanel(); else closePanel();
      });
      document.addEventListener('click', (e) => {
        if (panel.hidden) return;
        if (!panel.contains(e.target) && e.target !== btn) closePanel();
      }, true);
      document.addEventListener('keydown', (e) => { if (e.key === 'Escape') closePanel(); });
    }

    // Theme toggle functionality
    const themeToggle = document.getElementById('themeToggle');
    if (themeToggle) {
      // Load saved theme or default to light
      const savedTheme = localStorage.getItem('theme') || 'light';
      const isDark = savedTheme === 'dark';
      
      // Apply theme
      document.documentElement.setAttribute('data-theme', savedTheme);
      themeToggle.checked = isDark;
      
      // Handle toggle changes
      themeToggle.addEventListener('change', () => {
        const newTheme = themeToggle.checked ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
      });
    }

    // Auto refresh slider functionality
    const autoRefreshSlider = document.getElementById('autoRefreshSlider');
    const autoRefreshValue = document.getElementById('autoRefreshValue');
    if (autoRefreshSlider && autoRefreshValue) {
      // Load saved auto refresh interval or default to 6
      const savedInterval = localStorage.getItem('autoRefreshInterval') || '6';
      autoRefreshSlider.value = savedInterval;
      autoRefreshValue.textContent = savedInterval;
      
      // Handle slider changes
      autoRefreshSlider.addEventListener('input', () => {
        const value = autoRefreshSlider.value;
        autoRefreshValue.textContent = value;
        localStorage.setItem('autoRefreshInterval', value);
        
        // Send update to server
        fetch('/api/auto-refresh', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            interval: parseInt(value)
          })
        }).catch(error => {
          console.error('Error updating auto refresh interval:', error);
        });
      });
      
      // Set initial interval on server
      fetch('/api/auto-refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          interval: parseInt(savedInterval)
        })
      }).catch(error => {
        console.error('Error setting initial auto refresh interval:', error);
      });
    }

    // Alphabetical Index Navigation
    const alphaIndex = document.getElementById('alphaIndex');
    if (alphaIndex && tbody) {
      const rows = Array.from(tbody.querySelectorAll('tr'));
      let hideTimeout;
      let isInteracting = false;
      const availableLetters = new Set();
      
      // Find which letters have entries
      rows.forEach(row => {
        const titleCell = row.querySelector('td:first-child .series-title');
        if (titleCell) {
          const firstLetter = titleCell.textContent.trim().charAt(0).toUpperCase();
          if (firstLetter.match(/[A-Z]/)) {
            availableLetters.add(firstLetter);
          }
        }
      });
      
      // Auto-hide functionality
      const showIndex = () => {
        alphaIndex.classList.add('visible');
        clearTimeout(hideTimeout);
        hideTimeout = setTimeout(() => {
          if (!isInteracting) {
            alphaIndex.classList.remove('visible');
          }
        }, 2000); // Hide after 2 seconds of inactivity
      };
      
      const hideIndex = () => {
        if (!isInteracting) {
          alphaIndex.classList.remove('visible');
        }
      };
      
      const startInteracting = () => {
        isInteracting = true;
        showIndex();
      };
      
      const stopInteracting = () => {
        isInteracting = false;
        hideTimeout = setTimeout(hideIndex, 1000);
      };
      
      // Show on scroll
      let scrollTimeout;
      window.addEventListener('scroll', () => {
        showIndex();
        clearTimeout(scrollTimeout);
        scrollTimeout = setTimeout(() => {
          if (!isInteracting) {
            hideIndex();
          }
        }, 1500);
      });
      
      // Show on touch near the index area
      document.addEventListener('touchstart', (e) => {
        const touch = e.touches[0];
        const rightEdge = window.innerWidth - 100; // Show if touch within 100px of right edge
        if (touch.clientX > rightEdge) {
          showIndex();
        }
      });
      
      // Update letter states
      const letters = alphaIndex.querySelectorAll('.alpha-letter');
      letters.forEach(letter => {
        const letterValue = letter.dataset.letter;
        if (!availableLetters.has(letterValue)) {
          letter.classList.add('disabled');
        }
      });
      
      // Find first row starting with letter
      const findRowByLetter = (letter) => {
        return rows.find(row => {
          const titleCell = row.querySelector('td:first-child .series-title');
          if (titleCell) {
            const firstLetter = titleCell.textContent.trim().charAt(0).toUpperCase();
            return firstLetter === letter;
          }
          return false;
        });
      };
      
      // Smooth scroll to element
      const scrollToElement = (element) => {
        const offset = window.innerWidth <= 768 ? 80 : 120; // Account for mobile/desktop differences
        const elementPosition = element.getBoundingClientRect().top + window.scrollY - offset;
        
        window.scrollTo({
          top: elementPosition,
          behavior: 'smooth'
        });
      };
      
      // Handle letter clicks
      letters.forEach(letter => {
        if (!letter.classList.contains('disabled')) {
          letter.addEventListener('click', () => {
            const targetLetter = letter.dataset.letter;
            const targetRow = findRowByLetter(targetLetter);
            
            if (targetRow) {
              // Remove previous active state
              letters.forEach(l => l.classList.remove('active'));
              // Add active state to clicked letter
              letter.classList.add('active');
              
              // Scroll to the row
              scrollToElement(targetRow);
              
              // Remove active state after a short delay
              setTimeout(() => {
                letter.classList.remove('active');
              }, 1500);
            }
          });
          
          // Show index when hovering over letters
          letter.addEventListener('mouseenter', showIndex);
        }
      });
      
      // Handle touch events for drag navigation
      let isDragging = false;
      let startY = 0;
      
      const handleTouchStart = (e) => {
        isDragging = true;
        startY = e.touches[0].clientY;
        startInteracting();
        e.preventDefault();
      };
      
      const handleTouchMove = (e) => {
        if (!isDragging) return;
        
        const currentY = e.touches[0].clientY;
        const indexRect = alphaIndex.getBoundingClientRect();
        const relativeY = currentY - indexRect.top;
        
        // Find which letter we're over
        const letterHeight = indexRect.height / letters.length;
        const letterIndex = Math.floor(relativeY / letterHeight);
        
        if (letterIndex >= 0 && letterIndex < letters.length) {
          const targetLetter = letters[letterIndex];
          if (!targetLetter.classList.contains('disabled')) {
            const targetRow = findRowByLetter(targetLetter.dataset.letter);
            if (targetRow) {
              // Clear all active states
              letters.forEach(l => l.classList.remove('active'));
              // Set current letter as active
              targetLetter.classList.add('active');
              // Scroll to row
              scrollToElement(targetRow);
            }
          }
        }
        
        e.preventDefault();
      };
      
      const handleTouchEnd = (e) => {
        isDragging = false;
        stopInteracting();
        // Remove all active states after touch ends
        setTimeout(() => {
          letters.forEach(l => l.classList.remove('active'));
        }, 1000);
        e.preventDefault();
      };
      
      // Add touch event listeners
      alphaIndex.addEventListener('touchstart', handleTouchStart, { passive: false });
      alphaIndex.addEventListener('touchmove', handleTouchMove, { passive: false });
      alphaIndex.addEventListener('touchend', handleTouchEnd, { passive: false });
    }
  })();

  // Authentication functions
  function logout() {
    fetch('/logout', {
      method: 'POST',
      headers: {
        'Accept': 'application/json'
      }
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        window.location.href = '/login';
      }
    })
    .catch(error => {
      console.error('Logout error:', error);
      // Fallback to redirect
      window.location.href = '/logout';
    });
  }

  // User management functions
  function showUserManagement() {
    document.getElementById('userModal').style.display = 'block';
    loadUsers();
  }

  function closeUserManagement() {
    document.getElementById('userModal').style.display = 'none';
    clearCreateUserForm();
  }

  function clearCreateUserForm() {
    document.getElementById('createUserForm').reset();
    const message = document.getElementById('createUserMessage');
    message.style.display = 'none';
    message.className = 'message';
  }

  function showMessage(elementId, message, isSuccess) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.className = 'message ' + (isSuccess ? 'success' : 'error');
    element.style.display = 'block';
  }

  function loadUsers() {
    fetch('/api/users')
      .then(response => response.json())
      .then(data => {
        const usersContainer = document.getElementById('usersList');
        if (data.users && data.users.length > 0) {
          const table = document.createElement('table');
          table.className = 'users-table';
          
          const header = table.createTHead();
          const headerRow = header.insertRow();
          headerRow.innerHTML = '<th>Username</th><th>Role</th><th>Created</th><th>Actions</th>';
          
          const tbody = table.createTBody();
          data.users.forEach(user => {
            const row = tbody.insertRow();
            const createdDate = new Date(user.created_at).toLocaleDateString();
            row.innerHTML = '<td>' + user.username + '</td><td><span class="role-badge role-' + user.role + '">' + user.role + '</span></td><td>' + createdDate + '</td><td><button class="reset-btn" onclick="showResetPasswordModal(\'' + user.username + '\')">Reset Password</button></td>';
          });
          
          usersContainer.innerHTML = '';
          usersContainer.appendChild(table);
        } else {
          usersContainer.innerHTML = '<p>No users found.</p>';
        }
      })
      .catch(error => {
        console.error('Error loading users:', error);
        document.getElementById('usersList').innerHTML = '<p>Error loading users.</p>';
      });
  }

  // Handle create user form submission
  document.getElementById('createUserForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const username = document.getElementById('newUsername').value;
    const password = document.getElementById('newPassword').value;
    const role = document.getElementById('newRole').value;
    
    fetch('/api/users/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: username,
        password: password,
        role: role
      })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        showMessage('createUserMessage', data.message, true);
        clearCreateUserForm();
        loadUsers(); // Refresh the user list
      } else {
        showMessage('createUserMessage', data.message, false);
      }
    })
    .catch(error => {
      console.error('Error creating user:', error);
      showMessage('createUserMessage', 'Error creating user', false);
    });
  });

  // Close modal when clicking outside of it
  window.addEventListener('click', function(event) {
    const modal = document.getElementById('userModal');
    if (event.target === modal) {
      closeUserManagement();
    }
    
    const resetModal = document.getElementById('resetPasswordModal');
    if (event.target === resetModal) {
      closeResetPasswordModal();
    }
  });

  // Password reset functions
  function showResetPasswordModal(username) {
    document.getElementById('resetUsername').value = username;
    document.getElementById('resetNewPassword').value = '';
    document.getElementById('resetConfirmPassword').value = '';
    const message = document.getElementById('resetPasswordMessage');
    message.style.display = 'none';
    message.className = 'message';
    document.getElementById('resetPasswordModal').style.display = 'block';
  }

  function closeResetPasswordModal() {
    document.getElementById('resetPasswordModal').style.display = 'none';
  }

  // Handle reset password form submission
  document.getElementById('resetPasswordForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const username = document.getElementById('resetUsername').value;
    const newPassword = document.getElementById('resetNewPassword').value;
    const confirmPassword = document.getElementById('resetConfirmPassword').value;
    
    // Validate passwords match
    if (newPassword !== confirmPassword) {
      showMessage('resetPasswordMessage', 'Passwords do not match', false);
      return;
    }
    
    // Validate password length
    if (newPassword.length < 3) {
      showMessage('resetPasswordMessage', 'Password must be at least 3 characters long', false);
      return;
    }
    
    fetch('/api/users/reset-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: username,
        new_password: newPassword
      })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        showMessage('resetPasswordMessage', data.message, true);
        // Clear form fields
        document.getElementById('resetNewPassword').value = '';
        document.getElementById('resetConfirmPassword').value = '';
        // Close modal after a short delay
        setTimeout(() => {
          closeResetPasswordModal();
        }, 1500);
      } else {
        showMessage('resetPasswordMessage', data.message, false);
      }
    })
    .catch(error => {
      console.error('Error resetting password:', error);
      showMessage('resetPasswordMessage', 'Error resetting password', false);
    });
  });

  // Live refresh via Server-Sent Events
  (function() {
    const eventSource = new EventSource('/events');
    const spinner = document.getElementById('loadingSpinner');
    let activeJobs = 0;
    let lastActivity = Date.now();
    
    function updateSpinner() {
      if (activeJobs > 0) {
        spinner.classList.add('active');
        lastActivity = Date.now();
      } else {
        // Keep spinner visible for a short time after last activity
        setTimeout(() => {
          if (Date.now() - lastActivity > 3000) { // 3 seconds delay
            spinner.classList.remove('active');
          }
        }, 3000);
      }
    }
    
    eventSource.onmessage = function(event) {
      if (event.data === 'refresh') {
        // Show a brief notification
        const notification = document.createElement('div');
        notification.textContent = 'New entries added - refreshing...';
        notification.style.cssText = 'position: fixed; top: 20px; right: 20px; background: var(--bg); color: var(--text); border: 1px solid var(--line); padding: 0.75rem 1rem; border-radius: 0.5rem; box-shadow: 0 4px 12px rgba(0,0,0,0.15); z-index: 1000; font-size: 0.9rem;';
        document.body.appendChild(notification);
        
        // Refresh the page after a brief delay
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      } else if (event.data === 'connected') {
        // Initial connection - show spinner for initial activity
        activeJobs = 10; // Assume initial scraping activity
        updateSpinner();
        // Gradually reduce over time if no updates come
        setTimeout(() => {
          if (activeJobs >= 10) {
            activeJobs = Math.max(0, activeJobs - 5);
            updateSpinner();
          }
        }, 10000); // 10 seconds
      } else {
        // Try to parse scraper update JSON
        try {
          const update = JSON.parse(event.data);
          if (update.Status === 'running') {
            activeJobs = Math.max(activeJobs, 1);
            updateSpinner();
          } else if (update.Status === 'completed' || update.Status === 'failed') {
            activeJobs = Math.max(0, activeJobs - 1);
            updateSpinner();
          }
          // Any scraper update indicates activity
          lastActivity = Date.now();
        } catch (e) {
          // Not JSON, ignore
        }
      }
    };
    
    eventSource.onerror = function() {
      console.log('SSE connection lost, will retry automatically');
      activeJobs = 0;
      updateSpinner();
    };
    
    // Show spinner initially for a few seconds to indicate loading
    activeJobs = 1;
    updateSpinner();
    setTimeout(() => {
      if (activeJobs <= 1) {
        activeJobs = 0;
        updateSpinner();
      }
    }, 5000);
  })();

  // Refresh data function
  async function refreshData() {
    const btn = document.getElementById('refreshBtn');
    const status = document.getElementById('refreshStatus');
    
    // Disable button and show loading state
    btn.disabled = true;
    btn.innerHTML = '⏳ Refreshing...';
    status.textContent = 'Triggering data refresh...';
    status.style.color = '#3b82f6';
    
    try {
      const response = await fetch('/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        status.textContent = 'Refresh started for ' + data.seriesCount + ' series. The page will update automatically.';
        status.style.color = '#059669';
        
        // Listen for real-time updates
        if (typeof eventSource !== 'undefined') {
          let completedJobs = 0;
          const totalJobs = data.jobsQueued || 0;
          
          const originalHandler = eventSource.onmessage;
          eventSource.onmessage = function(event) {
            if (event.data === 'refresh') {
              // Page will reload automatically
              location.reload();
            } else if (event.data.startsWith('{')) {
              try {
                const update = JSON.parse(event.data);
                if (update.status === 'completed' || update.status === 'failed') {
                  completedJobs++;
                  status.textContent = 'Progress: ' + completedJobs + '/' + totalJobs + ' series updated';
                  
                  if (completedJobs >= totalJobs) {
                    setTimeout(() => {
                      status.textContent = 'Refresh complete! Page reloading...';
                      location.reload();
                    }, 1000);
                  }
                }
              } catch (e) {
                // Ignore parse errors
              }
            }
            
            // Call original handler
            if (originalHandler) originalHandler(event);
          };
        }
      } else {
        throw new Error('HTTP ' + response.status);
      }
    } catch (error) {
      console.error('Refresh error:', error);
      status.textContent = 'Error triggering refresh. Please try again.';
      status.style.color = '#dc2626';
    } finally {
      // Re-enable button after a delay
      setTimeout(() => {
        btn.disabled = false;
        btn.innerHTML = '🔄 Refresh All Data';
      }, 2000);
    }
  }
  </script>
</body>
</html>
`

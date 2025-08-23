package handlers

const IndexHTML = `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>syllabus</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="icon" href="/static/favicon.ico" type="image/x-icon">
<link rel="shortcut icon" href="/static/favicon.ico" type="image/x-icon">
<link rel="apple-touch-icon" href="/static/favicon.ico">
<style>
:root{
  --bg:#ffffff;--text:#111827;--muted:#6b7280;--line:#e5e7eb;--head-bg:#f9fafb;--row-hover:#f3f4f6;
  --aud:#0ea5e9;--amz:#f59e0b
}
[data-theme="dark"]{--bg:#2b2b2b;--text:#a9b7c6;--muted:#808080;--line:#3c3f41;--head-bg:#3c3f41;--row-hover:#3c3f41}

html,body{height:100%}
body{margin:0;display:flex;flex-direction:column;height:100vh;overflow:hidden;background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif}

/* ── top bar ─────────────────────────────────────────── */
.top-bar{height:75px;background:#f5f5f5;border-bottom:1px solid var(--line);display:flex;align-items:center;justify-content:space-between;padding:0 24px;box-shadow:0 2px 8px rgba(0,0,0,.05);z-index:200;flex-shrink:0}
.top-bar-left{display:flex;align-items:center;gap:16px}
.top-bar-logo{display:flex;align-items:center;gap:12px}
.top-bar-logo .logo{width:32px;height:32px;flex-shrink:0}
.top-bar-title{font-size:24px;font-weight:600;color:var(--text)}
.top-bar-search{flex:1;max-width:420px;margin:0 24px}
.search-input{width:100%;padding:10px 16px;border:1px solid var(--line);border-radius:8px;background:var(--bg);color:var(--text);font-size:14px}
.search-input:focus{outline:none;border-color:#3182ce;box-shadow:0 0 0 3px rgba(49,130,206,.1)}
.top-bar-right{display:flex;align-items:center;gap:16px}
.status-indicator{display:none;width:20px;height:20px;position:relative}
.status-dots{width:100%;height:100%;position:relative}
.status-dot{position:absolute;width:4px;height:4px;border-radius:50%;background:var(--aud);animation:statusPulse 1.5s infinite}
.status-dot:nth-child(1){top:2px;left:8px;animation-delay:0s}
.status-dot:nth-child(2){top:6px;left:14px;animation-delay:0.2s}
.status-dot:nth-child(3){top:14px;left:14px;animation-delay:0.4s}
.status-dot:nth-child(4){top:18px;left:8px;animation-delay:0.6s}
.status-dot:nth-child(5){top:14px;left:2px;animation-delay:0.8s}
.status-dot:nth-child(6){top:6px;left:2px;animation-delay:1.0s}
@keyframes statusPulse{0%,100%{opacity:0.3;transform:scale(1)}50%{opacity:1;transform:scale(1.2)}}
.settings-btn{background:none;border:1px solid var(--line);color:var(--muted);cursor:pointer;padding:8px;border-radius:6px;transition:.2s;background:#fff}
.settings-btn:hover{background:#f8fafc;color:var(--text)}
.user-box{position:relative;background:var(--bg);border:1px solid var(--line);border-radius:8px;padding:8px 12px;color:var(--text);font-size:14px;font-weight:500;cursor:pointer}
.user-dropdown{position:absolute;top:100%;right:0;margin-top:4px;background:var(--bg);border:1px solid var(--line);border-radius:8px;box-shadow:0 4px 12px rgba(0,0,0,.15);min-width:140px;display:none;z-index:300}
.user-dropdown-item{display:block;width:100%;padding:10px 16px;text-align:left;background:none;border:none;color:var(--text);font-size:14px;cursor:pointer;text-decoration:none}
.user-dropdown-item:hover{background:var(--row-hover)}
[data-theme="dark"] .user-dropdown{background:#313335;box-shadow:0 4px 12px rgba(0,0,0,.3)}
.user-list-item{display:flex;align-items:center;justify-content:space-between;padding:8px 12px;border:1px solid var(--line);border-radius:6px;margin-bottom:8px;background:var(--head-bg)}
.user-info{flex:1}
.user-name{font-weight:600;color:var(--text)}
.user-role{font-size:.85rem;color:var(--muted)}
.user-actions{display:flex;gap:4px}
.delete-user-btn{padding:4px 8px;background:#dc2626;color:white;border:none;border-radius:4px;cursor:pointer;font-size:12px}
.delete-user-btn:hover{background:#b91c1c}
.reset-password-btn{padding:4px 8px;background:#f59e0b;color:white;border:none;border-radius:4px;cursor:pointer;font-size:12px}
.reset-password-btn:hover{background:#d97706}
[data-theme="dark"] .user-list-item{background:#374151}
[data-theme="dark"] .top-bar{background:#3c3f41}

/* ── layout containers ───────────────────────────────── */
.main{flex:1;min-height:0;display:flex;flex-direction:column;overflow:hidden}
.container{flex:1;min-height:0;overflow:hidden;}

#desktopView{display:none}
#mobileView{display:block}
@media (min-width: 1024px){
  #desktopView{display:flex}
  #mobileView{display:none}
}

/* cards/tiles */
.card{border:1px solid var(--line);border-radius:12px;background:#fff;overflow:hidden}
[data-theme="dark"] .card{background:#313335}
.card-h{padding:14px 16px;border-bottom:1px solid var(--line);background:var(--head-bg);border-radius:12px 12px 0 0;font-weight:600}
.card-c{padding:16px}
.tile{border:1px solid var(--line);border-radius:12px;background:var(--head-bg);padding:14px 16px}
.tile .k{font-size:.85rem;color:var(--muted)}
.tile .v{font-size:1.6rem;font-weight:700;line-height:1.1;white-space:nowrap}

/* desktop layout */
.desktop-wrap{padding:24px;gap:16px;display:flex;flex-direction:column;min-height:0;height:100%}
.desktop-top-row{display:grid;grid-template-columns:1fr 200px 240px;gap:16px}
.desktop-body{display:grid;grid-template-columns:280px 1fr;gap:16px;min-height:0;flex:1}
.panel{min-height:0;display:flex;flex-direction:column}
.panel-scroll{min-height:0;flex:1;overflow:auto;padding-bottom:24px;scroll-padding-bottom:24px;margin-bottom:12px}

/* table */
.table{width:100%;border-collapse:separate;border-spacing:0;position:relative}
.table thead{position:sticky;top:0;z-index:10}
.table thead th{background:var(--head-bg);border-bottom:1px solid var(--line);text-align:left;font-weight:600;padding:14px 16px;cursor:pointer;user-select:none;position:relative}
.table thead th:first-child{border-top-left-radius:12px}
.table thead th:last-child{border-top-right-radius:12px}
.table thead th:hover{background:var(--row-hover)}
.table thead th.sortable{position:relative}
.table thead th.sort-asc::after{content:' ▲';color:var(--aud);font-size:10px}
.table thead th.sort-desc::after{content:' ▼';color:var(--aud);font-size:10px}
.table tbody td{padding:12px 16px;border-bottom:1px solid var(--line)}
.zebra tbody tr:nth-child(2n){background:rgba(0,0,0,.03)}
[data-theme="dark"] .zebra tbody tr:nth-child(2n){background:#2f3133}

/* pills/badges */
.badge{display:inline-flex;align-items:center;gap:6px;padding:2px 8px;border-radius:9999px;font-size:.78rem;font-weight:600;border:1px solid rgba(0,0,0,.06)}
.b-aud{background:rgba(14,165,233,.12);color:#0369a1}
.b-amz{background:rgba(245,158,11,.16);color:#92400e}
.next{display:inline-block;padding:2px 8px;border-radius:8px;font-size:.78rem;font-weight:600}
.next-aud{background:var(--aud);color:#fff}
.next-amz{background:var(--amz);color:#fff}
.next-none{background:#64748b;color:#fff}
.latest{display:inline-block;padding:2px 8px;border-radius:8px;font-size:.78rem;font-weight:600;background:#f1f5f9;color:#475569;border:1px solid #e2e8f0;text-align:center}
[data-theme="dark"] .latest{background:#374151;color:#9ca3af;border-color:#4b5563}
.next{text-align:center}
.next-none{background:#374151;color:#cbd5e1}
.linkpill{display:inline-flex;align-items:center;justify-content:center;font-size:.72rem;line-height:1;padding:.28rem .45rem;border-radius:999px;text-decoration:none;border:1px solid rgba(0,0,0,.06);min-width:1.8rem}
.link-aud{background:rgba(14,165,233,.08);color:var(--aud)}
.link-amz{background:rgba(245,158,11,.10);color:var(--amz)}

/* Icon styles */
.icon-headphones{
  display:inline-block;
  width:14px;
  height:14px;
  background:currentColor;
  mask:url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3e%3cpath d='M3 14h3a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2v-7a9 9 0 0 1 18 0v7a2 2 0 0 1-2 2h-2a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3'/%3e%3c/svg%3e") no-repeat center;
  mask-size:contain;
  -webkit-mask:url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3e%3cpath d='M3 14h3a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2v-7a9 9 0 0 1 18 0v7a2 2 0 0 1-2 2h-2a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3'/%3e%3c/svg%3e") no-repeat center;
  -webkit-mask-size:contain;
}
.icon-book{
  display:inline-block;
  width:14px;
  height:14px;
  background:currentColor;
  mask:url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3e%3cpath d='M4 19.5A2.5 2.5 0 0 1 6.5 17H20'/%3e%3cpath d='M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z'/%3e%3c/svg%3e") no-repeat center;
  mask-size:contain;
  -webkit-mask:url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3e%3cpath d='M4 19.5A2.5 2.5 0 0 1 6.5 17H20'/%3e%3cpath d='M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z'/%3e%3c/svg%3e") no-repeat center;
  -webkit-mask-size:contain;
}

/* sidebar checkbox list */
.checkrow{display:flex;align-items:center;gap:8px}
.checkrow input[type=checkbox]{width:16px;height:16px}

/* modal */
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.45);display:none;align-items:center;justify-content:center;z-index:500}
.modal-panel{background:#fff;border:1px solid var(--line);border-radius:12px;max-width:520px;width:92%;box-shadow:0 20px 60px rgba(0,0,0,.25)}
[data-theme="dark"] .modal-panel{background:#313335}
.modal-head{display:flex;align-items:center;justify-content:space-between;padding:14px 16px;border-bottom:1px solid var(--line);font-weight:700}
.modal-body{padding:16px}
.modal-row{display:flex;align-items:center;justify-content:space-between;gap:16px;padding:10px 0}
.switch{position:relative;width:46px;height:26px;flex-shrink:0}
.switch input{opacity:0;width:0;height:0}
.slider{position:absolute;cursor:pointer;inset:0;background:#cbd5e1;border-radius:9999px;transition:.2s}
.slider:before{content:"";position:absolute;height:20px;width:20px;left:3px;top:3px;background:white;border-radius:50%;transition:.2s;box-shadow:0 1px 2px rgba(0,0,0,.25)}
.switch input:checked + .slider{background:#16a34a}
.switch input:checked + .slider:before{transform:translateX(20px)}

/* mobile */
.mobile-wrap{padding:16px;display:flex;flex-direction:column;gap:12px;height:100%;overflow:auto}
.m-top{display:grid;grid-template-columns:1fr;gap:8px}
@media (min-width: 480px){ .m-top{grid-template-columns:1fr 1fr 1fr} }
.m-item{border:1px solid var(--line);border-radius:12px;background:#fff;padding:12px}
[data-theme="dark"] .m-item{background:#313335}
.m-title{font-weight:700}
.m-row{display:flex;align-items:center;gap:10px;color:var(--muted);font-size:.95rem;margin-top:6px;flex-wrap:wrap}
</style>
</head>
<body>
  <!-- Top bar -->
  <div class="top-bar">
    <div class="top-bar-left">
      <div class="top-bar-logo">
        <img src="/static/syllabus_logo.png" alt="syllabus Logo" class="logo">
        <span class="top-bar-title">syllabus</span>
      </div>
      <div class="top-bar-search">
        <input type="text" class="search-input" placeholder="Search series..." id="searchInput">
      </div>
    </div>
    <div class="top-bar-right">
      <div class="status-indicator" id="statusIndicator">
        <div class="status-dots">
          <div class="status-dot"></div>
          <div class="status-dot"></div>
          <div class="status-dot"></div>
          <div class="status-dot"></div>
          <div class="status-dot"></div>
          <div class="status-dot"></div>
        </div>
      </div>
      <button class="settings-btn" id="settingsBtn" onclick="openSettingsModal()" title="Settings">⚙️</button>
      <div class="user-box" id="userBox" onclick="toggleUserDropdown()">
        {{ if .Authenticated }}
          {{ if .User }}{{ .User.Username }}{{ else }}account{{ end }}
        {{ else }}
          guest
        {{ end }}
        <div class="user-dropdown" id="userDropdown">
          {{ if .Authenticated }}
            {{ if .User.IsAdmin }}
              <button class="user-dropdown-item" onclick="openUsersModal()">Users</button>
            {{ end }}
            <button class="user-dropdown-item" onclick="logout()">Logout</button>
          {{ else }}
            <a class="user-dropdown-item" href="/login">Login</a>
          {{ end }}
        </div>
      </div>
    </div>
  </div>

  <div class="main">
    <div class="container">

      <!-- ================= DESKTOP VIEW ================= -->
      <div id="desktopView" class="desktop-wrap">

        <!-- row: Soonest | Total series | With upcoming -->
        <div class="desktop-top-row">
          <div class="tile" title="Next book release across all series">
            <div class="k">Soonest release</div>
            <div class="v" id="soonestTop">—</div>
          </div>
          <div class="tile" title="Total number of series being tracked">
            <div class="k">Total series</div>
            <div class="v" id="tSeries">0</div>
          </div>
          <div class="tile" title="Number of series that have upcoming release dates">
            <div class="k">Series with upcoming dates</div>
            <div class="v" id="tUpcoming">0</div>
          </div>
        </div>

        <!-- Two-column body -->
        <div class="desktop-body">
          <!-- Left: compact checkbox filters -->
          <div class="card panel">
            <div class="card-h">Filters</div>
            <div class="card-c panel-scroll" style="padding-top:12px">
              <div class="checkrow"><input type="checkbox" id="fAudNext"><label for="fAudNext">Audible has <b>Next</b></label></div>
              <div class="checkrow"><input type="checkbox" id="fAmzNext"><label for="fAmzNext">Amazon has <b>Next</b></label></div>
              <div class="checkrow"><input type="checkbox" id="fAnyUpcoming"><label for="fAnyUpcoming">Any upcoming date</label></div>
              <div class="checkrow"><input type="checkbox" id="fNoNext"><label for="fNoNext">No upcoming date</label></div>
              <hr style="border:0;border-top:1px solid var(--line);margin:12px 0">
              <button id="clearFilters" style="width:100%;padding:.6rem .8rem;border:1px solid var(--line);border-radius:8px;background:#fff;cursor:pointer">Clear filters</button>
            </div>
          </div>

          <!-- Right: List -->
          <div class="card panel">
            <div class="panel-scroll" style="padding:0;margin:0">
              <table class="table zebra" id="seriesTable">
                <thead>
                  <tr>
                    <th class="sortable" data-sort="title" onclick="sortTable('title')">Series</th>
                    <th class="sortable" data-sort="audible" style="color:var(--aud)" onclick="sortTable('audible')">Audible</th>
                    <th class="sortable" data-sort="aud-latest" style="color:var(--aud)" onclick="sortTable('aud-latest')">Latest</th>
                    <th class="sortable" data-sort="aud-next" style="color:var(--aud)" onclick="sortTable('aud-next')">Next</th>
                    <th class="sortable" data-sort="amazon" style="color:var(--amz)" onclick="sortTable('amazon')">Amazon</th>
                    <th class="sortable" data-sort="amz-latest" style="color:var(--amz)" onclick="sortTable('amz-latest')">Latest</th>
                    <th class="sortable" data-sort="amz-next" style="color:var(--amz)" onclick="sortTable('amz-next')">Next</th>
                  </tr>
                </thead>
                <tbody id="seriesTbody">
                  {{ range .Rows }}
                  <tr class="row"
                      data-title="{{ .Title }}"
                      data-aud-count="{{ .AudibleCount }}"
                      data-amz-count="{{ .AmazonCount }}"
                      data-aud-latest="{{ .AudibleLatest }}"
                      data-amz-latest="{{ .AmazonLatest }}"
                      data-aud-next="{{ .AudibleNext }}"
                      data-amz-next="{{ .AmazonNext }}">
                    <td class="text">
                      <div style="font-weight:700;color:var(--text)">{{ .Title }}</div>
                    </td>
                    <td style="text-align:center;padding:12px 8px">
                      <div style="display:inline-flex;align-items:center;gap:6px">
                        {{ if .AudibleURL }}<a class="linkpill link-aud" href="{{ .AudibleURL }}" target="_blank" rel="noopener" title="View on Audible"><span class="icon-headphones"></span></a>{{ end }}
                        <span style="color:var(--aud);font-weight:600">{{ .AudibleCount }}</span>
                      </div>
                    </td>
                    <td><span class="latest" data-latest-pill-aud>{{ if .AudibleLatest }}{{ .AudibleLatest }}{{ else }}—{{ end }}</span></td>
                    <td><span class="next next-none" data-next-pill-aud><center>-</center></span></td>
                    <td style="text-align:center;padding:12px 8px">
                      <div style="display:inline-flex;align-items:center;gap:6px">
                        {{ if .AmazonURL }}<a class="linkpill link-amz" href="{{ .AmazonURL }}" target="_blank" rel="noopener" title="View on Amazon"><span class="icon-book"></span></a>{{ end }}
                        <span style="color:var(--amz);font-weight:600">{{ .AmazonCount }}</span>
                      </div>
                    </td>
                    <td><span class="latest" data-latest-pill-amz>{{ if .AmazonLatest }}{{ .AmazonLatest }}{{ else }}—{{ end }}</span></td>
                    <td><span class="next next-none" data-next-pill-amz><center>-</center></span></td>
                  </tr>
                  {{ end }}
                </tbody>
              </table>
              <div style="height:24px"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= MOBILE VIEW ================== -->
      <div id="mobileView" class="mobile-wrap">
        <div class="m-top">
          <div class="tile" title="Next book release across all series"><div class="k">Soonest</div><div class="v" id="soonestMobile">—</div></div>
          <div class="tile" title="Total number of series being tracked"><div class="k">Total</div><div class="v" id="tSeriesM">0</div></div>
          <div class="tile" title="Number of series that have upcoming release dates"><div class="k">With upcoming</div><div class="v" id="tUpcomingM">0</div></div>
        </div>

        {{ range .Rows }}
        <div class="m-item"
             data-title="{{ .Title }}"
             data-aud-count="{{ .AudibleCount }}"
             data-amz-count="{{ .AmazonCount }}"
             data-aud-latest="{{ .AudibleLatest }}"
             data-amz-latest="{{ .AmazonLatest }}"
             data-aud-next="{{ .AudibleNext }}"
             data-amz-next="{{ .AmazonNext }}">
          <div class="m-title">{{ .Title }}</div>
          <div class="m-row"><span class="icon-headphones" style="color:var(--aud)"></span>{{ .AudibleCount }} Latest <span class="latest" data-latest-pill-aud>{{ if .AudibleLatest }}{{ .AudibleLatest }}{{ else }}—{{ end }}</span></div>
          <div class="m-row"><span class="icon-book" style="color:var(--amz)"></span>{{ .AmazonCount }} Latest <span class="latest" data-latest-pill-amz>{{ if .AmazonLatest }}{{ .AmazonLatest }}{{ else }}—{{ end }}</span></div>
          <div class="m-row">Next (Au): <span class="next next-none" data-next-pill-aud><center>-</center></span></div>
          <div class="m-row">Next (Am): <span class="next next-none" data-next-pill-amz><center>-</center></span></div>
          <div class="m-row" style="gap:6px">
            {{ if .AudibleURL }}<a class="linkpill link-aud" href="{{ .AudibleURL }}" target="_blank" rel="noopener" title="View on Audible"><span class="icon-headphones"></span>{{ .AudibleCount }}</a>{{ end }}
            {{ if .AmazonURL }}<a class="linkpill link-amz" href="{{ .AmazonURL }}" target="_blank" rel="noopener" title="View on Amazon"><span class="icon-book"></span>{{ .AmazonCount }}</a>{{ end }}
          </div>
        </div>
        {{ end }}
      </div>

    </div>
  </div>

  <!-- Settings Modal -->
  <div id="settingsModal" class="modal-overlay" role="dialog" aria-modal="true" aria-labelledby="settingsTitle">
    <div class="modal-panel">
      <div class="modal-head">
        <div id="settingsTitle">Settings</div>
        <button id="settingsClose" class="settings-btn" aria-label="Close settings">✕</button>
      </div>
      <div class="modal-body">
        <div class="modal-row">
          <div>
            <div style="font-weight:600">Theme</div>
            <div style="color:var(--muted);font-size:.9rem">Switch between light and dark mode.</div>
          </div>
          <div style="display:flex;gap:8px;align-items:center">
            <label class="switch" aria-label="Toggle dark mode">
              <input type="checkbox" id="toggleTheme">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="modal-row">
          <div>
            <div style="font-weight:600">Show days remaining</div>
            <div style="color:var(--muted);font-size:.9rem">When on, Next shows whole-number days (e.g., "42d") and Soonest shows "in X days".</div>
          </div>
          <label class="switch" aria-label="Toggle days remaining">
            <input type="checkbox" id="toggleDays">
            <span class="slider"></span>
          </label>
        </div>
        <div class="modal-row">
          <div>
            <div style="font-weight:600">Force Scrape</div>
            <div style="color:var(--muted);font-size:.9rem">Enable mannual scrape.</div>
          </div>
          <div style="display:flex;gap:8px;align-items:center">
            <button id="forceScrapeBtn" style="padding:6px 12px;background:#16a34a;color:white;border:none;border-radius:6px;cursor:pointer;font-size:13px;font-weight:500">Force Scrape</button>
          </div>
        </div>
        <div class="modal-row">
          <div>
            <div style="font-weight:600">iCal Export</div>
            <div style="color:var(--muted);font-size:.9rem">Subscribe to upcoming release dates in your calendar app.</div>
          </div>
          <div style="display:flex;gap:8px">
            <button id="icalExportBtn" style="padding:6px 12px;background:var(--aud);color:white;border:none;border-radius:6px;cursor:pointer;font-size:13px;font-weight:500">Download</button>
            <button id="icalSubscribeBtn" style="padding:6px 12px;background:var(--amz);color:white;border:none;border-radius:6px;cursor:pointer;font-size:13px;font-weight:500">Subscribe</button>
          </div>
        </div>
        <div style="padding-top:16px;border-top:1px solid var(--line);display:flex;align-items:center;justify-content:space-between">
          <a href="https://github.com/michaeldvinci/syllabus" target="_blank" rel="noopener" style="color:var(--muted);text-decoration:none;display:flex;align-items:center;gap:4px;font-size:12px" title="View source code on GitHub">
            <svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"/>
            </svg>
            GitHub
          </a>
          <div style="font-style:italic;color:var(--muted);font-size:12px">Last scrape started: {{ .LastScrape }}</div>
        </div>
      </div>
    </div>
  </div>

  <!-- User Management Modal -->
  <div id="usersModal" class="modal-overlay" role="dialog" aria-modal="true" aria-labelledby="usersTitle">
    <div class="modal-panel" style="max-width:600px">
      <div class="modal-head">
        <div id="usersTitle">User Management</div>
        <button id="usersClose" class="settings-btn" aria-label="Close user management">✕</button>
      </div>
      <div class="modal-body">
        <div class="modal-row">
          <div>
            <div style="font-weight:600">Create New User</div>
            <div style="color:var(--muted);font-size:.9rem">Add a new user to the system.</div>
          </div>
          <button id="createUserBtn" style="padding:6px 12px;background:var(--aud);color:white;border:none;border-radius:6px;cursor:pointer;font-size:13px;font-weight:500">Create User</button>
        </div>
        <div style="padding-top:16px;border-top:1px solid var(--line)">
          <div style="font-weight:600;margin-bottom:12px">Existing Users</div>
          <div id="usersList" style="max-height:300px;overflow-y:auto">
            <div style="text-align:center;color:var(--muted);padding:20px">Loading users...</div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Create User Modal -->
  <div id="createUserModal" class="modal-overlay" role="dialog" aria-modal="true" aria-labelledby="createUserTitle">
    <div class="modal-panel">
      <div class="modal-head">
        <div id="createUserTitle">Create New User</div>
        <button id="createUserClose" class="settings-btn" aria-label="Close create user">✕</button>
      </div>
      <div class="modal-body">
        <div style="margin-bottom:16px">
          <label style="display:block;margin-bottom:4px;font-weight:600">Username</label>
          <input type="text" id="newUsername" style="width:100%;padding:8px;border:1px solid var(--line);border-radius:4px;background:var(--bg);color:var(--text)" placeholder="Enter username">
        </div>
        <div style="margin-bottom:16px">
          <label style="display:block;margin-bottom:4px;font-weight:600">Password</label>
          <input type="password" id="newPassword" style="width:100%;padding:8px;border:1px solid var(--line);border-radius:4px;background:var(--bg);color:var(--text)" placeholder="Enter password">
        </div>
        <div style="margin-bottom:16px">
          <label style="display:block;margin-bottom:4px;font-weight:600">Role</label>
          <select id="newUserRole" style="width:100%;padding:8px;border:1px solid var(--line);border-radius:4px;background:var(--bg);color:var(--text)">
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end">
          <button id="cancelCreateUser" style="padding:8px 16px;background:var(--bg);color:var(--text);border:1px solid var(--line);border-radius:6px;cursor:pointer">Cancel</button>
          <button id="confirmCreateUser" style="padding:8px 16px;background:var(--aud);color:white;border:none;border-radius:6px;cursor:pointer">Create</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Reset Password Modal -->
  <div id="resetPasswordModal" class="modal-overlay" role="dialog" aria-modal="true" aria-labelledby="resetPasswordTitle">
    <div class="modal-panel">
      <div class="modal-head">
        <div id="resetPasswordTitle">Reset Password</div>
        <button id="resetPasswordClose" class="settings-btn" aria-label="Close reset password">✕</button>
      </div>
      <div class="modal-body">
        <div style="margin-bottom:16px">
          <div style="font-weight:600;margin-bottom:8px">User: <span id="resetUsernameDisplay"></span></div>
        </div>
        <div style="margin-bottom:16px">
          <label style="display:block;margin-bottom:4px;font-weight:600">New Password</label>
          <input type="password" id="resetNewPassword" style="width:100%;padding:8px;border:1px solid var(--line);border-radius:4px;background:var(--bg);color:var(--text)" placeholder="Enter new password">
        </div>
        <div style="margin-bottom:16px">
          <label style="display:block;margin-bottom:4px;font-weight:600">Confirm New Password</label>
          <input type="password" id="resetConfirmPassword" style="width:100%;padding:8px;border:1px solid var(--line);border-radius:4px;background:var(--bg);color:var(--text)" placeholder="Confirm new password">
        </div>
        <div style="display:flex;gap:8px;justify-content:flex-end">
          <button id="cancelResetPassword" style="padding:8px 16px;background:var(--bg);color:var(--text);border:1px solid var(--line);border-radius:6px;cursor:pointer">Cancel</button>
          <button id="confirmResetPassword" style="padding:8px 16px;background:#f59e0b;color:white;border:none;border-radius:6px;cursor:pointer">Reset Password</button>
        </div>
      </div>
    </div>
  </div>

<script>
/* ── robust date helpers ───────────────────────── */
const parseAnyDate = (val) => {
  if (val == null) return null;
  const s = (''+val).trim();
  if (!s || s.toLowerCase()==='none' || s==='—' || s.toLowerCase()==='-' || s.toLowerCase()==='n/a') return null;
  const m = s.match(/\\b\\d{4}-\\d{2}-\\d{2}\\b/);
  if (m) return new Date(m[0] + 'T00:00:00');
  const d = new Date(s);
  return isNaN(d) ? null : d;
};

const daysUntil = (val) => {
  const d = parseAnyDate(val);
  if (!d) return null;
  const now = new Date(), today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const one = 1000*60*60*24;
  return Math.ceil((d - today)/one);
};

const fmt = (val)=> {
  const d = parseAnyDate(val);
  return d ? d.toLocaleDateString(undefined,{year:'numeric',month:'short',day:'2-digit'}) : '—';
};

/* ── settings state ─────────────────────────────── */
let SHOW_DAYS = false;
try { SHOW_DAYS = localStorage.getItem('syll_show_days') === '1'; } catch(e){}

/* ── UI helpers ─────────────────────────────────── */
function setPill(pill, date){
  if(!pill) return;
  const d = parseAnyDate(date);
  if(!d){
    pill.innerHTML = '<center>-</center>';
    pill.className = 'next next-none';
    return;
  }
  const left = daysUntil(d);
  pill.title = fmt(d); // tooltip always shows absolute date
  if (SHOW_DAYS && left != null) {
    pill.textContent = left>0 ? (left+'d') : 'soon';
  } else {
    pill.textContent = fmt(d);
  }
  // Color based on provider (Audible vs Amazon)
  if (left == null) {
    pill.className = 'next next-none';
  } else {
    const isAudible = pill.hasAttribute('data-next-pill-aud');
    pill.className = 'next ' + (isAudible ? 'next-aud' : 'next-amz');
  }
}

function setLatestPill(pill, date){
  if(!pill) return;
  const d = parseAnyDate(date);
  if(!d){
    pill.innerHTML = '<center>—</center>';
    return;
  }
  pill.title = fmt(d); // tooltip shows absolute date
  if (SHOW_DAYS) {
    const daysAgo = -daysUntil(d); // negative because it's in the past
    pill.textContent = daysAgo > 0 ? (daysAgo + 'd ago') : 'today';
  } else {
    pill.textContent = fmt(d);
  }
}

function computeTilesAndDecorate(){
  const dRows = Array.from(document.querySelectorAll('#seriesTbody tr'));
  const mRows = Array.from(document.querySelectorAll('#mobileView .m-item'));
  const all = dRows.length ? dRows : mRows;

  let total = all.length, upcoming = 0, soonest = null;

  const consider = (date) => {
    const d = parseAnyDate(date);
    if(!d) return;
    if(!soonest || d < soonest) soonest = d;
  };

  const decorate = (el) => {
    const audNext = el.dataset.audNext;
    const amzNext = el.dataset.amzNext;
    const audLatest = el.dataset.audLatest;
    const amzLatest = el.dataset.amzLatest;
    
    setPill(el.querySelector('[data-next-pill-aud]'), audNext);
    setPill(el.querySelector('[data-next-pill-amz]'), amzNext);
    setLatestPill(el.querySelector('[data-latest-pill-aud]'), audLatest);
    setLatestPill(el.querySelector('[data-latest-pill-amz]'), amzLatest);

    if(parseAnyDate(audNext) || parseAnyDate(amzNext)) upcoming++;
    consider(audNext); consider(amzNext);
  };
  all.forEach(decorate);

  const tSeries = document.querySelector('#tSeries'); if(tSeries) tSeries.textContent = total;
  const tUpcoming = document.querySelector('#tUpcoming'); if(tUpcoming) tUpcoming.textContent = upcoming;
  const soonTop = document.querySelector('#soonestTop');
  if(soonTop){
    if(soonest){
      const left = daysUntil(soonest);
      soonTop.textContent = SHOW_DAYS && left!=null ? ('in '+left+' days') : fmt(soonest);
    }else{
      soonTop.textContent = '—';
    }
  }

  const tSeriesM = document.querySelector('#tSeriesM'); if(tSeriesM) tSeriesM.textContent = total;
  const tUpcomingM = document.querySelector('#tUpcomingM'); if(tUpcomingM) tUpcomingM.textContent = upcoming;
  const soonM = document.querySelector('#soonestMobile');
  if(soonM){
    if(soonest){
      const left = daysUntil(soonest);
      soonM.textContent = SHOW_DAYS && left!=null ? ('in '+left+' days') : fmt(soonest);
    } else {
      soonM.textContent = '—';
    }
  }
}

/* ── filtering ──────────────────────────────────── */
function wireFilters(){
  function applyFilters(){
    const fAud  = document.querySelector('#fAudNext')?.checked;
    const fAmz  = document.querySelector('#fAmzNext')?.checked;
    const fAny  = document.querySelector('#fAnyUpcoming')?.checked;
    const fNone = document.querySelector('#fNoNext')?.checked;

    [Array.from(document.querySelectorAll('#seriesTbody tr')), Array.from(document.querySelectorAll('#mobileView .m-item'))].forEach(rows=>{
      rows.forEach(el=>{
        const hasAud = !!parseAnyDate(el.dataset.audNext);
        const hasAmz = !!parseAnyDate(el.dataset.amzNext);
        const any = hasAud || hasAmz;
        let show = true;
        if(fAud && !hasAud) show = false;
        if(fAmz && !hasAmz) show = false;
        if(fAny && !any) show = false;
        if(fNone && any) show = false;
        el.style.display = show ? '' : 'none';
      });
    });
  }

  ['#fAudNext','#fAmzNext','#fAnyUpcoming','#fNoNext'].forEach(id=>{
    const cb = document.querySelector(id); if(cb) cb.addEventListener('change', applyFilters);
  });
  const clearBtn = document.querySelector('#clearFilters');
  if(clearBtn){
    clearBtn.addEventListener('click', ()=>{
      ['#fAudNext','#fAmzNext','#fAnyUpcoming','#fNoNext'].forEach(id=>{ const cb=document.querySelector(id); if(cb) cb.checked=false; });
      applyFilters();
    });
  }
}

/* ── search ─────────────────────────────────────── */
function wireSearch(){
  const input = document.querySelector('#searchInput');
  if(!input) return;
  const dRows = Array.from(document.querySelectorAll('#seriesTbody tr'));
  const mRows = Array.from(document.querySelectorAll('#mobileView .m-item'));
  input.addEventListener('input', e=>{
    const term = e.target.value.trim().toLowerCase();
    [dRows, mRows].forEach(rows=>{
      rows.forEach(el=>{
        const ok = !term || (el.dataset.title||'').toLowerCase().includes(term);
        el.style.display = ok ? '' : 'none';
      });
    });
  });
}

/* ── status indicator ──────────────────────────── */
let ACTIVE_TASKS = 0;
let POLLING_ACTIVE = false;
let POLL_TIMEOUT = null;

function showStatusIndicator(){
  ACTIVE_TASKS++;
  const indicator = document.getElementById('statusIndicator');
  if(indicator) indicator.style.display = 'block';
}

function hideStatusIndicator(){
  ACTIVE_TASKS = Math.max(0, ACTIVE_TASKS - 1);
  if(ACTIVE_TASKS === 0){
    const indicator = document.getElementById('statusIndicator');
    if(indicator) indicator.style.display = 'none';
    stopPolling();
  }
}

function stopPolling(){
  POLLING_ACTIVE = false;
  if(POLL_TIMEOUT){
    clearTimeout(POLL_TIMEOUT);
    POLL_TIMEOUT = null;
  }
}

/* ── theme management ──────────────────────────── */
let CURRENT_THEME = 'light';
try { CURRENT_THEME = localStorage.getItem('syll_theme') || 'light'; } catch(e){}

function applyTheme(theme){
  if(theme === 'dark'){
    document.documentElement.setAttribute('data-theme', 'dark');
  } else {
    document.documentElement.removeAttribute('data-theme');
  }
  CURRENT_THEME = theme;
  try { localStorage.setItem('syll_theme', theme); } catch(e){}
}

/* ── settings modal ─────────────────────────────── */
function openSettingsModal(){
  const overlay = document.getElementById('settingsModal');
  const toggleDays = document.getElementById('toggleDays');
  const toggleTheme = document.getElementById('toggleTheme');
  if(!overlay) return;
  
  if(toggleDays) toggleDays.checked = !!SHOW_DAYS;
  if(toggleTheme) toggleTheme.checked = CURRENT_THEME === 'dark';
  
  overlay.style.display = 'flex';
  const closer = document.getElementById('settingsClose');
  const onKey = (ev)=>{ if(ev.key==='Escape'){ close(); } };
  function close(){
    overlay.style.display = 'none';
    document.removeEventListener('keydown', onKey);
  }
  closer?.addEventListener('click', close, { once:true });
  overlay.addEventListener('click', (e)=>{ if(e.target===overlay) close(); }, { once:true });
  document.addEventListener('keydown', onKey);
}
window.openSettingsModal = openSettingsModal;

function wireSettings(){
  const toggleDays = document.getElementById('toggleDays');
  const toggleTheme = document.getElementById('toggleTheme');
  const forceScrapeBtn = document.getElementById('forceScrapeBtn');
  
  if(toggleDays){
    toggleDays.addEventListener('change', ()=>{
      SHOW_DAYS = !!toggleDays.checked;
      try { localStorage.setItem('syll_show_days', SHOW_DAYS ? '1' : '0'); } catch(e){}
      computeTilesAndDecorate(); // re-render pill texts and tiles
    });
  }
  
  if(toggleTheme){
    toggleTheme.addEventListener('change', ()=>{
      const newTheme = toggleTheme.checked ? 'dark' : 'light';
      applyTheme(newTheme);
    });
  }
  
  if(forceScrapeBtn){
    forceScrapeBtn.addEventListener('click', ()=>{
      forceScrapeBtn.disabled = true;
      forceScrapeBtn.textContent = 'Scraping...';
      
      fetch('/refresh', {method: 'POST'})
        .then(response => {
          if(response.ok){
            // Start polling to monitor the background scrape
            startPolling();
            
            // Re-enable button after a delay
            setTimeout(()=>{
              forceScrapeBtn.disabled = false;
              forceScrapeBtn.textContent = 'Force Scrape';
            }, 2000);
          } else {
            throw new Error('Scrape failed');
          }
        })
        .catch(err => {
          console.error('Force scrape failed:', err);
          forceScrapeBtn.textContent = 'Failed - Try Again';
          setTimeout(()=>{
            forceScrapeBtn.disabled = false;
            forceScrapeBtn.textContent = 'Force Scrape';
          }, 3000);
        });
    });
  }
}

/* ── user dropdown ─────────────────────────────── */
function toggleUserDropdown(){
  const dropdown = document.getElementById('userDropdown');
  if(!dropdown) return;
  const isVisible = dropdown.style.display === 'block';
  dropdown.style.display = isVisible ? 'none' : 'block';
  if(!isVisible){
    const onClick = (e)=>{
      if(!e.target.closest('#userBox')){
        dropdown.style.display = 'none';
        document.removeEventListener('click', onClick);
      }
    };
    setTimeout(()=>document.addEventListener('click', onClick), 0);
  }
}
window.toggleUserDropdown = toggleUserDropdown;

function logout(){
  fetch('/logout', {method: 'POST'})
    .then(()=>window.location.reload())
    .catch(err=>console.error('Logout failed:', err));
}
window.logout = logout;

/* ── user management ───────────────────────────── */
function openUsersModal(){
  const overlay = document.getElementById('usersModal');
  if(!overlay) return;
  overlay.style.display = 'flex';
  loadUsers();
  
  const closer = document.getElementById('usersClose');
  const onKey = (ev)=>{ if(ev.key==='Escape'){ closeUsersModal(); } };
  function closeUsersModal(){
    overlay.style.display = 'none';
    document.removeEventListener('keydown', onKey);
  }
  closer?.addEventListener('click', closeUsersModal, { once:true });
  overlay.addEventListener('click', (e)=>{ if(e.target===overlay) closeUsersModal(); }, { once:true });
  document.addEventListener('keydown', onKey);
}
window.openUsersModal = openUsersModal;

function loadUsers(){
  const usersList = document.getElementById('usersList');
  if(!usersList) return;
  
  usersList.innerHTML = '<div style="text-align:center;color:var(--muted);padding:20px">Loading users...</div>';
  
  fetch('/api/users')
    .then(response => response.json())
    .then(data => {
      if(data.users && data.users.length > 0){
        usersList.innerHTML = data.users.map(user => 
          '<div class="user-list-item">' +
            '<div class="user-info">' +
              '<div class="user-name">' + user.username + '</div>' +
              '<div class="user-role">' + user.role + '</div>' +
            '</div>' +
            '<div class="user-actions">' +
              '<button class="reset-password-btn" onclick="openResetPasswordModal(\'' + user.username + '\')">' +
                'Reset Password' +
              '</button>' +
              '<button class="delete-user-btn" onclick="deleteUser(\'' + user.username + '\')" ' + 
              (user.username === 'admin' ? 'disabled title="Cannot delete admin user"' : '') + '>' +
                'Delete' +
              '</button>' +
            '</div>' +
          '</div>'
        ).join('');
      } else {
        usersList.innerHTML = '<div style="text-align:center;color:var(--muted);padding:20px">No users found.</div>';
      }
    })
    .catch(err => {
      console.error('Failed to load users:', err);
      usersList.innerHTML = '<div style="text-align:center;color:#dc2626;padding:20px">Failed to load users.</div>';
    });
}

function deleteUser(username){
  if(!confirm('Are you sure you want to delete user "' + username + '"?')) return;
  
  fetch('/api/users/delete', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({username})
  })
  .then(response => response.json())
  .then(data => {
    if(data.success){
      loadUsers(); // Reload the list
    } else {
      alert('Failed to delete user: ' + (data.message || 'Unknown error'));
    }
  })
  .catch(err => {
    console.error('Delete user failed:', err);
    alert('Failed to delete user.');
  });
}
window.deleteUser = deleteUser;

function openCreateUserModal(){
  const overlay = document.getElementById('createUserModal');
  if(!overlay) return;
  
  // Clear form
  document.getElementById('newUsername').value = '';
  document.getElementById('newPassword').value = '';
  document.getElementById('newUserRole').value = 'user';
  
  overlay.style.display = 'flex';
  
  const closer = document.getElementById('createUserClose');
  const cancel = document.getElementById('cancelCreateUser');
  const confirm = document.getElementById('confirmCreateUser');
  
  const onKey = (ev)=>{ if(ev.key==='Escape'){ closeCreateUserModal(); } };
  function closeCreateUserModal(){
    overlay.style.display = 'none';
    document.removeEventListener('keydown', onKey);
  }
  
  closer?.addEventListener('click', closeCreateUserModal, { once:true });
  cancel?.addEventListener('click', closeCreateUserModal, { once:true });
  overlay.addEventListener('click', (e)=>{ if(e.target===overlay) closeCreateUserModal(); }, { once:true });
  document.addEventListener('keydown', onKey);
  
  confirm?.addEventListener('click', ()=>{
    const username = document.getElementById('newUsername').value.trim();
    const password = document.getElementById('newPassword').value;
    const role = document.getElementById('newUserRole').value;
    
    if(!username || !password){
      alert('Username and password are required.');
      return;
    }
    
    fetch('/api/users/create', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username, password, role})
    })
    .then(response => response.json())
    .then(data => {
      if(data.success){
        closeCreateUserModal();
        loadUsers(); // Reload the list
      } else {
        alert('Failed to create user: ' + (data.message || 'Unknown error'));
      }
    })
    .catch(err => {
      console.error('Create user failed:', err);
      alert('Failed to create user.');
    });
  }, { once:true });
}

function openResetPasswordModal(username){
  const overlay = document.getElementById('resetPasswordModal');
  const usernameDisplay = document.getElementById('resetUsernameDisplay');
  if(!overlay || !usernameDisplay) return;
  
  // Set username and clear form
  usernameDisplay.textContent = username;
  document.getElementById('resetNewPassword').value = '';
  document.getElementById('resetConfirmPassword').value = '';
  
  overlay.style.display = 'flex';
  
  const closer = document.getElementById('resetPasswordClose');
  const cancel = document.getElementById('cancelResetPassword');
  const confirm = document.getElementById('confirmResetPassword');
  
  const onKey = (ev)=>{ if(ev.key==='Escape'){ closeResetPasswordModal(); } };
  function closeResetPasswordModal(){
    overlay.style.display = 'none';
    document.removeEventListener('keydown', onKey);
  }
  
  closer?.addEventListener('click', closeResetPasswordModal, { once:true });
  cancel?.addEventListener('click', closeResetPasswordModal, { once:true });
  overlay.addEventListener('click', (e)=>{ if(e.target===overlay) closeResetPasswordModal(); }, { once:true });
  document.addEventListener('keydown', onKey);
  
  confirm?.addEventListener('click', ()=>{
    const newPassword = document.getElementById('resetNewPassword').value;
    const confirmPassword = document.getElementById('resetConfirmPassword').value;
    
    if(!newPassword || !confirmPassword){
      alert('Both password fields are required.');
      return;
    }
    
    if(newPassword !== confirmPassword){
      alert('Passwords do not match.');
      return;
    }
    
    if(newPassword.length < 4){
      alert('Password must be at least 4 characters long.');
      return;
    }
    
    fetch('/api/users/reset-password', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username: username, newPassword: newPassword})
    })
    .then(response => response.json())
    .then(data => {
      if(data.success){
        alert('Password reset successfully for user: ' + username);
        closeResetPasswordModal();
      } else {
        alert('Failed to reset password: ' + (data.message || 'Unknown error'));
      }
    })
    .catch(err => {
      console.error('Reset password failed:', err);
      alert('Failed to reset password.');
    });
  }, { once:true });
}
window.openResetPasswordModal = openResetPasswordModal;

function wireUserManagement(){
  const createBtn = document.getElementById('createUserBtn');
  if(createBtn){
    createBtn.addEventListener('click', openCreateUserModal);
  }
}

/* ── table sorting ─────────────────────────────── */
let CURRENT_SORT = {column: null, direction: 'asc'};

function sortTable(column){
  const table = document.getElementById('seriesTable');
  const tbody = document.getElementById('seriesTbody');
  if(!table || !tbody) return;
  
  // Toggle sort direction if same column, otherwise default to asc
  if(CURRENT_SORT.column === column){
    CURRENT_SORT.direction = CURRENT_SORT.direction === 'asc' ? 'desc' : 'asc';
  } else {
    CURRENT_SORT.direction = 'asc';
  }
  CURRENT_SORT.column = column;
  
  // Update header indicators
  const headers = table.querySelectorAll('th.sortable');
  headers.forEach(th => {
    th.classList.remove('sort-asc', 'sort-desc');
    if(th.dataset.sort === column){
      th.classList.add(CURRENT_SORT.direction === 'asc' ? 'sort-asc' : 'sort-desc');
    }
  });
  
  // Get all rows and sort them
  const rows = Array.from(tbody.querySelectorAll('tr'));
  rows.sort((a, b) => {
    let aVal, bVal;
    
    switch(column){
      case 'title':
        aVal = a.dataset.title || '';
        bVal = b.dataset.title || '';
        return CURRENT_SORT.direction === 'asc' ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
        
      case 'audible':
        aVal = parseInt(a.dataset.audCount) || 0;
        bVal = parseInt(b.dataset.audCount) || 0;
        return CURRENT_SORT.direction === 'asc' ? aVal - bVal : bVal - aVal;
        
      case 'amazon':
        aVal = parseInt(a.dataset.amzCount) || 0;
        bVal = parseInt(b.dataset.amzCount) || 0;
        return CURRENT_SORT.direction === 'asc' ? aVal - bVal : bVal - aVal;
        
      case 'aud-latest':
        aVal = parseAnyDate(a.dataset.audLatest);
        bVal = parseAnyDate(b.dataset.audLatest);
        if(!aVal && !bVal) return 0;
        if(!aVal) return CURRENT_SORT.direction === 'asc' ? 1 : -1;
        if(!bVal) return CURRENT_SORT.direction === 'asc' ? -1 : 1;
        return CURRENT_SORT.direction === 'asc' ? aVal - bVal : bVal - aVal;
        
      case 'aud-next':
        aVal = parseAnyDate(a.dataset.audNext);
        bVal = parseAnyDate(b.dataset.audNext);
        if(!aVal && !bVal) return 0;
        if(!aVal) return CURRENT_SORT.direction === 'asc' ? 1 : -1;
        if(!bVal) return CURRENT_SORT.direction === 'asc' ? -1 : 1;
        return CURRENT_SORT.direction === 'asc' ? aVal - bVal : bVal - aVal;
        
      case 'amz-latest':
        aVal = parseAnyDate(a.dataset.amzLatest);
        bVal = parseAnyDate(b.dataset.amzLatest);
        if(!aVal && !bVal) return 0;
        if(!aVal) return CURRENT_SORT.direction === 'asc' ? 1 : -1;
        if(!bVal) return CURRENT_SORT.direction === 'asc' ? -1 : 1;
        return CURRENT_SORT.direction === 'asc' ? aVal - bVal : bVal - aVal;
        
      case 'amz-next':
        aVal = parseAnyDate(a.dataset.amzNext);
        bVal = parseAnyDate(b.dataset.amzNext);
        if(!aVal && !bVal) return 0;
        if(!aVal) return CURRENT_SORT.direction === 'asc' ? 1 : -1;
        if(!bVal) return CURRENT_SORT.direction === 'asc' ? -1 : 1;
        return CURRENT_SORT.direction === 'asc' ? aVal - bVal : bVal - aVal;
        
      default:
        return 0;
    }
  });
  
  // Re-append sorted rows
  rows.forEach(row => tbody.appendChild(row));
}
window.sortTable = sortTable;

/* ── ical export ───────────────────────────────── */
function wireIcalExport(){
  const downloadBtn = document.getElementById('icalExportBtn');
  const subscribeBtn = document.getElementById('icalSubscribeBtn');
  
  if(downloadBtn){
    downloadBtn.addEventListener('click', ()=>{
      window.open('/calendar.ics', '_blank');
    });
  }
  
  if(subscribeBtn){
    subscribeBtn.addEventListener('click', ()=>{
      const currentUrl = window.location.origin;
      const icalUrl = currentUrl + '/calendar.ics';
      
      // Try to copy to clipboard first
      if(navigator.clipboard){
        navigator.clipboard.writeText(icalUrl).then(()=>{
          alert('iCal subscription URL copied to clipboard:\n\n' + icalUrl + '\n\nPaste this into your calendar app to subscribe.');
        }).catch(()=>{
          // Fallback: show URL in alert
          alert('iCal subscription URL:\n\n' + icalUrl + '\n\nCopy this URL and paste it into your calendar app to subscribe.');
        });
      } else {
        // Fallback: show URL in alert
        alert('iCal subscription URL:\n\n' + icalUrl + '\n\nCopy this URL and paste it into your calendar app to subscribe.');
      }
    });
  }
}

/* ── background task monitoring ─────────────────── */
function startPolling(){
  if(POLLING_ACTIVE) return; // Already polling
  POLLING_ACTIVE = true;
  checkBackgroundTasks();
}

function checkBackgroundTasks(){
  if(!POLLING_ACTIVE) return; // Stop if polling was disabled
  
  fetch('/api/scrape-status')
    .then(response => response.json())
    .then(data => {
      if(data.activeJobs > 0){
        showStatusIndicator();
        if(POLLING_ACTIVE){ // Only schedule next check if still active
          POLL_TIMEOUT = setTimeout(checkBackgroundTasks, 2000);
        }
      } else {
        hideStatusIndicator(); // This will call stopPolling()
      }
    })
    .catch(() => {
      // If API fails, assume no active tasks and stop polling
      console.log('Scrape status API failed, stopping polling');
      hideStatusIndicator();
    });
}

/* ── boot ───────────────────────────────────────── */
function initDashboard(){
  // Apply saved theme on load
  applyTheme(CURRENT_THEME);
  
  wireFilters();
  wireSearch();
  wireSettings();
  wireIcalExport();
  wireUserManagement();
  computeTilesAndDecorate();
  
  // Start monitoring background tasks (only if there might be active tasks)
  startPolling();
}

document.addEventListener('DOMContentLoaded', initDashboard);

// Clean up polling when page is unloaded
window.addEventListener('beforeunload', ()=>{
  stopPolling();
});
</script>
</body>
</html>
`

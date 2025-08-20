package handlers

const IndexHTML = `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Audiobook / Ebook Series Tracker</title>
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
body { font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial; margin: 2rem; color: var(--text); background: var(--bg); }
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
.links { display: inline-flex; gap: .35rem; }
.pill { display: inline-flex; align-items: center; justify-content: center; font-size: .72rem; line-height: 1; padding: .28rem .45rem; border-radius: 999px; text-decoration: none; border: 1px solid rgba(0,0,0,.06); }
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
.topbar { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-bottom: .75rem; }
.settings-btn { display: inline-flex; align-items: center; justify-content: center; width: 2.25rem; height: 2.25rem; border-radius: .5rem; border: 1px solid var(--line); background: #fff; box-shadow: var(--head-shadow); cursor: pointer; font-size: 1.05rem; }
.settings-btn:hover { background: #f8fafc; }
.settings-btn:focus { outline: 2px solid #93c5fd; outline-offset: 2px; }
.settings-wrap { position: relative; }
.settings-panel { position: absolute; right: 0; top: 2.8rem; width: 320px; max-width: calc(100vw - 2rem); background: #fff; border: 1px solid var(--line); border-radius: .5rem; box-shadow: 0 10px 20px rgba(0,0,0,.08), 0 2px 6px rgba(0,0,0,.06); padding: .75rem; z-index: 10; }
.settings-panel .panel-section { padding: .5rem .25rem; }
.settings-panel .panel-heading { font-weight: 700; font-size: .85rem; color: var(--muted); margin-bottom: .25rem; text-transform: uppercase; letter-spacing: .02em; }
.settings-panel code { background: #f3f4f6; padding: .15rem .35rem; border-radius: .35rem; }
</style>
</head>
<body>
  <div class="topbar">
    <h1>Syllabus</h1>
    <div class="settings-wrap">
      <button class="settings-btn" id="settingsBtn" aria-expanded="false" aria-controls="settingsPanel" title="Settings" aria-label="Settings">⚙️</button>
      <div class="settings-panel" id="settingsPanel" hidden>
        <div class="panel-section">
          <div class="panel-heading">Generated at</div>
          <div class="panel-content"><code>{{ .Now }}</code></div>
        </div>
      </div>
    </div>
  </div>
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
  })();
  </script>
</body>
</html>
`
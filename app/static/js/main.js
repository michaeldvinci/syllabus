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
          alert('iCal subscription URL copied to clipboard:\\n\\n' + icalUrl + '\\n\\nPaste this into your calendar app to subscribe.');
        }).catch(()=>{
          // Fallback: show URL in alert
          alert('iCal subscription URL:\\n\\n' + icalUrl + '\\n\\nCopy this URL and paste it into your calendar app to subscribe.');
        });
      } else {
        // Fallback: show URL in alert
        alert('iCal subscription URL:\\n\\n' + icalUrl + '\\n\\nCopy this URL and paste it into your calendar app to subscribe.');
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
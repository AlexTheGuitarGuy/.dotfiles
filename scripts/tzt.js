#!/usr/bin/env node
const { execSync, spawnSync } = require('child_process');
const fs = require('fs');
const path = require('path');

function run(cmd, args, opts = {}) {
  const result = spawnSync(cmd, args, {
    encoding: 'utf-8',
    stdio: ['pipe', 'pipe', 'inherit'],
    ...opts,
  });
  return { ok: result.status === 0, stdout: (result.stdout || '').trim() };
}

function shell(cmd) {
  try {
    const stdout = execSync(cmd, {
      encoding: 'utf-8',
      stdio: ['pipe', 'pipe', 'inherit'],
    });
    return { ok: true, stdout: (stdout || '').trim() };
  } catch (e) {
    return { ok: false, stdout: (e.stdout || '').trim() };
  }
}

function pickDir(pattern) {
  if (pattern) {
    return shell(`zoxide query -l ${JSON.stringify(pattern)}`);
  }
  const r = shell('zoxide query -l');
  if (!r.ok || !r.stdout) return r;
  const fzf = spawnSync('fzf', {
    input: r.stdout,
    encoding: 'utf-8',
    stdio: ['pipe', 'pipe', 'inherit'],
  });
  return { ok: fzf.status === 0, stdout: (fzf.stdout || '').trim() };
}

function parseWorkspaceItems(raw) {
  if (!raw) return [];
  try {
    const data = JSON.parse(raw);
    const list =
      data?.result?.workspaces ??  // herdr CLI response format
      data?.workspaces ??          // alternative wrapper
      (Array.isArray(data) ? data : null);
    return Array.isArray(list) ? list : [];
  } catch {
    // NDJSON fallback
    const lines = raw.split('\n').filter(Boolean);
    const items = [];
    for (const line of lines) {
      try { items.push(JSON.parse(line)); } catch {}
    }
    return items;
  }
}

function findWorkspace(items, selectedName, selected) {
  return items.find((w) => {
    const id = w.workspace_id || w.id;
    if (!id) return false;
    return (
      w.label === selectedName ||
      w.name === selectedName ||
      w.title === selectedName ||
      w.displayName === selectedName ||
      w.display_name === selectedName ||
      w.handle === selectedName ||
      w.cwd === selected
    );
  });
}

const pattern = process.argv[2];
const picked = pickDir(pattern);

if (!picked.ok || !picked.stdout) process.exit(0);

const selected = picked.stdout.split('\n')[0];
const selectedName = path.basename(selected).replace(/[^a-zA-Z0-9_-]/g, '_');
const stateDir = path.join(process.env.HOME || '/home/' + process.env.USER, '.local/share/tzt');
const stateFile = path.join(stateDir, 'workspaces.json');

let wsId = null;
try {
  if (fs.existsSync(stateFile)) {
    const state = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
    if (state[selectedName]) wsId = state[selectedName];
  }
} catch {}

if (wsId) {
  const focusResult = run('herdr', ['workspace', 'focus', wsId]);
  if (focusResult.ok) {
    // All good — launch herdr if not inside it
    if (!process.env.HERDR_PANE_ID) {
      spawnSync('herdr', { stdio: 'inherit' });
    }
    process.exit(0);
  }
  // Stale ID — clear it and fall through to create
  wsId = null;
  try {
    if (fs.existsSync(stateFile)) {
      const state = JSON.parse(fs.readFileSync(stateFile, 'utf-8'));
      delete state[selectedName];
      fs.writeFileSync(stateFile, JSON.stringify(state));
    }
  } catch {}
}

// Create new workspace
{
  run('herdr', ['workspace', 'create', '--cwd', selected, '--label', selectedName, '--focus']);

  const list = run('herdr', ['workspace', 'list']);
  if (list.ok) {
    const items = parseWorkspaceItems(list.stdout);
    const match = findWorkspace(items, selectedName, selected);
    if (match) {
      const matchId = match.workspace_id || match.id;
      fs.mkdirSync(stateDir, { recursive: true });
      const existing = {};
      try {
        if (fs.existsSync(stateFile)) Object.assign(existing, JSON.parse(fs.readFileSync(stateFile, 'utf-8')));
      } catch {}
      existing[selectedName] = matchId;
      fs.writeFileSync(stateFile, JSON.stringify(existing));
    } else {
      // No match found — dump raw output for inspection
      try {
        fs.mkdirSync(stateDir, { recursive: true });
        fs.writeFileSync(path.join(stateDir, 'debug.json'),
          JSON.stringify({ name: selectedName, cwd: selected, itemCount: items.length, firstItem: items[0] || null, raw: list.stdout }, null, 2));
      } catch {}
    }
  }

  run('herdr', ['tab', 'create', '--cwd', selected, '--label', 'term-2', '--no-focus']);
  run('herdr', ['tab', 'create', '--cwd', selected, '--label', 'term-3', '--no-focus']);
  run('herdr', ['tab', 'create', '--cwd', selected, '--label', 'term-4', '--no-focus']);
}

if (!process.env.HERDR_PANE_ID) {
  spawnSync('herdr', { stdio: 'inherit' });
}

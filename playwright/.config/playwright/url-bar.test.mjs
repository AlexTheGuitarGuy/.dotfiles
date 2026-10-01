import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import vm from 'node:vm';

const ctx = { URL };
vm.runInNewContext(readFileSync(new URL('./url-bar.js', import.meta.url), 'utf8'), ctx);

test('shows origin, path and query', () => {
  assert.equal(
    ctx.displayUrl('https://finanzanalyse.integration.dvag/berater/haushalt/1/uebersicht?tab=2'),
    'https://finanzanalyse.integration.dvag/berater/haushalt/1/uebersicht?tab=2',
  );
});

test('drops the Keycloak redirect hash with the auth code', () => {
  assert.equal(
    ctx.displayUrl('https://finanzanalyse.integration.dvag/x#state=a&session_state=b&code=c'),
    'https://finanzanalyse.integration.dvag/x',
  );
});

test('removes secret query params but keeps the rest', () => {
  assert.equal(
    ctx.displayUrl('https://integration.auth.dvag/auth?client_id=web&code=c&state=s&session_state=x&access_token=t'),
    'https://integration.auth.dvag/auth?client_id=web',
  );
});

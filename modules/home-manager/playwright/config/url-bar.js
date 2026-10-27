var SECRET_PARAMS = ['code', 'state', 'session_state', 'access_token', 'id_token', 'refresh_token', 'token'];

function displayUrl(href) {
  var url = new URL(href);
  SECRET_PARAMS.forEach(function (name) {
    url.searchParams.delete(name);
  });
  return url.origin + url.pathname + url.search;
}

if (typeof document !== 'undefined' && window.top === window) {
  setInterval(function () {
    var bar = document.getElementById('__url_bar');
    if (!bar && document.body) {
      bar = document.createElement('div');
      bar.id = '__url_bar';
      bar.style.cssText =
        'position:fixed;left:0;right:0;bottom:0;z-index:2147483647;padding:4px 10px;' +
        'background:#202124;color:#e8eaed;font:13px/1.4 monospace;white-space:nowrap;' +
        'overflow:hidden;text-overflow:ellipsis;pointer-events:none';
      document.body.appendChild(bar);
    }
    if (bar) bar.textContent = displayUrl(location.href);
  }, 250);
}

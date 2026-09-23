local M = {}

M.servers = {
  bashls = {},
  cssls = {},
  html = {},
  jsonls = {},
  ts_ls = {},
  rust_analyzer = {},
  yamlls = {},
  dockerls = {},
  graphql = {},
  tailwindcss = {},
  lua_ls = {
    Lua = {
      workspace = { checkThirdParty = false },
      telemetry = { enable = false },
    },
  },
  gopls = {},
  buf_ls = {},
  sqlls = {
    settings = {
      sqls = {
        connections = {
          {
            driver = 'postgresql',
            dataSourceName = 'host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable',
          },
        },
      },
    },
  },
}

-- angularls is configured separately via vim.lsp.config() with a custom
-- root_dir, but still needs to be in Mason's ensure_installed list.
M.mason_extra = { 'angularls' }

M.mason_exclude = {}

M.mason_ensure_installed = function()
  local list = {}
  for name in pairs(M.servers) do
    if not M.mason_exclude[name] then
      table.insert(list, name)
    end
  end
  for _, name in ipairs(M.mason_extra) do
    table.insert(list, name)
  end
  table.sort(list)
  return list
end

-- Every server that should autostart via vim.lsp.enable(), including
-- angularls (configured separately).
M.enable_list = function()
  local list = {}
  for name in pairs(M.servers) do
    table.insert(list, name)
  end
  for _, name in ipairs(M.mason_extra) do
    table.insert(list, name)
  end
  table.sort(list)
  return list
end

return M

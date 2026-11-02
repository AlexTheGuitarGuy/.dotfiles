local M = {}

M.servers = {
  bashls = {},
  cssls = {},
  html = {},
  jsonls = {},
  ts_ls = {
    settings = {
      typescript = {
        updateImportsOnFileMove = { enabled = 'always' },
        suggest = { completeFunctionCalls = true },
      },
      javascript = {
        updateImportsOnFileMove = { enabled = 'always' },
        suggest = { completeFunctionCalls = true },
      },
    },
  },
  rust_analyzer = {},
  yamlls = {},
  dockerls = {},
  graphql = {},
  tailwindcss = {},
  eslint = {},
  lua_ls = {
    settings = {
      Lua = {
        workspace = { checkThirdParty = false },
        telemetry = { enable = false },
      },
    },
  },
  gopls = {
    settings = {
      gopls = {
        completeUnimported = true,
        usePlaceholders = true,
      },
    },
  },
  nixd = {},
  helm_ls = {},
  terraformls = {},
  gdscript = {},
  buf_ls = {},
  sqls = {
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

-- Every server that should autostart via vim.lsp.enable(), including
-- angularls (configured separately).
M.enable_list = function()
  local list = { 'angularls' }
  for name in pairs(M.servers) do
    table.insert(list, name)
  end
  table.sort(list)
  return list
end

return M

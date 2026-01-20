local servers = {
  bashls = {},
  cssls = {},
  html = {},
  jsonls = {},
  ts_ls = {},
  rust_analyzer = {},
  yamlls = {},
  dockerls = {},
  eslint_d = {},
  graphql = {},
  tailwindcss = {},
  -- svelte = {},
  lua_ls = {
    Lua = {
      workspace = { checkThirdParty = false },
      telemetry = { enable = false },
    },
  },
  gopls = {},
  buf_ls = {},

  -- csharp_ls = {},

  --[[ gopls = {},
  groovyls = {},
  gradle_ls = {}, ]]
}

function Format_without_lsp()
  vim.lsp.buf.format({
    async = true,
    filter = function(client)
      for key, _ in pairs(servers) do
        if client.name == key then
          return false
        end
      end

      return true
    end,
  })
end

local config = function()
  local lsp_zero = require('lsp-zero')
  lsp_zero.extend_lspconfig()

  lsp_zero.on_attach(function(client, bufnr)
    lsp_zero.default_keymaps({ buffer = bufnr })

    lsp_zero.on_attach(function(client, bufnr)
      lsp_zero.default_keymaps({ buffer = bufnr })

      vim.keymap.set('n', 'K', function()
        vim.lsp.buf.hover({
          border = 'rounded',
          title = ' Documentation ',
          -- Perfect padding (1 space all around)
          pad_top = 1,
          pad_bottom = 1,
          pad_left = 1,
          pad_right = 1,
          max_width = math.floor(vim.o.columns * 0.8),
          max_height = math.floor(vim.o.lines * 0.4),
        })
      end, { buffer = bufnr, desc = 'Hover docs' })
      vim.keymap.set('n', '<leader>la', function()
        vim.lsp.buf.code_action()
      end, opts)
      vim.keymap.set('n', '<leader>lr', function()
        vim.lsp.buf.rename()
      end, opts)
      vim.keymap.set('n', '<leader>lh', function()
        vim.diagnostic.open_float()
      end, opts)
      vim.keymap.set('n', 'gr', function()
        vim.lsp.buf.references()
      end, opts)
    end)

    local opts = { buffer = bufnr, remap = false }
  end)

  vim.lsp.config('angularls', {
    root_dir = vim.fs.root(0, { 'angular.json', 'project.json' }),
  })

  for name, opts in pairs(servers or {}) do
    vim.lsp.config(name, opts)
  end
end

return {
  'neovim/nvim-lspconfig',
  dependencies = {
    'hrsh7th/cmp-nvim-lsp',
  },
  event = { 'BufReadPre', 'BufNewFile' },
  cmd = 'LspInfo',
  config = config,
}

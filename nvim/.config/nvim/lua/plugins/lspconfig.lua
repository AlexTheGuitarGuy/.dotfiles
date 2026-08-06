local servers = require('core.lsp_servers').servers

local config = function()
  local lsp_zero = require('lsp-zero')
  lsp_zero.extend_lspconfig()

  vim.diagnostic.config({
    virtual_text = {
      prefix = '●',
      source = 'if_many',
    },
    signs = true,
    underline = true,
    update_in_insert = false,
    severity_sort = true,
  })

  lsp_zero.on_attach(function(client, bufnr)
    lsp_zero.default_keymaps({ buffer = bufnr })

    local opts = { buffer = bufnr, remap = false }

    vim.keymap.set('n', 'K', function()
      vim.lsp.buf.hover({
        border = 'rounded',
        title = ' Documentation ',
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

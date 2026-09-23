local lsp_servers = require('core.lsp_servers')
local servers = lsp_servers.servers

local config = function()
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

  vim.api.nvim_create_autocmd('LspAttach', {
    callback = function(args)
      local bufnr = args.buf
      local opts = { buffer = bufnr, remap = false }

      vim.keymap.set('n', 'K', function()
        vim.lsp.buf.hover({
          border = 'rounded',
          title = ' Documentation ',
          max_width = math.floor(vim.o.columns * 0.8),
          max_height = math.floor(vim.o.lines * 0.4),
        })
      end, { buffer = bufnr, desc = 'Hover docs' })
      vim.keymap.set('n', 'gd', function()
        vim.lsp.buf.definition()
      end, opts)
      vim.keymap.set('n', 'gD', function()
        vim.lsp.buf.declaration()
      end, opts)
      vim.keymap.set('n', 'gi', function()
        vim.lsp.buf.implementation()
      end, opts)
      vim.keymap.set('n', 'go', function()
        vim.lsp.buf.type_definition()
      end, opts)
      vim.keymap.set('n', 'gs', function()
        vim.lsp.buf.signature_help()
      end, opts)
      vim.keymap.set('n', 'gr', function()
        vim.lsp.buf.references()
      end, opts)
      vim.keymap.set('n', '<F2>', function()
        vim.lsp.buf.rename()
      end, opts)
      vim.keymap.set({ 'n', 'x' }, '<F3>', function()
        require('conform').format({ async = true, lsp_fallback = true })
      end, opts)
      vim.keymap.set('n', '<F4>', function()
        vim.lsp.buf.code_action()
      end, opts)
      vim.keymap.set('n', '<leader>la', function()
        vim.lsp.buf.code_action()
      end, opts)
      vim.keymap.set('n', '<leader>lr', function()
        vim.lsp.buf.rename()
      end, opts)
      vim.keymap.set('n', '<leader>lh', function()
        vim.diagnostic.open_float()
      end, opts)
    end,
  })

  vim.lsp.config('angularls', {
    root_dir = vim.fs.root(0, { 'angular.json', 'project.json' }),
  })

  for name, opts in pairs(servers or {}) do
    vim.lsp.config(name, opts)
  end

  vim.lsp.enable(lsp_servers.enable_list())
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

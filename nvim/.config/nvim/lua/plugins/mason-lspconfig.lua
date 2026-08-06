return {
  'williamboman/mason-lspconfig.nvim',
  dependencies = {
    'williamboman/mason.nvim',
    'neovim/nvim-lspconfig',
  },
  config = function()
    require('mason-lspconfig').setup({
      ensure_installed = require('core.lsp_servers').mason_ensure_installed(),
    })
  end,
}

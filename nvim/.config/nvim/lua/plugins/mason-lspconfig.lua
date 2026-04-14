return {
  'williamboman/mason-lspconfig.nvim',
  dependencies = {
    'williamboman/mason.nvim',
    'neovim/nvim-lspconfig',
  },
  config = function()
    require('mason-lspconfig').setup({
      ensure_installed = {
        'angularls',
        'cssls',
        'html',
        'jsonls',
        'sqlls',
        'ts_ls',
        'yamlls',
        'dockerls',
        'graphql',
        'lua_ls',
      },
    })
  end,
}

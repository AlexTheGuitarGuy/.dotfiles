return {
  'williamboman/mason-lspconfig.nvim',
  dependencies = {
    'williamboman/mason.nvim',
    'neovim/nvim-lspconfig',
  },
  config = function()
    require('mason-lspconfig').setup({
      'angularls',
      'cssls',
      'html',
      'jsonls',
      'sqlls',
      'tsserver',
      'yamlls',
      'dockerls',
      'eslint_d',
      'graphql',
      'lua_ls',
    })
  end,
}

return {
  'nvim-treesitter/nvim-treesitter-context',
  keys = {
    { '<leader>t', '<cmd>TSContext toggle<CR>', desc = 'Toggle Treesitter Context' },
  },
  config = function()
    require('treesitter-context').setup({
      enable = false,
    })
  end,
}

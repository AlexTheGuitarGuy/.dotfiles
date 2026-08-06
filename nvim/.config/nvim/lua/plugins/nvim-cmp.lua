local config = function()
  local cmp = require('cmp')
  local cmp_select = { behavior = cmp.SelectBehavior.Select }

  cmp.setup({
    sources = cmp.config.sources({
      { name = 'nvim_lsp' },
      { name = 'luasnip' },
      { name = 'buffer' },
      { name = 'path' },
    }),
    mapping = cmp.mapping.preset.insert({
      ['<CR>'] = cmp.mapping.confirm({ select = true }),
      ['<C-j>'] = cmp.mapping.select_next_item(cmp_select),
      ['<C-k>'] = cmp.mapping.select_prev_item(cmp_select),
    }),
  })
end

return {
  'hrsh7th/nvim-cmp',
  event = 'InsertEnter',
  dependencies = {
    'onsails/lspkind.nvim',
    'L3MON4D3/LuaSnip',
  },
  config = config,
}

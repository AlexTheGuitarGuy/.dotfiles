local config = function()
  vim.filetype.add({
    extension = {
      prisma = 'prisma',
    },
    pattern = {
      ['%.env%..+'] = 'sh',
    },
  })

  vim.api.nvim_create_autocmd('FileType', {
    pattern = '*',
    callback = function(args)
      pcall(vim.treesitter.start, args.buf)
    end,
  })

  require('nvim-treesitter').install({
    'typescript',
    'lua',
    'prisma',
    'nu',
    'markdown',
    'markdown_inline',
  })
end

return {
  'nvim-treesitter/nvim-treesitter',
  lazy = false,
  build = ':TSUpdate',
  config = config,
}

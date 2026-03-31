local config = function()
  vim.filetype.add({
    extension = {
      prisma = 'prisma',
    },
  })

  vim.api.nvim_create_autocmd('FileType', {
    pattern = '*',
    callback = function(args)
      pcall(vim.treesitter.start, args.buf)
    end,
  })

  require('nvim-treesitter.config').setup({
    indent = { enable = true },
    autotag = { enable = true },
    highlight = {
      enable = true,
      additional_vim_regex_highlighting = false,
    },
    ensure_installed = { 'typescript', 'lua', 'prisma', 'nu' },
  })
end

return {
  'nvim-treesitter/nvim-treesitter',
  lazy = false,
  build = ':TSUpdate',
  config = config,
}

local config = function()
  vim.filetype.add({
    extension = {
      prisma = 'prisma',
      bru = 'bruno',
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

  vim.api.nvim_create_autocmd('User', {
    pattern = 'TSUpdate',
    callback = function()
      require('nvim-treesitter.parsers').bruno = {
        install_info = {
          url = 'https://github.com/freska-fi/tree-sitter-bruno',
          revision = '82c9674c0b6c3815b90355364e990b2355f68cc1',
          queries = 'queries',
        },
      }
    end,
  })

  require('nvim-treesitter').install({
    'typescript',
    'javascript',
    'json',
    'graphql',
    'bruno',
    'lua',
    'prisma',
    'nu',
    'markdown',
    'markdown_inline',
    'gdscript',
    'godot_resource',
    'gdshader',
  })
end

return {
  'nvim-treesitter/nvim-treesitter',
  lazy = false,
  build = ':TSUpdate',
  config = config,
}

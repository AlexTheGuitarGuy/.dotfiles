local config = function()
  require('nvim-treesitter.config').setup({
    indent = { enable = true },
    autotag = { enable = true },
    highlight = {
      enable = true,
      additional_vim_regex_highlighting = false,
    },
    ensure_installed = { 'typescript', 'lua' },
  })

vim.treesitter.language.register("nu", "nushell")
end

return {
  'nvim-treesitter/nvim-treesitter',
  lazy = false,
  build = ':TSUpdate',
  config = config,
}

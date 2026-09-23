return {
  'catppuccin/nvim',
  name = 'catppuccin',
  priority = 1000,
  config = function(_, opts)
    require('catppuccin').setup(opts)
    vim.cmd.colorscheme('catppuccin')
  end,
  opts = {
    flavour = 'mocha',
    custom_highlights = function(colors)
      return {
        GitSignsCurrentLineBlame = { fg = colors.overlay2 },
      }
    end,
  },
}

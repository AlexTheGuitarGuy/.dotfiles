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
    integrations = {
      lualine = {
        all = function(colors)
          return {
            normal = { a = { bg = colors.mauve }, b = { fg = colors.mauve } },
            visual = { a = { bg = colors.pink }, b = { fg = colors.pink } },
            inactive = { a = { fg = colors.mauve } },
          }
        end,
      },
    },
    custom_highlights = function(colors)
      return {
        GitSignsCurrentLineBlame = { fg = colors.overlay2 },
      }
    end,
  },
}

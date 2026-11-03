local function accent()
  local ok, theme = pcall(function()
    local path = (vim.env.XDG_DATA_HOME or vim.env.HOME .. '/.local/share') .. '/dotfiles/theme.json'
    return vim.json.decode(table.concat(vim.fn.readfile(path), '\n'))
  end)
  return ok and theme.accent or 'blue'
end

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
          local main = colors[accent()]
          return {
            normal = { a = { bg = main }, b = { fg = main } },
            visual = { a = { bg = colors.pink }, b = { fg = colors.pink } },
            inactive = { a = { fg = main } },
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

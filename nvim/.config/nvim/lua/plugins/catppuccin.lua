return {
  "catppuccin/nvim",
  name = "catppuccin",
  priority = 1000,
  cmd = function()
    vim.cmd.colorscheme("catppuccin")
  end,
  opts = {
    flavour = "mocha",
  },
}

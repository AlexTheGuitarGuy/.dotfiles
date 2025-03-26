local M = {}

M.init = function()
  if vim.fn.exists("g:neovide") == 1 then
    vim.g.neovide_scale_factor = 0.8

    vim.o.guifont = "JetBrainsMono Nerd Font:style=Regular:h12"
  end
end

return M

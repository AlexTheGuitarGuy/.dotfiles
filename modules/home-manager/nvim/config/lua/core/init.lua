local M = {}

local highlight_group = vim.api.nvim_create_augroup('YankHighlight', { clear = true })
vim.api.nvim_create_autocmd('TextYankPost', {
  callback = function()
    vim.highlight.on_yank()
  end,
  group = highlight_group,
  pattern = '*',
})

local options = require('core.options')
local keymaps = require('core.keymaps')
local lazy = require('core.lazy')
local godot = require('core.godot')

M.init = function()
  options.init()
  keymaps.init()
  lazy.init()
  godot.init()
end

return M

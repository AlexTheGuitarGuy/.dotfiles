local M = {}

M.pipe = '/tmp/godot-nvim.pipe'

M.init = function()
  if vim.fn.filereadable(vim.fn.getcwd() .. '/project.godot') == 1 then
    pcall(vim.fn.serverstart, M.pipe)
  end
end

return M

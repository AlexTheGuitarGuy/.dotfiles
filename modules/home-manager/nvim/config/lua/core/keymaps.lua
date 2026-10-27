local M = {}

local options = { noremap = true, silent = true }
local apply_keymaps = function(mode, keymaps)
  for key, command in pairs(keymaps) do
    vim.keymap.set(mode, key, command, options)
  end
end

M.n = {
  -- No highlight
  ['<leader>H'] = ':noh<CR>',
  -- Let J stay in place
  ['J'] = 'mzJ`z',
  -- Let cursor stay in center
  ['<C-d>'] = '<C-d>zz',
  ['<C-u>'] = '<C-u>zz',
  -- Also let cursor stay in center when searching
  ['n'] = 'nzzzv',
  ['N'] = 'Nzzzv',
  -- Yank in clipboard
  ['<leader>y'] = '"+y',
  ['<leader>Y'] = '"+Y',
  -- Worst thing in the universe
  ['Q'] = '<nop>',
  -- Quick save and quit
  ['<C-s>'] = ':w!<CR>',
  ['<C-q>'] = ':q!<CR>',
  -- Better window navigation
  ['<C-h>'] = '<C-w>h',
  ['<C-j>'] = '<C-w>j',
  ['<C-k>'] = '<C-w>k',
  ['<C-l>'] = '<C-w>l',
  -- Resize with arrows
  ['<A-Up>'] = ':resize +2<CR>',
  ['<A-Down>'] = ':resize -2<CR>',
  ['<A-Left>'] = ':vertical resize -2<CR>',
  ['<A-Right>'] = ':vertical resize +2<CR>',
  -- Move text up and down
  ['<A-j>'] = ':m .+1<CR>==',
  ['<A-k>'] = ':m .-2<CR>==',

  -- Vertical split
  ['|'] = ':vsplit<CR>',

  -- Open current file with the OS's default handler (image viewer, etc.)
  -- gio reads the same mimeapps.list as xdg-open, but skips KDE's kde-open5
  -- detour, which doesn't always resolve the configured default silently.
  ['<leader>i'] = function()
    local opt = nil
    if vim.fn.has('mac') == 0 and vim.fn.has('win32') == 0 and vim.fn.executable('gio') == 1 then
      opt = { cmd = { 'gio', 'open' } }
    end
    vim.ui.open(vim.fn.expand('%:p'), opt)
  end,

  -- Open current file's directory in the OS's default file manager
  ['<leader>e'] = function()
    vim.ui.open(vim.fn.expand('%:p:h'))
  end,

  -- Convert line into branch name
  ['<leader>c'] = 'VuV:s/ /-/g<CR>:noh<CR>',
}
M.v = {
  -- Stay in indent mode
  ['<'] = '<gv',
  ['>'] = '>gv',
  -- Move text up and down
  -- Yank in clipboard
  ['y'] = '"+y',
  ['<A-j>'] = ":m '>+1<CR>gv=gv",
  ['<A-k>'] = ":m '<-2<CR>gv=gv",
  ['p'] = '"_dP',
}
M.x = {
  -- Move text up and down
  ['J'] = ":m '>+1<CR>gv=gv",
  ['K'] = ":m '<-2<CR>gv=gv",
  ['<A-j>'] = ":m '>+1<CR>gv=gv",
  ['<A-k>'] = ":m '<-2<CR>gv=gv",
}
M.i = {
  -- Press jk fast to exit insert mode
  ['jk'] = '<ESC>',
  ['kj'] = '<ESC>',
  ['JK'] = '<ESC>',
  ['KJ'] = '<ESC>',
}

M.init = function()
  vim.keymap.set('', '<Space>', '<Nop>', options)
  apply_keymaps('n', M.n)
  apply_keymaps('v', M.v)
  apply_keymaps('x', M.x)
  apply_keymaps('i', M.i)
end

return M

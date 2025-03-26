local M = {}

function _G.Tabline()
  local s = ""
  for i = 1, vim.fn.tabpagenr("$") do
    -- Add tab page number
    s = s .. "%" .. i .. "T"
    -- Add the file name with full path
    local buflist = vim.fn.tabpagebuflist(i)
    local winnr = vim.fn.tabpagewinnr(i)
    s = s .. " %#TabLineColor#" .. vim.fn.fnamemodify(vim.fn.bufname(buflist[winnr]), ":.") .. " "
    -- Highlight the current tab
    if i == vim.fn.tabpagenr() then
      s = s .. "%#TabLineSel#"
    else
      s = s .. "%#TabLine#"
    end
  end
  return s
end

-- Update tabline on buffer changes
vim.api.nvim_create_augroup("UpdateTabline", { clear = true })
vim.api.nvim_create_autocmd({ "BufEnter", "BufWinEnter", "TabEnter" }, {
  group = "UpdateTabline",
  callback = function()
    vim.o.tabline = "%!v:lua.Tabline()"
  end,
})

M.config = {
  backup = false, -- creates a backup file
  clipboard = "unnamedplus", -- allows neovim to access the system clipboard
  cmdheight = 2, -- more space in the neovim command line for displaying messages
  completeopt = { "menuone", "noselect" }, -- mostly just for cmp
  conceallevel = 0, -- so that `` is visible in markdown files
  fileencoding = "utf-8", -- the encoding written to a file
  hlsearch = true, -- highlight all matches on previous search pattern
  ignorecase = true, -- ignore case in search patterns
  mouse = "a", -- allow the mouse to be used in neovim
  pumheight = 10, -- pop up menu height
  showmode = false, -- we don't need to see things like -- INSERT -- anymore
  showtabline = 2, -- always show tabs
  tabline = Tabline(),
  smartcase = true, -- smart case
  smartindent = true, -- make indenting smarter again
  splitbelow = true, -- force all horizontal splits to go below current window
  splitright = true, -- force all vertical splits to go to the right of current window
  swapfile = false, -- creates a swapfile
  -- termguicolors = true,                    -- set term gui colors (most terminals support this)
  timeoutlen = 500, -- time to wait for a mapped sequence to complete (in milliseconds)
  undofile = true, -- enable persistent undo
  updatetime = 300, -- faster completion (4000ms default)
  writebackup = false, -- if a file is being edited by another program (or was written to file while editing with another program), it is not allowed to be edited
  expandtab = true, -- convert tabs to spaces
  shiftwidth = 2, -- the number of spaces inserted for each indentation
  tabstop = 2, -- insert 2 spaces for a tab
  cursorline = true, -- highlight the current line
  numberwidth = 4, -- set number column width to 2 {default 4}
  signcolumn = "yes", -- always show the sign column, otherwise it would shift the text each time
  wrap = true, -- display lines as one long line
  linebreak = true, -- companion to wrap, don't split words
  scrolloff = 8, -- minimal number of screen lines to keep above and below the cursor
  sidescrolloff = 8, -- minimal number of screen columns either side of cursor if wrap is `false`
  guifont = "monospace:h17", -- the font used in graphical neovim applications
  whichwrap = "bs<>[]hl", -- which "horizontal" keys are allowed to travel to prev/next line

  number = true,
  relativenumber = true,
  foldcolumn = "1", -- '0' is not bad
  foldlevel = 99, -- Using ufo provider need a large value, feel free to decrease the value
  foldlevelstart = 99,
  foldenable = true,
}

-- vim.opt.shortmess = "ilmnrx"                        -- flags to shorten vim messages, see :help 'shortmess'
vim.opt.shortmess:append("c") -- don't give |ins-completion-menu| messages
vim.opt.iskeyword:append("-") -- hyphenated words recognized by searches
vim.opt.formatoptions:remove({ "c", "r", "o" }) -- don't insert the current comment leader automatically for auto-wrapping comments using 'textwidth', hitting <Enter> in insert mode, or hitting 'o' or 'O' in normal mode.
vim.opt.runtimepath:remove("/usr/share/vim/vimfiles") -- separate vim plugins from neovim in case vim still in use

if vim.fn.has("wsl") == 1 then
  vim.g.clipboard = {
    name = "win32yank",
    copy = {
      ["+"] = "win32yank -i --crlf",
      ["*"] = "win32yank -i --crlf",
    },
    paste = {
      ["+"] = "win32yank -o --lf",
      ["*"] = "win32yank -o --lf",
    },
    cache_enabled = 0,
  }
end

M.init = function()
  for option, value in pairs(M.config) do
    vim.opt[option] = value
  end
end

return M

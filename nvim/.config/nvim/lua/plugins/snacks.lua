local keys = {
  { '<leader>a', function() Snacks.dashboard() end, desc = 'Open Dashboard' },
  { '<leader>z', function() Snacks.zen() end, desc = 'Toggle Zen Mode' },
  { '<leader>ff', function() Snacks.picker.files() end, desc = 'Find Files' },
  { '<leader>fw', function() Snacks.picker.grep() end, desc = 'Live Grep' },
  { '<leader>fu', function() Snacks.picker.grep_word() end, desc = 'Grep Word Under Cursor' },
  { '<leader>gf', function() Snacks.picker.git_files() end, desc = 'Git Files' },
  { '<leader>fr', function() Snacks.picker.recent() end, desc = 'Recent Files' },
}

return {
  'folke/snacks.nvim',
  priority = 1000,
  lazy = false,
  keys = keys,
  opts = {
    picker = {
      sources = {
        explorer = {
          hidden = true, -- show hidden files
          ignored = false, -- don't show gitignored files
          exclude = { -- exclude specific patterns
            '*.uid', -- glob pattern for files ending with .uid
            'server.pipe', -- exact filename match
          },
        },
        files = {
          hidden = true,
          exclude = { 'node_modules' },
        },
        grep = {
          hidden = true,
          exclude = { 'node_modules' },
        },
      },
    },
    dashboard = {
      preset = {
        header = [[
 ▄▀▀▀▀▄  ▄▀▀█▄   ▄▀▀▄ ▀▄  ▄▀▀█▄▄▄▄  ▄▀▀█▄   ▄▀▀▄ ▄▀▀▄  ▄▀▀█▀▄    ▄▀▀▄ ▄▀▄
█ █   ▐ ▐ ▄▀ ▀▄ █  █ █ █ ▐  ▄▀   ▐ ▐ ▄▀ ▀▄ █   █    █ █   █  █  █  █ ▀  █
   ▀▄     █▄▄▄█ ▐  █  ▀█   █▄▄▄▄▄    █▄▄▄█ ▐  █    █  ▐   █  ▐  ▐  █    █
▀▄   █   ▄▀   █   █   █    █    ▌   ▄▀   █    █   ▄▀      █       █    █
 █▀▀▀   █   ▄▀  ▄▀   █    ▄▀▄▄▄▄   █   ▄▀      ▀▄▀     ▄▀▀▀▀▀▄  ▄▀   ▄▀
 ▐      ▐   ▐   █    ▐    █    ▐   ▐   ▐              █       █ █    █
                ▐         ▐                           ▐       ▐ ▐    ▐]],
        keys = {
          { icon = ' ', key = 'f', desc = 'Find file', action = function() Snacks.picker.files() end },
          { icon = ' ', key = 'e', desc = 'New file', action = ':ene <BAR> startinsert<CR>' },
          { icon = ' ', key = 'r', desc = 'Recently used files', action = function() Snacks.picker.recent() end },
          { icon = '󱎸 ', key = 't', desc = 'Find text', action = function() Snacks.picker.grep() end },
          { icon = ' ', key = 'c', desc = 'Configuration', action = ':e $MYVIMRC<CR>' },
          { icon = '󰒲 ', key = 's', desc = 'Restore Session', section = 'session' },
          { icon = ' ', key = 'q', desc = 'Quit Neovim', action = ':qa<CR>' },
        },
      },
      sections = {
        { section = 'header' },
        { section = 'keys', gap = 1, padding = 1 },
        { section = 'startup' },
      },
    },
    zen = {
      win = {
        backdrop = { transparent = false, blend = 40 },
      },
    },
  },
}

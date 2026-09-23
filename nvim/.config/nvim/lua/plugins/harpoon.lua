local keys = {
  {
    '<leader>ha',
    function()
      require('harpoon'):list():add()
    end,
    mode = 'n',
  },
  {
    '<leader>hl',
    function()
      local h = require('harpoon')
      h.ui:toggle_quick_menu(h:list())
    end,
    mode = 'n',
  },
  {
    '<A-u>',
    function()
      require('harpoon'):list():select(1)
    end,
    mode = 'n',
  },
  {
    '<A-i>',
    function()
      require('harpoon'):list():select(2)
    end,
    mode = 'n',
  },
  {
    '<A-o>',
    function()
      require('harpoon'):list():select(3)
    end,
    mode = 'n',
  },
  {
    '<A-p>',
    function()
      require('harpoon'):list():select(4)
    end,
    mode = 'n',
  },
}

return {
  'ThePrimeagen/harpoon',
  branch = 'harpoon2',
  dependencies = {
    'plenary.nvim',
  },
  keys = keys,
}

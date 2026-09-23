return {
  'MagicDuck/grug-far.nvim',
  config = function()
    require('grug-far').setup({})
  end,
  keys = {
    {
      '<leader>S',
      function()
        require('grug-far').toggle_instance()
      end,
      desc = 'Toggle Grug Far',
      mode = 'n',
    },
    {
      '<leader>sw',
      function()
        require('grug-far').open({ prefills = { search = vim.fn.expand('<cword>') } })
      end,
      desc = 'Search current word',
      mode = 'n',
    },
    {
      '<leader>sw',
      function()
        require('grug-far').with_visual_selection()
      end,
      desc = 'Search current selection',
      mode = 'v',
    },
    {
      '<leader>sp',
      function()
        require('grug-far').open({ prefills = { search = vim.fn.expand('<cword>'), paths = vim.fn.expand('%') } })
      end,
      desc = 'Search on current file',
      mode = 'n',
    },
  },
}

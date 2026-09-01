local keys = {
  {
    '<leader>9v',
    function()
      require('99').visual()
    end,
    mode = 'v',
  },
  {
    '<leader>9s',
    function()
      require('99').search()
    end,
    mode = 'n',
  },
  {
    '<leader>9x',
    function()
      require('99').stop_all_requests()
    end,
    mode = 'n',
  },
  {
    '<leader>9m',
    function()
      require('99.extensions.telescope').select_model()
    end,
    mode = 'n',
  },
  {
    '<leader>9p',
    function()
      require('99.extensions.telescope').select_provider()
    end,
    mode = 'n',
  },
}

return {
  'ThePrimeagen/99',
  keys = keys,
  opts = function()
    local _99 = require('99')
    return {
      provider = _99.Providers.ClaudeCodeProvider,
      tmp_dir = './tmp',
      completion = { source = 'cmp' },
    }
  end,
}

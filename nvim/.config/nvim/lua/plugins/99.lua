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
      local pickers_util = require('99.extensions.pickers')
      pickers_util.get_models(nil, function(models, current)
        vim.ui.select(models, {
          prompt = '99: Select Model (current: ' .. current .. ')',
        }, function(choice)
          if choice then
            pickers_util.on_model_selected(choice)
          end
        end)
      end)
    end,
    mode = 'n',
  },
  {
    '<leader>9p',
    function()
      local pickers_util = require('99.extensions.pickers')
      local info = pickers_util.get_providers()
      vim.ui.select(info.names, {
        prompt = '99: Select Provider (current: ' .. info.current .. ')',
      }, function(choice)
        if choice then
          pickers_util.on_provider_selected(choice, info.lookup)
        end
      end)
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
      completion = { source = 'native' },
    }
  end,
}

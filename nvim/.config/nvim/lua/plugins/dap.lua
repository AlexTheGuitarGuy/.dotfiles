local keys = {
  {
    '<leader>dc',
    function()
      require('dap').continue()
    end,
    mode = 'n',
    desc = 'Debug: Continue/Start',
  },
  {
    '<leader>db',
    function()
      require('dap').toggle_breakpoint()
    end,
    mode = 'n',
    desc = 'Debug: Toggle Breakpoint',
  },
  {
    '<leader>do',
    function()
      require('dap').step_over()
    end,
    mode = 'n',
    desc = 'Debug: Step Over',
  },
  {
    '<leader>di',
    function()
      require('dap').step_into()
    end,
    mode = 'n',
    desc = 'Debug: Step Into',
  },
  {
    '<leader>dO',
    function()
      require('dap').step_out()
    end,
    mode = 'n',
    desc = 'Debug: Step Out',
  },
}

local function js_debug_adapter()
  return {
    type = 'server',
    host = 'localhost',
    port = '${port}',
    executable = {
      command = 'node',
      args = {
        vim.fn.stdpath('data') .. '/mason/packages/js-debug-adapter/js-debug/src/dapDebugServer.js',
        '${port}',
      },
    },
  }
end

local function config()
  local dap = require('dap')

  dap.adapters['pwa-node'] = js_debug_adapter()
  dap.adapters['pwa-chrome'] = js_debug_adapter()

  for _, language in ipairs({ 'typescript', 'javascript', 'typescriptreact', 'javascriptreact' }) do
    dap.configurations[language] = {
      {
        type = 'pwa-node',
        request = 'attach',
        name = 'Attach to process (--inspect)',
        processId = require('dap.utils').pick_process,
        cwd = '${workspaceFolder}',
        skipFiles = { '<node_internals>/**' },
      },
      {
        type = 'pwa-node',
        request = 'launch',
        name = 'Launch file',
        program = '${file}',
        cwd = '${workspaceFolder}',
        skipFiles = { '<node_internals>/**' },
      },
      {
        type = 'pwa-chrome',
        request = 'launch',
        name = 'Launch Chrome (dev server)',
        url = function()
          return vim.fn.input('Dev server URL: ', 'http://localhost:4200')
        end,
        webRoot = '${workspaceFolder}',
        sourceMaps = true,
      },
    }
  end
end

return {
  'mfussenegger/nvim-dap',
  keys = keys,
  config = config,
}

local function setup_keymaps(neotest)
  local opts = { remap = false }
  vim.keymap.set('n', '<leader>mr', neotest.run.run, opts)
  vim.keymap.set('n', '<leader>ms', neotest.run.stop, opts)
  vim.keymap.set('n', '<leader>mo', neotest.output.open, opts)
  vim.keymap.set('n', '<leader>mO', function()
    neotest.output.open({ enter = true })
  end, opts)
  vim.keymap.set('n', '<leader>mS', neotest.summary.toggle, opts)
  vim.keymap.set('n', '<leader>mf', function()
    neotest.run.run(vim.fn.expand('%'))
  end, opts)
  vim.keymap.set('n', '<leader>mn', function()
    neotest.jump.next({ status = 'failed' })
  end, opts)
  vim.keymap.set('n', '<leader>mp', function()
    neotest.jump.prev({ status = 'failed' })
  end, opts)
end

return {
  'nvim-neotest/neotest',
  dependencies = {
    'nvim-lua/plenary.nvim',
    'antoinemadec/FixCursorHold.nvim',
    'nvim-treesitter/nvim-treesitter',
    'nvim-neotest/neotest-jest',
    'marilari88/neotest-vitest',
  },
  config = function()
    local neotest = require('neotest')

    -- Neovim 0.12 changed iter_matches to always return TSNode[] per capture.
    -- Neotest hasn't been updated for this, so we intercept parse_positions and
    -- convert the string build_position into a normalizing wrapper function,
    -- which forces neotest to run on the main thread where our fix applies.
    local ts_lib = require('neotest.lib').treesitter
    local orig_parse_positions = ts_lib.parse_positions
    ts_lib.parse_positions = function(file_path, query, opts)
      if type(opts.build_position) == 'string' then
        local bp_str = opts.build_position
        local bp_func = assert(load('return ' .. bp_str))()
        opts = vim.tbl_extend('force', opts, {
          build_position = function(fp, src, captured_nodes, meta)
            local normalized = {}
            for k, v in pairs(captured_nodes) do
              normalized[k] = type(v) == 'table' and v[1] or v
            end
            return bp_func(fp, src, normalized, meta)
          end,
        })
      end
      return orig_parse_positions(file_path, query, opts)
    end

    neotest.setup({
      discovery = {
        enabled = false,
      },
      adapters = {
        require('neotest-jest')({
          jestCommand = 'node --expose-gc --no-compilation-cache ./node_modules/jest/bin/jest.js --logHeapUsage --colors --silent',
          jestConfigFile = function()
            local root = vim.fn.getcwd()
            for _, name in ipairs({ 'custom.jest.config.ts', 'jest.config.ts', 'jest.config.js', 'jest.config.mjs' }) do
              if vim.fn.filereadable(root .. '/' .. name) == 1 then
                return name
              end
            end
          end,
          env = { CI = true },
          filter_dir = function(name, rel_path, root)
            return name ~= 'node_modules'
          end,
          cwd = function(path)
            return vim.fn.getcwd()
          end,
        }),
        require('neotest-vitest'),
      },
      summary = {
        open = 'botright vsplit | vertical resize 80',
      },
    })

    setup_keymaps(neotest)
  end,
}


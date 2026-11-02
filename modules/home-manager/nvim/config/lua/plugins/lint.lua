return {
  'mfussenegger/nvim-lint',
  event = { 'BufReadPre', 'BufNewFile' },
  config = function()
    local lint = require('lint')
    lint.linters_by_ft = {
      nix = { 'statix', 'deadnix' },
      dockerfile = { 'hadolint' },
    }
    vim.api.nvim_create_autocmd({ 'BufWritePost', 'InsertLeave' }, {
      callback = function()
        lint.try_lint()
      end,
    })
  end,
}

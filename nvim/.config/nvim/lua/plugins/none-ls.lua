return {
  {
    'jay-babu/mason-null-ls.nvim',
    event = { 'BufReadPre', 'BufNewFile' },
    config = function()
      require('mason-null-ls').setup({
        ensure_installed = {
          'prettierd',
          'stylua',
          'eslint_d',
        },
      })
    end,
    lazy = true,
  },
  {
    'nvimtools/none-ls.nvim',
    dependencies = {
      'jay-babu/mason-null-ls.nvim',
      'nvimtools/none-ls-extras.nvim',
    },
    event = { 'BufReadPre', 'BufNewFile' },
    lazy = true,
    opts = function(_, opts)
      local nls = require('null-ls')
      local augroup = vim.api.nvim_create_augroup('LspFormatting', {})
      vim.diagnostic.config({
        virtual_text = {
          prefix = '●',
          source = 'if_many',
        },
        signs = true,
        underline = true,
        update_in_insert = false, -- Disable during insert for performance
        severity_sort = true,
      })
      opts.sources = vim.list_extend(opts.sources or {}, {
        nls.builtins.formatting.stylua,
        nls.builtins.formatting.prettierd,

        require('none-ls.code_actions.eslint_d'),
        require('none-ls.diagnostics.eslint_d'),
        require('none-ls.formatting.eslint_d'),
      })

      opts.on_attach = function(client, bufnr)
        if client.supports_method('textDocument/formatting') then
          vim.api.nvim_clear_autocmds({ group = augroup, buffer = bufnr })
          vim.api.nvim_create_autocmd('BufWritePost', {
            group = augroup,
            buffer = bufnr,
            callback = function()
              Format_without_lsp()
            end,
          })
        end
      end
    end,
  },
}

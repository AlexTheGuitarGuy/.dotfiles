return {
  'stevearc/conform.nvim',
  event = { 'BufReadPre', 'BufNewFile' },
  opts = {
    formatters_by_ft = {
      lua = { 'stylua' },
      javascript = { 'prettierd' },
      javascriptreact = { 'prettierd' },
      typescript = { 'prettierd' },
      typescriptreact = { 'prettierd' },
      json = { 'prettierd' },
      jsonc = { 'prettierd' },
      css = { 'prettierd' },
      html = { 'prettierd' },
      yaml = { 'prettierd' },
      markdown = { 'prettierd' },
      graphql = { 'prettierd' },
    },
    format_on_save = {
      lsp_fallback = true,
      timeout_ms = 500,
    },
  },
}

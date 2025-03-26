local opts = function()
  local ls = require("null-ls")
  local formatting = ls.builtins.formatting
  local diagnostics = ls.builtins.diagnostics
  local completion = ls.builtins.completion
  local code_actions = ls.builtins.code_actions

  -- Format on save
  --[[ local async_formatting = function(bufnr)
    bufnr = bufnr or vim.api.nvim_get_current_buf()

    vim.lsp.buf_request(
      bufnr,
      "textDocument/formatting",
      vim.lsp.util.make_formatting_params({}),
      function(err, res, ctx)
        if err then
          local err_msg = type(err) == "string" and err or err.message
          -- you can modify the log message / level (or ignore it completely)
          vim.notify("formatting: " .. err_msg, vim.log.levels.WARN)
          return
        end

        -- don't apply results if buffer is unloaded or has been modified
        if not vim.api.nvim_buf_is_loaded(bufnr) or vim.api.nvim_buf_get_option(bufnr, "modified") then
          return
        end

        if res then
          local client = vim.lsp.get_client_by_id(ctx.client_id)
          vim.lsp.util.apply_text_edits(res, bufnr, client and client.offset_encoding or "utf-16")
          vim.api.nvim_buf_call(bufnr, function()
            vim.cmd("silent noautocmd update")
          end)
        end
      end
    )
  end ]]

  local augroup = vim.api.nvim_create_augroup("LspFormatting", {})

  return {
    sources = {
      formatting.prettierd,
      formatting.black,
      formatting.stylua.with({
        filetypes = { "lua" },
        args = { "--indent-width", "2", "--indent-type", "Spaces", "-" },
      }),
      -- formatting.gci,
      formatting.gofumpt,
      formatting.buf,
      code_actions.refactoring,
      --[[ formatting.eslint_d,
      diagnostics.eslint_d,
      completion.spell, ]]
    },
    on_attach = function(client, bufnr)
      if client.supports_method("textDocument/formatting") then
        vim.api.nvim_clear_autocmds({ group = augroup, buffer = bufnr })
        vim.api.nvim_create_autocmd("BufWritePost", {
          group = augroup,
          buffer = bufnr,
          callback = function()
            -- async_formatting(bufnr)
            Format_without_lsp()
          end,
        })
      end
    end,
  }
end

return {
  "jose-elias-alvarez/null-ls.nvim",
  dependencies = {
    "plenary.nvim",
  },
  lazy = false,
  opts = opts,
}

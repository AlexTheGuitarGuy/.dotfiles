local config = function()
  local lsp_zero = require("lsp-zero")
  lsp_zero.extend_cmp()

  local cmp = require("cmp")
  local cmp_select = { behavior = cmp.SelectBehavior.Select }

  -- local lspkind = require("lspkind")
  cmp.setup({
    sources = cmp.config.sources({
      { name = "nvim_lsp" },
      { name = "luasnip" },
      { name = "buffer" },
      { name = "path" },
      { name = "cmp-tw2css" },
    }),
    mapping = cmp.mapping.preset.insert({
      ["<CR>"] = cmp.mapping.confirm({ select = true }),
      ["<C-j>"] = cmp.mapping.select_next_item(cmp_select),
      ["<C-k>"] = cmp.mapping.select_prev_item(cmp_select),
    }),
    formatting = {
      --[[ format = lspkind.cmp_format({
        maxwidth = 50,
        ellipsis_char = "...",
      }), ]]
      -- format = require("tailwindcss-colorizer-cmp").formatter,
    },
  })
end

return {
  "hrsh7th/nvim-cmp",
  event = "InsertEnter",
  dependencies = {
    "onsails/lspkind.nvim",
    "L3MON4D3/LuaSnip",
  },
  config = config,
}

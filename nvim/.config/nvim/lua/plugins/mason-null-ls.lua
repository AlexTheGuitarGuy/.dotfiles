return {
  "jay-babu/mason-null-ls.nvim",
  event = { "BufReadPre", "BufNewFile" },
  dependencies = {
    "williamboman/mason.nvim",
    "nvimtools/none-ls.nvim",
  },
  lazy = false,
  --[[ config = function()
    require("mason-null-ls").setup({
      "beautysh",
      "prettier",
      "prettierd",
      "eslint_d",
      "gitlint",
      "black",
      "stylua",
      "nxls",
      "rustywind",
    })
  end, ]]
}

return {
  --[[ "ThePrimeagen/refactoring.nvim",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "nvim-treesitter/nvim-treesitter",
  },
  config = function()
    require("refactoring").setup()
  end,
  keys = {
    ["<leader>rf"] = {
      function()
        vim.cmd("Refactor extract_to_file")
      end,
      "Extract to File",
      mode = "f",
    },
    ["<leader>rn"] = {
      function()
        vim.cmd("Refactor refactor_names")
      end,
      "Rename",
      mode = "n",
    },
    ["<leader>re"] = {
      function()
        vim.cmd("Refactor extract")
      end,
      "Extract",
      mode = "e",
    },
    ["<leader>rv"] = {
      function()
        vim.cmd("Refactor extract_var")
      end,
      "Extract Variable",
      mode = "v",
    },
    ["<leader>rb"] = {
      function()
        vim.cmd("Refactor extract_block")
      end,
      "Extract Block",
      mode = "b",
    },
    ["<leader>rI"] = {
      function()
        vim.cmd("Refactor inline_func")
      end,
      "Inline Function",
      mode = "I",
    },
    ["<leader>ri"] = {
      function()
        vim.cmd("Refactor inline_var")
      end,
      "Inline Variable",
      mode = "i",
    },
  }, ]]
}

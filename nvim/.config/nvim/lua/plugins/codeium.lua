return {
  "Exafunction/codeium.vim",
  event = "BufEnter",
  config = function()
    vim.keymap.set("i", "<M-l>", function()
      return vim.fn["codeium#Accept"]()
    end, { expr = true, silent = true })
    vim.keymap.set("n", "<leader>ca", ":CodeiumToggle<CR>")
  end,
}

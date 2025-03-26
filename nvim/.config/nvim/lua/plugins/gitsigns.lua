return {
  "lewis6991/gitsigns.nvim",
  opts = {
    signs = {
      add = { text = "│" },
      change = { text = "│" },
      delete = { text = "_" },
      topdelete = { text = "‾" },
      changedelete = { text = "~" },
      untracked = { text = "┆" },
    },
    signcolumn = true, -- Toggle with `:Gitsigns toggle_signs`
    numhl = false, -- Toggle with `:Gitsigns toggle_numhl`
    linehl = false, -- Toggle with `:Gitsigns toggle_linehl`
    word_diff = false, -- Toggle with `:Gitsigns toggle_word_diff`
    watch_gitdir = {
      follow_files = true,
    },
    attach_to_untracked = true,
    current_line_blame = false, -- Toggle with `:Gitsigns toggle_current_line_blame`
    current_line_blame_opts = {
      virt_text = true,
      virt_text_pos = "eol", -- 'eol' | 'overlay' | 'right_align'
      delay = 1000,
      ignore_whitespace = false,
      virt_text_priority = 100,
    },
    current_line_blame_formatter = "<author>, <author_time:%Y-%m-%d> - <summary>",
    sign_priority = 6,
    update_debounce = 100,
    status_formatter = nil, -- Use default
    max_file_length = 40000, -- Disable if file is longer than this (in lines)
    preview_config = {
      -- Options passed to nvim_open_win
      border = "single",
      style = "minimal",
      relative = "cursor",
      row = 0,
      col = 1,
    },
    on_attach = function(bufnr)
      vim.keymap.set(
        "n",
        "<leader>gp",
        require("gitsigns").prev_hunk,
        { buffer = bufnr, desc = "[G]o to [P]revious Hunk" }
      )
      vim.keymap.set("n", "<leader>gn", require("gitsigns").next_hunk, { buffer = bufnr, desc = "[G]o to [N]ext Hunk" })
      vim.keymap.set("n", "<leader>ph", require("gitsigns").preview_hunk, { buffer = bufnr, desc = "[P]review [H]unk" })
      vim.keymap.set("n", "<leader>gr", require("gitsigns").reset_hunk, { buffer = bufnr, desc = "[G]it [R]eset Hunk" })
      vim.keymap.set(
        "n",
        "<leader>gR",
        require("gitsigns").reset_buffer,
        { buffer = bufnr, desc = "[G]it [R]eset Buffer" }
      )
      vim.keymap.set("n", "<leader>gd", require("gitsigns").diffthis, { buffer = bufnr, desc = "[G]it [D]iff This" })
    end,
  },
}

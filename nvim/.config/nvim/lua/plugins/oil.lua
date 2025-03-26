return {
  "stevearc/oil.nvim",
  opts = {},
  -- Optional dependencies
  dependencies = { "nvim-tree/nvim-web-devicons" },
  config = function()
    vim.keymap.set("n", "-", ":Oil<CR>", { noremap = true, silent = true })
    local oil = require("oil")
    oil.setup({
      win_options = {
        signcolumn = "yes:2",
      },
      keymaps = {
        ["<C-s>"] = false,
        ["<C-l>"] = false,
        ["<C-h>"] = false,
        ["|"] = "actions.select_vsplit",
        ["Y"] = {
          callback = function()
            local entry = oil.get_cursor_entry()

            local dir = oil.get_current_dir()

            if not entry or not dir then
              return
            end

            local relpath = vim.fn.fnamemodify(dir, ":.")

            vim.fn.setreg("+", relpath .. entry.name)
          end,
        },
      },
    })
  end,
}

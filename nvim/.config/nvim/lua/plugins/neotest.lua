local function setup_keymaps(neotest)
  local opts = { remap = false }
  vim.keymap.set("n", "<leader>mr", neotest.run.run, opts)
  vim.keymap.set("n", "<leader>ms", neotest.run.stop, opts)
  vim.keymap.set("n", "<leader>mo", neotest.output.open, opts)
  vim.keymap.set("n", "<leader>mO", function()
    neotest.output.open({ enter = true })
  end, opts)
  vim.keymap.set("n", "<leader>mS", neotest.summary.toggle, opts)
  vim.keymap.set("n", "<leader>mf", function()
    neotest.run.run(vim.fn.expand("%"))
  end, opts)
  vim.keymap.set("n", "<leader>mn", function()
    neotest.jump.next({ status = "failed" })
  end, opts)
  vim.keymap.set("n", "<leader>mp", function()
    neotest.jump.prev({ status = "failed" })
  end, opts)
end

return {
  "nvim-neotest/neotest",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "antoinemadec/FixCursorHold.nvim",
    "nvim-treesitter/nvim-treesitter",
    "nvim-neotest/neotest-jest",
    "marilari88/neotest-vitest",
  },
  config = function()
    local neotest = require("neotest")

    neotest.setup({
      adapters = {
        require("neotest-jest")({
          jestCommand = "pnpm test --",
          jestConfigFile = "custom.jest.config.ts",
          env = { CI = true },
          cwd = function(path)
            return vim.fn.getcwd()
          end,
        }),
                require("neotest-vitest"),
      },
    })

    setup_keymaps(neotest)
  end,
}

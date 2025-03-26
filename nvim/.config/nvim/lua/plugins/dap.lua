local keys = function()
  local dap = require("dap")
  return {
    {
      "<F5>",
      function()
        dap.continue()
      end,
      mode = "n",
    },
    {
      "<F9>",
      function()
        dap.toggle_breakpoint()
      end,
      mode = "n",
    },
    {
      "<F10>",
      function()
        dap.step_over()
      end,
      mode = "n",
    },
    {
      "<F11>",
      function()
        dap.step_into()
      end,
      mode = "n",
    },
    {
      "<F12>",
      function()
        dap.step_out()
      end,
      mode = "n",
    },
  }
end

return {
  "mfussenegger/nvim-dap",
  keys = keys,
}

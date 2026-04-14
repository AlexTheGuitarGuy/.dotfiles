local keys = {
  { "<F5>",  function() require("dap").continue() end,          mode = "n" },
  { "<F9>",  function() require("dap").toggle_breakpoint() end, mode = "n" },
  { "<F10>", function() require("dap").step_over() end,         mode = "n" },
  { "<F11>", function() require("dap").step_into() end,         mode = "n" },
  { "<F12>", function() require("dap").step_out() end,          mode = "n" },
}

return {
  "mfussenegger/nvim-dap",
  keys = keys,
}

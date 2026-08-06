return {
  'folke/snacks.nvim',
  priority = 1000,
  opts = {
    picker = {
      sources = {
        explorer = {
          hidden = true, -- show hidden files
          ignored = false, -- don't show gitignored files
          exclude = { -- exclude specific patterns
            '*.uid', -- glob pattern for files ending with .uid
            'server.pipe', -- exact filename match
          },
        },
      },
    },
  },
}

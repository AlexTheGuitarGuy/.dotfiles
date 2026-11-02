local opts = { buffer = true }

vim.keymap.set('n', ']h', function()
  vim.fn.search('^#\\+\\s', 'W')
end, vim.tbl_extend('force', opts, { desc = 'Next heading' }))

vim.keymap.set('n', '[h', function()
  vim.fn.search('^#\\+\\s', 'bW')
end, vim.tbl_extend('force', opts, { desc = 'Previous heading' }))

vim.keymap.set('n', '<leader>x', function()
  local line = vim.fn.getline('.')
  if line:match('^%s*[-*] %[ %]') then
    line = line:gsub('%[ %]', '[x]', 1)
  elseif line:match('^%s*[-*] %[x%]') then
    line = line:gsub('%[x%]', '[ ]', 1)
  else
    return
  end
  vim.fn.setline('.', line)
end, vim.tbl_extend('force', opts, { desc = 'Toggle checkbox' }))

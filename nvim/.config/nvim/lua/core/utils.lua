local M = {}

function M.tabline()
  local s = ''
  for i = 1, vim.fn.tabpagenr('$') do
    -- Add tab page number
    s = s .. '%' .. i .. 'T'
    -- Add the file name with full path
    local buflist = vim.fn.tabpagebuflist(i)
    local winnr = vim.fn.tabpagewinnr(i)
    s = s .. ' %#TabLineColor#' .. vim.fn.fnamemodify(vim.fn.bufname(buflist[winnr]), ':.') .. ' '
    -- Highlight the current tab
    if i == vim.fn.tabpagenr() then
      s = s .. '%#TabLineSel#'
    else
      s = s .. '%#TabLine#'
    end
  end
  return s
end

return M

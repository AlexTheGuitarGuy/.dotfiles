local open_image_under_cursor = function()
  local node = require("nvim-tree.lib").get_node_at_cursor()
  if node and node.type == "file" then
    local file_extension = vim.fn.fnamemodify(node.name, ":e"):lower()
    local image_extensions = { "png", "jpg", "jpeg", "gif", "bmp", "svg" }

    if vim.tbl_contains(image_extensions, file_extension) then
      vim.fn.system(string.format("gwenview %s", vim.fn.shellescape(node.absolute_path)))
    else
      print("Not an image file")
    end
  else
    print("No file under cursor")
  end
end

--[[ return {
  "nvim-tree/nvim-tree.lua",
  config = function()
    require("nvim-tree").setup({
      update_focused_file = {
        enable = true,
        update_cwd = true,
      },
      view = {
        width = 40,
        side = "left",
        number = true,
        relativenumber = true,
      },
      diagnostics = {
        enable = true,
        show_on_dirs = true,
      },
      on_attach = function(bufnr)
        local api = require("nvim-tree.api")

        local function opts(desc)
          return { desc = "nvim-tree: " .. desc, buffer = bufnr, noremap = true, silent = true, nowait = true }
        end

        -- default mappings
        api.config.mappings.default_on_attach(bufnr)

        vim.keymap.set("n", "l", api.node.open.edit, opts("Open"))
        vim.keymap.set("n", "h", api.node.navigate.parent_close, opts("Close Directory"))
        vim.keymap.set("n", "v", api.node.open.vertical, opts("Open: Vertical Split"))
        vim.keymap.set("n", "gi", open_image_under_cursor, opts("Open Image Under Cursor"))
      end,
    })

    vim.keymap.set("n", "<leader>e", ":NvimTreeToggle<CR>", { noremap = true, silent = true })
    vim.keymap.set("n", "<leader>nf", ":NvimTreeFindFile<CR>", { noremap = true, silent = true })
  end,
} ]]
return {}

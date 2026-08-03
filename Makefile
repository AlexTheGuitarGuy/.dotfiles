.PHONY: stow-dotfiles unstow-dotfiles check

check:
	@echo "=== Dry-run: what stow --adopt would do ==="
	@stow -n --adopt --verbose=2 */ 2>&1 || true
	@echo "=== Existing dangling symlinks pointing into this repo ==="
	@find $(HOME) -maxdepth 4 -lname '*dotfiles*' -xtype l 2>/dev/null || true

stow-dotfiles: check
	@read -p "Proceed with stow --adopt? [y/N] " ans; [ "$$ans" = "y" ] || [ "$$ans" = "Y" ] || (echo "Aborted."; exit 1)
	@echo "Stowing dotfiles..."
	@stow --adopt */ --verbose=1
	@stow -D ansible-scripts/ scripts/

unstow-dotfiles:
	@echo "Unstowing dotfiles..."
	@stow -D */ --verbose=1

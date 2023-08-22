.PHONY: stow-dotfiles unstow-dotfiles

stow-dotfiles:
	@echo "Stowing dotfiles..."
	@stow */ --verbose=1
	@stow -D ansible-scripts/

unstow-dotfiles:
	@echo "Unstowing dotfiles..."
	@stow -D */ --verbose=1

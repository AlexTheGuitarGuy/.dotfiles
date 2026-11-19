{
  lib,
  buildGoModule,
  makeWrapper,
  mermaid-cli,
  chafa,
}:
# Terminal viewer for the investigate-code skill, copied from Daniel Snetcov's dotfiles.nix
# (commit 4fbb72b, since removed there).
buildGoModule {
  pname = "evidence";
  version = "0.1.0";

  src = ./src;

  vendorHash = "sha256-H35pwUB0i0Q8uqvdaoIOslCIl888hPlt1AU5oSqBCYU=";
  env.CGO_ENABLED = 0;

  nativeBuildInputs = [makeWrapper];
  postInstall = ''
    wrapProgram "$out/bin/evidence" --prefix PATH : ${lib.makeBinPath [mermaid-cli chafa]}
  '';

  meta.mainProgram = "evidence";
}

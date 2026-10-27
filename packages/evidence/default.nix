{
  lib,
  buildGoModule,
  fetchFromGitLab,
  makeWrapper,
  mermaid-cli,
  chafa,
}:
# Terminal viewer for the investigate-code skill, from Daniel Snetcov's dotfiles.nix.
buildGoModule {
  pname = "evidence";
  version = "0.1.0";

  src = fetchFromGitLab {
    owner = "Daniel-Snetcov";
    repo = "dotfiles.nix";
    rev = "4fbb72b996454dbd8a77c571ac48c00a4952839d";
    hash = "sha256-DGzBZJrE+2/6r/ptgEigFxPBuXEmd4zIC+zNTuDUSwA=";
  };
  sourceRoot = "source/packages/evidence";

  vendorHash = null;
  env.CGO_ENABLED = 0;

  nativeBuildInputs = [makeWrapper];
  postInstall = ''
    wrapProgram "$out/bin/evidence" --prefix PATH : ${lib.makeBinPath [mermaid-cli chafa]}
  '';

  meta.mainProgram = "evidence";
}

{
  python3Packages,
  fetchPypi,
}:
# Ren'Py language server; installed from the wheel, since its build backend needs a newer uv_build than nixpkgs has.
python3Packages.buildPythonApplication rec {
  pname = "renpy-language-server";
  version = "0.2.3";
  format = "wheel";

  src = fetchPypi {
    pname = "renpy_language_server";
    inherit version format;
    dist = "py3";
    python = "py3";
    hash = "sha256-kmsXb+PelYQTBJQ4z8UaOggy/TwzRg+8OUkvqV6BWEU=";
  };

  dependencies = with python3Packages; [pygls typer];
  pythonRelaxDeps = ["typer"];

  meta.mainProgram = "renpy-language-server";
}

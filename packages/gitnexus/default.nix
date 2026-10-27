{
  lib,
  stdenv,
  buildNpmPackage,
  fetchurl,
  autoPatchelfHook,
  python3,
  openssl,
}:
buildNpmPackage rec {
  pname = "gitnexus";
  version = "1.6.12";

  # The npm tarball ships the built dist/; package-lock.json is generated from it
  # (npm install --package-lock-only), since the tarball has none.
  src = fetchurl {
    url = "https://registry.npmjs.org/gitnexus/-/gitnexus-${version}.tgz";
    hash = "sha256-ljkabWbUSL8JeAhJUenO435UO3XkRc+a0s0TWNfRQeE=";
  };
  postPatch = ''
    cp ${./package-lock.json} package-lock.json
  '';
  npmDepsHash = "sha256-oqp1JttEF++/Me2eYUVdwQB8gAuj6LsRF2F1nAqF5Uo=";
  dontNpmBuild = true;

  # onnxruntime-node would download the CUDA provider from NuGet; embeddings run on the CPU.
  env.ONNXRUNTIME_NODE_INSTALL = "skip";

  # Drop prebuilt binaries for other platforms and the GPU providers, so autoPatchelf
  # only sees x86_64 CPU libraries.
  postInstall = ''
    cd $out/lib/node_modules/gitnexus
    find . -path '*/prebuilds/*' -mindepth 1 -maxdepth 4 -type d ! -name linux-x64 ! -name prebuilds -prune -exec rm -rf {} +
    rm -rf node_modules/@ladybugdb/core-{darwin,win32,linux-arm64}*
    find node_modules/onnxruntime-node/bin -mindepth 3 -maxdepth 3 -type d ! -path '*/linux/x64' -exec rm -rf {} +
    rm -f node_modules/onnxruntime-node/bin/*/linux/x64/libonnxruntime_providers_{cuda,tensorrt}.so
  '';

  nativeBuildInputs = [autoPatchelfHook python3];
  buildInputs = [stdenv.cc.cc.lib openssl];

  # Only ladybugdb needs patching (it links OpenSSL). The other prebuilt addons only need
  # libc and libstdc++, which node already loads, and patchelf breaks sharp's libvips.
  dontAutoPatchelf = true;
  postFixup = ''
    autoPatchelf $out/lib/node_modules/gitnexus/node_modules/@ladybugdb
  '';

  meta = {
    description = "Code knowledge graph CLI and MCP server";
    homepage = "https://github.com/abhigyanpatwari/GitNexus";
    mainProgram = "gitnexus";
  };
}

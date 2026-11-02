{
  lib,
  config,
  pkgs,
  ...
}: let
  live = config.lib.file.mkOutOfStoreSymlink;
  repo = "${config.dotfiles.repo}/modules/home-manager/opencode/config";
  link = dir:
    lib.mapAttrs' (name: _: lib.nameValuePair "opencode/${dir}${name}" {source = live "${repo}/${dir}${name}";})
    (lib.filterAttrs (_: type: type == "regular") (builtins.readDir ./config/${dir}));

  caveman = pkgs.fetchzip {
    url = "https://registry.npmjs.org/caveman-opencode-plugin/-/caveman-opencode-plugin-0.1.5.tgz";
    hash = "sha256-TZcRL+hm3EzD8/3z0dBLYpMfoOfGmJNy+aRndDuZC2w=";
  };
in {
  options.dotfiles.opencode.enable = lib.mkEnableOption "OpenCode";

  config = lib.mkIf config.dotfiles.opencode.enable {
    # nodejs for the npx-launched MCP servers in opencode.json
    home.packages = [pkgs.unstable.opencode pkgs.nodejs];

    # Each file is linked on its own: OpenCode installs its plugins into ~/.config/opencode
    # (node_modules, package.json), so the directory itself stays writable.
    xdg.configFile =
      removeAttrs (link "") ["opencode/opencode.json"]
      // link "agents/"
      // link "plugins/"
      // link "herdr-opencode/"
      // link "themes/"
      // {
        "opencode/plugins/caveman.js".text = ''
          import v1 from "${caveman}/dist/index.js";

          const session = "global";

          export default {
            id: "caveman",
            async setup(ctx) {
              const hooks = await v1();

              await ctx.session.hook("context", async (event) => {
                const out = {system: []};
                await hooks["experimental.chat.system.transform"]({sessionID: session}, out);
                for (const text of out.system) event.system.push({type: "text", text});
              });

              await ctx.command.transform((editor) => {
                for (const [name, description] of [
                  ["caveman-commit", "Generate commit messages in caveman style"],
                  ["caveman-review", "Review code in caveman style"],
                ]) {
                  editor.add({
                    name,
                    description,
                    execute: async ({sessionID, prompt, delivery}) => {
                      const out = {};
                      await hooks["command.execute.before"]({command: name, arguments: prompt.text ?? "", sessionID: session}, out);
                      const instruction = out.parts?.[0]?.text;
                      if (!instruction) return;
                      await ctx.session.prompt({...prompt, sessionID, text: `''${instruction}\n\n''${prompt.text ?? ""}`, delivery});
                    },
                  });
                }
              });
            },
          };
        '';
      };
  };
}

{
  description = "My simple statically generated website";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];
      imports = [inputs.treefmt-nix.flakeModule];
      perSystem = {
        config,
        pkgs,
        ...
      }: {
        packages = rec {
          generator = pkgs.callPackage ./default.nix {};
          website = pkgs.stdenv.mkDerivation {
            pname = "website";
            version = "0.0.1";
            src = ./.;
            buildInputs = [generator];
            buildPhase = ''
              generator
            '';
            installPhase = ''
              mkdir -p $out
              cp -r ./static/* $out
              cp -r ./generated/* $out
            '';
          };
          default = website;
        };
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
        treefmt.config = {
          projectRootFile = "flake.nix";
          programs = {
            alejandra.enable = true;
            deadnix.enable = true;
            gofumpt.enable = true;
            prettier.enable = true;
            statix.enable = true;
          };

          settings.formatter.prettier.options = ["--tab-width" "4"];
        };
      };
    };
}

{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, self, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux"];

      perSystem = {
        pkgs,
        ...
      }: {
        packages.website = pkgs.stdenvNoCC.mkDerivation {
          name = "website";
          src = ./.;
          nativeBuildInputs = [pkgs.zola];
          buildPhase = ''
            zola build
          '';
          installPhase = ''
            mv public $out
          '';
        };
        formatter = pkgs.alejandra;

        devShells.default = pkgs.mkShell {
          buildInputs = [pkgs.zola];
        };
      };
    };
}

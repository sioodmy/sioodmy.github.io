{
  # https://lukebentleyfox.net/posts/building-this-blog/
  # building zola is based on ^ blog post

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
      in {
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
        defaultPackage = pkgs.website;
        formatter = pkgs.alejandra;

        devShell = pkgs.mkShell {
          buildInputs = [pkgs.zola];
        };
      }
    );
}

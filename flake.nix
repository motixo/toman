{
  description = "Toman CLI Development and Release Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "toman";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;
          subPackages = [ "cmd/app" ];
          postInstall = ''mv $out/bin/app $out/bin/toman'';
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            goreleaser
            gopls
            go-tools
          ];

          shellHook = ''
            echo "--- Toman CLI Dev Environment ---"
            echo "To release, run: goreleaser release --snapshot --clean"
          '';
        };
      }
    );
}

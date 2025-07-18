{
  description = "http-Ostrich -- The ultimate fast ready to use HTTP server for easily distributing files inside a network";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };
  outputs =
    { self, nixpkgs, ... }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
      ]; # I don't have any darwin system to test

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          http-ostrich = pkgs.buildGoModule {
            pname = "http-ostrich";
            version = "0.9.1";
            src = ./.;
            vendorHash = "sha256-HObQqhTZLvt1XeiCN7i0e1WrXjbhuIaNzpI6c+FH7KY=";#nixpkgs.lib.fakeHash;

            buildInputs = [

            ];

          };
          default = http-ostrich;
        }
      );

      devShells = forAllSystems (system:
      let 
        pkgs = nixpkgsFor.${system};
      in 
      {
        default = pkgs.mkShell {
          packages = [
            pkgs.go
          ];
        };
      }
      );
    };
}

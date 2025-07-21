{
  description = "Ansicht, the mail lister";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };
      ansicht = pkgs.buildGoModule {
        pname = "ansicht";
        version = "0.0.1";
        src = ./.;
        # vendorHash = pkgs.lib.fakeHash;
        vendorHash = "sha256-0fohCgdRgu16HT/F3ixVrPVatw17CBv6dsx5WdWHqxM=";

        meta = {
          license = pkgs.lib.licenses.mit;
        };

        buildInputs = [ pkgs.notmuch ];

      };
    in {
      packages = {
        inherit ansicht;
        default = ansicht;
      };

      devShells.default = pkgs.mkShell {
        hardeningDisable = ["fortify"];
        buildInputs = with pkgs; [
          go
          notmuch
          gopls
          go-tools         # linter (`staticcheck`)
          delve            # debugger
          gdlv             # GUI for delve
          # golangci-lint    # linter (`golangci-lint run`), formatter
        ];

        shellHook = /*bash*/''
          echo "It's dangerous to go alone, take this!"
          echo "  go run main.go"
          echo "  dlv debug --headless main.go --listen=localhost:2345"
          echo "  gdlv connect localhost:2345"
          echo "  staticcheck -h"
        '';
      };
    });
}

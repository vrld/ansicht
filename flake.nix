{
  description = "Einsicht, the mail viewer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };
    in {
      devShells.default = pkgs.mkShell {
        hardeningDisable = ["fortify"];
        buildInputs = with pkgs; [
          go
          notmuch
          # uncomment when needed
          # gopls            # language server
          go-tools         # linter (`staticcheck`)
          delve            # debugger
          gdlv             # GUI for delve
          # golangci-lint    # linter (`golangci-lint run`), formatter
        ];

        shellHook = /*bash*/''
          echo "It's dangerous to go alone, take this!"
          echo "  go run main.go"
        '';
      };
    });
}

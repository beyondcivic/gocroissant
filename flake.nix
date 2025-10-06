{
  description = "Nix Flake Environment for Golang 1.24 Development";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go_1_24
          gopls
          golangci-lint
          powershell
          jq
          delve
          gcc
          gnumake
          zsh
        ];

        shellHook = ''


          # Check if we are already in a Nix shell
          if [ -z "$IN_NIX_SHELL" ]; then
            echo "Entering Nix Shell"
            exec nix develop
          fi

          # Set GOPATH to home directory's go folder
          export GOPATH=$HOME/go
          # Add GOPATH/bin to PATH
          export PATH=$GOPATH/bin:$PATH
          # Set CGO flags
          export CGO_ENABLED=1
          export CGO_CFLAGS="-O2"
          # Set zsh as default shell
          export SHELL=${pkgs.zsh}/bin/zsh
          
          echo "Go v1.24 development environment activated!"
          echo "Go version: $(go version)"
          echo "gopls version: $(gopls version)"
          echo "golangci-lint version: $(golangci-lint version)"
          echo "PowerShell version: $(pwsh -Version)"
          echo "jq version: $(jq --version)"
          echo "Delve version: $(dlv version)"
          echo "GCC version: $(gcc --version | head -n 1)"
          echo "Zsh version: $(zsh --version)"

          # Launch zsh
          exec zsh
        '';
        
        nativeBuildInputs = with pkgs; [
          pkg-config
        ];
      };
    });
}
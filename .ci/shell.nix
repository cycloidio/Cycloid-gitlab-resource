{ pkgs ? import <nixpkgs> { }, ... }:

pkgs.mkShell {
  packages = with pkgs; [
    # build tools
    go_1_25
    gcc
    golangci-lint

    # project tools
    just
    watchexec
  ];
}

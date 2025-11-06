{ pkgs ? import <nixpkgs> { }, ... }:

pkgs.mkShell {
  packages = with pkgs; [
    # build tools
    go
    gcc
    golangci-lint
    upx # executable compression

    # project tools
    just
    watchexec
  ];
}

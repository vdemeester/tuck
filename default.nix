let
  _pkgs = import <nixpkgs> {};
in
{ pkgs ? import (_pkgs.fetchFromGitHub { owner = "NixOS";
                                         repo = "nixpkgs-channels";
                                         rev = "0afb6d789c8bf74825e8cdf6a5d3b9ab8bde4f2d";
                                         sha256 = "147vhzrnwcy0v77kgbap31698qbda8rn09n5fnjp740svmkjpaiz";
                                       }) {}
}:

pkgs.stdenv.mkDerivation rec {
    name = "tuck";
    env = pkgs.buildEnv { name = name; paths = buildInputs; };
    buildInputs = [
        pkgs.go_1_8
        pkgs.vndr
        pkgs.gnumake
	pkgs.gotools
    ];
}

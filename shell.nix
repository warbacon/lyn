{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  name = "lyn";
  packages = with pkgs; [
    air
    go
    gopls
  ];
}

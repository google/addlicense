with import <nixpkgs> {};
mkShell {
  packages = [cowsay];
  shellHook = "echo \"Hello World!\" | cowsay";
}

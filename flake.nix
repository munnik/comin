{
  description = "Comin - GitOps for NixOS Machines";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }:
  let
    systems = [ "aarch64-linux" "x86_64-linux" ];
    forAllSystems = nixpkgs.lib.genAttrs systems;
    nixpkgsFor = forAllSystems (system: nixpkgs.legacyPackages.${system});
    optionsDocFor = forAllSystems (system:
      import ./nix/module-options-doc.nix (nixpkgsFor."${system}")
    );
  in {
    overlays.default = final: prev: {
      comin = final.callPackage ./nix/package.nix { };
    };

    packages = forAllSystems (system: {
      comin = nixpkgsFor."${system}".callPackage ./nix/package.nix { };
      default = self.packages."${system}".comin;
      generate-module-options = optionsDocFor."${system}".optionsDocCommonMarkGenerator;
    });
    checks = forAllSystems (system: {
      module-options-doc = optionsDocFor."${system}".checkOptionsDocCommonMark;
      # I don't understand why nix flake check does't build packages.default
      package = self.packages."${system}".comin;
    });

    nixosModules.comin = nixpkgs.lib.modules.importApply ./nix/module.nix { inherit self; };
    devShells.x86_64-linux.default = let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in pkgs.mkShell {
      buildInputs = [
        pkgs.go pkgs.godef pkgs.gopls
        pkgs.golangci-lint
      ];
    };
  };
}

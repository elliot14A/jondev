{
  description = "jondev - Build portifolios with ease";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

   outputs = {self, nixpkgs, flake-utils}: 
     flake-utils.lib.eachDefaultSystem (system: 
       let 
         pkgs = nixpkgs.legacyPackages.${system}; 
       in 
       {
         devShells.default = pkgs.mkShell {
           buildInputs = with pkgs; [
             go 
             air
             sqlite
             sqlite-interactive
             bun
           ];
           shellHook = ''
             echo "jondev - Build protifolios with ease";
             echo "Available tools:"
             echo "  * Bun $(bun --version)"
             echo "  * Go $(go version)"
             echo "  * SQLite $(sqlite3 --version)"
           '';
         };
     });
}

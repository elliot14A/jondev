{ pkgs, ... }: {

  languages = {
    go = {
      enable = true;
      package = pkgs.go;
    };
    typescript.enable = true;
  };

  packages = with pkgs; [
    air
    litecli
    sqlite
    sqlc
    go-migrate
    protobuf
    protoc-gen-go
    protoc-gen-go-grpc
    bun
  ];

  # Environment variables
  env = {
    DB_PATH = ".jondev/sqlite/jondev.db";
    MIGRATIONS_PATH = "./migrations";
    BUILD_DIR = "build";
    BINARY_NAME = "jondev-server";
    PROTO_DIR = "proto";
    GEN_DIR = "./proto/gen";
    NODE_ENV = "development";
  };

  # Development process management
  processes = {
    dev-server.exec = "air serve";
    dev-web.exec = "cd web-dashboard && bun run dev";
  };

  # Scripts/Commands
  scripts = {
    clean.exec = "rm -rf $BUILD_DIR && mkdir -p $BUILD_DIR";

    build-server.exec = "go build -o $BUILD_DIR/$BINARY_NAME main.go";

    build-web.exec = "cd web-dashboard && bun run build";

    build.exec = "build-server && build-web";

    generate-hash.exec = "go run main.go generate-hash";

    install.exec = ''
      cd web-dashboard && bun install
      go mod download
    '';

    db-init.exec = "touch $DB_PATH";

    migrate-create = {
      exec = ''
        read -p "Enter migration name: " name
        migrate create -ext sql -dir $MIGRATIONS_PATH -seq $name
      '';
      description = "Create a new migration file";
    };

    migrate-up = {
      exec = "migrate -database \"sqlite3://$DB_PATH\" -path $MIGRATIONS_PATH up";
      description = "Run all pending migrations";
    };

    migrate-down = {
      exec = "migrate -database \"sqlite3://$DB_PATH\" -path $MIGRATIONS_PATH down";
      description = "Rollback the last migration";
    };

    migrate-force = {
      exec = ''
        read -p "Enter version: " version
        migrate -database "sqlite3://$DB_PATH" -path $MIGRATIONS_PATH force $version
      '';
      description = "Force migration version";
    };

    migrate-version = {
      exec = "migrate -database \"sqlite3://$DB_PATH\" -path $MIGRATIONS_PATH version";
      description = "Show current migration version";
    };

    sqlc-gen = {
      exec = "sqlc generate";
      description = "Generate Go code from SQL";
    };

    proto = {
      exec = ''
        mkdir -p $GEN_DIR/hash/v1
        protoc --go_out=$GEN_DIR \
          --go_opt=module=github.com/elliot14A/jondev/gen \
          --go-grpc_out=$GEN_DIR \
          --go-grpc_opt=module=github.com/elliot14A/jondev/gen \
          $PROTO_DIR/v1/**/*
      '';
      description = "Generate protobuf code";
    };
  };

  # Pre-commit hooks
  pre-commit.hooks = {
    nixpkgs-fmt.enable = true;
    gofmt.enable = true;
    prettier.enable = true;
  };

  # Enter shell hook
  enterShell = ''
    echo "jondev - Build portfolios with ease"
    echo "Available tools:"
    echo "  * Bun $(bun --version)"
    echo "  * Go $(go version)"
    echo "  * SQLite $(sqlite3 --version)"
    echo "Development commands:"
    echo "  * devenv up        - Start development servers"
    echo "  * build - Build the project"
    echo "  * proto - Generate protobuf code"
    echo "  * sqlc-gen - Generate SQL code"
  '';
}

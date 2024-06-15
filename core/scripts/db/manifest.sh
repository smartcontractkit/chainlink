#!/bin/bash
root=$(git rev-parse --show-toplevel)
pushd $root
git log --format=""   --name-only --diff-filter=A -- core/store/migrate/migrations/*go core/store/migrate/migrations/*sql core/store/migrate/template/**/*.go core/store/migrate/template/**/*.sql > core/scripts/db/manifest.txt
popd
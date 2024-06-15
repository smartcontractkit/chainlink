#!/bin/bash
coreMigrationsPath=core/store/migrate/migrations
pluginMigrationsPath=core/store/migrate/plugins
root=$(git rev-parse --show-toplevel)
pushd $root > /dev/null
git log --format=""   --name-only --diff-filter=A -- $coreMigrationsPath/*go $coreMigrationsPath/*sql  $pluginMigrationsPath/**/*.sql > core/store/migrate/manifest.txt
popd > /dev/null
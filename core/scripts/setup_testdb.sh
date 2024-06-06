#/bin/sh

function exit_error {
    echo "Error: $1"
    exit 1
}
# Create a new user and database for development
# This script is intended to be run on a local development machine
tdir=$(mktemp -d -t db-dev-user)

username="chainlink_dev"
password="insecurepassword"
database="chainlink_development_test"
# here document for the SQL commands
cat << EOF > $tdir/db-dev-user.sql
DROP DATABASE IF EXISTS $database;
-- create a new user and database for development if they don't exist
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '$username') THEN
        CREATE ROLE $username WITH LOGIN PASSWORD '$password';
    END IF;
END \$\$;
SELECT 'CREATE DATABASE $database WITH OWNER $username;' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$database')\gexec

-- Grant all privileges on the database to the user
ALTER DATABASE $database OWNER TO $username;
GRANT ALL PRIVILEGES ON DATABASE "$database" TO "$username";
ALTER USER $username CREATEDB;

-- Create a pristine database for testing
SELECT 'CREATE DATABASE chainlink_test_pristine WITH OWNER $username;' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'chainlink_test_pristine')\gexec
EOF

# Print the SQL commands
echo "SQL commands to be run: $tdir/db-dev-user.sql"
echo "##########################################################################################################"
echo "##########################################################################################################"

cat $tdir/db-dev-user.sql
echo "##########################################################################################################"
echo "##########################################################################################################"
echo ""
# Run the SQL commands
psql -U postgres -h localhost -f $tdir/db-dev-user.sql || exit_error "Failed to create user $username and database $database"


#test the connection
PGPASSWORD=$password psql -U $username -h localhost -d $database -c "SELECT 1"  ||  exit_error "Connection failed for $username to $database"



db_url=$(echo "CL_DATABASE_URL=postgresql://$username:$password@localhost:5432/$database?sslmode=disable")
echo $db_url
repo=$(git rev-parse --show-toplevel)
pushd $repo
export $db_url 
make testdb || exit_error "Failed to create test database"
popd

# Set the database URL in the .dbenv file
dbenv=$repo/.dbenv
echo "\n!Success!\n"
echo "Datbase URL: $db_url"

echo "export $db_url" >> $dbenv
echo "Has been set in the $dbenv file"

echo "Either" 
echo "    source $dbenv"
echo "Or explicitly set environment variable in your shell"
echo "    export $db_url"
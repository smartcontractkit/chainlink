#/bin/sh

# Create a new user and database for development
# This script is intended to be run on a local development machine
tdir=$(mktemp -d -t db-dev-user)

username="chainlink_dev"
password="insecurepassword"
database="chainlink_development_test"
# here document for the SQL commands
cat << EOF > $tdir/db-dev-user.sql
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
psql -U postgres -h localhost -f $tdir/db-dev-user.sql


#test the connection
PGPASSWORD=$password psql -U $username -h localhost -d $database -c "SELECT 1" && echo "Connection successful" || echo "Connection failed"

db_url=$(echo "CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test")
echo $db_url
repo=$(git rev-parse --show-toplevel)
pushd $repo
export $db_url 
make testdb || echo "Failed to create test database"
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
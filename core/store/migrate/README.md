# Notes
- Node operators do not always run their migrations with 
super user priviledges so you cannot use ```CREATE EXTENSION```
- After adding a migration, run `tools/bin/db_schema_dump` to update `schema.sql`

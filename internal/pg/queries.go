package pg

const queryExistsTableWithName = `
SELECT tablename
FROM pg_tables
WHERE tablename = $1;`

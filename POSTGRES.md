# Postgres

installed: `brew install postgresql@17`
started service: `brew services start postgresql@17`
created db: `createdb heblang`

docs: https://www.postgresql.org/docs/current/

## Connect and query

connect: `psql heblang`

Top 100 most frequent words in `tanach_words`:

```sql
SELECT letters, COUNT(*) AS n
FROM tanach_words
GROUP BY letters
ORDER BY n DESC, letters
LIMIT 100;
```

quit psql: `\q`

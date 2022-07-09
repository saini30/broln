# Postgres support in broln

With the introduction of the `kvdb` interface, broln can support multiple database
backends. One of the supported backends is Postgres. This document
describes how it can be configured.

## Building broln with postgres support

Since `broln v0.14.1-beta` the necessary build tags to enable postgres support are
already enabled by default. The default release binaries or docker images can
be used. To build from source, simply run:

```shell
⛰  make install
```

## Configuring Postgres for broln

In order for broln to run on Postgres, an empty database should already exist. A
database can be created via the usual ways (psql, pgadmin, etc). A user with
access to this database is also required.

Creation of a schema and the tables is handled by broln automatically.

## Configuring broln for Postgres

broln is configured for Postgres through the following configuration options:

* `db.backend=postgres` to select the Postgres backend.
* `db.postgres.dsn=...` to set the database connection string that includes
  database, user and password.
* `db.postgres.timeout=...` to set the connection timeout. If not set, no
  timeout applies.

Example as follows:
```
[db]
db.backend=postgres
db.postgres.dsn=postgresql://dbuser:dbpass@127.0.0.1:5432/dbname
db.postgres.timeout=0
```
Connection timeout is disabled, to account for situations where the database
might be slow for unexpected reasons.

## Important note about replication

In case a replication architecture is planned, streaming replication should be avoided, as the master does not verify the replica is indeed identical, but it will only forward the edits queue, and let the slave catch up autonomously; synchronous mode, albeit slower, is paramount for `broln` data integrity across the copies, as it will finalize writes only after the slave confirmed successful replication.

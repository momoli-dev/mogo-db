# mogo-db

These are the common database utilities that can be used across various Go projects.

Most features are mainly focused on PostgreSQL with the support of PostGIS, using pgx as the database driver.

## Examples

Here are some code snippets demonstrating the usage of the package. You may also refer to the tests for more examples.

### New Connection

Start by creating the handler:

```go
conn := NewConn(ctx, &ConnParams{
    Addr: "postgres://user:pass@localhost:5432/dbname",
    HasPostgis: true,
})
defer conn.Close(ctx)
```

### Transaction Management

Transaction management is done through the `context.Context` object. They can be nested and each level will join the 
outer transaction. Outer most transaction will be committed or rolled back.

```go
// Functions are safe to call even if not in a transaction.
conn.BeginTx(ctx)
defer RollbackTx(ctx)

Tx(ctx)
InTx(ctx)
CommitTx(ctx)
```

### conn.WithTx Helper Method

`conn.WithTx` helper, provides a convenient way to run a function within a transaction

```go
conn.WithTx(ctx, func(ctx context.Context) error {
    // Do some database operations here
    if err != nil {
        // If an error is returned, the transaction will be rolled back
        return err
    }

    // If nil is returned, the transaction will be committed
    return nil
})
```

`conn.WithTx` can be nested, just like the transaction management functions.

```go
conn.WithTx(ctx, func(ctx context.Context) error {
    conn.WithTx(ctx, func(ctx context.Context) error {
        // Do some database operations here
        return nil
    })
    return nil
})
```

### Generic WithTx Function

Generic `WithTx` function is also provided for convenience, allowing the caller to return an additional result from
inner function.

```go
// A generic WithTx function is also provided for convenience
res, err := WithTx(ctx, conn, func(ctx context.Context) (T, error) {
    // Do some database operations here
    return res, nil
})
```


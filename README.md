# Database

A set of database types, driver and query builder for sql based databases.

## Nullable Types

database package contains nullable datatype for working with nullable data. nullable types implements **Scanners**, **Valuers**, **Marshaler** and **Unmarshaler** interfaces.

**Note:** You can use `Val` method to get variable nullable value.

**Note:** Slice types is a comma separated list of variable that stored as string in database. e.g.: "1,2,3,4"

### Available Nullable Types

```go
import "github.com/gomig/database/types"
var a types.NullBool
var a types.NullFloat32
var a types.Float32Slice
var a types.NullFloat64
var a types.Float64Slice
var a types.NullInt
var a types.IntSlice
var a types.NullInt8
var a types.Int8Slice
var a types.NullInt16
var a types.Int16Slice
var a types.NullInt32
var a types.Int32Slice
var a types.NullInt64
var a types.Int64Slice
var a types.NullString
var a types.StringSlice
var a types.NullTime
var a types.NullUInt
var a types.UIntSlice
var a types.NullUInt8
var a types.UInt8Slice
var a types.NullUInt16
var a types.UInt16Slice
var a types.NullUInt32
var a types.UInt32Slice
var a types.NullUInt64
var a types.UInt64Slice
```

## MySQL Driver

Create new MySQL connection. this function return a `"github.com/jmoiron/sqlx"` instance.

```go
// Signature:
NewMySQLConnector(host string, username string, password string, database string) (*sqlx.DB, error)

// Example:
import "github.com/gomig/database"
db, err := database.NewMySQLConnector("", "root", "root", "myDB")
```

## Postgres Driver

Create new Postgres connection. this function return a `"github.com/jmoiron/sqlx"` instance.

```go
// Signature:
NewPostgresConnector(host string, port string, user string, password string, database string) (*sqlx.DB, error)

// Example:
import "github.com/gomig/database"
db, err := database.NewPostgresConnector("localhost", "", "postgres", "", "")
```

## Query Builder

Make complex query use `Query` structure.

**Note:** You can use special `@in` keyword in your query and query builder make a `IN(params)` query for you.

### Query Builder Methods

#### Add

Add new query.

```go
// Signature:
Add(q Query)
```

#### Query

Get query string.

```go
// Signature:
Query() string
```

#### Params

Get query builder parameters.

```go
// Signature:
Params() []any
```

### Query Structure Fields

**Type** _(String)_: Determine query type `AND`, `OR`, etc.

**Query** _(String)_: Query string.

**Params** _[]any_: Query parameters.

**Closure** _bool_: Determine query is sub query or not.

```go
import "github.com/gomig/database"
import "fmt"
var qBuilder database.QueryBuilder
qBuilder.Add(database.Query{
    Query:  "firstname LIKE '%?%'",
    Params: []any{"john"},
})
qBuilder.Add(database.Query{
    Type: "AND",
    Query:  "role @in",
    Params: []any{"admin", "support", "user"},
})
qBuilder.Add(database.Query{
    Type: "AND",
    Query:  "age > ? AND age < ?",
    Params: []any{15, 30},
    Closure: true,
})
fmt.Print(qBuilder.Query()) // firstname LIKE '%?%' AND role IN(?, ?, ?) AND (age > ? AND age < ?)
fmt.Print(qBuilder.Params()) // ["john", "admin", "support", "user", 15, 30]
```

## Migration

A set of command for migrating and seeding SQL database. This package use pure SQL file for migrating.

Driver flag is optional and default driver used if this flag not passed. When using multiple database driver, this flag used by resolver function to get database driver by name.

Migration and seed files directory will get at register time and also can set by flags.

**Note:** This package use `"github.com/jmoiron/sqlx"` as database driver.

**Note:** All sub commands automatically registered when main migration command registered.

```bash
myApp migration [command] --driver --migration_dir --seed_dir
# or simply
myApp migration [command] -d -m -s
```

```go
// Signature:
MigrationCommand(resolver func(driver string) *sqlx.DB, defDriver string, migDir string, seedDir string) *cobra.Command

// Example
import "github.com/gomig/database/migration"
rootCmd.AddCommand(migration.MigrationCommand(myResolver, "--APP-DB", "./database/migrations", "./database/seeds"))
```

### Clear

Delete all database table. these command run `clean.sql` migration file.

```bash
myApp migration clear --driver
```

### Migrate

Migrate database.

```bash
myApp migration migrate --driver --migration_dir
```

### Migrated

Show migrated files list.

```bash
myApp migration migrated --driver
```

### Seed

Seed database.

```bash
myApp migration seed --driver --seed_dir
```

### Seeded

Show seeded files list.

```bash
myApp migration seeded --driver
```

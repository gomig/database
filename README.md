# Database

A set of database types, driver and query builder for sql based databases.

## Drivers

### MySQL Driver

Create new MySQL connection. this function return a `"github.com/jmoiron/sqlx"` instance.

```go
// Signature:
NewMySQLConnector(host string, username string, password string, database string) (*sqlx.DB, error)

// Example:
import "github.com/gomig/database"
db, err := database.NewMySQLConnector("", "root", "root", "myDB")
```

### Postgres Driver

Create new Postgres connection. this function return a `"github.com/jmoiron/sqlx"` instance.

```go
// Signature:
NewPostgresConnector(host string, port string, user string, password string, database string) (*sqlx.DB, error)

// Example:
import "github.com/gomig/database"
db, err := database.NewPostgresConnector("localhost", "", "postgres", "", "")
```

## Repository

Set of generic functions to work with database. For reading from database `Find` and `FindOne` function use `q` and `db` fields to map struct field to database column.

**Note:** `q` struct tag used to advanced field name in query.

**Note:** You must use `?` as placeholder. Repository functions will transform placeholder automatically to `$1, $2, etc..` automatically.

**Note:** You can implement `Decoder` interface to call struct `Decode() error` method after read by `Find` and `FindOne` functions.

**Note:** You can auto fill fields list by putting placeholder (`?`) in select field list.

```go
type User struct{
    Id    string `db:"id"`
    Name  string `db:"name"`
    Owner *string `q:"owners.name as owner" db:"owner"` // must used for custom field
}

users, err := database.Find[User](db, `SELECT ? FROM users LEFT JOIN owners ON users.owner_id = owners.id`, nil)
// this function generate following query string:
// SELECT id, name, owners.name as owner FROM users LEFT JOIN owners ON users.owner_id = owners.id
```

### Find

Read query results to struct slice. You can use `resolver` callback to manipulate record after read from database.

```go
// Signature:
func Find[T any](db *sqlx.DB, query string, resolver func(*T), args ...any) ([]T, error)
```

### FindOne

Read single result or return nil if not exists.

```go
// Signature:
func FindOne[T any](db *sqlx.DB, query string, resolver func(*T), args ...any) (*T, error);
```

### Count

Get count of documents.

```go
// Signature:
func Count(db *sqlx.DB, query string, args ...any) (int64, error);
```

### Insert

Insert struct to database. This function use `db` tag to map struct field to database column.

```go
// Signature:
func Insert(db *sqlx.DB, entity any, table string, driver Driver) (sql.Result, error);
```

### Update

Update struct in database. This function use `db` tag to map struct field to database column.

```go
// Signature:
func Update(db *sqlx.DB, entity any, table string, driver Driver, condition string, args ...any) (sql.Result, error)
```

## Query Builder

Make complex query use for sql `WHERE` command.

**Note:** You can use special `@in` keyword in your query and query builder make a `IN(param1, param2)` query for you.

```go
import "github.com/gomig/database"
import "fmt"

query := database.NewQuery(database.DriverPostgres)
query.And("firstname LIKE '%?%'", "John")
query.And("role @in", "admin", "support", "user")
query.OrClosure("age > ? AND age < ?", 15, 30)
fmt.Print(query.ToSQL(1)) // " firstname LIKE '%$1%' AND role IN ($2,$3,$4) OR (age > $5 AND age < $6)"
fmt.Print(query.Params()) // [John admin support user 15 30]
```

### And

Add new simple condition to query with `AND`.

```go
// Signature:
And(cond string, args ...any)
```

### Or

Add new simple condition to query with `OR`.

```go
// Signature:
Or(cond string, args ...any)
```

### AndClosure

Add new condition to query with `AND` in nested `()`.

```go
// Signature:
AndClosure(cond string, args ...any)
```

### OrClosure

Add new condition to query with `OR` in nested `()`.

```go
// Signature:
OrClosure(cond string, args ...any)
```

### ToSQL

Generate query with placeholder based on counter.

```go
// Signature:
ToSQL(counter int) string
```

### Params

Get list of query parameters.

```go
// Signature:
Params() []any
```

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

### Migrate Manually

```go
// Signature:
func ExecuteScripts(db *sqlx.DB, commands []MigrationScript) error


// Example
import "github.com/gomig/database/migration"
Commands := []migration.MigrationScript{
    {Name: "001-create-user-table", IsSeed: false, CMD: "CREATE TABLE IF NOT EXISTS ..."},
}
if err := migration.ExecuteScripts(myDb, commands); err != nil {
    panic(err.Error())
} else {
    fmt.Println("Migrated")
}
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

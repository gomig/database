# Database

A set of database types, driver and query builder for sql based databases.

## Drivers

### MySQL Driver

Create new MySQL connection. this function return a `"github.com/jmoiron/sqlx"` instance.

```go
// Signature:
NewMySQLConnector(host string, username string, password string, database string) (*sqlx.DB, error)

// Example:
import "github.com/gomig/database/v2"
db, err := database.NewMySQLConnector("", "root", "root", "myDB")
```

### Postgres Driver

Create new Postgres connection. this function return a `"github.com/jmoiron/sqlx"` instance.

```go
// Signature:
NewPostgresConnector(host string, port string, user string, password string, database string) (*sqlx.DB, error)

// Example:
import "github.com/gomig/database/v2"
db, err := database.NewPostgresConnector("localhost", "", "postgres", "", "")
```

## Repository

Set of generic functions to work with database. For reading from database `Find` and `FindOne` function use `q` and `db` fields to map struct field to database column.

**Note:** `q` struct tag used to advanced field name in query.

**Note:** You must use `?` as placeholder. Repository functions will transform placeholder automatically to `$1, $2` for postgres driver.

**Note:** You can implement `Decoder` interface to call struct `Decode() error` method after read by `Find` and `FindOne` functions.

**Note:** You can auto fill select columns from struct fields putting `SELECT ...` placeholder.

```go
type User struct{
    Id    string `db:"id"`
    Name  string `db:"name"`
    Owner *string `q:"owners.name as owner" db:"owner"` // must used for custom field
}

users, err := database.Find[User](db, `SELECT ... FROM users LEFT JOIN owners ON users.owner_id = owners.id`, nil)
// this function generate following query string:
// SELECT id, name, owners.name as owner FROM users LEFT JOIN owners ON users.owner_id = owners.id
```

### Find

Read query results to struct slice. You can use `resolver` callback to manipulate record after read from database.

```go
// Signature:
func Find[T any](db *sqlx.DB, driver database.Driver, query string, resolver func(*T), args ...any) ([]T, error)
```

### FindOne

Read single result or return nil if not exists.

```go
// Signature:
func FindOne[T any](db *sqlx.DB, driver database.Driver, query string, resolver func(*T), args ...any) (*T, error);
```

### Count

Get count of documents.

```go
// Signature:
func Count(db *sqlx.DB, driver database.Driver, query string, args ...any) (int64, error);
```

### Insert

Insert struct to database. This function use `db` tag to map struct field to database column.

```go
// Signature:
func Insert(db *sqlx.DB, entity any, table string, driver database.Driver) (sql.Result, error);
```

### InsertTx

Insert struct to database using transaction. This function use `db` tag to map struct field to database column.

```go
// Signature:
func InsertTx(tx *sql.Tx, entity any, table string, driver database.Driver) (sql.Result, error);
```

### Update

Update struct in database. This function use `db` tag to map struct field to database column.

```go
// Signature:
func Update(db *sqlx.DB, entity any, table string, driver database.Driver, condition string, args ...any) (sql.Result, error)
```

### UpdateTx

Update struct in database using transaction. This function use `db` tag to map struct field to database column.

```go
// Signature:
func UpdateTx(tx *sql.Tx, entity any, table string, driver database.Driver, condition string, args ...any) (sql.Result, error)
```

## Query Builder

Make complex query use for sql `WHERE` command.

**Note:** You can use special `@in` keyword in your query and query builder make a `IN(param1, param2)` query for you.

```go
import "github.com/gomig/database/v2"
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

### ToString

Generate query string and replace `@q` with `ToSQL()`.

```go
// Signature:
ToString(pattern string, counter int, params ...any) string

// example
import "github.com/gomig/database/v2"
query := database.NewQuery(database.DriverPostgres).
                And("name = ?", "John Doe").
                And("id = ?", 3)
sql := query.ToString("SELECT * FROM users WHERE @q ORDER BY %s %s;", 1, "name", "asc");
// SELECT * FROM users WHERE name = $1 AND id = $2 ORDER BY name asc;
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
import "github.com/gomig/database/v2/types"
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

Advanced migration for SQL based database.

**Note:** This package use `"github.com/jmoiron/sqlx"` as database driver.

```bash
myApp migration [command]
```

```go
// Signature:
MigrationCommand(db *sqlx.DB, root string) *cobra.Command

// Example
import "github.com/gomig/database/v2/migration"
rootCmd.AddCommand(migration.MigrationCommand(myDB, "./database"))
```

### Migration Script Structure

Each migration script or file can contains 4 main section and defined with `--- [SECTION <name>]` line. Each migration file can contains 4 section:

- **UP:** Scripts on this section will run with `migration migrate` command.
- **SCRIPT:** Scripts on this section will run with `migration script` command.
- **SEED:** Scripts on this section will run with `migration seed` command.
- **DOWN:** Scripts on this section will run with `migration down` command.

**Note:** For writing multiple SQL script in single section you could add `-- [br]` in end of your command.

### Usage

#### new

This command create a new timestamp based standard migration file.

Flags:

- `-d` or `--dir`: used to define directory of files.

```bash
myApp migration new "create user" -d "my sub/directory/path"
```

#### summery

Show summery of migration executed on database.

```bash
myApp migration summery
```

#### up

Run `UP` scripts.

Flags:

- `-d` or `--dir`: used to define directory of files.
- `-n` or `--name`: used to run special script only.

```bash
myApp migration up -n "create user"
```

#### script

Run `SCRIPT` scripts.

Flags:

- `-d` or `--dir`: used to define directory of files.
- `-n` or `--name`: used to run special script only.

```bash
myApp migration script -d "some\sub\dir"
```

#### seed

Run `SEED` scripts.

Flags:

- `-d` or `--dir`: used to define directory of files.
- `-n` or `--name`: used to run special script only.

```bash
myApp migration seed
```

#### down

Run `DOWN` scripts to rollback migrations.

Flags:

- `-d` or `--dir`: used to define directory of files.
- `-n` or `--name`: used to run special script only.

```bash
myApp migration down
```

### Helpers Function

#### Migrate

This function run "UP" scripts from _migrations list_ on database and return succeeded list as result.

```go
// Signature:
func Migrate(db *sqlx.DB, migrations []migration.MigrationsT, name string) ([]string, error)
```

#### Script

This function run "SCRIPT" scripts from _migrations list_ on database and return succeeded list as result.

```go
// Signature:
func Script(db *sqlx.DB, migrations []migration.MigrationsT, name string) ([]string, error)
```

#### Seed

This function run "SEED" scripts from _migrations list_ on database and return succeeded list as result.

```go
// Signature:
func Seed(db *sqlx.DB, migrations []migration.MigrationsT, name string) ([]string, error)
```

#### Rollback

This function run "DOWN" scripts from _migrations list_ on database and return succeeded list as result.

```go
// Signature:
func Seed(db *sqlx.DB, migrations []migration.MigrationsT, name string) ([]string, error)
```

#### ReadDirectory

This function read migration files to `[]migration.MigrationsT` entity.

```go
// Signature:
func ReadDirectory(dir string) (MigrationsT, error)
```

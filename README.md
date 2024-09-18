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

Set of generic functions to work with database.

**Note:** You must use `?` as placeholder. Repository functions will transform placeholder automatically to `$1, $2` for numeric args mode.

**Note:** SQL placeholders cast as numeric `$1, $2` by default. You can change this behavior with `NumericArgs(false)` method.

**Note:** You can use replace phrase in your query string using `@some` in your query and replace with dynamic value for cleaner code.

### Commander

Normalize sql placeholder and execute.

```go
import "github.com/gomig/database/v2"

// -> UPDATE users SET name = $1 WHERE id = $2;
result, err := database.NewCMD(myDatabase).
    Command(`UPDATE users SET name = ? WHERE @cond;`).
    Replace("@cond", "id = ?").
    Exec("John Doe", 8921)
```

**NumericArgs** specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder.

**Command** set sql command **(Required)**.

**Replace** replace phrase in query string before run.

**Exec** normalize command and exec.

### Counter

Count records.

```go
import "github.com/gomig/database/v2"

// -> SELECT COUNT(id) FROM users WHERE name ILIKE '%$1%';
count, err := database.NewCounter(myTx).
    Query(`SELECT COUNT(id) FROM users WHERE @cond;`).
    Replace("@cond", "name ILIKE '%?%'").
    Result("John")
```

**NumericArgs** specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder.

**Query** set sql query **(Required)**.

**Replace** replace phrase in query string before run.

**Result** get count, returns -1 on error.

### Finder

Find single or multiple record.

Finder use `q` and `db` struct tag to map struct field to database column. If q or db struct tag set to `"-"` field will ignored.

**Note:** `q` struct tag used to advanced field name in query.

**Note:** You can implement `Decoder` interface in your struct to manipulate record after read. You could register extra resolver on Finder.

**Note:** You can auto fill select columns from struct by `@fields` placeholder in sql select statement. e.g. `SELECT @fields FROM users`.

```go
type User struct{
    Id      int `db:"id"`
    Name    string `db:"name"`
    Address string `q:"addresses.address AS address" db:"address"`
}
```

```go
import (
    "time"
    "strings"
    "github.com/gomig/database/v2"
)

type User struct{
    Id      int     `db:"id"`
    Name    string  `db:"name"`
    QueryAt time.Time `q:"-" db:"query_at"` // ignore from auto fill and select manually
}

// -> SELECT id, name, NOW() AS query_at FROM users WHERE id = $1;
single, err := database.NewFinder[User](db).
    Query(`SELECT @fields, NOW() AS query_at FROM users WHERE id = ?;`).
    Resolve(func(user *User) error {
        user.Name = strings.ToUpper(user.Name)
        return nil
    }).
    Single(3)


// -> SELECT id, name, NOW() AS query_at FROM users;
all, err := database.NewFinder[User](db).
    Query(`SELECT @fields, @queryAt FROM users;`).
    Replace("@queryAt", "NOW() AS query_at").
    Resolve(func(user *User) error {
        user.Name = strings.ToUpper(user.Name)
        return nil
    }).
    Result()
```

**NumericArgs** specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder

**Query** set sql query **(Required)**.

**Replace** replace phrase in query string before run.

**Resolve** reginster new resolver to run on record after read.

**Single** get first result.

**Result** get multiple result.

### Inserter

Insert struct to database. Inserter use `db` struct tag to resolve fields. If field is private or `db` tag is empty or equals `"-"` field ignored.

```go
import (
    "strings"
    "github.com/gomig/database/v2"
)

type User struct{
    Id      int     `db:"id"`
    Name    string  `db:"name"`
    Temp    string `db:"-"` // not inserted to database
}

// -> INSERT INTO users (id, name) VALUES(?, ?);
result, err := database.NewInserter[User](db).
    NumericArgs(false).
    Table("users").
    Insert(User {
        Id: 6,
        Name: "Jack Ma",
    })
```

**NumericArgs** specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder.

**Table** table name **(Required)**.

**Insert** insert and return result

### Updater

Update struct to database. Updater use `db` struct tag to resolve fields. If field is private or `db` tag is empty or equals `"-"` field ignored.

```go
import (
    "strings"
    "github.com/gomig/database/v2"
)

type User struct{
    Id      int     `db:"-"`
    Name    string  `db:"name"`
}

john := User {
    Id: 6,
    Name: "Jack Ma",
}

// -> UPDATE users SET name = $1 WHERE id = $2;
result, err := database.NewUpdater[User](db).
    Table("users").
    Where("id = ?", john.Id).
    Update(john)
```

**NumericArgs** specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder.

**Table** table name **(Required)**.

**Where** update condition **(Required)**.

**Update** update and return result.

## Query Builder

Make complex query use for sql `WHERE` command.

**Note:** You can use special `@in` placeholder in your query to make a `IN(param1, param2)` query for you.

**Note:** You can use special `@where` placeholder in your query to replace with `WHERE Raw()` value.

**Note:** You can use special `@query` placeholder in your query to replace with `Raw()` value.

```go
import (
    "fmt"
    "github.com/gomig/database/v2"
)

query := database.NewQuery().
    And("firstname LIKE '%?%'", "John").
    AndIf(myConditionPassed, "role @in", "admin", "support", "user").
    OrClosure("age > ? AND age < ?", 15, 30).
    OrIf(false, "id = ?", 5). // ignored because condition (first argument) not true
    Replace("@sort", "name").
    Replace("@order", "ASC")

// -> firstname LIKE '%$5%' AND role IN ($6, $7, $8) OR (age > $9 AND age < $10)
raw := query.Raw()

// -> SELECT * users WHERE firstname LIKE '%$5%' AND role IN ($6, $7, $8) OR (age > $9 AND age < $10) ORDER BY name ASC;
cmd := query.SQL(`SELECT * FROM USERS @where ORDER BY @sort @order;`) //

// -> [John admin support user 15 30]
args := query.Args()
```

**And** add new simple condition to query with AND.

**AndIf** add new And condition if first parameter is true.

**Or** add new simple condition to query with OR.

**OrIf** add new Or condition if first parameter is true.

**AndClosure** add new condition to query with AND in nested `()`.

**AndClosureIf** add new AndClosure condition if first parameter is true.

**OrClosure** add new condition to query with OR in nested `()`.

**OrClosureIf** add new AndClosure condition if first parameter is true.

**NumericArgs** specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder.

**NumericStart** set numeric argument start for numeric args mode.

**Replace** replace phrase in query string before run.

**Raw** get raw generated query.

**SQL** use generated query in part of sql command. this method replace `@query` with `Raw()` and `@where` with `WHERE Raw()` value.

**Args** get list of arguments.

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

Advance stage based migration for SQL based database.

**Note:** This package use `"github.com/jmoiron/sqlx"` as database driver.

**Note:** You can pass stages list to automatically run on migrate.

```bash
myApp migration [command]
```

```go
// Signature:
MigrationCommand(db *sqlx.DB, root string, extension string, authExecute ...string) *cobra.Command

// Example
import "github.com/gomig/database/v2/migration"
rootCmd.AddCommand(migration.MigrationCommand(myDB, "./database", "sql", "up", "index", "views"))
```

### Migration Script Structure

Each migration script or file can contains multiple stage `--- [STAGE <name>]` line. Each migration can have **DOWN** stage to rollback migration.

**Note:** For writing multiple SQL script in single section you could add `-- [end]` in end of your command.

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

#### run

Run stages scripts. Run auto execute stages list in order if not stage passed.

Flags:

- `-d` or `--dir`: used to define directory of files.
- `-n` or `--name`: used to run special script only.

```bash
# run all (up stage) scripts and then run all (index stage) scripts
myApp migration run up index
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

#### ReadFS

Read migration from file system.

```go
// Signature
ReadFS(dir, ext string) (Files, error)
```

#### InitMigration

Prepare database to run migrations.

```go
// Signature:
func InitMigration(db *sqlx.DB) error
```

#### Migrate

Run migration on database.

```go
// Signature:
func Migrate(db *sqlx.DB, stage string, files ...File) ([]string, error)
```

#### Rollback

This function run "DOWN" scripts from _migrations list_ on database and return succeeded list as result.

```go
// Signature:
func Rollback(db *sqlx.DB, files ...File) ([]string, error)
```

#### StageMigrated

Get migrated items for stage.

```go
// Signature:
StageMigrated(db *sqlx.DB, stage string) (Migrations, error)
```

#### Migrated

Get migrated list.

```go
// Signature:
Migrated(db *sqlx.DB) (Migrations, error)
```

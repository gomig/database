package database

type Driver string

const DriverMySQL Driver = "mysql"
const DriverPostgres Driver = "postgres"

type IDecoder interface {
	Decode() error
}

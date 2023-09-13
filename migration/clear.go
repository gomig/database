package migration

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func clearCMD(resolver func(driver string) *sqlx.DB) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "clear"
	cmd.Short = "delete all database table"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var err error
		driver, err := cmd.Flags().GetString("driver")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		db := resolver(driver)
		if db == nil {
			fmt.Printf("failed: %s database driver not found\n", driver)
			return
		}

		_, err = db.Exec("SET FOREIGN_KEY_CHECKS=0;")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		res, err := db.Query("SHOW TABLES;")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		var tables []string
		for res.Next() {
			var table string
			res.Scan(&table)
			tables = append(tables, table)
		}

		for _, table := range tables {
			_, err := db.Exec("DROP TABLE IF EXISTS " + table)
			if err != nil {
				fmt.Printf("failed: %s\n", err.Error())
				return
			}
		}

		_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1;")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		fmt.Println("cleared!")
	}
	return cmd
}

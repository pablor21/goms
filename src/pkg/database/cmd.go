package database

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func AddDatabaseCmd(rootCmd *cobra.Command) {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
	}

	migrateUp := &cobra.Command{
		Use:   "up",
		Short: "Migrate the database up",
		Run: func(cmd *cobra.Command, args []string) {
			// get the number of steps to migrate up
			steps, _ := cmd.Flags().GetInt("steps")
			dbName, _ := cmd.Flags().GetString("db")
			if dbName == "" {
				dbName = "default"
			}
			_, err := GetConnection(dbName).Migrator().Up(context.Background(), steps)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	migrateUp.Flags().IntP("steps", "s", 0, "Number of steps to migrate up")
	migrateCmd.AddCommand(migrateUp)

	migrateDown := &cobra.Command{
		Use:   "down",
		Short: "Migrate the database down",
		Run: func(cmd *cobra.Command, args []string) {
			// get the number of steps to migrate down
			steps, _ := cmd.Flags().GetInt("steps")
			dbName, _ := cmd.Flags().GetString("db")
			if dbName == "" {
				dbName = "default"
			}
			_, err := GetConnection(dbName).Migrator().Down(context.Background(), steps)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	migrateDown.Flags().IntP("steps", "s", 0, "Number of steps to migrate down")
	migrateCmd.AddCommand(migrateDown)

	migrateList := &cobra.Command{
		Use:   "list",
		Short: "List applied migrations",
		Run: func(cmd *cobra.Command, args []string) {

			dbName, _ := cmd.Flags().GetString("db")
			if dbName == "" {
				dbName = "default"
			}
			list, err := GetConnection(dbName).Migrator().Applied(context.Background())

			if len(list) == 0 {
				color.Green("No migrations applied")
				return
			}

			color.Green("List of applied migrations")
			t := table.NewWriter()
			t.AppendHeader(table.Row{"#", "Name", "Run At"})

			if err != nil {
				fmt.Println(err)
			}
			for i, v := range list {
				t.AppendRow([]interface{}{i + 1, v.Name, v.RunAt})
				// fmt.Fprintln(cmd.OutOrStdout(), i+1, v.Name, v.RunAt)
			}
			fmt.Println(t.Render())
		},
	}
	migrateCmd.AddCommand(migrateList)

	migratePending := &cobra.Command{
		Use:   "pending",
		Short: "List pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			dbName, _ := cmd.Flags().GetString("db")
			if dbName == "" {
				dbName = "default"
			}
			list, err := GetConnection(dbName).Migrator().Pending(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(list) == 0 {
				color.Green("No pending migrations")
				return
			}

			color.Green("List of pending migrations")

			t := table.NewWriter()
			t.AppendHeader(table.Row{"#", "Name"})

			for i, v := range list {
				t.AppendRow([]interface{}{i + 1, v.Name()})
				// fmt.Fprintln(cmd.OutOrStdout(), i+1, v.Name)
			}
			fmt.Println(t.Render())

		},
	}
	migrateCmd.AddCommand(migratePending)

	migrateCreate := &cobra.Command{
		Use:   "create",
		Short: "Create a new migration",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				fmt.Println("Please provide a name for the migration")
				return
			}
			color.Green("Creating migration \"%s\"", name)

			dbName, _ := cmd.Flags().GetString("db")
			if dbName == "" {
				dbName = "default"
			}

			err := GetConnection(dbName).Migrator().Create(context.Background(), name)
			if err != nil {
				fmt.Println(err)
				return
			}
			color.Green("Migration \"%s\" created", name)
		},
	}
	migrateCreate.Flags().StringP("name", "n", "", "Name of the migration")
	migrateCmd.AddCommand(migrateCreate)

	migrateReset := &cobra.Command{
		Use:   "reset",
		Short: "Reset the database",
		Run: func(cmd *cobra.Command, args []string) {
			dbName, _ := cmd.Flags().GetString("db")
			if dbName == "" {
				dbName = "default"
			}

			err := GetConnection(dbName).Migrator().Reset(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	migrateCmd.AddCommand(migrateReset)

	rootCmd.AddCommand(migrateCmd)
}

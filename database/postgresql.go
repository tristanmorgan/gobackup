package database

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
)

// PostgreSQL database
//
// type: postgresql
// host: localhost
// port: 5432
// socket:
// database: test
// username:
// password:
type PostgreSQL struct {
	Base
	host        string
	port        string
	socket      string
	database    string
	username    string
	password    string
	dumpCommand string
	args        string
}

func (ctx PostgreSQL) perform() (err error) {
	viper := ctx.viper
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 5432)
	viper.SetDefault("params", "")

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.socket = viper.GetString("socket")
	ctx.database = viper.GetString("database")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.args = viper.GetString("args")

	// socket
	if len(ctx.socket) != 0 {
		ctx.host = ""
		ctx.port = ""
	}

	if err = ctx.prepare(); err != nil {
		return
	}

	err = ctx.dump()
	return
}

func (ctx *PostgreSQL) prepare() (err error) {
	// pg_dump command
	var dumpArgs []string
	if len(ctx.database) == 0 {
		return fmt.Errorf("PostgreSQL database config is required")
	}
	if len(ctx.host) > 0 {
		dumpArgs = append(dumpArgs, "--host="+ctx.host)
	}
	if len(ctx.port) > 0 {
		dumpArgs = append(dumpArgs, "--port="+ctx.port)
	}
	if len(ctx.socket) > 0 {
		host := filepath.Dir(ctx.socket)
		port := strings.TrimPrefix(filepath.Ext(ctx.socket), ".")
		dumpArgs = append(dumpArgs, "--host="+host, "--port="+port)
	}
	if len(ctx.username) > 0 {
		dumpArgs = append(dumpArgs, "--username="+ctx.username)
	}
	if len(ctx.args) > 0 {
		dumpArgs = append(dumpArgs, ctx.args)
	}

	ctx.dumpCommand = "pg_dump " + strings.Join(dumpArgs, " ") + " " + ctx.database

	return nil
}

func (ctx *PostgreSQL) dump() error {
	logger := logger.Tag("PostgreSQL")

	dumpFilePath := path.Join(ctx.dumpPath, ctx.database+".sql")
	logger.Info("-> Dumping PostgreSQL...")
	if len(ctx.password) > 0 {
		os.Setenv("PGPASSWORD", ctx.password)
	}
	_, err := helper.Exec(ctx.dumpCommand, "-f", dumpFilePath)
	if err != nil {
		return err
	}
	logger.Info("dump path:", dumpFilePath)
	return nil
}

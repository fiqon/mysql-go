package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"log"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

type Connection struct {
	Db  *sql.DB
	Opt Options
}

func ConnectionString(options Options) string {
	config := mysqlDriver.NewConfig()

	config.User = options.User
	config.Passwd = options.Password
	config.Net = options.Net
	config.Addr = options.Host
	config.DBName = options.DB
	config.TLSConfig = "cert"

	// max allowed packet var from mysql server
	config.MaxAllowedPacket = 67108864

	return config.FormatDSN()
}

func New(options Options) (Connection, error) {
	var conn Connection
	err := conn.Connect(options)

	return conn, err
}

func (conn *Connection) CheckConnection() error {
	return conn.Db.Ping()
}

func (conn *Connection) Connect(options Options) error {
	log.Println(ConnectionString(options))

	rootCertPool, err := x509.SystemCertPool()

	if err != nil {
		log.Panic(err)
	}

	mysqlDriver.RegisterTLSConfig("cert", &tls.Config{
		RootCAs:            rootCertPool,
	})

	db, err := sql.Open("mysql", ConnectionString(options))

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(options.MaxFileTimeMinutes))
	db.SetMaxOpenConns(options.MaxOpenConns)
	db.SetMaxIdleConns(options.MaxOpenConns)

	conn.Db = db

	if err := conn.CheckConnection(); err != nil {
		return err
	}

	conn.Opt = options

	return nil
}

func (conn *Connection) Close() error {
	return conn.Db.Close()
}

func (conn *Connection) Clone(dbName string) (Connection, error) {
	var newConn Connection

	opt := conn.Opt
	opt.DB = dbName

	err := newConn.Connect(opt)

	return newConn, err
}

func (conn *Connection) ChangeDB(dbName string) error {
	newConn, err := conn.Clone(dbName)

	conn.Db = newConn.Db

	return err
}

func (conn *Connection) Reconnect(dbName string) error {
	if conn.CheckConnection() != nil {
		return conn.Connect(conn.Opt)
	}

	return nil
}

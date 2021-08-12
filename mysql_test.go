package mysql_test

import (
	"mysql"
	"os"
	"testing"
)

func TestConnection(t *testing.T) {
	opt := mysql.NewOption(os.Getenv("mysqlUser"), os.Getenv("mysqlPassword"), os.Getenv("mysqlAddress"), os.Getenv("mysqlDbName"))

	var conn mysql.Connection
	if err := conn.Connect(opt); err != nil {
		t.Fail()
	}

	// conn.Close()
	defer conn.Close()

	if err := conn.ChangeDB(os.Getenv("newMysqlDbName")); err != nil {
		t.Fail()
	}

	newConn, err := conn.Clone(os.Getenv("mysqlDbName"))

	if err != nil || newConn.CheckConnection() != nil {
		t.Fail()
	}

	// newConn.Close()
	defer newConn.Close()

	otherConn, err := mysql.New(opt)

	if err != nil || otherConn.CheckConnection() != nil {
		t.Fail()
	}

	defer otherConn.Close()
}

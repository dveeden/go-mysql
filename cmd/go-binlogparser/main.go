package main

import (
	"flag"
	"os"

	"github.com/go-mysql-org/go-mysql/replication"
)

var (
	name   = flag.String("name", "", "binlog file name")
	offset = flag.Int64("offset", 0, "parse start offset")
)

func main() {
	flag.Parse()

	p := replication.NewBinlogParser()

	f := func(e *replication.BinlogEvent) error {
		e.Dump(os.Stdout)
		return nil
	}

	err := p.ParseFile(*name, *offset, f)
	if err != nil {
		println(err.Error())
	}
}

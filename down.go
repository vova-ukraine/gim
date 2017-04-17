package main

import (
	"fmt"
	"os"

	"strconv"

	"sort"

	"github.com/urfave/cli"
	"github.com/vova-ukraine/gim/core"
)

// TODO: Implement `gim down *` ability to down throw all migrations
func downCmd(c *cli.Context) error {
	var vi int64
	var err error
	fmt.Println("Down to custom migration version")
	if !c.Args().Present() {
		fmt.Println("Migration version undefined. Use `gim down <version>`")
		os.Exit(1)
	}

	v := c.Args().Get(0)

	if v == "*" {
		vi = -1
	} else {
		vi, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Printf("Invalid version format. Error: " + err.Error())
			os.Exit(1)
		}
	}

	cfg := loadConfigHelper()
	db := initDBHelper(cfg)

	md, err := core.LoadDBMigrations(db)
	if err != nil {
		fmt.Println("Unable to load applied migration from DB: " + err.Error())
		os.Exit(1)
	}

	if _, ok := md[vi]; !ok && vi > 0 {
		fmt.Println("Unable revert down to migration version `" + v + "`. No such applied version")
		os.Exit(1)
	}

	var vs = []int64{}
	for v, _ := range md {
		vs = append(vs, v)
	}
	sort.Sort(sort.Reverse(version(vs)))

	for _, v := range vs {
		if v == vi {
			break
		}
		fmt.Print("Reverting migration with version `" + strconv.FormatInt(v, 10) + "`...")
		err = core.RevertMigration(db, md[v])
		if err != nil {
			fmt.Println("failed.")
			fmt.Println("Unable to revert migration. Error:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("ok.")
	}

	return nil
}

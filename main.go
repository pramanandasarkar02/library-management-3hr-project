package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)


func main(){


	// cluster config
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.Consistency = gocql.Quorum



	// create session
	session, err := cluster.CreateSession()
	if err != nil{
		log.Fatal("Connection err: ", err)
	}
	defer session.Close()

	// create keyspace
	err = session.Query(`
		Create keyspace if not exists test_keyspace
		with replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}
	`).Exec()

	if err != nil {
		log.Fatal(err)
	}

	session.Close()

	cluster.Keyspace = "test_keyspace"

	session, err = cluster.CreateSession()
	if err != nil{
		log.Fatal(err)
	}
	defer session.Close()

	err = session.Query(`
		Create table if not exists demo_table(
			id UUID primary key,
			name text
		)
	`).Exec()

	if err != nil{
		log.Fatal(err)
	}

	id := gocql.TimeUUID()
	err = session.Query(`INSERT INTO demo_table (id, name) values (?, ?)`, id, "Pramananda",).Exec()
	if err != nil{
		log.Fatal(err)
	}

	var name string
	err = session.Query(`select name from demo_table where id = ?`, id).Scan(&name)

	if err != nil{
		log.Fatal(err)
	}


	fmt.Println("User name: ", name)
}
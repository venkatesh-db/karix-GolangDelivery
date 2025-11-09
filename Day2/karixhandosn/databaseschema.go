
package main 

import "fmt"

type  database interface{

	 query()
}

type mysqlschema struct {

	host string
	port int
	username string
	password string
	dbname string
	
}


func (m *mysqlschema) query() {
	
	fmt.Println("mysql query executed")

}

type  postgresschema struct {
	
	host string
	port int
	username string
	password string
	dbname string


}
func (p *postgresschema) query() {
	fmt.Println("postgress query executed")
}


type mongodbschema struct {
	
	host string
	port int
	username string
	password string
	dbname string

}

func (m *mongodbschema) query() {
	
	fmt.Println("mongodb query executed")
	
}

func databaseschema() {
	
	var db database // interface object 


	db = &mysqlschema{
		host: "localhost",
		port: 3306,
		username: "root",
		password: "root",
		dbname: "karix",
	}

	db.query()
	

}

func main(){

	databaseschema()
}





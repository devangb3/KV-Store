package database_integration

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/devangb3/KV-Store/config"
	"github.com/joho/godotenv"
	
)


func main(){
	_ = godotenv.Load();

	cfg, err := config.LoadConfig();
	if err != nil{
		log.Fatalf("error loading config : %v\n", err);
	}
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName);
	db, err := sql.Open("pgx", connStr);
	if err != nil{
		fmt.Printf("Error while openning connection : %v\n", err);
		return;
	}
	defer db.Close();

	if err:= db.Ping(); err!= nil{
		log.Fatalf("Error sending ping : %v\n", err);
		return;
	}
	fmt.Println("Successfully connectied to Postgres")
	
	if err := create_table(db); err != nil{
		log.Fatalf("Error creating Table")
		return;
	} else {
		log.Println("Successfully created table");
	}

	if err := insert_record(db, "Devang", "Davis"); err != nil{
		log.Fatalf("Error inserting record : %v", err);
		return;
	}else{
		log.Println("Successfully inserted record.");
	}

	if err := list_users(db); err != nil{
		log.Fatalf("Error listing users  : %v", err);
		return;
	}else{
		log.Println("End User list");
	}
}
func create_table(db *sql.DB) error {
	sqlQuery := "CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name varchar(20), city varchar(20))"
	_,err := db.Exec(sqlQuery);
	return err;
}
func insert_record(db *sql.DB, name string, city string) error{
	sqlQuery := "INSERT INTO USERS (name, city) VALUES ($1, $2)";
	_,err := db.Exec(sqlQuery, name, city);
	return err;
}
func list_users(db *sql.DB) error{
	sqlQuery := "SELECT * FROM USERS";
	val,err := db.Query(sqlQuery);
	if err!= nil{
		return err;
	}
	defer val.Close();
	fmt.Println("Users in DB: ");
	for val.Next(){
		var id int;
		var name string ;
		var city string;
		if err:= val.Scan(&id, &name, &city); err != nil{
			return err;
		}
		fmt.Printf("ID %v Name %v, City %v\n",id, name, city)
	}
	return val.Err();
}
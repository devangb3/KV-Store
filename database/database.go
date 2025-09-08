package database

import (
	"database/sql"
	"fmt"
	
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/devangb3/KV-Store/config"
	
)
type Store struct{
	db *sql.DB
}
type User struct{
	ID int
	Name string
	City string
}

func NewStore(cfg config.Config)(*Store, error){
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName);
	db, err := sql.Open("pgx", connStr);
	if err != nil{
		return nil, fmt.Errorf("Error while openning connection : %v\n", err);
	}
	if err:= db.Ping(); err!=nil{
		return nil, fmt.Errorf("Error pinging database %v\n", err)
	}
	return &Store{db:db}, nil;
}
func(s *Store) Close(){
	s.db.Close();
}

func(s *Store) CreateUsersTable() error {
	sqlQuery := "CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name varchar(20), city varchar(20))"
	_,err := s.db.Exec(sqlQuery);
	return err;
}
func(s *Store) InsertUser(name string, city string) error{
	sqlQuery := "INSERT INTO USERS (name, city) VALUES ($1, $2)";
	_,err := s.db.Exec(sqlQuery, name, city);
	return err;
}
func(s *Store) GetUsers() ([]User, error){
	sqlQuery := "SELECT * FROM USERS";
	val,err := s.db.Query(sqlQuery);
	if err!= nil{
		return nil,err;
	}
	defer val.Close();

	var users []User
	for val.Next(){
		var u User
		if err:= val.Scan(&u.ID, &u.Name, &u.City); err != nil{
			return nil,err;
		}
		users = append(users, u);
	}
	return users,val.Err();
}
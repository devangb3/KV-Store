package database

import "log"

type Record struct{
	Key string
	Value string
}

func(s *Store) InsertRecord(key string, value string) error{
	existingValue, _ := s.GetRecord(key);
	if existingValue == ""{
		_, err := s.db.Exec("INSERT INTO store(key, value) VALUES ($1, $2)", key, value)
		return err;
	}else{
		log.Printf("Overwriting existing record")
		_,err := s.db.Exec("UPDATE store set value = $1 WHERE key = $2", value, key);
		return err;
	}
}

func(s *Store) GetRecord(key string) (string, error){
	rows, err := s.db.Query("SELECT value from store where key = $1", key);
	if err != nil{
		return "", err;
	}
	defer rows.Close();
	var ans string;
	for rows.Next(){
		var r Record
		if err:= rows.Scan(&r.Value); err != nil{
			return "", err;
		}
		ans = r.Value;
	}
	return ans,rows.Err();
}

func(s *Store) DeleteRecord(key string) error{
	_,err := s.db.Exec("DELETE from store where key = $1", key);
	return err;
}
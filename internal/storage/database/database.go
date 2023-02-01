package database

import (
	"CodeWarsCheckerAnalysis/internal/config"
	"CodeWarsCheckerAnalysis/internal/models"
	"database/sql"
	_ "github.com/lib/pq"
)

func New(cfg config.Config) (Database, error) {
	db, err := sql.Open("postgres", cfg.GetConnStr())
	if err != nil {
		return Database{}, err
	}

	err = db.Ping()
	if err != nil {
		return Database{}, err
	}

	return Database{db}, nil
}

type Database struct {
	con *sql.DB
}

func (db *Database) InitTables() error {
	queryChallenges := `CREATE TABLE challenges (
    id    varchar(30) PRIMARY KEY NOT NULL,
    rank varchar(7) NOT NULL,
    lastUpdate timestamp with time zone DEFAULT now()
);
`
	_, err := db.con.Exec(queryChallenges)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertChallenge(c models.Challenge) error {
	_, err := db.GetChallenge(c.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			_, err = db.con.Exec("INSERT INTO challenges (id, rank, lastUpdate) VALUES ($1, $2, $3)", c.Id, c.Rank, c.LastUpdate)
		}
		return err
	}

	_, err = db.con.Exec("UPDATE challenges SET rank = $1, lastUpdate = $2 WHERE id = $3", c.Rank, c.LastUpdate, c.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetChallenge(id string) (models.Challenge, error) {
	var res models.Challenge

	row := db.con.QueryRow("select * FROM challenges WHERE id = $1", id)
	if row.Err() != nil {
		return res, nil
	}

	err := row.Scan(&res.Id, &res.Rank, &res.LastUpdate)
	return res, err
}

func (db *Database) GetAllChallenges() (map[string]models.Challenge, error) {
	challenges := make(map[string]models.Challenge)
	rows, err := db.con.Query("SELECT * FROM challenges")
	if err != nil {
		return challenges, err
	}
	defer rows.Close()

	for rows.Next() {
		c := models.Challenge{}
		err := rows.Scan(&c.Id, &c.Rank, &c.LastUpdate)
		if err != nil {
			return challenges, err
		}
		challenges[c.Id] = c
	}

	return challenges, nil
}

func (db *Database) Close() error {
	return db.con.Close()
}

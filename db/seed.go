package db

import "database/sql"

func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	users := []struct {
		name      string
		email     string
		gender    string
		birthDate string
	}{
		{"Alice Johnson", "alice@example.com", "female", "1995-03-12"},
		{"Bob Smith", "bob@example.com", "male", "1990-07-22"},
		{"Carol White", "carol@example.com", "female", "1992-11-05"},
		{"David Brown", "david@example.com", "male", "1988-01-30"},
		{"Eva Martinez", "eva@example.com", "female", "1997-06-18"},
		{"Frank Lee", "frank@example.com", "male", "1993-09-25"},
		{"Grace Kim", "grace@example.com", "female", "1996-04-14"},
		{"Henry Wilson", "henry@example.com", "male", "1991-12-03"},
		{"Irene Clark", "irene@example.com", "female", "1994-08-07"},
		{"Jack Davis", "jack@example.com", "male", "1989-02-19"},
		{"Karen Turner", "karen@example.com", "female", "1998-05-28"},
		{"Liam Harris", "liam@example.com", "male", "1999-10-11"},
		{"Mia Robinson", "mia@example.com", "female", "2000-03-31"},
		{"Nathan Scott", "nathan@example.com", "male", "1987-07-16"},
		{"Olivia Adams", "olivia@example.com", "female", "1995-01-09"},
		{"Paul Nelson", "paul@example.com", "male", "1992-06-23"},
		{"Quinn Baker", "quinn@example.com", "female", "1993-11-17"},
		{"Ryan Carter", "ryan@example.com", "male", "1990-04-02"},
		{"Sophia Mitchell", "sophia@example.com", "female", "1996-08-20"},
		{"Tom Evans", "tom@example.com", "male", "1988-12-15"},
		{"Uma Foster", "uma@example.com", "female", "1997-02-27"},
		{"Victor Gray", "victor@example.com", "male", "1991-09-08"},
	}

	for _, u := range users {
		if _, err := db.Exec(
			`INSERT INTO users (name, email, gender, birth_date) VALUES ($1,$2,$3,$4)`,
			u.name, u.email, u.gender, u.birthDate,
		); err != nil {
			return err
		}
	}

	friendships := [][2]int{
		{1, 3}, {3, 1},
		{1, 4}, {4, 1},
		{1, 5}, {5, 1},
		{2, 3}, {3, 2},
		{2, 4}, {4, 2},
		{2, 5}, {5, 2},
		{1, 2}, {2, 1},
		{18, 6}, {6, 18},
		{18, 7}, {7, 18},
		{18, 8}, {8, 18},
		{19, 6}, {6, 19},
		{19, 7}, {7, 19},
		{19, 8}, {8, 19},
		{9, 10}, {10, 9},
		{11, 12}, {12, 11},
		{13, 14}, {14, 13},
		{15, 16}, {16, 15},
		{17, 20}, {20, 17},
		{21, 22}, {22, 21},
	}

	for _, f := range friendships {
		if _, err := db.Exec(
			`INSERT INTO user_friends (user_id, friend_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			f[0], f[1],
		); err != nil {
			return err
		}
	}

	return nil
}

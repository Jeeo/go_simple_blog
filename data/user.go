package data

type User struct {
	ID       int
	Username string
	Password string
}

// GetByUserSessionID Gets a user based on provided userID
func GetByUserSessionID(userID string) (user User, err error) {
	err = DB.QueryRow("Select * from Users where id = $1", userID).Scan(&user.ID, &user.Username, &user.Password)

	return user, err
}

// Get Gets a user based on provided username
func Get(username string) (user User, err error) {
	err = DB.QueryRow("Select * from Users where username = $1", username).Scan(&user.ID, &user.Username, &user.Password)

	return
}

// Create Creates a user
func (user *User) Create() (err error) {
	statement := "insert into Users (username, password) values ($1, $2) returning id, username, password"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.Username, user.Password).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return
	}

	return
}

// Update Updates a user
func (user *User) Update() (err error) {
	statement := "update users set username = $1, password = $2 where id = $3"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password, user.ID)

	return
}

// Delete Deletes a user
func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID)

	return
}

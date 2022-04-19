package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/hanna-yhchen/q-notes/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// Insert inserts a new user into database.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	statement := `INSERT INTO users(name, email, hashed_password, created)
VALUES(?, ?, ?, UTC_TIMESTAMP())`
	if _, err = m.DB.Exec(statement, name, email, string(hashedPassword)); err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Authenticate verifies whether a user exists with the provided email and password.
// Return the user ID if they do exist.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var (
		id             int
		hashedPassword []byte
	)

	statement := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	if err := m.DB.QueryRow(statement, email).Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// Get fetches details for a specific user by ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	statement := `SELECT id, name, email, created, active FROM users WHERE id = ?`
	if err := m.DB.QueryRow(statement, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}

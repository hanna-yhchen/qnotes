package models

import "time"

type Note struct {
	ID int
	Title string
	Body  string
	Created time.Time
}
package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (user_id, event_id) VALUES ($1,$2) returning id"
	err := m.DB.QueryRowContext(ctx, query, attendee.EventId, attendee.UserId).Scan(&attendee.Id)
	if err != nil {
		return nil, err
	}
	return attendee, nil
}

func (m *AttendeeModel) GetByEventAndAttendee(eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendees where event_id = $1 AND user_id = $2"

	var attendee Attendee
	err := m.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.Id, &attendee.UserId, &attendee.EventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &attendee, nil
}

func (m *AttendeeModel) GetAttendeesByEvent(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select u.id,u.name,u.email from users u join attendees a on u.id = a.user_id where a.event_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (m *AttendeeModel) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM attendees WHERE user_id = $1 AND event_id = $2"
	_, err := m.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}

	return nil

}

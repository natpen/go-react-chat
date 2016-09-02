package main

import (
	"database/sql"
	"github.com/twinj/uuid"
	"log"
	"time"
)

func GetUserLastActive(name string) (time.Time, error) {
	var lastActive time.Time
	err := db.QueryRow("SELECT last_active FROM users WHERE name = $1", name).Scan(&lastActive)

	switch {
	case err == sql.ErrNoRows:
		return time.Now(), err
	case err != nil:
		log.Fatal(err)
		return time.Now(), err
	default:
		return lastActive, nil
	}
	// TODO: row.close necessary here?
}

func AddUserOrUpdateLastActive(name string) error {
	_, err := db.Query(`
    INSERT INTO users (name, last_active)
    VALUES ($1, $2)
    ON CONFLICT (name)
    DO UPDATE SET last_active = $2
    `, name, time.Now())
	return err
}

func StoreMessage(message Message) error {
	_, err := db.Query(`
    INSERT INTO messages (id, user_id, timestamp, text)
    VALUES (
      $1,
      (SELECT id
       FROM users
       WHERE name = $2),
      $3,
      $4
    )
    `, message.Id, message.Name, message.Timestamp, message.Text)
	return err
}

func ConvertMessageRowsToMessageArray(rows *sql.Rows) []Message {
	var (
		mIdString  string
		mId        uuid.Uuid
		mTimestamp time.Time
		mName      string
		mText      string
	)
	messages := make([]Message, 0)

	for rows.Next() {
		err := rows.Scan(&mIdString, &mTimestamp, &mText, &mName)
		if err != nil {
			log.Fatal(err)
		}
		mId, _ = uuid.Parse(mIdString)
		messages = append(messages, Message{mId, "user-message", mName, mTimestamp, mText})
	}
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return messages
}

func GetMessagesForUser(name string, maxTimestamp time.Time) []Message {

	rows, err := db.Query(`
    SELECT messages.id, messages.timestamp, messages.text, users.name
    FROM messages
    INNER JOIN users
      ON messages.user_id = users.id
    WHERE messages.timestamp > LEAST(
      (SELECT MIN(timestamp)
       FROM (SELECT timestamp
       FROM messages
       WHERE timestamp < $1
       ORDER BY timestamp DESC
       LIMIT $2) AS zzz),
      (SELECT COALESCE(last_active, now())
       FROM users
       WHERE name = $3))
    ORDER BY messages.timestamp ASC
  `, maxTimestamp, messagesPerPage+1, name)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return ConvertMessageRowsToMessageArray(rows)
}

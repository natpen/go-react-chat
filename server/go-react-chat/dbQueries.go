package main

import (
	"database/sql"
	"github.com/twinj/uuid"
	"log"
	"time"
)

func GetUserLastActive(name string) (time.Time, error) {
	stmt, err := db.Prepare("SELECT last_active FROM users WHERE name = $1")
	if err != nil {
		log.Fatal(err)
	}
	var lastActive time.Time
	err = stmt.QueryRow(name).Scan(&lastActive)

	switch {
	case err == sql.ErrNoRows:
		return time.Now(), err
	case err != nil:
		log.Fatal(err)
		return time.Now(), err
	default:
		return lastActive, nil
	}
}

func AddUserOrUpdateLastActive(name string) (int64, error) {
	stmt, err := db.Prepare(`
    INSERT INTO users (name, last_active)
    VALUES ($1, $2)
    ON CONFLICT (name)
    DO UPDATE SET last_active = $2
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowCount, err
}

func StoreMessage(message Message) (int64, error) {
	stmt, err := db.Prepare(`
    INSERT INTO messages (id, user_id, timestamp, text)
    VALUES (
      $1,
      (SELECT id
       FROM users
       WHERE name = $2),
      $3,
      $4
    )
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(message.Id, message.Name, message.Timestamp, message.Text)
	if err != nil {
		log.Fatal(err)
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowCount, err
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

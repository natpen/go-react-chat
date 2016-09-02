### Local Setup

Tested versions given in parantheses - no guarantees for other versions!

 1. Get [Go](https://golang.org/) (1.7)
 2. Get [PostgreSQL](https://www.postgresql.org/) (9.5.4)
 3. Get [Node](https://nodejs.org) (6.4.0)
 4. Get [npm](https://npmjs.com) (3.10.7)
 5. Clone this repository
 6. `cd` into base project directory
 7. Create .env file with contents: `DATABASE_URL: <put your Postgres url here!>`
 8. `make init && make run`
 9. Navigate your favorite browser to http://localhost:8000

### Improvements

 * add failed join error messages to ui
 * formalize message type support (e.g., user-message, system-message, user-joined, user-left, user-started-typing, etc.)
 * error handling audit
 * optimistic state mutation on message submission
 * tests
 * figure out how Go does versioned vendor lib mgmt (see Makefile for overly-simplistic approach currently used [spoiler alert: it's just `go get ...`])
 * add db migration support
 * audit use of go compilation strategy in Makefile (`go build` vs. `go install`, etc.)
 * use compacter data structure for across-the-wire communications
 * have chatroom store the message_id of the beginning of the page for new joins (this would improve eliminate a subquery currently in use for fetching the initial batch of messages for a user, and thus make it FASTER! Ooooooh, ahhhh.)

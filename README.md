# QNotes
A simple note-taking web app built with Go.

## Demo
The sample website is deployed on Heroku:
https://qnotes-go.herokuapp.com/ <br>
Please sign up a new user to explore or just log in with the test account:
- email: test@example.com
- password: abc123

## Features
- User Authentication
- View/Create/Edit/Delete Notes (in Plain Text) with Authentication
- Form Validation
- RESTful Routes (e.g. /note/{id} for showing a specific note page)

TODO:
- [ ] Search Notes
- [ ] Markdown Support

## Technologies
Use of Go's standard library:
- HTTP server: [net/http](https://pkg.go.dev/net/http)
- HTML rendering: [html/template](https://pkg.go.dev/html/template)
- Command line configuration: [flag](https://pkg.go.dev/flag)

Use of third party package:
- Cookie-based session manager: [golangcollege/sessions](https://github.com/golangcollege/sessions)
- Routing & middleware: [go-chi/chi](https://github.com/go-chi/chi)
- MySQL driver: [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- CSRF protection: [justinas/nosurf](https://github.com/justinas/nosurf)
- Layout tool: [Bootstrap 5](https://getbootstrap.com/)
- Alert tool: [SweetAlert2](https://sweetalert2.github.io/)

Others:
- Use Javascript Fetch API to send DELETE request

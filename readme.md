1. Echo framework
2. GORM for migration and DB process
3. Golang-jwt for auth
4. robfig/cron package for managing cron jobs
5. MySQL for DB to store data
6. Validator package for validation.
7. goDotenv for env vars
8. go-mail for sending mail notifications

cmd --> api implementations are here
internal --> DBs, models, things that cannot be reused elsewhere
main.go --> entry point of go application

# Initializing Go project
go mod init github.com/Wanjie-Ryan/Go-Budget

# Getting env variables
go get github.com/joho/godotenv

# Initalizing Echo
go get github.com/labstack/echo/v4

handlers folder --> controllers
middlewares --> anything regarding middlewares
request --> houses any struct that is used to constructu reeequest coming into the application
services --> communicate with DB
common folder houses anything that we want to call globally in our app
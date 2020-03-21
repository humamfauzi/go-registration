# go-registration
Registration and User 

File structure
  main.go
  router.go
  notification.go
  auth.go
  datasource.go
  datawrite.go
  utils.go
  config
  |
    live.config.json
    staging.config.json
    local.staging.json
  external_connection
    mysql.go

# Main
Handle execution and serving http server

# Router
List all possible routing and types that needed to perform operation

# Notification
All notification utilities such as sending email and push notification handled herre

# Auth
Handle all encryption, decryption, and validation

# Datasource
Handle data source connection; can be a database, external URL or any external source

# Utils
Handle all utilty accross all major function

# Config
Where all config stored; config included but not limited to credentials, external url and reuseable data
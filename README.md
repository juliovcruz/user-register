# user-register

# How to run

Execute the commands:
```
docker build -t user-register .
docker run -p 8080:8080 user-register
``` 

Now you can access the application at the URL `http://localhost:8080`.
Swagger is on `http://localhost:8080/swagger/index.html`

# Development

Execute

```
go run cmd/api/main.go
```
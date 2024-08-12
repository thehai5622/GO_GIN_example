# RESTful API with Go and Gin framework Example
In project, CRUD with 'task' use MySQL on Docker and running on Local. Have enjoy it :)))

# Env example
Create file ".env" in Root folder project if it doesn't exist and content is below
```bash
PORT=8080
APP_ENV=local

DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=blueprint
DB_USERNAME=melkey
DB_PASSWORD=password1234
DB_ROOT_PASSWORD=password4321
```

# Database: MySQL
If you use Window and don't have "make" command. Opening the Command Prompt as Administrator and run
```bash
choco install make -y
```
On task command in Makefile, run command "Create DB container".
Connect to DB by "Navicat" or something else, and SQL file in "database" folder

# Makefile
Recommend: Use Bash terminal to avoid command structure errors
build the application
```bash
make build
```

run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

clean up binary from the last build
```bash
make clean
```

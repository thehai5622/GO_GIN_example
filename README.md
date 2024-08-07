# RESTful API with Gin framework Example
In project, me use MySQL on Docker and running on Local. Have enjoy it :)))

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

# Makefile
build the application
```bash
make build
```

run the application
```bash
make run
```

clean up binary from the last build
```bash
make clean
```

PS: If you use Window and don't have "make" command. Opening the Command Prompt as Administrator and run
```bash
choco install make -y
```
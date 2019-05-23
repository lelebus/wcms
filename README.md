# Wine Cellar Management System

Description

---

## Prerequisites

### PostgreSQL

Follow [this guide](https://www.postgresql.org/download/) for downloading the DB, by choosing your Operating System

### Yarn & Parcel-Bundler

Read [this file](ui/README.md) for a guide


## Deployment

### Configuration

- start new server for postgres and connect to it
- create postgres user
- create database
- give all privileges to the user (if it is not the owner)
- setup __config.json__ file in the application folder:
``` 
{
	"port": "APPLICATION_PORT",
	"postgreSQL": {
		"host": "CONNECTION_HOST",
		"port": "CONNECTION_PORT",
		"name": "DB_NAME",
		"user": "DB_USERNAME",
		"password": "DB_USERNAME_PASSWORD"
    }
}
```

### Run WCMS

- Lauch the application through the launch file **wcms**
- Connect to the application through a web browser:
_navigate to `locahost:APPLICATION_PORT/` if the application was lauched on the local machine_

---

## Built With

* [VueJS](https://v1.vuejs.org/)
* [Bulma](https://bulma.io/)

## API

[API Documentation]()

## Authors

- **Gabriele De Candido** - _Back-End_ - [lelebus](https://github.com/lelebus)
- **Fabio Endrizzi** - _Front-End_ - [jcte02](https://github.com/jcte02)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

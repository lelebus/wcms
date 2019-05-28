# Wine Cellar Management System

The goal of the web-application is an efficient management of a wine cellar.

The user is able to insert new wines manually, by filling in the fields with the necessary information. Alternatively, it is possible to import all wines through a CSV or EXCEL file, matching the downloadable template.

Wines are organized into catalogs, which can be either automatic or manual.
For **automatic** catalogs, all wines that match the chosen filters are added (e.g. Catalog "Italian Wines" adds all italian wines). **Customized** catalogs, instead, allow the user to add wines manually. 

One of the benefits of using the WCMS is not being forced to work on excel files. Therefore, there is also the functionality to download an automatically generated file, ready for printing. 

---

## Prerequisites

### PostgreSQL

Follow [this guide](https://www.postgresql.org/download/) to download the Database System

## Development

### Golang

Follow [this guide] (https://golang.org/doc/install) to download the Go Programming Language

### Yarn with Parcel-Bundler

Read [this file](ui/README.md) for a guide

## Deployment

### Configuration

- Start new server for postgres and connect to it
- Create postgres user
- Create database
- Give all privileges to the user (if it is not the owner)
- Setup __config.json__ file in the application folder:
  
``` json
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

- Launch the application through the launch file **wcms**
- Connect to the application through a web browser:

_navigate to `HOST:APPLICATION_PORT/`_

---

## Built With

* [Bulma](https://bulma.io/)
* [VueJS](https://v1.vuejs.org/)

## API

[API Documentation](https://documenter.getpostman.com/view/3497129/S1TSaKhe?version=latest)

## Authors

- **Gabriele De Candido** - _Back-End_ - [lelebus](https://github.com/lelebus)
- **Fabio Endrizzi** - _Front-End_ - [jcte02](https://github.com/jcte02)

## License

This project is licensed under the MIT License - [see the details]](LICENSE.md)
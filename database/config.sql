

CREATE TABLE winery (
  name text PRIMARY KEY
);

CREATE TABLE origin (
  territory text,
  region text,
  country text NOT NULL, 
  PRIMARY KEY (territory, region, country)
);

CREATE TABLE wine (
  id SERIAL PRIMARY KEY UNIQUE,
  storage_area text NOT NULL,
  type text NOT NULL,
  size float NOT NULL,
  name text NOT NULL,
  winery text NOT NULL references winery(name),
  year int NOT NULL,
  territory text,
  region text,
  country text NOT NULL,
  price money NOT NULL,
  catalogs int[],
  details text,
  internal_notes text,
  is_active boolean,
  FOREIGN KEY (territory, region, country) references origin(territory,region,country)
);

CREATE TABLE catalog (
  id SERIAL PRIMARY KEY,
  name text NOT NULL UNIQUE,
  level int NOT NULL,
  parent int references catalog(id),
  type text[],
  size float[],
  year int[],
  territory text[],
  region text[],
  country text[],
  winery text[],
  wines int[],
  is_customized boolean
);

CREATE TABLE purchase (
  id int NOT NULL PRIMARY KEY,
  wine int NOT NULL PRIMARY KEY,
  date timestamp NOT NULL PRIMARY KEY,
  supplier text,
  quantity integer NOT NULL,
  cost money NOT NULL
);


CREATE TABLE winery (
  name text PRIMARY KEY
);

CREATE TABLE origin (
  territory text,
  region text,
  nation text NOT NULL, 
  PRIMARY KEY (territory, region, nation)
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
  nation text NOT NULL,
  price money NOT NULL,
  catalog []text,
  details text,
  internal_notes text,
  is_active boolean,
  FOREIGN KEY (territory, region, nation) references origin(territory,region,nation)
);


CREATE TABLE winery (
  winery text NOT NULL PRIMARY KEY,
);


CREATE TABLE origin (
  territory text FOREIGN KEY,
  region text FOREIGN KEY,
  nation text NOT NULL FOREIGN KEY,
);

CREATE TABLE catalog (
  id int PRIMARY KEY,
  name text NOT NULL,
  level integer NOT NULL CHECK (level >= 0 AND level <= 3),
  parent int FOREIGN KEY REFERENCES catalog(id)
  type text[],
  size float[],
  year int[],
  territory text[],
  region text[],
  country text[],
  winery text[],
  storage text[],
  wine integer[],
  FOREIGN KEY (territory, region, country) REFERENCES territory(name, region, country)
);

CREATE TABLE purchase (
  id int PRIMARY KEY,
  wine int NOT NULL references wine(id),
  date timestamp NOT NULL, --'2016-06-22 19:10:25-07'
  supplier text,
  quantity int NOT NULL,
  cost money NOT NULL
);

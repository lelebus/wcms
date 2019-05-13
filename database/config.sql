

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
  id int NOT NULL PRIMARY KEY,
  name text NOT NULL UNIQUE,
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
  customized boolean
);

CREATE TABLE purchase (
  id int NOT NULL PRIMARY KEY,
  wine int NOT NULL PRIMARY KEY,
  date timestamp NOT NULL PRIMARY KEY,
  supplier text,
  quantity integer NOT NULL,
  cost money NOT NULL,
  FOREIGN KEY (wine) REFERENCES wine(id) ON DELETE RESTRICT
);
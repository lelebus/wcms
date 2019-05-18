CREATE TABLE IF NOT EXISTS wine (
  id SERIAL PRIMARY KEY UNIQUE,
  storage_area text NOT NULL,
  type text NOT NULL,
  size float NOT NULL,
  name text NOT NULL,
  winery text NOT NULL,
  year int NOT NULL,
  territory text,
  region text,
  country text NOT NULL,
  price money NOT NULL,
  catalogs int[],
  details text,
  internal_notes text,
  is_active boolean
);

CREATE TABLE IF NOT EXISTS catalog (
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

CREATE TABLE IF NOT EXISTS purchase (
  id int PRIMARY KEY,
  wine int NOT NULL references wine(id),
  date timestamp NOT NULL, --'22-06-2006'
  supplier text,
  quantity integer NOT NULL,
  cost float NOT NULL
);

INSERT INTO catalog (id, name, level, parent) VALUES (0, 'root', 0, 0);
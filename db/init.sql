CREATE TABLE IF NOT EXISTS wine (
  id SERIAL PRIMARY KEY UNIQUE,
  position int,
  storage_area text NOT NULL,
  type text NOT NULL,
  size float NOT NULL,
  name text NOT NULL,
  winery text NOT NULL,
  year int,
  territory text,
  region text,
  country text NOT NULL,
  price float NOT NULL,
  catalogs int[],
  details text,
  internal_notes text
);

CREATE TABLE IF NOT EXISTS catalog (
  id SERIAL PRIMARY KEY,
  position int,
  name text NOT NULL UNIQUE,
  level int NOT NULL,
  parent int references catalog(id) ON DELETE CASCADE,
  type text[],
  size float[],
  territory text[],
  region text[],
  country text[],
  winery text[],
  wines int[],
  is_customized boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS purchase (
  id int PRIMARY KEY,
  wine int NOT NULL references wine(id),
  date timestamp NOT NULL, --'22-06-2006'
  supplier text,
  quantity integer NOT NULL,
  cost float NOT NULL
);

DO $$ BEGIN
IF NOT EXISTS (SELECT id FROM catalog WHERE id = 0) THEN
INSERT INTO catalog (id, name, level, parent, is_customized) VALUES (0, 'root', 0, 0, true);
END IF;
END $$;
CREATE DATABASE rinhadb;

\c rinhadb

CREATE TABLE IF NOT EXISTS people (
  person_id uuid DEFAULT gen_random_uuid(), 
  -- obrigatório, string de até 100 caracteres.
  name VARCHAR(100) NOT NULL,
  -- obrigatório, único, string de até 32 caracteres.
  surname VARCHAR(32) NOT NULL,
  -- obrigatório, string para data no formato AAAA-MM-DD (ano, mês, dia).
  birthdate DATE NOT NULL,

  PRIMARY KEY (person_id)
);

-- com cada elemento sendo obrigatório e de até 32 caracteres.
CREATE TABLE IF NOT EXISTS languages (
  language_id uuid DEFAULT gen_random_uuid(),
  name VARCHAR(32) UNIQUE,

  PRIMARY KEY (language_id)
);

-- opcional, vetor de string
CREATE TABLE IF NOT EXISTS stack (
  person_id uuid references people(person_id),
  language_id uuid references languages(language_id),

  PRIMARY KEY (person_id, language_id)
);


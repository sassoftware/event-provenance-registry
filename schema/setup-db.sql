CREATE TABLE "Event" (
  "ID"  varchar  NOT NULL UNIQUE,
  "name"  varchar  NOT NULL,
  "version"  varchar  NOT NULL,
  "release"  varchar NOT NULL,
  "platformID"  varchar NOT NULL,
  "package"  varchar  NOT NULL,
  "description"  varchar NOT NULL,
  "payload"  jsonb NOT NULL,
  "event_receiver_id"  varchar NOT NULL,
  "success"  boolean  NOT NULL,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  
  CONSTRAINT "Event_pk" PRIMARY KEY ("ID")
);

CREATE TABLE "Event_receiver" (
  "ID"  varchar  NOT NULL  UNIQUE,
  "name"  varchar  NOT NULL,
  "type"  varchar  NOT NULL,
  "version"  varchar  NOT NULL,
  "description"  varchar  NOT NULL,
  "created_at" timestamp NOT NULL,
  "schema" jsonb NOT NULL,
  "fingerprint" varchar NOT NULL,
  CONSTRAINT Event_receiver_pk PRIMARY KEY ("ID"),
  UNIQUE NULLS NOT DISTINCT ("name", "type", "version")
);

CREATE TABLE "Event_receiver_group" (
  "ID"  varchar  NOT NULL  UNIQUE,
  "name"  varchar NOT NULL,
  "type" varchar NOT NULL,
  "version"  varchar NOT NULL,
  "description"  varchar NOT NULL,
  "enabled"  boolean  NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  CONSTRAINT Event_receiver_group_pk PRIMARY KEY ("ID", "type", "name", "version")
);

CREATE TABLE "Event_receiver_group_to_event_receiver" (
  "ID"  varchar  NOT NULL  UNIQUE,
  "event_receiver_group"  varchar,
  "event_receiver"  varchar,
  UNIQUE ("event_receiver_group", "event_receiver"),

  CONSTRAINT "Event_receiver_group_to_event_receiver_fk0" FOREIGN KEY ("event_receiver_group") REFERENCES "Event_receiver_group"("ID"),
  CONSTRAINT "Event_receiver_group_to_event_receiver_fk1" FOREIGN KEY ("event_receiver") REFERENCES "Event_receiver"("ID")
);

ALTER TABLE "Event" ADD CONSTRAINT "Event_fk0" FOREIGN KEY ("event_receiver_id") REFERENCES "Event_receiver"("ID");


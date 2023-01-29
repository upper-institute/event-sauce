CREATE TABLE events (
  partition_key uuid NOT NULL,
  sorting_key bigint NOT NULL,
  event_payload bytea NOT NULL,
  PRIMARY KEY (partition_key, sorting_key)
);

CREATE INDEX ix_events_sorting_key on events(sorting_key);
-- name: GetLastEvent :one
SELECT * FROM events
WHERE partition_key = $1 ORDER BY sorting_key DESC LIMIT 1;

-- name: RangeEvents :many
SELECT sorting_key, event_payload FROM events
WHERE partition_key = $1 AND sorting_key >= @start_sorting_key::bigint
ORDER BY sorting_key ASC LIMIT $2;

-- name: RangeEventsWithStop :many
SELECT sorting_key, event_payload FROM events
WHERE partition_key = $1 AND sorting_key BETWEEN @start_sorting_key::bigint AND @stop_sorting_key::bigint
ORDER BY sorting_key ASC LIMIT $2;

-- name: InsertEventSequence :exec
WITH sk(sorting_key) AS (
  SELECT
    CASE WHEN max(events.sorting_key)+1 = $2
      THEN $2
      ELSE NULL
    END
  FROM events
  WHERE events.partition_key = $1
)
INSERT INTO events (
  partition_key, sorting_key, event_payload
) VALUES (
  $1, sk.sorting_key, $3
);

-- name: UpsertEvent :exec
INSERT INTO events (
  partition_key, sorting_key, event_payload
) VALUES (
  $1, $2, $3
) ON CONFLICT (partition_key, sorting_key)
  DO UPDATE SET event_payload = $3;
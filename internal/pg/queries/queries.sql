-- name: InsertTrip :one
INSERT INTO AKT_TRIP
    ("destination", "owner_email", "owner_name", "start_at", "end_at") VALUES
    ( $1, $2, $3, $4, $5 )
RETURNING "id";

-- name: GetTrip :one
SELECT "id", "destination", "owner_email", "owner_name", "confirmed", "start_at", "end_at"
FROM AKT_TRIP
WHERE id = $1;

-- name: UpdateTrip :exec
UPDATE AKT_TRIP
SET
    "destination" = $1,
    "end_at" = $2,
    "start_at" = $3,
    "confirmed" = $4
WHERE id = $5;

-- name: GetParticipant :one
SELECT "id", "trip_id", "email", "confirmed"
FROM AKT_PARTICIPANT
WHERE id = $1;

-- name: ConfirmParticipant :exec
UPDATE AKT_PARTICIPANT
SET "confirmed" = true
WHERE id = $1;

-- name: GetParticipants :many
SELECT "id", "trip_id", "email", "confirmed"
FROM AKT_PARTICIPANT
WHERE id = $1;

-- name: InviteParticipantToTrip :one
INSERT INTO AKT_PARTICIPANT 
    ("trip_id", "email") VALUES
    ( $1, $2 )
RETURNING "id";

-- name: InviteParticipantsToTrip :copyfrom
INSERT INTO AKT_PARTICIPANT 
    ("trip_id", "email") VALUES
    ( $1, $2 );

-- name: CreateActivity :one
INSERT INTO AKT_ACTIVITY
    ("trip_id", "title", "date") VALUES
    ( $1, $2, $3 )
RETURNING "id";

-- name: GetTripActivities :many
SELECT "id", "trip_id", "title", "date"
FROM AKT_ACTIVITY
WHERE trip_id = $1;

-- name: CreateTripLink :one
INSERT INTO AKT_LINK
    ("trip_id", "title", "url") VALUES
    ( $1, $2, $3 )
RETURNING "id";

-- name: GetTripLinks :many
SELECT "id", "trip_id", "title", "url"
FROM AKT_LINK
WHERE trip_id = $1;

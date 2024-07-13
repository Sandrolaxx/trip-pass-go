package pg

import (
	"context"
	"fmt"
	"trip-pass-go/internal/api/spec"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (query *Queries) CreateTrip(
	ctx context.Context,
	pool *pgxpool.Pool,
	params spec.CreateTripRequest,
) (uuid.UUID, error) {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pg: failed to begin trx for CreateTrip: %w", err)
	}

	defer tx.Rollback(ctx) //Estranho, mas só é chamado se não chegar no commit

	qtx := query.WithTx(tx)

	tripId, err := qtx.InsertTrip(ctx, InsertTripParams{
		Destination: params.Destination,
		OwnerEmail:  string(params.OwnerEmail),
		OwnerName:   params.OwnerName,
		StartAt:     pgtype.Timestamp{Valid: true, Time: params.StartAt},
		EndAt:       pgtype.Timestamp{Valid: true, Time: params.EndAt},
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pg: failed to insert trip for CreateTrip: %w", err)
	}

	participants := make([]InviteParticipantsToTripParams, len(params.GuestsEmails))

	for i := 0; i < len(params.GuestsEmails); i++ {
		participants[i] = InviteParticipantsToTripParams{
			TripID: tripId,
			Email:  string(params.GuestsEmails[i]),
		}
	}

	if _, err := qtx.InviteParticipantsToTrip(ctx, participants); err != nil {
		return uuid.UUID{}, fmt.Errorf("pg: failed to insert participants for CreateTrip: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.UUID{}, fmt.Errorf("pg: failed to commit trx for CreateTrip: %w", err)
	}

	return tripId, nil
}

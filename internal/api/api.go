package api

import (
	"context"
	"errors"
	"net/http"
	"trip-pass-go/internal/api/spec"
	"trip-pass-go/internal/pg"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type store interface {
	GetParticipant(ctx context.Context, participantId uuid.UUID) (pg.AktParticipant, error)
	ConfirmParticipant(ctx context.Context, participantId uuid.UUID) error
}

type API struct {
	store store
	logger *zap.Logger
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger) API {
	return API{ pg.New(pool), logger }
}

func (api API) ParticipantConfirm(
	writer http.ResponseWriter,
	res http.Request,
	participantId string) spec.Response {

	id, err := uuid.Parse(participantId)
	if err != nil {
		return *spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "uuid inválido!"})
	}

	participant, err := api.store.GetParticipant(res.Context(), id)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return *spec.PatchParticipantsParticipantIDConfirmJSON400Response(
				spec.Error{Message: "participant não encontrado"},
			)
		}

		api.logger.Error("failed to get participant", zap.Error(err), zap.String("participantId", participantId))

		return *spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "something went wrong, try again"},
		)
	}

	if participant.Confirmed {
		return *spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "participant já confirmado"},
		) 
	}
	
	if err := api.store.ConfirmParticipant(res.Context(), id); err != nil {
		api.logger.Error("failed to confirm participant", zap.Error(err), zap.String("participantId", participantId))

		return *spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "something went wrong, try again"},
		)
	}

	return *spec.PatchParticipantsParticipantIDConfirmJSON204Response(nil)
}

func (api API) CreateTrip(
	writer http.ResponseWriter,
	res http.Request) spec.Response {
	panic("not implemented")
}

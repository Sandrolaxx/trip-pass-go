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
	store  store
	logger *zap.Logger
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger) API {
	return API{pg.New(pool), logger}
}

func (api API) ParticipantConfirm(
	writer http.ResponseWriter,
	req *http.Request,
	participantId string) *spec.Response {

	//print(req.Header.Get("id")) caso utilizar passando param header

	id, err := uuid.Parse(participantId)

	if err != nil {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "uuid inválido!"})
	}

	participant, err := api.store.GetParticipant(req.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
				spec.Error{Message: "participant não encontrado"},
			)
		}

		api.logger.Error("failed to get participant", zap.Error(err), zap.String("participantId", participantId))

		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "something went wrong, try again"},
		)
	}

	if participant.Confirmed {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "participant já confirmado"},
		)
	}

	if err := api.store.ConfirmParticipant(req.Context(), id); err != nil {
		api.logger.Error("failed to confirm participant", zap.Error(err), zap.String("participantId", participantId))

		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "something went wrong, try again"},
		)
	}

	return spec.PatchParticipantsParticipantIDConfirmJSON204Response(nil)
}

func (api API) CreateTrip(
	writer http.ResponseWriter,
	req *http.Request) *spec.Response {
	panic("not implemented")
}

// GetTripsTripID implements spec.ServerInterface.
func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// GetTripsTripIDActivities implements spec.ServerInterface.
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// GetTripsTripIDConfirm implements spec.ServerInterface.
func (api API) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// GetTripsTripIDLinks implements spec.ServerInterface.
func (api API) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// GetTripsTripIDParticipants implements spec.ServerInterface.
func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// PostTripsTripIDActivities implements spec.ServerInterface.
func (api API) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// PostTripsTripIDInvites implements spec.ServerInterface.
func (api API) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// PostTripsTripIDLinks implements spec.ServerInterface.
func (api API) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

// PutTripsTripID implements spec.ServerInterface.
func (api API) PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("unimplemented")
}

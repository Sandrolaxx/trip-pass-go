package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"trip-pass-go/internal/api/spec"
	"trip-pass-go/internal/pg"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type store interface {
	CreateTrip(context.Context, *pgxpool.Pool, spec.CreateTripRequest) (uuid.UUID, error)
	GetTrip(ctx context.Context, tripId uuid.UUID) (pg.AktTrip, error)

	GetParticipant(ctx context.Context, participantId uuid.UUID) (pg.AktParticipant, error)
	ConfirmParticipant(ctx context.Context, participantId uuid.UUID) error
}

type mailer interface {
	SendConfirmTripEmailToTripOwner(tripId uuid.UUID) error
}

type API struct {
	store     store
	logger    *zap.Logger
	validator *validator.Validate
	pool      *pgxpool.Pool
	mailer    mailer
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger, mailer mailer) API {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return API{pg.New(pool), logger, validate, pool, mailer}
}

func (api API) PatchParticipantsConfirm(
	writer http.ResponseWriter,
	req *http.Request,
	params spec.PatchParticipantsConfirmParams) *spec.Response {

	id, err := uuid.Parse(params.ID)

	if err != nil {
		return spec.PatchParticipantsConfirmJSON400Response(spec.Error{Message: "uuid inválido!"})
	}

	participant, err := api.store.GetParticipant(req.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PatchParticipantsConfirmJSON400Response(
				spec.Error{Message: "participant não encontrado"},
			)
		}

		api.logger.Error("failed to get participant", zap.Error(err), zap.String("participantId", params.ID))

		return spec.PatchParticipantsConfirmJSON400Response(
			spec.Error{Message: "something went wrong, try again"},
		)
	}

	if participant.Confirmed {
		return spec.PatchParticipantsConfirmJSON400Response(
			spec.Error{Message: "participant já confirmado"},
		)
	}

	if err := api.store.ConfirmParticipant(req.Context(), id); err != nil {
		api.logger.Error("failed to confirm participant", zap.Error(err), zap.String("participantId", params.ID))

		return spec.PatchParticipantsConfirmJSON400Response(
			spec.Error{Message: "something went wrong, try again"},
		)
	}

	return spec.PatchParticipantsConfirmJSON204Response(nil)
}

func (api API) PostTrips(
	writer http.ResponseWriter,
	req *http.Request) *spec.Response {
	var body spec.CreateTripRequest

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "invalid JSON"})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "invalid input: " + err.Error()})
	}

	tripId, err := api.store.CreateTrip(req.Context(), api.pool, body)

	if err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "failed to create trip, try againd"})
	}

	//Criando uma GoRotine para disparo de e-mail assincrono
	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripOwner(tripId); err != nil {
			api.logger.Error("failed to send email on PostTrips", zap.Error(err), zap.String("trip_id", tripId.String()))
		}
	}()

	return spec.PostTripsJSON201Response(spec.CreateTripResponse{TripID: tripId.String()})
}

// GetTrips implements spec.ServerInterface.
func (api API) GetTrips(writer http.ResponseWriter, req *http.Request, params spec.GetTripsParams) *spec.Response {
	id, err := uuid.Parse(params.ID)

	if err != nil {
		return spec.GetTripsJSON400Response(spec.Error{Message: "uuid inválido!"})
	}

	trip, err := api.store.GetTrip(req.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.GetTripsJSON400Response(spec.Error{Message: "trip não encontrada"})
		}

		api.logger.Error("failed to get trip", zap.Error(err), zap.String("tripId", params.ID))

		return spec.GetTripsJSON400Response(spec.Error{Message: "something went wrong, try again"})
	}

	tripResponse := spec.GetTripDetailsResponseTripObj{
		ID:          trip.ID.String(),
		Confirmed:   trip.Confirmed,
		Destination: trip.Destination,
		EndAt:       trip.EndAt.Time,
		StartAt:     trip.StartAt.Time,
	}

	return spec.GetTripsJSON200Response(spec.GetTripDetailsResponse{Trip: tripResponse})
}

// GetTripsActivities implements spec.ServerInterface.
func (api API) GetTripsActivities(w http.ResponseWriter, r *http.Request, params spec.GetTripsActivitiesParams) *spec.Response {
	panic("unimplemented")
}

// GetTripsConfirm implements spec.ServerInterface.
func (api API) GetTripsConfirm(w http.ResponseWriter, r *http.Request, params spec.GetTripsConfirmParams) *spec.Response {
	panic("unimplemented")
}

// GetTripsLinks implements spec.ServerInterface.
func (api API) GetTripsLinks(w http.ResponseWriter, r *http.Request, params spec.GetTripsLinksParams) *spec.Response {
	panic("unimplemented")
}

// GetTripsParticipants implements spec.ServerInterface.
func (api API) GetTripsParticipants(w http.ResponseWriter, r *http.Request, params spec.GetTripsParticipantsParams) *spec.Response {
	panic("unimplemented")
}

// PostTripsActivities implements spec.ServerInterface.
func (api API) PostTripsActivities(w http.ResponseWriter, r *http.Request, params spec.PostTripsActivitiesParams) *spec.Response {
	panic("unimplemented")
}

// PostTripsInvites implements spec.ServerInterface.
func (api API) PostTripsInvites(w http.ResponseWriter, r *http.Request, params spec.PostTripsInvitesParams) *spec.Response {
	panic("unimplemented")
}

// PostTripsLinks implements spec.ServerInterface.
func (api API) PostTripsLinks(w http.ResponseWriter, r *http.Request, params spec.PostTripsLinksParams) *spec.Response {
	panic("unimplemented")
}

// PutTrips implements spec.ServerInterface.
func (api API) PutTrips(w http.ResponseWriter, r *http.Request, params spec.PutTripsParams) *spec.Response {
	panic("unimplemented")
}

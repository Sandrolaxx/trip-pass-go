package mailpit

import (
	"context"
	"fmt"
	"time"
	"trip-pass-go/internal/pg"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
)

type store interface {
	GetTrip(context.Context, uuid.UUID) (pg.AktTrip, error)
}

type Mailpit struct {
	store store
}

func NewMailTrip(pool *pgxpool.Pool) Mailpit {
	return Mailpit{pg.New(pool)}
}

func (mp Mailpit) SendConfirmTripEmailToTripOwner(tripId uuid.UUID) error {
	ctx := context.Background()

	trip, err := mp.store.GetTrip(ctx, tripId)

	if err != nil {
		return fmt.Errorf("mailpit failed to GetTrip for SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg := mail.NewMsg()

	if err := msg.From("sandrolax@gmail.com"); err != nil {
		return fmt.Errorf("mailpit failed to set From Email for SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := msg.To(trip.OwnerEmail); err != nil {
		return fmt.Errorf("mailpit failed to set To Email for SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg.Subject("Confirme sua viagem")

	msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
		Olá, %s!
		
		A sua viagem para %s que começa no dia %s precisa ser confirmada.

		Clique no botão abaixo para confirmar.
	`, trip.OwnerName, trip.Destination, trip.StartAt.Time.Format(time.DateOnly)))

	client, err := mail.NewClient("mailpit", mail.WithTLSPolicy(mail.NoTLS), mail.WithPort(1025))

	if err != nil {
		return fmt.Errorf("mailpit failed to create email client for SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("mailpit failed to send email for SendConfirmTripEmailToTripOwner: %w", err)
	}

	return nil
}

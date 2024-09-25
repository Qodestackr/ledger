package v2

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/formancehq/ledger/internal/controller/system"

	ledger "github.com/formancehq/ledger/internal"

	"github.com/formancehq/go-libs/api"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func createLedger(systemController system.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		configuration := ledger.Configuration{}
		data, err := io.ReadAll(r.Body)
		if err != nil && !errors.Is(err, io.EOF) {
			api.InternalServerError(w, r, err)
			return
		}

		if len(data) > 0 {
			if err := json.Unmarshal(data, &configuration); err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
		}

		if err := systemController.CreateLedger(r.Context(), chi.URLParam(r, "ledger"), configuration); err != nil {
			switch {
			case errors.Is(err, system.ErrLedgerAlreadyExists):
				api.BadRequest(w, ErrValidation, err)
			case errors.Is(err, ledger.ErrInvalidLedgerName{}) ||
				errors.Is(err, ledger.ErrInvalidBucketName{}):
				api.BadRequest(w, ErrValidation, err)
			default:
				api.InternalServerError(w, r, err)
			}
			return
		}
		api.NoContent(w)
	}
}
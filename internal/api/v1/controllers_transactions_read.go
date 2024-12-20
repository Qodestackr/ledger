package v1

import (
	"net/http"
	"strconv"

	"github.com/formancehq/go-libs/v2/api"
	"github.com/formancehq/go-libs/v2/collectionutils"
	"github.com/formancehq/go-libs/v2/platform/postgres"
	"github.com/formancehq/ledger/internal/api/common"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/go-chi/chi/v5"
)

func readTransaction(w http.ResponseWriter, r *http.Request) {
	l := common.LedgerFromContext(r.Context())

	txId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.BadRequest(w, common.ErrValidation, err)
		return
	}

	query := ledgercontroller.NewGetTransactionQuery(int(txId))
	if collectionutils.Contains(r.URL.Query()["expand"], "volumes") {
		query = query.WithExpandVolumes()
	}
	if collectionutils.Contains(r.URL.Query()["expand"], "effectiveVolumes") {
		query = query.WithExpandEffectiveVolumes()
	}

	tx, err := l.GetTransaction(r.Context(), query)
	if err != nil {
		switch {
		case postgres.IsNotFoundError(err):
			api.NotFound(w, err)
		default:
			common.HandleCommonErrors(w, r, err)
		}
		return
	}

	api.Ok(w, mapTransactionToV1(*tx))
}

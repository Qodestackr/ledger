package v1

import (
	"fmt"
	"net/http"

	"github.com/formancehq/go-libs/v2/api"
	"github.com/formancehq/ledger/internal/api/common"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
)

func countTransactions(w http.ResponseWriter, r *http.Request) {

	options, err := getPaginatedQueryOptionsOfPITFilterWithVolumes(r)
	if err != nil {
		api.BadRequest(w, common.ErrValidation, err)
		return
	}
	options.QueryBuilder = buildGetTransactionsQuery(r)

	count, err := common.LedgerFromContext(r.Context()).
		CountTransactions(r.Context(), ledgercontroller.NewListTransactionsQuery(*options))
	if err != nil {
		common.HandleCommonErrors(w, r, err)
		return
	}

	w.Header().Set("Count", fmt.Sprint(count))
	api.NoContent(w)
}

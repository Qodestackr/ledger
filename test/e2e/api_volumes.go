//go:build it

package test_suite

import (
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/go-libs/testing/platform/pgtesting"
	. "github.com/formancehq/ledger/pkg/testserver"
	"github.com/formancehq/stack/ledger/client/models/components"
	"github.com/formancehq/stack/ledger/client/models/operations"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Transaction struct {
	Amount        int64
	Asset         string
	Source        string
	Destination   string
	EffectiveDate time.Time
}

var now = time.Now()

var _ = Context("Ledger accounts list API tests", func() {
	var (
		db  = pgtesting.UsePostgresDatabase(pgServer)
		ctx = logging.TestingContext()
	)

	testServer := NewTestServer(func() Configuration {
		return Configuration{
			PostgresConfiguration: db.GetValue().ConnectionOptions(),
			Output:                GinkgoWriter,
			Debug:                 debug,
			NatsURL:               natsServer.GetValue().ClientURL(),
		}
	})

	transactions := []Transaction{
		{Amount: 100, Asset: "USD", Source: "world", Destination: "account:user1", EffectiveDate: now.Add(-4 * time.Hour)},        //user1:100, world:-100
		{Amount: 125, Asset: "USD", Source: "world", Destination: "account:user2", EffectiveDate: now.Add(-3 * time.Hour)},        //user1:100, user2:125, world:-225
		{Amount: 75, Asset: "USD", Source: "account:user1", Destination: "account:user2", EffectiveDate: now.Add(-2 * time.Hour)}, //user1:25, user2:200, world:-200
		{Amount: 175, Asset: "USD", Source: "world", Destination: "account:user1", EffectiveDate: now.Add(-1 * time.Hour)},        //user1:200, user2:200, world:-400
		{Amount: 50, Asset: "USD", Source: "account:user2", Destination: "bank", EffectiveDate: now},                              //user1:200, user2:150, world:-400, bank:50
		{Amount: 100, Asset: "USD", Source: "account:user2", Destination: "account:user1", EffectiveDate: now.Add(1 * time.Hour)}, //user1:300, user2:50, world:-400, bank:50
		{Amount: 150, Asset: "USD", Source: "account:user1", Destination: "bank", EffectiveDate: now.Add(2 * time.Hour)},          //user1:150, user2:50, world:-400, bank:200
	}

	BeforeEach(func() {
		err := CreateLedger(ctx, testServer.GetValue(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())

		for _, transaction := range transactions {
			_, err := CreateTransaction(
				ctx,
				testServer.GetValue(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: components.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []components.V2Posting{
							{
								Amount:      big.NewInt(transaction.Amount),
								Asset:       transaction.Asset,
								Source:      transaction.Source,
								Destination: transaction.Destination,
							},
						},
						Timestamp: &transaction.EffectiveDate,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	When("Get current Volumes and Balances From origin of time till now (insertion-date)", func() {
		It("should be ok", func() {
			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					InsertionDate: pointer.For(true),
					Ledger:        "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(response.Data)).To(Equal(4))
			for _, volume := range response.Data {
				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(150)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(50)))
				}
				if volume.Account == "bank" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}
				if volume.Account == "world" {
					Expect(volume.Balance).To(Equal(big.NewInt(-400)))
				}
			}
		})
	})

	When("Get Volumes and Balances From oot til oot+2 hours (effectiveDate) ", func() {
		It("should be ok", func() {

			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					StartTime: pointer.For(now.Add(-4 * time.Hour)),
					EndTime:   pointer.For(now.Add(-2 * time.Hour)),
					Ledger:    "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(response.Data)).To(Equal(3))
			for _, volume := range response.Data {
				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(25)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}
				if volume.Account == "world" {
					Expect(volume.Balance).To(Equal(big.NewInt(-225)))
				}
			}
		})
	})

	When("Get Volumes and Balances Filter by address account", func() {
		It("should be ok", func() {
			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					InsertionDate: pointer.For(true),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "account:",
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.Data).To(HaveLen(2))
			for _, volume := range response.Data {
				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(150)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(50)))
				}
			}
		})
	})

	When("Get Volumes and Balances Filter by address account a,d and end-time now effective", func() {
		It("should be ok", func() {
			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "account:",
						},
					},
					EndTime: pointer.For(now),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.Data).To(HaveLen(2))

			for _, volume := range response.Data {
				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(150)))
				}
			}
		})
	})

	When("Get Volumes and Balances Filter by address account which doesn't exist", func() {
		It("should be ok", func() {
			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "foo:",
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.Data).To(HaveLen(0))
		})
	})

	When("Get Volumes and Balances Filter With futures dates empty", func() {
		It("should be ok", func() {
			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					StartTime: pointer.For(time.Now().Add(8 * time.Hour)),
					EndTime:   pointer.For(time.Now().Add(12 * time.Hour)),
					Ledger:    "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(response.Data)).To(Equal(0))
		})
	})

	When("Get Volumes and Balances Filter by address account aggregation by level 1", func() {
		It("should be ok", func() {
			response, err := GetVolumesWithBalances(
				ctx,
				testServer.GetValue(),
				operations.V2GetVolumesWithBalancesRequest{
					InsertionDate: pointer.For(true),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "account:",
						},
					},
					GroupBy: pointer.For(int64(1)),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(response.Data)).To(Equal(1))
			for _, volume := range response.Data {
				if volume.Account == "account" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}
			}
		})
	})
})
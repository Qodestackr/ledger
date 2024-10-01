//go:build it

package test_suite

import (
	"github.com/formancehq/go-libs/logging"
	. "github.com/formancehq/go-libs/testing/api"
	"github.com/formancehq/go-libs/testing/platform/pgtesting"
	. "github.com/formancehq/ledger/pkg/testserver"
	"github.com/formancehq/stack/ledger/client/models/components"
	"github.com/formancehq/stack/ledger/client/models/operations"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ledgerevents "github.com/formancehq/ledger/pkg/events"
)

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
	BeforeEach(func() {
		err := CreateLedger(ctx, testServer.GetValue(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
	})
	When("creating a transaction on a ledger", func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
			tx        *components.V2Transaction
			events    chan *nats.Msg
			err       error
		)
		BeforeEach(func() {
			events = testServer.GetValue().Subscribe()
			tx, err = CreateTransaction(
				ctx,
				testServer.GetValue(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: components.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []components.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
						Timestamp: &timestamp,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})
		When("transferring funds from destination to another account", func() {
			BeforeEach(func() {
				_, err := CreateTransaction(
					ctx,
					testServer.GetValue(),
					operations.V2CreateTransactionRequest{
						V2PostTransaction: components.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []components.V2Posting{
								{
									Amount:      big.NewInt(100),
									Asset:       "USD",
									Source:      "alice",
									Destination: "foo",
								},
							},
							Timestamp: &timestamp,
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
			})
			When("trying to revert the original transaction", func() {
				var (
					force bool
					err   error
				)
				revertTx := func() {
					_, err = RevertTransaction(
						ctx,
						testServer.GetValue(),
						operations.V2RevertTransactionRequest{
							Force:  pointer.For(force),
							ID:     tx.ID,
							Ledger: "default",
						},
					)
				}
				JustBeforeEach(revertTx)
				It("Should fail", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(HaveErrorCode(string(components.V2ErrorsEnumInsufficientFund)))
				})
				Context("With forcing", func() {
					BeforeEach(func() {
						force = true
					})
					It("Should be ok", func() {
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
		})
		When("reverting it", func() {
			BeforeEach(func() {
				_, err := RevertTransaction(
					ctx,
					testServer.GetValue(),
					operations.V2RevertTransactionRequest{
						Ledger: "default",
						ID:     tx.ID,
					},
				)
				Expect(err).To(Succeed())
			})
			It("should trigger a new event", func() {
				Eventually(events).Should(Receive(Event(ledgerevents.EventTypeRevertedTransaction)))
			})
			It("should revert the original transaction", func() {
				response, err := GetTransaction(
					ctx,
					testServer.GetValue(),
					operations.V2GetTransactionRequest{
						Ledger: "default",
						ID:     tx.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(response.Reverted).To(BeTrue())
			})
			When("trying to revert again", func() {
				It("should be rejected", func() {
					_, err := RevertTransaction(
						ctx,
						testServer.GetValue(),
						operations.V2RevertTransactionRequest{
							Ledger: "default",
							ID:     tx.ID,
						},
					)
					Expect(err).NotTo(BeNil())
					Expect(err).To(HaveErrorCode(string(components.V2ErrorsEnumAlreadyRevert)))
				})
			})
		})
		When("reverting it at effective date", func() {
			BeforeEach(func() {
				_, err := RevertTransaction(
					ctx,
					testServer.GetValue(),
					operations.V2RevertTransactionRequest{
						Ledger:          "default",
						ID:              tx.ID,
						AtEffectiveDate: pointer.For(true),
					},
				)
				Expect(err).To(Succeed())
			})
			It("should revert the original transaction at date of the original tx", func() {
				response, err := GetTransaction(
					ctx,
					testServer.GetValue(),
					operations.V2GetTransactionRequest{
						Ledger: "default",
						ID:     tx.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(response.Reverted).To(BeTrue())
				Expect(response.Timestamp).To(Equal(tx.Timestamp))
			})
		})
		When("reverting with dryRun", func() {
			BeforeEach(func() {
				_, err := RevertTransaction(
					ctx,
					testServer.GetValue(),
					operations.V2RevertTransactionRequest{
						Ledger: "default",
						ID:     tx.ID,
						DryRun: pointer.For(true),
					},
				)
				Expect(err).To(Succeed())
			})
			It("should not revert the transaction", func() {
				tx, err := GetTransaction(ctx, testServer.GetValue(), operations.V2GetTransactionRequest{
					Ledger: "default",
					ID:     tx.ID,
				})
				Expect(err).To(BeNil())
				Expect(tx.Reverted).To(BeFalse())
			})
		})
	})
})
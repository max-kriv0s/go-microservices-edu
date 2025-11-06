//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции sightings")

		cancel()
	})

	Describe("Get", func() {
		var partUUID string

		BeforeEach(func() {
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой запчасти в MongoDB")
		})

		It("должен успешно возвращать запчасть по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			part := resp.GetPart()

			Expect(err).ToNot(HaveOccurred())
			Expect(part).ToNot(BeNil())
			Expect(part.GetUuid()).To(Equal(partUUID))
			Expect(part.GetName()).ToNot(BeEmpty())
			Expect(part.GetDescription()).ToNot(BeEmpty())
			Expect(part.GetPrice()).To(BeNumerically(">", 0))
			Expect(part.GetDimensions()).ToNot(BeNil())
			Expect(part.GetTags()).ToNot(BeEmpty())
			Expect(part.GetMetadata()).ToNot(BeEmpty())
		})
	})

	Describe("ListParts", func() {
		var partUUID string
		BeforeEach(func() {
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой запчасти в MongoDB")
		})

		It("должен успешно возвращать список запчастей", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})

			parts := resp.GetParts()

			logger.Info(ctx, "ListParts", zap.Any("parts", parts), zap.Error(err))

			Expect(err).ToNot(HaveOccurred())
			Expect(parts).ToNot(BeNil())
			Expect(parts).To(HaveLen(1))
			Expect(parts[0].GetUuid()).To(Equal(partUUID))
		})
	})
})

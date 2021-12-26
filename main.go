package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	colRepo "github.com/nikitasadok/database-system/column/repository/mongodb"
	//	colUcase "github.com/nikitasadok/database-system/column/usecase"
	dbDelivery "github.com/nikitasadok/database-system/database/delivery/http"
	dbRepo "github.com/nikitasadok/database-system/database/repository/mongodb"
	dbUcase "github.com/nikitasadok/database-system/database/usecase"
	_ "github.com/nikitasadok/database-system/docs"
	rowDelivery "github.com/nikitasadok/database-system/row/delivery/http"
	rowRepo "github.com/nikitasadok/database-system/row/repository"
	rowUcase "github.com/nikitasadok/database-system/row/usecase"
	tDelivery "github.com/nikitasadok/database-system/table/delivery/http"
	tRepo "github.com/nikitasadok/database-system/table/repository/mongodb"
	tUcase "github.com/nikitasadok/database-system/table/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// db := domain.NewDatabaseWithName("testdb")
	// t := domain.NewTable("testtable")

	// for i := 0; i < 10; i++ {
	// 	t.AddColumn(domain.NewColumnData(fmt.Sprintf("colDatabases-%d", i+1), domain.INTEGER))
	// }

	// db.AddTable(t)

	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}

	colDatabases := conn.Database("it").Collection("databases")
	colTables := conn.Database("it").Collection("tables")
	colColumnData := conn.Database("it").Collection("colData")
	colRows := conn.Database("it").Collection("rows")

	dbRepo := dbRepo.NewMongoDatabaseRepository(colDatabases, colTables)
	dbUcase := dbUcase.NewDatabaseUsecase(dbRepo, 10*time.Second)

	colRepo := colRepo.NewMongoColumnRepository(colColumnData)
	// colUcase := colUcase.

	tableRepo := tRepo.NewMongoDatabaseRepository(colTables, colDatabases)

	rowRepo := rowRepo.NewMongoRowRepository(colTables, colRows)
	rowUcase := rowUcase.NewRowUsecase(rowRepo, colRepo, 10*time.Second)
	tableUcase := tUcase.NewTableUsecase(tableRepo, colRepo, rowRepo, 10*time.Second)

	r := gin.Default()
	dbDelivery.NewDatabaseHandler(&r.RouterGroup, dbUcase)
	tDelivery.NewTableHandler(&r.RouterGroup, tableUcase, rowUcase)
	rowDelivery.NewRowHandler(&r.RouterGroup, tableUcase, rowUcase)
	//go tDeliveryGRPC.NewTableGRPCServer(9090, tableUcase)
	time.Sleep(time.Second)



	//go clients.TablesProduct(os.Args[1], 9090)
	r.Run()

}

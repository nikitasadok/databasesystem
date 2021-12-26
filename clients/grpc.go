package clients

import (
	"context"
	"fmt"
	table "github.com/nikitasadok/database-system/table/delivery/grpc"
	"google.golang.org/grpc"
	"os"
	"time"
)

func TablesProduct(host string, port int) {
	args := os.Args[2:]

	if len(args) != 3 {
		fmt.Println("wrong number of arguments supplied, should be 3. tableID1, tableID2, dbID")
		return
	}

	id1 := args[0]
	id2 := args[1]
	dbID := args[2]

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(host + ":"+  fmt.Sprint(port), opts...)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	client := table.NewTableProductProviderClient(conn)
	res, err := client.GetTableProduct(context.Background(), &table.Filter{
		Id1:  id1,
		Id2:  id2,
		DbId: dbID,
	})
	time.Sleep(1 * time.Second)
	// fmt.Printf("%+v\n", res)
	for i, r := range res.Rows {
		fmt.Println("ROW", i)
		for _, c := range r.Cells {
			fmt.Printf("%+v\n", c)
		}
	}
}

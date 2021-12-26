package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nikitasadok/database-system/domain"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	UnimplementedTableProductProviderServer
	TUsecase domain.TableUsecase
}

func NewTableGRPCServer(port int, ucase domain.TableUsecase) {
	lis, err := net.Listen("tcp", "127.0.0.1:" + fmt.Sprint(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterTableProductProviderServer(grpcServer, &Server{
		TUsecase:                              ucase,
	})
	grpcServer.Serve(lis)
}

func (s*Server) GetTableProduct(ctx context.Context, filter *Filter) (*TableProduct, error) {
	res, err := s.TUsecase.GetTableProduct(ctx, filter.Id1, filter.Id2, filter.DbId)
	if err != nil {
		return nil, errors.New("cannot find table product: " + err.Error())
	}

	var rowsOut []*Row
	for i := range res.Rows {
		var cellsTmp []*Cell
		for j := range res.Rows[i].Cells {
			v ,_ := json.Marshal(res.Rows[i].Cells[j].GetValue())
			cellsTmp = append(cellsTmp, &Cell{
				Value:         v,
				ColData:       &ColumnData{
					Name:          res.Rows[i].Cells[j].GetColumnName(),
					Datatype:      res.Rows[i].Cells[j].GetDatatype(),
				},
			})
		}
		rowsOut = append(rowsOut, &Row{
			Cells:         cellsTmp,
		})
	}

	return &TableProduct{
		Rows: rowsOut,
	}, nil
}

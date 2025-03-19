package queries

import (
	"context"
	"errors"
	"grpc/internal/application/storage"
	"grpc/internal/domain/news"
	"grpc/internal/domain/repository"
	"grpc/internal/shared/dto"
	"grpc/pkg/mocks"
	"reflect"
	"testing"
)

func TestListHandler_Handle(t *testing.T) {
	type fields struct {
		Repo       repository.NewsRepositoryInterface
		Transactor storage.TransactorInterface
	}
	type args struct {
		ctx context.Context
		req dto.ListRequest
	}

	repo := mocks.NewRepositoryInterface(t)
	transactor := mocks.NewTransactorInterface(t)
	ctx := context.Background()

	emptyRequest := dto.ListRequest{}
	invalidRequest := dto.ListRequest{
		Page: -1,
	}

	repo.EXPECT().GetList(ctx, emptyRequest).RunAndReturn(func(ctx context.Context, request dto.ListRequest) ([]news.News, error) {
		return []news.News{}, nil
	})
	repo.EXPECT().GetList(ctx, invalidRequest).RunAndReturn(
		func(ctx context.Context, request dto.ListRequest) ([]news.News, error) {
			return []news.News{}, errors.New("invalid request")
		})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []news.News
		wantErr bool
	}{
		{
			name: "Correct request",
			fields: fields{
				Repo:       repo,
				Transactor: transactor,
			},
			args: args{
				ctx: context.Background(),
				req: emptyRequest,
			},
			want:    []news.News{},
			wantErr: false,
		},
		{
			name: "Invalid request",
			fields: fields{
				Repo:       repo,
				Transactor: transactor,
			},
			args: args{
				ctx: context.Background(),
				req: invalidRequest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := ListHandler{
				Repo:       tt.fields.Repo,
				Transactor: tt.fields.Transactor,
			}
			got, err := l.Handle(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

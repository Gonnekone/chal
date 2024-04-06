package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/Gonnekone/challenge/internal/domain/models"
	"github.com/Gonnekone/challenge/internal/lib/logger/handlers/slogdiscard"
)

func Test_fetchCarInfo(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		regNum  string
		want    models.Car
		wantErr bool
	}{
		{
			name: "X123XX150",
			ctx:    context.Background(),
			regNum: "X123XX150",
			want: models.Car {
				RegNum: "X123XX150",
				Mark:   "Lada",
				Model:  "Vesta",
				Owner: models.People {
					Name:      "John",
					Surname:   "Doe",
					Patronymic: "Smith",
				},
			},
			wantErr: false,
		},
		{
			name: "A456BC789",
			ctx:    context.Background(),
			regNum: "A456BC789",
			want: models.Car {
				RegNum: "A456BC789",
				Mark:   "Toyota",
				Model:  "Corolla",
				Year:   2015,
				Owner: models.People {
					Name:    "Alice",
					Surname: "Johnson",
				},
			},
			wantErr: false,
		},
		{
			name: "H789GF123",
			ctx:    context.Background(),
			regNum: "H789GF123",
			want: models.Car {
				RegNum: "H789GF123",
				Mark:   "BMW",
				Model:  "X5",
				Year:   2019,
				Owner: models.People {
					Name:    "Bob",
					Surname: "Brown",
					Patronymic: "Lee",
				},
			},
			wantErr: false,
		},
		{
			name: "Z456BC789",
			ctx:    context.Background(),
			regNum: "Z456BC789",
			want: models.Car {
				RegNum: "Z456BC789",
				Mark:   "Toyota",
				Model:  "Mark 2",
				Year:   1999,
				Owner: models.People{
					Name:    "Alice",
					Surname: "Johnson",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchCarInfo(tt.ctx, slogdiscard.NewDiscardLogger(), tt.regNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchCarInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchCarInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

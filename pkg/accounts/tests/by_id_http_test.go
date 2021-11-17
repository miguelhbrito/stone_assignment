package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stone_assignment/pkg/accounts"
	"github.com/stone_assignment/pkg/api/entity"
	"github.com/stone_assignment/pkg/api/response"
	"github.com/stone_assignment/pkg/mcontext"
	"github.com/stone_assignment/pkg/merrors"
	"github.com/stretchr/testify/assert"
)

func TestByIdAccountHTPP_Handler(t *testing.T) {
	tests := []struct {
		name      string
		accountId string
		manager   accounts.Account
		h         accounts.ByIdAccountHTPP
		want      http.HandlerFunc
	}{
		{
			name:      "Success",
			accountId: "1",
			manager: AccountCustomMock{
				GetByIdMock: func(mctx mcontext.Context, id string) (entity.Account, error) {
					return entity.Account{
						Id: "any_id",
					}, nil
				},
			},
			want: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(response.Account{
					Id:        "any_id",
					CreatedAt: time.Time{}.String(),
				})
			},
		},
		{
			name:      "Error to get account by Id",
			accountId: "1",
			manager: AccountCustomMock{
				GetByIdMock: func(mctx mcontext.Context, id string) (entity.Account, error) {
					return entity.Account{}, errors.New("some error")
				},
			},
			want: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				data, _ := json.Marshal(merrors.HTTPError{Msg: errors.New("some error").Error()})
				_, _ = w.Write(data)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := accounts.NewByIdAccountHTPP(tt.manager)

			r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", tt.accountId), nil)

			w := httptest.NewRecorder()

			tt.want.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.Handler()(g, r)

			assert.Equal(t, w.Code, g.Result().StatusCode, fmt.Sprintf("expected status code %v ", w.Code))

			assert.Equal(t, w.Body.String(), g.Body.String(), "body was not equal as expected")
		})
	}
}

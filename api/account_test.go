package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Clayagiffeb/Simple_Bank/db/mock"
	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/Clayagiffeb/Simple_Bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {

	// Here we generate cases to cover 100% of the cases
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) { // build the stubs for Getaccount: Expect the Getaccount() run with any context with the ID exactly 1 time
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)     //check response status
				requireBodyMatchAccount(t, recorder.Body, account) // check for Body, but the Body is represented as byte, we will write a function for this
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) { // build the stubs for Getaccount: Expect the Getaccount() run with any context with the ID exactly 1 time
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code) //check response status

			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) { // build the stubs for Getaccount: Expect the Getaccount() run with any context with the ID exactly 1 time
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code) //check response status

			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) { // build the stubs for Getaccount: Expect the Getaccount() run with any context with the ID exactly 1 time
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code) //check response status

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build the stubs
			tc.buildStubs(store)
			// start the test by creating new server
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// send request to server and responne to recoreder
			server.router.ServeHTTP(recorder, request)
			//check Response
			tc.checkResponse(t, recorder)
		}) // run in parallel

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body) // read all data from the body
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount) //parses the data intoo gotAccount
	require.NoError(t, err)
	require.Equal(t, gotAccount, account) //check for equal account
}

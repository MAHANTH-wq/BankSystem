package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/mahanth/simplebank/db/mock"
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/pb"
	"github.com/mahanth/simplebank/util"
	"github.com/mahanth/simplebank/worker"
	mockwk "github.com/mahanth/simplebank/worker/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)

	if !ok {
		return false
	}
	err := util.CheckPassword(expected.password, actualArg.CreateUserParams.HashedPassword)

	if err != nil {
		return false
	}
	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	//call the AfterCreate function to ensure it is called

	err = actualArg.AfterCreate(expected.user)

	if err != nil {
		return false
	}
	return true
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v and password %s", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg:      arg,
		password: password,
		user:     user,
	}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser()
	user.Username = strings.ToLower(user.Username)

	testCases := []struct {
		name          string
		request       *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, resp *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			request: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Password: password,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username:       user.Username,
						HashedPassword: user.HashedPassword,
						FullName:       user.FullName,
						Email:          user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, resp *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				createdUser := resp.GetUser()
				require.Equal(t, user.Username, createdUser.GetUsername())
				require.Equal(t, user.FullName, createdUser.GetFullName())
				require.Equal(t, user.Email, createdUser.GetEmail())
			},
		},
		{
			name: "UnknownError",
			request: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Password: password,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				taskDistributor.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, resp *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unknown, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)
			// Create a mock task distributor
			taskCtrl := gomock.NewController(t)
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)

			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)

			res, err := server.CreateUser(context.Background(), tc.request)
			tc.checkResponse(t, res, err)
		})
	}

}

func randomUser() (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, _ := util.HashPassword(password)

	return db.User{
		Username:       util.RandomString(6),
		FullName:       util.RandomString(6),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}, password

}

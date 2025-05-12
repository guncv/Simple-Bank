package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/guncv/Simple-Bank/db/mock"
	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	worker "github.com/guncv/Simple-Bank/worker"
	mockworker "github.com/guncv/Simple-Bank/worker/mock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.Users, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.Users
}

func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, actualArg.CreateUserParams.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.CreateUserParams.HashedPassword = actualArg.CreateUserParams.HashedPassword
	if !reflect.DeepEqual(e.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	actualArg.AfterCreate(e.user)
	return true
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.Users) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor)
		checkResponse func(t *testing.T, p *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
					AfterCreate: func(user db.Users) error {
						return nil
					},
				}

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, p *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, p)
				createdUser := p.User
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, p *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
			},
		},
		// {
		// 	name: "FailedToSendVerifyEmail",
		// 	req: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		FullName: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
		// 		store.EXPECT().
		// 			CreateUserTx(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.CreateUserTxResult{}, nil)

		// 		taskDistributor.EXPECT().
		// 			DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
		// 			Times(1).
		// 			Return(sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, p *pb.CreateUserResponse, err error) {
		// 		require.Error(t, err)
		// 		require.Nil(t, p)
		// 	},
		// },
		{
			name: "DuplicateUsername",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, p *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
			},
		},
		{
			name: "InvalidArgument",
			req: &pb.CreateUserRequest{
				Username: "",
				Password: "",
				FullName: "",
				Email:    "",
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, p *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storectrl := gomock.NewController(t)
			defer storectrl.Finish()

			taskDistributorctrl := gomock.NewController(t)
			defer taskDistributorctrl.Finish()

			store := mockdb.NewMockStore(storectrl)
			taskDistributor := mockworker.NewMockTaskDistributor(taskDistributorctrl)
			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)
			response, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, response, err)
		})
	}
}

package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/guncv/Simple-Bank/db/mock"
	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/token"
	"github.com/guncv/Simple-Bank/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqUpdateUserTxParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
}

func (e eqUpdateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, actualArg.HashedPassword.String)
	if err != nil {
		return false
	}

	if actualArg.Username != e.arg.Username ||
		actualArg.FullName != e.arg.FullName ||
		actualArg.Email != e.arg.Email ||
		!actualArg.HashedPassword.Valid {
		return false
	}
	return true
}

func (e eqUpdateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateUserTxParams(arg db.UpdateUserParams, password string) gomock.Matcher {
	return eqUpdateUserTxParamsMatcher{arg, password}
}

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	newPassword := util.RandomString(8)
	invalidEmail := "invalid-email"
	invalidFullName := "invalid-full-name"
	invalidPassword := "pass"

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, token token.Maker) context.Context
		checkResponse func(t *testing.T, p *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Password: &newPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				hashedPassword, err := util.HashPassword(newPassword)
				require.NoError(t, err)

				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newFullName,
						Valid:  true,
					},
					Email: sql.NullString{},
					HashedPassword: sql.NullString{
						String: hashedPassword,
						Valid:  true,
					},
				}

				updatedUser := user
				updatedUser.FullName = newFullName
				updatedUser.Email = newEmail
				updatedUser.HashedPassword = hashedPassword

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserTxParams(arg, newPassword)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return newContextWithBearerToken(t, token, user.Username, util.Role(user.Role), time.Minute)
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, p)
				updatedUser := p.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newFullName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		{
			name: "UpdateAnotherUser",
			req: &pb.UpdateUserRequest{
				Username: util.RandomOwner(),
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return newContextWithBearerToken(t, token, user.Username, util.Role(user.Role), time.Minute)
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
				require.Equal(t, codes.PermissionDenied, status.Code(err))
			},
		},
		{
			name: "ErrorUpdateUser",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
				Password: &newPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				hashedPassword, err := util.HashPassword(newPassword)
				require.NoError(t, err)

				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newFullName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
					HashedPassword: sql.NullString{
						String: hashedPassword,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserTxParams(arg, newPassword)).
					Times(1).
					Return(db.Users{}, errors.New("email already exists"))
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return newContextWithBearerToken(t, token, user.Username, util.Role(user.Role), time.Minute)
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				log.Printf("error: %v", err)
				require.Error(t, err)
				require.Nil(t, p)
				require.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			name: "NotFound",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Users{}, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return newContextWithBearerToken(t, token, user.Username, util.Role(user.Role), time.Minute)
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
				require.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "InvalidArgument",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &invalidFullName,
				Email:    &invalidEmail,
				Password: &invalidPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return newContextWithBearerToken(t, token, user.Username, util.Role(user.Role), time.Minute)
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
				require.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Unauthorized",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
				require.Equal(t, codes.Unauthenticated, status.Code(err))
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, token token.Maker) context.Context {
				return newContextWithBearerToken(t, token, user.Username, util.Role(user.Role), -time.Minute)
			},
			checkResponse: func(t *testing.T, p *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, p)
				require.Equal(t, codes.Unauthenticated, status.Code(err))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil)
			ctx := tc.buildContext(t, server.tokenMaker)
			response, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, response, err)
		})
	}
}

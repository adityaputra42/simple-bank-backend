package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/util"
	"simple-bank/worker"
	mockwk "simple-bank/worker/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserTxMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxMatcher) Matches(x interface{}) bool {

	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = arg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}

	err = arg.AfterCreate(expected.user)

	return err == nil
}

func (e eqCreateUserTxMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParam(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxMatcher{arg: arg, password: password, user: user}
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "Ok",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserParam(arg, password, user)).Times(1).Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).Times(1).Return(nil)

			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},

		// {
		// 	name: "InternalError",
		// 	body: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		FullName: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
		// 		store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserTxResult{}, sql.ErrConnDone)
		// 		taskDistributor.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
		// 		require.Error(t, err)

		// 	},
		// },
		// {
		// 	name: "DuplicateUsername",
		// 	body: &pb.CreateUserRequest{
		// 		Username: user.Username,
		// 		Password: password,
		// 		FullName: user.FullName,
		// 		Email:    user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserTxResult{User: db.User{}}, &pq.Error{Code: "23505"})

		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
		// 		require.Error(t, err)

		// 	},
		// },
		{
			name: "InvalidUsername",
			body: &pb.CreateUserRequest{
				Username: "invalid-user#1",
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},

			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)

			},
		},
		{
			name: "InvalidEmail",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    "invalid-user#1",
			},

			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)

			},
		},
		{
			name: "PasswordTooShort",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				Password: "123",
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)

			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDst := mockwk.NewMockTaskDistributor(taskCtrl)
			// build stub
			tc.buildStubs(store, taskDst)
			// Start test server and send request
			server := NewTestServer(t, store, taskDst)
			res, err := server.CreateUser(context.Background(), tc.body)
			// check response

			tc.checkResponse(t, res, err)
		})

	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return user, password
}

func RequireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(io.Reader(body))
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	fmt.Println("db user => ", user)
	fmt.Println("got user => ", gotUser)
	require.Equal(t, user.Email, gotUser.Email)
}

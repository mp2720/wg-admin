package db_test

// import (
// 	"context"
// 	"errors"
// 	"math"
// 	"mp2720/wg-admin/wg-admin/db"
// 	"mp2720/wg-admin/wg-admin/storage/data"
// 	"mp2720/wg-admin/wg-admin/utils"
// 	"mp2720/wg-admin/wg-admin/utils/testutils"
// 	"sync"
// 	"testing"
// )
//
// func Test_UserRepo(t *testing.T) {
// 	// var allUsers []data.User
// 	// var allUsersLock sync.Mutex
//
// 	var wg sync.WaitGroup
//
// 	userRepo := db.NewUserRepo()
//
// 	const (
// 		runners        = 128
// 		usersPerRunner = 1000
//
// 		nameMinSize = 1
// 		nameMaxSize = 2048
// 		fareMinSize = 0
// 		fareMaxSize = 2048
//
// 		addressCntMin = 0
// 		addressCntMax = math.MaxInt64
// 	)
// 	nameRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
//
// 	wg.Add(runners)
//
// 	for range runners {
// 		go func() {
// 			var users []data.User
//
// 			// add random users
// 			for range usersPerRunner {
// 				user, err := data.NewUser(
// 					testutils.RandString(nameMinSize, nameMaxSize, nameRunes),
// 					testutils.RandBool(),
// 					nil,
// 					testutils.RandString(fareMinSize, fareMinSize, nil),
// 					testutils.RandRangeInt64(addressCntMin, addressCntMax),
// 				)
// 				if err != nil {
// 					panic(err)
// 				}
//
// 				if err = userRepo.Create(context.Background(), &user); err != nil {
// 					if errors.As(err, &utils.ErrAlreadyExists{}) {
// 						// ok
// 						continue
// 					}
// 					panic(err)
// 				}
//
// 				users = append(users, user)
// 			}
//
// 			// delete with 0.5 probability
// 			var usersRem []data.User
// 			for _, user := range users {
// 				if testutils.RandBool() {
// 					// leave
// 					usersRem = append(usersRem, user)
// 					continue
// 				}
//
// 				// delete
// 				if err := userRepo.Delete(context.Background(), user.Name); err != nil {
// 					panic(err)
// 				}
// 			}
//
// 			// update with 0.5 probability
// 			for i, user := range usersRem {
// 				if testutils.RandBool() {
// 					continue
// 				}
//
// 				randClock := testutils.NewRandomClock(false)
//
// 				updatedUser, err := userRepo.Update(
// 					context.Background(),
// 					user.Name,
// 					func(ctx context.Context, user *data.User) error {
// 						err := user.Update(data.UserPatch{
// 							Name: testutils.RandValueOrNil(
// 								testutils.RandString(nameMinSize, nameMinSize, nameRunes),
// 							),
// 							IsAdmin:  testutils.RandValueOrNil(testutils.RandBool()),
// 							IsBanned: testutils.RandValueOrNil(testutils.RandBool()),
// 							PrivateKey: testutils.RandValueOrNil(
// 								testutils.MustGenerateWireguardPrivateKey(),
// 							),
// 							Fare: testutils.RandValueOrNil(
// 								testutils.RandString(fareMinSize, fareMaxSize, nil),
// 							),
// 							AddressCount: testutils.RandValueOrNil(
// 								testutils.RandRangeInt64(addressCntMin, addressCntMax),
// 							),
// 							MaxAddresses: testutils.RandValueOrNil(
// 								testutils.RandRangeInt64(addressCntMin, addressCntMax),
// 							),
// 							TokenIssuedAt: testutils.RandValueOrNil(randClock.Now()),
// 							LastSeenAt:    testutils.RandValueOrNil(randClock.Now()),
// 							PaidByTime:    testutils.RandValueOrNil(randClock.Now()),
// 						})
// 						if err != nil {
// 							return err
// 						}
//
// 						return nil
// 					},
// 				)
// 				if err != nil {
// 					if errors.Is(err, data.ErrUserExceedsAddressLimit) {
// 						// ok
// 						continue
// 					}
// 					panic(err)
// 				}
//
// 				users[i] = updatedUser
// 			}
//
//             // Get by name
//
// 			wg.Done()
// 		}()
// 	}
//
// 	wg.Wait()
// }

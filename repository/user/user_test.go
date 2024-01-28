package repository_test

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"project/delivery/models"
	"project/domain/entity"
	repository "project/repository/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_GetById(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		stub     func(sqlmock.Sqlmock)
		wantUser *entity.User
		wantErr  error
	}{
		// ...

		{
			name: "success",
			id:   1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL`

				mockSQL.ExpectQuery(expectedQuery).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
						AddRow(1, "John Doe", "john.doe@example.com"))
			},
			wantUser: &entity.User{
				Id:    1,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   2,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL`

				mockSQL.ExpectQuery(expectedQuery).
					WithArgs(2).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantUser: nil,
			wantErr:  nil,
		},
		{
			name: "error",
			id:   3,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM "users" WHERE "users"."id" = \$1 AND "users"."deleted_at" IS NULL`

				mockSQL.ExpectQuery(expectedQuery).
					WithArgs(3).
					WillReturnError(errors.New("new error"))
			},
			wantUser: nil,
			wantErr:  errors.New("new Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.stub(mockSQL)

			u := repository.NewUserRepository(gormDB)

			result, err := u.GetById(tt.id)

			assert.Equal(t, tt.wantUser, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
func Test_CheckPermission(t *testing.T) {
	type args struct {
		user *entity.User
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       bool
		wantErr    error
	}{
		{
			name: "user is blocked",
			args: args{user: &entity.User{Name: "Abhishek", Email: "abhishek@gmail.com", Phone: "8592007305", Password: "1234"}},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."phone" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
					WithArgs("8592007305").
					WillReturnRows(sqlmock.NewRows([]string{"Permission"}).AddRow(false))
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "user is not blocked",
			args: args{user: &entity.User{Name: "Abhishek", Email: "abhishek@gmail.com", Phone: "8592007305", Password: "1234"}},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."phone" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
					WithArgs("8592007305").
					WillReturnRows(sqlmock.NewRows([]string{"Permission"}).AddRow(true))
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "error in fetching",
			args: args{user: &entity.User{Name: "Abhishek", Email: "abhishek@gmail.com", Phone: "8592007305", Password: "1234"}},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."phone" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
					WithArgs("8592007305").
					WillReturnRows(sqlmock.NewRows([]string{"Permission"}).AddRow(false)).WillReturnError(errors.New("error in fetching block detail"))
			},
			want:    false,
			wantErr: errors.New("error in fetching block detail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSql, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.beforeTest(mockSql)
			u := repository.NewUserRepository(gormDB)

			got, err := u.CheckPermission(tt.args.user)
			assert.Equal(t, err, tt.wantErr)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userDatabase_isBlocked()=%v want %v", got, tt.want)
			}

		})
	}
}

func Test_GetAddressById(t *testing.T) {
	type args struct {
		addressId int
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       *entity.UserAddress
		wantErr    error
	}{
		{
			name: "success",
			args: args{addressId: 1},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_addresses" WHERE id=$1 AND "user_addresses"."deleted_at" IS NULL ORDER BY "user_addresses"."id" LIMIT 1`)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "address", "state", "country", "pin", "type"}).AddRow(1, 1, "test", "test", "test", "123456", "test"))
			},
			want: &entity.UserAddress{Id: 1, User_id: 1, Address: "test", State: "test", Country: "test", Pin: "123456", Type: "test"},

			wantErr: nil,
		},
		{
			name: "error",
			args: args{addressId: 1},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_addresses" WHERE id=$1 AND "user_addresses"."deleted_at" IS NULL ORDER BY "user_addresses"."id" LIMIT 1`)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{}).AddRow()).WillReturnError(errors.New("error fetching data"))
			},
			want:   nil,
			wantErr: errors.New("error fetching data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSql, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.beforeTest(mockSql)
			u := repository.NewUserRepository(gormDB)

			got, err := u.GetAddressById(tt.args.addressId)
			fmt.Printf("got: %+v\n", got)
			fmt.Printf("want: %+v\n", tt.want)
			assert.Equal(t, tt.wantErr, err)

			fmt.Println("Actual error:", err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.GetByAddresId()= got %v want %v", got, tt.want)
			}
		})
	}
}

func Test_GetSignUpByPhone(t *testing.T){
	type args struct{
		phone string
	}

	tests:=[]struct{
		name string
		args args
		beforeTest func(sqlmock.Sqlmock)
		want *models.Signup
		wantErr error
	}{
		{
			name:"success",
			args:args{phone: "8585124716"} ,
			beforeTest: func(s sqlmock.Sqlmock){
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "signups" WHERE "signups"."phone" = $1 ORDER BY "signups"."name" LIMIT 1`)).
				WithArgs("8585124716").
				WillReturnRows(sqlmock.NewRows([]string{"name","email","phone","password","referalcode"}).AddRow("abhi","abhishek@gmail.com","8585124716","12345",""))
			},
			want: &models.Signup{Name: "abhi",Email: "abhishek@gmail.com",Phone: "8585124716",Password: "12345",ReferalCode: ""},
			wantErr: nil,
		},
		{
			name:"error",
			args: args{phone: "8585124716"},
			beforeTest: func(s sqlmock.Sqlmock){
				s.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "signups" WHERE "signups"."phone" = $1 ORDER BY "signups"."name" LIMIT 1`)).
				WithArgs("8585124716").
				WillReturnRows(sqlmock.NewRows([]string{}).AddRow()).WillReturnError(errors.New("error fetching data"))
			},
			want: nil,
			wantErr: errors.New("error fetching data") ,
		},
	}
	for _,tt:=range tests{
		t.Run(tt.name,func(t *testing.T){
			mockDB,mockSQL,_:=sqlmock.New()
			defer mockDB.Close()
			 
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.beforeTest(mockSQL)
			u:=repository.NewUserRepository(gormDB)

			got,err:=u.GetSignupByPhone(tt.args.phone)
			fmt.Printf("got: %+v\n", got)
			fmt.Printf("want: %+v\n", tt.want)
			assert.Equal(t, tt.wantErr, err)

			fmt.Println("Actual error:", err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.GetByAddresId()= got %v want %v", got, tt.want)
			}
		})
	}
}

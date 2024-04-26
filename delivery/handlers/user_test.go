package handlers

// import (
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"project/delivery/models"
// 	mock "project/mock/mockusecase"
// 	"strings"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_UserLogin(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userUseCase := mock.NewMockUserUseCase(ctrl)
// 	cartUseCase := mock.NewMockCartUsecase(ctrl)
// 	productUseCase := mock.NewMockProductUseCase(ctrl)
// 	UserHandler := NewUserhandler(userUseCase, productUseCase, cartUseCase)

// 	tests := []struct {
// 		name           string
// 		formData       string
// 		beforeTest     func(mock.MockUserUseCase, string, string)
// 		expectedStatus int
// 		expectedJSON   string
// 	}{
// 		{
// 			name:     "login succesfull",
// 			formData: "phone=1234567890&password=mysecretpassword",
// 			beforeTest: func(muc mock.MockUserUseCase, phone, password string) {
// 				muc.EXPECT().ExecuteLoginWithPassword(phone, password).Times(1).Return(1, nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedJSON:   `{"message": "user logged in succesfully and cookie stored"}`,
// 		},
// 		{
// 			name:     "login failed",
// 			formData: "phone=1234567890&password=mysecretpassword",
// 			beforeTest: func(muc mock.MockUserUseCase, phone, password string) {
// 				muc.EXPECT().ExecuteLoginWithPassword(phone, password).Times(1).Return(0, errors.New("authentication failed"))
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedJSON:   `{"error":"authentication failed"}`,
// 		},
// 	}
// 	{
// 		for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 				fmt.Printf("Testing: %s\n", tt.name)

// 				tt.beforeTest(*userUseCase, "1234567890", "mysecretpassword")
// 				router := gin.Default()
// 				router.POST("/login", UserHandler.LoginWithPassword)

// 				req, err := http.NewRequest("POST", "/login", strings.NewReader(tt.formData))
// 				assert.NoError(t, err)
// 				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 				w := httptest.NewRecorder()
// 				router.ServeHTTP(w, req)
// 				assert.Equal(t, tt.expectedStatus, w.Code)
// 				assert.JSONEq(t, tt.expectedJSON, w.Body.String())
// 				fmt.Printf("Request Data: %s\n", tt.formData)
// 				fmt.Printf("Actual Response: %s\n", w.Body.String())
// 			})
// 		}
// 	}
// }

// func Test_UserSignUpWithOtp(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userUseCase := mock.NewMockUserUseCase(ctrl)
// 	cartUseCase := mock.NewMockCartUsecase(ctrl)
// 	productUseCase := mock.NewMockProductUseCase(ctrl)
// 	UserHandler := NewUserhandler(userUseCase, productUseCase, cartUseCase)

// 	tests := []struct {
// 		name           string
// 		input          models.Signup
// 		beforeTest     func(mock.MockUserUseCase, models.Signup)
// 		expectedStatus int
// 		expectedJSON   string
// 	}{
// 		{
// 			name: "user signup",
// 			input: models.Signup{
// 				Name:        "Abhi",
// 				Email:       "abhishek@gmail.com",
// 				Phone:       "9947117079",
// 				Password:    "12345",
// 				ReferalCode: "",
// 			},
// 			beforeTest: func(muc mock.MockUserUseCase, signup models.Signup) {
// 				muc.EXPECT().ExecuteSignupWithOtp(signup).Times(1).Return("mocked_key", nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedJSON:   `{"otp send to ": "9947117079", "key": "mocked_key"}`,
// 		},
// 	}
// 	{
// 		for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 				tt.beforeTest(*userUseCase, tt.input)
// 				router := gin.Default()
// 				router.POST("/signup", UserHandler.SignupWithOtp)

// 				req, err := http.NewRequest("POST", "/signup", strings.NewReader(`
// 					{
// 						"name":"Abhi",
// 						"email":"abhishek@gmail.com",
// 						"phone":"9947117079",
// 						"password":"12345",
// 						"referalcode":""
// 					}
// 					`))
// 				assert.NoError(t, err)
// 				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 				w := httptest.NewRecorder()
// 				router.ServeHTTP(w, req)
// 				assert.Equal(t, tt.expectedStatus, w.Code)
// 				assert.JSONEq(t, tt.expectedJSON, w.Body.String())
// 				fmt.Printf("Request Data: %s\n", tt.input)
// 				fmt.Printf("Actual Response: %s\n", w.Body.String())
// 			})
// 		}
// 	}
// }

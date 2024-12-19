// usecases/user_usecase.go
package usecases

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/errors"
	"github.com/shayja/go-template-api/internal/services"
	"github.com/shayja/go-template-api/internal/utils"
)

type UserRepository interface {
	GetUserById(id string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByMobile(mobile string) (*entities.User, error)
	ValidatePassword(passwordHash string, plainPassword string) error
	CreateUser(user *entities.User) (*entities.User, error)
	SaveOTP(otp *entities.OTP) error
	ValidateOTP(mobile string, otp string) (bool, error)
	GetOTP(mobile string) (*entities.OTP, error)
}

type UserInteractor struct {
	UserRepository UserRepository
	SMSService     *services.SMSService // Add SMSService dependency
}

func (uc *UserInteractor) GetUserById(id string) (*entities.User, error) {
	return uc.UserRepository.GetUserById(id)
}

func (uc *UserInteractor) Login(input *entities.AuthenticationInput) (string, error) {
	// Validate username and password
	if len(input.Username) < 2 {
		return "", fmt.Errorf("username is required")
	}

	if len(input.Password) < 2 {
		return "", fmt.Errorf("password is required")
	}

	// Fetch user by username
	user, err := uc.UserRepository.GetUserByUsername(input.Username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Validate password
	err = uc.UserRepository.ValidatePassword(user.Password, input.Password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT
	jwt, err := utils.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return jwt, nil
}



func (uc *UserInteractor) GetUserByUsername(username string) (*entities.User, error) {
	return uc.UserRepository.GetUserByUsername(username)
}

func (uc *UserInteractor) GetUserByMobile(mobile string) (*entities.User, error) {
	return uc.UserRepository.GetUserByMobile(mobile)
}

func (uc *UserInteractor) ValidatePassword(passwordHash string, plainPassword string) error {
	return uc.UserRepository.ValidatePassword(passwordHash, plainPassword)
}

func (uc *UserInteractor) RegisterUser(request *entities.UserRequest) (*entities.User, error) {
    user := &entities.User{
		FirstName: request.FirstName,
		LastName: request.LastName,
		Email: strings.ToLower(request.Email), 
		Username: request.Username, 
		Password: request.Password, 
		Mobile: request.Mobile,
	}
	log.Printf("user: %s", user.Password)
    return uc.UserRepository.CreateUser(user)
}

// GenerateAndSendOTP generates an OTP, saves it, and sends it to the user's mobile number
func (uc *UserInteractor) GenerateAndSendOTP(mobile string) error {
	
	// Fetch the user associated with the mobile number
	user, err_notfound := uc.UserRepository.GetUserByMobile(mobile)
	if err_notfound != nil || user == nil {
		return errors.ErrUserNotFound
	}

	// Generate a random OTP
	otpCode := GenerateOTP()

	// Set OTP expiration time (e.g., 5 minutes from now)
	expiration := time.Now().Add(5 * time.Minute)

	otp := &entities.OTP{
		UserId:    user.Id,
		Mobile:    mobile,
		OTP:       otpCode,
		Expiration: expiration,
		CreatedAt: time.Now(),
	}

	// Save the OTP in the repository
	err := uc.UserRepository.SaveOTP(otp)
	if err != nil {
		return err
	}

	log.Printf("Sending SMS to %s with OTP: %s", mobile, otp)

	// Send the OTP via SMSService
	err = uc.SMSService.SendSMS(mobile, "Your OTP is: "+otpCode)
	if err != nil {
		return err
	}

	return nil
}


// ResendOTP resends the previously generated OTP if it is still valid, or generates and sends a new one
func (uc *UserInteractor) ResendOTP(mobile string) error {
	// Retrieve the existing OTP and its expiration time
	existingOTP, err := uc.UserRepository.GetOTP(mobile)
	if err != nil {
		return err
	}

	// Check if the existing OTP is still valid
	if time.Now().Before(existingOTP.Expiration) {
		// Resend the existing OTP
		err = uc.SMSService.SendSMS(mobile, "Your OTP is: "+existingOTP.OTP)
		if err != nil {
			return err
		}
		return nil
	}

	// If the existing OTP is expired, generate a new one
	return uc.GenerateAndSendOTP(mobile)
}

// VerifyOTP validates the provided OTP for the given mobile number
func (uc *UserInteractor) VerifyOTP(mobile string, otp string) (*entities.User, error) {
	// Check if the OTP is valid
	isValid, err := uc.UserRepository.ValidateOTP(mobile, otp)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errors.ErrInvalidOTP
	}

	// Fetch the user associated with the mobile number
	user, err := uc.UserRepository.GetUserByMobile(mobile)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateOTP generates a random 6-digit OTP as a string
func GenerateOTP() string {
	const otpLength = 6
	otp := make([]byte, otpLength)

	for i := range otp {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp[i] = byte(num.Int64()) + '0'
	}

	return string(otp)
}

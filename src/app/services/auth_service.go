package services

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pablor21/goms/app/config"
	app_context "github.com/pablor21/goms/app/context"
	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/app/mappers"
	"github.com/pablor21/goms/app/models"
	admin_mailer "github.com/pablor21/goms/app/server/apps/admin/mailer"
	"github.com/pablor21/goms/pkg/auth"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/errors"
	"github.com/pablor21/goms/pkg/interactions/request"
	"github.com/pablor21/goms/pkg/logger"
	"github.com/pablor21/goms/pkg/random"
	"github.com/pablor21/goms/pkg/security"
	"github.com/pablor21/goms/pkg/storage"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, req dtos.UserLoginInput) (dtos.UserLoginResponseData, context.Context, error)
	AuthenitcateContext(ctx context.Context, user *models.User) (context.Context, error)
	GetContextUser(ctx context.Context) (dtos.UserDTO, context.Context, error)
	GetUserByID(ctx context.Context, id int64) (dtos.UserDTO, error)
	UpdateProfile(ctx context.Context, req dtos.UpdateProfileInput) (dtos.UserDTO, error)
	RequestOTP(ctx context.Context, req dtos.RequestOTPInput) (dtos.RequestOTPResponseData, error)
	LoginWithOTP(ctx context.Context, req dtos.LoginWithOTPInput) (dtos.UserLoginResponseData, context.Context, error)
	LoginWithOTPToken(ctx context.Context, token string) (dtos.UserLoginResponseData, context.Context, error)
}

var ErrUnauthorized = errors.NewAppError("Invalid credentials", http.StatusUnauthorized)
var ErrUnauthenticated = errors.NewAppError("Unauthenticated", http.StatusUnauthorized)

type UserCtxKeyType string

var UserCtxKey UserCtxKeyType = "user"

var _authService AuthService

type defaultAuthService struct {
}

func GetAuthService() AuthService {
	if _authService == nil {
		_authService = &defaultAuthService{}
	}
	return _authService
}

func (s *defaultAuthService) GetDbConnection() *database.GormConnection {
	return database.GetConnection("default").(*database.GormConnection)
}

func (s *defaultAuthService) Login(ctx context.Context, req dtos.UserLoginInput) (res dtos.UserLoginResponseData, resCtx context.Context, err error) {
	resCtx = ctx
	db, ctx, err := s.GetDbConnection().GetContextTx(ctx)
	if err != nil {
		return
	}
	defer db.Rollback()
	// repo := repositories.Use(s.GetDbConnection().Conn()).User
	// repo.WithContext(ctx).Where(repo.Email.Eq(req.Username)).Or(repo.Metadata.RawExpr()).Or(repo.PhoneNumber.Eq(req.Username)).Where(repo.Status.Eq(string(models.UserStatusActive))).Joins(repo.Avatar.TagEntries.Tag).First()

	var user *models.User
	var q *gorm.DB

	if valid.IsEmail(req.Username) {
		q = db.Where("email = ?", req.Username)
	} else if valid.IsE164(req.Username) {
		q = db.Where("phone_number = ?", req.Username)
	} else {
		// simulate slow response (to avoid the caller to know if the user exists, since the check password is slow)
		time.Sleep(1 * time.Second)
		err = ErrUnauthorized
		return
	}

	q.Where("status = ?", models.UserStatusActive).Joins("Avatar")
	err = q.First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// simulate slow response (to avoid the caller to know if the user exists, since the check password is slow)
			time.Sleep(1 * time.Second)
			err = ErrUnauthorized
		}
		return
	}

	if !user.CheckPassword(req.Password) {
		err = ErrUnauthorized
		return
	}

	res.User = mappers.GetAuthMapper().MapUserToDTO(ctx, user, true)

	resCtx, err = s.AuthenitcateContext(resCtx, user)
	if err != nil {
		return
	}

	db.Commit()
	return
}

func (s *defaultAuthService) LoginWithOTP(ctx context.Context, req dtos.LoginWithOTPInput) (res dtos.UserLoginResponseData, resCtx context.Context, err error) {
	resCtx = ctx
	db := s.GetDbConnection().Conn()

	// check if the input has token
	if req.Token != "" {
		return s.LoginWithOTPToken(ctx, req.Token)
	}

	var validOtps []*models.OTP
	err = db.Where("username = ? AND valid_until >= ?", req.Username, req.Code, time.Now().UTC()).Order("created_at desc").Find(&validOtps).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrUnauthorized
		}
		return
	}

	if len(validOtps) == 0 {
		err = ErrUnauthorized
		return
	}

	// check if the user exists
	var user *models.User
	err = db.Table((&models.User{}).TableName()+" as u").Where("u.id = ? AND u.status = ?", validOtps[0].UserID, models.UserStatusActive).Joins("Avatar").First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrUnauthorized
		}
		return
	}

	var otp *models.OTP
	for _, tryOtp := range validOtps {
		//check the token
		if tryOtp.CheckCode(req.Code) {
			otp = tryOtp
			break
		}
	}

	if otp == nil {
		err = ErrUnauthorized
		return
	}

	// delete all the otps
	_ = db.Where("username = ?", req.Username).Delete(&models.OTP{})

	res.User = mappers.GetAuthMapper().MapUserToDTO(ctx, user, true)

	resCtx, err = s.AuthenitcateContext(resCtx, user)
	if err != nil {
		return
	}

	return
}

func (s *defaultAuthService) LoginWithOTPToken(ctx context.Context, token string) (res dtos.UserLoginResponseData, resCtx context.Context, err error) {
	resCtx = ctx
	// decrypt the token
	b, err := security.Decrypt(token)
	if err != nil {
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(b), &data)
	if err != nil {
		return
	}

	username, ok := data["username"].(string)
	if !ok {
		err = errors.NewAppError("Invalid token", http.StatusBadRequest)
		return
	}

	code, ok := data["code"].(string)
	if !ok {
		err = errors.NewAppError("Invalid token", http.StatusBadRequest)
		return
	}

	return s.LoginWithOTP(ctx, dtos.LoginWithOTPInput{
		Username: username,
		Code:     code,
	})
}

func (s *defaultAuthService) AuthenitcateContext(ctx context.Context, user *models.User) (resCtx context.Context, err error) {
	resCtx = ctx
	if user == nil {
		return resCtx, ErrUnauthenticated
	}

	// create a new principal
	principal := auth.NewPrincipal()
	principal.SetId(user.ID).SetRoles([]string{string(user.Role)}).SetClaims(map[string]interface{}{
		auth.ClaimEmail:     user.Email,
		auth.ClaimFirstName: user.FirstName,
		auth.ClaimLastName:  user.LastName,
		auth.ClaimLang:      user.Lang,
		auth.ClaimStatus:    string(user.Status),
	})
	resCtx = auth.SetContextPrincipal(resCtx, principal)
	return resCtx, nil
}

func (s *defaultAuthService) RequestOTP(ctx context.Context, req dtos.RequestOTPInput) (res dtos.RequestOTPResponseData, err error) {

	if !valid.IsEmail(req.Username) && !valid.IsE164(req.Username) {
		err = errors.ErrBadRequest
		return
	}

	// check if the user has an OTP
	var existingOTP *models.OTP
	err = s.GetDbConnection().Conn().Where("username = ? AND valid_until >= ?", req.Username, time.Now().UTC()).Order("created_at desc").First(&existingOTP).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
		err = nil
		existingOTP = nil
	}

	// delete expired otps
	_ = s.GetDbConnection().Conn().Where("valid_until < ?", time.Now().UTC()).Delete(&models.OTP{}).Error

	// retry period
	retryIn := config.GetConfig().OTP.ResendDelay

	if existingOTP != nil {
		timeLimit := existingOTP.CreatedAt.Add(time.Duration(retryIn) * time.Second)
		// check if the otp has been created before retry time
		if timeLimit.After(time.Now().UTC()) {
			err = errors.NewAppError("Too many tries", http.StatusTooManyRequests).SetDetails(map[string]interface{}{
				// "retryAfterSeconds": time.Now().UTC().Sub(otp.CreatedAt.Add(time.Duration(retryIn) * time.Second)).Seconds(),
				"retryAfterSeconds": int(timeLimit.Sub(time.Now().UTC()).Seconds()),
			})
			return
		}
	}

	db, ctx, err := s.GetDbConnection().GetContextTx(ctx)
	if err != nil {
		return
	}
	defer db.Rollback()

	err = validation.ValidateStructWithContext(ctx, &req, validation.Field(&req.Username, validation.Required))
	if err != nil {
		err = errors.NewValidationError(err)
		return
	}

	err = validation.ValidateStructWithContext(ctx, &req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Client, validation.Required, validation.In(dtos.ClientTypeAdmin, dtos.ClientTypeFrontend)),
	)

	if err != nil {
		err = errors.NewValidationError(err)
		return
	}

	var u *models.User
	var q *gorm.DB
	if valid.IsEmail(req.Username) {
		q = db.Where("email = ?", req.Username)
	} else if valid.IsE164(req.Username) {
		q = db.Where("phone_number = ?", req.Username)
	} else {
		// simulate slow response (to avoid the caller to know if the user exists, since the check password is slow)
		// time.Sleep(1 * time.Second)
		err = ErrUnauthorized
		return
	}
	err = q.Where("status = ?", models.UserStatusActive).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrUnauthorized
		}
		return
	}

	var otp *models.OTP
	expires := time.Now().UTC().Add(time.Duration(config.GetConfig().OTP.Lifetime) * time.Second)
	// create a new OTP
	code := random.GenerateRandomNumber(config.GetConfig().OTP.Length)

	otp = &models.OTP{
		Username:      req.Username,
		UserID:        u.ID,
		MaxAttempts:   config.GetConfig().OTP.MaxAttempts,
		AttemptsCount: 0,
		ValidUntil:    &expires,
	}

	err = otp.SetCode(code)
	if err != nil {
		return
	}

	err = db.Create(&otp).Error
	if err != nil {
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"username": otp.Username,
		"code":     code,
	})
	if err != nil {
		return
	}

	token, err := security.Encrypt(string(b))
	if err != nil {
		return
	}
	// retryIn = otp.CreatedAt.Add(time.Duration(config.GetConfig().OTP.Lifetime) * time.Second).Sub(time.Now().UTC()).Seconds()
	// create the email data
	mailData := dtos.OTPResultMailData{
		User:              mappers.GetAuthMapper().MapUserToDTO(ctx, u, true),
		Username:          req.Username,
		Expiration:        *otp.ValidUntil,
		ExpirationMinutes: int(time.Duration(config.GetConfig().OTP.Lifetime) * time.Second / time.Minute),
		// RetryIn: retryIn,
		Token: token,
		Code:  code,
		Link:  "https://" + app_context.GetServerContext(ctx).Request().Host + app_context.GetUrlGenerator(ctx)("api.auth.otp-login-token") + "?token=" + token,
	}

	// send otp via email
	switch req.Client {
	case dtos.ClientTypeFrontend:
		// send email
		break
	case dtos.ClientTypeAdmin:
		go func() {
			err := admin_mailer.SendOtp(context.TODO(), mailData)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to send OTP email")
			}
		}() // send email in background
	}
	err = db.Commit()

	res.Username = req.Username
	res.ResendAfter = int64(retryIn)

	return
}

func (s *defaultAuthService) GetContextUser(ctx context.Context) (user dtos.UserDTO, resCtx context.Context, err error) {
	resCtx = ctx
	if user, ok := ctx.Value(UserCtxKey).(dtos.UserDTO); ok {
		return user, resCtx, nil
	}

	// check if user is authenticated
	if principal := auth.GetContextPrincipal(ctx); principal != nil {
		id := principal.GetID().(int64)
		if id > 0 {
			user, err := s.GetUserByID(ctx, id)
			resCtx = context.WithValue(ctx, UserCtxKey, user)
			return user, resCtx, err
		} else {
			return user, resCtx, ErrUnauthenticated
		}
	}
	return user, resCtx, ErrUnauthenticated
}

func (s *defaultAuthService) GetUserByID(ctx context.Context, id int64) (user dtos.UserDTO, err error) {
	db, ctx, err := s.GetDbConnection().GetContextTx(ctx)
	if err != nil {
		return
	}
	defer db.Rollback()

	var u *models.User
	err = db.Where("users.id = ?", id).Joins("Avatar").First(&u).Error
	if err != nil {
		return user, errors.NewAppError("User not found", http.StatusNotFound)
	}

	user = mappers.GetAuthMapper().MapUserToDTO(ctx, u, true)
	return
}

func (s *defaultAuthService) UpdateProfile(ctx context.Context, req dtos.UpdateProfileInput) (user dtos.UserDTO, err error) {
	db, ctx, err := s.GetDbConnection().GetContextTx(ctx)
	if err != nil {
		return
	}
	defer db.Rollback()

	principal := auth.GetContextPrincipal(ctx)
	if principal == nil {
		return user, ErrUnauthenticated
	}

	id := principal.GetID().(int64)
	if id <= 0 {
		return user, ErrUnauthenticated
	}

	var u *models.User
	err = db.Where("users.id = ?", id).Joins("Avatar").First(&u).Error
	if err != nil {
		return user, errors.NewAppError("User not found", http.StatusNotFound)
	}

	err = validation.ValidateStructWithContext(ctx, &req,
		// validation.Field(&req.FirstName, validation.Required, validation.Length(2, 50)),
		// validation.Field(&req.LastName, validation.Required, validation.Length(2, 50)),
		validation.Field(&req.Avatar, validation.By(func(value interface{}) error {
			if value == nil || value.(request.FileInput).FileHeader == nil {
				return nil
			}

			//valid mime types
			validMimeTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
			if !valid.IsIn(value.(request.FileInput).FileHeader.Header.Get("Content-Type"), validMimeTypes...) {
				return &errors.FieldErrorDetail{
					Code:    "invalid_mime_type",
					Message: "Invalid mime type",
					Params: map[string]interface{}{
						"valid_mimes": validMimeTypes,
						"mime_type":   value.(request.FileInput).FileHeader.Header.Get("Content-Type"),
					},
				}
			}

			return nil
		})),
	)
	if err != nil {
		err = errors.NewValidationError(err)
		return
	}

	//TODO: validate input
	updates := map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"lang":       req.Lang,
	}

	if req.Avatar.FileHeader != nil {
		u.Avatar, err = saveAvatar(ctx, u, req.Avatar)
		if err != nil {
			return
		}
		updates["avatar_asset_id"] = u.Avatar.ID
	}

	err = db.Model(&u).Omit("Avatar").Updates(updates).Error
	if err != nil {
		return
	}

	user = mappers.GetAuthMapper().MapUserToDTO(ctx, u, true)
	db.Commit()
	return
}

func saveAvatar(ctx context.Context, user *models.User, file request.FileInput) (res *models.Asset, err error) {
	// var assetId int64 = 0
	if user.Avatar != nil {
		// assetId = user.Avatar.ID
		// delete old avatar
		err = GetAssetService()._DeleteAsset(ctx, user.Avatar)
		if err != nil {
			return
		}
	}
	return GetAssetService()._SaveAsset(ctx, dtos.AssetInput{
		// ID:          assetId,
		OwnerId:     user.ID,
		OwnerType:   "users",
		Section:     "user_avatars",
		AssetType:   models.AssetTypeImage,
		StorageName: storage.StorageNameDefault,
		File:        &file,
	})
}

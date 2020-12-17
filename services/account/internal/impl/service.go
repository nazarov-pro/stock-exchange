package impl

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/protobuf/proto"
	"github.com/nazarov-pro/stock-exchange/pkg/util/crypt"
	"github.com/nazarov-pro/stock-exchange/pkg/util/gen"
	"github.com/nazarov-pro/stock-exchange/pkg/util/time"
	"github.com/nazarov-pro/stock-exchange/services/account/domain"
	"github.com/nazarov-pro/stock-exchange/services/account/domain/pb"
	"github.com/nazarov-pro/stock-exchange/services/account/internal/kafka"
	"github.com/nazarov-pro/stock-exchange/services/account/internal/config"
)

// SimpleService simple server(default one)
type SimpleService struct {
	repository domain.Repository
	logger     log.Logger
}

// New creating a new instance of SimpleService
func New(repo domain.Repository, logger log.Logger) domain.Service {
	kafka.CreateTopics()
	return &SimpleService{
		repository: repo,
		logger:     logger,
	}
}

// Register registration opf the account
func (svc *SimpleService) Register(ctx context.Context, req *domain.RegisterAccountRequest) (*domain.Account, error) {
	logger := log.With(svc.logger, "method", "RegisterAccount")
	level.Info(logger).Log("username", req.Username, "email", req.Email)
	acc, err := svc.repository.FindByUsernameOrEmail(ctx, req.Username, req.Email)

	switch err {
	case domain.ErrAccountNotFound: // account not found is acceptable
		break
	case nil:
		if acc.Username == req.Username {
			err = domain.ErrExistedUsername
		} else if acc.Email == req.Email {
			err = domain.ErrExistedEmail
		} else {
			level.Error(logger).Log("err", err, "msg", "this type of error is not expected")
		}

		level.Error(logger).Log("err", err)
		return nil, err
	default:
		level.Error(logger).Log("err", err)
		return nil, err
	}

	pwd, err := crypt.HashPassword(req.Password)

	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	acc = &domain.Account{
		Username:       req.Username,
		Password:       pwd,
		Email:          req.Email,
		Status:         domain.Registered,
		ActivationCode: gen.NewUUID(),
		CreatedDate:    time.Epoch(),
	}
	if err := svc.repository.Create(ctx, acc); err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}
	data, err := proto.Marshal(generateActivationMessage(acc))
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	go kafka.SendEmail(acc.Email, data)

	return acc, nil
}

// Activate activate account by email and activationCode
func (svc *SimpleService) Activate(ctx context.Context, req *domain.ActivateAccountRequest) error {
	logger := log.With(svc.logger, "method", "ActivateAccount")
	level.Info(logger).Log("email", req.Email, "activationCode", req.ActivationCode)
	err := svc.repository.UpdateStatus(ctx, req.Email, req.ActivationCode, domain.Registered, domain.Activated)
	if err == nil {
		data, err := proto.Marshal(generateSuccessfullyActivatedMessage(req.Email))
		if err != nil {
			level.Error(logger).Log("err", err)
			return err
		}

		go kafka.SendEmail(req.Email, data)
		return nil
	}
	return err
}

func generateActivationMessage(acc *domain.Account) *pb.SendEmail {
	activationLink := fmt.Sprintf(
		"%s/accounts/activate?email=%s&activationCode=%s", 
		config.Config.GetString("app.activationLinkBaseUri"), acc.Email, acc.ActivationCode,
	)
	return &pb.SendEmail{
		Recipients: []string{acc.Email},
		Subject:    "Email Activation",
		Content: fmt.Sprintf("Hello %s, Welcome to %s. You can activate your account by followint the link below %s",
			acc.Username, config.Config.GetString("app.productName"), activationLink),
	}
}

func generateSuccessfullyActivatedMessage(email string) *pb.SendEmail {
	return &pb.SendEmail{
		Recipients: []string{email},
		Subject:    "Account Activated",
		Content: fmt.Sprintf("Hello %s, Welcome to %s. Your account successfully activated.",
			email, config.Config.GetString("app.productName")),
	}
}

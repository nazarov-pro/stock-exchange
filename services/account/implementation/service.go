package implementation

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/protobuf/proto"
	"github.com/nazarov-pro/stock-exchange/pkg/util/crypt"
	"github.com/nazarov-pro/stock-exchange/pkg/util/gen"
	"github.com/nazarov-pro/stock-exchange/pkg/util/time"
	"github.com/nazarov-pro/stock-exchange/services/account"
	"github.com/nazarov-pro/stock-exchange/services/account/pb"
	"github.com/nazarov-pro/stock-exchange/services/account/transport/kafka"
)

// SimpleService simple server(default one)
type SimpleService struct {
	repository account.Repository
	logger     log.Logger
}

// New creating a new instance of SimpleService
func New(repo account.Repository, logger log.Logger) account.Service {
	kafka.CreateTopics()
	return &SimpleService{
		repository: repo,
		logger:     logger,
	}
}

// Register registration opf the account
func (svc *SimpleService) Register(ctx context.Context, req *account.RegisterAccountRequest) (*account.Account, error) {
	logger := log.With(svc.logger, "method", "RegisterAccount")
	level.Info(logger).Log("username", req.Username, "email", req.Email)
	acc, err := svc.repository.FindByUsernameOrEmail(ctx, req.Username, req.Email)

	switch err {
	case account.ErrAccountNotFound: // account not found is acceptable
		break
	case nil:
		if acc.Username == req.Username {
			err = account.ErrExistedUsername
		} else if acc.Email == req.Email {
			err = account.ErrExistedEmail
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

	acc = &account.Account{
		Username:       req.Username,
		Password:       pwd,
		Email:          req.Email,
		Status:         account.Registered,
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

	err = kafka.SendEmail(acc.Email, data)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	return acc, nil
}

// Activate activate account by email and activationCode
func (svc *SimpleService) Activate(ctx context.Context, req *account.ActivateAccountRequest) error {
	logger := log.With(svc.logger, "method", "ActivateAccount")
	level.Info(logger).Log("email", req.Email, "activationCode", req.ActivationCode)
	err := svc.repository.UpdateStatus(ctx, req.Email, req.ActivationCode, account.Registered, account.Activated)
	if err == nil {
		data, err := proto.Marshal(generateSuccessfullyActivatedMessage(req.Email))
		if err != nil {
			level.Error(logger).Log("err", err)
			return err
		}

		return kafka.SendEmail(req.Email, data)
	}
	return err
}

func generateActivationMessage(acc *account.Account) *pb.SendEmail {
	activationLink := fmt.Sprintf("http://localhost:8080/accounts/activate?email=%s&activationCode=%s", acc.Email, acc.ActivationCode)
	return &pb.SendEmail{
		Recipients: []string{acc.Email},
		Subject:    "Email Activation",
		Content: fmt.Sprintf("Hello %s, Welcome to %s. You can activate your account by followint the link below %s",
			acc.Username, "Stock Exchange", activationLink),
	}
}

func generateSuccessfullyActivatedMessage(email string) *pb.SendEmail {
	return &pb.SendEmail{
		Recipients: []string{email},
		Subject:    "Account Activated",
		Content: fmt.Sprintf("Hello %s, Welcome to %s. Your account successfully activated.",
			email, "Stock Exchange"),
	}
}

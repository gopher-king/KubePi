package user

import (
	"errors"
	v1User "github.com/KubeOperator/ekko/internal/model/v1/user"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service interface {
	common.DBService
	Create(u *v1User.User, options common.DBOptions) error
	Get(name string, options common.DBOptions) (*v1User.User, error)
	GetByEmail(email string, options common.DBOptions) (*v1User.User, error)
	List(options common.DBOptions) ([]v1User.User, error)
	Delete(name string, options common.DBOptions) error
	Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1User.User, int, error)
	Update(name string, u *v1User.User, options common.DBOptions) error
}

func NewService() Service {
	return &service{
	}
}

type service struct {
	common.DefaultDBService
}

func (u *service) Update(name string, us *v1User.User, options common.DBOptions) error {
	cu, err := u.Get(name, options)
	if err != nil {
		return err
	}
	if cu.CreatedBy == "system" {
		return errors.New("can not delete this resource,because it created by system")
	}
	db := u.GetDB(options)
	us.UUID = cu.UUID
	us.CreateAt = cu.CreateAt
	us.UpdateAt = time.Now()
	return db.Update(us)
}

func (u *service) Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1User.User, int, error) {
	db := u.GetDB(options)
	query := db.Select()
	if num != 0 && size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	count, err := query.Count(&v1User.User{})
	if err != nil {
		return nil, 0, err
	}
	users := make([]v1User.User, 0)
	if err := query.Find(&users); err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *service) Get(name string, options common.DBOptions) (*v1User.User, error) {
	db := u.GetDB(options)
	var us v1User.User
	if err := db.One("Name", name, &us); err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *service) GetByEmail(email string, options common.DBOptions) (*v1User.User, error) {
	db := u.GetDB(options)
	var us v1User.User
	if err := db.One("Email", email, &us); err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *service) List(options common.DBOptions) ([]v1User.User, error) {
	db := u.GetDB(options)
	users := make([]v1User.User, 0)
	if err := db.All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *service) Delete(name string, options common.DBOptions) error {
	db := u.GetDB(options)
	item, err := u.Get(name, options)
	if err != nil {
		return err
	}
	if item.CreatedBy == "system" {
		return errors.New("can not delete this resource,because it created by system")
	}
	return db.DeleteStruct(item)
}

func (u *service) Create(us *v1User.User, options common.DBOptions) error {
	db := u.GetDB(options)
	us.UUID = uuid.New().String()
	us.CreateAt = time.Now()
	us.UpdateAt = time.Now()
	if us.Spec.Authenticate.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(us.Spec.Authenticate.Password), bcrypt.DefaultCost) //加密处理
		us.Spec.Authenticate.Password = string(hash)
	}
	return db.Save(us)
}
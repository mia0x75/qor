package validations_test

import (
	"regexp"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/qor/test/utils"
	"github.com/qor/qor/validations"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name       string
	CompanyID  int
	Company    Company
	CreditCard CreditCard
	Addresses  []Address
	Languages  []Language `gorm:"many2many:user_languages"`
}

func (user *User) Validate(db *gorm.DB) {
	if user.Name == "invalid" {
		validations.AddErrorForColumn(db, user, "Name", "invalid user name")
	}
}

type Company struct {
	gorm.Model
	Name string
}

func (company *Company) Validate(db *gorm.DB) {
	if company.Name == "invalid" {
		validations.AddError(db, company, "invalid company name")
	}
}

type CreditCard struct {
	gorm.Model
	UserID int
	Number string
}

func (card *CreditCard) Validate(db *gorm.DB) {
	if !regexp.MustCompile("^(\\d){13,16}$").MatchString(card.Number) {
		validations.AddErrorForColumn(db, card, "Number", "invalid card number")
	}
}

type Address struct {
	gorm.Model
	UserID  int
	Address string
}

func (address *Address) Validate(db *gorm.DB) {
	if address.Address == "invalid" {
		validations.AddErrorForColumn(db, address, "Address", "invalid address")
	}
}

type Language struct {
	gorm.Model
	Code string
}

func (language *Language) Validate(db *gorm.DB) {
	if language.Code == "invalid" {
		validations.AddErrorForColumn(db, language, "Code", "invalid language")
	}
}

func init() {
	db = utils.TestDB()
	validations.RegisterCallbacks(db)
	db.AutoMigrate(&User{}, &Company{}, &CreditCard{}, &Address{}, &Language{})
}

func TestSaveInvalidUesr(t *testing.T) {
	user := User{Name: "invalid"}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid user")
	}
}

func TestSaveInvalidCompany(t *testing.T) {
	user := User{
		Name:    "valid",
		Company: Company{Name: "invalid"},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid company")
	}
}

func TestSaveInvalidCreditCard(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "invalid"},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid credit card")
	}
}

func TestSaveInvalidAddresses(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "4111111111111111"},
		Addresses:  []Address{{Address: "invalid"}},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid addresses")
	}
}

func TestSaveInvalidLanguage(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "4111111111111111"},
		Addresses:  []Address{{Address: "valid"}},
		Languages:  []Language{{Code: "invalid"}},
	}

	if result := db.Save(&user); result.Error == nil {
		t.Errorf("Should get error when save invalid language")
	}
}

func TestSaveAllValidData(t *testing.T) {
	user := User{
		Name:       "valid",
		Company:    Company{Name: "valid"},
		CreditCard: CreditCard{Number: "4111111111111111"},
		Addresses:  []Address{{Address: "valid1"}, {Address: "valid2"}},
		Languages:  []Language{{Code: "valid1"}, {Code: "valid2"}},
	}

	if result := db.Save(&user); result.Error != nil {
		t.Errorf("Should get no error when save valid data, but got: %v", result.Error)
	}
}
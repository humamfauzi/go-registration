package registration

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserMap map[string]string

// func main() {
// 	db, err := gorm.Open("mysql", "root:@/try1?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	db.AutoMigrate(&User{})

// 	newUser := User{
// 		Email:    "c@c.c",
// 		Password: "asd",
// 		Name:     "jijop",
// 		Token:    "alskdjliowwdnaskd-asdlkasjdlo",
// 	}
// 	db.Create(&newUser)
// 	fmt.Println(newUser)

// 	var user User
// 	// db.First(&user, 2)
// 	// db.Delete(&user)
// 	search := []string{"b@b.b", "askdjwoa"}
// 	db.First(&user, "email = ?, password = ?", search)

// 	db.Model(&user).Update("Password", "secure")

// }

// type WriteUser interface {
// 	CreateUser(user User)
// 	UpdateUser(userId int, userProfile User) error
// 	DeleteUser(userId int) error
// }

// type ReadUser interface {
// 	GetAllUser() ([]User, error)
// 	GetUser(userId int) (User, error)
// 	FindUser(userProfile User) (User, error)
// 	FindAllUser(userProfile User) ([]User, error)
// }

// type UserConn struct {
// 	MySql *gorm.DB // for mysql connection
// }

// func (uc *UserConn) CreateUser(user User) {
// 	uc.MySql.Create(&user)
// }

// func (uc *UserConn) UpdateUser(oldUserProfile User, newUserProfile UserMap) {
// 	for column, value := range newUserProfile {
// 		uc.MySql.Model(&oldUserProfile).Update(column, value)
// 	}
// }

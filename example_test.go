package clerk_test

import "github.com/arthurwhite/clerk"

type User struct {
	Name string
}

var (
	db   clerk.DB
	data = &struct {
		Users []*User
	}{}
)

func init() {
	var err error
	if db, err = clerk.New("data.gob", data); err != nil {
		panic(err)
	}
}

func Example() {
	data.Users = append(data.Users, &User{"Crowe"}, &User{"Jones"}, &User{"Owen"})
	if err := db.Save(); err != nil {
		panic(err)
	}
}

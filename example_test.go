package clerk_test

import "github.com/arthurwhite/clerk"

type User struct {
	Name string
}

var (
	db   *clerk.DB
	data = &struct {
		Users []*User
	}{}
)

func Example() {
	var err error
	if db, err = clerk.New("data.gob", data); err != nil {
		panic(err)
	}
	defer db.Remove() // Remove database at exit.

	data.Users = append(data.Users, &User{"Crowe"}, &User{"Jones"}, &User{"Owen"})
	if err := db.Save(); err != nil {
		panic(err)
	}
}

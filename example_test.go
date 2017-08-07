package clerk_test

import "github.com/arthurwhite/clerk"

type User struct {
	Name string
}

var db = new(struct {
	clerk.DB
	Users []*User
})

func init() {
	clerk.Init("example.gob", db)
}

func Example() {
	db.Lock()
	defer db.Unlock()

	db.Users = append(db.Users, &User{"Crowe"}, &User{"Jones"}, &User{"Owen"})

	if err := db.Save(); err != nil {
		panic(err)
	}
}

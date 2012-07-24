package mysql

import (
	"database/sql"
	"fmt"
	"github.com/kr/pretty"
	"github.com/xiam/gosexy/db"
	"math/rand"
	"testing"
)

const myHost = "10.0.0.11"
const myDatabase = "gotest"
const myUser = "gouser"
const myPassword = "gopass"

func TestMyTruncate(t *testing.T) {

	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	collections := sess.Collections()

	for _, name := range collections {
		col := sess.Collection(name)
		col.Truncate()
		if col.Count() != 0 {
			t.Errorf("Could not truncate '%s'.", name)
		}
	}

}

func TestMyAppend(t *testing.T) {

	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	col := sess.Collection("people")

	col.Truncate()

	names := []string{"Juan", "José", "Pedro", "María", "Roberto", "Manuel", "Miguel"}

	for i := 0; i < len(names); i++ {
		col.Append(db.Item{"name": names[i]})
	}

	if col.Count() != len(names) {
		t.Error("Could not append all items.")
	}

}

func TestMyFind(t *testing.T) {

	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	col := sess.Collection("people")

	result := col.Find(db.Cond{"name": "José"})

	if result["name"] != "José" {
		t.Error("Could not find a recently appended item.")
	}

}

func TestMyDelete(t *testing.T) {
	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	col := sess.Collection("people")

	col.Remove(db.Cond{"name": "Juan"})

	result := col.Find(db.Cond{"name": "Juan"})

	if len(result) > 0 {
		t.Error("Could not remove a recently appended item.")
	}
}

func TestMyUpdate(t *testing.T) {
	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	sess.Use("test")

	col := sess.Collection("people")

	col.Update(db.Cond{"name": "José"}, db.Set{"name": "Joseph"})

	result := col.Find(db.Cond{"name": "Joseph"})

	if len(result) == 0 {
		t.Error("Could not update a recently appended item.")
	}
}

func TestMyPopulate(t *testing.T) {
	var i int

	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	sess.Use("test")

	places := []string{"Alaska", "Nebraska", "Alaska", "Acapulco", "Rome", "Singapore", "Alabama", "Cancún"}

	for i = 0; i < len(places); i++ {
		sess.Collection("places").Append(db.Item{
			"code_id": i,
			"name":    places[i],
		})
	}

	people := sess.Collection("people").FindAll(
		db.Fields{"id", "name"},
	)

	for i = 0; i < len(people); i++ {
		person := people[i]

		// Has 5 children.
		for j := 0; j < 5; j++ {
			sess.Collection("children").Append(db.Item{
				"name":      fmt.Sprintf("%s's child %d", person["name"], j+1),
				"parent_id": person["id"],
			})
		}

		// Lives in
		sess.Collection("people").Update(
			db.Cond{"id": person["id"]},
			db.Set{"place_code_id": int(rand.Float32() * float32(len(places)))},
		)

		// Has visited
		for k := 0; k < 3; k++ {
			place := sess.Collection("places").Find(db.Cond{
				"code_id": int(rand.Float32() * float32(len(places))),
			})
			sess.Collection("visits").Append(db.Item{
				"place_id":  place["id"],
				"person_id": person["id"],
			})
		}
	}

}

func TestMyRelation(t *testing.T) {
	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	col := sess.Collection("people")

	result := col.FindAll(
		db.Relate{
			"lives_in": db.On{
				sess.Collection("places"),
				db.Cond{"code_id": "{place_code_id}"},
			},
		},
		db.RelateAll{
			"has_children": db.On{
				sess.Collection("children"),
				db.Cond{"parent_id": "{id}"},
			},
			"has_visited": db.On{
				sess.Collection("visits"),
				db.Cond{"person_id": "{id}"},
				db.Relate{
					"place": db.On{
						sess.Collection("places"),
						db.Cond{"id": "{place_id}"},
					},
				},
			},
		},
	)

	fmt.Printf("%# v\n", pretty.Formatter(result))
}

func TestCustom(t *testing.T) {
	sess := Session(db.DataSource{Host: myHost, Database: myDatabase, User: myUser, Password: myPassword})

	err := sess.Open()
	defer sess.Close()

	if err != nil {
		panic(err)
	}

	_, err = sess.Driver().(*sql.DB).Query("SELECT NOW()")

	if err != nil {
		panic(err)
	}

}
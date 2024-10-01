package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Person struct {
	Name      string
	Age       int
	IsBlocked bool
}

type User struct {
	Name      string `json:"at"` // eksportiruemye nazwaniya (Bas harp bilen baslanmaly)
	Age       int    `json:"age"`
	IsBlocked bool   `json:"is_blocked"`
	Books     []Book `json:"books"`
}

type Book struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func main() {
	sjsonPackage()
	gjsonPackage()
	addModifier()
	marshal()
	unMarshal()
	//jsonTest()
}
func addModifier() {
	jsonVal := `{
    "name": {"first": "Tom", "last": "Anderson"},
    "age": 30,
    "children": ["Sarah", "Alex", "Jack"],
    "fav.movie": "Harry Potter",
    "friends": [
      {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["instagram", "facebook", "telegram"]},
      {"first": "Mark", "last": "Roger", "age": 34, "nets": ["instagram", "facebook"]},
      {"first": "Jane","last": "Murphy", "age": 44, "nets": ["instagram","telegram"]}
    ]
  }`

	gjson.AddModifier("case", func(jsonVal, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(jsonVal)
		} else if arg == "lower" {
			return strings.ToLower(jsonVal)
		}
		return jsonVal
	})
	if !gjson.Valid(jsonVal) { // json validate edip beryar. json valid bolmasada kabir gte funksiyalar isleyar
		panic("Jason is not valid")
	}
	result := gjson.Get(jsonVal, "friends.2.nets|@case:upper")
	fmt.Println("case: upper\n", result)
	result = gjson.Get(jsonVal, "friends.0.nets|@case:lower")
	fmt.Println("case: lower\n", result)

	fmt.Println(gjson.Parse(jsonVal).Get("name"))
	for _, val := range result.Array() {
		fmt.Println(val.String())
	}

	for _, val := range gjson.Get(jsonVal, "friends.1").Array() {
		fmt.Println(val.String())
	}

	fmt.Println()

	res, ok := gjson.Parse(jsonVal).Value().(map[string]interface{})
	if !ok {
		panic("Not okay parsing to map")
	}
	fmt.Println(res)

}

func sjsonPackage() {
	jsonVal := `{
    "name": {"first": "Tom", "last": "Anderson"},
    "age": 30,
    "children": ["Sarah", "Alex", "Jack"],
    "fav.movie": "Harry Potter",
    "friends": [
      {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["instagram", "facebook", "telegram"]},
      {"first": "Mark", "last": "Roger", "age": 34, "nets": ["instagram", "facebook"]},
      {"first": "Jane","last": "Murphy", "age": 44, "nets": ["instagram","telegram"]}
    ]
  }`

	result1, _ := sjson.Set(jsonVal, "name.first", "Artur") // update value in json and return json
	fmt.Println(result1)

	result2, _ := sjson.Set(jsonVal, "children.3", "Artur") // len(slice) den uly index gorkezsen taze element gosyar
	fmt.Println(result2)

	result3, _ := sjson.Set(jsonVal, "children.-1", "Tom") // sonky elementi update edyar
	fmt.Println(result3)

	newValue, _ := sjson.Set(jsonVal, "fav_books", map[string]interface{}{"hello": "world"}) // taze field gosyar
	fmt.Println(newValue)

	delValue, _ := sjson.Delete(jsonVal, "friends")
	fmt.Println(delValue)
}

func gjsonPackage() {
	newJson1 := `{"name":{"first":"Nick", "last":"Portman"}, "age":20, "phone":[1111,2222] }`
	value := gjson.Get(newJson1, "name.first")
	fmt.Println(value) // Nick

	fmt.Printf("%#v\n", value.String())
	value = gjson.Get(newJson1, "phone.0")
	fmt.Printf("%v\n", value)
	fmt.Println()

	newjson2 := `{
    "name": {"first": "Tom", "last": "Anderson"},
    "age": 30,
    "children": ["Sarah", "Alex", "Jack"],
    "fav.movie": "Harry Potter",
    "friends": [
      {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["instagram", "facebook", "telegram"]},
      {"first": "Mark", "last": "Roger", "age": 34, "nets": ["instagram", "facebook"]},
      {"first": "Jane", "last": "Murphy", "age": 44, "nets": ["instagram","telegram"]}
    ]
  }`

	value2 := gjson.Get(newjson2, "children")
	fmt.Printf("Children : %v\n", value2)

	fmt.Println(gjson.Get(newjson2, "children.0"))
	fmt.Println("length: ", gjson.Get(newjson2, "children.#"))
	fmt.Println(gjson.Get(newjson2, "chil*.2")) // chil* ozi gozlap tapyar
	fmt.Println(gjson.Get(newjson2, "fav\\.movie"))
	fmt.Println("Friends: ", gjson.Get(newjson2, "fri*"))
	fmt.Println("Friend: ", gjson.Get(newjson2, "fri*.0.last"))
	fmt.Println("Friend's nets: ", gjson.Get(newjson2, "fri*.2.nets.1"))

	fmt.Println(gjson.Get(newjson2, `friends.#(last=="Roger").first`)) // friends obyektin icinden last==Roger bolandaky first elementi alyar

	fmt.Println("First found friend who is 44 =>", gjson.Get(newjson2, `friends.#(age==44).first`)) // age==47 bolanlary alyar, dine birinji gozlap tapanyny alyar
	fmt.Println("All found friends who is 44 =>", gjson.Get(newjson2, `friends.#(age==44)#.first`))

}

func unMarshal() {
	byt := []byte(`{"name":"Artur","age":20,"is_blocked":true,"books":[{"name":"bookName","year":1990},{"name":"bookName","year":1890}]}`)
	var datMap map[string]interface{}

	if err := json.Unmarshal(byt, &datMap); err != nil {
		panic(err)
	}

	fmt.Println(datMap)
	fmt.Println(datMap["name"])
	bookName := datMap["books"].([]interface{})[0].(map[string]interface{})["name"]
	fmt.Println(bookName)
	fmt.Println()

	byt2 := []byte(`{"at":"Artur","age":20,"is_blocked":true,"books":[{"name":"bookName","year":1990},{"name":"bookName","year":1890}]}`)
	var datMap2 User

	if err := json.Unmarshal(byt2, &datMap2); err != nil {
		panic(err)
	}
	fmt.Println(datMap2)
	fmt.Println("At", datMap2.Name)
}

func marshal() {
	boolVal, _ := json.Marshal(true)
	fmt.Println(string(boolVal))

	stringVal, _ := json.Marshal("String")
	fmt.Println(string(stringVal))

	intVal, _ := json.Marshal(13245)
	fmt.Println(string(intVal))

	slice := []string{"A", "B", "C"}
	sliceVal, _ := json.Marshal(slice)
	fmt.Println(string(sliceVal))
	fmt.Println()

	mapString := map[int]string{1: "first", 2: "second", 3: "third"}
	val, _ := json.Marshal(mapString)
	fmt.Println(string(val))
	fmt.Println()

	mapInterface := map[string]interface{}{"first": 1, "second": "two", "third": true, "fourth": []int{1, 2, 3}}
	mapVal, _ := json.Marshal(mapInterface)
	fmt.Printf("TYpe %T, value: %#v\n", mapVal, string(mapVal))
	fmt.Println()

	person := Person{
		Name:      "John",
		Age:       32,
		IsBlocked: false,
	}

	personVal, _ := json.Marshal(person)
	fmt.Println(string(personVal))
	fmt.Println()

	var books []Book

	book1 := Book{Name: "bookName", Year: 1990}
	book2 := Book{Name: "bookName", Year: 1890}

	books = append(books, book1, book2)

	user := User{
		Name:      "Artur",
		Age:       20,
		IsBlocked: true,
		Books:     books,
	}

	userVal, _ := json.Marshal(user)
	fmt.Println(string(userVal))
}

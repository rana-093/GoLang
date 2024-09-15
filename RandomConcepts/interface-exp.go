package RandomConcepts

import "fmt"

type Animal interface {
	Speak() string
}

type Cat struct {
	name string
	id   int
}

type Dog struct {
	name        string
	dateOfBirth int
}

func (c Cat) Speak() string {
	return fmt.Sprintf("Hi %s, I am %s. I call meow meow", c.name, c.id)
}

func (d Dog) Speak() string {
	return fmt.Sprintf("Hi I am %s, I call ghew ghew!", d.name)
}

func TestInterface() error {
	animals := []Animal{
		Cat{
			name: "Alibi",
			id:   12,
		},
		Dog{
			name:        "Ok Here",
			dateOfBirth: 123,
		},
	}

	for _, animal := range animals {
		animal.Speak()
	}

	return nil
}

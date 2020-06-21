package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"shoppinglistserver/log"
	"testing"
)

func initTest() {
	log.InitLogger("INFO")
	storage := "./test.sqlite"
	_ = os.Remove(storage)
	err := InitStorage(storage)
	if err != nil {
		panic(err)
	}
}

func Test_GeneralBehaviour(t *testing.T) {
	initTest()
	err := New("One")
	assert.NoError(t, err)

	err = New("Two")
	assert.NoError(t, err)

	err = New("Three")
	assert.NoError(t, err)

	err = ToggleChecked(1)
	assert.NoError(t, err)

	all, err := GetAll()
	assert.NoError(t, err)

	// Should be: Two(0), Three(0), One(1)

	assert.Len(t, all, 3, "Should contain 3 items")
	assert.Equal(t, "Two", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[0].Checked, "Should be equal")
	assert.Equal(t, "Three", all[1].Name, "Should be equal")
	assert.Equal(t, 1, all[1].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[1].Checked, "Should be equal")
	assert.Equal(t, "One", all[2].Name, "Should be equal")
	assert.Equal(t, 2, all[2].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[2].Checked, "Should be equal")

	err = ToggleChecked(1)
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)

	// Should be: Two(0), Three(0), One(0)

	assert.Len(t, all, 3, "Should contain 3 items")
	assert.Equal(t, "Two", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[0].Checked, "Should be equal")
	assert.Equal(t, "Three", all[1].Name, "Should be equal")
	assert.Equal(t, 1, all[1].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[1].Checked, "Should be equal")
	assert.Equal(t, "One", all[2].Name, "Should be equal")
	assert.Equal(t, 2, all[2].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[2].Checked, "Should be equal")

	err = ToggleChecked(1)
	assert.NoError(t, err)

	// Status: Two(0), Three(0), One(1)

	err = ToggleChecked(2)
	assert.NoError(t, err)

	// Status: Three(0), One(1), Two(1)

	err = ToggleChecked(3)
	assert.NoError(t, err)

	// Status: One(1), Two(1), Three(1)

	all, err = GetAll()
	assert.NoError(t, err)

	// Should be: One(1), Two(1), Three(1)

	assert.Len(t, all, 3, "Should contain 3 items")
	assert.Equal(t, "One", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[0].Checked, "Should be equal")
	assert.Equal(t, "Two", all[1].Name, "Should be equal")
	assert.Equal(t, 1, all[1].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[1].Checked, "Should be equal")
	assert.Equal(t, "Three", all[2].Name, "Should be equal")
	assert.Equal(t, 2, all[2].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[2].Checked, "Should be equal")

	err = ToggleChecked(1)
	assert.NoError(t, err)

	// Status: One(0), Two(1), Three(1)

	err = ToggleChecked(2)
	assert.NoError(t, err)

	// Status: One(0), Two(0), Three(1)

	err = ToggleChecked(3)
	assert.NoError(t, err)

	// Status: One(0), Two(0), Three(0)

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Len(t, all, 3, "Should contain 3 items")
	assert.Equal(t, "One", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[0].Checked, "Should be equal")
	assert.Equal(t, "Two", all[1].Name, "Should be equal")
	assert.Equal(t, 1, all[1].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[1].Checked, "Should be equal")
	assert.Equal(t, "Three", all[2].Name, "Should be equal")
	assert.Equal(t, 2, all[2].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[2].Checked, "Should be equal")

	err = ToggleChecked(3)
	assert.NoError(t, err)
	err = ToggleChecked(2)
	assert.NoError(t, err)
	err = ToggleChecked(1)
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Len(t, all, 3, "Should contain 3 items")
	assert.Equal(t, "Three", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[0].Checked, "Should be equal")
	assert.Equal(t, "Two", all[1].Name, "Should be equal")
	assert.Equal(t, 1, all[1].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[1].Checked, "Should be equal")
	assert.Equal(t, "One", all[2].Name, "Should be equal")
	assert.Equal(t, 2, all[2].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[2].Checked, "Should be equal")

	err = ToggleChecked(1)
	assert.NoError(t, err)

	err = DeleteAllChecked()
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Len(t, all, 1, "Should contain only 1 item")
	assert.Equal(t, "One", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[0].Checked, "Should be equal")

	err = DeleteOne(1)
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 0, "Should contain no items")

	err = New("New")
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)

	newItem := all[0]
	assert.Equal(t, 0, newItem.ListOrder, "Should be equal")
	assert.Equal(t, 0, newItem.Checked, "Should be equal")

	err = ToggleChecked(newItem.Id)
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[0].Checked, "Should be equal")

	err = New("Other")
	assert.NoError(t, err)

	all, err = GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 2, "Should be equal")
	assert.Equal(t, "Other", all[0].Name, "Should be equal")
	assert.Equal(t, 0, all[0].ListOrder, "Should be equal")
	assert.Equal(t, 0, all[0].Checked, "Should be equal")
	assert.Equal(t, "New", all[1].Name, "Should be equal")
	assert.Equal(t, 1, all[1].ListOrder, "Should be equal")
	assert.Equal(t, 1, all[1].Checked, "Should be equal")
}

func Test_UpdateBadId(t *testing.T) {
	initTest()
	err := New("One")
	assert.NoError(t, err)

	all, err := GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	err = Update("New", 10)
	assert.Error(t, err, "Should fail")

	all, err = GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)
	assert.Equal(t, "One", all[0].Name, "Name should not have changed")
}

func Test_DeleteBadId(t *testing.T) {
	initTest()
	err := New("One")
	assert.NoError(t, err)

	all, err := GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	err = DeleteOne(10)
	assert.Error(t, err, "Should fail")

	err = DeleteOne(10)
	assert.Error(t, err, "Should fail")

	all, err = GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)
	assert.Equal(t, "One", all[0].Name, "Name should not have changed")
}

func Test_StrangePath(t *testing.T) {
	initTest()
	err := New("Patatas")
	assert.NoError(t, err)

	err = New("Cacahuetes")
	assert.NoError(t, err)

	// Patatas(0), Cacahuetes(0)

	err = ToggleChecked(1)
	assert.NoError(t, err)

	// Cacahuetes(0), Patatas(1)

	err = ToggleChecked(2)
	assert.NoError(t, err)

	// Patatas(1), Cacahuetes(1)

	err = ToggleChecked(2)
	assert.NoError(t, err)

	// Cacahuetes(0), Patatas(1)

	all, err := GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 2)

	assert.Equal(t, all[0].Id, 2)
	assert.Equal(t, all[0].Name, "Cacahuetes")
	assert.Equal(t, all[0].ListOrder, 0)
	assert.Equal(t, all[1].Id, 1)
	assert.Equal(t, all[1].Name, "Patatas")
	assert.Equal(t, all[1].ListOrder, 1)
}

func Test_Reorder(t *testing.T) {
	initTest()

	assert.NoError(t, New("One"))
	assert.NoError(t, New("Two"))
	assert.NoError(t, New("Three"))
	assert.NoError(t, New("Four"))

	assert.NoError(t, MoveToPosition(1, 3))

	all, err := GetAll()
	assert.NoError(t, err)

	assert.Equal(t, "Two", all[0].Name)
	assert.Equal(t, 0, all[0].ListOrder)
	assert.Equal(t, "Three", all[1].Name)
	assert.Equal(t, 1, all[1].ListOrder)
	assert.Equal(t, "Four", all[2].Name)
	assert.Equal(t, 2, all[2].ListOrder)
	assert.Equal(t, "One", all[3].Name)
	assert.Equal(t, 3, all[3].ListOrder)

	assert.NoError(t, MoveToPosition(1, 0))

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Equal(t, "One", all[0].Name)
	assert.Equal(t, 0, all[0].ListOrder)
	assert.Equal(t, "Two", all[1].Name)
	assert.Equal(t, 1, all[1].ListOrder)
	assert.Equal(t, "Three", all[2].Name)
	assert.Equal(t, 2, all[2].ListOrder)
	assert.Equal(t, "Four", all[3].Name)
	assert.Equal(t, 3, all[3].ListOrder)

	assert.NoError(t, MoveToPosition(2, 2))

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Equal(t, "One", all[0].Name)
	assert.Equal(t, 0, all[0].ListOrder)
	assert.Equal(t, "Three", all[1].Name)
	assert.Equal(t, 1, all[1].ListOrder)
	assert.Equal(t, "Two", all[2].Name)
	assert.Equal(t, 2, all[2].ListOrder)
	assert.Equal(t, "Four", all[3].Name)
	assert.Equal(t, 3, all[3].ListOrder)

	assert.NoError(t, MoveToPosition(3, 2))

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Equal(t, "One", all[0].Name)
	assert.Equal(t, 0, all[0].ListOrder)
	assert.Equal(t, "Two", all[1].Name)
	assert.Equal(t, 1, all[1].ListOrder)
	assert.Equal(t, "Three", all[2].Name)
	assert.Equal(t, 2, all[2].ListOrder)
	assert.Equal(t, "Four", all[3].Name)
	assert.Equal(t, 3, all[3].ListOrder)

}

func Test_Reorder2(t *testing.T) {
	initTest()

	assert.NoError(t, New("A"))
	assert.NoError(t, New("B"))
	assert.NoError(t, New("Y"))
	assert.NoError(t, New("Z"))

	assert.NoError(t, ToggleChecked(3))
	assert.NoError(t, ToggleChecked(4))

	all, err := GetAll()
	assert.NoError(t, err)

	assert.Equal(t, "A", all[0].Name)
	assert.Equal(t, 0, all[0].ListOrder)
	assert.Equal(t, "B", all[1].Name)
	assert.Equal(t, 1, all[1].ListOrder)
	assert.Equal(t, "Y", all[2].Name)
	assert.Equal(t, 2, all[2].ListOrder)
	assert.Equal(t, "Z", all[3].Name)
	assert.Equal(t, 3, all[3].ListOrder)

	assert.NoError(t, MoveToPosition(3, 3))

	all, err = GetAll()
	assert.NoError(t, err)

	assert.Equal(t, "A", all[0].Name)
	assert.Equal(t, 0, all[0].ListOrder)
	assert.Equal(t, "B", all[1].Name)
	assert.Equal(t, 1, all[1].ListOrder)
	assert.Equal(t, "Z", all[2].Name)
	assert.Equal(t, 2, all[2].ListOrder)
	assert.Equal(t, "Y", all[3].Name)
	assert.Equal(t, 3, all[3].ListOrder)
}

func Test_AddSameItemTwice(t *testing.T) {
	initTest()
	assert.NoError(t, New("A"))
	assert.NoError(t, ToggleChecked(1))
	assert.NoError(t, New("A"))
	assert.Error(t, New("A"))

}

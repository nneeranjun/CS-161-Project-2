package proj2

// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder, and we will be Very Upset.

import (
	_ "encoding/hex"
	_ "encoding/json"
	_ "errors"
	"github.com/cs161-staff/userlib"
	_ "github.com/google/uuid"
	"reflect"
	_ "strconv"
	_ "strings"
	"testing"
)

func clear() {
	// Wipes the storage so one test does not affect another
	userlib.DatastoreClear()
	userlib.KeystoreClear()
}

func TestInit(t *testing.T) {
	clear()
	t.Log("Initialization test")

	// You can set this to false!
	userlib.SetDebugStatus(true)

	u, err := InitUser("alice", "fubar")
	if err != nil {
		// t.Error says the test fails
		t.Error("Failed to initialize user", err)
		return
	}
	t.Log(make([]byte, 16))
	// t.Log() only produces output if you run with "go test -v"
	t.Log("Got user", u)
	// If you want to comment the line above,
	// write _ = u here to make the compiler happy
	// You probably want many more tests here.
	u1, err1 := GetUser("alice", "fubar")
	if err1 != nil {
		// t.Error says the test fails
		t.Error("Failed to login user, error:", err1)
		return
	}
	t.Log("Logged in user", u1)
}

func TestInitNew(t *testing.T) {
	clear()
	t.Log("Initialization test")

	userlib.SetDebugStatus(true)

	_, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Error: Failed to initialize the user", err)
		return
	}

	_, err2 := InitUser("alice", "fubar")
	if err2 == nil {
		t.Error("Error: Username not unique", err)
		return
	}

	_, err3 := GetUser("alice", "fubar")
	if err3 != nil {
		t.Error("Error: Failed to get User", err)
		return
	}

	_, err4 := GetUser("alice", "incorrect_password")
	if err4 == nil {
		t.Error("Error: Allowed login with wrong password", err)
		return
	}
}

func TestStorage(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	//v := []byte("This is a test")
	//u.StoreFile("file1", v)
	//print("Length of map: ")
	////print(len(u.FileMap))
	//u.StoreFile("file2", []byte("This is not a test"))
	//print("Length of map: ")
	////print(len(u.FileMap))
	file, err := u.LoadFile("file1")
	if err != nil {
		print(err)
	} else {
		print(file)
	}

/*
	v2, err2 := u.LoadFile("file1")
	if err2 != nil {
		t.Error("Failed to upload and download", err2)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Downloaded file is not the same", v, v2)
		return
	}
*/

}

func TestStorageNew(t *testing.T) {
	clear()
	user1, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	data := []byte("Testing storage")
	user1.StoreFile("file1", data)

	returnData, err2 := user1.LoadFile("file1")
	if err2 != nil {
		t.Error("Error: Could not load file", err2)
		return
	}
	if !reflect.DeepEqual(data, returnData) {
		t.Error("Error: Did not load file properly")
		return
	}

	user2, err := InitUser("Connor", "McGregor")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	file1 := []byte("File 2 testing data")
	file2 := []byte("Something different")

	user2.StoreFile("File1", file1)
	user2.StoreFile("File1", file2)

	returnData, err = user2.LoadFile("File1")
	if err != nil {
		t.Error("Error: Loaded data is not the same", err)
		return
	}

	if !reflect.DeepEqual(returnData, file2) {
		t.Error("Error: Loaded data is not the same")
		return
	}

	fileNew := append(file2, file1...)
	err = user2.AppendFile("File1", file1)
	if err != nil {
		t.Error("Error: Could not append file", err)
	}
	fileCompare, err := user2.LoadFile("File1")
	if !reflect.DeepEqual(fileCompare, fileNew) {
		t.Error("Error: Could not load file properly")
		return
	}

	user3, _ := InitUser("Khabib", "RussianLastName")
	newFile1 := []byte("File1data")
	newFile2 := []byte("File2data")
	appendedFile := append(newFile1, newFile2...)
	user3.StoreFile("File1", newFile1)
	user3.StoreFile("File2", newFile2)
	err = user3.AppendFile("File1", newFile2)
	newnewFile1, _ := user3.LoadFile("File1")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(appendedFile, newnewFile1) {
		t.Error("Error: Loaded file doesn't match")
	}
	user4, _ := InitUser("Nilay", "Neeranjun")
	user5, _ := InitUser("Kobe", "Bryant")
	user4 = user4
	user5 = user5

	file1 = []byte("This is a TEST FILE OHHHHHHUHUHDUSHIUSHDIUSHDIUSHDUSHDHSDHSIUDH YEAH")
	user1.StoreFile("file1", file1)

	user2.Username = "NILAY"
	val, err := user2.LoadFile("file1")
	val = val
	if err == nil {
		t.Error("Error: YOU CAN'T CHANGE A USERNAME")
	}

	clear()
	nilay, _ := InitUser("Nilay", "Neeranjun")
	albert, _ := InitUser("Albert", "Zhang")

	file1 = []byte("TEST FILE")
	nilay.StoreFile("File1", file1)

	magic_string, err := nilay.ShareFile("File1", "Albert")
	if err != nil {
		t.Error("COULD NOT SHARE FILE !!!!!", err)
		return
	}
	err = albert.ReceiveFile("STORE FILE", "Nilay", magic_string)
	if err != nil {
		t.Error("COULD NOT LOAD RECEIVED FILE")
		return
	}

	mutatedFile := []byte("CHANGED CONTENTS")
	albert.StoreFile("File1", mutatedFile)

	receivedFile, err := albert.LoadFile("STORE FILE")
	receivedFile = receivedFile
	if err != nil {
		t.Error("DID NOT LOAD FILE IN THE CORRECT MANNER")
		return
	}
}


func TestInvalidFile(t *testing.T) {
	clear()

	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	_, err2 := u.LoadFile("this file does not exist")
	if err2 == nil {
		t.Error("Downloaded a nonexistent file", err2)
		return
	}


}


func TestShare(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}
	u2, err2 := InitUser("bob", "foobar")
	if err2 != nil {
		t.Error("Failed to initialize bob", err2)
		return
	}

	v := []byte("This is a test")
	u.StoreFile("file1", v)
	
	var v2 []byte
	var magic_string string

	v, err = u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
		return
	}

	magic_string, err = u.ShareFile("file1", "bob")
	if err != nil {
		t.Error("Failed to share the a file", err)
		return
	}
	err = u2.ReceiveFile("file2", "alice", magic_string)
	if err != nil {
		t.Error("Failed to receive the share message", err)
		return
	}

	v2, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing", err)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Shared file is not the same", v, v2)
		return
	}



}

func TestDataStoreCoverage(t *testing.T) {
	clear()
	//Set Users Alice and Bob
	u1, err := InitUser("alice", "password")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	u1.StoreFile("file1", []byte("hello"))

	datastoremap := userlib.DatastoreGetMap()

	for index,_ := range datastoremap{
		datastoremap[index] = []byte("Evil")
	}

	lf, err := u1.LoadFile("file1")
	if err == nil{
		t.Error("Failed to check modification", err)
		return
	}

	_ = lf


}

func TestAppendFileNew(t *testing.T) {
	clear()

	//Set Users Alice and Bob
	u1, err := InitUser("alice", "password")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	u2, err := InitUser("bob", "password")
	if err != nil {
		t.Error("COULD NOT initialize user", err)
		return
	}

	u1.StoreFile("File1", []byte("File1 CONTENTS"))

	err = u2.AppendFile("File1",[]byte("NEWTESTDATAASDSD"))
	if err == nil{
		t.Error("STILL APPENDED")
		return
	}

	magic_string, err := u1.ShareFile("File1", "bob")
	if err != nil {
		t.Error("COULD NOT SHARE TEH FREAKING FILE")
		return
	}

	//Bob recievies a corrupted magic string
	err = u2.ReceiveFile("NEWFILENAME", "alice", "WRONG MAGIC STRING")
	if err == nil {
		t.Error("Still could receive file even with wrong magic string")
		return
	}

	err = u2.ReceiveFile("NEWFILENAME", "alice", magic_string)
	if err != nil {
		t.Error("COULD NOT RECEIVE FILE NORMALLY", err)
		return
	}

	err = u2.AppendFile("NEWFILENAME", []byte("APPENDING SOME BYTES YEAAAA"))
	if err != nil{
		t.Error("COULD NOT APPEND FILE PROPERLY HAHAHHAHAHHA")
		return
	}
	//Alice receives the changes
	receivedFile, err := u1.LoadFile("File1")
	if err != nil {
		t.Error("ERROROROROOROROROOROR")
		return
	}
	if !reflect.DeepEqual(receivedFile, []byte("File1 CONTENTSAPPENDING SOME BYTES YEAAAA")){
		t.Error("Data did not append correctly")
		return
	}

}





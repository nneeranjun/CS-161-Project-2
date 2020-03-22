package proj2

// CS 161 Project 2 Spring 2020
// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder. We will be very upset.

import (
	// You neet to add with
	// go get github.com/cs161-staff/userlib
	"github.com/cs161-staff/userlib"

	// Life is much easier with json:  You are
	// going to want to use this so you can easily
	// turn complex structures into strings etc...
	"encoding/json"

	// Likewise useful for debugging, etc...
	"encoding/hex"

	// UUIDs are generated right based on the cryptographic PRNG
	// so lets make life easier and use those too...
	//
	// You need to add with "go get github.com/google/uuid"
	"github.com/google/uuid"

	// Useful for debug messages, or string manipulation for datastore keys.
	"strings"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"

	// if you are looking for fmt, we don't give you fmt, but you can use userlib.DebugMsg.
	// see someUsefulThings() below:
)

// This serves two purposes: 
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.
// Of course, this function can be deleted.
func someUsefulThings() {
	// Creates a random UUID
	f := uuid.New()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Will actually work with go structures as well
	d, _ := json.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g uuid.UUID
	json.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// This creates an error type
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// And a random RSA key.  In this case, ignoring the error
	// return value
	var pk userlib.PKEEncKey
        var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)
}

// Helper function: Takes the first 16 bytes and
// converts it into the UUID type
func bytesToUUID(data []byte) (ret uuid.UUID) {
	for x := range ret {
		ret[x] = data[x]
	}
	return
}

// The structure definition for a user record
type User2 struct {
	Username string
	MacUsername[] byte
	SaltHKDF[] byte
	SaltPassword[] byte
	Password[] byte
	PrivateKeys map[string]string
	//TODO: Maybe add files and access tokens???


	// You can add other fields here if you want...
	// Note for JSON to marshal/unmarshal, the fields need to
	// be public (start with a capital letter)
}

type User struct {
	Username string
	Password string
	SecretKeyEnc userlib.PKEDecKey
	DSSignKey userlib.DSSignKey
	MAC[] byte
	FileMap map[string][]byte
	//TODO: Add files and access tokens


	// You can add other fields here if you want...
	// Note for JSON to marshal/unmarshal, the fields need to
	// be public (start with a capital letter)
}



// This creates a user.  It will only be called once for a user
// (unless the keystore and datastore are cleared during testing purposes)

// It should store a copy of the userdata, suitably encrypted, in the
// datastore and should store the user's public key in the keystore.

// The datastore may corrupt or completely erase the stored
// information, but nobody outside should be able to get at the stored

// You are not allowed to use any global storage other than the
// keystore and the datastore functions in the userlib library.

// You can assume the password has strong entropy, EXCEPT
// the attackers may possess a precomputed tables containing 
// hashes of common passwords downloaded from the internet.

func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata

	//TODO: This is a toy implementation.
	/*
	userdata.Username = username //setting username
	userdata.PrivateKeys = make(map[string]string) //initializing map
	//generate random salt for hkdf and password
	var saltHKDF = userlib.RandomBytes(20)
	var saltPass = userlib.RandomBytes(20)

	userdata.SaltHKDF = saltHKDF
	userdata.SaltPassword = saltPass
	//set password to be some sort of hash of the password with the salt
	var hashedPass = userlib.Argon2Key([]byte (password), saltPass, 32) //hashing password
	userdata.Password = hashedPass //setting hashed password

	var hkdfKey = userlib.Argon2Key([]byte (password), saltHKDF, 32) //key generated to generate more keys using HKDF

	var macKey, _ = userlib.HashKDF(hkdfKey, []byte("mac")) //MAC key used for MAC-ing other keys
	var symmEncKey, _ = userlib.HashKDF(hkdfKey, []byte("symmetric_encryption")) //symmetric key used for encryption
	macKey = macKey[:16]
	symmEncKey = symmEncKey[:16]

	//Generate public & private keys
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()

	var macSK, _ = userlib.HMACEval(macKey, []byte("This is " + username +(sk.KeyType))) //MAC used for secret key generated above (Do we include SK here)?
	var encSK = string(userlib.SymEnc(symmEncKey, userlib.RandomBytes(16), sk.PrivKey.D.Bytes())) //encrypting the secret key so attackers can't see
	userdata.PrivateKeys[encSK] = string(macSK) //mapping private key to mac for verification

	e := userlib.KeystoreSet(username, pk) //storing public key in Keystore
	//error if username already exists
	if e != nil {
		return nil, e
	}
	//Jsonify user struct data and store in DataStore
	var all_zeroes = make([]byte, 16)
	var macUsername, _ = userlib.HMACEval(all_zeroes, []byte(username)) //Hash (MAC) username so that we can use bytesToUUID
	userdata.MacUsername = macUsername
	var UUID = bytesToUUID(macUsername)                             //generate UUID from Hash of username
	//update our struct here

	var data, _ = json.Marshal(userdata)
	userlib.DatastoreSet(UUID, data)

	//End of toy implementation
	*/

	//******************* START OF NEW IMPLEMENTATION ********************************************************************
	var hkdfKey = userlib.Argon2Key([]byte (password), []byte(username), 32) //key generated to generate more keys using HKDF

	//MAC key generated for MAC'ing
	var macKey, _ = userlib.HashKDF(hkdfKey, []byte("mac"))
	macKey = macKey[:16]

	//Symmetric key generated used for symmetric encryption
	var symmetricKey, _ = userlib.HashKDF(hkdfKey, []byte("symmetric key"))
	symmetricKey = symmetricKey[:16]

	//Generate public & private keys for public key crypto. RSA Encryption guarantees confidentiality for asymmetric-keys.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()

	//Generate public & private keys for digital signatures. RSA Signatures guarantee integrity + authenticity for asymmetric-keys.
	var vk userlib.DSVerifyKey
	var dssk userlib.DSSignKey
	dssk, vk, _ = userlib.DSKeyGen()

	errorPKE := userlib.KeystoreSet(username + " " + "public key", pk) //storing public key in Keystore
	errorDS := userlib.KeystoreSet(username + " " + "verify key", vk) //storing verify key in Keystore

	//error if username already exists
	if errorPKE != nil {
		return nil, errorPKE
	} else if errorDS != nil {
		return nil, errorDS
	}

	var macUsername, _ = userlib.HMACEval(macKey, []byte(username)) //Hash (MAC) username so that we can use bytesToUUID


	//Storing username, secret key, and signing key in user struct
	userdata.Username = username
	userdata.Password = password
	userdata.DSSignKey = dssk
	userdata.SecretKeyEnc = sk

	// TODO: Make empty File map? For now, assuming we add file map in StoreFile.
	//userdata.FileMap = make(map[string][]byte)

	var UUID = bytesToUUID(macUsername)


	//Marshal the userdata struct, so it's JSON encoded.
	var data, _ = json.Marshal(userdata)
	//encrypt user data
	var encryptedData = userlib.SymEnc(symmetricKey, userlib.RandomBytes(16), data)
	//mac user data
	var MAC, _ = userlib.HMACEval(macKey, encryptedData)


	var dataPlusMAC =  append(encryptedData[:], MAC[:]...) //appending MAC to encrypted user struct
	//print(len(x) == len(encryptedData) + len(MAC))
	//println("TOTAL:")
	//println(x)

	userlib.DatastoreSet(UUID, dataPlusMAC)
	//********************************** END OF NEW IMPLEMENTATION *****************************************************************
	return &userdata, nil
}



// This fetches the user information from the Datastore.  It should
// fail with an error if the user/password is invalid, or if the user
// data was corrupted, or if the user can't be found.
func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata

	/*var all_zeroes = make([]byte, 16)

	var macUsername, _ = userlib.HMACEval(all_zeroes, []byte(username)) //Hash (MAC) username so that we can use bytesToUUID
	userdata.MacUsername = macUsername
	var UUID = bytesToUUID(macUsername)

	//var hkdfKey = userlib.Argon2Key([]byte (password), saltHKDF, 32) //key generated to generate more keys using HKDF
	//var macKey, _ = userlib.HashKDF(hkdfKey, []byte("mac")) //MAC key used for MAC-ing other keys

	var userStruct, ok = userlib.DatastoreGet(UUID)
	if !ok {
		return nil, errors.New("username does not exist")
	}

	_ = json.Unmarshal(userStruct, &userdata)
	print("username: " + userdata.Username)

	//var userStruct_username = userStruct.Username
	//
	//if username != userStruct. or
*/
	//******************* START OF NEW IMPLEMENTATION ********************************************************************
	var hkdfKey = userlib.Argon2Key([]byte (password), []byte(username), 32) //key generated to generate more keys using HKDF

	var macKey, _ = userlib.HashKDF(hkdfKey, []byte("mac")) //mac key to recover UUID
	macKey = macKey[:16]

	var symmKey, _ = userlib.HashKDF(hkdfKey, []byte("symmetric key"))
	symmKey = symmKey[:16]

	var macUsername, _ = userlib.HMACEval(macKey, []byte(username)) //Hash (MAC) username so that we can use bytesToUUID
	var UUID = bytesToUUID(macUsername)
	var data, ok = userlib.DatastoreGet(UUID)
	if ok == false {
		return nil, errors.New("Username/Password invalid")
	}
	//TODO: MAKE SURE Math is right
	var ciphertext = data[0: len(data) - 64]
	var macRec = data[len(data) - 64: ] //Mac received from Datastore
	var macComp, _ = userlib.HMACEval(macKey, ciphertext) //recompute MAC

	if userlib.HMACEqual(macRec, macComp) == false {
		return nil, errors.New("data corrupted")
	}

	//Data is now verified, can decrypt data

	var decryptedData = userlib.SymDec(symmKey, ciphertext)
	_ = json.Unmarshal(decryptedData, userdataptr)
	//******************* END OF NEW IMPLEMENTATION ********************************************************************

	return userdataptr, nil
}

// This stores a file in the datastore.
//
// The plaintext of the filename + the plaintext and length of the filename 
// should NOT be revealed to the datastore!
func (userdata *User) StoreFile(filename string, data []byte) {

	//TODO: This is a toy implementation.
	//We will hash the filename using SHA256 to hide the filename length and then
	// encrypt the file contents and sign this file using the private key generated
	// from the Digital Signature library of the person who is uploading it.

	//encrypt and MAC the file

	//if err != nil {
	//
	//} else {
	//	if userData.
	//}


	username := []byte(userdata.Username)
	password := []byte(userdata.Password)
	fileMap := make(map[string][]byte)

	//make mac key
	var hkdfKey = userlib.Argon2Key(password, username, 32)
	var macKey, _ = userlib.HashKDF(hkdfKey, []byte("mac"))
	macKey = macKey[:16]
	//make symmetric key
	var symmKey, _ = userlib.HashKDF(hkdfKey, []byte("symmetric key"))
	symmKey = symmKey[:16]

	// Use hash function to hide the filename length.
	zeroKey := make([]byte, 16) //byte array of 16 0's
	macFilename, _ := userlib.HMACEval(zeroKey, []byte(filename)) // HashFunction(filename) = HMAC(0, filename)

	//TODO: If the file has been shared with others, the file must stay shared.
	//Marshal the file contents, so it's JSON encoded.
	var marshalFileData, _ = json.Marshal(data)
	//encrypt file contents
	var encryptedFileData = userlib.SymEnc(symmKey, userlib.RandomBytes(16), marshalFileData)
	//mac file contents
	var fileDataMAC, _ = userlib.HMACEval(macKey, encryptedFileData)
	var encMACFile =  append(encryptedFileData[:], fileDataMAC[:]...) //appending MAC to encrypted file contents
	fileMap[string(macFilename)] = encMACFile
	//update FileMap
	userdata.FileMap = fileMap
	//get user data from DataStore.
	oldUserData, err := GetUser(userdata.Username, userdata.Password)
	if err != nil {
		print("Error occurred. ")
		return
	}
	oldUserData = userdata

	var macUsername, _ = userlib.HMACEval(macKey, username) //Hash (MAC) username so that we can use bytesToUUID
	var UUID = bytesToUUID(macUsername)
	//Marshal the userdata struct, so it's JSON encoded.
	var newUserData, _ = json.Marshal(oldUserData)
	//encrypt user data
	var encryptedData = userlib.SymEnc(symmKey, userlib.RandomBytes(16), newUserData)
	//mac user data
	var MAC, _ = userlib.HMACEval(macKey, encryptedData)
	var dataPlusMAC =  append(encryptedData[:], MAC[:]...) //appending MAC to encrypted user struct
	userlib.DatastoreSet(UUID, dataPlusMAC)
	//End of toy implementation


	return
}

// This adds on to an existing file.
//
// Append should be efficient, you shouldn't rewrite or reencrypt the
// existing file, but only whatever additional information and
// metadata you need.
func (userdata *User) AppendFile(filename string, data []byte) (err error) {
	return
}

// This loads a file from the Datastore.
//
// It should give an error if the file is corrupted in any way.
func (userdata *User) LoadFile(filename string) (data []byte, err error) {

	//TODO: This is a toy implementation.
	UUID, _ := uuid.FromBytes([]byte(filename + userdata.Username)[:16])
	packaged_data, ok := userlib.DatastoreGet(UUID)
	if !ok {
		return nil, errors.New(strings.ToTitle("File not found!"))
	}
	json.Unmarshal(packaged_data, &data)
	return data, nil
	//End of toy implementation

	return
}

// This creates a sharing record, which is a key pointing to something
// in the datastore to share with the recipient.

// This enables the recipient to access the encrypted file as well
// for reading/appending.

// Note that neither the recipient NOR the datastore should gain any
// information about what the sender calls the file.  Only the
// recipient can access the sharing record, and only the recipient
// should be able to know the sender.
func (userdata *User) ShareFile(filename string, recipient string) (
	magic_string string, err error) {

	return
}

// Note recipient's filename can be different from the sender's filename.
// The recipient should not be able to discover the sender's view on
// what the filename even is!  However, the recipient must ensure that
// it is authentically from the sender.
func (userdata *User) ReceiveFile(filename string, sender string,
	magic_string string) error {
	return nil
}

// Removes target user's access.
func (userdata *User) RevokeFile(filename string, target_username string) (err error) {
	return
}

package handlers

import "testing"

func TestInfoMap_GetKey(t *testing.T) {
	i := NewInfoMap()
	k := i.GetKey()
	if k != 1 {
		t.Fatalf("call to initial GetKey did not return 1, but %v", k)
	}

	k = i.GetKey()
	if k != 2 {
		t.Fatalf("second call to GetKey did not return 2, but %v", k)
	}
}

func TestInfoMap_Load_GetKey(t *testing.T) {
	i := NewInfoMap()
	k := i.GetKey()

	password := "angryMonkey"
	i.Store(k, password)

	loadedPassword := i.Load(k)

	if loadedPassword != password {
		t.Fatalf("Load call for an key stored should not be altered, instead got %v", loadedPassword)
	}
}

func TestInfoMap_Load_NoGetKey(t *testing.T) {
	i := NewInfoMap()
	k := i.GetKey()

	hashedPassword := i.Load(k)

	if hashedPassword != "" {
		t.Fatalf("Load call for an key not stored should be empty, not %v", hashedPassword)
	}
}

func TestInfoMap_Store(t *testing.T) {
	//i := NewInfoMap()

}
package smp

import (
	"testing"
)

func TestGoodPerson(t *testing.T) {
	pub := NewPublic()
	msg := []byte("El pueblo unido")

	alice := NewPerson(pub, msg)
	bob := NewPerson(pub, msg)

	alice.FirstKeyReceive(bob.FirstKeySend())
	bob.FirstKeyReceive(alice.FirstKeySend())

	alice.SecondReceive(bob.SecondSend())
	bob.SecondReceive(alice.SecondSend())

	alice.FinalReceive(bob.FinalSend())
	bob.FinalReceive(alice.FinalSend())

	if !bob.Check() {
		t.Fatal()
	}

	if !alice.Check() {
		t.Fatal()
	}
}

func TestBadPerson(t *testing.T) {
	alicesecret := []byte("No pasaran")
	bobsecret := []byte("la tumba de facismo")
	pub := NewPublic()

	alice := NewPerson(pub, alicesecret)
	bob := NewPerson(pub, bobsecret)

	alice.FirstKeyReceive(bob.FirstKeySend())
	bob.FirstKeyReceive(alice.FirstKeySend())

	alice.SecondReceive(bob.SecondSend())
	bob.SecondReceive(alice.SecondSend())

	alice.FinalReceive(bob.FinalSend())
	bob.FinalReceive(alice.FinalSend())

	if bob.Check() {
	   t.Fatal()
	}

	if alice.Check() {
		t.Fatal()
	}
}

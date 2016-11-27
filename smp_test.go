package smp

import (
	"math/big"
	"testing"
)

func TestGood(t *testing.T) {
	secret := []byte("lo primero es ganar la guerra")
	pub := NewPublic()
	alice := NewAlice(pub, secret)
	bob := NewBob(pub, secret)

	g2b, g3b := alice.One()
	g2b, g3b, pb, qb := bob.Two(g2b, g3b)
	pa, qa, ra := alice.Three(g2b, g3b, pb, qb)
	rb, err := bob.Four(pa, qa, ra)
	if err != nil {
		t.Fatal()
	}
	err = alice.Five(rb)
	if err != nil {
		t.Fatal()
	}
}

func TestBad(t *testing.T) {
	alicesecret := []byte("No pasaran")
	bobsecret := []byte("la tumba de facismo")
	pub := NewPublic()
	alice := NewAlice(pub, alicesecret)
	bob := NewBob(pub, bobsecret)

	g2b, g3b := alice.One()
	g2b, g3b, pb, qb := bob.Two(g2b, g3b)
	pa, qa, ra := alice.Three(g2b, g3b, pb, qb)
	_, err := bob.Four(pa, qa, ra)
	if err == nil {
		t.Fatal()
	}
}

func TestBobLies(t *testing.T) {
	alicesecret := []byte("toda la juventud unida")
	bobsecret := []byte("trabaja y lucha por la revolución")
	pub := NewPublic()
	alice := NewAlice(pub, alicesecret)
	bob := NewBob(pub, bobsecret)

	g2b, g3b := alice.One()
	g2b, g3b, pb, qb := bob.Two(g2b, g3b)
	alice.Three(g2b, g3b, pb, qb)

	// bob pretends he does not get an error and lies to alice
	lie := big.NewInt(666)

	err := alice.Five(lie)
	if err == nil {
		t.Fatal()
	}
}

func TestPerson(t *testing.T) {
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

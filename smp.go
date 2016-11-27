// Referenced from
// https://otr.cypherpunks.ca/Protocol-v2-3.1.0.html

package smp

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var primebits int = 1536

func Pow(base, exp, mod *big.Int) *big.Int {
	return new(big.Int).Exp(base, exp, mod)
}

func Mul(a, b, mod *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(a, b), mod)
}

func Div(a, b, mod *big.Int) *big.Int {
	return Mul(a, new(big.Int).ModInverse(b, mod), mod)
}

func Eq(a, b *big.Int) bool {
	res := a.Cmp(b)
	if res == 0 {
		return true
	} else {
		return false
	}
}

func randInt(max *big.Int, err error) (*big.Int, error) {
	if err != nil {
		return nil, err
	}
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type Public struct {
	Prime *big.Int
	Base  *big.Int
}

func NewPublic() *Public {
	prime, err := rand.Prime(rand.Reader, primebits)
	base, err := randInt(prime, err)
	if err != nil {
		panic(err)
	}
	return &Public{Prime: prime, Base: base}
}

type Person struct {
	g1 *big.Int
	p  *big.Int

	secret *big.Int

	g2 *big.Int
	g3 *big.Int

	qa *big.Int
	qb *big.Int
	pa *big.Int
	pb *big.Int
	ra *big.Int
	rb *big.Int
	rab *big.Int

	exp1 *big.Int
	exp2 *big.Int
	exp3 *big.Int
}

func NewPerson(pub *Public, secret []byte) *Person {
	exp1, err := randInt(pub.Prime, nil)
	exp2, err := randInt(pub.Prime, err)
	exp3, err := randInt(pub.Prime, err)
	if err != nil {
		panic(err)
	}
	secretInt := new(big.Int).SetBytes(secret)
	return &Person{g1: pub.Base, p: pub.Prime, exp1: exp1, exp2: exp2, exp3: exp3,
		secret: secretInt}
}

// with person funcs, p is Alice
func (p *Person) FirstKeySend() (one, two *big.Int) {
	one = Pow(p.g1, p.exp1, p.p)
	two = Pow(p.g1, p.exp2, p.p)
	return one, two
}

func (p *Person) FirstKeyReceive(one, two *big.Int) error {
	if one  == big.NewInt(1) || two == big.NewInt(1) {
		return errors.New("Bad DHKE value received.")
	}
	p.g2 = Pow(one, p.exp1, p.p)
	p.g3 = Pow(two, p.exp2, p.p)
	return nil
}

func (p *Person) SecondSend() (pa, qa *big.Int) {
	p.pa = Pow(p.g3, p.exp3, p.p)
	p.qa = Mul(Pow(p.g1, p.exp3, p.p), Pow(p.g2, p.secret, p.p), p.p)
	return p.pa, p.qa
}

func (p *Person) SecondReceive(pb, qb *big.Int) {
	p.pb = pb
	p.qb = qb
}

func (p *Person) FinalSend() (ra *big.Int) {
	p.ra = Pow(Div(p.qa, p.qb, p.p), p.exp2, p.p)
	return p.ra
}

func (p *Person) FinalReceive(rb *big.Int) {
	p.rb = rb
}

// nb: to make the protocol symmetric, we check for both
// rab == pa/pb and rab == pb/pa
// is this dangerous?
func (p *Person) Check() bool {
	p.rab = Pow(p.rb, p.exp2, p.p)
	if Eq(p.rab, Div(p.pa, p.pb, p.p)) || Eq(p.rab, Div(p.pb, p.pa, p.p)) {
		return true
	}
	return false
}

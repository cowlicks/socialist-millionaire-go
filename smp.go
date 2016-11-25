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

	exp1 *big.Int
	exp2 *big.Int
	exp3 *big.Int
}

func NewPublic() *Public {
	prime, err := rand.Prime(rand.Reader, primebits)
	base, err := randInt(prime, err)
	if err != nil {
		panic(err)
	}
	return &Public{Prime: prime, Base: base}
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

type Alice struct {
	*Person

	a2 *big.Int
	a3 *big.Int

	s *big.Int
}

func NewAlice(pub *Public, secret []byte) *Alice {
	person := NewPerson(pub, secret)
	return &Alice{person, person.exp1, person.exp2, person.exp3}
}

type Bob struct {
	*Person

	b2 *big.Int
	b3 *big.Int
	r  *big.Int
}

func NewBob(pub *Public, secret []byte) *Bob {
	person := NewPerson(pub, secret)
	return &Bob{person, person.exp1, person.exp2, person.exp3}
}

func (a *Alice) One() (g2a, g3a *big.Int) {
	g2a = Pow(a.g1, a.a2, a.p)
	g3a = Pow(a.g1, a.a3, a.p)
	return g2a, g3a
}

func (b *Bob) Two(g2a, g3a *big.Int) (g2b, g3b, pb, qb *big.Int) {
	if g2a == big.NewInt(1) {
		panic("shit")
	}
	g2b = Pow(b.g1, b.b2, b.p)
	g3b = Pow(b.g1, b.b3, b.p)
	b.g2 = Pow(g2a, b.b2, b.p)
	b.g3 = Pow(g3a, b.b3, b.p)

	b.pb = Pow(b.g3, b.r, b.p)
	b.qb = Mul(Pow(b.g1, b.r, b.p), Pow(b.g2, b.secret, b.p), b.p)

	return g2b, g3b, b.pb, b.qb
}

func (a *Alice) Three(g2b, g3b, pb, qb *big.Int) (pa, qa, ra *big.Int) {
	a.pb = pb
	a.g2 = Pow(g2b, a.a2, a.p)
	a.g3 = Pow(g3b, a.a3, a.p)

	a.pa = Pow(a.g3, a.s, a.p)
	a.qa = Mul(Pow(a.g1, a.s, a.p), Pow(a.g2, a.secret, a.p), a.p)

	ra = Pow(Div(a.qa, qb, a.p), a.a3, a.p)

	return a.pa, a.qa, ra
}

func (b *Bob) Four(pa, qa, ra *big.Int) (rb *big.Int, err error) {
	//fmt.Println(pa, qa, ra)
	rb = Pow(Div(qa, b.qb, b.p), b.b3, b.p)
	rab := Pow(ra, b.b3, b.p)
	if Eq(rab, Div(pa, b.pb, b.p)) {
		return rb, nil
	}
	return nil, errors.New("bad match")
}

func (a *Alice) Five(rb *big.Int) error {
	rab := Pow(rb, a.a3, a.p)
	if Eq(rab, Div(a.pa, a.pb, a.p)) {
		return nil
	}
	return errors.New("bad match")
}

# socialist-millionaire-go
Socialist Millionaires' Protocol (SMP) implementation in go

This is usually used when two parties want to check they have the same shared
secret over an insecure channel. The protocol cannot be Man-in-the-Middled, and
eavesdroppers cannot determine whether or not Alice and Bob shared the same
secret.

Usage example:
```go


// Public information shared in the exchange
pub := NewPublic()
// The secret message
msg := []byte("El pueblo unido")

// Create the parties
alice := NewPerson(pub, msg)
bob := NewPerson(pub, msg)

// Alice and Bob secure some values with Diffie-Hellman exchanges
alice.FirstKeyReceive(bob.FirstKeySend())
bob.FirstKeyReceive(alice.FirstKeySend())

// They do another exchange of derived values
alice.SecondReceive(bob.SecondSend())
bob.SecondReceive(alice.SecondSend())

// Final exchange
alice.FinalReceive(bob.FinalSend())
bob.FinalReceive(alice.FinalSend())

// They check if their secret messages were the same
fmt.Println(bob.Check())   // true
fmt.Println(alice.Check()  // true

```

package hashpwd

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	raw    string
	hashed string
	cost   int
	err    error
}

// uses default cost
func New(rawPwd string) *BcryptHasher {
	h := &BcryptHasher{raw: rawPwd}
	h.cost = bcrypt.DefaultCost
	return h
}

// set a new cost
func (b *BcryptHasher) SetCost(cost int) {
	b.cost = cost
}

func (b *BcryptHasher) HashPwd() {
	bx, err := bcrypt.GenerateFromPassword([]byte(b.raw), b.cost)
	if err != nil {
		b.err = err
	}
	b.hashed = string(bx)
}

func CompareHashAndPwd(hashed, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}

// return hashed password
func (b *BcryptHasher) Hashed() string {
	return b.hashed
}

// return the error encountered during hashing
func (b *BcryptHasher) Error() error {
	return b.err
}

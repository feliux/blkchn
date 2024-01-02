package wallet

import "fmt"

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}

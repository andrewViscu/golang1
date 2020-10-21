package middleware
 
import (
	"unicode"
	"golang.org/x/crypto/bcrypt"
)
 
// Password validates plain password against the rules defined below.
//
// upp: at least one upper case letter.
// low: at least one lower case letter.
// num: at least one digit.
// tot: at least eight characters long.
// No empty string or whitespace.

func Password(pass string) bool {
	var (
		upp, low, num bool
		tot           uint8
	)
 
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		default:
			return false
		}
	}
 
	if !upp || !low || !num || tot < 8 {
		return false 
	}
 
	return true
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
package util

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func TemplateFunc(ts int) string {
	t := time.Unix(int64(ts), 0)
	return t.Format("2006-01-02 15:04:05")
}
func CreateUUIDStr() string {
	return uuid.NewString()
}
func CreateUUIDBin() uuid.UUID {
	return uuid.New()
}

// func CreateUUID() (uuid string) {
// 	u := new([16]byte)
// 	_, err := rand.Read(u[:])
// 	if err != nil {
// 		log.Fatalln("Cannot generate UUID", err)
// 	}

// 	// 0x40 is reserved variant from RFC 4122
// 	u[8] = (u[8] | 0x40) & 0x7F
// 	// Set the four most significant bits (bits 12 through 15) of the
// 	// time_hi_and_version field to the 4-bit version number.
// 	u[6] = (u[6] & 0xF) | (0x4 << 4)
// 	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
// 	return
// }

func UUIDStrToBin(uuidStr string) (b uuid.UUID) {
	t, err := uuid.Parse(uuidStr)
	if err != nil {
		return
	}
	return t
}
func BINToUUIDStr(u uuid.UUID) string {
	return u.String()
}

func CreateSalt() (salt string) {
	// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 所有可能的字符集合
	b := make([]byte, 64)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func IsValidEmail(email string) bool {
	// 基本的电子邮件地址验证规则
	// 匹配 alphanumeric 和 !#$%&'*+-/=?^_`{|}~ 字符，以及点号 . 和 @
	// 并要求至少一个字符前有点号，例如 a@b.c
	// 同时，@ 后面至少有一个点号，例如 a@b.c，不能是 a@b 或 a@b.c.d
	// 对于 TLD（顶级域名），至少有两个字符，例如 .com、.org 等
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return emailRegex.MatchString(email)
}

func IsValidMD5(s string) bool {
	matched, err := regexp.MatchString(`^[a-f0-9]{32}$`, s)
	return err == nil && matched
}

func IsValidSalt(s string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]{64}$`, s)
	return err == nil && matched
}

func IsValidUUIDSTR(s string) bool {
	return uuid.Validate(s) == nil

	// // UUID的正则表达式模式
	// uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// // 检查字符串是否符合UUID的正则表达式模式
	// if !uuidPattern.MatchString(s) {
	// 	return false
	// }
	// return true
}
func IsValidUUID(b uuid.UUID) bool {
	return len(b) == 16
}

func RandomName() string {
	data, err := os.ReadFile("./utils/corpus3.txt")
	if err != nil {
		fmt.Println("file corpus3.txt open failed")
		return ""
	}
	lines := strings.Split(string(data), "\n")
	return lines[rand.Intn(len(lines))]
}
func RandomEmail() string {
	email := make([]byte, 14)
	for i := range 10 {
		email[i] = chars[rand.Intn(len(chars))]
	}
	email[10] = '@'
	for i := range 3 {
		email[i+11] = chars[rand.Intn(len(chars))]
	}
	domains := []string{".com", ".org", ".test", "top", ".club"}
	return string(email) + domains[rand.Intn(len(domains))]

}

func RandomContent() string {
	data, err := os.ReadFile("./utils/corpus1.txt")
	if err != nil {
		fmt.Println("file corpus1.txt open failed")
		return ""
	}
	lines := strings.Split(string(data), "\n")
	r_scale := rand.Intn(10) + 1
	res := ""
	for range r_scale {
		res += lines[rand.Intn(len(lines))]
	}
	return res
}

func RandomThread() (title, content string) {
	data, err := os.ReadFile("./utils/corpus2.txt")
	if err != nil {
		fmt.Println("file corpus2.txt open failed")
		return
	}
	lines := strings.Split(string(data), "\n")
	randomLine := lines[rand.Intn(len(lines))]
	parts := strings.Split(randomLine, "##")
	title = parts[0] + RandomContent()
	content = parts[1] + RandomContent()
	return
}

func IfElse(condition bool, trueExpr, falseExpr interface{}) interface{} {
	if condition {
		return trueExpr
	}
	return falseExpr
}

func MapKeys[K, V comparable](m map[K]V) (keys []K) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

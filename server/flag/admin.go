package flag

import (
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"server/global"
	"server/model/apptypes"
	"server/model/database"
	"server/utils"
	"syscall"
)

// Admin 创建一个管理员用户
func Admin() error {
	var user database.User
	fmt.Print("Enter email: ")

	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return fmt.Errorf("failed to read email: %v", err)
	}
	user.Email = email

	// 获取标准输入的fd
	fd := int(syscall.Stdin)
	// 关闭回显，确保密码不会在终端显示
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	fmt.Print("Enter password: ")
	password, err := readPassword()
	fmt.Println()
	if err != nil {
		return err
	}

	fmt.Print("Confirm password: ")
	rePassword, err := readPassword()
	fmt.Println()
	if err != nil {
		return err
	}

	if password != rePassword {
		return fmt.Errorf("passwords do not match")
	}

	if len(password) < 8 || len(password) > 32 {
		return fmt.Errorf("password length must be between 8 and 32")
	}
	user.Password = utils.BcryptHash(password)
	user.Username = global.Config.Website.Name
	user.UUID = uuid.Must(uuid.NewV4())
	user.RoleID = apptypes.Admin
	user.Avatar = "/images/avatar.png"
	user.Address = global.Config.Website.Address

	if err := global.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func readPassword() (string, error) {
	var password string
	var buf [1]byte

	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			return "", err
		}
		char := buf[0]
		if char == '\n' || char == '\r' {
			break
		}
		password += string(char)
	}

	return password, nil
}

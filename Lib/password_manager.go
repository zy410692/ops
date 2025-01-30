package Lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type PasswordEntry struct {
	Site     string    `json:"site"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

type PasswordManager struct {
	entries []PasswordEntry
	key     []byte
}

func NewPasswordManager() *PasswordManager {
	// 使用固定的加密密钥（在实际应用中应该更安全地存储）
	key := []byte("123456789012345678901234567890zy")
	return &PasswordManager{
		entries: make([]PasswordEntry, 0),
		key:     key,
	}
}

func (pm *PasswordManager) AddPassword(site string, password ...string) error {
	// 检查是否已存在
	if exists, _ := pm.VerifyPassword(site); exists {
		return fmt.Errorf("网站 %s 的密码已存在，请使用update命令更新", site)
	}

	var newPassword string
	// 根据参数数量决定使用自动生成的密码还是指定的密码
	if len(password) == 0 {
		newPassword = GeneratePassword(8)
	} else if len(password) == 1 {
		newPassword = password[0]
	} else {
		return fmt.Errorf("参数错误：只能提供0个或1个密码参数")
	}

	now := time.Now()
	entry := PasswordEntry{
		Site:     site,
		Password: newPassword,
		Created:  now,
		Modified: now,
	}

	pm.entries = append(pm.entries, entry)
	return pm.saveToFile()
}

func (pm *PasswordManager) ListPasswords() []PasswordEntry {
	return pm.entries
}

func (pm *PasswordManager) saveToFile() error {
	data, err := json.Marshal(pm.entries)
	if err != nil {
		return err
	}

	// 加密数据
	encrypted, err := pm.encrypt(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("passwords.enc", []byte(encrypted), 0644)
}

func (pm *PasswordManager) LoadFromFile() error {
	data, err := ioutil.ReadFile("passwords.enc")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// 解密数据
	decrypted, err := pm.decrypt(string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal(decrypted, &pm.entries)
}

func (pm *PasswordManager) encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(pm.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (pm *PasswordManager) decrypt(encrypted string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(pm.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("密文太短")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// 删除密码
func (pm *PasswordManager) DeletePassword(site string) error {
	found := false
	newEntries := make([]PasswordEntry, 0)

	for _, entry := range pm.entries {
		if entry.Site != site {
			newEntries = append(newEntries, entry)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("未找到网站 %s 的密码", site)
	}

	pm.entries = newEntries
	return pm.saveToFile()
}

// 修改密码
func (pm *PasswordManager) UpdatePassword(site string, password ...string) error {
	found := false
	var newPassword string

	// 根据参数数量决定使用自动生成的密码还是指定的密码
	if len(password) == 0 {
		newPassword = GeneratePassword(8)
	} else if len(password) == 1 {
		newPassword = password[0]
	} else {
		return fmt.Errorf("参数错误：只能提供0个或1个密码参数")
	}

	for i := range pm.entries {
		if pm.entries[i].Site == site {
			pm.entries[i].Password = newPassword
			pm.entries[i].Modified = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到网站 %s 的密码", site)
	}

	return pm.saveToFile()
}

// 验证密码是否存在
func (pm *PasswordManager) VerifyPassword(site string) (bool, *PasswordEntry) {
	for _, entry := range pm.entries {
		if entry.Site == site {
			return true, &entry
		}
	}
	return false, nil
}

// DeleteAllPasswords 删除所有密码
func (pm *PasswordManager) DeleteAllPasswords() error {
	// 清空密码列表
	pm.entries = make([]PasswordEntry, 0)
	// 保存更改到文件
	return pm.saveToFile()
}

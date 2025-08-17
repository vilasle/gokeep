package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vilasle/gokeep/internal/repository/client"
)

type view interface {
	View() string
}

type loginPasswordView struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (v *loginPasswordView) View() string {
	return fmt.Sprintf("Login: %s, Password: %s", v.Login, v.Password)
}

type bankCardView struct {
	Number  string `json:"number"`
	Expires string `json:"expires"`
	CVV     string `json:"cvv"`
}

func (v *bankCardView) View() string {
	return fmt.Sprintf("Number: %s, Expires: %s, CVV: %s", v.Number, v.Expires, v.CVV)
}
func (c *Client) showList(ls []client.GetResponse) {
	layout := "[ %d ] ID: %d, Description: %s\n"
	for i, data := range ls {
		fmt.Printf(layout, i+1, data.ID, data.View)
	}
}

func (c *Client) showFullEntity(data client.GetResponse, tData client.PrivateDataType) error {
	content, err := c.decrypt(data.DEK, data.Data)
	if err != nil {
		return err
	}

	if tData == client.TypeTextData ||
		tData == client.TypeBinaryData {
		path := filepath.Join(c.workspace.UploadDirectory.Path, data.View)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
		fmt.Printf("File %s saved to %s\n", data.View, path)
		return nil
	}

	var v view
	switch tData {
	case client.TypeLoginPassword:
		v = &loginPasswordView{}
	case client.TypeBankCard:
		v = &bankCardView{}
	default:
		return fmt.Errorf("unknown type of data: %s", tData)
	}

	if err := json.Unmarshal(content, v); err != nil {
		return err
	}

	fmt.Println(v.View())

	return nil
}

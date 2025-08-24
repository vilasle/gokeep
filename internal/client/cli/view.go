package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/repository/client"
)

type view interface {
	View() string
}

type loginPasswordView struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (v *loginPasswordView) View() string {
	return fmt.Sprintf("ID: %d\nLogin: %s\nPassword: %s\n", v.ID, v.Login, v.Password)
}

type bankCardView struct {
	ID      int    `json:"-"`
	Number  string `json:"number"`
	Expires string `json:"expires"`
	CVV     int    `json:"cvv"`
}

func (v *bankCardView) View() string {
	return fmt.Sprintf("ID: %d\nNumber: %s\nExpires: %s\nCVV: %d\n", v.ID, v.Number, v.Expires, v.CVV)
}

func (c *Client) showFullEntity(data client.GetResponse, tData model.Type) error {
	content, err := c.decrypt(data.DEK, data.Data)
	if err != nil {
		return err
	}

	if tData == model.TypePlainText ||
		tData == model.TypeBinaryData {
		path := filepath.Join(c.workspace.UploadDirectory.Path, data.View)
		fd, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fd.Close()
		if _, err := fd.Write(content); err != nil {
			return err
		}

		fmt.Printf("File %s saved to %s\n", data.View, path)
		return nil
	}

	var v view
	switch tData {
	case model.TypeUsepass:
		v = &loginPasswordView{ID: data.ID}
	case model.TypeBankCard:
		v = &bankCardView{ID: data.ID}
	default:
		return fmt.Errorf("unknown type of data: %s", tData)
	}

	if err := json.Unmarshal(content, v); err != nil {
		return err
	}

	fmt.Println(v.View())

	return nil
}

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	CVV     string `json:"cvv"`
}

func (v *bankCardView) View() string {
	return fmt.Sprintf("ID: %d\nNumber: %s\nExpires: %s\nCVV: %s\n", v.ID, v.Number, v.Expires, v.CVV)
}

func (c *Client) showFullEntities(tData model.Type, data ...client.GetResponse) error {
	for _, entity := range data {
		content, err := c.decrypt(entity.DEK, entity.Data)
		if err != nil {
			return err
		}

		if tData == model.TypePlainText ||
			tData == model.TypeBinaryData {
			path := filepath.Join(c.workspace.UploadDirectory.Path, entity.View)
			fd, err := os.Create(path)
			if err != nil {
				return err
			}
			defer fd.Close()
			if _, err := fd.Write(content); err != nil {
				return err
			}

			fmt.Printf("File %s saved to %s\n", entity.View, path)
			continue
		}

		var v view
		switch tData {
		case model.TypeUsepass:
			lp := string(content)
			lps := strings.Split(lp, "\n")
			if len(lps) < 2 {
				return fmt.Errorf("invalid login password data: %s", lp)
			}
			v = &loginPasswordView{ID: entity.ID, Login: lps[0], Password: lps[1]}
		case model.TypeBankCard:
			lp := string(content)
			lps := strings.Split(lp, "\n")
			if len(lps) < 3 {
				return fmt.Errorf("invalid bank card data: %s", lp)
			}
			v = &bankCardView{ID: entity.ID, Number: lps[0], CVV: lps[1], Expires: lps[2]}
		default:
			return fmt.Errorf("unknown type of data: %d", tData)
		}

		fmt.Println(v.View())

	}

	return nil
}

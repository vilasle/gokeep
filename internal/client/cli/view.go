package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/repository/client"
)

type viewer interface {
	View() string
}

type loginPasswordView struct {
	ID       int
	Login    string
	Password string
	Metadata []client.MetadataValue
}

func (v *loginPasswordView) View() string {
	var metadata string
	for _, m := range v.Metadata {
		metadata += fmt.Sprintf("%s=%s;", m.Key, m.Value)
	}

	return fmt.Sprintf("ID: %d\nLogin: %s\nPassword: %s\nMetadata: %s",
		v.ID, v.Login, v.Password, metadata)
}

type bankCardView struct {
	ID       int
	Number   string
	Expires  string
	CVV      string
	Metadata []client.MetadataValue
}

func (v *bankCardView) View() string {
	var metadata string
	for _, m := range v.Metadata {
		metadata += fmt.Sprintf("%s=%s;", m.Key, m.Value)
	}
	return fmt.Sprintf("ID: %d\nNumber: %s\nExpires: %s\nCVV: %s\nMetadata: %s",
		v.ID, v.Number, v.Expires, v.CVV, metadata)
}

type entityView struct {
	ID          int
	Description string
	Metadata    []client.MetadataValue
}

func (v *entityView) View() string {
	var metadata string
	for _, m := range v.Metadata {
		metadata += fmt.Sprintf("%s=%s;", m.Key, m.Value)
	}
	return fmt.Sprintf("ID: %d|Desc: %s|Metadata: %s", v.ID, v.Description, metadata)
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

			fmt.Printf("File '%s' saved to '%s'\n", entity.View, path)
			continue
		}

		var v viewer
		switch tData {
		case model.TypeUsepass:
			lp := string(content)
			lps := strings.Split(lp, "\n")
			if len(lps) < 2 {
				return fmt.Errorf("invalid login password data: %s", lp)
			}
			v = &loginPasswordView{ID: entity.ID, Login: lps[0], Password: lps[1], Metadata: entity.Metadata}
		case model.TypeBankCard:
			lp := string(content)
			lps := strings.Split(lp, "\n")
			if len(lps) < 3 {
				return fmt.Errorf("invalid bank card data: %s", lp)
			}
			v = &bankCardView{ID: entity.ID, Number: lps[0], CVV: lps[1], Expires: lps[2], Metadata: entity.Metadata}
		default:
			return fmt.Errorf("unknown type of data: %d", tData)
		}

		fmt.Println(v.View())

	}

	return nil
}

func (c *Client) showListOfEntities(tData model.Type, data ...client.GetResponse) error {
	for _, entity := range data {
		ev := entityView{ID: entity.ID, Description: entity.View, Metadata: entity.Metadata}
		fmt.Println(ev.View())
	}

	return nil
}

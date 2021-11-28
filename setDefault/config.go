package setDefault

import (
	"barista.run/colors"
	"barista.run/oauth"
	"barista.run/pango/icons/fontawesome"
	"barista.run/pango/icons/material"
	"barista.run/pango/icons/mdi"
	"barista.run/pango/icons/typicons"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func setupOauthEncryption() error {
	const service = "barista-sample-bar"
	var username string
	if u, err := user.Current(); err == nil {
		username = u.Username
	} else {
		username = fmt.Sprintf("user-%d", os.Getuid())
	}
	var secretBytes []byte
	// IMPORTANT: The oauth tokens used by some modules are very sensitive, so
	// we encrypt them with a random key and store that random key using
	// libsecret (gnome-keyring or equivalent). If no secret provider is
	// available, there is no way to store tokens (since the version of
	// sample-bar used for setup-oauth will have a different key from the one
	// running in i3bar). See also https://github.com/zalando/go-keyring#linux.
	secret, err := keyring.Get(service, username)
	if err == nil {
		secretBytes, err = base64.RawURLEncoding.DecodeString(secret)
	}
	if err != nil {
		secretBytes = make([]byte, 64)
		_, err := rand.Read(secretBytes)
		if err != nil {
			return err
		}
		secret = base64.RawURLEncoding.EncodeToString(secretBytes)
		if err = keyring.Set(service, username, secret); err != nil {
			log.Fatal(err)
		}
	}
	oauth.SetEncryptionKey(secretBytes)
	return nil
}

func Config() {
		homePath, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		viper.SetConfigFile(filepath.Join(homePath, ".config/i3/gobar/config.yaml"))
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}

		err = material.Load(filepath.Join(homePath, ".config/i3/gobar/Github/material-design-icons"))
		if err != nil {
			log.Println(err)
		}

		err = mdi.Load(filepath.Join(homePath, ".config/i3/gobar/Github/MaterialDesign-Webfont"))
		if err != nil {
			log.Println(err)
		}

		err = typicons.Load(filepath.Join(homePath, ".config/i3/gobar/Github/typicons.font"))
		if err != nil {
			log.Println(err)
		}

		err = fontawesome.Load(filepath.Join(homePath, ".config/i3/gobar/Github/Font-Awesome"))
		if err != nil {
			log.Println(err)
		}

		colors.LoadBarConfig()
		bg := colors.Scheme("background")
		fg := colors.Scheme("statusline")
		if fg != nil && bg != nil {
			iconColor := fg.Colorful().BlendHcl(bg.Colorful(), 0.5).Clamped()
			colors.Set("dim-icon", iconColor)
			_, _, v := fg.Colorful().Hsv()
			if v < 0.3 {
				v = 0.3
			}
			colors.Set("bad", colorful.Hcl(40, 1.0, v).Clamped())
			colors.Set("degraded", colorful.Hcl(90, 1.0, v).Clamped())
			colors.Set("good", colorful.Hcl(120, 1.0, v).Clamped())
		}

		if err := setupOauthEncryption(); err != nil {
			panic(fmt.Sprintf("Could not setup oauth token encryption: %v", err))
		}
}
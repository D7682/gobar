package main

// Used the simple.go file

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/format"
	"barista.run/modules/clock"
	"barista.run/modules/cputemp"
	"barista.run/modules/diskspace"
	"barista.run/modules/meminfo"
	"barista.run/modules/sysinfo"
	"barista.run/modules/volume"
	"barista.run/modules/volume/alsa"
	"barista.run/modules/weather"
	_ "barista.run/modules/weather"
	"barista.run/modules/weather/openweathermap"
	_ "barista.run/modules/weather/openweathermap"
	"barista.run/oauth"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/fontawesome"
	"barista.run/pango/icons/material"
	"barista.run/pango/icons/mdi"
	"barista.run/pango/icons/typicons"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/martinlindhe/unit"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	startTaskManager  = click.RunLeft("urxvt", "-e", "htop")
	spacer            = pango.Text(" ").XXSmall()
	gsuiteOauthConfig = []byte(`{"installed": {
	"client_id":"%%GOOGLE_CLIENT_ID%%",
	"project_id":"i3-barista",
	"auth_uri":"https://accounts.google.com/o/oauth2/auth",
	"token_uri":"https://www.googleapis.com/oauth2/v3/token",
	"auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
	"client_secret":"%%GOOGLE_CLIENT_SECRET%%",
	"redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
	}}`)
)

func SetConfig() {
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

func GetVolume() (bar.Module, error) {
	volumeModule := volume.New(alsa.DefaultMixer()).Output(func(v volume.Volume) bar.Output {
		var iconName string
		if v.Mute {
			return outputs.Pango(pango.Icon("fas fa-volume-mute"), spacer, "MUT").Color(colors.Scheme("degraded"))
		}
		iconName = "off"
		pct := v.Pct()
		if pct > 66 {
			iconName = "up"
		} else if pct > 33 {
			iconName = "down"
		}

		return outputs.Pango(
			pango.Icon("fas fa-volume-"+iconName),
			spacer,
			pango.Textf("%2d%%", pct),
		)
	})

	return volumeModule, nil
}

func GetDiskSpace() (bar.Module, error) {
	diskspacemod := diskspace.New("/").Output(func(i diskspace.Info) bar.Output {
		return outputs.Textf("%s/%s avail", format.IBytesize(i.Used()), format.IBytesize(i.Total))
	})
	return diskspacemod, nil
}

// GetMem gets the amount of free ram left on the hard drive.
func GetMem() (bar.Module, error) {
	freeMemModule := meminfo.New().Output(func(m meminfo.Info) bar.Output {
		out := outputs.Pango(pango.Icon("material-memory"), format.IBytesize(m.Available()))
		freeGigs := m.Available().Gigabytes()
		switch {
		case freeGigs < 0.5:
			out.Urgent(true)
		case freeGigs < 1:
			out.Color(colors.Scheme("bad"))
		case freeGigs < 2:
			out.Color(colors.Scheme("degraded"))
		case freeGigs > 12:
			out.Color(colors.Scheme("good"))
		}
		out.OnClick(startTaskManager)
		return out
	})
	return freeMemModule, nil
}

// GetLoad is used to get the average cpu load.
func GetLoad() (bar.Module, error) {
	loadAvg := sysinfo.New().Output(func(s sysinfo.Info) bar.Output {
		out := outputs.Textf("%0.2f %0.2f", s.Loads[0], s.Loads[2])
		// Load averages are unusually high for a few minutes after boot.
		if s.Uptime < 10*time.Minute {
			// so don't add colours until 10 minutes after system start.
			return out
		}
		switch {
		case s.Loads[0] > 128, s.Loads[2] > 64:
			out.Urgent(true)
		case s.Loads[0] > 64, s.Loads[2] > 32:
			out.Color(colors.Scheme("bad"))
		case s.Loads[0] > 32, s.Loads[2] > 16:
			out.Color(colors.Scheme("degraded"))
		}
		out.OnClick(startTaskManager)
		return out
	})
	return loadAvg, nil
}
func GetWeather() (bar.Module, error) {
	weatherModule := weather.New(openweathermap.New(viper.GetString("openweather.key")).CityID(viper.GetString("openweather.cityid"))).Output(func(w weather.Weather) bar.Output {
		var iconName string
		switch w.Condition {
		case weather.Thunderstorm,
			weather.TropicalStorm,
			weather.Hurricane:
			iconName = "stormy"
		case weather.Drizzle,
			weather.Hail:
			iconName = "shower"
		case weather.Rain:
			iconName = "downpour"
		case weather.Mist,
			weather.Smoke,
			weather.Whirls,
			weather.Haze,
			weather.Fog:
			iconName = "windy"
		case weather.Clear:
			if !w.Sunset.IsZero() && time.Now().After(w.Sunset) {
				iconName = "night"
			} else {
				iconName = "sunny"
			}
		case weather.PartlyCloudy:
			iconName = "partly-cloudy"
		case weather.Tornado,
			weather.Windy:
			iconName = "windy"
		}

		if iconName == "" {
			iconName = "warning-outline"
		} else {
			iconName = "weather-" + iconName
		}

		return outputs.Pango(
			pango.Icon("typecn-"+iconName), spacer,
			pango.Textf("%.1f°F", w.Temperature.Fahrenheit()),
			pango.Textf(" (provided by %s)", w.Attribution).XSmall(),
		)
	})

	return weatherModule, nil

}

func GetTime() (bar.Module, error) {
	localtime := clock.Local().Output(time.Second, func(now time.Time) bar.Output {
		return outputs.Pango(
			pango.Icon("material-today").Color(colors.Scheme("dim-icon")),
			now.Format(" Mon, Jan-02-2006 "),
			pango.Icon("material-access-time").Color(colors.Scheme("dim-icon")),
			now.Format(" 03:04:05PM"),
		).OnClick(click.RunLeft("gsimplecal"))
	})
	return localtime, nil
}

func GetCpuTemp() (bar.Module, error) {
	cpuModule := cputemp.New().RefreshInterval(2 * time.Second).Output(func(temp unit.Temperature) bar.Output {
		out := outputs.Pango(
			pango.Icon("mdi-fan"), spacer,
			pango.Textf("%2d℃", int(temp.Celsius())),
		)
		switch {
		case temp.Celsius() > 90:
			out.Urgent(true)
		case temp.Celsius() > 70:
			out.Color(colors.Scheme("bad"))
		case temp.Celsius() > 60:
			out.Color(colors.Scheme("degraded"))
		}
		return out
	})
	return cpuModule, nil
}

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
		keyring.Set(service, username, secret)
	}
	oauth.SetEncryptionKey(secretBytes)
	return nil
}

func main() {
	SetConfig()
	modules := []func() (bar.Module, error){
		GetTime,
		GetWeather,
		GetVolume,
		GetMem,
		GetDiskSpace,
		GetCpuTemp,
	}

	for _, module := range modules {
			barModule, err := module()
			if err != nil {
				log.Fatal(err)
			}
			barista.Add(barModule)
	}

	if err := barista.Run(); err != nil {
		log.Fatal(err)
	}
}

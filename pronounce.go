package translate

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const (
	mp3FileName = "./.translate.mp3"
)

func playMp3File(mp3File string) error {
	f, err := os.Open(mp3File)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func Pronounce(words string, useAmerican bool) {
	if words == "" {
		return
	}

	u := &url.URL{
		Scheme: "https",
		Host:   "dict.youdao.com",
		Path:   "/dictvoice",
	}
	q := u.Query()
	q.Add("audio", words)
	if useAmerican {
		q.Add("type", "2")
	} else {
		q.Add("type", "1")
	}
	u.RawQuery = q.Encode()

	http.DefaultClient.Timeout = 4 * time.Second
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatal(string(b))
	}
	if err := os.WriteFile(mp3FileName, b, 0666); err != nil {
		log.Fatal(err)
	}

	if err := playMp3File(mp3FileName); err != nil {
		log.Fatal(err)
	}
}

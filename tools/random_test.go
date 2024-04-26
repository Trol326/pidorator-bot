package tools

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"testing"
	"time"

	randX "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

const (
	playersAmount int = 11
	rollsAmount   int = 100
)

type player struct {
	Name  string
	Score int
}

func Test_GetRandomInt32(t *testing.T) {
	var players = make([]player, 0)
	for i := 0; i < playersAmount; i++ {
		n := fmt.Sprintf("Player%d", i+1)
		p := player{Name: n, Score: 0}
		players = append(players, p)
	}
	var savedPlayers = make([]player, 0)
	savedPlayers = append(savedPlayers, players...)

	random1(players)
	_ = copy(players, savedPlayers)

	random2(players)
	_ = copy(players, savedPlayers)

	random3(players)
	_ = copy(players, savedPlayers)

	random4(players)
	_ = copy(players, savedPlayers)

	random5(players)
	_ = copy(players, savedPlayers)

	random6(players)
	_ = copy(players, savedPlayers)

	random7(players)
	_ = copy(players, savedPlayers)

	random8(players)
	_ = copy(players, savedPlayers)

	random9(players)
	_ = copy(players, savedPlayers)

	random10(players)
	t.Logf("\n\nResult: %v", players)
	_ = copy(players, savedPlayers)

}

// math uniform
func random1(players []player) {
	n := len(players)
	src := randX.NewSource(uint64(time.Now().Unix()))
	dist := distuv.Uniform{Min: 0, Max: float64(n), Src: src}
	for i := 0; i < rollsAmount; i++ {
		num := int(dist.Rand())
		if num >= n {
			num = n - 1
		}
		players[num].Score++
	}

	hist(players, "Математическое равномерное распределение", "Random1")
}

// seed from startup time
func random2(players []player) {
	n := len(players)
	random := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < rollsAmount; i++ {
		num := random.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Сид от времени запуска", "Random2")
}

// seed from current time
func random3(players []player) {
	n := len(players)
	for i := 0; i < rollsAmount; i++ {
		random := rand.New(rand.NewSource(time.Now().Unix() + int64(86400*i)))
		num := random.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Сид от времени крутки", "Random3")
}

// seed from time + const mod
func random4(players []player) {
	n := len(players)
	mod := rand.Int63()
	for i := 0; i < rollsAmount; i++ {
		random := rand.New(rand.NewSource(time.Now().Unix() + int64(86400*i) + mod))
		num := random.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Сид от времени крутки + сгенеренного при запуске числа", "Random4")
}

// seed from time + random mod
func random5(players []player) {
	n := len(players)
	for i := 0; i < rollsAmount; i++ {
		mod := rand.Int63()
		random := rand.New(rand.NewSource(time.Now().Unix() + int64(86400*i) + mod))
		num := random.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Сид от времени крутки + рандомного числа", "Random5")
}

// seed from number generated on startup
func random6(players []player) {
	n := len(players)
	mod := rand.Int63()
	random := rand.New(rand.NewSource(mod))
	for i := 0; i < rollsAmount; i++ {
		num := random.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Сид от числа сгенеренного при запуске", "Random6")
}

// seed from random number
func random7(players []player) {
	n := len(players)
	for i := 0; i < rollsAmount; i++ {
		mod := rand.Int63()
		random := rand.New(rand.NewSource(mod))
		num := random.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Сид от рандомного числа", "Random7")
}

// no seed
func random8(players []player) {
	n := len(players)
	for i := 0; i < rollsAmount; i++ {
		num := rand.Int31n(int32(n))
		players[num].Score++
	}

	hist(players, "Дефолтный рандом(сид неизвестен)", "Random8")
}

// no seed + float + cut
func random9(players []player) {
	n := len(players)
	for i := 0; i < rollsAmount; i++ {
		num := int(rand.Float64() * float64(n))
		if num >= n {
			num = n - 1
		}
		players[num].Score++
	}

	hist(players, "Дефолтный рандом(сид неизвестен) из float с обрезкой", "Random9")
}

// no seed + float + cut
func random10(players []player) {
	n := len(players)
	for i := 0; i < rollsAmount; i++ {
		num := int(math.Round(rand.Float64() * float64(n)))
		if num >= n {
			num = n - 1
		}
		players[num].Score++
	}

	hist(players, "Дефолтный рандом(сид неизвестен) из float с округлением", "Random10")
}

func hist(dist []player, title, name string) {
	n := len(dist)
	vals := make(plotter.XYs, n)
	for i := 0; i < n; i++ {
		vals[i].X = float64(i + 1)
		vals[i].Y = float64(dist[i].Score)
	}

	plt := plot.New()
	plt.Title.Text = fmt.Sprintf("%s\n%s\nPlayers: %d; Rolls: %d", name, title, playersAmount, rollsAmount)
	plt.X.Label.Text = "Player"
	plt.Y.Label.Text = "Score"

	hist, err := plotter.NewHistogram(vals, n)
	if err != nil {
		log.Println("Cannot plot:", err)
	}
	hist.FillColor = color.RGBA{R: 255, G: 127, B: 80, A: 255} // coral color
	plt.Add(hist)

	t := time.Now()
	fileName := fmt.Sprintf("tests/test_%s_%d-%d-%d_%d-%d-%d", name, t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
	err = plt.Save(400, 200, fileName+".png")
	if err != nil {
		log.Panic(err)
	}
}

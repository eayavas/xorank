package logic

import (
	"math"

	"github.com/eayavas/XORank/internal/models"
)

const KFactor = 32.0

func Calculate(winner, loser *models.Item) {
	// Expected score
	// 1 / (1 + 10^((Rb - Ra) / 400))
	expectedWin := 1.0 / (1.0 + math.Pow(10, (loser.Rating-winner.Rating)/400.0))

	winner.Rating += KFactor * (1.0 - expectedWin)
	loser.Rating -= KFactor * (1.0 - expectedWin)

	winner.Wins++
	loser.Losses++
}

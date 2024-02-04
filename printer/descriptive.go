package printer

import (
	"fmt"

	"github.com/marianogappa/cheesse/core"
)

type DescriptivePrinter struct{}

func (p DescriptivePrinter) PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error) {
	return genericGamePrinter(gameSteps, gameCharacteristics, p)
}

func piece(gameStep core.GameStep, gameCharacteristics GameCharacteristics, renderFileIfPawn bool) string {
	return gameStep.StepAction.FromPiece.PieceType.ToDescriptive(gameCharacteristics.descriptiveUseKt != nil && *gameCharacteristics.descriptiveUseKt)
}

func enPassant(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsEnPassantCapture || gameCharacteristics.usesEnPassantSymbol == nil {
		return ""
	}
	return fmt.Sprintf(" %v", *gameCharacteristics.usesEnPassantSymbol)
}

func promotion(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsPromotion || gameCharacteristics.usesPromotionSymbol == nil {
		return ""
	}
	useKt := gameCharacteristics.descriptiveUseKt != nil && *gameCharacteristics.descriptiveUseKt
	switch *gameCharacteristics.usesPromotionSymbol {
	case "Q":
		return gameStep.StepAction.PromotionPieceType.ToDescriptive(useKt)
	case "=":
		return fmt.Sprintf("=%v", gameStep.StepAction.PromotionPieceType.ToDescriptive(useKt))
	case "(":
		return fmt.Sprintf("(%v)", gameStep.StepAction.PromotionPieceType.ToDescriptive(useKt))
	case "/":
		return fmt.Sprintf("/%v", gameStep.StepAction.PromotionPieceType.ToDescriptive(useKt))
	}
	return ""
}

func castle(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsCastle || gameCharacteristics.usesCastlingSymbol == nil {
		return ""
	}
	switch *gameCharacteristics.usesCastlingSymbol {
	case "O-O":
		if gameStep.StepAction.IsKingsideCastle {
			return "O-O"
		}
		return "O-O-O"
	case "0-0":
		if gameStep.StepAction.IsKingsideCastle {
			return "0-0"
		}
		return "0-0-0"
	}
	return ""
}

func check(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	// TODO: Double check
	// TODO: Discover check
	if !gameStep.StepGame.IsCheck || gameStep.StepGame.IsCheckmate || gameCharacteristics.usesCheckSymbol == nil {
		return ""
	}
	return fmt.Sprintf("%v", *gameCharacteristics.usesCheckSymbol)
}

func checkmate(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepGame.IsCheckmate || gameCharacteristics.usesCheckmateSymbol == nil {
		return ""
	}
	return fmt.Sprintf("%v", *gameCharacteristics.usesCheckmateSymbol)
}

func capture(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsCapture {
		return ""
	}
	captureSymbol := ""
	if gameCharacteristics.usesCaptureSymbol != nil {
		captureSymbol = *gameCharacteristics.usesCaptureSymbol
	}
	return fmt.Sprintf(
		"%v%v%v%v%v%v%v",
		piece(gameStep, gameCharacteristics, true),
		captureSymbol,
		gameStep.StepAction.CapturedPiece.XY.ToDescriptive(gameStep.StepGame.Turn()),
		enPassant(gameStep, gameCharacteristics),
		promotion(gameStep, gameCharacteristics),
		algCheck(gameStep, gameCharacteristics),
		algCheckmate(gameStep, gameCharacteristics),
	)
}

func move(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	return fmt.Sprintf(
		"%v-%v%v%v%v%v",
		piece(gameStep, gameCharacteristics, false),
		gameStep.StepAction.ToXY.ToDescriptive(gameStep.StepGame.Turn().Opponent()), // N.B. because game is after doing the action
		enPassant(gameStep, gameCharacteristics),
		promotion(gameStep, gameCharacteristics),
		check(gameStep, gameCharacteristics),
		checkmate(gameStep, gameCharacteristics),
	)
}

func resign(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsResign || gameCharacteristics.usesEndGameSymbol == nil {
		return ""
	}
	return fmt.Sprintf("%v", *gameCharacteristics.usesEndGameSymbol)
}

func (p DescriptivePrinter) PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error) {
	gameCharacteristics = applyDefaultGameCharacteristics(gameCharacteristics)
	// TODO: Disambiguation
	if gameStep.StepAction.IsCastle {
		return castle(gameStep, gameCharacteristics), nil
	}
	if gameStep.StepAction.IsResign {
		return resign(gameStep, gameCharacteristics), nil
	}
	if gameStep.StepAction.IsCapture {
		return capture(gameStep, gameCharacteristics), nil
	}
	return move(gameStep, gameCharacteristics), nil
}

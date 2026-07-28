package printer

import (
	"fmt"

	"github.com/marianogappa/cheesse/core"
)

type AlgebraicPrinter struct{}

func (p AlgebraicPrinter) PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error) {
	return genericGamePrinter(gameSteps, gameCharacteristics, p)
}

func algPiece(gameStep core.GameStep, gameCharacteristics GameCharacteristics, renderFileIfPawn bool) string {
	if !gameCharacteristics.isFigurine && gameStep.StepAction.FromPiece.PieceType == core.PiecePawn && renderFileIfPawn {
		return gameStep.StepAction.FromPiece.XY.ToAlgebraic()[0:1]
	}
	if gameCharacteristics.isFigurine {
		return gameStep.StepAction.FromPiece.PieceType.ToFigurine()
	}
	return gameStep.StepAction.FromPiece.PieceType.ToAlgebraic()
}

func algDisambiguate(gameStep core.GameStep) string {
	action := gameStep.StepAction
	if action.FromPiece.PieceType == core.PiecePawn || action.FromPiece.PieceType == core.PieceKing {
		return ""
	}
	var sameFile, sameRank bool
	ambiguous := false
	for _, other := range gameStep.StepPreMoveGame.Actions {
		if other.IsResign || other.FromPiece.PieceType != action.FromPiece.PieceType || other.ToXY != action.ToXY || other.FromPiece.XY == action.FromPiece.XY {
			continue
		}
		ambiguous = true
		if other.FromPiece.XY.X == action.FromPiece.XY.X {
			sameFile = true
		}
		if other.FromPiece.XY.Y == action.FromPiece.XY.Y {
			sameRank = true
		}
	}
	if !ambiguous {
		return ""
	}
	from := action.FromPiece.XY.ToAlgebraic()
	if sameFile && sameRank {
		return from
	}
	if sameFile {
		return from[1:2]
	}
	return from[0:1]
}

func algEnPassant(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsEnPassantCapture || gameCharacteristics.usesEnPassantSymbol == nil {
		return ""
	}
	return fmt.Sprintf(" %v", *gameCharacteristics.usesEnPassantSymbol)
}

func algPromotion(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsPromotion || gameCharacteristics.usesPromotionSymbol == nil {
		return ""
	}
	switch *gameCharacteristics.usesPromotionSymbol {
	case "Q":
		return gameStep.StepAction.PromotionPieceType.ToAlgebraic()
	case "=":
		return fmt.Sprintf("=%v", gameStep.StepAction.PromotionPieceType.ToAlgebraic())
	case "(":
		return fmt.Sprintf("(%v)", gameStep.StepAction.PromotionPieceType.ToAlgebraic())
	case "/":
		return fmt.Sprintf("/%v", gameStep.StepAction.PromotionPieceType.ToAlgebraic())
	}
	return ""
}

func algCastle(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
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

func algCheck(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	// TODO: Double check
	// TODO: Discover check
	if !gameStep.StepGame.IsCheck || gameStep.StepGame.IsCheckmate || gameCharacteristics.usesCheckSymbol == nil {
		return ""
	}
	return fmt.Sprintf("%v", *gameCharacteristics.usesCheckSymbol)
}

func algCheckmate(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepGame.IsCheckmate || gameCharacteristics.usesCheckmateSymbol == nil {
		return ""
	}
	return fmt.Sprintf("%v", *gameCharacteristics.usesCheckmateSymbol)
}

func algCapture(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsCapture {
		return ""
	}
	captureSymbol := ""
	if gameCharacteristics.usesCaptureSymbol != nil {
		captureSymbol = *gameCharacteristics.usesCaptureSymbol
	}
	return fmt.Sprintf(
		"%v%v%v%v%v%v%v%v",
		algPiece(gameStep, gameCharacteristics, true),
		algDisambiguate(gameStep),
		captureSymbol,
		gameStep.StepAction.ToXY.ToAlgebraic(),
		algEnPassant(gameStep, gameCharacteristics),
		algPromotion(gameStep, gameCharacteristics),
		algCheck(gameStep, gameCharacteristics),
		algCheckmate(gameStep, gameCharacteristics),
	)
}

func algMove(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	return fmt.Sprintf(
		"%v%v%v%v%v%v%v",
		algPiece(gameStep, gameCharacteristics, false),
		algDisambiguate(gameStep),
		gameStep.StepAction.ToXY.ToAlgebraic(),
		algEnPassant(gameStep, gameCharacteristics),
		algPromotion(gameStep, gameCharacteristics),
		algCheck(gameStep, gameCharacteristics),
		algCheckmate(gameStep, gameCharacteristics),
	)
}

func algResign(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	if !gameStep.StepAction.IsResign || gameCharacteristics.usesEndGameSymbol == nil {
		return ""
	}
	return fmt.Sprintf("%v", *gameCharacteristics.usesEndGameSymbol)
}

func (p AlgebraicPrinter) PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error) {
	gameCharacteristics = applyDefaultGameCharacteristics(gameCharacteristics)
	if gameStep.StepAction.IsCastle {
		return algCastle(gameStep, gameCharacteristics), nil
	}
	if gameStep.StepAction.IsResign {
		return algResign(gameStep, gameCharacteristics), nil
	}
	if gameStep.StepAction.IsCapture {
		return algCapture(gameStep, gameCharacteristics), nil
	}
	return algMove(gameStep, gameCharacteristics), nil
}

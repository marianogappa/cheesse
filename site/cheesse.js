// cheesse.js — thin wrapper around the cheesse WASM globals.
// Every function JSON-encodes the request, calls the wasm global, and JSON-decodes the response.
const enc = new TextEncoder()
const dec = new TextDecoder()

function cheesse(fn, obj) {
  const raw = obj === undefined ? fn() : fn(enc.encode(JSON.stringify(obj)))
  return JSON.parse(dec.decode(raw))
}

const cheesse_api = {
  defaultGame: () => cheesse(cheesseDefaultGame).game,

  parseGame: (fenString) => cheesse(cheesseParseGame, { game: { fenString } }).game,

  doAction: (fenString, fromSquare, toSquare, promotionPieceType) => {
    const r = cheesse(cheesseDoAction, {
      game: { fenString },
      action: { fromSquare, toSquare, promotionPieceType: promotionPieceType || '' }
    })
    if (r.error) return { success: false, error: r.error }
    return { success: true, game: r.game, action: r.action }
  },

  parseNotation: (notationString, fenString) => {
    const game = fenString ? { fenString } : {}
    const r = cheesse(cheesseParseNotation, { game, notationString })
    if (r.error) return { success: false }
    const pr = r.parseResult
    return {
      success: pr.parseWasSuccessful,
      notationName: pr.notationName,
      parseWasSuccessful: pr.parseWasSuccessful,
      validActionCount: pr.validActionCount,
      actions: (pr.steps || []).map(s => s.actionString),
      boards: (pr.steps || []).map(s => s.game.fenString),
      steps: pr.steps || []
    }
  },

  convertNotation: (notationString, targetNotation, fenString) => {
    const game = fenString ? { fenString } : {}
    const r = cheesse(cheesseConvertNotation, { game, notationString, targetNotation })
    if (r.error) return { success: false }
    const pr = r.parseResult
    return {
      success: pr.parseWasSuccessful,
      notationName: pr.notationName,
      parseWasSuccessful: pr.parseWasSuccessful,
      validActionCount: pr.validActionCount,
      actions: (pr.steps || []).map(s => s.actionString),
      boards: (pr.steps || []).map(s => s.game.fenString),
      steps: pr.steps || []
    }
  },

  aiMove: (fenString, mode) => {
    const r = cheesse(cheesseAIMove, { game: { fenString }, mode })
    if (r.error) return { success: false, error: r.error }
    return {
      success: r.moveAvailable,
      moveAvailable: r.moveAvailable,
      game: r.game,
      action: r.action,
      board: r.game.fenString,
      isCheckmate: r.game.isCheckmate,
      isDraw: r.game.isDraw || r.game.isStalemate,
      isGameOver: r.game.isGameOver
    }
  }
}

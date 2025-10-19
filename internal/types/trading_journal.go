package types

// TradingSession represents the trading session time zones
type TradingSession string

const (
	TradingSessionAsia   TradingSession = "asia"
	TradingSessionLondon TradingSession = "london"
	TradingSessionNewYork TradingSession = "new_york"
)

// IsValid checks if the trading session is valid
func (s TradingSession) IsValid() bool {
	switch s {
	case TradingSessionAsia, TradingSessionLondon, TradingSessionNewYork:
		return true
	}
	return false
}

// TradeType represents the type of trade
type TradeType string

const (
	TradeTypeSwing    TradeType = "swing"
	TradeTypeIntraday TradeType = "intraday"
)

// IsValid checks if the trade type is valid
func (t TradeType) IsValid() bool {
	switch t {
	case TradeTypeSwing, TradeTypeIntraday:
		return true
	}
	return false
}

// TradeDirection represents the direction of the trade
type TradeDirection string

const (
	TradeDirectionBuy  TradeDirection = "buy"
	TradeDirectionSell TradeDirection = "sell"
)

// IsValid checks if the trade direction is valid
func (d TradeDirection) IsValid() bool {
	switch d {
	case TradeDirectionBuy, TradeDirectionSell:
		return true
	}
	return false
}

// EntryType represents the type of entry order
type EntryType string

const (
	EntryTypeMarket EntryType = "market"
	EntryTypeLimit  EntryType = "limit"
)

// IsValid checks if the entry type is valid
func (e EntryType) IsValid() bool {
	switch e {
	case EntryTypeMarket, EntryTypeLimit:
		return true
	}
	return false
}

// TradeResult represents the outcome of a trade
type TradeResult string

const (
	TradeResultTakeProfit TradeResult = "TP"  // Take Profit
	TradeResultStopLoss   TradeResult = "SL"  // Stop Loss
	TradeResultBreakEven  TradeResult = "BE"  // Break Even
)

// IsValid checks if the trade result is valid
func (r TradeResult) IsValid() bool {
	switch r {
	case TradeResultTakeProfit, TradeResultStopLoss, TradeResultBreakEven:
		return true
	}
	return false
}

// TimeFrame represents common forex timeframes
type TimeFrame string

const (
	TimeFrame1M  TimeFrame = "1M"
	TimeFrame5M  TimeFrame = "5M"
	TimeFrame15M TimeFrame = "15M"
	TimeFrame30M TimeFrame = "30M"
	TimeFrame1H  TimeFrame = "1H"
	TimeFrame4H  TimeFrame = "4H"
	TimeFrame1D  TimeFrame = "1D"
	TimeFrame1W  TimeFrame = "1W"
	TimeFrame1MO TimeFrame = "1MO"
)

// IsValid checks if the timeframe is valid
func (tf TimeFrame) IsValid() bool {
	switch tf {
	case TimeFrame1M, TimeFrame5M, TimeFrame15M, TimeFrame30M,
		TimeFrame1H, TimeFrame4H, TimeFrame1D, TimeFrame1W, TimeFrame1MO:
		return true
	}
	return false
}

// CurrencyPair represents common forex currency pairs
type CurrencyPair string

const (
	// Major pairs
	CurrencyPairEURUSD CurrencyPair = "EURUSD"
	CurrencyPairGBPUSD CurrencyPair = "GBPUSD"
	CurrencyPairUSDJPY CurrencyPair = "USDJPY"
	CurrencyPairUSDCHF CurrencyPair = "USDCHF"
	CurrencyPairAUDUSD CurrencyPair = "AUDUSD"
	CurrencyPairUSDCAD CurrencyPair = "USDCAD"
	CurrencyPairNZDUSD CurrencyPair = "NZDUSD"

	// Minor pairs
	CurrencyPairEURGBP CurrencyPair = "EURGBP"
	CurrencyPairEURJPY CurrencyPair = "EURJPY"
	CurrencyPairGBPJPY CurrencyPair = "GBPJPY"
	CurrencyPairEURCHF CurrencyPair = "EURCHF"
	CurrencyPairEURAUD CurrencyPair = "EURAUD"
	CurrencyPairEURCAD CurrencyPair = "EURCAD"
	CurrencyPairGBPCHF CurrencyPair = "GBPCHF"
	CurrencyPairGBPAUD CurrencyPair = "GBPAUD"
	CurrencyPairGBPCAD CurrencyPair = "GBPCAD"

	// Exotic pairs
	CurrencyPairUSDTRY CurrencyPair = "USDTRY"
	CurrencyPairUSDMXN CurrencyPair = "USDMXN"
	CurrencyPairUSDZAR CurrencyPair = "USDZAR"
	CurrencyPairUSDNOK CurrencyPair = "USDNOK"
	CurrencyPairUSDSEK CurrencyPair = "USDSEK"
)

// IsValid checks if the currency pair is valid
func (cp CurrencyPair) IsValid() bool {
	switch cp {
	case CurrencyPairEURUSD, CurrencyPairGBPUSD, CurrencyPairUSDJPY, CurrencyPairUSDCHF,
		CurrencyPairAUDUSD, CurrencyPairUSDCAD, CurrencyPairNZDUSD,
		CurrencyPairEURGBP, CurrencyPairEURJPY, CurrencyPairGBPJPY, CurrencyPairEURCHF,
		CurrencyPairEURAUD, CurrencyPairEURCAD, CurrencyPairGBPCHF, CurrencyPairGBPAUD,
		CurrencyPairGBPCAD, CurrencyPairUSDTRY, CurrencyPairUSDMXN, CurrencyPairUSDZAR,
		CurrencyPairUSDNOK, CurrencyPairUSDSEK:
		return true
	}
	return false
}

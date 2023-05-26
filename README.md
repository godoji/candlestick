# Candlestick

This repository provides convenient data structures for representing financial market candlestick and indicator data in Go. These structures are designed to be easily serialized and deserialized, making it convenient for storing and transferring data.

## Installation

To include this package in your Go project, use the `go get` command:

```shell
go get github.com/godoji/candlestick
```

## Overview

### Candles

- `Candle`: Represents a candlestick in a financial chart. The fields include:
    - `Open`: The opening price.
    - `High`: The highest price during the time period.
    - `Low`: The lowest price during the time period.
    - `Close`: The closing price.
    - `Volume`: The trading volume during the time period.
    - `TakerVolume`: The volume contributed by takers.
    - `NumberOfTrades`: The number of trades occurred in the time period.
    - `Time`: The time when the candle was observed.
    - `Missing`: A boolean indicating if data for the candle is missing.
- `DataSetMeta`: Contains metadata about a dataset. The fields include:
    - `UID`: The unique identifier for the dataset.
    - `Block`: The block number of the dataset.
    - `Complete`: A boolean indicating if the dataset is complete.
    - `LastUpdate`: The time of the last update to the dataset.
    - `Symbol`: The symbol for the asset in the dataset.
    - `Interval`: The time interval for the candles in the dataset.
- `CandleSet`: Represents a set of candlesticks. The fields include:
    - `Candles`: An array of Candle structs.
    - `Meta`: A DataSetMeta struct containing metadata about the dataset.

### Indicators

- `SeriesType`: Defines the type of series as a string. Possible values include "BAR_CHART", "LINE_CHART", and "PATH_CHART".

- `AxisType`: Defines the type of axis as a string. Possible values include "PRICE_AXIS" and "CUSTOM_AXIS".

- `Indicator`: Represents a set of indicator series. The fields include:
    - `Series`: A map where each key is a string and each value is a pointer to an IndicatorSeries.
    - `Meta`: An IndicatorMeta containing metadata about the indicator.

- `IndicatorValue`: Represents an indicator value. The fields include:
    - `Value`: The value of the indicator.
    - `Missing`: A boolean indicating if the indicator value is missing.

- `IndicatorSeries`: Represents a series of indicator values. The fields include:
    - `Values`: An array of IndicatorValue structs.
    - `Kind`: The kind of series, represented as a SeriesType.
    - `Axis`: The axis type for the series, represented as an AxisType.

- `IndicatorMeta`: Contains metadata about an indicator. The fields include:
    - `UID`: The unique identifier for the indicator.
    - `Block`: The block number of the indicator.
    - `Complete`: A boolean indicating if the indicator is complete.
    - `LastUpdate`: The time of the last update to the indicator.
    - `Symbol`: The symbol for the asset in the indicator.
    - `Interval`: The time interval for the indicator.
    - `BaseInterval`: The base time interval for the indicator.
    - `Name`: The name of the indicator.
    - `Parameters`: A slice of integers representing parameters for the indicator.

### Exchange and Asset Info

- `AssetInfo` struct: Represents information about a trading asset. It includes the following fields:
    - `Symbol`: The symbol of the asset.
    - `Identifier`: The identifier of the asset.
    - `PairInfo`: Pair information for the base and quote assets.
    - `PrecisionInfo`: Precision information for the base and quote assets.
    - `TradingConstraints`: Constraints applicable for trading the asset.
    - `OnboardingDate`: The onboarding date for the asset.
    - `SplitHistory`: Any split history for the asset.

- `AssetSplit` struct: Represents a split in the asset. It includes the following fields:
    - `Time`: The time of the split.
    - `SplitRatio`: The split ratio.

- `TradeConstraints` struct: Represents the constraints applicable for trading the asset. It includes the following fields:
    - `MaxPrice`: The maximum price.
    - `MinPrice`: The minimum price.
    - `TickSize`: The tick size (smallest increment by which the price can change).
    - `QuantityRestrictions`: Quantity restrictions.
    - `MinNotionalValue`: The minimum notional value of a trade.

- `ExchangeInfo` struct: Represents information about an exchange. It includes the following fields:
    - `ExchangeName`: The name of the exchange.
    - `ExchangeID`: The ID of the exchange.
    - `BrokerID`: The ID of the broker.
    - `LastUpdateTimestamp`: The last update timestamp.
    - `TradingSymbols`: The available trading symbols.
    - `Resolutions`: The resolutions available for chart data.

- `BrokerInfo` struct: Represents basic information about a broker. It includes the following field:
    - `BrokerName`: The name of the broker.

- `ExchangeList` struct: Represents a list of exchanges and brokers. It includes the following fields:
    - `Exchanges`: A slice of `ExchangeInfo` objects.
    - `BrokerInfo`: A map with broker names as keys and `BrokerInfo` objects as values.

### Trade and Order Info

- `PositionSide` type: A type alias for string used to represent the position side. It can be either "LONG" (buying with the expectation that the asset will increase in value) or "SHORT" (selling with the expectation that the asset will decrease in value).

- `OrderSide` type: A type alias for string used to represent the order side. It can be either "BUY" (indicating a buy order) or "SELL" (indicating a sell order).

- `OrderKind` type: A type alias for string used to represent the kind of order. It can be "SL" (Stop Loss), "TP" (Take Profit), "LIMIT" (Limit order), or "MARKET" (Market order).

- `Position` struct: Represents a position in a trading context. It includes the following fields:
    - `EntryPrice`: The entry price.
    - `Amount`: The amount of the asset.
    - `Symbol`: The symbol of the asset.
    - `LastUpdateTimestamp`: The last update timestamp.
    - `LastPrice`: The last price of the asset.

- `Order` struct: Represents a trading order. It includes the following fields:
    - `OrderID`: The order ID.
    - `Symbol`: The symbol of the order.
    - `Price`: The price of the order.
    - `Side`: The order side.
    - `Kind`: The kind of order.
    - `PositionSide`: The position side.
    - `CloseOrder`: A boolean indicating whether it's a close

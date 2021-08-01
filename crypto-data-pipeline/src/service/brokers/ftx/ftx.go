package ftx

import (
	"context"
	"fmt"

	"github.com/go-numb/go-ftx/realtime"
)

func GetTickerData(ticker []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"ticker"}, ticker, nil)

	for {
		v := <-ch
		fmt.Printf("%s	%+v\n", v.Symbol, v.Ticker)
	}
}

package bench_test

import (
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/unedtamps/orbook/bench"
	"github.com/unedtamps/orbook/book"
)

func BenchmarkFastBook(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const batchSize = 4096

	b.Run("AddLimit", func(b *testing.B) {
		rng := rand.New(rand.NewSource(101))
		reqs := make([]book.OrderRequest, batchSize)
		for i := 0; i < batchSize; i++ {
			side := book.Buy
			offset := int64(1 + rng.Intn(200))
			price := mid - offset
			if i%2 == 1 {
				side = book.Sell
				price = mid + offset
			}
			reqs[i] = book.OrderRequest{
				ID: book.OrderID(10_000 + i), TraderID: book.TraderID(1 + i%16),
				Side: side, Type: book.Limit, Price: price, Qty: 100,
			}
		}
		bk, err := book.NewFastBook(cfg, batchSize+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k*2), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k*2), Qty: 100}, evBuf[:0])
		}
		for i := 0; i < batchSize; i++ {
			evBuf = bk.PlaceOrder(reqs[i], evBuf[:0])
		}
		for i := 0; i < batchSize; i++ {
			evBuf = bk.CancelOrder(reqs[i].ID, evBuf[:0])
		}
		idx := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == batchSize {
				b.StopTimer()
				for j := 0; j < batchSize; j++ {
					evBuf = bk.CancelOrder(reqs[j].ID, evBuf[:0])
				}
				idx = 0
				b.StartTimer()
			}
			evBuf = bk.PlaceOrder(reqs[idx], evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("Cancel", func(b *testing.B) {
		rng := rand.New(rand.NewSource(202))
		reqs := make([]book.OrderRequest, batchSize)
		for i := 0; i < batchSize; i++ {
			side := book.Buy
			offset := int64(1 + rng.Intn(200))
			price := mid - offset
			if i%2 == 1 {
				side = book.Sell
				price = mid + offset
			}
			reqs[i] = book.OrderRequest{
				ID: book.OrderID(10_000 + i), TraderID: book.TraderID(1 + i%16),
				Side: side, Type: book.Limit, Price: price, Qty: 100,
			}
		}
		cancelOrder := make([]book.OrderID, batchSize)
		for i, p := range rng.Perm(batchSize) {
			cancelOrder[i] = reqs[p].ID
		}
		bk, err := book.NewFastBook(cfg, batchSize+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k*2), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k*2), Qty: 100}, evBuf[:0])
		}
		for i := 0; i < batchSize; i++ {
			evBuf = bk.PlaceOrder(reqs[i], evBuf[:0])
		}
		idx := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == batchSize {
				b.StopTimer()
				for j := 0; j < batchSize; j++ {
					evBuf = bk.PlaceOrder(reqs[j], evBuf[:0])
				}
				idx = 0
				b.StartTimer()
			}
			evBuf = bk.CancelOrder(cancelOrder[idx], evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("MatchSingle", func(b *testing.B) {
		bk, err := book.NewFastBook(cfg, batchSize+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		seed := func() {
			for j := 0; j < batchSize; j++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1 + j), TraderID: 1, Side: book.Sell, Type: book.Limit, Price: mid, Qty: 100}, evBuf[:0])
			}
		}
		seed()
		idx := 0
		var takerID book.OrderID = 100_000
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == batchSize {
				b.StopTimer()
				seed()
				idx = 0
				b.StartTimer()
			}
			takerID++
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: takerID, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: mid, Qty: 100}, evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("MatchSweep", func(b *testing.B) {
		const sweepLevels, numSweeps = 10, 500
		const totalAsks = sweepLevels * numSweeps
		bk, err := book.NewFastBook(cfg, totalAsks+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 64)
		seed := func() {
			for j := 0; j < totalAsks; j++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1 + j), TraderID: 1, Side: book.Sell, Type: book.Limit, Price: mid + int64(j), Qty: 100}, evBuf[:0])
			}
		}
		seed()
		idx := 0
		var takerID book.OrderID = 1_000_000
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == numSweeps {
				b.StopTimer()
				seed()
				idx = 0
				b.StartTimer()
			}
			takerID++
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: takerID, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: mid + int64(totalAsks), Qty: sweepLevels * 100}, evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("Market", func(b *testing.B) {
		const sweepLevels, numSweeps = 10, 500
		const totalAsks = sweepLevels * numSweeps
		bk, err := book.NewFastBook(cfg, totalAsks+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 64)
		seed := func() {
			for j := 0; j < totalAsks; j++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1 + j), TraderID: 1, Side: book.Sell, Type: book.Limit, Price: mid + int64(j), Qty: 100}, evBuf[:0])
			}
		}
		seed()
		idx := 0
		var takerID book.OrderID = 1_000_000
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == numSweeps {
				b.StopTimer()
				seed()
				idx = 0
				b.StartTimer()
			}
			takerID++
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: takerID, TraderID: 2, Side: book.Buy, Type: book.Market, Qty: sweepLevels * 100}, evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("BestBidAsk", func(b *testing.B) {
		bk, err := book.NewFastBook(cfg, 1024)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k), Qty: 100}, evBuf[:0])
		}
		var bp, bq, ap, aq int64
		var bok, aok bool
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bp, bq, bok = bk.BestBid()
			ap, aq, aok = bk.BestAsk()
		}
		sinkInt64 = bp + bq + ap + aq
		sinkBool = bok && aok
	})

	b.Run("Depth", func(b *testing.B) {
		bk, err := book.NewFastBook(cfg, 1024)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k), Qty: 100}, evBuf[:0])
		}
		bidsBuf := make([]book.DepthLevel, 0, 16)
		asksBuf := make([]book.DepthLevel, 0, 16)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bidsBuf, asksBuf = bk.Depth(10, bidsBuf, asksBuf)
		}
		sinkInt = len(bidsBuf) + len(asksBuf)
	})

	b.Run("MixedWorkload", func(b *testing.B) {
		w := getSharedWorkload()
		bk, err := book.NewFastBook(w.Config, len(w.Warmup)+32768)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 64)
		ops := w.Ops
		nOps := len(ops)
		applyWarmup := func() {
			for i := range w.Warmup {
				evBuf = bk.PlaceOrder(w.Warmup[i].Req, evBuf[:0])
			}
		}
		applyWarmup()
		for i := 0; i < nOps; i++ {
			if ops[i].Kind == bench.OpCancel {
				evBuf = bk.CancelOrder(ops[i].CancelID, evBuf[:0])
			} else {
				evBuf = bk.PlaceOrder(ops[i].Req, evBuf[:0])
			}
		}
		bk.Reset()
		applyWarmup()

		idx := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == nOps {
				b.StopTimer()
				bk.Reset()
				applyWarmup()
				idx = 0
				b.StartTimer()
			}
			if ops[idx].Kind == bench.OpCancel {
				evBuf = bk.CancelOrder(ops[idx].CancelID, evBuf[:0])
			} else {
				evBuf = bk.PlaceOrder(ops[idx].Req, evBuf[:0])
			}
			idx++
		}
		sinkInt = len(evBuf)
	})
}

func BenchmarkArenaBook(b *testing.B) {
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	const batchSize = 4096

	b.Run("AddLimit", func(b *testing.B) {
		rng := rand.New(rand.NewSource(101))
		reqs := make([]book.OrderRequest, batchSize)
		for i := 0; i < batchSize; i++ {
			side := book.Buy
			offset := int64(1 + rng.Intn(200))
			price := mid - offset
			if i%2 == 1 {
				side = book.Sell
				price = mid + offset
			}
			reqs[i] = book.OrderRequest{
				ID: book.OrderID(10_000 + i), TraderID: book.TraderID(1 + i%16),
				Side: side, Type: book.Limit, Price: price, Qty: 100,
			}
		}
		bk, err := book.NewArenaBook(cfg, batchSize+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k*2), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k*2), Qty: 100}, evBuf[:0])
		}
		for i := 0; i < batchSize; i++ {
			evBuf = bk.PlaceOrder(reqs[i], evBuf[:0])
		}
		for i := 0; i < batchSize; i++ {
			evBuf = bk.CancelOrder(reqs[i].ID, evBuf[:0])
		}
		idx := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == batchSize {
				b.StopTimer()
				for j := 0; j < batchSize; j++ {
					evBuf = bk.CancelOrder(reqs[j].ID, evBuf[:0])
				}
				idx = 0
				b.StartTimer()
			}
			evBuf = bk.PlaceOrder(reqs[idx], evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("Cancel", func(b *testing.B) {
		rng := rand.New(rand.NewSource(202))
		reqs := make([]book.OrderRequest, batchSize)
		for i := 0; i < batchSize; i++ {
			side := book.Buy
			offset := int64(1 + rng.Intn(200))
			price := mid - offset
			if i%2 == 1 {
				side = book.Sell
				price = mid + offset
			}
			reqs[i] = book.OrderRequest{
				ID: book.OrderID(10_000 + i), TraderID: book.TraderID(1 + i%16),
				Side: side, Type: book.Limit, Price: price, Qty: 100,
			}
		}
		cancelOrder := make([]book.OrderID, batchSize)
		for i, p := range rng.Perm(batchSize) {
			cancelOrder[i] = reqs[p].ID
		}
		bk, err := book.NewArenaBook(cfg, batchSize+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k*2), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k*2), Qty: 100}, evBuf[:0])
		}
		for i := 0; i < batchSize; i++ {
			evBuf = bk.PlaceOrder(reqs[i], evBuf[:0])
		}
		idx := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == batchSize {
				b.StopTimer()
				for j := 0; j < batchSize; j++ {
					evBuf = bk.PlaceOrder(reqs[j], evBuf[:0])
				}
				idx = 0
				b.StartTimer()
			}
			evBuf = bk.CancelOrder(cancelOrder[idx], evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("MatchSingle", func(b *testing.B) {
		bk, err := book.NewArenaBook(cfg, batchSize+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		seed := func() {
			for j := 0; j < batchSize; j++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1 + j), TraderID: 1, Side: book.Sell, Type: book.Limit, Price: mid, Qty: 100}, evBuf[:0])
			}
		}
		seed()
		idx := 0
		var takerID book.OrderID = 100_000
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == batchSize {
				b.StopTimer()
				seed()
				idx = 0
				b.StartTimer()
			}
			takerID++
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: takerID, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: mid, Qty: 100}, evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("MatchSweep", func(b *testing.B) {
		const sweepLevels, numSweeps = 10, 500
		const totalAsks = sweepLevels * numSweeps
		bk, err := book.NewArenaBook(cfg, totalAsks+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 64)
		seed := func() {
			for j := 0; j < totalAsks; j++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1 + j), TraderID: 1, Side: book.Sell, Type: book.Limit, Price: mid + int64(j), Qty: 100}, evBuf[:0])
			}
		}
		seed()
		idx := 0
		var takerID book.OrderID = 1_000_000
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == numSweeps {
				b.StopTimer()
				seed()
				idx = 0
				b.StartTimer()
			}
			takerID++
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: takerID, TraderID: 2, Side: book.Buy, Type: book.Limit, Price: mid + int64(totalAsks), Qty: sweepLevels * 100}, evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("Market", func(b *testing.B) {
		const sweepLevels, numSweeps = 10, 500
		const totalAsks = sweepLevels * numSweeps
		bk, err := book.NewArenaBook(cfg, totalAsks+512)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 64)
		seed := func() {
			for j := 0; j < totalAsks; j++ {
				evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1 + j), TraderID: 1, Side: book.Sell, Type: book.Limit, Price: mid + int64(j), Qty: 100}, evBuf[:0])
			}
		}
		seed()
		idx := 0
		var takerID book.OrderID = 1_000_000
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == numSweeps {
				b.StopTimer()
				seed()
				idx = 0
				b.StartTimer()
			}
			takerID++
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: takerID, TraderID: 2, Side: book.Buy, Type: book.Market, Qty: sweepLevels * 100}, evBuf[:0])
			idx++
		}
		sinkInt = len(evBuf)
	})

	b.Run("BestBidAsk", func(b *testing.B) {
		bk, err := book.NewArenaBook(cfg, 1024)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k), Qty: 100}, evBuf[:0])
		}
		var bp, bq, ap, aq int64
		var bok, aok bool
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bp, bq, bok = bk.BestBid()
			ap, aq, aok = bk.BestAsk()
		}
		sinkInt64 = bp + bq + ap + aq
		sinkBool = bok && aok
	})

	b.Run("Depth", func(b *testing.B) {
		bk, err := book.NewArenaBook(cfg, 1024)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 16)
		for k := 1; k <= 100; k++ {
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(k), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(k), Qty: 100}, evBuf[:0])
			evBuf = bk.PlaceOrder(book.OrderRequest{ID: book.OrderID(1000 + k), TraderID: 2, Side: book.Sell, Type: book.Limit, Price: mid + int64(k), Qty: 100}, evBuf[:0])
		}
		bidsBuf := make([]book.DepthLevel, 0, 16)
		asksBuf := make([]book.DepthLevel, 0, 16)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bidsBuf, asksBuf = bk.Depth(10, bidsBuf, asksBuf)
		}
		sinkInt = len(bidsBuf) + len(asksBuf)
	})

	b.Run("MixedWorkload", func(b *testing.B) {
		w := getSharedWorkload()
		bk, err := book.NewArenaBook(w.Config, len(w.Warmup)+32768)
		if err != nil {
			b.Fatal(err)
		}
		evBuf := make([]book.Event, 0, 64)
		ops := w.Ops
		nOps := len(ops)
		applyWarmup := func() {
			for i := range w.Warmup {
				evBuf = bk.PlaceOrder(w.Warmup[i].Req, evBuf[:0])
			}
		}
		applyWarmup()
		for i := 0; i < nOps; i++ {
			if ops[i].Kind == bench.OpCancel {
				evBuf = bk.CancelOrder(ops[i].CancelID, evBuf[:0])
			} else {
				evBuf = bk.PlaceOrder(ops[i].Req, evBuf[:0])
			}
		}
		bk.Reset()
		applyWarmup()

		idx := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if idx == nOps {
				b.StopTimer()
				bk.Reset()
				applyWarmup()
				idx = 0
				b.StartTimer()
			}
			if ops[idx].Kind == bench.OpCancel {
				evBuf = bk.CancelOrder(ops[idx].CancelID, evBuf[:0])
			} else {
				evBuf = bk.PlaceOrder(ops[idx].Req, evBuf[:0])
			}
			idx++
		}
		sinkInt = len(evBuf)
	})
}

func TestFastAndArenaLatencyAndGC(t *testing.T) {
	w := bench.GenerateMixedWorkload(bench.DefaultConfig(), 20261122, 2000, 200_000)

	measure := func(name string, place func(book.OrderRequest, []book.Event) []book.Event, cancel func(book.OrderID, []book.Event) []book.Event, reset func()) {
		evBuf := make([]book.Event, 0, 64)
		for i := range w.Warmup {
			evBuf = place(w.Warmup[i].Req, evBuf[:0])
		}
		t0 := time.Now()
		for i := range w.Ops {
			op := &w.Ops[i]
			if op.Kind == bench.OpCancel {
				evBuf = cancel(op.CancelID, evBuf[:0])
			} else {
				evBuf = place(op.Req, evBuf[:0])
			}
		}
		elapsed := time.Since(t0)
		reset()

		for i := range w.Warmup {
			evBuf = place(w.Warmup[i].Req, evBuf[:0])
		}
		h := hdrhistogram.New(1, 50_000_000, 3)
		for i := range w.Ops {
			op := &w.Ops[i]
			s := time.Now()
			if op.Kind == bench.OpCancel {
				evBuf = cancel(op.CancelID, evBuf[:0])
			} else {
				evBuf = place(op.Req, evBuf[:0])
			}
			ns := time.Since(s).Nanoseconds()
			if ns < 1 {
				ns = 1
			}
			_ = h.RecordValue(ns)
		}
		t.Logf("%-12s throughput=%.0f ops/s avg=%.1fns p50=%dns p90=%dns p99=%dns p99.9=%dns max=%dns",
			name, float64(len(w.Ops))/elapsed.Seconds(),
			float64(elapsed.Nanoseconds())/float64(len(w.Ops)),
			h.ValueAtQuantile(50), h.ValueAtQuantile(90), h.ValueAtQuantile(99), h.ValueAtQuantile(99.9), h.Max())
	}

	fb, err := book.NewFastBook(w.Config, 36000)
	if err != nil {
		t.Fatal(err)
	}
	measure("FastBook", fb.PlaceOrder, fb.CancelOrder, fb.Reset)

	ab, err := book.NewArenaBook(w.Config, 36000)
	if err != nil {
		t.Fatal(err)
	}
	measure("ArenaBook", ab.PlaceOrder, ab.CancelOrder, ab.Reset)

	const restingN = 250_000
	cfg := bench.DefaultConfig()
	mid := (cfg.MinPrice + cfg.MaxPrice) / 2
	evBuf := make([]book.Event, 0, 16)

	fbGC, _ := book.NewFastBook(cfg, restingN)
	for i := 1; i <= restingN; i++ {
		evBuf = fbGC.PlaceOrder(book.OrderRequest{ID: book.OrderID(i), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(1+i%500), Qty: 100}, evBuf[:0])
	}
	runtime.GC()
	t1 := time.Now()
	for r := 0; r < 5; r++ {
		runtime.GC()
	}
	fastGCTime := time.Since(t1) / 5
	runtime.KeepAlive(fbGC)
	fbGC = nil
	runtime.GC()

	abGC, _ := book.NewArenaBook(cfg, restingN)
	for i := 1; i <= restingN; i++ {
		evBuf = abGC.PlaceOrder(book.OrderRequest{ID: book.OrderID(i), TraderID: 1, Side: book.Buy, Type: book.Limit, Price: mid - int64(1+i%500), Qty: 100}, evBuf[:0])
	}
	runtime.GC()
	t2 := time.Now()
	for r := 0; r < 5; r++ {
		runtime.GC()
	}
	arenaGCTime := time.Since(t2) / 5
	runtime.KeepAlive(abGC)

	t.Logf("GC cycle with %d resting orders: FastBook(pointers)=%v vs ArenaBook(noscan)=%v (%.2fx faster)",
		restingN, fastGCTime, arenaGCTime, float64(fastGCTime)/float64(arenaGCTime))
}

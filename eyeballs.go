package happy_eyeballs

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const defaultTimeDelta = time.Millisecond * 250

type Dialer interface {
	Dial(ctx context.Context, address string) error
}

type NewDialer func() Dialer

type Logger func(data string)

func nilLogger(data string) {}

type ConnectSpec struct {
	TimeDelta time.Duration
	Addresses []string
	NewDialer NewDialer
	LogInfo   Logger
	LogErr    Logger
}

func (c ConnectSpec) WithDefaults() ConnectSpec {
	if c.TimeDelta == 0 {
		c.TimeDelta = defaultTimeDelta
	}
	if c.LogInfo == nil {
		c.LogInfo = nilLogger
	}
	if c.LogErr == nil {
		c.LogErr = nilLogger
	}
	return c
}

func Dial(cfg ConnectSpec) (Dialer, error) {
	ctx := context.Background()
	return DialWithContext(ctx, cfg)
}

func DialWithContext(ctx context.Context, cfg ConnectSpec) (d Dialer, err error) {
	cfg = cfg.WithDefaults()
	if len(cfg.Addresses) == 0 {
		return nil, fmt.Errorf(`Nothing to connect, addresses list it empty`)
	}

	wg := sync.WaitGroup{}
	ctxExternal, cancelExternal := context.WithCancel(ctx)
	defer func() {
		cancelExternal()
		if d != nil {
			err = nil
		}
	}()

	for _, addr := range cfg.Addresses {
		select {
		case <-ctxExternal.Done():
			return d, ctxExternal.Err()
		default:
		}

		wg.Add(1)
		ctxInternal, cancelInternal := context.WithTimeout(ctx, cfg.TimeDelta)
		go func(addr string) {
			defer func() {
				wg.Done()
				cancelInternal()
			}()
			dialer := cfg.NewDialer()
			cfg.LogInfo(fmt.Sprintf(`Dialing "%s": ...`, addr))
			if err = dialer.Dial(ctx, addr); err != nil {
				cfg.LogErr(fmt.Sprintf(`Dial "%s" failed: %v`, addr, err))
				return
			}
			cfg.LogInfo(fmt.Sprintf(`Dial "%s" successful: returning`, addr))
			d = dialer
			cancelExternal()
		}(addr)

		<-ctxInternal.Done()
		cancelInternal()
	}

	wg.Wait()
	return
}

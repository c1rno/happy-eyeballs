package happy_eyeballs

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestDial_1(t *testing.T) {
	testDial(t, nil, 100, 200, false)
}
func TestDial_2(t *testing.T) {
	testDial(t, nil, 200, 300, false)
}
func TestDial_3(t *testing.T) {
	testDial(t, nil, 300, 400, false)
}
func TestDial_4(t *testing.T) {
	for i := 0; i < 5; i += 1 {
		testDial(t, nil, 100, 400, false)
	}
}
func TestDial_5(t *testing.T) {
	testDial(t, nil, 999, 1000, false)
}
func TestDial_6(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	testDial(t, ctx, 1999, 2000, true)
	cancel()
}

type fakeDialer struct {
	t     *testing.T
	delay time.Duration
}

func (f *fakeDialer) Read(b []byte) (int, error) {
	f.t.Logf("Reading: %s", string(b))
	return 0, nil
}

func (f *fakeDialer) Write(b []byte) (int, error) {
	f.t.Logf("Writing: %s", string(b))
	return 0, nil
}

func (f *fakeDialer) Dial(ctx context.Context, address string) (err error) {
	time.Sleep(f.delay)
	if address == "localhost" {
		return nil
	}
	return errors.New("fakeDialer returnErr=true")
}

func testDial(t *testing.T, ctx context.Context, min, max int64, expectErr bool) {
	defer t.Log("Done")

	cfg := ConnectSpec{
		Addresses: []string{
			"google.com",
			"facebook.com",
			"reddit.com",
			"localhost",
			"linkedin.com",
			"twitter",
		},
		NewDialer: func() Dialer {
			return &fakeDialer{
				t:     t,
				delay: time.Millisecond * time.Duration(rand.Int63n(max-min)+min),
			}
		},
		LogInfo: func(data string) {
			t.Logf("[%s] %s", time.Now(), data)
		},
		LogErr: func(data string) {
			t.Logf("[%s] %s", time.Now(), data)
		},
	}

	var (
		d   Dialer
		err error
	)
	if ctx != nil {
		d, err = DialWithContext(ctx, cfg)
	} else {
		d, err = Dial(cfg)
	}

	if expectErr {
		if err == nil {
			t.Fatalf("Expected err is nil")
		}
		t.Logf("(1) Expected err: %v", err)
		return
	}

	if err != nil {
		t.Fatalf("(1) Unexpected err: %v", err)
	}
	if d == nil {
		t.Fatalf("Returned dialer is nil")
	}

	fake := d.(*fakeDialer)
	_, err = fake.Write([]byte("ping"))
	if err != nil {
		t.Fatalf("(2) Unexpected err: %v", err)
	}
	_, err = fake.Write([]byte("pong"))
	if err != nil {
		t.Fatalf("(3) Unexpected err: %v", err)
	}
}

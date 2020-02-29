package happy_eyeballs

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestDial_1(t *testing.T) {
	testDial(t, 100, 200)
}
func TestDial_2(t *testing.T) {
	testDial(t, 200, 300)
}
func TestDial_3(t *testing.T) {
	testDial(t, 300, 400)
}
func TestDial_4(t *testing.T) {
	for i := 0; i < 5; i += 1 {
		testDial(t, 100, 400)
	}
}
func TestDial_5(t *testing.T) {
	testDial(t, 999, 1000)
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

func testDial(t *testing.T, min, max int64) {
	d, err := Dial(ConnectSpec{
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
	})
	if d == nil {
		t.Fatalf("Returned dialer is nil")
	}
	if err != nil {
		t.Fatalf("(1) Unexpected err: %v", err)
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
	t.Log("Done")
}

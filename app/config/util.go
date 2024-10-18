package config

import (
	"time"

	"github.com/filecoin-project/go-address"
)

// Time is a wrapper type for Duration
// for decoding and encoding from/to TOML
type Time time.Time

// UnmarshalText implements interface for TOML decoding
func (t *Time) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*t = Time(time.Time{})
		return nil
	}
	d, err := time.Parse(string(text), time.DateOnly)
	if err != nil {
		return err
	}
	*t = Time(d)
	return err
}

func (t Time) MarshalText() ([]byte, error) {
	if t.Time().IsZero() {
		return []byte{}, nil
	}
	d := t.Time().Format(time.DateOnly)
	return []byte(d), nil
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

// Address is a wrapper type for Address
// for decoding and encoding from/to TOML
type Address address.Address

// UnmarshalText implements interface for TOML decoding
func (addr *Address) UnmarshalText(text []byte) error {
	d, err := address.NewFromString(string(text))
	if err != nil {
		return err
	}
	*addr = Address(d)
	return err
}

func (addr Address) MarshalText() ([]byte, error) {
	if address.Address(addr) == address.Undef {
		return []byte{}, nil
	}
	return []byte(address.Address(addr).String()), nil
}

func (addr Address) String() string {
	return address.Address(addr).String()
}

func (addr Address) Address() address.Address {
	return address.Address(addr)
}

func (addr Address) Empty() bool {
	return address.Address(addr).Empty()
}

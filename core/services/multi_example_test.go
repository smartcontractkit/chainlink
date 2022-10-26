package services

import (
	"context"
	"fmt"
)

type Healthy string

func (h Healthy) Start(ctx context.Context) error {
	fmt.Println(h, "started")
	return nil
}

func (h Healthy) Close() error {
	fmt.Println(h, "closed")
	return nil
}

type CloseFailure string

func (c CloseFailure) Start(ctx context.Context) error {
	fmt.Println(c, "started")
	return nil
}

func (c CloseFailure) Close() error {
	fmt.Println(c, "close failure")
	return fmt.Errorf("failed to close: %s", c)
}

type WontStart string

func (f WontStart) Start(ctx context.Context) error {
	fmt.Println(f, "start failure")
	return fmt.Errorf("failed to start: %s", f)
}

func (f WontStart) Close() error {
	fmt.Println(f, "close failure")
	return fmt.Errorf("cannot call Close after failed Start: %s", f)
}

func ExampleMultiStart() {
	ctx := context.Background()

	a := Healthy("a")
	b := CloseFailure("b")
	c := WontStart("c")

	var ms MultiStart
	if err := ms.Start(ctx, a, b, c); err != nil {
		fmt.Println(err)
	}

	// Output:
	// a started
	// b started
	// c start failure
	// b close failure
	// a closed
	// failed to start: c; failed to close: b
}

func ExampleMultiClose() {
	ctx := context.Background()

	f1 := CloseFailure("f")
	f2 := CloseFailure("f")

	var ms MultiStart
	if err := ms.Start(ctx, f1, f2); err != nil {
		fmt.Println(err)
		return
	}
	mc := MultiClose{f1, f2}
	if err := mc.Close(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// f started
	// f started
	// f close failure
	// f close failure
	// failed to close: f; failed to close: f
}

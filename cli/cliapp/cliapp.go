package cliapp

import (
	"context"
	"fmt"
	"time"

	"carmaintenance/internal/core"
)

// func Run(ctx context.Context, dbPath *string, specsDir *string) error {
func Run(ctx context.Context, backend *core.Backend) error {

	for {
		fmt.Println("working...")
		time.Sleep(1 * time.Second)
	}

	// // err = go func(ctx context.Context) error {
	// // func(ctx context.Context) error {
	// func(ctx context.Context) {
	// 	for {
	// 		// // TODO: Main app loop
	// 		// fmt.Println("working...")
	// 		// time.Sleep(1 * time.Second)
	//
	// 		select {
	// 		case <-ctx.Done():
	// 			fmt.Println("context cancelled, stopping app logic")
	// 			// return nil
	// 			return
	// 		// default:
	// 		// 	// do real work
	// 		// 	// TODO: Run the app logic here
	// 		// 	// if err := doWork(db); err != nil {
	// 		// 	// 	return err
	// 		// 	// }
	// 		// }
	// 		default:
	// 			fmt.Println("working...")
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	}
	// 	// return nil
	// }(ctx)
	//
	// fmt.Printf("algo")
	// <-ctx.Done() // Wait for interrupt signal
	// fmt.Printf("mais")

	return nil
}


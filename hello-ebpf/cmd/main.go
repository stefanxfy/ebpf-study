package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type packet_event PacketMonitor ../ebpf/packet_monitor.c -- -I/usr/include/x86_64-linux-gnu

// PacketEvent represents the structure sent from eBPF program
type PacketEvent struct {
	SrcIP    uint32
	DstIP    uint32
	SrcPort  uint16
	DstPort  uint16
	Protocol uint8
	Length   uint32
}

func main() {
	// Remove memory limit for eBPF
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Failed to remove memlock:", err)
	}
	_ = os.RemoveAll("/sys/fs/bpf/hello")
	bpfPath := "./target/hello.o"

	bpfSpec, err := ebpf.LoadCollectionSpec(bpfPath)
	if err != nil {
		var verifierError *ebpf.VerifierError
		if errors.As(err, &verifierError) {
			log.Printf("Verifier error: %+v\n", verifierError)
		}
		log.Printf("Failed to load eBPF spec: %v\n", err)
		os.Exit(1)
	}

	// Load the eBPF collection
	coll, err := ebpf.NewCollection(bpfSpec)
	if err != nil {
		var verifierError *ebpf.VerifierError
		if errors.As(err, &verifierError) {
			log.Printf("Verifier error: %+v\n", verifierError)
		}
		log.Printf("Failed to load eBPF collection: %v\n", err)
		os.Exit(1)
	}
	defer coll.Close()

	// Get the tracepoint program
	prog := coll.Programs["handle_tp"]
	if prog == nil {
		log.Fatal("Program 'handle_tp' not found in eBPF collection")
	}

	// Attach to the sys_enter_write tracepoint
	l, err := link.Tracepoint("syscalls", "sys_enter_write", prog, nil)
	if err != nil {
		log.Printf("Failed to attach to tracepoint: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	log.Println("eBPF program attached successfully to sys_enter_write tracepoint")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nReceived interrupt, shutting down...")
		cancel()
	}()

	// Wait for context cancellation
	<-ctx.Done()

	println(len(bpfSpec.Programs))
	fmt.Println("Monitoring stopped")
}

package communication

import (
	"fmt"
//========= REDACTED =========
	"net"

	"github.com/dedis/protobuf"
	"github.com/2_alt_hw2/Peerster/logger"
	"github.com/2_alt_hw2/Peerster/types"

//========= REDACTED =========
//========= REDACTED =========
)

// UDPWriter represents a connexion through UDP which encapsulates a socket
// This abstraction prevents reading from multiple locations in the code
type UDPWriter struct {
	conn *net.UDPConn
}

// UDPDecoder is a callback function that handles incoming packets
// It can returns false to stop the listening loop
type UDPDecoder func(string, []byte) bool

// NewUDPListener creates a goroutine that listens for incoming UDP datagrams and forwards them
// to the decoder given in parameter
func NewUDPListener(address string, decoder UDPDecoder, log *logger.Logger) (*UDPWriter, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		running := true
		for running {
			buffer := make([]byte, 16384)
			nBytes, gossipAddr, err := udpConn.ReadFromUDP(buffer[:])
			if err == nil {
				running = decoder(gossipAddr.String(), buffer[0:nBytes])
			} else {
				log.Warn(fmt.Sprintf("Unable to read from socket (%v): %v", address, err))
			}
		}

		err := udpConn.Close()
		if err != nil {
			log.Warn(fmt.Sprintf("Unable to to close %v: %v", udpConn, err))
		}
	}()

	return &UDPWriter{udpConn}, nil
}

// WriteTo sends a packet to an address
func (u UDPWriter) WriteTo(packet []byte, addr *net.UDPAddr) (int, error) {
	return u.conn.WriteToUDP(packet, addr)
}

// NewClientPacketDecoder is a preconfigured packet decoder for ClientPacket type
func NewClientPacketDecoder(log *logger.Logger) (<-chan types.Message, UDPDecoder) {
	clientCh := make(chan types.Message, 100)
	callback := func(_ string, buffer []byte) bool {
		packet := types.Message{}
		err := protobuf.Decode(buffer, &packet)
		if err != nil {
			log.Warn(fmt.Sprintf("Error in ClientPacket decoding: %v", err))
			return clientCh != nil // Keep going if channel is still open
		}

		if clientCh != nil {
			clientCh <- packet
			return true // Keep going
		}
		return false // Stop reading packets
	}
	return clientCh, callback
}

// NewGossipPacketDecoder is a preconfigured packet decoder for GossipPacket type
func NewGossipPacketDecoder(log *logger.Logger) (<-chan types.GossipPacket, UDPDecoder) {
	gossipCh := make(chan types.GossipPacket, 100)
	callback := func(gossipAddr string, buffer []byte) bool {
		packet := types.GossipPacket{}
		err := protobuf.Decode(buffer, &packet)
		if err != nil {
			log.Warn(fmt.Sprintf("Error in GossipPacket decoding: %v", err))
			return gossipCh != nil // Keep going if channel is still open
		}

		if gossipCh != nil {
			packet.GossipAddr = &gossipAddr
			gossipCh <- packet
			return true // Keep going
		}
		return false // Stop reading packets
	}
	return gossipCh, callback
}

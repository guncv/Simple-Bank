package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) (*Metadata, error) {
	metaData := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %v", md)
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			metaData.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			metaData.UserAgent = userAgent[0]
		}

		if clientIp := md.Get(xForwardedForHeader); len(clientIp) > 0 {
			metaData.ClientIp = clientIp[0]
		}
	}

	peer, ok := peer.FromContext(ctx)
	if ok {
		if peer.Addr != nil {
			log.Printf("peer: %v", peer)
			metaData.ClientIp = peer.Addr.String()
		}
	}

	log.Printf("metaData: %v", metaData)

	return metaData, nil
}

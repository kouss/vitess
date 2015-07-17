// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorpcbinlogplayer

import (
	"time"

	"golang.org/x/net/context"

	"github.com/youtube/vitess/go/netutil"
	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcwrap/bsonrpc"
	"github.com/youtube/vitess/go/vt/binlog/binlogplayer"
	"github.com/youtube/vitess/go/vt/binlog/proto"
	"github.com/youtube/vitess/go/vt/topo"
)

// client implements a Client over go rpc
type client struct {
	*rpcplus.Client
}

func (client *client) Dial(endPoint topo.EndPoint, connTimeout time.Duration) error {
	addr := netutil.JoinHostPort(endPoint.Host, endPoint.NamedPortMap["vt"])
	var err error
	client.Client, err = bsonrpc.DialHTTP("tcp", addr, connTimeout, nil)
	return err
}

func (client *client) Close() {
	client.Client.Close()
}

func (client *client) ServeUpdateStream(ctx context.Context, req *proto.UpdateStreamRequest) (chan *proto.StreamEvent, binlogplayer.ErrFunc, error) {
	result := make(chan *proto.StreamEvent, 10)
	responseChan := make(chan *proto.StreamEvent, 10)
	resp := client.Client.StreamGo("UpdateStream.ServeUpdateStream", req, responseChan)
	var finalError error
	go func() {
		defer close(result)
		for {
			select {
			case <-ctx.Done():
				finalError = ctx.Err()
				return
			case r, ok := <-responseChan:
				if !ok {
					// no more results from the server
					finalError = resp.Error
					return
				}
				result <- r
			}
		}
	}()
	return result, func() error {
		return finalError
	}, nil
}

func (client *client) StreamKeyRange(ctx context.Context, req *proto.KeyRangeRequest) (chan *proto.BinlogTransaction, binlogplayer.ErrFunc, error) {
	result := make(chan *proto.BinlogTransaction, 10)
	responseChan := make(chan *proto.BinlogTransaction, 10)
	resp := client.Client.StreamGo("UpdateStream.StreamKeyRange", req, responseChan)
	var finalError error
	go func() {
		defer close(result)
		for {
			select {
			case <-ctx.Done():
				finalError = ctx.Err()
				return
			case r, ok := <-responseChan:
				if !ok {
					// no more results from the server
					finalError = resp.Error
					return
				}
				result <- r
			}
		}
	}()
	return result, func() error {
		return finalError
	}, nil
}

func (client *client) StreamTables(ctx context.Context, req *proto.TablesRequest) (chan *proto.BinlogTransaction, binlogplayer.ErrFunc, error) {
	result := make(chan *proto.BinlogTransaction, 10)
	responseChan := make(chan *proto.BinlogTransaction, 10)
	resp := client.Client.StreamGo("UpdateStream.StreamTables", req, responseChan)
	var finalError error
	go func() {
		defer close(result)
		for {
			select {
			case <-ctx.Done():
				finalError = ctx.Err()
				return
			case r, ok := <-responseChan:
				if !ok {
					// no more results from the server
					finalError = resp.Error
					return
				}
				result <- r
			}
		}
	}()
	return result, func() error {
		return finalError
	}, nil
}

// Registration as a factory
func init() {
	binlogplayer.RegisterClientFactory("gorpc", func() binlogplayer.Client {
		return &client{}
	})
}
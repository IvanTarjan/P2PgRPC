package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/IvanTarjan/P2PgRPC/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string
	listenAddr string
	logger *zap.SugaredLogger
	peerLock sync.RWMutex
	peers map[proto.NodeClient]*proto.Version
	proto.UnimplementedNodeServer
}

func NewNode() *Node{
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	return &Node{
		peers: make(map[proto.NodeClient]*proto.Version),
		version: "blocker-0.1",
		logger: logger.Sugar(),
	}
}


func (n *Node) Start(listenAddr string, bootstrapNodes []string) error{
	n.listenAddr = listenAddr
	var(
		opts = []grpc.ServerOption{}
		grpServer = grpc.NewServer(opts...)	
	)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil{
		return err
	}
	proto.RegisterNodeServer(grpServer, n)
	n.logger.Infow("node started...", "port", n.listenAddr)
	
	// Bootstrap network with list of already known nodes
	if len(bootstrapNodes) > 0{
		go n.bootstrapNetwork(bootstrapNodes)
	}

	return grpServer.Serve(ln)
}


func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error){
	peer, _ := peer.FromContext(ctx)
	fmt.Println("Received transaction from: ", peer)
	return &proto.Ack{}, nil
}

func (n *Node) Handshake(ctx context.Context, version *proto.Version) (*proto.Version, error){
	c, err := makeNodeClient(version.ListenAddr)
	if err != nil {
		return nil , err
	}
	n.addPeer(c, version)
	return n.getVersion(), nil
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version){
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	// Handle the logic where we decide if we should add the peer or not
	
	n.peers[c] = v

	// connect all peers from received peerlist
	if len(v.PeerList) > 0 {
		go n.bootstrapNetwork(v.PeerList)
	}

	n.logger.Debugw("new peer successfully connected",
	"we", n.listenAddr,
	"addr", v.ListenAddr,
	"height", v.Height)
}

func (n *Node) deletePeer(c proto.NodeClient){
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) bootstrapNetwork(addrs []string) error{
	for _, addr := range addrs{
		if !n.canConnectWith(addr){
			continue
		}
		n.logger.Debugw("dialing remote node", "we", n.listenAddr, "remote", addr)
		c, v, err := n.dialRemoteNode(addr)
		if err != nil {
			return err
		}
		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) dialRemoteNode(addr string) (proto.NodeClient, *proto.Version, error){
	c, err := makeNodeClient(addr)
		if err != nil{
			return  nil, nil, err
		}
		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Error("Handshake error: ", err)
			return nil, nil, err
		}
		return c, v, nil
}

func (n *Node) getVersion() *proto.Version{
	return &proto.Version{
		Version: "blocker-0.1",
		Height: 0,
		ListenAddr: n.listenAddr,
		PeerList: n.getPeerList(),
	}
}

func (n *Node) canConnectWith(addr string) bool{
	if n.listenAddr == addr{
		return false
	}
	connectedPeers := n.getPeerList()
	for _, connectdAddr := range connectedPeers{
		if connectdAddr == addr{
			return false
		}
	}
	return true
}

func (n *Node) getPeerList() []string{
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()
	peers := []string{}
	for _, v := range n.peers{
		peers = append(peers, v.ListenAddr)
	}
	return peers
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error){
	c, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(c), nil
}
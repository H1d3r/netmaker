package peer

import (
	"fmt"
	"log"
	"net"

	"github.com/gravitl/netmaker/nm-proxy/common"
	"github.com/gravitl/netmaker/nm-proxy/proxy"
	"github.com/gravitl/netmaker/nm-proxy/wg"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Conn struct {
	Config ConnConfig
	Proxy  proxy.Proxy
}

// ConnConfig is a peer Connection configuration
type ConnConfig struct {

	// Key is a public key of a remote peer
	Key string
	// LocalKey is a public key of a local peer
	LocalKey string

	ProxyConfig     proxy.Config
	AllowedIPs      string
	LocalWgPort     int
	RemoteProxyIP   net.IP
	RemoteWgPort    int
	RemoteProxyPort int
}

func AddNewPeer(wgInterface *wg.WGIface, peer *wgtypes.PeerConfig) error {

	c := proxy.Config{
		Port:        peer.Endpoint.Port,
		RemoteKey:   peer.PublicKey.String(),
		WgInterface: wgInterface,
		AllowedIps:  peer.AllowedIPs,
	}
	p := proxy.NewProxy(c)
	remoteConn, err := net.Dial("udp", fmt.Sprintf("%s:%d", peer.Endpoint.IP.String(), common.NmProxyPort))
	if err != nil {
		return err
	}
	log.Printf("----> Established Remote Conn with RPeer: %s, LAddr: %s ----> RAddr: %s", peer.PublicKey, remoteConn.LocalAddr().String(), remoteConn.RemoteAddr().String())
	log.Printf("Starting proxy for Peer: %s\n", peer.PublicKey.String())
	err = p.Start(remoteConn)
	if err != nil {
		return err
	}
	connConf := common.ConnConfig{
		Key:             peer.PublicKey.String(),
		LocalKey:        wgInterface.Device.PublicKey.String(),
		LocalWgPort:     wgInterface.Device.ListenPort,
		RemoteProxyIP:   net.ParseIP(peer.Endpoint.IP.String()),
		RemoteWgPort:    peer.Endpoint.Port,
		RemoteProxyPort: common.NmProxyPort,
	}
	peerProxy := common.Proxy{
		Ctx:    p.Ctx,
		Cancel: p.Cancel,
		Config: common.Config{
			Port:        peer.Endpoint.Port,
			RemoteKey:   peer.PublicKey.String(),
			WgInterface: wgInterface,
			AllowedIps:  peer.AllowedIPs,
		},

		RemoteConn: remoteConn,
		LocalConn:  p.LocalConn,
	}
	peerConn := common.Conn{
		Config: connConf,
		Proxy:  peerProxy,
	}
	if _, ok := common.WgIFaceMap[wgInterface.Name]; ok {
		common.WgIFaceMap[wgInterface.Name][peer.PublicKey.String()] = &peerConn
	} else {
		common.WgIFaceMap[wgInterface.Name] = make(map[string]*common.Conn)
		common.WgIFaceMap[wgInterface.Name][peer.PublicKey.String()] = &peerConn
	}
	return nil
}

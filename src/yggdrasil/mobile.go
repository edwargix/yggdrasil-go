// +build mobile

package yggdrasil

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	hjson "github.com/hjson/hjson-go"
	"github.com/mitchellh/mapstructure"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"
	"github.com/yggdrasil-network/yggdrasil-go/src/crypto"
	"github.com/yggdrasil-network/yggdrasil-go/src/util"
)

// This file is meant to "plug the gap" for mobile support, as Gomobile will
// not create headers for Swift/Obj-C etc if they have complex (non-native)
// types. Therefore for iOS we will expose some nice simple functions. Note
// that in the case of iOS we handle reading/writing to/from TUN in Swift
// therefore we use the "dummy" TUN interface instead.

func (c *Core) StartAutoconfigure() error {
	logger := log.New(os.Stdout, "", 0)
	nc := config.GenerateConfig(true)
	nc.IfName = "dummy"
	nc.AdminListen = "tcp://localhost:9001"
	nc.Peers = []string{}
	if hostname, err := os.Hostname(); err == nil {
		nc.NodeInfo = map[string]interface{}{"name": hostname}
	}
	ifceExpr, err := regexp.Compile(".*")
	if err == nil {
		c.ifceExpr = append(c.ifceExpr, ifceExpr)
	}
	if err := c.Start(nc, logger); err != nil {
		return err
	}
	return nil
}

func (c *Core) StartJSON(configjson []byte) error {
	logger := log.New(os.Stdout, "", 0)
	nc := config.GenerateConfig(false)
	var dat map[string]interface{}
	if err := hjson.Unmarshal(configjson, &dat); err != nil {
		return err
	}
	if err := mapstructure.Decode(dat, &nc); err != nil {
		return err
	}
	nc.IfName = "dummy"
	for _, ll := range nc.MulticastInterfaces {
		ifceExpr, err := regexp.Compile(ll)
		if err != nil {
			panic(err)
		}
		c.AddMulticastInterfaceExpr(ifceExpr)
	}
	if err := c.Start(nc, logger); err != nil {
		return err
	}
	return nil
}

func GenerateConfigJSON() []byte {
	nc := config.GenerateConfig(false)
	nc.IfName = "dummy"
	if json, err := json.Marshal(nc); err == nil {
		return json
	} else {
		return nil
	}
}

func (c *Core) GetAddressString() string {
	return c.GetAddress().String()
}

func (c *Core) GetSubnetString() string {
	return c.GetSubnet().String()
}

func (c *Core) RouterRecvPacket() ([]byte, error) {
	packet := <-c.router.tun.recv
	return packet, nil
}

func (c *Core) RouterSendPacket(buf []byte) error {
	packet := append(util.GetBytes(), buf[:]...)
	c.router.tun.send <- packet
	return nil
}

func (c *Core) AWDLCreateInterface(boxPubKey []byte, sigPubKey []byte, name string) {
	var box crypto.BoxPubKey
	var sig crypto.SigPubKey
	copy(box[:crypto.BoxPubKeyLen], boxPubKey[:])
	copy(sig[:crypto.SigPubKeyLen], sigPubKey[:])
	c.awdl.create(&box, &sig, name)
}

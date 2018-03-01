package main

import (
	"fmt"
	"seth/accounts"
	"seth/log"

	cli "gopkg.in/urfave/cli.v1"
)

//NodeCli cli for node
type NodeCli struct {
}

// init nodecli init
func (n *NodeCli) init(c *cli.Context) error {
	return nil
}

// Help help cmd for nodecli
func (n *NodeCli) Help(c *cli.Context) error {
	return nil
}

// NewAccount new account output address of account,publickey&privatekey
func (n *NodeCli) NewAccount(c *cli.Context) error {
	address, publickey, privatekey := accounts.NewAccount()
	log.Info("address:%s;publickey:%s;privatekey:%s", address, publickey, privatekey)
	fmt.Printf("address:%s\n", address)
	fmt.Printf("publickey:%s\n", publickey)
	fmt.Printf("privatekey:%s\n", privatekey)
	return nil
}

// Start start cmd for nodecli
func (n *NodeCli) Start(c *cli.Context) error {
	return nil
}

// Clear Clear cmd for nodecli
func (n *NodeCli) Clear(c *cli.Context) error {
	return nil
}

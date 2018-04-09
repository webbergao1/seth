package main

import (
	"encoding/json"
	"fmt"
	"os"
	"seth/accounts"
	"seth/config"
	"seth/core"
	"seth/database"
	_ "seth/database/leveldb"
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

// InitGenesis init genesis block
func (n *NodeCli) InitGenesis(c *cli.Context) error {
	genesisparam := c.Args().First()
	var genesis *core.Genesis
	switch genesisparam {
	case "":
		genesis = core.DefaultGenesis()
	case core.TagMainNetGenesis:
		genesis = core.DefaultGenesis()
	case core.TagTestNetGenesis:
		genesis = core.TestnetGenesis()
	case core.TagDeveloperNetGenesis:
		genesis = core.DevelopernetGenesis()
	default:
		file, err := os.Open(genesisparam)
		if err != nil {
			log.Fatal("Failed to read genesis file: %v", err)
			return err
		}
		defer file.Close()
		genesis = new(core.Genesis)
		if err := json.NewDecoder(file).Decode(genesis); err != nil {
			log.Fatal("invalid genesis file: %v", err)
			return err
		}
	}
	datapath := config.ResolvePath("chaindata")
	db, err := database.GetDatabase(database.LevelDBName)
	if err != nil {
		log.Fatal("get database error: %v", err)
		return err
	}
	err = db.Open(datapath, 0, 0)
	if err != nil {
		log.Fatal("open database error: %v", err)
		return err
	}
	defer db.Close()
	hash, err := genesis.SetupGensisBlock(db)
	log.Info("genesis block hash:%s", hash.Hex())
	if err == core.ErrHasGenesisBlock {
		log.Error("error: %v;exist genesis block hash:%s", err, hash.Hex())
	}
	return err
}

// Start start cmd for nodecli
func (n *NodeCli) Start(c *cli.Context) error {
	return nil
}

// Clear Clear cmd for nodecli
func (n *NodeCli) Clear(c *cli.Context) error {
	return nil
}

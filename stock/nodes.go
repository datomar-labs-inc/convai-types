package stock

import (
	ctypes "github.com/datomar-labs-inc/convai-types"
)

var StockNodeBranch = ctypes.DBNode{
	Name:    "Branch",
	TypeID:  "branch",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

var StockNodeHalt = ctypes.DBNode{
	Name:    "Halt",
	TypeID:  "halt",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

var StockNodeResume = ctypes.DBNode{
	Name:    "Resume",
	TypeID:  "resume",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

type SetDataConfig struct {
	Items []SetDataItem `json:"items"`
}

type SetDataItem struct {
	Path  string `json:"path"`  // Path is the context level/memory container name path, eg: user.data.likesCheese
	Value string `json:"value"` // Value is what the data should be set to, usually a templated string
}

var StockNodeSetData = ctypes.DBNode{
	Name:    "Set Data",
	TypeID:  "set_data",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

var StockNodeDeleteData = ctypes.DBNode{
	Name:    "Delete Data",
	TypeID:  "delete_data",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

var StockNodeWebhook = ctypes.DBNode{
	Name:    "Webhook",
	TypeID:  "webhook",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

var StockNodeHTTPRequest = ctypes.DBNode{
	Name:    "HTTP Request",
	TypeID:  "http",
	Version: "0.0.1",
	Style: ctypes.NodeStyle{
		Color: "#ffffff",
	},

	// TODO write documentation
	Documentation: ``,
}

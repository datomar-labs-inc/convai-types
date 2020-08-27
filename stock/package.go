package stock

import (
	"github.com/google/uuid"

	ctypes "github.com/datomar-labs-inc/convai-types"
)

var StockPackage = ctypes.Package{
	DBPackage: ctypes.DBPackage{
		ID:             uuid.Nil,
		Name:           "Stock",
		Description:    "Functionality provided by default in Convai",
		OrganizationID: uuid.Nil,
		BaseURL:        "https://stock.convai.dev",
		SigningKey:     "cheezits",
	},
	Nodes: []ctypes.DBNode{
		StockNodeBranch,
		StockNodeSetData,
		StockNodeDeleteData,
		StockNodeHalt,
		StockNodeResume,
		StockNodeHTTPRequest,
		StockNodeWebhook,
	},
	Links: []ctypes.DBLink{
		StockLinkBasic,
		StockLinkDataEquals,
		StockLinkError,
		StockLinkPriority,
	},
	Events:     nil,
	Dispatches: nil,
}

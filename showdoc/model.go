package showdoc

type Catalog struct {
	CatId       int64
	CatName     string
	ItemId      int64
	SNumber     int64
	AddTime     int64
	ParentCatId int64
	Level       int64
}

type CatalogTree struct {
	ID          int64
	Name        string
	Level       int64
	SubCatalogs []CatalogTree
}

type Item struct {
	ItemId          int64
	ItemName        string
	ItemDescription string
	Uid             int64
	UserName        string
	Password        string
	AddTime         int64
	LastUpdateTime  int64
	ItemDomain      string
	ItemType        int64
	IsArchived      bool
	IsDel           bool
}

type Page struct {
	PageId         int64
	AuthorId       int64
	AuthorUserName string
	ItemId         int64
	CatId          int64
	PageTitle      string
	PageContent    string
	SNumber        int64
	AddTime        int64
	PageComments   string
	IsDel          bool
}

type PageHistory struct {
	PageHistoryId  int64
	PageId         int64
	AuthorUid      int64
	AuthorUserName string
	ItemId         int64
	CatId          int64
	PageTitle      string
	PageContent    string
	SNumber        int64
	AddTime        int64
	PageComments   string
}
